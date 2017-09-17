package storage

import "github.com/kris-nova/terraformctl/parser"

// TerraformCtlCacher defines the methods used to interact with the caching layer of terraformctl.
type TerraformCtlCacher interface {

	// Get is used to attempt to lookup a configuration by a name, if the configuration is not found Get will error.
	Get(name string) (*parser.TerraformConfiguration, error)

	// Set is used to set a configuration by name. Upsert logic is respected and the implementation will update or create as necessary.
	Set(name string, config *parser.TerraformConfiguration) error

	// List will return all known configurations in the cache.
	List() ([]*parser.TerraformConfiguration, error)

	// Lock will lock the caching layer so no other writes may be made until the layer is unlocked.
	Lock()

	// Unlock will unlock the caching layer so writes can continue to be made.
	Unlock()

	// Synchronize will attempt to synchronize a cache with an existing persistency layer that might already have configuration information in it.
	Synchronize(persister TerraformCtlPersister) error
}

// PersisterCancel is used to send over the PersisterCancel channel to tell terraformctl to cancel persisting the cache.
type PersisterCancel struct {
}

// TerraformCtlCacher defines the methods used to interact with the persistent layer of terraformctl.
type TerraformCtlPersister interface {

	// ConcurrentPersist is what will start a concurrent persisting process that will run and persist as new data comes into a defined cache.
	ConcurrentPersist(cacher TerraformCtlCacher) (chan *PersisterCancel, error)

	// List will return all known configurations in the cache.
	List() ([]*parser.TerraformConfiguration, error)

}

