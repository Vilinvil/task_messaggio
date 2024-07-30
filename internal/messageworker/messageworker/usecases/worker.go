package usecases

import (
	"context"
	"sync"

	"github.com/Vilinvil/task_messaggio/internal/messageworker/messageworker/repository"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/google/uuid"
)

var _ BrokerMessage = (*repository.BrokerMessageKafka)(nil)

type BrokerMessage interface {
	StartConsumption(ctx context.Context) (
		chConsumptionMessages <-chan repository.MessagePayloadWithCommitFunc,
		chErrConsumptionMessage <-chan error,
	)
}

var _ MessageRepository = (*repository.MessagePg)(nil)

type MessageRepository interface {
	SetStatusMessage(ctx context.Context, messageID *uuid.UUID, status repository.StatusMessage) error
}

type MessageWorker struct {
	messageRepository MessageRepository
	brokerMessage     BrokerMessage
	logger            *mylogger.MyLogger
	chErr             chan error
	onceChErr         *sync.Once
}

func NewMessageWorker(brokerMessage BrokerMessage, messageRepository MessageRepository,
	logger *mylogger.MyLogger,
) *MessageWorker {
	return &MessageWorker{ //nolint:exhaustruct
		brokerMessage:     brokerMessage,
		messageRepository: messageRepository,
		logger:            logger,
		onceChErr:         &sync.Once{},
	}
}

// JobMessages one time starts job with messages. At the next call of this method only returning
// existing channel happens.
func (s *MessageWorker) JobMessages(ctx context.Context,
	handlerFunc func(ctx context.Context, payload *models.MessagePayload) error,
) <-chan error {
	if s.chErr != nil {
		return s.chErr
	}

	s.onceChErr.Do(func() {
		go func() {
			s.chErr = make(chan error)

			chConsumptionMessages, chErrConsumptionMessage := s.brokerMessage.StartConsumption(ctx)

			for msgWithCommitFunc := range chConsumptionMessages {
				err := handlerFunc(ctx, &msgWithCommitFunc.MessagePayload)
				if err != nil {
					s.chErr <- err

					return
				}

				err = s.messageRepository.SetStatusMessage(
					ctx, &msgWithCommitFunc.MessagePayload.ID, repository.StatusMessageDone)
				if err != nil {
					s.chErr <- err

					return
				}

				err = msgWithCommitFunc.CommitFunc()
				if err != nil {
					s.chErr <- err

					return
				}
			}

			err := <-chErrConsumptionMessage
			if err != nil {
				s.chErr <- err

				return
			}

			s.chErr <- nil
		}()
	})

	return s.chErr
}
