package controller

import (
	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/kris-nova/terraformctl/parser"
)

type TerraformRunner struct {
	configuration *parser.TerraformConfiguration
}

func NewTerraformRunner(configuration *parser.TerraformConfiguration) *TerraformRunner {
	return &TerraformRunner{
		configuration: configuration,
	}
}

func (t *TerraformRunner) Apply() error {
	logger.Info("Calling Terraform apply on configuration [%s]", t.configuration.Name)
	return nil
}
