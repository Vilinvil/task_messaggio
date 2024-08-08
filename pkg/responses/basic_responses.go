package responses

import (
	"errors"
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

const BaseContentType = "application/json"

func SendResponse(w http.ResponseWriter, logger *mylogger.MyLogger, response Response) {
	w.Header().Set("Content-Type", BaseContentType)

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

func SendErrResponse(w http.ResponseWriter, logger *mylogger.MyLogger, err error) {
	w.Header().Set("Content-Type", "application/json")

	myErr := &myerrors.Error{} //nolint:exhaustruct
	if errors.As(err, &myErr) && myErr.IsClientError() {
		SendResponse(w, logger, myErr)

		return
	}

	SendResponse(w, logger, myerrors.ErrInternalServer)
}
