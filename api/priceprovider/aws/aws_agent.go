package aws

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/glog"

	"github.com/jmccarty3/spotlight/api"
	"github.com/jmccarty3/spotlight/api/priceprovider"
)

// ProviderName for aws price provider
const ProviderName = "aws"

// AWSProvider implements PriceProvider
type AWSProvider struct {
	client *ec2.EC2
}

func init() {
	priceprovider.RegisterPriceProvider(ProviderName, func(config map[string]string) priceprovider.PriceProvider {
		creds := credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvProvider{},
				&ec2rolecreds.EC2RoleProvider{
					Client: ec2metadata.New(session.New(&aws.Config{})),
				},
				&credentials.SharedCredentialsProvider{},
			})
		return newAWSProvider(config, creds)
	})
}

func newAWSProvider(config map[string]string, creds *credentials.Credentials) *AWSProvider {
	client := &AWSProvider{
		client: ec2.New(session.New(), &aws.Config{Region: aws.String("us-east-1")}),
	}

	return client
}

func (p *AWSProvider) GetPricesToDate(region string, zones, types []string, startDate time.Time) []*api.DataPoint {
	glog.V(3).Infoln("Getting Prices to date", startDate)
	return p.GetPrices(region, zones, types, startDate, time.Now())
}

func (p *AWSProvider) GetPrices(region string, zones, types []string, startDate, endDate time.Time) []*api.DataPoint {
	var instanceTypes []*string

	for _, t := range types {
		instanceTypes = append(instanceTypes, aws.String(t))
	}

	var results []*api.DataPoint
	for _, z := range zones {
		results = append(results, p.getPriceByAz(region, z, instanceTypes, startDate, endDate)...)
	}

	return results
}

//TODO Return errors
func (p *AWSProvider) getPriceByAz(region, az string, types []*string, startTime, endTime time.Time) []*api.DataPoint {
	params := &ec2.DescribeSpotPriceHistoryInput{
		AvailabilityZone: aws.String(az),
		EndTime:          aws.Time(endTime),
		InstanceTypes:    types,
		StartTime:        aws.Time(startTime),
		MaxResults:       aws.Int64(1000),
	}

	var results []*api.DataPoint

	glog.V(3).Infoln("Requesting Spot Prices for ", region, az, types, startTime, endTime)
	pageCount := 0
	lastTime := time.Now()
	p.client.DescribeSpotPriceHistoryPages(params, func(page *ec2.DescribeSpotPriceHistoryOutput, lastPage bool) bool {
		pageResults := make([]*api.DataPoint, len(page.SpotPriceHistory))
		pageCount++
		glog.V(4).Infoln("Page Number:", pageCount, "Size:", len(page.SpotPriceHistory))
		for i, p := range page.SpotPriceHistory {
			pageResults[i] = convertToDataPoint(p, &region)
			lastTime = *p.Timestamp
		}
		glog.V(4).Infoln("Last Time for Page:", lastTime)
		results = append(results, pageResults...)
		return len(page.SpotPriceHistory) > 0 && startTime.Before(lastTime) //TODO Compare NextToken
	})

	glog.V(3).Infoln("Finished gathering prices. Last Time:", lastTime)
	return results
}

func convertToDataPoint(history *ec2.SpotPrice, region *string) *api.DataPoint {
	price, _ := strconv.ParseFloat(*history.SpotPrice, 32)
	return &api.DataPoint{
		Region:       *region,
		InstanceType: *history.InstanceType,
		Zone:         *history.AvailabilityZone,
		Price:        float32(price),
		Timestamp:    *history.Timestamp,
	}
}
