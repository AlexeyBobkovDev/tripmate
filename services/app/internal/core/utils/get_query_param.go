package core_utils

import "net/http"

func GetQueryParam(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}
