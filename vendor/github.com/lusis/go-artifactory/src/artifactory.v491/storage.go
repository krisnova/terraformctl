package artifactory

import (
	"encoding/json"
)

type CreatedStorageItem struct {
	URI               string            `json:"uri"`
	DownloadURI       string            `json:"downloadUri"`
	Repo              string            `json:"repo"`
	Created           string            `json:"created"`
	CreatedBy         string            `json:"createdBy"`
	Size              string            `json:"size"`
	MimeType          string            `json:"mimeType"`
	Checksums         ArtifactChecksums `json:"checksums"`
	OriginalChecksums ArtifactChecksums `json:"originalChecksums"`
}

type ArtifactChecksums struct {
	MD5  string `json:"md5"`
	SHA1 string `json:"sha1"`
}

type FileList struct {
	URI     string         `json:"uri"`
	Created string         `json:"created"`
	Files   []FileListItem `json:"files"`
}

type FileListItem struct {
	URI          string `json:"uri"`
	Size         int    `json:"size"`
	LastModified string `json:"lastModified"`
	Folder       bool   `json:"folder"`
	SHA1         string `json:"sha1"`
}

func (c *ArtifactoryClient) ListFiles(repo string) (fileList FileList, err error) {
	d, err := c.HttpRequest(ArtifactoryRequest{
		Verb: "GET",
		Path: "/api/storage/" + repo + "?list&deep=1",
	})
	if err != nil {
		return fileList, err
	}
	err = json.Unmarshal(d, &fileList)
	return fileList, err
}
