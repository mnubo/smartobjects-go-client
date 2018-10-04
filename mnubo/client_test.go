package mnubo

import (
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	m := NewClient(os.Getenv("MNUBO_CLIENT_ID"), os.Getenv("MNUBO_CLIENT_SECRET"), os.Getenv("MNUBO_HOST"))
	n := NewClientWithToken("TOKEN", "HOST")

	at, err := m.getAccessToken()
	if err != nil {
		t.Errorf("Unable to get access token: %s", err)
	}
	if at.ExpiresIn <= 0 {
		t.Errorf("Access Token expiration timestamp is invalid %+v", at)
	}

	if n.ClientToken != "TOKEN" || n.Host != "HOST" {
		t.Errorf("Creating client with token should set ClientToken and Host: %+v", n)
	}
}
