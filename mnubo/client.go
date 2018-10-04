package mnubo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CompressionConfig struct {
	request  bool
	response bool
}

type Mnubo struct {
	ClientId     string
	ClientSecret string
	ClientToken  string
	Host         string
	AccessToken  AccessToken
	Compression  CompressionConfig
}

type ClientRequest struct {
	authorization string
	method        string
	path          string
	contentType   string
	payload       []byte
}

type AccessToken struct {
	Value     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	Scope     string `json:"scope"`
	Jti       string `json:"jti"`
}

func NewClient(id string, secret string, host string) *Mnubo {
	return &Mnubo{
		ClientId:     id,
		ClientSecret: secret,
		Host:         host,
	}
}

func NewClientWithToken(token string, host string) *Mnubo {
	return &Mnubo{
		ClientToken: token,
		Host:        host,
	}
}

func (m *Mnubo) getAccessToken() (AccessToken, error) {
	return m.getAccessTokenWithScope("ALL")
}

func (m *Mnubo) getAccessTokenWithScope(scope string) (AccessToken, error) {
	payload := fmt.Sprintf("grant_type=client_credentials&scope=%s", scope)
	data := []byte(fmt.Sprintf("%s:%s", m.ClientId, m.ClientSecret))

	cr := ClientRequest{
		authorization: fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(data)),
		method:        "POST",
		path:          "/oauth/token",
		contentType:   "application/x-www-form-urlencoded",
		payload:       []byte(payload),
	}
	at := AccessToken{}
	body, err := m.doRequest(cr)

	if err == nil {
		err = json.Unmarshal(body, &at)
		m.AccessToken = at
		return at, err
	}
	return at, err
}

func (m *Mnubo) doRequest(cr ClientRequest) ([]byte, error) {
	req, err := http.NewRequest(cr.method, m.Host+cr.path, bytes.NewReader(cr.payload))

	req.Header.Add("Content-Type", cr.contentType)
	req.Header.Add("X-MNUBO-SDK", "Go")

	if cr.authorization != "" {
		req.Header.Add("Authorization", cr.authorization)
	}

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
		return body, nil
	}

	return nil, fmt.Errorf("request Error: %s", body)
}

func (m *Mnubo) doRequestWithAuthentication(cr ClientRequest, response interface{}) error {
	if m.ClientToken != "" {
		cr.authorization = fmt.Sprintf("Bearer %s", m.ClientToken)
	} else {
		cr.authorization = fmt.Sprintf("Bearer %s", m.AccessToken.Value)
	}

	data, err := m.doRequest(cr)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, response)
}
