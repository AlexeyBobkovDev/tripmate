package trips_transport_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsTransportHTTP_GetTrip_InvalidIDs(t *testing.T) {
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
			mockService := mocks.NewMockTripsService(ctrl)
			handler := NewTripsHTTPHandler(mockService)

			mockService.EXPECT().
				GetTrip(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodGet, fmt.Sprintf("/trips/%s/%s", tt.ownerID, tt.tripID), nil)
			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)

			req.SetPathValue("ownerID", tt.ownerID)
			req.SetPathValue("tripID", tt.tripID)
			req = req.WithContext(ctx)

			handler.GetTrip(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestTripsTransportHTTP_GetTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	expectedTrip := testTripResponse()

	mockService.EXPECT().
		GetTrip(gomock.Any(), 1, 1).
		Return(&expectedTrip, nil)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/trips/1/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")
	req = req.WithContext(ctx)

	handler.GetTrip(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body TripResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

	assert.Equal(t, testExpectedTripResponse().ID, body.ID)
	assert.Equal(t, testExpectedTripResponse().OwnerID, body.OwnerID)
	assert.Equal(t, testExpectedTripResponse().Name, body.Name)
}

func TestTripsTransportHTTP_GetTrip_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	mockService.EXPECT().
		GetTrip(gomock.Any(), 1, 1).
		Return(nil, errors.Join(pgx.ErrNoRows, core_errors.ErrNotFound))

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/trips/1/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")
	req = req.WithContext(ctx)

	handler.GetTrip(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
