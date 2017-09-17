package artifactory

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

// ClientConfig is the configuration for an ArtifactoryClient
type ClientConfig struct {
	BaseURL    string
	Username   string
	Password   string
	Token      string
	AuthMethod string
	VerifySSL  bool
	Client     *http.Client
	Transport  *http.Transport
}

// Client is a client for interacting with Artifactory
type Client struct {
	Client    *http.Client
	Config    *ClientConfig
	Transport *http.Transport
}

// NewClient returns a new ArtifactoryClient with the provided ClientConfig
func NewClient(config *ClientConfig) (c Client) {
	verifySSL := func() bool {
		return !config.VerifySSL
	}
	if config.Transport == nil {
		config.Transport = &http.Transport{}
	}
	config.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: verifySSL()}
	if config.Client == nil {
		config.Client = &http.Client{}
	}
	config.Client.Transport = config.Transport
	return Client{Client: config.Client, Config: config, Transport: config.Transport}
}

func clientConfigFrom(from string) (c *ClientConfig) {
	conf := ClientConfig{}
	switch from {
	case "environment":
		if os.Getenv("ARTIFACTORY_URL") == "" {
			fmt.Printf("You must set the environment variable ARTIFACTORY_URL")
			os.Exit(1)
		} else {
			conf.BaseURL = os.Getenv("ARTIFACTORY_URL")
		}
		if os.Getenv("ARTIFACTORY_TOKEN") == "" {
			if os.Getenv("ARTIFACTORY_USERNAME") == "" || os.Getenv("ARTIFACTORY_PASSWORD") == "" {
				fmt.Printf("You must set the environment variables ARTIFACTORY_USERNAME/ARTIFACTORY_PASSWORD\n")
				os.Exit(1)
			} else {
				conf.AuthMethod = "basic"
			}
		} else {
			conf.AuthMethod = "token"
		}
	}
	if conf.AuthMethod == "token" {
		conf.Token = os.Getenv("ARTIFACTORY_TOKEN")
	} else {
		conf.Username = os.Getenv("ARTIFACTORY_USERNAME")
		conf.Password = os.Getenv("ARTIFACTORY_PASSWORD")
	}
	return &conf
}

// NewClientFromEnv returns a new ArtifactoryClient the is automatically configured from environment variables
func NewClientFromEnv() (c Client) {
	config := clientConfigFrom("environment")
	client := NewClient(config)
	return client
}
