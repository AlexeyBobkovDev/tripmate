//go:build integration

package routes_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
)

func TestRoutesRepository_CreateRoute_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)
	routeCreate := routeCreateFixture(2)
	routeCreate.Status = domain.RouteStatusEdited

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreate)
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	assertRouteMatchesCreate(t, routeCreate, createdRoute, trip)
	assertRouteRowsPersisted(t, createdRoute)
}

func TestRoutesRepository_CreateRoute_Integration_RollsBackOnFailure(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)

	routeCreate := routeCreateFixture(2)
	// The leg references a stop that does not exist in the generated route,
	// so the stop IDs mapping has no value for it and the insert must fail.
	routeCreate.Days[1].Legs[0].ToStopID = 999999

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreate)
	require.Error(t, err)
	require.Nil(t, createdRoute)

	var routeCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM app.routes WHERE trip_id=$1;`,
		trip.ID,
	).Scan(&routeCounter))
	assert.Zero(t, routeCounter, "route creation must be rolled back atomically")
}

func TestRoutesRepository_CreateRoute_Integration_InvalidTrip_ForeignKeyViolation(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	routeCreate := routeCreateFixture(1)

	createdRoute, err := routesRepository.CreateRoute(
		ctx,
		1,
		&domain.Trip{ID: 999999, Version: 1},
		routeCreate,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, core_postgres_pool_errors.ErrViolatedForeignKey)
	require.Nil(t, createdRoute)
}
