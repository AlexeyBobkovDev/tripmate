package core_http_response

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestHTTPResponseHandler_JSONResponse_MarshableResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       any
	}{
		{
			name:       "StatusOK with string body",
			statusCode: http.StatusOK,
			body:       "operation was successful",
		},
		{
			name:       "StatusCreated with string body",
			statusCode: http.StatusCreated,
			body:       "the user was created",
		},
		{
			name:       "StatusCreated with body",
			statusCode: http.StatusCreated,
			body: map[string]any{
				"id":      float64(1),
				"message": "the user is successfully created",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseHandler, responseRecorder := SetupResponseHandler(t, "DEBUG", t.TempDir())
			responseHandler.JSONResponse(
				tt.statusCode,
				tt.body,
			)

			assert.Equal(t, tt.statusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

			var response any
			err := json.NewDecoder(responseRecorder.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, tt.body, response)
		})
	}
}

func TestHTTPResponseHandler_JSONResponse_UnmarshableResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       any
	}{
		{
			name:       "unmarshable value",
			statusCode: http.StatusInternalServerError,
			body:       make(chan int),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseHandler, responseRecorder := SetupResponseHandler(t, "DEBUG", t.TempDir())
			responseHandler.JSONResponse(
				tt.statusCode,
				tt.body,
			)

			assert.Equal(t, tt.statusCode, responseRecorder.Code)
			assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

			assert.Empty(t, responseRecorder.Body.String())
		})
	}
}

func TestHTTPResponseHandler_ErrorResponse_UnclassifiedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{
			name: "database error",
			err:  errors.New("database"),
			msg:  "some error occurred in this app",
		},
		{
			name: "validation error",
			err:  errors.New("validation error"),
			msg:  "some error occurred in this app",
		},
	}
	standardErrorResponse := ErrorResponse{
		Error:   errors.New("internal server error").Error(),
		Message: "an unexpected error occurred",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseHandler, responseRecorder := SetupResponseHandler(t, "DEBUG", t.TempDir())
			responseHandler.ErrorResponse(tt.err, tt.msg)

			assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
			assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

			var response ErrorResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, standardErrorResponse, response)
		})
	}
}

func TestHTTPResponseHandler_ErrorResponse_ClassifiedErrors(t *testing.T) {
	tests := []struct {
		name               string
		err                error
		msg                string
		expectedHTTPStatus int
	}{
		{
			name:               "bad request",
			err:                core_errors.NewError(errors.New("bad request"), core_errors.CategoryBadRequest),
			msg:                "bad request",
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name:               "not found",
			err:                core_errors.NewError(errors.New("not found"), core_errors.CategoryNotFound),
			msg:                "not found",
			expectedHTTPStatus: http.StatusNotFound,
		},
		{
			name:               "conflict",
			err:                core_errors.NewError(errors.New("conflict"), core_errors.CategoryConflict),
			msg:                "conflict",
			expectedHTTPStatus: http.StatusConflict,
		},
		{
			name:               "database error",
			err:                core_errors.NewError(errors.New("database error"), core_errors.CategoryDatabase),
			msg:                "database error",
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name:               "unexpected category",
			err:                core_errors.NewError(errors.New("unexpected category"), core_errors.Category(200)),
			msg:                "unexpected category",
			expectedHTTPStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			standardErrorResponse := ErrorResponse{
				Error:   tt.err.Error(),
				Message: tt.msg,
			}
			responseHandler, responseRecorder := SetupResponseHandler(t, "DEBUG", t.TempDir())
			responseHandler.ErrorResponse(tt.err, tt.msg)

			assert.Equal(t, tt.expectedHTTPStatus, responseRecorder.Code)
			assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

			var response ErrorResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, standardErrorResponse, response)
		})
	}
}

func TestHTTPResponseHandler_NoContent(t *testing.T) {
	responseHandler, responseRecorder := SetupResponseHandler(t, "DEBUG", t.TempDir())
	responseHandler.NoContent()

	assert.Equal(t, http.StatusNoContent, responseRecorder.Code)

	assert.Empty(t, responseRecorder.Body.String())

	assert.Equal(t, http.StatusNoContent, responseRecorder.Code)
}

func TestHTTPResponseHandler_errorResponse(t *testing.T) {
	responseHandler, responseRecorder := SetupResponseHandler(t, "DEBUG", t.TempDir())
	responseHandler.errorResponse(
		http.StatusBadRequest,
		errors.New("bad request"),
		"bad request",
	)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "application/json", responseRecorder.Header().Get("Content-Type"))

	var response ErrorResponse
	err := json.NewDecoder(responseRecorder.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, ErrorResponse{
		Error:   errors.New("bad request").Error(),
		Message: "bad request",
	}, response)
}
