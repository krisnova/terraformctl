package artifactory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Gavc struct {
	GroupID    string
	ArtifactID string
	Version    string
	Classifier string
	Repos      []string
}

func (c *ArtifactoryClient) GAVCSearch(coords *Gavc) (files []FileInfo, e error) {
	url := "/api/search/gavc"
	params := make(map[string]string)
	if &coords.GroupID != nil {
		params["g"] = coords.GroupID
	}
	if &coords.ArtifactID != nil {
		params["a"] = coords.ArtifactID
	}
	if &coords.Version != nil {
		params["v"] = coords.Version
	}
	if &coords.Classifier != nil {
		params["c"] = coords.Classifier
	}
	if &coords.Repos != nil {
		params["repos"] = strings.Join(coords.Repos, ",")
	}
	d, err := c.Get(url, params)
	if err != nil {
		return files, err
	} else {
		var dat GavcSearchResults
		err := json.Unmarshal(d, &dat)
		if err != nil {
			return files, err
		} else {
			files = dat.Results
			return files, nil
		}
	}
}

func (c *ArtifactoryClient) DockerSearch(name string) (files []FileInfo, e error) {
	var request ArtifactoryRequest
	params := make(map[string]string)
	params["docker.repoName"] = fmt.Sprintf("*%s*", name)
	request.Verb = "GET"
	request.Path = "/api/search/prop"
	request.QueryParams = params
	request.ContentType = "application/json"
	data, err := c.HttpRequest(request)
	if err != nil {
		return files, err
	}
	var dat GavcSearchResults
	uerr := json.Unmarshal(data, &dat)
	if uerr != nil {
		return files, uerr
	}
	files = dat.Results
	return files, nil
}

func (c *ArtifactoryClient) VagrantSearch(name string) (files []AQLFileInfo, e error) {
	var request ArtifactoryRequest
	request.Verb = "POST"
	request.Path = "/api/search/aql"
	aqlString := fmt.Sprintf(`items.find(
  {
    "$and":[
      {"$or":[
        {"@box_name":{"$match":"*%s*"}}
      ]},
      {"$rf":[
        {"$or":[
          {"property.key":{"$eq":"box_name"}},
          {"property.key":{"$eq":"box_version"}},
          {"property.key":{"$eq":"box_provider"}}
        ]}
      ]}
    ]
  }
).include("updated","created_by","repo","type","actual_md5","property.key","size","original_sha1","name","modified_by","original_md5","property.value","path","modified","id","actual_sha1","created","depth")`, name)
	request.Body = bytes.NewReader([]byte(aqlString))
	request.ContentType = "text/plain"
	data, err := c.HttpRequest(request)
	if err != nil {
		return files, err
	}
	var dat AQLResults
	uerr := json.Unmarshal(data, &dat)
	if uerr != nil {
		return files, uerr
	}
	files = dat.Results
	return files, nil
}
