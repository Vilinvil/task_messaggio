package usecases_test

import (
	"context"
	"errors"
	"math/rand"
	"sync"
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

	return usecases.NewMessageService(mockMessageRepository, mockBrokerMessageRepository, logger)
}

func TestMessageService_AddMessage(t *testing.T) {
	t.Parallel()

	uuid.SetRand(rand.New(rand.NewSource(1))) //nolint:gosec

	muRand := sync.Mutex{}

	nopLogger := mylogger.NewNop()

	type TestCase struct {
		name                            string
		inputValue                      string
		behaviorMessageRepository       func(m *mocks.MockMessageRepository)
		behaviorBrokerMessageRepository func(m *mocks.MockBrokerMessageRepository)
		expectedResponse                uuid.UUID
		expectedErr                     error
	}

	testCases := []TestCase{
		{
			name:       "basic test",
			inputValue: "basic message",
			behaviorMessageRepository: func(m *mocks.MockMessageRepository) {
				m.EXPECT().AddMessage(gomock.Any(), &models.MessagePayload{
					ID:    uuid.MustParse("52fdfc07-2182-454f-963f-5f0f9a621d72"),
					Value: "basic message",
				})
			},
			behaviorBrokerMessageRepository: func(m *mocks.MockBrokerMessageRepository) {
				m.EXPECT().WriteMessage(gomock.Any(), &models.MessagePayload{
					ID:    uuid.MustParse("52fdfc07-2182-454f-963f-5f0f9a621d72"),
					Value: "basic message",
				})
			},
			expectedResponse: uuid.MustParse("52fdfc07-2182-454f-963f-5f0f9a621d72"),
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

			muRand.Lock()

			receivedResponse, err := messageService.AddMessage(ctx, testCase.inputValue)

			muRand.Unlock()

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
