package blobPersist

import (
	"github.com/kris-nova/terraformctl/storage"
)

type BlobPersist struct {
}

// NewBlobPersist will return a new Azure blob storage implementation to be used for the persistent layer of terraformctl.
func NewBlobPersist() *BlobPersist {
	return &BlobPersist{}
}

// ConcurrentPersist is what will start a concurrent persisting process that will run and persist as new data comes into a defined cache.
func (c *BlobPersist) ConcurrentPersist(cacher storage.TerraformCtlCacher) (chan *storage.PersisterCancel, error) {
	cancelChan := make(chan *storage.PersisterCancel)
	go func() {

	}()
	return cancelChan, nil
}
