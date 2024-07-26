package usecases

import (
	"context"
	"fmt"

	"github.com/Vilinvil/task_messaggio/internal/message/message/repository"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"

	"github.com/google/uuid"
)

var _ MessageRepository = (*repository.MessagePg)(nil)

type MessageRepository interface {
	AddMessage(ctx context.Context, preMessage *models.PreMessage) error
}

type MessageService struct {
	repository MessageRepository
	logger     *mylogger.MyLogger
}

func NewMessageService(repository MessageRepository, logger *mylogger.MyLogger) *MessageService {
	return &MessageService{
		repository: repository,
		logger:     logger,
	}
}

func (m *MessageService) AddMessage(ctx context.Context, value string) (uuid.UUID, error) {
	logger := m.logger.EnrichReqID(ctx)

	preMessage := models.NewPreMessage(value)

	err := preMessage.Validate()
	if err != nil {
		logger.Error(err)

		return uuid.UUID{}, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	err = m.repository.AddMessage(ctx, preMessage)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return preMessage.ID, nil
}
