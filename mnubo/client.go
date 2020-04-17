package mnubo

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultTimeout = time.Second * 10

	DefaultBackoffMaxInterval = time.Minute * 5
)

// CompressionConfig is used to compress requests and / or response to / from the SmartObjects platform.
type CompressionConfig struct {
	Request  bool
	Response bool
}

// ExponentialBackoffConfig is used to configure exponential backoff.
// You should not need to configure this, but it's available if you require further tweaking.
type ExponentialBackoffConfig struct {
	// After MaxElapsedTime the ExponentialBackOff stops.
	// It never stops if MaxElapsedTime == 0.
	MaxElapsedTime time.Duration
	// Callback called between retries when the platform is not available.
	// It will not be called if value is nil.
	NotifyOnError func(error, time.Duration)
}

// Mnubo is the main object representing all available endpoints in SmartObjects.
type Mnubo struct {
	ClientId           string
	ClientSecret       string
	ClientToken        string
	Host               string
	AccessToken        AccessToken
	Timeout            time.Duration // Timeout for HTTP requests sent to SmartObjects.
	Compression        CompressionConfig
	ExponentialBackoff ExponentialBackoffConfig
	Model              *Model
	Events             *Events
	Objects            *Objects
	Owners             *Owners
	Search             *Search
	CustomTransport    *http.Transport
}

// ClientRequest is an internal structure to help with making HTTP requests to SmartObjects.
type ClientRequest struct {
	authorization   string
	method          string
	path            string
	contentType     string
	urlQuery        url.Values
	payload         []byte
	skipCompression bool
}

// AccessToken represents a token obtained after validating client id and secret.
type AccessToken struct {
	Value     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	ExpiresAt time.Time
	Scope     string `json:"scope"`
	Jti       string `json:"jti"`
}

type smartObjectsNotAvailableError struct {
}

func (e *smartObjectsNotAvailableError) Error() string {
	return "SmartObjects platform is not available"
}

func smartObjectsNotAvailable() *smartObjectsNotAvailableError {
	return &smartObjectsNotAvailableError{}
}

// hasExpired returns true if an access token has expired.
func (at *AccessToken) hasExpired() bool {
	now := time.Now()
	return at.ExpiresAt.Before(now)
}

// NewClient creates a new Mnubo structure based on id, secret and host.
func NewClient(id string, secret string, host string) *Mnubo {
	m := &Mnubo{
		ClientId:     id,
		ClientSecret: secret,
		Host:         host,
	}
	m.initClient()
	return m
}

// NewClientWithToken creates a new Mnubo structure based on a static token.
func NewClientWithToken(token string, host string) *Mnubo {
	m := &Mnubo{
		ClientToken: token,
		Host:        host,
	}
	m.initClient()
	return m
}

// initClient initializes internal wrappers for SmartObjects main endpoints.
func (m *Mnubo) initClient() {
	m.Model = NewModel(m)
	m.Events = NewEvents(m)
	m.Objects = NewObjects(m)
	m.Owners = NewOwners(m)
	m.Search = NewSearch(m)
	m.Timeout = DefaultTimeout
	m.ExponentialBackoff = ExponentialBackoffConfig{
		MaxElapsedTime: DefaultBackoffMaxInterval,
	}
}

// isUsingStaticToken returns true if the client was initialized with its own static token
// ie: not using client id / secret.
func (m *Mnubo) isUsingStaticToken() bool {
	return m.ClientToken != ""
}

// GetAccessToken fetches a new AccessToken with scope ALL.
func (m *Mnubo) GetAccessToken() (AccessToken, error) {
	return m.GetAccessTokenWithScope("ALL")
}

// GetAccessTokenWithScope fetches a new AccessToken with specified scope.
func (m *Mnubo) GetAccessTokenWithScope(scope string) (AccessToken, error) {
	payload := fmt.Sprintf("grant_type=client_credentials&scope=%s", scope)
	data := []byte(fmt.Sprintf("%s:%s", m.ClientId, m.ClientSecret))

	cr := ClientRequest{
		authorization:   fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(data)),
		method:          "POST",
		path:            "/oauth/token",
		contentType:     "application/x-www-form-urlencoded",
		skipCompression: true,
		payload:         []byte(payload),
	}
	at := AccessToken{}
	err := m.doRequest(cr, &at)
	now := time.Now()

	if err == nil {
		if err != nil {
			return at, fmt.Errorf("unable to unmarshall body %t", err)
		}
		dur, err := time.ParseDuration(fmt.Sprintf("%dms", at.ExpiresIn))
		at.ExpiresAt = now.Add(dur)
		m.AccessToken = at
		return at, err
	}
	return at, err
}

