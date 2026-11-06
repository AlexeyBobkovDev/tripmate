//go:build integration

package routes_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
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
}

func TestRoutesRepository_DeleteRoute_Integration_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)

	err := routesRepository.DeleteRoute(ctx, 1, 999999)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestRoutesRepository_DeleteRoute_Integration_WrongTrip(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreateFixture(1))
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	err = routesRepository.DeleteRoute(ctx, trip.ID+1, createdRoute.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)

	retrievedRoute, err := routesRepository.GetRoute(ctx, trip.ID, createdRoute.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedRoute)
	assert.Equal(t, createdRoute.ID, retrievedRoute.ID)
}
