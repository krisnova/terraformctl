package controller

import (
	"fmt"
	"github.com/kris-nova/kubicorn/cutil/hang"
	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/kris-nova/terraformctl/storage"
)

// TerraformControlLoop is the Kubernetes controller that will ultimately be ran in a container. The controller
// is responsible for mutating infrastructure according to whatever definitions are currently set in the configured
// storage layer.
type TerraformControlLoop struct {
	optionsChannel chan *TerraformControlLoopOptions
	options        *TerraformControlLoopOptions
	cacher         storage.TerraformCtlCacher
	persister      storage.TerraformCtlPersister
}

// TerraformControlLoopOptions is a data structure that is used to configure, and update the control loop.
// Change these options and pass a new struct into UpdateOptions to change the control loops behavior.
type TerraformControlLoopOptions struct {

	// Stop is a bool that can be set to 1 to tell the control loop to stop running.
	// Use this bit to perform a clean exit of the control loop.
	Stop bool
}

// NewTerraformControlLoop will initialize a new TerraformCtlController and return it.
func NewTerraformControlLoop(options *TerraformControlLoopOptions, cacher storage.TerraformCtlCacher, persister storage.TerraformCtlPersister) *TerraformControlLoop {
	t := &TerraformControlLoop{
		options:        options,
		optionsChannel: make(chan *TerraformControlLoopOptions),
		cacher:         cacher,
		persister:      persister,
	}
	return t
}

// UpdateOptions will accept a new TerraformControlLoopOptions struct and will attempt to send it
// over a channel to the control loop. Note that the function will hang until the control loop has
// read the new options into memory.
func (t *TerraformControlLoop) UpdateOptions(options *TerraformControlLoopOptions) {
	t.options = options
	t.optionsChannel <- options
}

// Run is the method that will start running the infrastructure controller concurrently.
func (t *TerraformControlLoop) Run() chan error {
	errorChan := make(chan error)

	// Initialize persistence layer
	persistentCancelChan, err := t.persister.ConcurrentPersist(t.cacher)
	if err != nil {
		errorChan <- fmt.Errorf("Unable to run concurrent persister with error: %v", err)
		return errorChan
	}

	// Build hanger for the control loop
	hg := hang.Hanger{
		Ratio: 1,
	}

	go func() {
		// The main loop for the control loop. We will loop while options.Stop is set to false.
		for t.options.Stop == false {

			// List all known clusters
			configurations, err := t.cacher.List()
			if err != nil {
				errorChan <- fmt.Errorf("Unable to list clusters: %v", err)
				hg.Hang()
			}

			// Reconcile each cluster with terraform
			for _, config := range configurations {
				terraformRunner := NewTerraformRunner(config)

				// Check if exists
				runApply := false

				if config.GetApplyHash() == "" {
					logger.Info("New configuration [%s]", config.Name)
					runApply = true
				}

				//logger.Info("Found existing configuration [%s]", config.Name)
				newHash, err := config.Hash()
				if err != nil {
					errorChan <- fmt.Errorf("Unable to calculate hash for configuration [%s] with error: %v", config.Name, err)
					hg.Hang()
					continue
				}
				if config.GetApplyHash() != newHash {
					logger.Info("Delta in configuration hash [%s] [%s]", config.GetApplyHash(), newHash)
					runApply = true
				}

				if runApply {
					err = terraformRunner.Apply()
					if err != nil {
						errorChan <- fmt.Errorf("Unable to run terraform apply for configuration [%s] with error: %v", config.Name, err)
						hg.Hang()
						continue
					}
					applyHash, err := config.Hash()
					if err != nil {
						errorChan <- fmt.Errorf("Unable to calculate apply hash with error: %v", err)
					}

					// Update the hash and save
					config.SetApplyHash(applyHash)
					t.cacher.Lock()
					t.cacher.Set(config.Name, config)
					t.cacher.Unlock()
					logger.Info("Saved configuration [%s]", config.Name)
				}
			}
			// Reset ratio
			hg.Ratio = 1
		}
		// Lock the cache
		t.cacher.Lock()

		// Cancel the persistence layer
		persistentCancelChan <- &storage.PersisterCancel{}

		errorChan <- nil

	}()
	return errorChan
}
