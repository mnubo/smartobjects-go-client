package mnubo

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	n := NewClientWithToken("TOKEN", "HOST")

	at, err := m.GetAccessToken()
	now := time.Now()

	if err != nil {
		t.Errorf("unable to get access token: %s", err)
	}
	if at.ExpiresIn <= 0 {
		t.Errorf("access token expiration timestamp is invalid %+v", at)
	}

	if n.ClientToken != "TOKEN" || n.Host != "HOST" {
		t.Errorf("creating client with token should set ClientToken and Host: %+v", n)
	}

	if at.ExpiresAt.Before(now) {
		t.Errorf("access token expiration time %s should be in after now %s", at.ExpiresAt, now)
	}

	at.ExpiresAt = time.Now()
	m.AccessToken = at
}

func TestAccessToken(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	at, _ := m.GetAccessToken()
	now := time.Now()

	if m.AccessToken.hasExpired() == true {
		t.Errorf("access token should not expire so soon")
	}

	firstTokenValue := m.AccessToken.Value
	eat := at
	eat.ExpiresAt = now
	m.AccessToken = eat

	if m.AccessToken.hasExpired() == false {
		t.Errorf("access token should expire after a while")
	}

	var results interface{}
	cr := ClientRequest{
		method: "GET",
		path:   "test",
	}

	m.doRequestWithAuthentication(cr, &results)
	secondTokenValue := m.AccessToken.Value

	if firstTokenValue == secondTokenValue {
		t.Errorf("authentication should re-fetch token after expiration")
	}
}

func TestCompression(t *testing.T) {
	var results SearchResults

	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))

	compression := CompressionConfig{
		Request:  true,
		Response: true,
	}
	m.Compression = compression

	err := m.Search.CreateBasicQueryWithString(`{ "from": "event", "select": [ { "count": "*" } ] }`, &results)

	if err != nil {
		t.Errorf("error while running the query: %t", err)
	}

	if len(results.Rows) != 1 || len(results.Rows[0]) != 1 {
		t.Errorf("expecting results to have a count in firt row and cell")
	}
}

func TestExponentialBackoff(t *testing.T) {
	simFailures := 4
	var sendError = simFailures
	var gotNotified = 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sendError > 0 {
			http.Error(w, "", http.StatusServiceUnavailable)
			sendError--
		} else {
			buff := bytes.NewBufferString("{\"access_token\":\"valid\"}")
			w.Write(buff.Bytes())
		}
	}))
	defer ts.Close()
	m.Host = ts.URL
	m.ExponentialBackoff = ExponentialBackoffConfig{
		MaxElapsedTime: time.Second * 5,
		NotifyOnError: func(e error, duration time.Duration) {
			gotNotified++
		},
	}

	_, err := m.GetAccessToken()

	if err != nil {
		t.Errorf("unable to call client: %+v", err)
	}

	if gotNotified != simFailures {
		t.Errorf("expecting %d failures, got notified only %d times", simFailures, gotNotified)
	}
}

func TestClientErrors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "This is an error returned by the platform.", http.StatusBadRequest)
	}))
	defer ts.Close()
	m.Host = ts.URL

	_, err := m.GetAccessToken()

	if err == nil {
		t.Errorf("The client should bubble up the errors from the platform")
	}

	got := strings.TrimSpace(fmt.Sprintf("%s", err))
	expect := "The server responded with StatusCode: 400 - Body: This is an error returned by the platform."

	if strings.Compare(got, expect) != 0 {
		t.Errorf("expected: '%s', got: '%s'", expect, got)
	}
}
