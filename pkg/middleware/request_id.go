package middleware

import (
	"net/http"

	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
)

func AddReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := mylogger.AddRequestIDToCtx(r.Context())

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
