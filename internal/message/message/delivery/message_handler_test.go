package delivery_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/Vilinvil/task_messaggio/internal/message/message/delivery"
	"github.com/Vilinvil/task_messaggio/internal/message/message/mocks"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func NewMessageHandler(ctrl *gomock.Controller,
	behaviorMessageService func(m *mocks.MockMessageService), logger *mylogger.MyLogger,
) *delivery.MessageHandler {
	mockMessageService := mocks.NewMockMessageService(ctrl)

	behaviorMessageService(mockMessageService)

	return delivery.NewMessageHandler(mockMessageService, logger)
}

func TestAddMessage(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductService func(m *mocks.MockMessageService)
		body                   io.Reader
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name: "basic work",
			body: strings.NewReader("basic message"),
			behaviorProductService: func(m *mocks.MockMessageService) {
				m.EXPECT().AddMessage(gomock.Any(), "basic message").Return(uuid.UUID{1}, nil)
			},
			expectedResponse: delivery.ResponseMessageAddSuccessful,
		},
		{
			name: "русское сообщение",
			body: strings.NewReader("%D1%80%D1%83%D1%81%D1%81%D0%BA%D0%BE%D0%B5" +
				"%20%D1%81%D0%BE%D0%BE%D0%B1%D1%89%D0%B5%D0%BD%D0%B8%D0%B5"),
			behaviorProductService: func(m *mocks.MockMessageService) {
				m.EXPECT().AddMessage(gomock.Any(), "русское сообщение").Return(uuid.UUID{1}, nil)
			},
			expectedResponse: delivery.ResponseMessageAddSuccessful,
		},
		{
			name: "ошибка декодирования с",
			body: strings.NewReader("%%%"),
			behaviorProductService: func(m *mocks.MockMessageService) {
				m.EXPECT()
			},
			expectedResponse: delivery.ErrBadFormatMessage,
		},
		{
			name: "внутренняя ошибка сервиса",
			body: strings.NewReader("basic message"),
			behaviorProductService: func(m *mocks.MockMessageService) {
				m.EXPECT().AddMessage(gomock.Any(), "basic message").Return(
					uuid.UUID{}, fmt.Errorf("тестовая внутреняя ошибка"), //nolint
				)
			},
			expectedResponse: myerrors.ErrInternalServer,
		},
		{
			name: "ошибка считывания тела запроса",
			body: iotest.ErrReader(fmt.Errorf("тестовая ошибка")), //nolint
			behaviorProductService: func(m *mocks.MockMessageService) {
				m.EXPECT()
			},
			expectedResponse: myerrors.ErrInternalServer,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			messageHandler := NewMessageHandler(ctrl, testCase.behaviorProductService, nopLogger)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/message", testCase.body)

			messageHandler.AddMessage(w, req)

			err := test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestGetMessageStatistic(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	type TestCase struct {
		name                   string
		behaviorProductService func(m *mocks.MockMessageService)
		body                   io.Reader
		expectedResponse       any
	}

	testCases := [...]TestCase{
		{
			name: "basic work",
			body: strings.NewReader("basic message"),
			behaviorProductService: func(m *mocks.MockMessageService) {
				m.EXPECT().GetMessageStatistic(gomock.Any()).Return(
					&models.MessageStatistic{Total: 0, Handled: 0}, nil)
			},
			expectedResponse: models.NewResponseMessageStatistic(
				http.StatusOK, &models.MessageStatistic{Total: 0, Handled: 0}),
		},
		{
			name: "внутренняя ошибка сервиса",
			body: strings.NewReader("basic message"),
			behaviorProductService: func(m *mocks.MockMessageService) {
				m.EXPECT().GetMessageStatistic(gomock.Any()).Return(
					nil, fmt.Errorf("тестовая внутреняя ошибка"), //nolint
				)
			},
			expectedResponse: myerrors.ErrInternalServer,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			messageHandler := NewMessageHandler(ctrl, testCase.behaviorProductService, nopLogger)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/message/statistic", testCase.body)

			messageHandler.GetMessageStatistic(w, req)

			err := test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Fatalf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}
