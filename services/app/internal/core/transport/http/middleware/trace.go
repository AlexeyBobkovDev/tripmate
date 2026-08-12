package core_middleware

import (
	"net/http"

	"go.uber.org/zap"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
)

func TraceMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			logger := core_logger.FromContext(r.Context())
			requestID := r.Header.Get("X-Request-ID")

			logger.Debug("start processing request", zap.String("requestID", requestID))

			next.ServeHTTP(rw, r)

			logger.Debug("stop processing request", zap.String("requestID", requestID))
		})
	}
}
