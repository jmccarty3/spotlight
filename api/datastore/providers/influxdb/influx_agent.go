package influx

import (
	"errors"
	"fmt"
	"os"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/jmccarty3/spotlight/api"
	"github.com/jmccarty3/spotlight/api/datastore"
)

const DataStoreName = "influxdb"
const DefaultUserName = ""
const DefaultPassword = ""
const DefaultSpotPriceDB = "spotPrices"
const DefaultPredictionDB = "spotPredictions"
const DefaultBestPriceDB = "bestSpotPrediction"

type databaseSettings struct {
	spotPriceDB  string
	predictionDB string
	bestPriceDB  string
	precision    string
}

// InfluxStore implements DataStore
type InfluxStore struct {
	databaseSettings
	client client.Client
}

func getDefaultSettings() databaseSettings {
	return databaseSettings{
		spotPriceDB:  DefaultSpotPriceDB,
		predictionDB: DefaultPredictionDB,
		bestPriceDB:  DefaultBestPriceDB,
		precision:    "s",
	}
}

func getHostname(config map[string]string) string {
	host := config["influx_host"]

	//TODO Make this not terrible
	if host == "" {
		host = os.Getenv("INFLUXDB_HOST")
	}

	if host == "" {
		host = "localhost"
	}

	return host
}

func getPort(config map[string]string) string {
	return "8086"
}

func ensureDatabaseExists(c client.Client, database string) error {
	query := fmt.Sprintf("CREATE DATABASE %s", database)
	_, err := c.Query(client.NewQuery(query, database, "s"))
	return err
}

func newInfluxStore(config map[string]string) datastore.DataStore {
	client, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s:%s", getHostname(config), getPort(config)),
		Username: DefaultUserName,
		Password: DefaultPassword,
	})

	settings := getDefaultSettings()

	ensureDatabaseExists(client, settings.bestPriceDB)
	ensureDatabaseExists(client, settings.predictionDB)
	ensureDatabaseExists(client, settings.spotPriceDB)

	return &InfluxStore{
		databaseSettings: settings,
		client:           client,
	}
}

//Register with DataProvider
func init() {
	datastore.RegisterDataStore(DataStoreName, newInfluxStore)
}

func (s *InfluxStore) StoreDataPoints(points []*api.DataPoint) error {
	return s.writeData(s.spotPriceDB, "spot_prices", points...)
}

func (s *InfluxStore) VerifyInitialData(uint) error {
	return errors.New("NotImplemented")
}

func (s *InfluxStore) StorePrediction(prediction *api.DataPoint) error {
	return s.writeData(s.predictionDB, "predictions", prediction)
}

func (s *InfluxStore) StoreBestPrice(price *api.DataPoint) error {
	return s.writeData(s.bestPriceDB, "best_prices", price)
}

func (s *InfluxStore) createBatchPoints(db string) (client.BatchPoints, error) {
	return client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: s.precision,
	})
}

func convertDataPoint(name string, point *api.DataPoint) (*client.Point, error) {
	tags := map[string]string{
		"region":        point.Region,
		"az":            point.Zone,
		"instance-type": point.InstanceType,
	}

	fields := map[string]interface{}{"price": point.Price}
	//fmt.Println("Converted Point:", point.Timestamp)
	return client.NewPoint(name, tags, fields, point.Timestamp)
}

func (s *InfluxStore) writeData(db, metric string, points ...*api.DataPoint) error {
	bp, err := s.createBatchPoints(db)

	if err != nil {
		return err
	}

	for _, point := range points {
		if p, err := convertDataPoint(metric, point); err == nil {
			bp.AddPoint(p)
		} else {
			return err
		}
	}

	return s.client.Write(bp)
}
