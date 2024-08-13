package mylogger_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"testing"
	"testing/iotest"

	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/Vilinvil/task_messaggio/pkg/utils"
	"google.golang.org/grpc/metadata"
)

func TestAddRequestIDToCtx(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	amountTimes := 100

	for i := range amountTimes {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			baseCtx := context.Background()

			if reqID := mylogger.GetRequestIDFromCtx(baseCtx); reqID != "" {
				t.Fatalf("not expected reqID: %s", reqID)
			}

			ctxWithReqestID, err := mylogger.AddRequestIDToCtx(baseCtx, rand.Reader)
			utils.EqualErrors(t, nil, err)

			if reqID := mylogger.GetRequestIDFromCtx(ctxWithReqestID); reqID == "" {
				t.Fatalf("expected reqID but get empty string")
			}
		})
	}
}

func TestAddRequestErrNoLogger(t *testing.T) { //nolint:paralleltest
	baseCtx := context.Background()

	ctxWithRequestID, err := mylogger.AddRequestIDToCtx(baseCtx, rand.Reader)
	utils.EqualErrors(t, mylogger.ErrNoLogger, err)
	utils.PlainEqual(t, nil, ctxWithRequestID)
}

func TestAddRequestErrRandReader(t *testing.T) {
	t.Parallel()

	_ = mylogger.NewNop()

	errInternalReader := fmt.Errorf("errInternalReader") //nolint

	errReader := iotest.ErrReader(errInternalReader)

	baseCtx := context.Background()

	ctxWithReqestID, err := mylogger.AddRequestIDToCtx(baseCtx, errReader)
	utils.EqualErrors(t, errInternalReader, err)
	utils.PlainEqual(t, nil, ctxWithReqestID)
}

func TestGetRequestIDFromMDCtx(t *testing.T) {
	t.Parallel()

	baseCtx := context.Background()
	ctxWithoutReqID := baseCtx

	expectedReqID := "1234_test_req_id"
	ctxWithReqID := metadata.NewIncomingContext(baseCtx,
		mylogger.NewMDFromRequestIDCtx(mylogger.SetRequestIDToCtx(baseCtx, expectedReqID)))

	utils.PlainEqual(t, "", mylogger.GetRequestIDFromMDCtx(ctxWithoutReqID))
	utils.PlainEqual(t, expectedReqID, mylogger.GetRequestIDFromMDCtx(ctxWithReqID))
}
