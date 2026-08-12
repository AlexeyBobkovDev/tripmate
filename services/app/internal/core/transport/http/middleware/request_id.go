package core_middleware

import (
	"net/http"

	"github.com/google/uuid"
)

func RequestIDMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestID := uuid.NewString()
			r.Header.Add("X-Request-ID", requestID)

			next.ServeHTTP(rw, r)
		})
	}
}
