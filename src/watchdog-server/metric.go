package main

import (
	. "constants"
	"custom-pkg/aws/cw"
	. "custom-pkg/logger"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	EC2 "github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"strconv"
	"strings"
)

func (w *Watchdog) disk(instanceId string, tags []*EC2.Tag) error {
	w.metricType = METRIC_TYPE_DISK
	resp, err := cw.GetMetric("Linux/System", "DiskUtilization", instanceId)
	if err != nil {
		return err
	}

	if err = w.check(resp.MetricDataResults, tags, METRIC_TYPE_DISK); err != nil {
		return err
	}
	return err
}

func (w *Watchdog) memory(instanceId string, tags []*EC2.Tag) error {
	w.metricType = METRIC_TYPE_MEMORY
	resp, err := cw.GetMetric("Linux/System", "MemoryUtilization", instanceId)
	if err != nil {
		return err
	}

	if err = w.check(resp.MetricDataResults, tags, METRIC_TYPE_MEMORY); err != nil {
		return err
	}
	return err
}

func (w *Watchdog) cpu(instanceId string, tags []*EC2.Tag) error {
	w.metricType = METRIC_TYPE_CPU
	resp, err := cw.GetMetric("AWS/EC2", "CPUUtilization", instanceId)
	if err != nil {
		return err
	}

	if err = w.check(resp.MetricDataResults, tags, METRIC_TYPE_CPU); err != nil {
		return err
	}
	return err
}

func (w *Watchdog) check(metricDataResults []*cloudwatch.MetricDataResult, tags []*EC2.Tag, metricType int) error {
	var err error
	for _, val := range tags {
		if *val.Key == "Name" {
			w.instName = *val.Value
			break
		}
	}

	// 서비스별 임계치 설정
	var criticalValue float64

	if os.Getenv("APP_ENV") == "staging" || os.Getenv("APP_ENV") == "dev" {
		if strings.Contains(w.instName, "staging") == false {
			return err
		}
		criticalValue = DEFAULT_STAGING_CRITICAL_VALUE
	} else {
		criticalValue = DEFAULT_PRODUCTION_CRITICAL_VALUE
	}

	switch metricType {
	case METRIC_TYPE_DISK:
		if strings.Contains(w.instName, TAG_ES) == true && strings.Contains(w.instName, "services") == false {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_ES]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_ES]
			}
		} else if strings.Contains(w.instName, TAG_DEVOPS) {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_DEVOPS]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_DEVOPS]
			}
		} else if strings.Contains(w.instName, TAG_CACHE) == true {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_CACHE]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_CACHE]
			}
		}
		break

	case METRIC_TYPE_CPU:
		if strings.Contains(w.instName, TAG_SOCKET) == true {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_SOCKET]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_SOCKET]
			}
		} else if strings.Contains(w.instName, TAG_DEVOPS) == true {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_DEVOPS]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_DEVOPS]
			}
		} else if strings.Contains(w.instName, TAG_ES) == true && strings.Contains(w.instName, "services") == false {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_ES]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_ES]
			}
		} else if strings.Contains(w.instName, TAG_CACHE) == true {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_CACHE]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_CACHE]
			}
		}
		break

	case METRIC_TYPE_MEMORY:
		if strings.Contains(w.instName, TAG_DEVOPS) {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_DEVOPS]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_DEVOPS]
			}
		} else if strings.Contains(w.instName, TAG_ES) == true && strings.Contains(w.instName, "services") == false {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_ES]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_ES]
			}
		} else if strings.Contains(w.instName, TAG_CACHE) == true {
			if _, ok := w.consts.CriticalValues()[metricType][TAG_CACHE]; ok {
				criticalValue = w.consts.CriticalValues()[metricType][TAG_CACHE]
			}
		}
		break
	default:
		break
	}

	for _, metricData := range metricDataResults {
		if metricData.Timestamps != nil && metricData.Values != nil {
			var count float64
			var total float64
			for index, _ := range metricData.Timestamps {
				if *metricData.StatusCode == "Complete" {
					count++
					total = total + *metricData.Values[index]
					//Log.Debug(w.instName, *metricData.Timestamps[index], *metricData.Values[index])
				}
			}
			average := total / count
			if average > criticalValue {
				averageStr := strconv.FormatFloat(average, 'f', 0, 64)
				criticalValueStr := strconv.FormatFloat(criticalValue, 'f', 0, 64)
				w.sm.values[w.metricType] = w.sm.values[w.metricType] + w.instName + " " + averageStr + "% MAX(" + criticalValueStr + "%)\n"
				Log.Warning(w.instName, *metricData.Label, averageStr+"%", "MAX("+criticalValueStr+"%)")
			}
		}
	}
	return err
}
