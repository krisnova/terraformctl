package controller

import "github.com/kris-nova/kubicorn/apis/cluster"

// GenerateTerraformFileFromKubicornAPI will build a terraform configuration file from a kubicorn API
func GenerateTerraformConfigurationFileFromKubicornAPI(cluster *cluster.Cluster) ([]byte, error) {
	return []byte(""), nil
}

// GenerateTerraformVarsFileFromKubicornAPI will build a terraform vars file from a kubicorn API
func GenerateTerraformFileFromKubicornAPI(cluster *cluster.Cluster) ([]byte, error) {
	return []byte(""), nil
}
