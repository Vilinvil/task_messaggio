package mylogger

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/mylogger"
	"math"
	"math/big"
)

type keyCtx string

const (
	requestIDKey keyCtx = "requestIDKey"
)

func SetRequestIDToCtx(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func AddRequestIDToCtx(ctx context.Context) (context.Context, error) {
	logger, err := mylogger.Get()
	if err != nil {
		return nil, err
	}

	bigInt, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	requestID := fmt.Sprintf("%d", bigInt.Int64())

	return SetRequestIDToCtx(ctx, requestID), nil
}

func GetRequestIDFromCtx(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}
