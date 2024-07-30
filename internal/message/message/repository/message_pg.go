package repository

import (
	"context"

	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessagePg struct {
	pool   *pgxpool.Pool
	logger *mylogger.MyLogger
}

func NewMessagePg(ctx context.Context, urlDataBase string, logger *mylogger.MyLogger) (*MessagePg, error) {
	pool, err := repository.NewPgxPool(ctx, urlDataBase)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	return &MessagePg{
		pool:   pool,
		logger: logger,
	}, nil
}

func (m *MessagePg) insertMessage(ctx context.Context, tx pgx.Tx, preMessage *models.MessagePayload) error {
	logger := m.logger.EnrichReqID(ctx)

	SQLInsertMessage := `INSERT INTO public."message" (id, value) VALUES ($1, $2)`

	_, err := tx.Exec(ctx, SQLInsertMessage, preMessage.ID, preMessage.Value)
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (m *MessagePg) addCounterMessage(ctx context.Context, tx pgx.Tx, delta int) error {
	logger := m.logger.EnrichReqID(ctx)

	SQLAddCounterMessage := `UPDATE public."counter_message" SET total = total + $1`

	_, err := tx.Exec(ctx, SQLAddCounterMessage, delta)
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (m *MessagePg) AddMessage(ctx context.Context, preMessage *models.MessagePayload) error {
	err := pgx.BeginFunc(ctx, m.pool, func(tx pgx.Tx) error {
		err := m.insertMessage(ctx, tx, preMessage)
		if err != nil {
			return err
		}

		err = m.addCounterMessage(ctx, tx, 1)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *MessagePg) GetMessageStatistic(ctx context.Context) (*models.MessageStatistic, error) {
	logger := m.logger.EnrichReqID(ctx)

	SQLGetMessageStatistic := `SELECT total, handled FROM public."counter_message"`
	messageStatistic := new(models.MessageStatistic)

	err := m.pool.QueryRow(ctx, SQLGetMessageStatistic).Scan(
		&messageStatistic.Total, &messageStatistic.Handled)
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	return messageStatistic, nil
}
