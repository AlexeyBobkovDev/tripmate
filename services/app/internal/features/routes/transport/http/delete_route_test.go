package routes_transport_http

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestRoutesTransportHTTP_DeleteRoute_InvalidIDs(t *testing.T) {
	tests := []struct {
		name    string
		tripID  string
		routeID string
	}{
		{
			name:    "empty tripID",
			tripID:  "",
			routeID: "1",
		},
		{
			name:    "invalid tripID",
			tripID:  "abc",
			routeID: "1",
		},
		{
			name:    "zero tripID",
			tripID:  "0",
			routeID: "1",
		},
		{
			name:    "negative tripID",
			tripID:  "-1",
			routeID: "1",
		},
		{
			name:    "empty routeID",
			tripID:  "1",
			routeID: "",
		},
		{
			name:    "invalid routeID",
			tripID:  "1",
			routeID: "abc",
		},
		{
			name:    "zero routeID",
			tripID:  "1",
			routeID: "0",
		},
		{
			name:    "negative routeID",
			tripID:  "1",
			routeID: "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRoutesService := mocks.NewMockRoutesService(ctrl)
			routesTransport := NewRoutesHTTPHandler(mockRoutesService)

			mockRoutesService.
				EXPECT().
				DeleteRoute(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/routes/1/1", nil)
			req.SetPathValue("tripID", tt.tripID)
			req.SetPathValue("routeID", tt.routeID)

			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)
			req = req.WithContext(ctx)

			routesTransport.DeleteRoute(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestRoutesTransportHTTP_DeleteRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesService := mocks.NewMockRoutesService(ctrl)
	routesTransport := NewRoutesHTTPHandler(mockRoutesService)

	mockRoutesService.
		EXPECT().
		DeleteRoute(gomock.Any(), 1, 1).
		Return(nil)

	rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/routes/1/1", nil)
	req.SetPathValue("tripID", "1")
	req.SetPathValue("routeID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	routesTransport.DeleteRoute(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestRoutesTransportHTTP_DeleteRoute_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesService := mocks.NewMockRoutesService(ctrl)
	routesTransport := NewRoutesHTTPHandler(mockRoutesService)

	mockRoutesService.
		EXPECT().
		DeleteRoute(gomock.Any(), 1, 1).
		Return(errors.New("service error"))

	rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/routes/1/1", nil)
	req.SetPathValue("tripID", "1")
	req.SetPathValue("routeID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	routesTransport.DeleteRoute(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestRoutesTransportHTTP_DeleteRoute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesService := mocks.NewMockRoutesService(ctrl)
	routesTransport := NewRoutesHTTPHandler(mockRoutesService)

	mockRoutesService.
		EXPECT().
		DeleteRoute(gomock.Any(), 1, 1).
		Return(core_errors.ErrNotFound)

	rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/routes/1/1", nil)
	req.SetPathValue("tripID", "1")
	req.SetPathValue("routeID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	routesTransport.DeleteRoute(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