// doGzip compressed data using gzip BestSpeed algorithm.
func doGzip(w io.Writer, data []byte) error {
	gw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		return err
	}
	if _, err := gw.Write(data); err != nil {
		return err
	}
	if err := gw.Flush(); err != nil {
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}
	return nil
}

// doGunzip uncompressed gzipped data.
func doGunzip(w io.Writer, data []byte) error {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	defer gr.Close()
	if err != nil {
		return err
	}
	ud, err := ioutil.ReadAll(gr)
	if err != nil {
		return err
	}
	w.Write(ud)
	return nil
}

func doHttpRequest(client *http.Client, req *http.Request, response interface{}) func() error {
	wrappedFunc := func() error {
		res, err := client.Do(req)
		if err != nil {
			return backoff.Permanent(err)
		}
		defer res.Body.Close()

		var body []byte
		body, err = ioutil.ReadAll(res.Body)

		if err != nil {
			return backoff.Permanent(err)
		}

		if res.Header.Get("Content-Encoding") == "gzip" {
			var w bytes.Buffer
			err := doGunzip(&w, body)

			if err != nil {
				return backoff.Permanent(err)
			}

			body = w.Bytes()
		}

		ct := res.Header.Get("Content-Type")
		if strings.Contains(ct, "application/json") {
			err := json.Unmarshal(body, response)
			if err != nil {
				return backoff.Permanent(err)
			}
		} else if strings.Contains(ct, "text/plain") {
			response = string(body)
		}

		if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
			return nil
		} else if res.StatusCode == http.StatusServiceUnavailable {
			return smartObjectsNotAvailable()
		}

		return backoff.Permanent(errors.New(fmt.Sprintf("The server responded with StatusCode: %d - Body: %s", res.StatusCode, response)))
	}

	return wrappedFunc
}

// doRequest is the main internal helper to send request to the SmartObjects platform.
// It handles compression / decompression based on client configuration.
func (m *Mnubo) doRequest(cr ClientRequest, response interface{}) error {
	var payload []byte

	if m.Compression.Request && !cr.skipCompression {
		var w bytes.Buffer
		err := doGzip(&w, cr.payload)
		if err != nil {
			return fmt.Errorf("unable to gzip request: %t", err)
		}
		payload = w.Bytes()
	} else {
		payload = cr.payload
	}

	req, err := http.NewRequest(cr.method, m.Host+cr.path, bytes.NewReader(payload))

	req.Header.Add("Content-Type", cr.contentType)
	req.Header.Add("X-MNUBO-SDK", "Go")

	if cr.authorization != "" {
		req.Header.Add("Authorization", cr.authorization)
	}

	if cr.urlQuery != nil {
		req.URL.RawQuery = cr.urlQuery.Encode()
	}

	if m.Compression.Request {
		req.Header.Add("Content-Encoding", "gzip")
	}

	if m.Compression.Response {
		req.Header.Add("Accept-Encoding", "gzip")
	}

	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: m.Timeout,
	}
	if m.CustomTransport != nil {
		client.Transport = m.CustomTransport
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = m.ExponentialBackoff.MaxElapsedTime
	return backoff.RetryNotify(doHttpRequest(client, req, response), b, m.ExponentialBackoff.NotifyOnError)
}

// doRequestWithAuthentication is the main helper to make requests requiring authentication.
func (m *Mnubo) doRequestWithAuthentication(cr ClientRequest, response interface{}) error {
	if m.isUsingStaticToken() {
		cr.authorization = fmt.Sprintf("Bearer %s", m.ClientToken)
	} else {
		if m.AccessToken.hasExpired() {
			_, err := m.GetAccessToken()

			if err != nil {
				return err
			}
		}
		cr.authorization = fmt.Sprintf("Bearer %s", m.AccessToken.Value)
	}

	err := m.doRequest(cr, response)

	return err
}
