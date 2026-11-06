package core_server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_middleware "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/middleware"
)

func TestNewAPIRouter(t *testing.T) {
	router := NewAPIRouter(APIVersionV1)

	require.NotNil(t, router)
	assert.Equal(t, APIVersionV1, router.APIVersion)
}

func TestAPIVersion_ToString(t *testing.T) {
	tests := []struct {
		name       string
		apiVersion APIVersion
		want       string
	}{
		{name: "version one", apiVersion: APIVersionV1, want: "v1"},
		{name: "version two", apiVersion: APIVersionV2, want: "v2"},
		{name: "version three", apiVersion: APIVersionV3, want: "v3"},
		{name: "custom version", apiVersion: APIVersion("v4"), want: "v4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.apiVersion.ToString())
		})
	}
}

func TestAPIRouter_RegisterRoutes(t *testing.T) {
	t.Run("registers single route and invokes handler", func(t *testing.T) {
		router := NewAPIRouter(APIVersionV1)

		handlerCalled := false
		router.RegisterRoutes(
			&Route{
				Method: http.MethodGet,
				Path:   "/health",
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handlerCalled = true
					w.WriteHeader(http.StatusOK)
				}),
			},
		)

		server := httptest.NewServer(router.mux)
		defer server.Close()

		response, err := http.Get(server.URL + "/health")
		require.NoError(t, err)
		defer func() {
			_ = response.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.True(t, handlerCalled)
	})

	t.Run("panics on nil route", func(t *testing.T) {
		router := NewAPIRouter(APIVersionV1)

		require.PanicsWithValue(t, "cannot register a new route or handler", func() {
			router.RegisterRoutes(nil)
		})
	})

	t.Run("panics on nil handler", func(t *testing.T) {
		router := NewAPIRouter(APIVersionV1)

		require.PanicsWithValue(t, "cannot register a new route or handler", func() {
			router.RegisterRoutes(&Route{
				Method:  http.MethodGet,
				Path:    "/health",
				Handler: nil,
			})
		})
	})

	t.Run("applies middlewares to registered handler", func(t *testing.T) {
		router := NewAPIRouter(APIVersionV1)

		var handlerCalled bool
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		router.Middlewares = []core_middleware.Middleware{
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-Test-Header", "middleware-value")
					next.ServeHTTP(w, r)
				})
			},
		}
		router.RegisterRoutes(&Route{
			Method:  http.MethodGet,
			Path:    "/middleware",
			Handler: handler,
		})

		server := httptest.NewServer(router.mux)
		defer server.Close()

		response, err := http.Get(server.URL + "/middleware")
		require.NoError(t, err)
		defer func() {
			_ = response.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.True(t, handlerCalled)
		assert.Equal(t, "middleware-value", response.Header.Get("X-Test-Header"))
	})
}
