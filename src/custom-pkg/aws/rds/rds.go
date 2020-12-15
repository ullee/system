package rds

import (
	. "constants"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"os"
)

func getSession() (*rds.RDS, error) {
	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(EC2_REGION),
			Credentials: credentials.NewStaticCredentials(S3_ACCESS_KEY, S3_SECRET_KEY, ""),
		},
	})

	if err != nil {
		return nil, err
	}

	return rds.New(s), err
}

func GetDBInstances() (*rds.DescribeDBInstancesOutput, error) {
	db, err := getSession()
	if err != nil {
		return nil, err
	}

	var params *rds.DescribeDBInstancesInput
	if os.Getenv("APP_ENV") == "production" {
		params = &rds.DescribeDBInstancesInput{
			Filters: []*rds.Filter{
				{
					Name:   aws.String("instance-state-name"),
					Values: []*string{aws.String("running")},
				},
			},
		}
	} else {
		params = &rds.DescribeDBInstancesInput{
			Filters: []*rds.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: []*string{aws.String("*staging*")},
				},
				{
					Name:   aws.String("instance-state-name"),
					Values: []*string{aws.String("running")},
				},
			},
		}
	}
	result, err := db.DescribeDBInstances(params)
	if err != nil {
		return nil, err
	}
	return result, err
}
