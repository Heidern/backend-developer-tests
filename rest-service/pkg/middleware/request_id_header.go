package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

const requestIDHeader = "X-Request-Id"

func RequestIDHeader(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if w.Header().Get(requestIDHeader) == "" {
			requestID := middleware.GetReqID(r.Context())
			w.Header().Add(requestIDHeader, requestID)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
