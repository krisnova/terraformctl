package artifactory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPermissions(t *testing.T) {
	responseFile, err := os.Open("assets/test/permissions.json")
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
	perms, err := client.GetPermissionTargets()
	assert.NoError(t, err, "should not return an error")
	assert.Len(t, perms, 5, "should have five targets")
	assert.Equal(t, "snapshot-write", perms[0].Name, "Should have the snapshot-write target")
	assert.Equal(t, "https://artifactory/artifactory/api/security/permissions/snapshot-write", perms[0].URI, "should have a uri")
	for _, p := range perms {
		assert.NotNil(t, p.Name, "Name should not be empty")
		assert.NotNil(t, p.URI, "Uri should not be empty")
	}
}

func TestGetPermissionDetails(t *testing.T) {
	responseFile, err := os.Open("assets/test/permissions_details.json")
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
	perms, err := client.GetPermissionTargetDetails("release-commiter", make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, "release-commiter", perms.Name, "Should be release-commiter")
	assert.Equal(t, "**", perms.IncludesPattern, "Includes should be **")
	assert.Equal(t, "", perms.ExcludesPattern, "Excludes should be nil")
	assert.Len(t, perms.Repositories, 3, "Should have 3 repositories")
	assert.Contains(t, perms.Repositories, "docker-local-v2", "Should have repos")
	assert.NotNil(t, perms.Principals.Users, "should have a user principal")
	assert.Contains(t, perms.Principals.Users, "admin", "Should have the admin user")
	assert.Len(t, perms.Principals.Users["admin"], 5, "should have 5 permissions")
	assert.Contains(t, perms.Principals.Users["admin"], "m", "Should have the m permission")
	assert.NotNil(t, perms.Principals.Groups, "should have a group principal")
	groups := []string{}
	for g := range perms.Principals.Groups {
		groups = append(groups, g)
	}
	assert.Contains(t, groups, "java-committers", "Should have the java committers group")
	assert.Len(t, perms.Principals.Groups["java-committers"], 4, "should have 4 permissions")
	assert.Contains(t, perms.Principals.Groups["java-committers"], "r", "Should have the r permission")
}

func TestCreatePermissionTarget(t *testing.T) {
	var buf bytes.Buffer
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		req, _ := ioutil.ReadAll(r.Body)
		buf.Write(req)
		fmt.Fprintf(w, "")
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

	permTarget := PermissionTargetDetails{
		Name:            "release-commiter",
		IncludesPattern: "**",
		ExcludesPattern: "",
		Repositories:    []string{"docker-local-v2", "libs-release-local", "plugins-release-local"},
		Principals: Principals{
			Users:  map[string][]string{"admin": []string{"r", "d", "w", "n", "m"}},
			Groups: map[string][]string{"java-committers": []string{"r", "d", "w", "n"}},
		},
	}

	expectedJSON, _ := json.Marshal(permTarget)
	err := client.CreatePermissionTarget("release-commiter", permTarget, make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, string(expectedJSON), string(buf.Bytes()), "should send permission target json")
}
