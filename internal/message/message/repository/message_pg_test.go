package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Vilinvil/task_messaggio/internal/message/message/repository"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/Vilinvil/task_messaggio/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
)

func TestGetMessageStatistic(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	type TestCase struct {
		name             string
		behaviourPool    func(mock pgxmock.PgxPoolIface)
		expectedResponse *models.MessageStatistic
		expectedErr      error
	}

	errInternalPgx := fmt.Errorf("internal pgx error") //nolint
	testBigNum := 1 << 60

	testCases := []TestCase{
		{
			name: "basic test",
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT total, handled`).WillReturnRows(pgxmock.NewRows([]string{
					"total", "handled",
				}).AddRow(0, 0))
			},
			expectedResponse: &models.MessageStatistic{Total: 0, Handled: 0},
			expectedErr:      nil,
		},
		{
			name: "big numbers",
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT total, handled`).WillReturnRows(pgxmock.NewRows([]string{
					"total", "handled",
				}).AddRow(testBigNum, testBigNum))
			},
			expectedResponse: &models.MessageStatistic{Total: testBigNum, Handled: testBigNum},
			expectedErr:      nil,
		},
		{
			name: "err no rows",
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT total, handled`).WillReturnError(pgx.ErrNoRows)
			},
			expectedResponse: nil,
			expectedErr:      pgx.ErrNoRows,
		},
		{
			name: "internal pgx err",
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT total, handled`).WillReturnError(errInternalPgx)
			},
			expectedResponse: nil,
			expectedErr:      errInternalPgx,
		},
	}

	baseCtx := context.Background()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			utils.EqualErrors(t, nil, err)

			defer mock.Close()

			testCase.behaviourPool(mock)

			messagePg := repository.NewMessagePg(mock, nopLogger)

			receivedMsgStatistic, receivedErr := messagePg.GetMessageStatistic(baseCtx)

			utils.EqualErrors(t, testCase.expectedErr, receivedErr)
			utils.DeepEqual(t, testCase.expectedResponse, receivedMsgStatistic)
			utils.PgxPoolExpectationWereMet(t, mock)
		})
	}
}

func TestAddMessage(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	type TestCase struct {
		name            string
		behaviourPool   func(mock pgxmock.PgxPoolIface)
		inputPreMessage *models.MessagePayload
		expectedErr     error
	}

	errInsertionMessage := fmt.Errorf("internal insertion err")   //nolint
	errUpdateCounter := fmt.Errorf("internal update counter err") //nolint

	testCases := []TestCase{
		{
			name: "basic test",
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO public."message"`).WithArgs(models.DummyUUID, "basic value").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(`UPDATE public."counter_message"`).WithArgs(1).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
				mock.ExpectRollback()
			},
			inputPreMessage: &models.MessagePayload{ID: models.DummyUUID, Value: "basic value"},
			expectedErr:     nil,
		},
		{
			name: "err insertion message",
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO public."message"`).WithArgs(models.DummyUUID, "basic value").
					WillReturnError(errInsertionMessage)
				mock.ExpectRollback()
				mock.ExpectRollback()
			},
			inputPreMessage: &models.MessagePayload{ID: models.DummyUUID, Value: "basic value"},
			expectedErr:     errInsertionMessage,
		},
		{
			name: "err update counter",
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO public."message"`).WithArgs(models.DummyUUID, "basic value").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectExec(`UPDATE public."counter_message"`).WithArgs(1).
					WillReturnError(errUpdateCounter)
				mock.ExpectRollback()
				mock.ExpectRollback()
			},
			inputPreMessage: &models.MessagePayload{ID: models.DummyUUID, Value: "basic value"},
			expectedErr:     errUpdateCounter,
		},
	}

	baseCtx := context.Background()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			utils.EqualErrors(t, nil, err)

			defer mock.Close()

			testCase.behaviourPool(mock)

			messagePg := repository.NewMessagePg(mock, nopLogger)

			receivedErr := messagePg.AddMessage(baseCtx, testCase.inputPreMessage)

			utils.EqualErrors(t, testCase.expectedErr, receivedErr)
			utils.PgxPoolExpectationWereMet(t, mock)
		})
	}
}
