package artifactory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Artifact struct {
	Info   FileInfo
	Client *ArtifactoryClient
}

type ArtifactProperties map[string][]string
type AQLPropertes struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}
type FileInfo struct {
	Uri          string             `json:"uri"`
	DownloadUri  string             `json:"downloadUri"`
	Repo         string             `json:"repo"`
	Path         string             `json:"path"`
	RemoteUrl    string             `json:"remoteUrl,omitempty"`
	Created      string             `json:"created"`
	CreatedBy    string             `json:"createdBy"`
	LastModified string             `json:"lastModified"`
	ModifiedBy   string             `json:"modifiedBy"`
	MimeType     string             `json:"mimeType"`
	Size         string             `json:"size"`
	Properties   ArtifactProperties `json:"properties"`
	Checksums    struct {
		SHA1 string `json:"sha1"`
		MD5  string `json:"md5"`
	} `json:"checksums"`
	OriginalChecksums struct {
		SHA1 string `json:"sha1"`
		MD5  string `json:"md5"`
	} `json:"originalChecksums,omitempty"`
}

type AQLResults struct {
	Results []AQLFileInfo `json:"results"`
	Range   struct {
		StartPos int `json:"start_pos"`
		EndPos   int `json:"end_pos"`
		Total    int `json:"total"`
		Limit    int `json:"limit"`
	} `json:"range"`
}
type AQLFileInfo struct {
	Repo         string         `json:"repo,omitempty"`
	Path         string         `json:"path,omitempty"`
	Name         string         `json:"name,omitempty"`
	Type         string         `json:"type,omitempty"`
	Created      string         `json:"created,omitempty"`
	CreatedBy    string         `json:"created_by,omitempty"`
	Modified     string         `json:"modified,omitempty"`
	ModifiedBy   string         `json:"modified_by,omitempty"`
	Depth        int            `json:"depth,omitempty"`
	Size         int64          `json:"size,omitempty"`
	Properties   []AQLPropertes `json:"properties,omitempty"`
	ActualMD5    string         `json:"actual_md5,omitempty"`
	ActualSHA1   string         `json:"actual_sha1,omitempty"`
	OriginalSHA1 string         `json:"original_sha1,omitempty"`
}

func (c *Artifact) Download() ([]byte, error) {
	return c.Client.RetrieveArtifact(c.Info.Repo, c.Info.Path)
}

func (c *Artifact) Delete() error {
	_, err := c.Client.DeleteArtifact(c.Info.Repo, c.Info.Path)
	return err
}

func (c *ArtifactoryClient) GetFileInfo(path string) (a Artifact, err error) {
	a.Client = c
	var res FileInfo
	d, err := c.HttpRequest(ArtifactoryRequest{
		Verb: "GET",
		Path: "/api/storage/" + path,
	})
	if err != nil {
		return a, err
	} else {
		e := json.Unmarshal(d, &res)
		if e != nil {
			return a, e
		} else {
			a.Info = res
			return a, nil
		}
	}
}

func (c *ArtifactoryClient) DeleteArtifact(repo, path string) ([]byte, error) {
	return c.HttpRequest(ArtifactoryRequest{
		Verb: "DELETE",
		Path: "/" + repo + "/" + path,
	})

}

func (c *ArtifactoryClient) RetrieveArtifact(repo string, path string) ([]byte, error) {
	return c.HttpRequest(ArtifactoryRequest{
		Verb: "GET",
		Path: "/" + repo + "/" + path,
	})
}

func (c *ArtifactoryClient) DeployArtifact(repoKey string, filename string, path string, properties map[string]string) (CreatedStorageItem, error) {
	var res CreatedStorageItem
	var fileProps []string
	var finalUrl string
	finalUrl = "/" + repoKey + "/"
	if &path != nil {
		finalUrl = finalUrl + path
	}
	baseFile := filepath.Base(filename)
	finalUrl = finalUrl + "/" + baseFile
	if len(properties) > 0 {
		finalUrl = finalUrl + ";"
		for k, v := range properties {
			fileProps = append(fileProps, k+"="+v)
		}
		finalUrl = finalUrl + strings.Join(fileProps, ";")
	}
	data, err := os.Open(filename)
	if err != nil {
		return res, err
	}
	defer data.Close()
	d, err := c.HttpRequest(ArtifactoryRequest{
		Verb: "PUT",
		Path: finalUrl,
		Body: data,
	})
	if err != nil {
		return res, err
	} else {
		e := json.Unmarshal(d, &res)
		if e != nil {
			return res, e
		} else {
			return res, nil
		}
	}
}
