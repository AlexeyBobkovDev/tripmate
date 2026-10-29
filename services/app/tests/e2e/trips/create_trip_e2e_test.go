//go:build e2e

package trips_e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestTripsTransportHTTP_CreateTrip_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "E2E Moscow Trip")

	assert.NotZero(t, trip.ID)
	assert.NotZero(t, trip.Version)
	assert.NotZero(t, trip.CreatedAt)
	assert.NotZero(t, trip.UpdatedAt)

	assert.Equal(t, user.ID, trip.OwnerID)
	assert.Equal(t, "E2E Moscow Trip", trip.Name)
	assert.Equal(t, "Moscow", trip.PlanningArea)
	assert.Equal(t, "Europe/Moscow", trip.TimeZone)
	assert.Nil(t, trip.DeletedAt)
}

func TestTripsTransportHTTP_CreateTrip_E2E_BadRequest_InvalidOwnerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	body := newTripRequestBody(-1, "Invalid owner trip")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/trips/%d", baseURL, -1),
		marshalJSON(t, body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	errResponse := decodeErrorResponse(t, response)
	require.Contains(t, errResponse.Error, core_errors.ErrBadRequest.Error())
}

func TestTripsTransportHTTP_CreateTrip_E2E_BadRequest_InvalidJSON(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/trips/%d", baseURL, 1),
		strings.NewReader(`{"owner_id":`),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	errResponse := decodeErrorResponse(t, response)
	require.Contains(t, errResponse.Error, core_errors.ErrInvalidArgument.Error())
}
