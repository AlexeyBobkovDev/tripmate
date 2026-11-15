//go:build integration

package routes_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	trips_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/repository/postgres"
)

func TestRoutesRepository_DeleteRoute_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreateFixture(2))
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	routeDayIDs := make([]int, len(createdRoute.Days))
	for ind, day := range createdRoute.Days {
		routeDayIDs[ind] = day.ID
	}

	err = routesRepository.DeleteRoute(ctx, trip.ID, createdRoute.ID)
	require.NoError(t, err)

	var routeCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM app.routes WHERE id=$1;`,
		createdRoute.ID,
	).Scan(&routeCounter))
	assert.Zero(t, routeCounter)

	for _, routeDayID := range routeDayIDs {
		var dayCounter int
		require.NoError(t, poolForIntegrationTesting.QueryRow(
			ctx,
			`SELECT COUNT(*) FROM app.route_days WHERE id=$1;`,
			routeDayID,
		).Scan(&dayCounter))
		assert.Zero(t, dayCounter)

		var stopCounter int
		require.NoError(t, poolForIntegrationTesting.QueryRow(
			ctx,
			`SELECT COUNT(*) FROM app.route_stops WHERE route_day_id=$1;`,
			routeDayID,
		).Scan(&stopCounter))
		assert.Zero(t, stopCounter)

		var legCounter int
		require.NoError(t, poolForIntegrationTesting.QueryRow(
			ctx,
			`SELECT COUNT(*) FROM app.route_legs WHERE route_day_id=$1;`,
			routeDayID,
		).Scan(&legCounter))
		assert.Zero(t, legCounter)
	}

	retrievedRoute, err := routesRepository.GetRoute(ctx, trip.ID, createdRoute.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, retrievedRoute)

	routes, err := routesRepository.GetRoutes(ctx, trip.ID)
	require.NoError(t, err)
	assert.Empty(t, routes)

	require.ErrorIs(
		t,
		routesRepository.DeleteRoute(ctx, trip.ID, createdRoute.ID),
		core_errors.ErrNotFound,
	)
}

func TestRoutesRepository_DeleteRoute_Integration_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)

	err := routesRepository.DeleteRoute(ctx, 1, 999999)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestRoutesRepository_DeleteRoute_Integration_OnlySelectedRouteDeleted(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)

	routeToDelete, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreateFixture(2))
	require.NoError(t, err)
	require.NotNil(t, routeToDelete)

	routeToKeepCreate := routeCreateFixture(1)
	routeToKeep, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeToKeepCreate)
	require.NoError(t, err)
	require.NotNil(t, routeToKeep)

	require.NoError(t, routesRepository.DeleteRoute(ctx, trip.ID, routeToDelete.ID))

	deletedRoute, err := routesRepository.GetRoute(ctx, trip.ID, routeToDelete.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, deletedRoute)

	keptRoute, err := routesRepository.GetRoute(ctx, trip.ID, routeToKeep.ID)
	require.NoError(t, err)
	require.NotNil(t, keptRoute)
	assertRouteMatchesCreate(t, routeToKeepCreate, keptRoute, trip)
}

func TestRoutesRepository_DeleteRoute_Integration_WrongTrip(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)
	otherTrip := newTripForRoutes(t, ctx)

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreateFixture(1))
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	err = routesRepository.DeleteRoute(ctx, otherTrip.ID, createdRoute.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)

	retrievedRoute, err := routesRepository.GetRoute(ctx, trip.ID, createdRoute.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedRoute)
	assert.Equal(t, createdRoute.ID, retrievedRoute.ID)
}

func TestRoutesRepository_DeleteRoute_Integration_TripDeleted(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	tripsRepository := trips_postgres_repository.NewTripsRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreateFixture(1))
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	require.NoError(t, tripsRepository.DeleteTrip(ctx, trip.OwnerID, trip.ID))

	retrievedRoute, err := routesRepository.GetRoute(ctx, trip.ID, createdRoute.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, retrievedRoute)

	require.ErrorIs(
		t,
		routesRepository.DeleteRoute(ctx, trip.ID, createdRoute.ID),
		core_errors.ErrNotFound,
	)
}
