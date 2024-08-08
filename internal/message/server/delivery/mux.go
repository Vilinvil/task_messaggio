package delivery

import (
	"context"
	"net/http"

	// /docs - need for handler with swagger documentations
	_ "github.com/Vilinvil/task_messaggio/docs"
	messagedelivery "github.com/Vilinvil/task_messaggio/internal/message/message/delivery"
	"github.com/Vilinvil/task_messaggio/internal/message/message/repository"
	"github.com/Vilinvil/task_messaggio/internal/message/message/usecases"
	"github.com/Vilinvil/task_messaggio/pkg/delivery"
	mymiddleware "github.com/Vilinvil/task_messaggio/pkg/middleware"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewMux(ctx context.Context, urlDataBase, brokerAddr string, logger *mylogger.MyLogger) (http.Handler, error) {
	router := http.NewServeMux()

	messagePg, err := repository.NewMessagePg(ctx, urlDataBase, logger)
	if err != nil {
		return nil, err
	}

	brokerMessage, err := repository.NewBrokerMessageKafka(brokerAddr, logger)
	if err != nil {
		return nil, err
	}

	messageService := usecases.NewMessageService(messagePg, brokerMessage, uuid.New, logger)

	messageHandler := messagedelivery.NewMessageHandler(messageService, logger)

	router.HandleFunc("GET /healthcheck", delivery.HealthCheckHandler)
	router.HandleFunc("POST /message", messageHandler.AddMessage)
	router.HandleFunc("GET /message/statistic", messageHandler.GetMessageStatistic)

	router.HandleFunc("GET /swagger/", httpSwagger.Handler())

	mux := http.NewServeMux()
	mux.Handle("/", mymiddleware.Panic(middleware.Context(ctx,
		mymiddleware.AddReqID(
			mymiddleware.AccessLogMiddleware(
				middleware.AddAPIName(router, middleware.APINameV1),
				logger))),
		logger))

	return mux, nil
}
