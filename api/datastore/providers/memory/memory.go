package memory

import (
	"errors"

	"github.com/jmccarty3/spotlight/api"
	"github.com/jmccarty3/spotlight/api/datastore"
)

const StoreName = "memory"

// MemoryProvider Implememts DataStore
type MemoryProvider struct {
}

func newMemoryProvider(config map[string]string) datastore.DataStore {
	return &MemoryProvider{}
}

func init() {
	datastore.RegisterDataStore(StoreName, newMemoryProvider)
}

func (s *MemoryProvider) StoreDataPoints([]*api.DataPoint) error {
	return errors.New("NotImplemented")
}

func (s *MemoryProvider) VerifyInitialData(uint) error {
	return errors.New("NotImplemented")
}

func (s *MemoryProvider) StorePrediction(*api.DataPoint) error {
	return errors.New("NotImplemented")
}

func (s *MemoryProvider) StoreBestPrice(*api.DataPoint) error {
	return errors.New("NotImplemented")
}
