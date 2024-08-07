package responses_test

import (
	"fmt"
	"github.com/Vilinvil/task_messaggio/pkg/models"
	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/Vilinvil/task_messaggio/pkg/responses"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/utils/test"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ErrMarshallerResponse struct{}

func (e ErrMarshallerResponse) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("Internal err MarshallJSON") //nolint
}

func (e ErrMarshallerResponse) Status() int {
	return http.StatusInternalServerError
}

const contentTypeTextPlain = "text/plain; charset=utf-8"

func TestSendResponse(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	type TestCase struct {
		name                string
		inputResponse       responses.Response
		expectedContentType string
		expectedResponse    any
	}

	testCases := []TestCase{
		{
			name:                "basic test",
			inputResponse:       responses.NewResponseSuccessful("basic test"),
			expectedContentType: responses.BaseContentType,
			expectedResponse:    responses.NewResponseSuccessful("basic test"),
		},
		{
			name: "response message statistic",
			inputResponse: models.NewResponseMessageStatistic(
				http.StatusOK,
				&models.MessageStatistic{Total: 20, Handled: 10},
			),
			expectedContentType: responses.BaseContentType,
			expectedResponse: models.NewResponseMessageStatistic(
				http.StatusOK,
				&models.MessageStatistic{Total: 20, Handled: 10},
			),
		},
		{
			name:                "error in MarshallJSON",
			inputResponse:       ErrMarshallerResponse{},
			expectedContentType: contentTypeTextPlain,
			expectedResponse:    myerrors.ErrInternalServer.Error() + "\n",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			w := httptest.NewRecorder()

			responses.SendResponse(w, nopLogger, testCase.inputResponse)

			if receivedContentType := w.Header().Get("Content-Type"); receivedContentType != testCase.expectedContentType {
				t.Errorf("Expected content-type: %s received: %s", testCase.expectedContentType, receivedContentType)
			}

			if expectedCode := testCase.inputResponse.Status(); w.Code != expectedCode {
				t.Errorf("Expected http code: %d received: %d", expectedCode, w.Code)
			}

			err := test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Errorf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}

func TestSendErrResponse(t *testing.T) {
	t.Parallel()

	nopLogger := mylogger.NewNop()

	type TestCase struct {
		name             string
		err              error
		expectedResponse responses.Response
	}

	testCases := []TestCase{
		{
			name:             "not myerror type of error",
			err:              fmt.Errorf("not myerror type of error"),
			expectedResponse: myerrors.ErrInternalServer,
		},
		{
			name:             "internal server error",
			err:              myerrors.NewInternalServerError("test internal server err"),
			expectedResponse: myerrors.ErrInternalServer,
		},
		{
			name:             "client bad request error",
			err:              myerrors.NewBadRequestError("not enough arguments"),
			expectedResponse: myerrors.NewBadRequestError("not enough arguments"),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			w := httptest.NewRecorder()

			responses.SendErrResponse(w, nopLogger, testCase.err)

			if receivedContentType := w.Header().Get("Content-Type"); receivedContentType != responses.BaseContentType {
				t.Errorf("Expected content-type: %s received: %s", responses.BaseContentType, receivedContentType)
			}

			if expectedCode := testCase.expectedResponse.Status(); w.Code != expectedCode {
				t.Errorf("Expected http code: %d received: %d", expectedCode, w.Code)
			}

			err := test.CompareHTTPTestResult(w, testCase.expectedResponse)
			if err != nil {
				t.Errorf("Failed CompareHTTPTestResult %+v", err)
			}
		})
	}
}
