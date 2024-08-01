package middleware

import (
	"net/http"

	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/Vilinvil/task_messaggio/pkg/responses"
)

func Panic(next http.Handler, logger *mylogger.MyLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered: %+v\n", err)
				responses.SendResponse(w, logger, myerrors.ErrInternalServer)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
