package artifactory

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientCustomTransport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "pong")
	}))
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	defer server.Close()
	conf := &ClientConfig{
		BaseURL:   "http://127.0.0.1:8080/",
		Username:  "username",
		Password:  "password",
		VerifySSL: false,
		Transport: transport,
	}

	client := NewClient(conf)
	res, err := client.Get("/ping", make(map[string]string))
	assert.Nil(t, err, "should not return an error")
	assert.NotNil(t, client.Transport)
	assert.Equal(t, "pong", string(res), "should return the testmsg")
}

func TestClientHTTPVerifySSLTrue(t *testing.T) {
	conf := &ClientConfig{VerifySSL: true}
	client := NewClient(conf)
	assert.False(t, client.Transport.TLSClientConfig.InsecureSkipVerify)
}

func TestClientHTTPVerifySSLFalse(t *testing.T) {
	conf := &ClientConfig{VerifySSL: false}
	client := NewClient(conf)
	assert.True(t, client.Transport.TLSClientConfig.InsecureSkipVerify)
}

func TestClientFromEnvWithBasicAuth(t *testing.T) {
	os.Setenv("ARTIFACTORY_URL", "http://artifactory.server.com")
	os.Setenv("ARTIFACTORY_USERNAME", "admin")
	os.Setenv("ARTIFACTORY_PASSWORD", "password")
	os.Setenv("ARTIFACTORY_TOKEN", "")
	client := NewClientFromEnv()
	assert.NotNil(t, client)
	assert.Equal(t, "http://artifactory.server.com", client.Config.BaseURL)
	assert.Equal(t, "basic", client.Config.AuthMethod)
	assert.Equal(t, "admin", client.Config.Username)
	assert.Equal(t, "password", client.Config.Password)
}

func TestClientFromEnvWithTokenAuth(t *testing.T) {
	os.Setenv("ARTIFACTORY_URL", "http://artifactory.server.com")
	os.Setenv("ARTIFACTORY_TOKEN", "someToken")
	client := NewClientFromEnv()
	assert.NotNil(t, client)
	assert.Equal(t, "http://artifactory.server.com", client.Config.BaseURL)
	assert.Equal(t, "token", client.Config.AuthMethod)
	assert.Equal(t, "someToken", client.Config.Token)
}
