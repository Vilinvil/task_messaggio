package mylogger

import (
	"context"
	"math/rand"
	"strconv"
)

type keyCtx string

const (
	requestIDKey keyCtx = "requestIDKey"
)

func SetRequestIDToCtx(ctx context.Context, requestID string) context.Context {
	ctx = context.WithValue(ctx, requestIDKey, requestID)

	return ctx
}

func AddRequestIDToCtx(ctx context.Context) context.Context {
	requestID := strconv.Itoa(rand.Int()) //nolint:gosec

	return SetRequestIDToCtx(ctx, requestID)
}

func GetRequestIDFromCtx(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}
