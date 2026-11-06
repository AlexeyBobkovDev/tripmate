//go:build integration

package routes_integration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestRoutesService_CreateRoute_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Create route integration test")

	createdRoute, err := testDeps.routesService.CreateRoute(ctx, user.ID, trip.ID)
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	expectedGeneratedRoute := testutils.RouteGeneratorResponseTest(trip.StartDate)
	assertRouteMatchesGenerated(t, expectedGeneratedRoute, createdRoute, trip)

	retrievedRoute, err := testDeps.routesService.GetRoute(ctx, trip.ID, createdRoute.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedRoute)

	assert.Equal(t, createdRoute.ID, retrievedRoute.ID)
	assertRouteMatchesGenerated(t, expectedGeneratedRoute, retrievedRoute, trip)
}

func TestRoutesService_CreateRoute_Integration_InvalidOwnerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createdRoute, err := testDeps.routesService.CreateRoute(ctx, -1, 1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, createdRoute)
}

func TestRoutesService_CreateRoute_Integration_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	createdRoute, err := testDeps.routesService.CreateRoute(ctx, user.ID, -1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, createdRoute)
}

func TestRoutesService_CreateRoute_Integration_TripNotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	createdRoute, err := testDeps.routesService.CreateRoute(ctx, user.ID, 2_000_000_000)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, createdRoute)
}

func TestRoutesService_CreateRoute_Integration_WrongOwner(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	owner := createTestUser(t, ctx)
	anotherUser := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, owner.ID, "Route wrong owner integration test")

	createdRoute, err := testDeps.routesService.CreateRoute(ctx, anotherUser.ID, trip.ID)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, createdRoute)
}

func TestRoutesService_CreateRoute_Integration_GeneratorFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Route generator failure integration test")

	generatorErr := errors.New("route generator is unavailable")
	setGenerateRouteFunc(t, func(
		ctx context.Context,
		in *route_generator_api_v1.GenerateRouteRequest,
		opts ...grpc.CallOption,
	) (*route_generator_api_v1.GenerateRouteResponse, error) {
		return nil, generatorErr
	})

	createdRoute, err := testDeps.routesService.CreateRoute(ctx, user.ID, trip.ID)

	require.ErrorIs(t, err, generatorErr)
	require.Nil(t, createdRoute)
}

func TestRoutesService_CreateRoute_Integration_GeneratorReceivesTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID, "Route generator request integration test")

	var actualRequest *route_generator_api_v1.GenerateRouteRequest
	setGenerateRouteFunc(t, func(
		ctx context.Context,
		in *route_generator_api_v1.GenerateRouteRequest,
		opts ...grpc.CallOption,
	) (*route_generator_api_v1.GenerateRouteResponse, error) {
		actualRequest = in

		return testutils.RouteGeneratorResponseTest(in.StartDate.AsTime()), nil
	})

	createdRoute, err := testDeps.routesService.CreateRoute(ctx, user.ID, trip.ID)
	require.NoError(t, err)
	require.NotNil(t, createdRoute)

	require.NotNil(t, actualRequest)
	assert.Equal(t, trip.Name, actualRequest.Name)
	assert.Equal(t, trip.PlanningArea, actualRequest.PlanningArea)
	assert.Equal(t, trip.StartCoordinates.Latitude, actualRequest.StartCoordinates.Latitude)
	assert.Equal(t, trip.StartCoordinates.Longitude, actualRequest.StartCoordinates.Longitude)
	assert.True(t, trip.StartDate.Equal(actualRequest.StartDate.AsTime()))
	assert.True(t, trip.EndDate.Equal(actualRequest.EndDate.AsTime()))
	assert.Equal(t, trip.Budget.AmountMinorUnits, actualRequest.Budget.AmountMinorUnits)
	assert.Equal(t, trip.Budget.Currency, actualRequest.Budget.Currency)
	assert.Equal(t, len(trip.Interests), len(actualRequest.Interests))
	assert.Equal(t, trip.TimeZone, actualRequest.TimeZone)
}
