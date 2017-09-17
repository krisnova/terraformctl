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

func TestGetRepositories(t *testing.T) {
	responseFile, err := os.Open("assets/test/repositories.json")
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
	repos, err := client.GetRepos("all")
	assert.NoError(t, err, "should not return an error")
	assert.Len(t, repos, 21, "should have twenty-one targets")
	assert.Equal(t, "dev-docker-local", repos[0].Key, "Should have the dev-docker-local repo")
	assert.Equal(t, "https://artifactory.domain/artifactory/dev-docker-local", repos[0].URL, "should have a uri")
	for _, r := range repos {
		assert.NotNil(t, r.Key, "Name should not be empty")
		assert.NotNil(t, r.URL, "Uri should not be empty")
		assert.NotNil(t, r.Rtype, "Type should not be empty")
	}
}

func TestGetRepositoryWithLocalRepoConfig(t *testing.T) {
	responseFile, err := os.Open("assets/test/local_repository_config.json")
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
	repo, err := client.GetRepo("libs-release-local", make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, LocalRepoMimeType, repo.MimeType(), "Repo mime type should be local repo")
	assert.Equal(t, "libs-release-local", repo.(LocalRepoConfig).Key, "Repo key should be libs-release-local")
	assert.Equal(t, "local", repo.(LocalRepoConfig).RClass, "Repo rclass should be local")
	assert.Equal(t, "maven", repo.(LocalRepoConfig).PackageType, "Repo package type should be maven")
	assert.Equal(t, "Local repository for in-house libraries", repo.(LocalRepoConfig).Description, "Repo description should be Local repository for in-house libraries")
	assert.Equal(t, "", repo.(LocalRepoConfig).Notes, "Repo notes should be empty")
	assert.Equal(t, "**/*", repo.(LocalRepoConfig).IncludesPattern, "Repo includes pattern should be **/*")
	assert.Equal(t, "", repo.(LocalRepoConfig).ExcludesPattern, "Repo excludes pattern should be empty")
	assert.Equal(t, "maven-2-default", repo.(LocalRepoConfig).LayoutRef, "Repo layout ref should be maven-2-default")
	assert.Equal(t, true, repo.(LocalRepoConfig).HandleReleases, "Repo handle releases should be true")
	assert.Equal(t, false, repo.(LocalRepoConfig).HandleSnapshots, "Repo handle snapshots should be false")
	assert.Equal(t, 0, repo.(LocalRepoConfig).MaxUniqueSnapshots, "Repo max unique snapshots should be 0")
	assert.Equal(t, true, repo.(LocalRepoConfig).SuppressPomConsistencyChecks, "Repo suppress pom consistency checks should be true")
	assert.Equal(t, false, repo.(LocalRepoConfig).BlackedOut, "Repo blacked out should be false")
	assert.Equal(t, []string{"artifactory"}, repo.(LocalRepoConfig).PropertySets, "Repo property sets should be array with artifactory")
	assert.Equal(t, false, repo.(LocalRepoConfig).DebianTrivialLayout, "Repo debian trivial layout should be false")
	assert.Equal(t, "client-checksums", repo.(LocalRepoConfig).ChecksumPolicyType, "Repo checksum policy type should be client-checksums")
	assert.Equal(t, 0, repo.(LocalRepoConfig).MaxUniqueTags, "Repo max unique tags should be 0")
	assert.Equal(t, "unique", repo.(LocalRepoConfig).SnapshotVersionBehavior, "Repo snapshot version behavior should be unique")
	assert.Equal(t, false, repo.(LocalRepoConfig).ArchiveBrowsingEnabled, "Repo archive browsing enabled should be false")
	assert.Equal(t, false, repo.(LocalRepoConfig).CalculateYumMetadata, "Repo calculate yum metadata should be false")
	assert.Equal(t, 0, repo.(LocalRepoConfig).YumRootDepth, "Repo yum root depth should be 0")
	assert.Equal(t, "V1", repo.(LocalRepoConfig).DockerAPIVersion, "Repo docker api version should be V1")
	assert.Equal(t, false, repo.(LocalRepoConfig).EnableFileListsIndexing, "Repo enable file lists indexing should be false")
}

