//go:build integration

package routes_postgres_repository

import (
	"context"
	"slices"
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

	for _, route := range routes {
		assert.NotEqualf(
			t,
			otherTripRoute.ID,
			route.ID,
			"routes of other trips must not be returned",
		)
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

func TestRoutesRepository_GetRoutes_Integration_NonExistentTrip_Empty(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	routesRepository := NewRoutesRepository(poolForIntegrationTesting)

	routes, err := routesRepository.GetRoutes(ctx, 999999)
	require.NoError(t, err)
	assert.Empty(t, routes)
}

// TestRoutesRepository_GetRoutes_Integration_Ordering inserts days and stops in
// reverse order and expects every returned route to be ordered by day_number
// and by stop position.
func TestRoutesRepository_GetRoutes_Integration_Ordering(t *testing.T) {
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

	routes, err := routesRepository.GetRoutes(ctx, trip.ID)
	require.NoError(t, err)
	require.Len(t, routes, 1)

	assertRouteMatchesCreate(t, routeCreate, routes[0], trip)
	assertRouteDaysOrdered(t, routes[0].Days)
}
