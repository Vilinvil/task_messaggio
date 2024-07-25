package responses

import (
	"net/http"

	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
)

type Marshaller interface {
	MarshalJSON() ([]byte, error)
}

type Response interface {
	Marshaller
	Status() int
}

func SendResponse(w http.ResponseWriter, logger *mylogger.MyLogger, response Response) {
	w.Header().Set("Content-Type", "application/json")

	responseSend, err := response.MarshalJSON()
	if err != nil {
		logger.Error(err)
		http.Error(w, myerrors.ErrInternalServer.Error(), myerrors.ErrInternalServer.Status())

		return
	}

	w.WriteHeader(response.Status())

	_, err = w.Write(responseSend)
	if err != nil {
		logger.Error(err)
		http.Error(w, myerrors.ErrInternalServer.Error(), myerrors.ErrInternalServer.Status())
	}
}

//easyjson:json
type ResponseSuccessful struct {
	status int
	Body   string `json:"body"`
}

func NewResponseSuccessful(body string) *ResponseSuccessful {
	return &ResponseSuccessful{
		status: http.StatusOK,
		Body:   body,
	}
}

func (r *ResponseSuccessful) Status() int {
	return r.status
}
