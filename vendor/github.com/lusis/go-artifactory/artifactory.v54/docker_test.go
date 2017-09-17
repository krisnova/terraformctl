package artifactory

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDockerRepoImages(t *testing.T) {
	responseFile, err := os.Open("assets/test/docker_repositories.json")
	if err != nil {
		t.Fatalf("Unable to read test data: %s", err.Error())
	}
	defer func() { _ = responseFile.Close() }()
	responseBody, _ := ioutil.ReadAll(responseFile)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(responseBody))
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
	repos, err := client.GetDockerRepoImages("docker", make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Len(t, repos, 6, "should have six images")
	assert.Equal(t, "docker-dev", repos[0], "Should have the docker-dev image")
}

func TestGetDockerRepoImageTags(t *testing.T) {
	responseFile, err := os.Open("assets/test/docker_repository_tags.json")
	if err != nil {
		t.Fatalf("Unable to read test data: %s", err.Error())
	}
	defer func() { _ = responseFile.Close() }()
	responseBody, _ := ioutil.ReadAll(responseFile)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(responseBody))
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
	tags, err := client.GetDockerRepoImageTags("docker", "docker-dev", make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Len(t, tags, 5, "should have five tags")
	assert.Equal(t, "0.1.0", tags[0], "Should have the 0.1.0 image")
}
