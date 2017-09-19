package controller

import (
	"fmt"
	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/kris-nova/terraformctl/parser"
	"github.com/kris-nova/terraformctl/terraform"
	"io/ioutil"
	"os"
)

// TerraformRunner maps a terraformctl TerraformConfiguration to the TerraformSDK
type TerraformRunner struct {
	configuration *parser.TerraformConfiguration
}

// NewTerraformRunner initializes a new TerraformRunner struct with specified configuration.
func NewTerraformRunner(configuration *parser.TerraformConfiguration) *TerraformRunner {
	return &TerraformRunner{
		configuration: configuration,
	}
}

// Apply will use the TerraformSDK to run a set of procedural terraform commands on the loaded TerraformConfiguration
// in an attempt to reconcile the infrastructure.
func (t *TerraformRunner) Reconcile() error {
	logger.Info("Calling Terraform apply on configuration [%s]", t.configuration.Name)

	// Init a temp directory with our terraform configuration.
	fpath := os.TempDir() + "/main.tf"
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("Unable to write temp file: %v", err)
	}

	// Write the main.tf file
	bytes, err := t.configuration.Bytes()
	if err != nil {
		return fmt.Errorf("Unable to read bytes for Terraform configuration: %v", err)
	}
	ioutil.WriteFile(fpath, bytes, 0664)
	f.Close()

	// -----------------------------------------------------------------------------------------------------------------
	// Init
	// ----
	exitCode, err := terraform.NewTerraformCommand([]string{
		"init", // Subcommand
		fpath,  // Directory
	}).Run()

	logger.Debug("Exit code: %d", exitCode)

	return nil
}
