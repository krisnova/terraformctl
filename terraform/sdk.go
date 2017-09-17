package terraform

import "github.com/kris-nova/kubicorn/cutil/logger"

type TerraformCommand struct {
	arguments []string
}

func NewTerraformCommand(args []string) *TerraformCommand {
	return &TerraformCommand{
		arguments: args,
	}
}

func (t *TerraformCommand) Run() error {
	logger.Info("Running Terraform command [%s]", t.arguments[0])
	return nil
}
