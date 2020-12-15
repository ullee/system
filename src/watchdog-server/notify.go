package main

import (
	. "constants"
	"custom-pkg/slack"
	"os"
)

func (w *Watchdog) setMessage(metricType int) {
	switch metricType {
	case METRIC_TYPE_DISK:
		w.sm.subject = "Disk 임계치 초과 알림"
		w.sm.title = "Disk 사용률이 임계치를 초과 하였습니다."
		break
	case METRIC_TYPE_CPU:
		w.sm.subject = "CPU 임계치 초과 알림"
		w.sm.title = "CPU 사용률이 임계치를 초과 하였습니다."
		break
	case METRIC_TYPE_MEMORY:
		w.sm.subject = "Memory 임계치 초과 알림"
		w.sm.title = "Memory 사용률이 임계치를 초과 하였습니다."
		break
	}
	if os.Getenv("APP_ENV") == "production" {
		w.sm.channel = SLACK_CHANNEL_INFRA
	} else if os.Getenv("APP_ENV") == "staging" {
		w.sm.channel = SLACK_CHANNEL_INFRA_PRIVATE
	} else {
		w.sm.channel = SLACK_CHANNEL_DEV
	}
	w.sm.username = "system"
}

func (w *Watchdog) notify() error {
	var err error
	for metricType, message := range w.sm.values {
		w.setMessage(metricType)
		attachment := slack.Attachment{}
		attachment.AddField(slack.Field{
			Title: w.sm.title,
			Value: message,
		})
		payload := slack.Payload{
			Text:        w.sm.subject,
			Username:    w.sm.username,
			Channel:     w.sm.channel,
			IconEmoji:   ":monkey_face:",
			Attachments: []slack.Attachment{attachment},
		}
		statusCode, err := slack.Send(payload)
		if err != nil || statusCode != 200 {
			return err
		}
	}
	return err
}