func TestGetRepositoryWithRemoteRepoConfig(t *testing.T) {
	responseFile, err := os.Open("assets/test/remote_repository_config.json")
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
	repo, err := client.GetRepo("jcenter", make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, RemoteRepoMimeType, repo.MimeType(), "Repo mime type should be remote repo")
	assert.Equal(t, "jcenter", repo.(RemoteRepoConfig).Key, "Repo key should be jvrnyrt")
	assert.Equal(t, "remote", repo.(RemoteRepoConfig).RClass, "Repo rclass should be remote")
	assert.Equal(t, "maven", repo.(RemoteRepoConfig).PackageType, "Repo package type should be maven")
	assert.Equal(t, "Bintray Central Java repository", repo.(RemoteRepoConfig).Description, "Repo description should be Bintray Central Java repository")
	assert.Equal(t, "", repo.(RemoteRepoConfig).Notes, "Repo notes should be empty")
	assert.Equal(t, "**/*", repo.(RemoteRepoConfig).IncludesPattern, "Repo includes pattern should be **/*")
	assert.Equal(t, "", repo.(RemoteRepoConfig).ExcludesPattern, "Repo excludes pattern should be empty")
	assert.Equal(t, "maven-2-default", repo.(RemoteRepoConfig).LayoutRef, "Repo layout ref should be maven-2-default")
	assert.Equal(t, true, repo.(RemoteRepoConfig).HandleReleases, "Repo handle releases should be true")
	assert.Equal(t, false, repo.(RemoteRepoConfig).HandleSnapshots, "Repo handle snapshots should be false")
	assert.Equal(t, 0, repo.(RemoteRepoConfig).MaxUniqueSnapshots, "Repo max unique snapshots should be 0")
	assert.Equal(t, true, repo.(RemoteRepoConfig).SuppressPomConsistencyChecks, "Repo suppress pom consistency checks should be true")
	assert.Equal(t, false, repo.(RemoteRepoConfig).BlackedOut, "Repo blacked out should be false")
	assert.Equal(t, []string{"artifactory"}, repo.(RemoteRepoConfig).PropertySets, "Repo property sets should be array with artifactory")
	assert.Equal(t, "http://jcenter.bintray.com", repo.(RemoteRepoConfig).URL, "Repo url should be http://jcenter.bintray.com")
	assert.Equal(t, "", repo.(RemoteRepoConfig).Username, "Repo username should be empty")
	assert.Equal(t, "", repo.(RemoteRepoConfig).Password, "Repo password should be empty")
	assert.Equal(t, "proxy", repo.(RemoteRepoConfig).Proxy, "Repo proxy should be proxy")
	assert.Equal(t, "", repo.(RemoteRepoConfig).RemoteRepoChecksumPolicyType, "Repo remote repo checksum policy type should be empty")
	assert.Equal(t, false, repo.(RemoteRepoConfig).HardFail, "Repo hardfail should be false")
	assert.Equal(t, false, repo.(RemoteRepoConfig).Offline, "Repo offline should be false")
	assert.Equal(t, true, repo.(RemoteRepoConfig).StoreArtifactsLocally, "Repo store artifacts locally should be true")
	assert.Equal(t, 15000, repo.(RemoteRepoConfig).SocketTimeoutMillis, "Repo socket timeout millis should be 15000")
	assert.Equal(t, "", repo.(RemoteRepoConfig).LocalAddress, "Repo local address should be empty")
	assert.Equal(t, 43200, repo.(RemoteRepoConfig).RetrivialCachePeriodSecs, "Repo retrieval cache period secs should be 43200")
	assert.Equal(t, 7200, repo.(RemoteRepoConfig).FailedRetrievalCachePeriodSecs, "Repo failed retrieval cache period secs should be 7200")
	assert.Equal(t, 7200, repo.(RemoteRepoConfig).MissedRetrievalCachePeriodSecs, "Repo missed retrieval cache period secs should be 7200")
	assert.Equal(t, false, repo.(RemoteRepoConfig).UnusedArtifactsCleanupEnabled, "Repo unused artifact cleanup enabled should be false")
	assert.Equal(t, 0, repo.(RemoteRepoConfig).UnusedArtifactsCleanupPeriodHours, "Repo unused artifacts cleanup period hours should be 0")
	assert.Equal(t, false, repo.(RemoteRepoConfig).FetchJarsEagerly, "Repo fetch jars eagerly should be false")
	assert.Equal(t, false, repo.(RemoteRepoConfig).FetchSourcesEagerly, "Repo fetch sources eagerly should be false")
	assert.Equal(t, false, repo.(RemoteRepoConfig).ShareConfiguration, "Repo share configuration should be false")
	assert.Equal(t, false, repo.(RemoteRepoConfig).SynchronizeProperties, "Repo synchronize properties should be false")
	assert.Equal(t, false, repo.(RemoteRepoConfig).BlockMismatchingMimeTypes, "Repo block mismatching mime types should be false")
	assert.Equal(t, false, repo.(RemoteRepoConfig).AllowAnyHostAuth, "Repo allow any host auth should be false")
	assert.Equal(t, false, repo.(RemoteRepoConfig).EnableCookieManagement, "Repo enable cookie management should be false")
	assert.Equal(t, "http://bower.registry.com", repo.(RemoteRepoConfig).BowerRegistryURL, "Repo bower registry url should be http://bower.registry.com")
	assert.Equal(t, "", repo.(RemoteRepoConfig).VcsType, "Repo vcs type should be empty")
	assert.Equal(t, "", repo.(RemoteRepoConfig).VcsGitProvider, "Repo vcs git provider should be empty")
	assert.Equal(t, "", repo.(RemoteRepoConfig).VcsGitDownloader, "Repo vcs git downloader should be empty")
	assert.Equal(t, "", repo.(RemoteRepoConfig).ClientTLSCertificate, "Repo client tls certificate should be empty")
}

