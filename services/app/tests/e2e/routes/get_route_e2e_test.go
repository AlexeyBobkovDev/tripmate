//go:build e2e

package routes_e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutesTransportHTTP_GetRoute_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Get route E2E test")
	createdRoute := createRoute(t, ctx, user.ID, trip.ID)

	retrievedRoute, statusCode := getRoute(t, ctx, trip.ID, createdRoute.ID)
	require.Equal(t, http.StatusOK, statusCode)
	require.NotNil(t, retrievedRoute)

	assert.Equal(t, createdRoute.ID, retrievedRoute.ID)
	assert.Equal(t, createdRoute.TripID, retrievedRoute.TripID)
	assert.Equal(t, createdRoute.TripVersion, retrievedRoute.TripVersion)
	assert.Equal(t, createdRoute.Version, retrievedRoute.Version)
	assert.Equal(t, createdRoute.Status, retrievedRoute.Status)
	assert.Equal(t, createdRoute.CreatedAt, retrievedRoute.CreatedAt)
	assert.Equal(t, createdRoute.UpdatedAt, retrievedRoute.UpdatedAt)
	assertRouteResponseMatchesGenerated(t, expectedGeneratedRoute(trip), retrievedRoute, trip)
}

func TestRoutesTransportHTTP_GetRoute_E2E_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	invalidTripIDs := []string{"abc", "0", "-1"}

	for _, tripID := range invalidTripIDs {
		t.Run(tripID, func(t *testing.T) {
			response := doRequest(
				t,
				ctx,
				http.MethodGet,
				fmt.Sprintf("%s/routes/%s/1", baseURL, tripID),
				nil,
			)
			defer response.Body.Close()

			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	}
}

func TestRoutesTransportHTTP_GetRoute_E2E_InvalidRouteID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Get route invalid route ID E2E test")

	invalidRouteIDs := []string{"abc", "0", "-1"}

	for _, routeID := range invalidRouteIDs {
		t.Run(routeID, func(t *testing.T) {
			response := doRequest(
				t,
				ctx,
				http.MethodGet,
				fmt.Sprintf("%s/routes/%d/%s", baseURL, trip.ID, routeID),
				nil,
			)
			defer response.Body.Close()

			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	}
}

func TestRoutesTransportHTTP_GetRoute_E2E_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Get route not found E2E test")

	route, statusCode := getRoute(t, ctx, trip.ID, 2_000_000_000)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Nil(t, route)
}

func TestRoutesTransportHTTP_GetRoute_E2E_WrongTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Get route wrong trip E2E test")
	anotherTrip := createTrip(t, ctx, user.ID, "Get route another trip E2E test")
	createdRoute := createRoute(t, ctx, user.ID, trip.ID)

	route, statusCode := getRoute(t, ctx, anotherTrip.ID, createdRoute.ID)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Nil(t, route)
}

func TestRoutesTransportHTTP_GetRoute_E2E_UnknownTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	route, statusCode := getRoute(t, ctx, 2_000_000_000, 1)

	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Nil(t, route)
}
