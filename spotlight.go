package main // import "github.com/jmccarty3/spotlight"

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/jmccarty3/spotlight/api/datastore"
	_ "github.com/jmccarty3/spotlight/api/datastore/providers"
	"github.com/jmccarty3/spotlight/api/priceprovider"
	_ "github.com/jmccarty3/spotlight/api/priceprovider/aws"
)

var (
	argInfluxHost   = flag.String("influxdb-host", "", "Host address of InfluxDB")
	argInfluxPort   = flag.Uint("influxdb-port", 8086, "InfluxDB Port")
	argAwsRegion    = flag.String("region", "us-east-1", "AWS Region")
	argDataBackstop = flag.Uint("data-backstop", 3, "Minimum amount of days worth of data required")
)

func main() {
	flag.Parse()
	glog.Info("Starting up")

	factory := priceprovider.GetPriceProvider("aws")
	var config map[string]string
	fmt.Println("Got Factory")
	aws := factory(config)

	date := time.Now().AddDate(0, 0, -1)
	fmt.Println("Calling Get Prices", date)
	vals := aws.GetPricesToDate("us-east-1", []string{"us-east-1c", "us-east-1a"}, []string{"m3.2xlarge", "c4.large", "c4.2xlarge", "m4.2xlarge"}, date)

	store := datastore.GetDataStore("influxdb")(config)
	err := store.StoreDataPoints(vals)

	fmt.Println("Data Added to Store")

	if err != nil {
		glog.Errorln("Error storing data:", err)
	}
}
