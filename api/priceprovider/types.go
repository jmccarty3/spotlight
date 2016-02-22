package priceprovider // import "github.com/jmccarty3/spotlight/api/priceprovider"

import (
	"sync"
	"time"

	"github.com/jmccarty3/spotlight/api"
)

// PriceProvider Interface to generate price data
type PriceProvider interface {
	GetPricesToDate(region string, zones, types []string, startDat time.Time) []*api.DataPoint
	GetPrices(region string, zones, types []string, startDate, endDate time.Time) []*api.DataPoint
}

// Factory is a function to returns a priceprovider.PriceProvider interface
type Factory func(map[string]string) PriceProvider

var providersMutex sync.Mutex
var providers = make(map[string]Factory)

// RegisterPriceProvider creates an instance of a named PriceProvider
func RegisterPriceProvider(name string, factory Factory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()

	//glog.Infof("Registering Price Provider: ", name)
	//TODO Check for double registration
	providers[name] = factory
}

func GetPriceProvider(name string) Factory {
	providersMutex.Lock()
	defer providersMutex.Unlock()

	return providers[name]
}
