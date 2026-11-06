//go:build e2e

package routes_e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	routes_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/transport/http"
)

func TestRoutesTransportHTTP_GetRoutes_E2E_Empty(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Get routes empty E2E test")

	response := doRequest(
		t,
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/routes/%d", baseURL, trip.ID),
		nil,
	)
	defer response.Body.Close()

	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "[]", strings.TrimSpace(string(body)))

	var routes []routes_transport_http.RouteResponse
	require.NoError(t, json.Unmarshal(body, &routes))
	assert.NotNil(t, routes)
	assert.Empty(t, routes)
}

func TestRoutesTransportHTTP_GetRoutes_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Get routes E2E test")

	firstRoute := createRoute(t, ctx, user.ID, trip.ID)
	secondRoute := createRoute(t, ctx, user.ID, trip.ID)

	routes, statusCode := getRoutes(t, ctx, trip.ID)
	require.Equal(t, http.StatusOK, statusCode)
	require.Len(t, routes, 2)

	gotByID := make(map[int]bool, len(routes))
	for _, route := range routes {
		gotByID[route.ID] = true
		assert.Equal(t, trip.ID, route.TripID)
		assert.Equal(t, trip.Version, route.TripVersion)
	}

	assert.True(t, gotByID[firstRoute.ID])
	assert.True(t, gotByID[secondRoute.ID])
}

func TestRoutesTransportHTTP_GetRoutes_E2E_OnlyReturnsTripRoutes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Get routes owned trip E2E test")
	anotherTrip := createTrip(t, ctx, user.ID, "Get routes another trip E2E test")

	ownedRoute := createRoute(t, ctx, user.ID, trip.ID)
	_ = createRoute(t, ctx, user.ID, anotherTrip.ID)

	routes, statusCode := getRoutes(t, ctx, trip.ID)
	require.Equal(t, http.StatusOK, statusCode)
	require.Len(t, routes, 1)

	assert.Equal(t, ownedRoute.ID, routes[0].ID)
	assert.Equal(t, trip.ID, routes[0].TripID)
}

func TestRoutesTransportHTTP_GetRoutes_E2E_UnknownTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routes, statusCode := getRoutes(t, ctx, 2_000_000_000)

	require.Equal(t, http.StatusOK, statusCode)
	assert.Empty(t, routes)
}

func TestRoutesTransportHTTP_GetRoutes_E2E_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	invalidTripIDs := []string{"abc", "0", "-1"}

	for _, tripID := range invalidTripIDs {
		t.Run(tripID, func(t *testing.T) {
			response := doRequest(
				t,
				ctx,
				http.MethodGet,
				fmt.Sprintf("%s/routes/%s", baseURL, tripID),
				nil,
			)
			defer response.Body.Close()

			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	}
}

func TestRoutesTransportHTTP_GetRoutes_E2E_ReturnsGeneratedRouteContent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Get routes content E2E test")
	createdRoute := createRoute(t, ctx, user.ID, trip.ID)

	routes, statusCode := getRoutes(t, ctx, trip.ID)
	require.Equal(t, http.StatusOK, statusCode)
	require.Len(t, routes, 1)

	require.Equal(t, createdRoute.ID, routes[0].ID)
	assertRouteResponseMatchesGenerated(t, expectedGeneratedRoute(trip), &routes[0], trip)
}
