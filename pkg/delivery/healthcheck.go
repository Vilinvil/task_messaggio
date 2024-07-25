package delivery

import (
	"net/http"

	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
	"github.com/Vilinvil/task_messaggio/pkg/responses"
)

func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	logger, err := mylogger.Get()
	if err != nil {
		http.Error(w, myerrors.ErrInternalServer.Error(), myerrors.ErrInternalServer.Status())

		return
	}

	responses.SendResponse(w, logger, responses.NewResponseSuccessful("OK"))
}
