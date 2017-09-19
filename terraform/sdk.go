package terraform

import (
	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/kris-nova/terraformctl/terraform/tfmain"
)

// TerraformCommand represents a terraform command a user would type on the command line
type TerraformCommand struct {
	arguments []string
}

// NewTerraformCommand will initialize a Terraform command based on a slice of arguments that will map to a Terraform command.
// Example:
//   [
//      "terraform",
//      "init",
//      "/path/to/my/config/dir/"
//   ]
func NewTerraformCommand(args []string) *TerraformCommand {
	return &TerraformCommand{
		arguments: args,
	}
}

// Run will attempt to run the command and return an exit code that would be how Terraform would have exited.
func (t *TerraformCommand) Run() (int, error) {
	logger.Info("Running Terraform command [%s]", t.arguments[1])

	exitCode := tfmain.HackedMain(t.arguments)
	logger.Info("Terraform exit code [%d]", exitCode)

	return exitCode, nil
}