func TestGetRepositoryWithVirtualRepoConfig(t *testing.T) {
	responseFile, err := os.Open("assets/test/virtual_repository_config.json")
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
	repo, err := client.GetRepo("jcenter", make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, VirtualRepoMimeType, repo.MimeType(), "Repo mime type should be virtual repo")
	assert.Equal(t, "libs-release", repo.(VirtualRepoConfig).Key, "Repo key should be libs-release")
	assert.Equal(t, "virtual", repo.(VirtualRepoConfig).RClass, "Repo rclass should be virtual")
	assert.Equal(t, "maven", repo.(VirtualRepoConfig).PackageType, "Repo package type should be maven")
	assert.Equal(t, "", repo.(VirtualRepoConfig).Description, "Repo description should be empty")
	assert.Equal(t, "", repo.(VirtualRepoConfig).Notes, "Repo notes should be empty")
	assert.Equal(t, "**/*", repo.(VirtualRepoConfig).IncludesPattern, "Repo includes pattern should be **/*")
	assert.Equal(t, "", repo.(VirtualRepoConfig).ExcludesPattern, "Repo excludes pattern should be empty")
	assert.Equal(t, []string{"libs-release-local", "ext-release-local", "remote-repos"}, repo.(VirtualRepoConfig).Repositories, "Repo repositories should be array with libs-release-local, ext-release-lcocal, and remote-repos")
	assert.Equal(t, false, repo.(VirtualRepoConfig).DebianTrivialLayout, "Repo debian trivial layout should be false")
	assert.Equal(t, false, repo.(VirtualRepoConfig).ArtifactoryRequestsCanRetrieveRemoteArtifacts, "Repo artifactory requests can retrieve remote artifacts should be false")
	assert.Equal(t, "", repo.(VirtualRepoConfig).KeyPair, "Repo key pair should be empty")
	assert.Equal(t, "discard_active_reference", repo.(VirtualRepoConfig).PomRepositoryReferencesCleanupPolicy, "Repo pom repository references cleanup policy should be discard_active_reference")
	assert.Equal(t, "libs-release-local", repo.(VirtualRepoConfig).DefaultDeploymentRepo, "Repo default deployment repo should be libs-release-local")
}

