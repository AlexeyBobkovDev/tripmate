//go:build e2e

package trips_e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestTripsTransportHTTP_DeleteTrip_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Trip To Delete")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/trips/%d/%d", baseURL, user.ID, trip.ID),
		nil,
	)
	require.NoError(t, err)

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusNoContent, response.StatusCode)

	getReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/trips/%d/%d", baseURL, user.ID, trip.ID),
		nil,
	)
	require.NoError(t, err)

	getResponse, err := newHTTPClient().Do(getReq)
	require.NoError(t, err)
	defer getResponse.Body.Close()

	require.Equal(t, http.StatusNotFound, getResponse.StatusCode)

	errResponse := decodeErrorResponse(t, getResponse)
	require.Contains(t, errResponse.Error, core_errors.ErrNotFound.Error())
}

func TestTripsTransportHTTP_DeleteTrip_E2E_BadRequest_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/trips/%d/-1", baseURL, user.ID),
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
