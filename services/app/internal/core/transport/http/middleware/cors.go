package core_middleware

import (
	"net/http"
)

func CORSMiddleware(allowedOriginsList map[string]struct{}) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if r.Method == http.MethodOptions {
				if _, ok := allowedOriginsList[origin]; ok {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Method", "GET, POST, PUT, PATCH, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
