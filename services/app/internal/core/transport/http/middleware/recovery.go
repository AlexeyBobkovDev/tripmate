package core_middleware

import (
	"net/http"

	"go.uber.org/zap"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
)

func RecoveryMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				logger := core_logger.FromContext(r.Context())

				if p := recover(); p != nil {
					logger.Error(
						"panic occurred",
						zap.Any("panic", p),
						zap.Stack("stack"),
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
