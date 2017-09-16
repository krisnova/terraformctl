package storage

import "github.com/kris-nova/kubicorn/apis/cluster"

// TerraformCtlCacher defines the methods used to interact with the caching layer of terraformctl.
type TerraformCtlCacher interface {

	// Get is used to attempt to lookup a cluster by a name, if the cluster is not found Get will error.
	Get(name string) (*cluster.Cluster, error)

	// Set is used to set a cluster by name. Upsert logic is respected and the implementation will update or create as necessary.
	Set(name string, cluster *cluster.Cluster) error

	// List will return all known clusters in the cache.
	List() ([]*cluster.Cluster, error)

	// Lock will lock the caching layer so no other writes may be made until the layer is unlocked.
	Lock()

	// Unlock will unlock the caching layer so writes can continue to be made.
	Unlock()

	// Synchronize will attempt to synchronize a cache with an existing persistency layer that might already have cluster information in it.
	Synchronize(persister TerraformCtlPersister) error
}

// PersisterCancel is used to send over the PersisterCancel channel to tell terraformctl to cancel persisting the cache.
type PersisterCancel struct {
}

// TerraformCtlCacher defines the methods used to interact with the persistent layer of terraformctl.
type TerraformCtlPersister interface {

	// ConcurrentPersist is what will start a concurrent persisting process that will run and persist as new data comes into a defined cache.
	ConcurrentPersist(cacher TerraformCtlCacher) (chan *PersisterCancel, error)
}
