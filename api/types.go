package api // import "github.com/jmccarty3/spotlight/api"

import "time"

type Query struct {
	region, instance_type string
	zones                 []string
	max_price             float32
}

type DataPoint struct {
	Region       string
	InstanceType string
	Zone         string
	Price        float32
	Timestamp    time.Time
}

type Analyzer interface {
	GetPrice(Query)
	GetBestPrice([]Query)
}
