package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Vilinvil/task_messaggio/internal/messageworker/messageworker/repository"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/Vilinvil/task_messaggio/pkg/utils"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
)

func TestSetStatusMessage(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	type TestCase struct {
		name               string
		inputUUID          *uuid.UUID
		inputStatusMessage repository.StatusMessage
		behaviourPool      func(mock pgxmock.PgxPoolIface)
		expectedErr        error
	}

	errInternalUpdateStatus := fmt.Errorf("errInternalUpdateStatus")   //nolint
	errInternalUpdateCounter := fmt.Errorf("errInternalUpdateCounter") //nolint

	testCases := []TestCase{
		{
			name:               "basic test",
			inputUUID:          &models.DummyUUID,
			inputStatusMessage: repository.StatusMessageDone,
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE public."message"`).
					WithArgs(repository.StatusMessageDone, &models.DummyUUID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(`UPDATE public."counter_message"`).WithArgs(1).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
				mock.ExpectRollback()
			},
			expectedErr: nil,
		},
		{
			name:               "test status pending",
			inputUUID:          &models.DummyUUID,
			inputStatusMessage: repository.StatusMessagePending,
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE public."message"`).
					WithArgs(repository.StatusMessagePending, &models.DummyUUID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(`UPDATE public."counter_message"`).WithArgs(1).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectCommit()
				mock.ExpectRollback()
			},
			expectedErr: nil,
		},
		{
			name:               "errInternalUpdateStatus",
			inputUUID:          &models.DummyUUID,
			inputStatusMessage: repository.StatusMessagePending,
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE public."message"`).
					WithArgs(repository.StatusMessagePending, &models.DummyUUID).
					WillReturnError(errInternalUpdateStatus)
				mock.ExpectRollback()
				mock.ExpectRollback()
			},
			expectedErr: errInternalUpdateStatus,
		},
		{
			name:               "errInternalUpdateCounter",
			inputUUID:          &models.DummyUUID,
			inputStatusMessage: repository.StatusMessagePending,
			behaviourPool: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE public."message"`).
					WithArgs(repository.StatusMessagePending, &models.DummyUUID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
				mock.ExpectExec(`UPDATE public."counter_message"`).WithArgs(1).
					WillReturnError(errInternalUpdateCounter)
				mock.ExpectRollback()
				mock.ExpectRollback()
			},
			expectedErr: errInternalUpdateCounter,
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

			receivedErr := messagePg.SetStatusMessage(baseCtx, testCase.inputUUID, testCase.inputStatusMessage)

			utils.EqualErrors(t, testCase.expectedErr, receivedErr)
			utils.PgxPoolExpectationWereMet(t, mock)
		})
	}
}
