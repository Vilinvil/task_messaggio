package usecases

import (
	"context"

	"github.com/Vilinvil/task_messaggio/internal/messageworker/messageworker/repository"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"

	"github.com/segmentio/kafka-go"
)

var _ BrokerMessage = (*repository.BrokerMessageKafka)(nil)

type BrokerMessage interface {
	FetchMessage() (*models.MessagePayload, error)
	CommitMessage(ctx context.Context, messages ...kafka.Message) error
}

type SleepingMessageWorker struct {
	logger *mylogger.MyLogger
}

func NewSleepingMessageWorker(logger *mylogger.MyLogger) *SleepingMessageWorker {
	return &SleepingMessageWorker{
		logger: logger,
	}
}

func (s *SleepingMessageWorker) JobMessages() error {
	for {

	}
}
