package myerrors

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInternalServer = NewInternalServerError("Внутренняя ошибка на сервере")

//easyjson:json
type Error struct {
	Err    string `json:"reason"` //nolint:tagliatelle
	status int
}

func New(status int, err string) *Error {
	return &Error{Err: err, status: status}
}

func NewBadRequestError(err string) *Error {
	return New(http.StatusBadRequest, err)
}

func NewInternalServerError(err string) *Error {
	return New(http.StatusInternalServerError, err)
}

func (e *Error) Error() string {
	return e.Err
}

func (e *Error) Status() int {
	return e.status
}

func (e *Error) IsClientError() bool {
	return e.status >= http.StatusBadRequest && e.status < http.StatusInternalServerError
}

func (e *Error) ConvertToGRPC() error {
	switch e.status {
	case http.StatusBadRequest:
		return status.Error(codes.InvalidArgument, e.Err)
	case http.StatusNotFound:
		return status.Error(codes.NotFound, e.Err)
	default:
		return status.Error(codes.Internal, ErrInternalServer.Err)
	}
}
