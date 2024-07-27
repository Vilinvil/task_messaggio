package delivery

import (
	"context"
	"fmt"
	"net/http"

	messagedelivery "github.com/Vilinvil/task_messaggio/internal/message/message/delivery"
	repository2 "github.com/Vilinvil/task_messaggio/internal/message/message/repository"
	"github.com/Vilinvil/task_messaggio/internal/message/message/usecases"
	"github.com/Vilinvil/task_messaggio/pkg/delivery"
	mymiddleware "github.com/Vilinvil/task_messaggio/pkg/middleware"
	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"

	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/repository"
)

func NewMux(ctx context.Context, urlDataBase string, logger *mylogger.MyLogger) (http.Handler, error) {
	router := http.NewServeMux()

	pool, err := repository.NewPgxPool(ctx, urlDataBase)
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	messagePg := repository2.NewMessagePg(pool, logger)

	messageService := usecases.NewMessageService(messagePg, logger)

	messageHandler := messagedelivery.NewMessageHandler(messageService, logger)

	router.HandleFunc("GET /healthcheck", delivery.HealthCheckHandler)
	router.HandleFunc("POST /message", messageHandler.AddMessage)
	router.HandleFunc("GET /message/statistic", messageHandler.GetMessageStatistic)

	mux := http.NewServeMux()
	mux.Handle("/", mymiddleware.Panic(middleware.Context(ctx,
		mymiddleware.AddReqID(
			mymiddleware.AccessLogMiddleware(
				middleware.AddAPIName(router, middleware.APINameV1),
				logger))),
		logger))

	return mux, nil
}
