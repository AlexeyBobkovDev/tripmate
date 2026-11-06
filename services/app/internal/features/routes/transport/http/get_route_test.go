package routes_transport_http

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestRoutesTransportHTTP_GetRoute_InvalidIDs(t *testing.T) {
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
				GetRoute(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodGet, "/routes/1/1", nil)
			req.SetPathValue("tripID", tt.tripID)
			req.SetPathValue("routeID", tt.routeID)

			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)

			req = req.WithContext(ctx)

			routesTransport.GetRoute(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestRoutesTransportHTTP_GetRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesService := mocks.NewMockRoutesService(ctrl)
	routesTransport := NewRoutesHTTPHandler(mockRoutesService)

	expectedRoute := testRouteResponse()

	mockRoutesService.
		EXPECT().
		GetRoute(
			gomock.Any(),
			1,
			1,
		).
		Return(&expectedRoute, nil)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/routes/1/1", nil)
	req.SetPathValue("tripID", "1")
	req.SetPathValue("routeID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	routesTransport.GetRoute(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body RouteResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

	assert.Equal(t, testExpectedRouteResponse().ID, body.ID)
	assert.Equal(t, testExpectedRouteResponse().TripID, body.TripID)
	assert.Equal(t, testExpectedRouteResponse().Version, body.Version)
}

func TestRoutesTransportHTTP_GetRoute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesService := mocks.NewMockRoutesService(ctrl)
	routesTransport := NewRoutesHTTPHandler(mockRoutesService)

	mockRoutesService.
		EXPECT().
		GetRoute(
			gomock.Any(),
			1,
			1,
		).
		Return(nil, errors.Join(pgx.ErrNoRows, core_errors.ErrNotFound))

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/routes/1/1", nil)
	req.SetPathValue("tripID", "1")
	req.SetPathValue("routeID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	routesTransport.GetRoute(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
