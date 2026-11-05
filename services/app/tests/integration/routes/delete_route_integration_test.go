//go:build integration

package routes_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestRoutesService_DeleteRoute_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Delete route integration test")
	route := createTestRoute(t, ctx, trip)

	err := testDeps.routesService.DeleteRoute(ctx, trip.ID, route.ID)
	require.NoError(t, err)

	retrievedRoute, err := testDeps.routesService.GetRoute(ctx, trip.ID, route.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedRoute)
}

func TestRoutesService_DeleteRoute_Integration_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := testDeps.routesService.DeleteRoute(ctx, -1, 1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
}

func TestRoutesService_DeleteRoute_Integration_InvalidRouteID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Delete route invalid route id integration test")

	err := testDeps.routesService.DeleteRoute(ctx, trip.ID, -1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
}

func TestRoutesService_DeleteRoute_Integration_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Delete route not found integration test")

	err := testDeps.routesService.DeleteRoute(ctx, trip.ID, 2_000_000_000)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestRoutesService_DeleteRoute_Integration_AlreadyDeleted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Delete route already deleted integration test")
	route := createTestRoute(t, ctx, trip)

	require.NoError(t, testDeps.routesService.DeleteRoute(ctx, trip.ID, route.ID))

	err := testDeps.routesService.DeleteRoute(ctx, trip.ID, route.ID)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestRoutesService_DeleteRoute_Integration_WrongTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Delete route wrong trip integration test")
	anotherTrip := createTestTrip(t, ctx, user.ID, "Delete route another trip integration test")

	route := createTestRoute(t, ctx, trip)

	err := testDeps.routesService.DeleteRoute(ctx, anotherTrip.ID, route.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)

	retrievedRoute, err := testDeps.routesService.GetRoute(ctx, trip.ID, route.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedRoute)
	require.Equal(t, route.ID, retrievedRoute.ID)
}
