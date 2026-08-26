package core_middleware

import (
	"net/http"
	"slices"
	"strings"
)

func CORSMiddleware(config CORSConfig) Middleware {
	allowedOrigins := make(map[string]struct{}, len(config.Origins))
	for _, allowedOrigin := range config.Origins {
		allowedOrigins[allowedOrigin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Vary", "Origin")
			requestMethod := r.Header.Get("Access-Control-Request-Method")
			if requestMethod != "" && !slices.Contains(config.Methods, requestMethod) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			origin := r.Header.Get("Origin")
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			if _, ok := allowedOrigins[origin]; !ok {
				http.Error(
					w,
					"forbidden",
					http.StatusForbidden,
				)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.Methods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.Headers, ", "))

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
