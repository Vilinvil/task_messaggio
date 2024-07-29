package repository

import (
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
)

type BrokerMessageKafka struct {
	logger *mylogger.MyLogger
}

func NewBrokerMessageKafka(logger *mylogger.MyLogger) *BrokerMessageKafka {
	return &BrokerMessageKafka{
		logger: logger,
	}
}

func (b *BrokerMessageKafka) FetchMessage() (*models.MessagePayload, error) {
	// TODO implement
	return nil, nil
}
