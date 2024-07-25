package delivery

import (
	"context"
	"net/http"

	"github.com/Vilinvil/task_messaggio/pkg/delivery"
	mymiddleware "github.com/Vilinvil/task_messaggio/pkg/middleware"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
)

func NewMux(ctx context.Context, logger *mylogger.MyLogger) (http.Handler, error) {
	router := http.NewServeMux()

	router.HandleFunc("GET /healthcheck", delivery.HealthCheckHandler)

	mux := http.NewServeMux()
	mux.Handle("/", mymiddleware.Panic(middleware.Context(ctx,
		mymiddleware.AddReqID(
			mymiddleware.AccessLogMiddleware(
				middleware.AddAPIName(router, middleware.APINameV1),
				logger))),
		logger))

	return mux, nil
}
