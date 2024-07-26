package repository

import (
	"context"
	"fmt"

	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessagePg struct {
	pool   *pgxpool.Pool
	logger *mylogger.MyLogger
}

func NewMessagePg(pool *pgxpool.Pool, logger *mylogger.MyLogger) *MessagePg {
	return &MessagePg{
		pool:   pool,
		logger: logger,
	}
}

func (m *MessagePg) insertMessage(ctx context.Context, tx pgx.Tx, preMessage *models.PreMessage) error {
	logger := m.logger.EnrichReqID(ctx)

	SQLInsertMessage := `INSERT INTO public."message" (id, value) VALUES ($1, $2)`

	_, err := tx.Exec(ctx, SQLInsertMessage, preMessage.ID, preMessage.Value)
	if err != nil {
		logger.Error(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (m *MessagePg) addCounterMessage(ctx context.Context, tx pgx.Tx, delta int) error {
	logger := m.logger.EnrichReqID(ctx)

	SQLAddCounterMessage := `UPDATE public."counter_message" SET total = total + $1`

	_, err := tx.Exec(ctx, SQLAddCounterMessage, delta)
	if err != nil {
		logger.Error(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}

func (m *MessagePg) AddMessage(ctx context.Context, preMessage *models.PreMessage) error {
	logger := m.logger.EnrichReqID(ctx)

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
		logger.Error(err)

		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	return nil
}
