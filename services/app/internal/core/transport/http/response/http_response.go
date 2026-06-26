package core_http_response

import (
	"encoding/json"
	"net/http"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(
	logger *core_logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: logger,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.rw.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(h.rw).Encode(&responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(msg string, err error, statusCode int) {
	h.rw.WriteHeader(statusCode)

	response := map[string]error{
		msg: err,
	}
	if err := json.NewEncoder(h.rw).Encode(&response); err != nil {
		h.log.Error("write error response", zap.Error(err))
	}
}
