package mylogger

import (
	"context"
	"crypto/rand"
	"io"
	"math"
	"math/big"
	"strconv"

	"google.golang.org/grpc/metadata"
)

type keyCtx string

const (
	requestIDKey keyCtx = "request_id_key"
)

func SetRequestIDToCtx(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func AddRequestIDToCtx(ctx context.Context, readerRandom io.Reader) (context.Context, error) {
	logger, err := Get()
	if err != nil {
		return nil, err
	}

	bigInt, err := rand.Int(readerRandom, big.NewInt(math.MaxInt))
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	requestID := strconv.FormatInt(bigInt.Int64(), 10)

	return SetRequestIDToCtx(ctx, requestID), nil
}

func GetRequestIDFromCtx(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}

func NewMDFromRequestIDCtx(ctx context.Context) metadata.MD {
	return metadata.Pairs(string(requestIDKey), GetRequestIDFromCtx(ctx))
}

func GetRequestIDFromMDCtx(ctx context.Context) string {
	slStr := metadata.ValueFromIncomingContext(ctx, string(requestIDKey))
	if len(slStr) < 1 {
		return ""
	}

	return slStr[0]
}
