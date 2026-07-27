package core_middleware

import (
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
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

func LoggerMiddleware(logger *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			r = r.WithContext(core_logger.ToContext(r.Context(), logger))

			next.ServeHTTP(rw, r)
		})
	}
}

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

func RecoveryMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				logger := core_logger.FromContext(r.Context())

				if p := recover(); r != nil {
					logger.Error(
						"panic occurred",
						zap.Any("panic", p),
					)
					http.Error(
						rw,
						"internal server error",
						http.StatusInternalServerError,
					)
				}
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
