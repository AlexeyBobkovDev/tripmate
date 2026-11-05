package routes_transport_http

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestRoutesTransportHTTP_GetRoutes_InvalidTripID(t *testing.T) {
	tests := []struct {
		name   string
		tripID string
	}{
		{
			name:   "empty tripID",
			tripID: "",
		},
		{
			name:   "invalid tripID",
			tripID: "abc",
		},
		{
			name:   "zero tripID",
			tripID: "0",
		},
		{
			name:   "negative tripID",
			tripID: "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRoutesService := mocks.NewMockRoutesService(ctrl)
			routesTransport := NewRoutesHTTPHandler(mockRoutesService)

			mockRoutesService.
				EXPECT().
				GetRoutes(gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodGet, "/routes/1", nil)
			req.SetPathValue("tripID", tt.tripID)

			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)
			req = req.WithContext(ctx)

			routesTransport.GetRoutes(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestRoutesTransportHTTP_GetRoutes_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesService := mocks.NewMockRoutesService(ctrl)
	routesTransport := NewRoutesHTTPHandler(mockRoutesService)

	expectedRoute := testRouteResponse()
	expectedRoutes := []*domain.Route{&expectedRoute}

	mockRoutesService.
		EXPECT().
		GetRoutes(
			gomock.Any(),
			1,
		).
		Return(expectedRoutes, nil)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/routes/1", nil)
	req.SetPathValue("tripID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	routesTransport.GetRoutes(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body []RouteResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

	require.Len(t, body, 1)
	assert.Equal(t, testExpectedRouteResponse().ID, body[0].ID)
	assert.Equal(t, testExpectedRouteResponse().TripID, body[0].TripID)
	assert.Equal(t, testExpectedRouteResponse().Version, body[0].Version)
}

func TestRoutesTransportHTTP_GetRoutes_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesService := mocks.NewMockRoutesService(ctrl)
	routesTransport := NewRoutesHTTPHandler(mockRoutesService)

	mockRoutesService.
		EXPECT().
		GetRoutes(
			gomock.Any(),
			1,
		).
		Return(nil, errors.New("service error"))

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/routes/1", nil)
	req.SetPathValue("tripID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	routesTransport.GetRoutes(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
