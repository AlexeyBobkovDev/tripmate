//go:build e2e

package trips_e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	trips_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http"
)

func TestTripsTransportHTTP_GetTrips_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	first := createTrip(t, ctx, user.ID, "First E2E Trip")
	second := createTrip(t, ctx, user.ID, "Second E2E Trip")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/trips/%d", baseURL, user.ID),
		nil,
	)
	require.NoError(t, err)

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusOK, response.StatusCode)

	var got []trips_transport_http.TripResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&got))

	require.Len(t, got, 2)

	names := []string{got[0].Name, got[1].Name}
	assert.Contains(t, names, first.Name)
	assert.Contains(t, names, second.Name)

	for _, trip := range got {
		assert.Equal(t, user.ID, trip.OwnerID)
		assert.NotZero(t, trip.ID)
		assert.Nil(t, trip.DeletedAt)
	}
}

func TestTripsTransportHTTP_GetTrips_E2E_EmptyList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/trips/%d", baseURL, user.ID),
		nil,
	)
	require.NoError(t, err)

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusOK, response.StatusCode)

	var got []trips_transport_http.TripResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&got))
	assert.Empty(t, got)
}

func TestTripsTransportHTTP_GetTrips_E2E_BadRequest_InvalidOwnerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL+"/trips/-1",
		nil,
	)
	require.NoError(t, err)

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	errResponse := decodeErrorResponse(t, response)
	require.Contains(t, errResponse.Error, core_errors.ErrBadRequest.Error())
}
