package memoryCache

import (
	"github.com/kris-nova/kubicorn/apis/cluster"
	"github.com/kris-nova/terraformctl/storage"
)

type MemoryCache struct {
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{}
}

// Get is used to attempt to lookup a cluster by a name, if the cluster is not found Get will error.
func (m *MemoryCache) Get(name string) (*cluster.Cluster, error) {
	return &cluster.Cluster{}, nil
}

// Set is used to set a cluster by name. Upsert logic is respected and the implementation will update or create as necessary.
func (m *MemoryCache) Set(name string, cluster *cluster.Cluster) error {
	return nil
}

// List will return all known clusters in the cache.
func (m *MemoryCache) List() ([]*cluster.Cluster, error) {
	var list []*cluster.Cluster
	return list, nil
}

// Lock will lock the caching layer so no other writes may be made until the layer is unlocked.
func (m *MemoryCache) Lock() {

}

// Unlock will unlock the caching layer so writes can continue to be made.
func (m *MemoryCache) Unlock() {

}

// Synchronize will attempt to synchronize a cache with an existing persistency layer that might already have cluster information in it.
func (m *MemoryCache) Synchronize(persister storage.TerraformCtlPersister) error {
	return nil
}
