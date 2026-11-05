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

	createdRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, routeCreateFixture(1))
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	retrievedRoute, err := routesRepository.GetRoute(ctx, trip.ID+1, createdRoute.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedRoute)
}
