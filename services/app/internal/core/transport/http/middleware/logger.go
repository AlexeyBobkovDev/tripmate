package core_middleware

import (
	"net/http"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
)

func LoggerMiddleware(logger *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			r = r.WithContext(core_logger.ToContext(r.Context(), logger))

			next.ServeHTTP(rw, r)
		})
	}
}
