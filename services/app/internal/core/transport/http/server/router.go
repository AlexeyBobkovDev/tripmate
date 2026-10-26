package core_server

import (
	"fmt"
	"net/http"

	core_middleware "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/middleware"
)

type APIVersion string

func (v APIVersion) ToString() string {
	return string(v)
}

const (
	APIVersionV1 APIVersion = "v1"
	APIVersionV2 APIVersion = "v2"
	APIVersionV3 APIVersion = "v3"
)

type APIRouter struct {
	mux         *http.ServeMux
	APIVersion  APIVersion
	Middlewares []core_middleware.Middleware
}

func NewAPIRouter(
	APIVersion APIVersion,
) *APIRouter {
	return &APIRouter{
		mux:        http.NewServeMux(),
		APIVersion: APIVersion,
	}
}

func (r *APIRouter) RegisterRoutes(routes ...*Route) {
	for _, route := range routes {
		if route == nil || route.Handler == nil {
			panic("cannot register a new route or handler")
		}

		path := fmt.Sprintf("%s %s", route.Method, route.Path)

		handler := core_middleware.ChainMiddleware(
			http.HandlerFunc(route.Handler),
			r.Middlewares...,
		)
		r.mux.Handle(
			path,
			handler,
		)
	}
}
