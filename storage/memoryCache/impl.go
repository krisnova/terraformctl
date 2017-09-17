package memoryCache

import (
	"fmt"
	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/kris-nova/terraformctl/parser"
	"github.com/kris-nova/terraformctl/storage"
	"sync"
)

type MemoryCache struct {
	configurations map[string]*parser.TerraformConfiguration
	mutex          sync.Mutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		configurations: make(map[string]*parser.TerraformConfiguration),
		mutex:          sync.Mutex{},
	}
}

// Get is used to attempt to lookup a configuration by a name, if the configuration is not found Get will error.
func (m *MemoryCache) Get(name string) (*parser.TerraformConfiguration, error) {
	if config, ok := m.configurations[name]; ok {
		return config, nil
	}
	return nil, fmt.Errorf("Configuration [%s] not found", name)
}

// Set is used to set a configuration by name. Upsert logic is respected and the implementation will update or create as necessary.
func (m *MemoryCache) Set(name string, config *parser.TerraformConfiguration) error {
	if existing, ok := m.configurations[name]; ok {
		// Always pass the hash through
		config.SetApplyHash(existing.GetApplyHash())
	}
	m.configurations[name] = config
	return nil
}

// List will return all known configurations in the cache.
func (m *MemoryCache) List() ([]*parser.TerraformConfiguration, error) {
	var list []*parser.TerraformConfiguration
	for _, c := range m.configurations {
		list = append(list, c)
	}
	return list, nil
}

// Lock will lock the caching layer so no other writes may be made until the layer is unlocked.
func (m *MemoryCache) Lock() {
	m.mutex.Lock()
}

// Unlock will unlock the caching layer so writes can continue to be made.
func (m *MemoryCache) Unlock() {
	m.mutex.Unlock()
}

// Synchronize will attempt to synchronize a cache with an existing persistency layer that might already have configuration information in it.
func (m *MemoryCache) Synchronize(persister storage.TerraformCtlPersister) error {
	logger.Info("Synchronizing with persistent storage")
	list, err := persister.List()
	if err != nil {
		return fmt.Errorf("Unable to list configurations from persistence layer while synchronizing: %v", err)
	}
	for _, config := range list {
		logger.Debug("Synchronizing configuration [%s]", config.Name)
		m.configurations[config.Name] = config
	}
	return nil
}
