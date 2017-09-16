package controller

type TerraformRunner struct {
	configurationFile []byte
	variablesFile     []byte
}

func NewTerraformRunner(configurationFile, variablesFile []byte) *TerraformRunner {
	return &TerraformRunner{
		configurationFile: configurationFile,
		variablesFile:     variablesFile,
	}
}

func (t *TerraformRunner) Apply() error {
	return nil
}
