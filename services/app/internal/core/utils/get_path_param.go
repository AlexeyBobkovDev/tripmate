package core_utils

import "net/http"

func GetPathParam(r *http.Request, name string) string {
	param := r.PathValue(name)
	return param
}
