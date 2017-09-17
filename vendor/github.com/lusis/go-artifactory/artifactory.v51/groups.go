package artifactory

import (
	"encoding/json"
)

// Group represents the json response for a group in Artifactory
type Group struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// GroupDetails represents the json response for a group's details in artifactory
type GroupDetails struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	AutoJoin        bool   `json:"autoJoin,omitempty"`
	Realm           string `json:"realm,omitempty"`
	RealmAttributes string `json:"realmAttributes,omitempty"`
}

// GetGroups gets a list of groups from artifactory
func (c *Client) GetGroups() ([]Group, error) {
	var res []Group
	d, e := c.Get("/api/security/groups", make(map[string]string))
	if e != nil {
		return res, e
	}
	err := json.Unmarshal(d, &res)
	if err != nil {
		return res, err
	}
	return res, e
}

// GetGroupDetails returns details for a Group
func (c *Client) GetGroupDetails(u string) (GroupDetails, error) {
	var res GroupDetails
	d, e := c.Get("/api/security/groups/"+u, make(map[string]string))
	if e != nil {
		return res, e
	}
	err := json.Unmarshal(d, &res)
	if err != nil {
		return res, err
	}
	return res, e
}

// CreateGroup creates a group in artifactory
func (c *Client) CreateGroup(gname string, g GroupDetails) error {
	j, jerr := json.Marshal(g)
	if jerr != nil {
		return jerr
	}
	o := make(map[string]string)
	_, err := c.Put("/api/security/groups/"+gname, string(j), o)
	return err
}
