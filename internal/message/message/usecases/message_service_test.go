package usecases_test

import (
	"context"
	"fmt"
	"github.com/Vilinvil/task_messaggio/pkg/utils"
	"testing"

	"github.com/Vilinvil/task_messaggio/internal/message/message/mocks"
	"github.com/Vilinvil/task_messaggio/internal/message/message/usecases"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func NewMessageService(ctrl *gomock.Controller,
	behaviorMessageRepository func(m *mocks.MockMessageRepository),
	behaviorBrokerMessageRepository func(m *mocks.MockBrokerMessageRepository),
	logger *mylogger.MyLogger,
) *usecases.MessageService {
	mockMessageRepository := mocks.NewMockMessageRepository(ctrl)
	mockBrokerMessageRepository := mocks.NewMockBrokerMessageRepository(ctrl)

	behaviorMessageRepository(mockMessageRepository)
	behaviorBrokerMessageRepository(mockBrokerMessageRepository)

	return usecases.NewMessageService(
		mockMessageRepository, mockBrokerMessageRepository, models.DummyGeneratorUUID, logger)
}

func TestAddMessage(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	type TestCase struct {
		name                            string
		inputValue                      string
		behaviorMessageRepository       func(m *mocks.MockMessageRepository)
		behaviorBrokerMessageRepository func(m *mocks.MockBrokerMessageRepository)
		expectedResponse                uuid.UUID
		expectedErr                     error
	}

	errTestInternal := fmt.Errorf("тестовая внутреняя ошибка") //nolint

	testCases := []TestCase{
		{
			name:       "basic test",
			inputValue: "basic message",
			behaviorMessageRepository: func(m *mocks.MockMessageRepository) {
				m.EXPECT().AddMessage(gomock.Any(), &models.MessagePayload{
					ID:    models.DummyUUID,
					Value: "basic message",
				})
			},
			behaviorBrokerMessageRepository: func(m *mocks.MockBrokerMessageRepository) {
				m.EXPECT().WriteMessage(gomock.Any(), &models.MessagePayload{
					ID:    models.DummyUUID,
					Value: "basic message",
				})
			},
			expectedResponse: models.DummyUUID,
			expectedErr:      nil,
		},
		{
			name:       "validation err zero len",
			inputValue: "",
			behaviorMessageRepository: func(m *mocks.MockMessageRepository) {
				m.EXPECT()
			},
			behaviorBrokerMessageRepository: func(m *mocks.MockBrokerMessageRepository) {
				m.EXPECT()
			},
			expectedResponse: uuid.UUID{},
			expectedErr:      models.ErrLenMessage,
		},
		{
			name:       "internal error in messageRepository.AddMessage(...)",
			inputValue: "basic message",
			behaviorMessageRepository: func(m *mocks.MockMessageRepository) {
				m.EXPECT().AddMessage(gomock.Any(), &models.MessagePayload{
					ID:    models.DummyUUID,
					Value: "basic message",
				}).Return(errTestInternal)
			},
			behaviorBrokerMessageRepository: func(m *mocks.MockBrokerMessageRepository) {
				m.EXPECT()
			},
			expectedResponse: uuid.UUID{},
			expectedErr:      errTestInternal,
		},

		{
			name:       "internal error in brokerMessageRepository.WriteMessage(...)",
			inputValue: "basic message",
			behaviorMessageRepository: func(m *mocks.MockMessageRepository) {
				m.EXPECT().AddMessage(gomock.Any(), &models.MessagePayload{
					ID:    models.DummyUUID,
					Value: "basic message",
				}).Return(nil)
			},
			behaviorBrokerMessageRepository: func(m *mocks.MockBrokerMessageRepository) {
				m.EXPECT().WriteMessage(gomock.Any(), &models.MessagePayload{
					ID:    models.DummyUUID,
					Value: "basic message",
				}).Return(errTestInternal)
			},
			expectedResponse: uuid.UUID{},
			expectedErr:      errTestInternal,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			messageService := NewMessageService(
				ctrl, testCase.behaviorMessageRepository, testCase.behaviorBrokerMessageRepository, nopLogger)

			ctx := context.Background()

			receivedResponse, err := messageService.AddMessage(ctx, testCase.inputValue)

			utils.PlainEqual(t, receivedResponse, testCase.expectedResponse)

			utils.EqualErrors(t, err, testCase.expectedErr)
		})
	}
}
