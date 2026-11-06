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

func TestRoutesService_GetRoute_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Get route integration test")
	createdRoute := createTestRoute(t, ctx, trip)

	retrievedRoute, err := testDeps.routesService.GetRoute(ctx, trip.ID, createdRoute.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedRoute)

	assert.Equal(t, createdRoute.ID, retrievedRoute.ID)
	assert.Equal(t, createdRoute.TripID, retrievedRoute.TripID)
	assert.Equal(t, createdRoute.TripVersion, retrievedRoute.TripVersion)
	assert.Equal(t, createdRoute.Version, retrievedRoute.Version)
	assert.Equal(t, createdRoute.Status, retrievedRoute.Status)
	assert.Equal(t, createdRoute.CreatedAt, retrievedRoute.CreatedAt)
	assert.Equal(t, createdRoute.UpdatedAt, retrievedRoute.UpdatedAt)
	assertRouteMatchesGenerated(
		t,
		testutils.RouteGeneratorResponseTest(trip.StartDate),
		retrievedRoute,
		trip,
	)
}

func TestRoutesService_GetRoute_Integration_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	route, err := testDeps.routesService.GetRoute(ctx, -1, 1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, route)
}

func TestRoutesService_GetRoute_Integration_InvalidRouteID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Get route invalid route id integration test")

	route, err := testDeps.routesService.GetRoute(ctx, trip.ID, -1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, route)
}

func TestRoutesService_GetRoute_Integration_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Get route not found integration test")

	route, err := testDeps.routesService.GetRoute(ctx, trip.ID, 2_000_000_000)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, route)
}

func TestRoutesService_GetRoute_Integration_WrongTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Get route wrong trip integration test")
	anotherTrip := createTestTrip(t, ctx, user.ID, "Get route another trip integration test")

	createdRoute := createTestRoute(t, ctx, trip)

	route, err := testDeps.routesService.GetRoute(ctx, anotherTrip.ID, createdRoute.ID)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, route)
}
