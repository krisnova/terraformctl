package artifactory

import "encoding/json"

// DockerImages represents the list of docker images in a docker repo
type DockerImages struct {
	Repositories []string `json:"repositories,omitempty"`
}

// DockerImageTags represents the list of tags for an image in a docker repo
type DockerImageTags struct {
	Name string   `json:"name,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

// GetDockerRepoImages returns the docker images in the named repo
func (c *Client) GetDockerRepoImages(key string, q map[string]string) ([]string, error) {
	var dat DockerImages

	d, err := c.Get("/api/docker/"+key+"/v2/_catalog", q)
	if err != nil {
		return dat.Repositories, err
	}

	err = json.Unmarshal(d, &dat)
	if err != nil {
		return dat.Repositories, err
	}

	return dat.Repositories, nil
}

// GetDockerRepoImageTags returns the docker images in the named repo
func (c *Client) GetDockerRepoImageTags(key, image string, q map[string]string) ([]string, error) {
	var dat DockerImageTags

	d, err := c.Get("/api/docker/"+key+"/v2/"+image+"/tags/list", q)
	if err != nil {
		return dat.Tags, err
	}

	err = json.Unmarshal(d, &dat)
	if err != nil {
		return dat.Tags, err
	}

	return dat.Tags, nil
}
