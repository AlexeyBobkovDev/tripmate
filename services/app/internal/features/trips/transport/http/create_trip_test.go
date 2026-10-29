package trips_transport_http

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsTransportHTTP_CreateTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTripsService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockTripsService)

	tripResponse := testTripResponse()

	mockTripsService.
		EXPECT().
		CreateTrip(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(&tripResponse, nil)

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodPost,
		"/trips/1",
		testTripRequestBodyWithDefaults(),
	)

	req.SetPathValue("ownerID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var trip TripResponse
	if err := json.NewDecoder(rec.Body).Decode(&trip); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testExpectedTripResponse().ID, trip.ID)
	assert.Equal(t, testExpectedTripResponse().OwnerID, trip.OwnerID)
	assert.Equal(t, testExpectedTripResponse().Version, trip.Version)
	assert.Equal(t, testExpectedTripResponse().Name, trip.Name)
	assert.Equal(t, testExpectedTripResponse().PlanningArea, trip.PlanningArea)
	assert.Equal(t, testExpectedTripResponse().StartLocation, trip.StartLocation)
	assert.True(t, testExpectedTripResponse().StartDate.Equal(trip.StartDate))
	assert.True(t, testExpectedTripResponse().EndDate.Equal(trip.EndDate))
	assert.Equal(t, testExpectedTripResponse().Budget, trip.Budget)
	assert.Equal(t, testExpectedTripResponse().Interests, trip.Interests)
	assert.Equal(t, testExpectedTripResponse().TimeZone, trip.TimeZone)
	assert.True(t, testExpectedTripResponse().CreatedAt.Equal(trip.CreatedAt))
	assert.True(t, testExpectedTripResponse().UpdatedAt.Equal(trip.UpdatedAt))
	assert.Equal(t, testExpectedTripResponse().DeletedAt, trip.DeletedAt)
}

func TestTripsTransportHTTP_CreateTrip_ServiceError(t *testing.T) {
	serviceError := core_postgres_pool_errors.ErrViolatedForeignKey

	ctrl := gomock.NewController(t)
	mockTripsService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockTripsService)

	mockTripsService.
		EXPECT().
		CreateTrip(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, serviceError)

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodPost,
		"/trips/1",
		testTripRequestBodyWithDefaults(),
	)

	req.SetPathValue("ownerID", "1")

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var returnedServiceError core_http_response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&returnedServiceError); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "failed to create trip", returnedServiceError.Message)
	assert.Equal(t, core_postgres_pool_errors.ErrViolatedForeignKey.Error(), returnedServiceError.Error)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestTripsTransportHTTP_CreateTrip_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTripsService := mocks.NewMockTripsService(ctrl)
	tripsTransport := NewTripsHTTPHandler(mockTripsService)

	mockTripsService.
		EXPECT().
		CreateTrip(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodPost,
		"/trips",
		strings.NewReader(`invalid json`),
	)

	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	tripsTransport.CreateTrip(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestTripsTransportHTTP_CreateTrip_ValidationError(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "empty name",
			body: `{
				"owner_id": 1,
				"name": "",
				"planning_area": "Moscow",
				"start_location": {"name": "Red Square", "latitude": 55.7539, "longitude": 37.6208},
				"start_date": "2026-11-14",
				"end_date": "2026-11-16",
				"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
				"interests": ["history"],
				"time_zone": "Europe/Moscow"
			}`,
		},
		{
			name: "invalid owner ID",
			body: `{
				"owner_id": 0,
				"name": "Trip",
				"planning_area": "Moscow",
				"start_location": {"name": "Red Square", "latitude": 55.7539, "longitude": 37.6208},
				"start_date": "2026-11-14",
				"end_date": "2026-11-16",
				"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
				"interests": ["history"],
				"time_zone": "Europe/Moscow"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTripsService := mocks.NewMockTripsService(ctrl)
			tripsTransport := NewTripsHTTPHandler(mockTripsService)

			mockTripsService.
				EXPECT().
				CreateTrip(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(
				t,
				http.MethodPost,
				"/trips",
				strings.NewReader(tt.body),
			)

			logger := testutils.SetupLogger(
				t,
				testutils.TestLevel,
				t.TempDir(),
			)

			ctx := core_logger.ToContext(req.Context(), logger)
			req = req.WithContext(ctx)

			tripsTransport.CreateTrip(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		})
	}
}
