package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
)

var ErrNonExistResponseStatus = myerrors.NewInternalServerError("отсутствует статус http ответа")

const NonExistResponseStatus = -1

type WriterWithStatus struct {
	http.ResponseWriter
	Status int
}

func (w *WriterWithStatus) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func AccessLogMiddleware(next http.Handler, logger *mylogger.MyLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writerWithStatus := &WriterWithStatus{ResponseWriter: w, Status: NonExistResponseStatus}

		start := time.Now()

		next.ServeHTTP(writerWithStatus, r)

		duration := time.Since(start)

		logger := logger.EnrichReqID(r.Context())

		path := r.URL.Path
		method := r.Method
		statusStr := strconv.Itoa(writerWithStatus.Status)

		if writerWithStatus.Status == NonExistResponseStatus {
			logger.Warn(ErrNonExistResponseStatus)

			statusStr = ""
		}

		logger.Infof(
			"path: %s method: %s status: %s duration: %v remoreAddr: %s",
			path, method, statusStr, duration, r.RemoteAddr)
	})
}
