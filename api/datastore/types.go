package datastore

import (
	"sync"

	"github.com/jmccarty3/spotlight/api"
)

type DataStore interface {
	StoreDataPoints([]*api.DataPoint) error
	VerifyInitialData(uint) error
	StorePrediction(*api.DataPoint) error
	StoreBestPrice(*api.DataPoint) error
}

type Factory func(map[string]string) DataStore

var providersMutex sync.Mutex
var providers = make(map[string]Factory)

func RegisterDataStore(name string, factory Factory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()

	providers[name] = factory
}

func GetDataStore(name string) Factory {
	providersMutex.Lock()
	defer providersMutex.Unlock()

	return providers[name]
}
