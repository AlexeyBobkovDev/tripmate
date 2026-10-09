package core_http_response

import (
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/testify/v2/require"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
)

func SetupResponseHandler(t *testing.T, loggerLevel, loggerFolder string) (*HTTPResponseHandler, *httptest.ResponseRecorder) {
	logger, err := core_logger.NewLogger(core_logger.Config{
		Level:  loggerLevel,
		Folder: loggerFolder,
	})
	require.NoError(t, err)
	responseRecorder := httptest.NewRecorder()
	responseHandler := NewHTTPResponseHandler(responseRecorder, logger)
	return responseHandler, responseRecorder
}
