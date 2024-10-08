package myerrors_test

import (
	"net/http"
	"testing"

	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestError_ConvertToGRPC(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		inputErr    *myerrors.Error
		expectedErr error
	}

	testCases := []TestCase{
		{
			name:        "bad request",
			inputErr:    myerrors.NewBadRequestError("err bad request"),
			expectedErr: status.Error(codes.InvalidArgument, "err bad request"),
		},
		{
			name:        "err not found",
			inputErr:    myerrors.New(http.StatusNotFound, "err not found"),
			expectedErr: status.Error(codes.NotFound, "err not found"),
		},
		{
			name:        "err internal",
			inputErr:    myerrors.NewInternalServerError("err internal"),
			expectedErr: status.Error(codes.Internal, myerrors.ErrInternalServer.Err),
		},
		{
			name:        "unknown status err is internal err",
			inputErr:    myerrors.New(999, "999 status err"),
			expectedErr: status.Error(codes.Internal, myerrors.ErrInternalServer.Err),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			receivedErr := testCase.inputErr.ConvertToGRPC()
			utils.DeepEqual(t, receivedErr, testCase.expectedErr)
		})
	}
}
