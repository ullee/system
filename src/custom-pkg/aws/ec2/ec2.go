package ec2

import (
	. "constants"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
)

func getSession() (*ec2.EC2, error) {
	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(EC2_REGION),
			Credentials: credentials.NewStaticCredentials(EC2_ACCESS_KEY, EC2_SECRET_KEY, ""),
		},
	})

	if err != nil {
		return nil, err
	}

	return ec2.New(s), err
}

func GetInstances() (*ec2.DescribeInstancesOutput, error) {
	ec, err := getSession()
	if err != nil {
		return nil, err
	}

	var params *ec2.DescribeInstancesInput
	if os.Getenv("APP_ENV") == "production" {
		params = &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Environment"),
					Values: []*string{aws.String("production")},
				},
				{
					Name:   aws.String("instance-state-name"),
					Values: []*string{aws.String("running")},
				},
			},
		}
	} else {
		params = &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Environment"),
					Values: []*string{aws.String("staging")},
				},
				{
					Name:   aws.String("instance-state-name"),
					Values: []*string{aws.String("running")},
				},
			},
		}
	}
	result, err := ec.DescribeInstances(params)
	if err != nil {
		return nil, err
	}
	return result, err
}
