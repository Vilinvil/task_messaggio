package mylogger_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/Vilinvil/task_messaggio/pkg/utils"
	"strconv"
	"testing"
	"testing/iotest"
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

func TestAddRequestErrNoLogger(t *testing.T) {
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
