package core_server

import (
	"net/http"
	"testing"
)

func TestRoute_Fields(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	route := &Route{
		Method:  http.MethodGet,
		Path:    "/trips/{ownerID}",
		Handler: handler,
	}

	if route.Method != http.MethodGet {
		t.Fatalf("route.Method = %q, want %q", route.Method, http.MethodGet)
	}
	if route.Path != "/trips/{ownerID}" {
		t.Fatalf("route.Path = %q, want %q", route.Path, "/trips/{ownerID}")
	}
	if route.Handler == nil {
		t.Fatal("route.Handler must not be nil")
	}
}
