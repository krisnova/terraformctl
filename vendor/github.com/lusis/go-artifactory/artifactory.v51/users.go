package artifactory

import (
	"bytes"
	"encoding/json"
	"errors"
)

// User represents a user in artifactory
type User struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// UserDetails represents the details of a user in artifactory
type UserDetails struct {
	Name                     string   `json:"name,omitempty"`
	Email                    string   `json:"email"`
	Password                 string   `json:"password"`
	Admin                    bool     `json:"admin,omitempty"`
	ProfileUpdatable         bool     `json:"profileUpdatable,omitempty"`
	InternalPasswordDisabled bool     `json:"internalPasswordDisabled,omitempty"`
	LastLoggedIn             string   `json:"lastLoggedIn,omitempty"`
	Realm                    string   `json:"realm,omitempty"`
	Groups                   []string `json:"groups,omitempty"`
}

// GetUsers returns all users
func (c *Client) GetUsers() ([]User, error) {
	var res []User
	d, e := c.HTTPRequest(Request{
		Verb: "GET",
		Path: "/api/security/users",
	})
	if e != nil {
		return res, e
	}
	err := json.Unmarshal(d, &res)
	if err != nil {
		return res, err
	}
	return res, e
}

// GetUserDetails returns details for the named user
func (c *Client) GetUserDetails(u string) (UserDetails, error) {
	var res UserDetails
	d, e := c.HTTPRequest(Request{
		Verb: "GET",
		Path: "/api/security/users/" + u,
	})
	if e != nil {
		return res, e
	}
	err := json.Unmarshal(d, &res)
	if err != nil {
		return res, err
	}
	return res, e
}

// CreateUser creates a user with the specified details
func (c *Client) CreateUser(uname string, u UserDetails) error {
	if &u.Email == nil || &u.Password == nil {
		return errors.New("Email and password are required to create users")
	}
	j, jerr := json.Marshal(u)
	if jerr != nil {
		return jerr
	}
	o := make(map[string]string)
	_, err := c.HTTPRequest(Request{
		Verb:        "PUT",
		Path:        "/api/security/users/" + uname,
		Body:        bytes.NewReader(j),
		QueryParams: o,
	})
	return err
}

// DeleteUser deletes a user
func (c *Client) DeleteUser(uname string) error {
	_, err := c.HTTPRequest(Request{
		Verb: "DELETE",
		Path: "/api/security/users/" + uname,
	})
	return err
}

// GetUserEncryptedPassword returns the current user's encrypted password
func (c *Client) GetUserEncryptedPassword() (string, error) {
	d, err := c.HTTPRequest(Request{
		Verb: "GET",
		Path: "/api/security/encryptedPassword",
	})
	return string(d), err
}

// GetUserAPIKey returns the current user's api key
func (c *Client) GetUserAPIKey() (string, error) {
	var res UserAPIKey
	d, err := c.HTTPRequest(Request{
		Verb: "GET",
		Path: "/api/security/apiKey",
	})
	if err != nil {
		return "", err
	}
	jsonErr := json.Unmarshal(d, &res)
	if jsonErr != nil {
		return "", jsonErr
	}
	return res.APIKey, nil
}

// CreateUserAPIKey creates an apikey for the current user
func (c *Client) CreateUserAPIKey() (string, error) {
	var res UserAPIKey
	d, err := c.HTTPRequest(Request{
		Verb: "POST",
		Path: "/api/security/apiKey",
	})
	if err != nil {
		return "", err
	}
	jsonErr := json.Unmarshal(d, &res)
	if jsonErr != nil {
		return "", jsonErr
	}
	return res.APIKey, nil
}
