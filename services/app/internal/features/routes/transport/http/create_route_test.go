package routes_transport_http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestRoutesTransportHTTP_CreateRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesService := mocks.NewMockRoutesService(ctrl)
	routesTransport := NewRoutesHTTPHandler(mockRoutesService)

	expectedRoute := testRouteResponse()

	mockRoutesService.
		EXPECT().
		CreateRoute(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(&expectedRoute, nil)

	rec, req := testutils.SetupHTTP(t, http.MethodPost, "/routes/1/1", nil)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	routesTransport.CreateRoute(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var route RouteResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&route))

	assert.Equal(t, testExpectedRouteResponse().ID, route.ID)
	assert.Equal(t, testExpectedRouteResponse().TripID, route.TripID)
	assert.Equal(t, testExpectedRouteResponse().TripVersion, route.TripVersion)
	assert.Equal(t, testExpectedRouteResponse().Version, route.Version)
	assert.Equal(t, testExpectedRouteResponse().Status, route.Status)
	assert.True(t, testExpectedRouteResponse().CreatedAt.Equal(route.CreatedAt))
	assert.True(t, testExpectedRouteResponse().UpdatedAt.Equal(route.UpdatedAt))
}

func TestRoutesTransportHTTP_CreateRoute_ServiceError(t *testing.T) {
	serviceError := core_postgres_pool_errors.ErrViolatedForeignKey

	ctrl := gomock.NewController(t)
	mockRoutesService := mocks.NewMockRoutesService(ctrl)
	routesTransport := NewRoutesHTTPHandler(mockRoutesService)

	mockRoutesService.
		EXPECT().
		CreateRoute(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, serviceError)

	rec, req := testutils.SetupHTTP(t, http.MethodPost, "/routes/1/1", nil)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	routesTransport.CreateRoute(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var returnedServiceError core_http_response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&returnedServiceError); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "failed to create route", returnedServiceError.Message)
	assert.Equal(t, serviceError.Error(), returnedServiceError.Error)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestRoutesTransportHTTP_CreateRoute_InvalidIDs(t *testing.T) {
	tests := []struct {
		name    string
		ownerID string
		tripID  string
	}{
		{
			name:    "empty ownerID",
			ownerID: "",
			tripID:  "1",
		},
		{
			name:    "invalid ownerID",
			ownerID: "abc",
			tripID:  "1",
		},
		{
			name:    "zero ownerID",
			ownerID: "0",
			tripID:  "1",
		},
		{
			name:    "negative ownerID",
			ownerID: "-1",
			tripID:  "1",
		},
		{
			name:    "empty tripID",
			ownerID: "1",
			tripID:  "",
		},
		{
			name:    "invalid tripID",
			ownerID: "1",
			tripID:  "abc",
		},
		{
			name:    "zero tripID",
			ownerID: "1",
			tripID:  "0",
		},
		{
			name:    "negative tripID",
			ownerID: "1",
			tripID:  "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRoutesService := mocks.NewMockRoutesService(ctrl)
			routesTransport := NewRoutesHTTPHandler(mockRoutesService)

			mockRoutesService.
				EXPECT().
				CreateRoute(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodPost, "/routes/1/1", nil)

			req.SetPathValue("ownerID", tt.ownerID)
			req.SetPathValue("tripID", tt.tripID)

			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)
			req = req.WithContext(ctx)

			routesTransport.CreateRoute(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		})
	}
}
