package notification

import (
	log "github.com/sirupsen/logrus"
)

type loggingHandler struct {
}

func (p *loggingHandler) Handle(pageInfo PageInfo) {
	log.WithFields(log.Fields{
		"name":   pageInfo.Name(),
		"url":    pageInfo.URL(),
		"status": pageInfo.Status(),
	}).Info("Send notification")
}

func NewLoggingHandler() Handler {
	return &loggingHandler{}
}
