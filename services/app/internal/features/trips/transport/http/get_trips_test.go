package trips_transport_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsTransportHTTP_GetTrips_InvalidOwnerID(t *testing.T) {
	tests := []struct {
		name    string
		ownerID string
	}{
		{
			name:    "empty ownerID",
			ownerID: "",
		},
		{
			name:    "invalid ownerID",
			ownerID: "abc",
		},
		{
			name:    "zero ownerID",
			ownerID: "0",
		},
		{
			name:    "negative ownerID",
			ownerID: "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := mocks.NewMockTripsService(ctrl)
			handler := NewTripsHTTPHandler(mockService)

			mockService.EXPECT().
				GetTrips(gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodGet, fmt.Sprintf("/trips/%s", tt.ownerID), nil)
			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)

			req.SetPathValue("ownerID", tt.ownerID)
			req = req.WithContext(ctx)

			handler.GetTrips(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestTripsTransportHTTP_GetTrips_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	expectedTrip := testTripResponse()
	expectedTrips := []domain.Trip{expectedTrip}

	mockService.EXPECT().
		GetTrips(gomock.Any(), 1).
		Return(&expectedTrips, nil)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/trips/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req = req.WithContext(ctx)

	handler.GetTrips(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body []TripResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

	assert.Len(t, body, 1)
	assert.Equal(t, testExpectedTripResponse().ID, body[0].ID)
	assert.Equal(t, testExpectedTripResponse().OwnerID, body[0].OwnerID)
	assert.Equal(t, testExpectedTripResponse().Name, body[0].Name)
}

func TestTripsTransportHTTP_GetTrips_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	mockService.EXPECT().
		GetTrips(gomock.Any(), 1).
		Return(nil, errors.New("service error"))

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/trips/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req = req.WithContext(ctx)

	handler.GetTrips(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
