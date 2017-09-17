package artifactory

// GetGeneralConfiguration returns the general artifactory configuration
func (c *Client) GetGeneralConfiguration() (s string, e error) {
	d, e := c.Get("/api/system/configuration", make(map[string]string))
	return string(d), e
}
