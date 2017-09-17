package artifactory

// GetGeneralConfiguration returns the general artifactory configuration
func (c *Client) GetGeneralConfiguration() (s string, e error) {
	d, e := c.Get("/api/system/configuration", make(map[string]string))
	return string(d), e
}

// GetSystemHealthPing returns a simple status response about the state of Artifactory
func (c *Client) GetSystemHealthPing() (s string, e error) {
	d, e := c.Get("/api/system/ping", make(map[string]string))
	return string(d), e
}
