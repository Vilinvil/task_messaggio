package usecases

import (
	"context"

	"github.com/Vilinvil/task_messaggio/internal/message/message/repository"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/google/uuid"
)

var _ MessageRepository = (*repository.MessagePg)(nil)

type MessageRepository interface {
	GetMessageStatistic(ctx context.Context) (*models.MessageStatistic, error)
	AddMessage(ctx context.Context, preMessage *models.MessagePayload) error
}

var _ BrokerMessageRepository = (*repository.BrokerMessageKafka)(nil)

type BrokerMessageRepository interface {
	WriteMessage(ctx context.Context, msgPayload *models.MessagePayload) error
}

type MessageService struct {
	messageRepository       MessageRepository
	brokerMessageRepository BrokerMessageRepository
	generatorUUID           models.GeneratorUUID
	logger                  *mylogger.MyLogger
}

func NewMessageService(messageRepository MessageRepository,
	brokerMessageRepository BrokerMessageRepository, generatorUUID models.GeneratorUUID, logger *mylogger.MyLogger,
) *MessageService {
	return &MessageService{
		messageRepository:       messageRepository,
		brokerMessageRepository: brokerMessageRepository,
		generatorUUID:           generatorUUID,
		logger:                  logger,
	}
}

func (m *MessageService) AddMessage(ctx context.Context, value string) (uuid.UUID, error) {
	logger := m.logger.EnrichReqID(ctx)

	preMessage := models.NewMessagePayload(value, m.generatorUUID)

	err := preMessage.Validate()
	if err != nil {
		logger.Error(err)

		return uuid.UUID{}, err
	}

	err = m.messageRepository.AddMessage(ctx, preMessage)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = m.brokerMessageRepository.WriteMessage(ctx, preMessage)
	if err != nil {
		return uuid.UUID{}, err
	}

	return preMessage.ID, nil
}

func (m *MessageService) GetMessageStatistic(ctx context.Context) (*models.MessageStatistic, error) {
	return m.messageRepository.GetMessageStatistic(ctx)
}
