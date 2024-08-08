package server

import (
	"context"
	"strings"
	"time"

	"github.com/Vilinvil/task_messaggio/internal/messageworker/config"
	"github.com/Vilinvil/task_messaggio/internal/messageworker/messageworker/repository"
	"github.com/Vilinvil/task_messaggio/internal/messageworker/messageworker/usecases"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	pkgrepository "github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
)

var _ MessageWorker = (*usecases.MessageWorker)(nil)

type MessageWorker interface {
	JobMessages(ctx context.Context, handlerFunc usecases.MessagePayloadHandler) <-chan error
}

func NewSleepingHandler(durationSleep time.Duration) func(ctx context.Context, payload *models.MessagePayload) error {
	return func(ctx context.Context, _ *models.MessagePayload) error {
		timer := time.NewTimer(durationSleep)

		select {
		case <-timer.C:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type Server struct{}

func (s *Server) Run(ctx context.Context, config *config.Config) error {
	logger, err := mylogger.New(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "), config.ProductionMode)
	if err != nil {
		return err
	}

	defer logger.Sync() //nolint:errcheck

	logger.Debugln(config)

	brokerMessageKafka, err := repository.NewBrokerMessageKafka(config.BrokerAddr, logger)
	if err != nil {
		return err
	}

	pool, err := pkgrepository.NewPgxPool(ctx, config.URLDataBase)
	if err != nil {
		logger.Error(err)

		return err
	}

	messageRepository := repository.NewMessagePg(pool, logger)

	messageWorker := usecases.NewMessageWorker(brokerMessageKafka, messageRepository, logger)

	sleepingHandler := NewSleepingHandler(config.TimeOnTask)

	chErr := messageWorker.JobMessages(ctx, sleepingHandler)

	logger.Infof("ServerWorker start job messages")

	return <-chErr
}
