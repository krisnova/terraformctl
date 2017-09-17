package artifactory

import (
	"bytes"
	"encoding/json"
	"errors"
)

type User struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

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

func (c *ArtifactoryClient) GetUsers() ([]User, error) {
	var res []User
	d, e := c.HttpRequest(ArtifactoryRequest{
		Verb: "GET",
		Path: "/api/security/users",
	})
	if e != nil {
		return res, e
	} else {
		err := json.Unmarshal(d, &res)
		if err != nil {
			return res, err
		} else {
			return res, e
		}
	}
}

func (c *ArtifactoryClient) GetUserDetails(u string) (UserDetails, error) {
	var res UserDetails
	d, e := c.HttpRequest(ArtifactoryRequest{
		Verb: "GET",
		Path: "/api/security/users/" + u,
	})
	if e != nil {
		return res, e
	} else {
		err := json.Unmarshal(d, &res)
		if err != nil {
			return res, err
		} else {
			return res, e
		}
	}
}

func (c *ArtifactoryClient) CreateUser(uname string, u UserDetails) error {
	if &u.Email == nil || &u.Password == nil {
		return errors.New("Email and password are required to create users")
	}
	j, jerr := json.Marshal(u)
	if jerr != nil {
		return jerr
	}
	o := make(map[string]string)
	_, err := c.HttpRequest(ArtifactoryRequest{
		Verb:        "PUT",
		Path:        "/api/security/users/" + uname,
		Body:        bytes.NewReader(j),
		QueryParams: o,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ArtifactoryClient) DeleteUser(uname string) error {
	_, err := c.HttpRequest(ArtifactoryRequest{
		Verb: "DELETE",
		Path: "/api/security/users/" + uname,
	})
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *ArtifactoryClient) GetUserEncryptedPassword() (s string, err error) {
	d, err := c.HttpRequest(ArtifactoryRequest{
		Verb: "GET",
		Path: "/api/security/encryptedPassword",
	})
	return string(d), err
}

func (c *ArtifactoryClient) GetUserApiKey() (s string, err error) {
	var res UserApiKey
	d, err := c.HttpRequest(ArtifactoryRequest{
		Verb: "GET",
		Path: "/api/security/apiKey",
	})
	if err != nil {
		return s, err
	} else {
		err := json.Unmarshal(d, &res)
		if err != nil {
			return s, err
		} else {
			return res.ApiKey, nil
		}
	}
}

func (c *ArtifactoryClient) CreateUserApiKey() (s string, err error) {
	var res UserApiKey
	d, err := c.HttpRequest(ArtifactoryRequest{
		Verb: "POST",
		Path: "/api/security/apiKey",
	})
	if err != nil {
		return s, err
	} else {
		err := json.Unmarshal(d, &res)
		if err != nil {
			return s, err
		}
		return res.ApiKey, nil
	}
}
