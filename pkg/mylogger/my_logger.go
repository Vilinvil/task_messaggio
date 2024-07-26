package mylogger

import (
	"context"
	"fmt"
	"sync"

	"github.com/Vilinvil/task_messaggio/pkg/myerrors"

	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger //nolint:gochecknoglobals
	once   sync.Once          //nolint:gochecknoglobals

	ErrNoLogger = myerrors.NewInternalServerError("my_logger.Get для отсутствующего логгера")
)

type MyLogger struct {
	*zap.SugaredLogger
}

func NewNop() *MyLogger {
	once.Do(func() {
		logger = zap.NewNop().Sugar()
	})

	return &MyLogger{logger}
}

func New(outputPaths []string, errorOutputPaths []string,
	productionMode bool, options ...zap.Option,
) (*MyLogger, error) {
	var err error

	once.Do(func() {
		var config zap.Config

		if productionMode {
			config = zap.NewProductionConfig()
		} else {
			config = zap.NewDevelopmentConfig()
		}

		config.OutputPaths = outputPaths
		config.ErrorOutputPaths = errorOutputPaths

		zapLogger, innerErr := config.Build(options...)
		if innerErr != nil {
			err = innerErr

			return
		}

		logger = zapLogger.Sugar()
	})

	if err != nil {
		return nil, err
	}

	return &MyLogger{logger}, nil
}

func Get() (*MyLogger, error) {
	if logger == nil {
		fmt.Println(ErrNoLogger)

		return nil, fmt.Errorf(myerrors.ErrTemplate, ErrNoLogger)
	}

	return &MyLogger{logger}, nil
}

func (m *MyLogger) EnrichReqID(ctx context.Context) *MyLogger {
	return &MyLogger{m.With(
		zap.String("requestID", GetRequestIDFromCtx(ctx)),
	)}
}