func TestGetRepositoryWithGenericRepoConfig(t *testing.T) {
	responseFile, err := os.Open("assets/test/generic_repository_config.json")
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
	repo, err := client.GetRepo("libs-release-local", make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, "", repo.MimeType(), "Repo mime type should be empty")
	assert.Equal(t, "libs-release-local", repo.(GenericRepoConfig).Key, "Repo key should be libs-release-local")
	assert.Equal(t, "", repo.(GenericRepoConfig).RClass, "Repo rclass should be empty")
	assert.Equal(t, "maven", repo.(GenericRepoConfig).PackageType, "Repo package type should be maven")
	assert.Equal(t, "Local repository for in-house libraries", repo.(GenericRepoConfig).Description, "Repo description should be Local repository for in-house libraries")
	assert.Equal(t, "", repo.(GenericRepoConfig).Notes, "Repo notes should be empty")
	assert.Equal(t, "**/*", repo.(GenericRepoConfig).IncludesPattern, "Repo includes pattern should be **/*")
	assert.Equal(t, "", repo.(GenericRepoConfig).ExcludesPattern, "Repo excludes pattern should be empty")
	assert.Equal(t, "maven-2-default", repo.(GenericRepoConfig).LayoutRef, "Repo layout ref should be maven-2-default")
	assert.Equal(t, true, repo.(GenericRepoConfig).HandleReleases, "Repo handle releases should be true")
	assert.Equal(t, false, repo.(GenericRepoConfig).HandleSnapshots, "Repo handle snapshots should be false")
	assert.Equal(t, 0, repo.(GenericRepoConfig).MaxUniqueSnapshots, "Repo max unique snapshots should be 0")
	assert.Equal(t, true, repo.(GenericRepoConfig).SuppressPomConsistencyChecks, "Repo suppress pom consistency checks should be true")
	assert.Equal(t, false, repo.(GenericRepoConfig).BlackedOut, "Repo blacked out should be false")
	assert.Equal(t, []string{"artifactory"}, repo.(GenericRepoConfig).PropertySets, "Repo property sets should be array with artifactory")
}

func TestCreateRepository(t *testing.T) {
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

	repoConf := LocalRepoConfig{
		GenericRepoConfig: GenericRepoConfig{
			Key:                          "libs-release-local",
			RClass:                       "local",
			PackageType:                  "maven",
			Description:                  "Local repository for in-house libraries",
			Notes:                        "",
			IncludesPattern:              "**/*",
			ExcludesPattern:              "",
			LayoutRef:                    "maven-2-default",
			HandleReleases:               true,
			HandleSnapshots:              false,
			MaxUniqueSnapshots:           0,
			SuppressPomConsistencyChecks: true,
			BlackedOut:                   false,
			PropertySets:                 []string{"artifactory"},
		},
		DebianTrivialLayout:     false,
		ChecksumPolicyType:      "client-checksums",
		MaxUniqueTags:           0,
		SnapshotVersionBehavior: "unique",
		ArchiveBrowsingEnabled:  false,
		CalculateYumMetadata:    false,
		YumRootDepth:            0,
		DockerAPIVersion:        "V1",
		EnableFileListsIndexing: false,
	}

	expectedJSON, _ := json.Marshal(repoConf)
	err := client.CreateRepo("libs-release-local", repoConf, make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, string(expectedJSON), string(buf.Bytes()), "should send repos json")
}

func TestUpdateRepository(t *testing.T) {
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

	repoConf := LocalRepoConfig{
		GenericRepoConfig: GenericRepoConfig{
			Key:                          "libs-release-local",
			RClass:                       "local",
			PackageType:                  "maven",
			Description:                  "Local repository for in-house libraries",
			Notes:                        "",
			IncludesPattern:              "**/*",
			ExcludesPattern:              "",
			LayoutRef:                    "maven-2-default",
			HandleReleases:               true,
			HandleSnapshots:              false,
			MaxUniqueSnapshots:           0,
			SuppressPomConsistencyChecks: true,
			BlackedOut:                   false,
			PropertySets:                 []string{"artifactory"},
		},
	}

	expectedJSON, _ := json.Marshal(repoConf)
	err := client.UpdateRepo("libs-release-local", repoConf, make(map[string]string))
	assert.NoError(t, err, "should not return an error")
	assert.Equal(t, string(expectedJSON), string(buf.Bytes()), "should send repo json")
}
