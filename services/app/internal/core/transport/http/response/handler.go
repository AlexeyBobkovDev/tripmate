package core_http_response

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
)

type HTTPResponseHandler struct {
	rw  http.ResponseWriter
	log *core_logger.Logger
}

func NewHTTPResponseHandler(
	rw http.ResponseWriter,
	log *core_logger.Logger,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		rw:  rw,
		log: log,
	}
}

func (h *HTTPResponseHandler) JSONResponse(
	statusCode int,
	responseBody any,
) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)
	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	var codedError *core_errors.CodedError
	if errors.As(err, &codedError) {
		switch codedError.Category {
		case core_errors.CategoryBadRequest:
			statusCode = http.StatusBadRequest
			logFunc = h.log.Warn
		case core_errors.CategoryNotFound:
			statusCode = http.StatusNotFound
			logFunc = h.log.Debug
		case core_errors.CategoryConflict:
			statusCode = http.StatusConflict
			logFunc = h.log.Warn
		case core_errors.CategoryDatabase:
			statusCode = http.StatusInternalServerError
			logFunc = h.log.Error
		default:
			statusCode = http.StatusInternalServerError
			logFunc = h.log.Error
		}

		logFunc(msg, zap.Error(err))
		h.errorResponse(statusCode, err, msg)
		return
	}
	h.log.Error(
		"INVALID ERROR WHICH CAN NOT BE DECODED TO CodedError error",
		zap.Error(err),
	)
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}

	h.JSONResponse(
		statusCode,
		response,
	)
}
