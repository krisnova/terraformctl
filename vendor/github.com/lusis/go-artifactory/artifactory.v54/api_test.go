package artifactory

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserAPIKey(t *testing.T) {
	apiKey := UserAPIKey{
		APIKey: "testAPIKey",
	}
	body, _ := json.Marshal(apiKey)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(body))
	}))
	defer server.Close()

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	conf := &ClientConfig{
		BaseURL:   "http://127.0.0.1:8080/",
		Username:  "username",
		Password:  "password",
		VerifySSL: false,
		Transport: transport,
	}

	client := NewClient(conf)
	key, err := client.GetUserAPIKey()
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, "testAPIKey", key, "key should be testAPIKey")
}

func TestCreateUserAPIKey(t *testing.T) {
	apiKey := UserAPIKey{
		APIKey: "testAPIKey",
	}
	body, _ := json.Marshal(apiKey)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(body))
	}))
	defer server.Close()

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	conf := &ClientConfig{
		BaseURL:   "http://127.0.0.1:8080/",
		Username:  "username",
		Password:  "password",
		VerifySSL: false,
		Transport: transport,
	}

	client := NewClient(conf)
	key, err := client.CreateUserAPIKey()
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, "testAPIKey", key, "key should be testAPIKey")
}
