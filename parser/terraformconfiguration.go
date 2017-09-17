package parser

import (
	"encoding/json"
	"fmt"
	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/mitchellh/hashstructure"
	"io/ioutil"
	"os"
	"strings"
)

// TerraformConfiguration represents a terraform directory on a local filesystem.
type TerraformConfiguration struct {
	Name      string
	tfBytes   []byte
	applyHash string
}

// NewTerraformConfigurationFromPath will attempt to parse a path on a local filesystem and build a TerraformConfiguration struct from it.
func NewTerraformConfigurationFromPath(name, path string) (*TerraformConfiguration, error) {
	if path == "." {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("Unable to parse working directory with error: %v", err)
		}
		path = wd
	}
	if !strings.HasPrefix(path, "/") {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("Unable to parse working directory with error: %v", err)
		}
		abs := fmt.Sprintf("%s/%s", wd, path)
		path = abs
	}
	fi, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to stat file [%s] with error: %v", path, err)
	}
	var filesToRead []string
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// Directory
		logger.Debug("Parsing directory [%s] for terraform configuration files", path)

		if !strings.HasSuffix(path, "/") {
			path = fmt.Sprintf("%s/", path)
		}

		allFiles, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("Unable to read directory: %v", err)
		}
		for _, file := range allFiles {
			if strings.HasSuffix(file.Name(), ".tf") {
				filesToRead = append(filesToRead, fmt.Sprintf("%s%s", path, file.Name()))
			}
		}
	case mode.IsRegular():
		// File
		logger.Debug("Parsing file [%s] for terraform configuration", path)
		filesToRead[0] = path
	default:
		return nil, fmt.Errorf("Invalid path [%s]", path)
	}

	var tfBytes []byte
	for _, file := range filesToRead {
		logger.Debug("Parsing file [%s]", file)
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("Unable to read file [%s] with error: %v", file, err)
		}
		tfBytes = append(tfBytes, bytes...)
	}
	return &TerraformConfiguration{
		Name:    name,
		tfBytes: tfBytes,
	}, nil
}

// NewTerraformConfigurationFromString will create a new terraform configuration from an encoded string that could be exchanged over gRPC.
func NewTerraformConfigurationFromEncodedString(encodedConfig string) (*TerraformConfiguration, error) {
	t := &TerraformConfiguration{}
	err := json.Unmarshal([]byte(encodedConfig), t)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal string with error: %v", err)
	}
	return t, nil
}

// Bytes will return the terraform configuration as a slice of bytes. Note this is NOT the terraform configuration.
func (t *TerraformConfiguration) Bytes() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("Unable to JSON marshal with error: %v", err)
	}
	return bytes, nil
}

// String will return the terraform configuration as a string. Note this is NOT the terraform configuration.
func (t *TerraformConfiguration) String() (string, error) {
	bytes, err := t.Bytes()
	if err != nil {
		return "", fmt.Errorf("Unable to get Bytes with error: %v", err)
	}
	return string(bytes), nil
}

func (t *TerraformConfiguration) Hash() (string, error) {
	hash, err := hashstructure.Hash(t, nil)
	if err != nil {
		return "", fmt.Errorf("Unable to hash terraform configuration: %v", err)
	}
	return fmt.Sprintf("%d", hash), nil
}

func (t *TerraformConfiguration) SetApplyHash(hash string) {
	t.applyHash = hash
}

func (t *TerraformConfiguration) GetApplyHash() string {
	return t.applyHash
}
