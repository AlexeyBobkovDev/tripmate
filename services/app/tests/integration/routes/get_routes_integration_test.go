//go:build integration

package routes_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestRoutesService_GetRoutes_Integration_Empty(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Get routes empty integration test")

	routes, err := testDeps.routesService.GetRoutes(ctx, trip.ID)
	require.NoError(t, err)
	assert.Empty(t, routes)
}

func TestRoutesService_GetRoutes_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Get routes integration test")

	firstRoute := createTestRoute(t, ctx, trip)
	secondRoute := createTestRoute(t, ctx, trip)

	routes, err := testDeps.routesService.GetRoutes(ctx, trip.ID)
	require.NoError(t, err)
	require.NotNil(t, routes)
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

func TestRoutesService_GetRoutes_Integration_OnlyReturnsTripRoutes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Get routes owned trip integration test")
	anotherTrip := createTestTrip(t, ctx, user.ID, "Get routes another trip integration test")

	ownedRoute := createTestRoute(t, ctx, trip)
	_ = createTestRoute(t, ctx, anotherTrip)

	routes, err := testDeps.routesService.GetRoutes(ctx, trip.ID)
	require.NoError(t, err)
	require.NotNil(t, routes)
	require.Len(t, routes, 1)

	assert.Equal(t, ownedRoute.ID, routes[0].ID)
	assert.Equal(t, trip.ID, routes[0].TripID)
}

func TestRoutesService_GetRoutes_Integration_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routes, err := testDeps.routesService.GetRoutes(ctx, -1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, routes)
}

func TestRoutesService_GetRoutes_Integration_UnknownTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routes, err := testDeps.routesService.GetRoutes(ctx, 2_000_000_000)
	require.NoError(t, err)
	assert.Empty(t, routes)
}

func TestRoutesService_GetRoutes_Integration_ReturnsGeneratedRouteContent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Get routes content integration test")

	createdRoute := createTestRoute(t, ctx, trip)

	routes, err := testDeps.routesService.GetRoutes(ctx, trip.ID)
	require.NoError(t, err)
	require.Len(t, routes, 1)

	require.Equal(t, createdRoute.ID, routes[0].ID)
	assertRouteMatchesGenerated(
		t,
		testutils.RouteGeneratorResponseTest(trip.StartDate),
		routes[0],
		trip,
	)
}
