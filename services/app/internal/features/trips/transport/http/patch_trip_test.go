package trips_transport_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsTransportHTTP_PatchTrip_InvalidIDs(t *testing.T) {
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
				PatchTrip(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodPatch, fmt.Sprintf("/trips/%s/%s", tt.ownerID, tt.tripID), nil)
			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)

			req.SetPathValue("ownerID", tt.ownerID)
			req.SetPathValue("tripID", tt.tripID)
			req = req.WithContext(ctx)

			handler.PatchTrip(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestTripsTransportHTTP_PatchTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	tripResp := testTripResponse()

	mockService.EXPECT().
		PatchTrip(
			gomock.Any(),
			1,
			1,
			gomock.Any(),
		).
		Return(&tripResp, nil)

	body := `{
		"version": 1,
		"name": "Updated Trip"
	}`

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/trips/1/1", strings.NewReader(body))
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")
	req = req.WithContext(ctx)

	handler.PatchTrip(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var actual TripResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&actual))
	assert.Equal(t, tripResp.ID, actual.ID)
}

func TestTripsTransportHTTP_PatchTrip_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	mockService.EXPECT().
		PatchTrip(
			gomock.Any(),
			1,
			1,
			gomock.Any(),
		).
		Return(nil, errors.New("some service error"))

	body := `{
		"version": 1,
		"name": "Updated Trip"
	}`

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/trips/1/1", strings.NewReader(body))
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")
	req = req.WithContext(ctx)

	handler.PatchTrip(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestTripsTransportHTTP_PatchTrip_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	conflictErr := errors.Join(errors.New("version conflict"), core_errors.ErrConflict)
	mockService.EXPECT().
		PatchTrip(
			gomock.Any(),
			1,
			1,
			gomock.Any(),
		).
		Return(nil, conflictErr)

	body := `{
		"version": 1,
		"name": "Updated Trip"
	}`

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/trips/1/1", strings.NewReader(body))
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")
	req = req.WithContext(ctx)

	handler.PatchTrip(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var errResp core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&errResp))
	assert.Equal(t, conflictErr.Error(), errResp.Error)
}

func TestTripsTransportHTTP_PatchTrip_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	notFoundErr := errors.Join(errors.New("not found"), core_errors.ErrNotFound)
	mockService.EXPECT().
		PatchTrip(
			gomock.Any(),
			1,
			1,
			gomock.Any(),
		).
		Return(nil, notFoundErr)

	body := `{
		"version": 1,
		"name": "Updated Trip"
	}`

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/trips/1/1", strings.NewReader(body))
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")
	req = req.WithContext(ctx)

	handler.PatchTrip(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&errResp))
	assert.Equal(t, notFoundErr.Error(), errResp.Error)
}
