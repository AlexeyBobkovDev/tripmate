package testutils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
)

func SetupHTTP(t *testing.T, method, path string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	req := httptest.NewRequest(method, path, body)
	logger := SetupLogger(t, TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	return rec, req
}
