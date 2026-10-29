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

func TestTripsTransportHTTP_PatchTrip_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	createdTrip := createTrip(t, ctx, user.ID, "Original Trip Name")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/trips/%d/%d", baseURL, user.ID, createdTrip.ID),
		newTripPatchBody(t, user.ID, createdTrip.Version, "Updated Trip Name"),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusOK, response.StatusCode)

	var patched trips_transport_http.TripResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&patched))

	assert.Equal(t, createdTrip.ID, patched.ID)
	assert.Equal(t, createdTrip.OwnerID, patched.OwnerID)
	assert.Equal(t, createdTrip.Version+1, patched.Version)
	assert.Equal(t, "Updated Trip Name", patched.Name)
	assert.Equal(t, createdTrip.PlanningArea, patched.PlanningArea)
	assert.Nil(t, patched.DeletedAt)
}

func TestTripsTransportHTTP_PatchTrip_E2E_Conflict_MismatchingVersions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Trip With Version")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/trips/%d/%d", baseURL, user.ID, trip.ID),
		newTripPatchBody(t, user.ID, trip.Version+1, "Should Not Be Applied"),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusConflict, response.StatusCode)

	errResponse := decodeErrorResponse(t, response)
	require.Contains(t, errResponse.Error, core_errors.ErrConflict.Error())
}

func TestTripsTransportHTTP_PatchTrip_E2E_BadRequest_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/trips/%d/-1", baseURL, user.ID),
		newTripPatchBody(t, user.ID, 1, "Invalid ID"),
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

func TestTripsTransportHTTP_PatchTrip_E2E_NotFound_DeletedTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Trip To Delete")

	deleteReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/trips/%d/%d", baseURL, user.ID, trip.ID),
		nil,
	)
	require.NoError(t, err)

	deleteResponse, err := newHTTPClient().Do(deleteReq)
	require.NoError(t, err)
	defer deleteResponse.Body.Close()
	require.Equal(t, http.StatusNoContent, deleteResponse.StatusCode)

	time.Sleep(50 * time.Millisecond)

	patchReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/trips/%d/%d", baseURL, user.ID, trip.ID),
		newTripPatchBody(t, user.ID, trip.Version, "Cannot Update Deleted Trip"),
	)
	require.NoError(t, err)
	patchReq.Header.Set("Content-Type", "application/json")

	patchResponse, err := newHTTPClient().Do(patchReq)
	require.NoError(t, err)
	defer patchResponse.Body.Close()

	require.Equal(t, http.StatusNotFound, patchResponse.StatusCode)

	errResponse := decodeErrorResponse(t, patchResponse)
	require.Contains(t, errResponse.Error, core_errors.ErrNotFound.Error())
}
