//go:build integration

package routes_postgres_repository

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestRoutesRepository_GetRoute_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)
	routeCreate := routeCreateFixture(2)

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreate)
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	retrievedRoute, err := routesRepository.GetRoute(ctx, trip.ID, createdRoute.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedRoute)

	assert.Equal(t, createdRoute.ID, retrievedRoute.ID)
	assertSameInstant(t, createdRoute.CreatedAt, retrievedRoute.CreatedAt)
	assertSameInstant(t, createdRoute.UpdatedAt, retrievedRoute.UpdatedAt)
	assertRouteMatchesCreate(t, routeCreate, retrievedRoute, trip)
	assertRouteRowsPersisted(t, retrievedRoute)
}

func TestRoutesRepository_GetRoute_Integration_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)

	retrievedRoute, err := routesRepository.GetRoute(ctx, 1, 999999)
	require.Error(t, err)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedRoute)
}

func TestRoutesRepository_GetRoute_Integration_WrongTrip(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)
	otherTrip := newTripForRoutes(t, ctx)

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreateFixture(1))
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	retrievedRoute, err := routesRepository.GetRoute(ctx, otherTrip.ID, createdRoute.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedRoute)

	// The route must still be readable through its own trip.
	ownerRoute, err := routesRepository.GetRoute(ctx, trip.ID, createdRoute.ID)
	require.NoError(t, err)
	require.NotNil(t, ownerRoute)
	assert.Equal(t, createdRoute.ID, ownerRoute.ID)
}

// TestRoutesRepository_GetRoute_Integration_Ordering inserts days and stops in
// reverse order and expects the repository to return them in domain order:
// days by day_number and stops by position.
func TestRoutesRepository_GetRoute_Integration_Ordering(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)

	routeCreate := routeCreateFixture(3)
	slices.Reverse(routeCreate.Days)

	firstGeneratedDay := &routeCreate.Days[0]
	firstGeneratedDay.Stops[0], firstGeneratedDay.Stops[1] = firstGeneratedDay.Stops[1], firstGeneratedDay.Stops[0]

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreate)
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	retrievedRoute, err := routesRepository.GetRoute(ctx, trip.ID, createdRoute.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedRoute)

	assertRouteMatchesCreate(t, routeCreate, retrievedRoute, trip)
	assertRouteDaysOrdered(t, retrievedRoute.Days)
}
