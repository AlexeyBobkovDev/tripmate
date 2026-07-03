package core_server

import (
	"fmt"
	"net/http"
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
	mux        *http.ServeMux
	Path       string
	APIVersion APIVersion
}

func (r *APIRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		path := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.mux.HandleFunc(
			path,
			route.Handler,
		)
	}
}
