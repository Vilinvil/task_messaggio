package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Vilinvil/task_messaggio/internal/message/message/mocks"
	"github.com/Vilinvil/task_messaggio/internal/message/message/usecases"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

var (
	testUUID          = uuid.MustParse("52fdfc07-2182-454f-963f-5f0f9a621d72") //nolint:gochecknoglobals
	testGeneratorUUID = func() uuid.UUID {                                     //nolint:gochecknoglobals
		return testUUID
	}
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

	return usecases.NewMessageService(mockMessageRepository, mockBrokerMessageRepository, testGeneratorUUID, logger)
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
					ID:    testUUID,
					Value: "basic message",
				})
			},
			behaviorBrokerMessageRepository: func(m *mocks.MockBrokerMessageRepository) {
				m.EXPECT().WriteMessage(gomock.Any(), &models.MessagePayload{
					ID:    testUUID,
					Value: "basic message",
				})
			},
			expectedResponse: testUUID,
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
					ID:    testUUID,
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
					ID:    testUUID,
					Value: "basic message",
				}).Return(nil)
			},
			behaviorBrokerMessageRepository: func(m *mocks.MockBrokerMessageRepository) {
				m.EXPECT().WriteMessage(gomock.Any(), &models.MessagePayload{
					ID:    testUUID,
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

			if receivedResponse != testCase.expectedResponse {
				t.Errorf("receivedResponse: %v NOT EQUAL expectedResponse: %v ",
					receivedResponse, testCase.expectedResponse)
			}

			if !errors.Is(err, testCase.expectedErr) {
				t.Errorf("receivedErr: %v NOT EQUAL expectedErr: %v", err, testCase.expectedErr)
			}
		})
	}
}
