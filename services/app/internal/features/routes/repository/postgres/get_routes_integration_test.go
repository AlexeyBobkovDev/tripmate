//go:build integration

package routes_postgres_repository

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func TestRoutesRepository_GetRoutes_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)

	firstRouteCreate := routeCreateFixture(2)
	firstCreatedRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, firstRouteCreate)
	require.NoError(t, err)
	require.NotNil(t, firstCreatedRoute)

	secondRouteCreate := routeCreateFixture(1)
	secondRouteCreate.Status = domain.RouteStatusEdited
	secondCreatedRoute, err := routesRepository.CreateRoute(ctx, trip.OwnerID, trip, secondRouteCreate)
	require.NoError(t, err)
	require.NotNil(t, secondCreatedRoute)

	otherTrip := newTripForRoutes(t, ctx)
	otherTripRouteCreate := routeCreateFixture(1)
	otherTripRoute, err := routesRepository.CreateRoute(ctx, otherTrip.OwnerID, otherTrip, otherTripRouteCreate)
	require.NoError(t, err)
	require.NotNil(t, otherTripRoute)

	routes, err := routesRepository.GetRoutes(ctx, trip.ID)
	require.NoError(t, err)
	require.Len(t, routes, 2)

	sort.Slice(routes, func(i, j int) bool {
		return routes[i].ID < routes[j].ID
	})

	expectedRoutes := []struct {
		created   *domain.Route
		generated *domain.RouteCreate
		trip      *domain.Trip
	}{
		{created: firstCreatedRoute, generated: firstRouteCreate, trip: trip},
		{created: secondCreatedRoute, generated: secondRouteCreate, trip: trip},
	}

	sort.Slice(expectedRoutes, func(i, j int) bool {
		return expectedRoutes[i].created.ID < expectedRoutes[j].created.ID
	})

	for ind, expected := range expectedRoutes {
		assert.Equal(t, expected.created.ID, routes[ind].ID)
		assertRouteMatchesCreate(t, expected.generated, routes[ind], expected.trip)
	}
}

func TestRoutesRepository_GetRoutes_Integration_Empty(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)
	trip := newTripForRoutes(t, ctx)

	routes, err := routesRepository.GetRoutes(ctx, trip.ID)
	require.NoError(t, err)
	assert.Empty(t, routes)
}
