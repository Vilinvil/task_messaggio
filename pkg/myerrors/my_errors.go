package myerrors

import "net/http"

var ErrInternalServer = NewInternalServerError("Внутренняя ошибка на сервере")

//easyjson:json
type Error struct {
	Err    string `json:"err"`
	status int
}

func New(err string, status int) *Error {
	return &Error{Err: err, status: status}
}

func NewBadRequestError(err string) *Error {
	return New(err, http.StatusBadRequest)
}

func NewInternalServerError(err string) *Error {
	return New(err, http.StatusInternalServerError)
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
