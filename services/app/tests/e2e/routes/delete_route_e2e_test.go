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

func TestRoutesTransportHTTP_DeleteRoute_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Delete route E2E test")
	route := createRoute(t, ctx, user.ID, trip.ID)

	statusCode := deleteRoute(t, ctx, trip.ID, route.ID)
	require.Equal(t, http.StatusNoContent, statusCode)

	retrievedRoute, statusCode := getRoute(t, ctx, trip.ID, route.ID)
	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Nil(t, retrievedRoute)

	routes, statusCode := getRoutes(t, ctx, trip.ID)
	require.Equal(t, http.StatusOK, statusCode)
	assert.Empty(t, routes)
}

func TestRoutesTransportHTTP_DeleteRoute_E2E_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	invalidTripIDs := []string{"abc", "0", "-1"}

	for _, tripID := range invalidTripIDs {
		t.Run(tripID, func(t *testing.T) {
			response := doRequest(
				t,
				ctx,
				http.MethodDelete,
				fmt.Sprintf("%s/routes/%s/1", baseURL, tripID),
				nil,
			)
			defer response.Body.Close()

			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	}
}

func TestRoutesTransportHTTP_DeleteRoute_E2E_InvalidRouteID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Delete route invalid route ID E2E test")

	invalidRouteIDs := []string{"abc", "0", "-1"}

	for _, routeID := range invalidRouteIDs {
		t.Run(routeID, func(t *testing.T) {
			response := doRequest(
				t,
				ctx,
				http.MethodDelete,
				fmt.Sprintf("%s/routes/%d/%s", baseURL, trip.ID, routeID),
				nil,
			)
			defer response.Body.Close()

			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	}
}

func TestRoutesTransportHTTP_DeleteRoute_E2E_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Delete route not found E2E test")

	statusCode := deleteRoute(t, ctx, trip.ID, 2_000_000_000)

	assert.Equal(t, http.StatusNotFound, statusCode)
}

func TestRoutesTransportHTTP_DeleteRoute_E2E_AlreadyDeleted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Delete route already deleted E2E test")
	route := createRoute(t, ctx, user.ID, trip.ID)

	require.Equal(t, http.StatusNoContent, deleteRoute(t, ctx, trip.ID, route.ID))
	assert.Equal(t, http.StatusNotFound, deleteRoute(t, ctx, trip.ID, route.ID))
}

func TestRoutesTransportHTTP_DeleteRoute_E2E_WrongTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Delete route wrong trip E2E test")
	anotherTrip := createTrip(t, ctx, user.ID, "Delete route another trip E2E test")
	route := createRoute(t, ctx, user.ID, trip.ID)

	statusCode := deleteRoute(t, ctx, anotherTrip.ID, route.ID)
	require.Equal(t, http.StatusNotFound, statusCode)

	retrievedRoute, statusCode := getRoute(t, ctx, trip.ID, route.ID)
	require.Equal(t, http.StatusOK, statusCode)
	require.NotNil(t, retrievedRoute)
	assert.Equal(t, route.ID, retrievedRoute.ID)
}
