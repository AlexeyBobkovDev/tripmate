package trips_transport_http

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsTransportHTTP_DeleteTrip_InvalidIDs(t *testing.T) {
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
				DeleteTrip(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodDelete, fmt.Sprintf("/trips/%s/%s", tt.ownerID, tt.tripID), nil)
			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)

			req.SetPathValue("ownerID", tt.ownerID)
			req.SetPathValue("tripID", tt.tripID)
			req = req.WithContext(ctx)

			handler.DeleteTrip(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestTripsTransportHTTP_DeleteTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	mockService.EXPECT().
		DeleteTrip(gomock.Any(), 1, 1).
		Return(nil)

	rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/trips/1/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")
	req = req.WithContext(ctx)

	handler.DeleteTrip(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestTripsTransportHTTP_DeleteTrip_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	mockService.EXPECT().
		DeleteTrip(gomock.Any(), 1, 1).
		Return(errors.New("service error"))

	rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/trips/1/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")
	req = req.WithContext(ctx)

	handler.DeleteTrip(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestTripsTransportHTTP_DeleteTrip_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTripsService(ctrl)
	handler := NewTripsHTTPHandler(mockService)

	mockService.EXPECT().
		DeleteTrip(gomock.Any(), 1, 1).
		Return(core_errors.ErrNotFound)

	rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/trips/1/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("ownerID", "1")
	req.SetPathValue("tripID", "1")
	req = req.WithContext(ctx)

	handler.DeleteTrip(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
