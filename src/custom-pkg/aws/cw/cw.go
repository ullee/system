package cw

import (
	. "constants"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var count int64 = 0

func getSession() (*cloudwatch.CloudWatch, error) {
	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(CW_REGION),
			Credentials: credentials.NewStaticCredentials(CW_ACCESS_KEY, CW_SECRET_KEY, ""),
		},
	})

	if err != nil {
		return nil, err
	}

	return cloudwatch.New(s), err
}

func GetDimensions(metadata map[string]string) (ret []*cloudwatch.Dimension) {

	var _ret []*cloudwatch.Dimension

	instanceIdName := "InstanceId"
	instanceIdValue, ok := metadata["instanceId"]
	if ok {
		dim := cloudwatch.Dimension{
			Name:  aws.String(instanceIdName),
			Value: aws.String(instanceIdValue),
		}
		_ret = append(_ret, &dim)
	}

	return _ret
}

func AddMetric(name, unit string, value float64, dimensions []*cloudwatch.Dimension, metricData []*cloudwatch.MetricDatum) (ret []*cloudwatch.MetricDatum, err error) {
	_metric := cloudwatch.MetricDatum{
		MetricName: aws.String(name),
		Unit:       aws.String(unit),
		Value:      aws.Float64(value),
		Dimensions: dimensions,
	}
	metricData = append(metricData, &_metric)
	return metricData, nil
}

func GetMetric(nameSpace, metricName, instanceId string) (*cloudwatch.GetMetricDataOutput, error) {

	cw, err := getSession()
	if err != nil {
		return nil, err
	}

	dataInput := &cloudwatch.GetMetricDataInput{
		StartTime: aws.Time(time.Now().Add(-5 * time.Minute)),
		EndTime:   aws.Time(time.Now()),
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			buildMetricDataQuery(nameSpace, metricName, instanceId),
		},
	}
	dataOutput, err := cw.GetMetricData(dataInput)
	if err != nil {
		return nil, err
	}
	return dataOutput, err
}

func buildMetricDataQuery(nameSpace, metricName, instanceId string) *cloudwatch.MetricDataQuery {
	count++
	return &cloudwatch.MetricDataQuery{
		Id: aws.String("id" + strconv.FormatInt(count, 10)),
		MetricStat: &cloudwatch.MetricStat{
			Period: aws.Int64(60),
			Stat:   aws.String("Average"),
			Metric: &cloudwatch.Metric{
				MetricName: aws.String(metricName),
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("InstanceId"),
						Value: aws.String(instanceId),
					},
				},
				Namespace: aws.String(nameSpace),
			},
		},
	}
}

func PutMetric(metricData []*cloudwatch.MetricDatum, namespace string) error {

	cw, err := getSession()
	if err != nil {
		return err
	}

	metricInput := &cloudwatch.PutMetricDataInput{
		MetricData: metricData,
		Namespace:  aws.String(namespace),
	}

	resp, err := cw.PutMetricData(metricInput)
	if err != nil {
		return err
	}
	fmt.Println(awsutil.StringValue(resp))
	return nil
}

/*
  Metadata struct:
  {
    "devpayProductCodes" : null,
	"privateIp" : "10.0.5.89",
	"availabilityZone" : "us-west-1a",
	"version" : "2010-08-31",
	"region" : "us-west-1",
	"instanceId" : "i-e0iag2b",
	"billingProducts" : null,
	"accountId" : "208372078340",
	"instanceType" : "m3.xlarge",
	"imageId" : "ami-43f91b07",
	"kernelId" : null,
    "ramdiskId" : null,
    "pendingTime" : "2015-06-30T08:28:48Z",
    "architecture" : "x86_64"
  }
*/
func GetInstanceMetadata() (metadata map[string]string, err error) {
	var data map[string]string
	resp, err := http.Get("http://169.254.169.254/latest/dynamic/instance-identity/document")
	if err != nil {
		return data, fmt.Errorf("can't reach metadata endpoint - %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, fmt.Errorf("can't read metadata response body - %s", err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, fmt.Errorf("can't json parsing metadata response body - %s", err)
	}

	return data, err
}
