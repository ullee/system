package main

import (
	. "constants"
	"custom-pkg/aws/ec2"
	. "custom-pkg/logger"
	"os"
)

const (
	SLACK_CHANNEL_DEV           = "slack channel dev"
	SLACK_CHANNEL_INFRA_PRIVATE = "slack channel staging"
	SLACK_CHANNEL_INFRA         = "slack channel production"
)

type Watchdog struct {
	metricType int
	instName   string
	sm         SlackMessage
	consts     Consts
}

type SlackMessage struct {
	subject  string
	username string
	channel  string
	title    string
	values   map[int]string
}

func main() {

	w := Watchdog{}
	instances, err := ec2.GetInstances()
	if err != nil {
		Log.Error(err)
		os.Exit(-1)
	}

	/* TODO RDS Access 권한 필요
	dbInstances, err := rds.GetDBInstances()
	if err != nil {
		Log.Error(err)
		os.Exit(-1)
	}
	*/

	w.sm.values = map[int]string{}
	for idx, _ := range instances.Reservations {
		for _, inst := range instances.Reservations[idx].Instances {
			if err = w.disk(*inst.InstanceId, inst.Tags); err != nil {
				Log.Error(err)
			}
			if err = w.cpu(*inst.InstanceId, inst.Tags); err != nil {
				Log.Error(err)
			}
			if err = w.memory(*inst.InstanceId, inst.Tags); err != nil {
				Log.Error(err)
			}
		}
	}

	if err := w.notify(); err != nil {
		Log.Error(err)
	}
}
