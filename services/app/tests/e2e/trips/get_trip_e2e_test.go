//go:build e2e

package trips_e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	trips_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http"
)

func TestTripsTransportHTTP_GetTrip_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	createdTrip := createTrip(t, ctx, user.ID, "Get Trip E2E")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/trips/%d/%d", baseURL, user.ID, createdTrip.ID),
		nil,
	)
	require.NoError(t, err)

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusOK, response.StatusCode)

	var got trips_transport_http.TripResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&got))

	require.Equal(t, createdTrip.ID, got.ID)
	require.Equal(t, createdTrip.OwnerID, got.OwnerID)
	require.Equal(t, createdTrip.Version, got.Version)
	require.Equal(t, createdTrip.Name, got.Name)
	require.Equal(t, createdTrip.PlanningArea, got.PlanningArea)
	require.Equal(t, createdTrip.TimeZone, got.TimeZone)
	require.Nil(t, got.DeletedAt)
}

func TestTripsTransportHTTP_GetTrip_E2E_BadRequest_InvalidIDs(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "invalid owner ID",
			url:  baseURL + "/trips/-1/1",
		},
		{
			name: "invalid trip ID",
			url:  baseURL + "/trips/1/-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, tt.url, nil)
			require.NoError(t, err)

			response, err := newHTTPClient().Do(req)
			require.NoError(t, err)
			defer response.Body.Close()

			require.Equal(t, http.StatusBadRequest, response.StatusCode)

			errResponse := decodeErrorResponse(t, response)
			require.Contains(t, errResponse.Error, core_errors.ErrBadRequest.Error())
		})
	}
}

func TestTripsTransportHTTP_GetTrip_E2E_NotFound_WrongOwner(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Private Trip")

	// Use a valid positive ID that does not belong to the trip owner.
	wrongOwnerID := user.ID + 100000

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/trips/%d/%d", baseURL, wrongOwnerID, trip.ID),
		nil,
	)
	require.NoError(t, err)

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusNotFound, response.StatusCode)

	errResponse := decodeErrorResponse(t, response)
	require.Contains(t, errResponse.Error, core_errors.ErrNotFound.Error())
}
