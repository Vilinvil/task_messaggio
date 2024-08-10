package repository

import (
	"context"

	"github.com/Vilinvil/task_messaggio/pkg/dbpool"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	StatusMessagePending StatusMessage = "pending"
	StatusMessageDone    StatusMessage = "done"
)

type StatusMessage string

type MessagePg struct {
	pool   dbpool.PgxPool
	logger *mylogger.MyLogger
}

func NewMessagePg(pool dbpool.PgxPool, logger *mylogger.MyLogger) *MessagePg {
	return &MessagePg{
		pool:   pool,
		logger: logger,
	}
}

func (m *MessagePg) updateStatus(ctx context.Context, tx pgx.Tx, messageID *uuid.UUID, status StatusMessage) error {
	logger := m.logger.EnrichReqID(ctx)

	SQLUpdateStatus := `UPDATE public."message" SET status = $1 WHERE id = $2`

	_, err := tx.Exec(ctx, SQLUpdateStatus, status, messageID)
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (m *MessagePg) updateCounterMessageHandled(ctx context.Context, tx pgx.Tx, delta int) error {
	logger := m.logger.EnrichReqID(ctx)

	SQLUpdateCounterMessageHandled := `UPDATE public."counter_message" SET handled = handled + $1`

	_, err := tx.Exec(ctx, SQLUpdateCounterMessageHandled, delta)
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (m *MessagePg) SetStatusMessage(ctx context.Context, messageID *uuid.UUID, status StatusMessage) error {
	err := pgx.BeginFunc(ctx, m.pool, func(tx pgx.Tx) error {
		err := m.updateStatus(ctx, tx, messageID, status)
		if err != nil {
			return err
		}

		err = m.updateCounterMessageHandled(ctx, tx, 1)
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
