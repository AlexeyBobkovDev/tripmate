package core_server

import "net/http"

type Route struct {
	Method  string
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

func (r Route) GetRoutes() []Route {
	return []Route{}
}
