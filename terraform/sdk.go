package terraform

import (
	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/kris-nova/terraformctl/terraform/tfmain"
)

type TerraformCommand struct {
	arguments []string
}

func NewTerraformCommand(args []string) *TerraformCommand {
	return &TerraformCommand{
		arguments: args,
	}
}

func (t *TerraformCommand) Run() (int, error) {
	logger.Info("Running Terraform command [%s]", t.arguments[0])

	exitCode := tfmain.HackedMain(t.arguments)
	logger.Info("Terraform exit code [%d]", exitCode)

	return exitCode, nil
}
