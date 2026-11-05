//go:build integration

package routes_integration_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

var testUserSequence atomic.Int64

type fakeRouteGeneratorClient struct {
	generateRouteFunc func(
		ctx context.Context,
		in *route_generator_api_v1.GenerateRouteRequest,
		opts ...grpc.CallOption,
	) (*route_generator_api_v1.GenerateRouteResponse, error)
}

func (c *fakeRouteGeneratorClient) GenerateRoute(
	ctx context.Context,
	in *route_generator_api_v1.GenerateRouteRequest,
	opts ...grpc.CallOption,
) (*route_generator_api_v1.GenerateRouteResponse, error) {
	if c.generateRouteFunc != nil {
		return c.generateRouteFunc(ctx, in, opts...)
	}

	return testutils.RouteGeneratorResponseTest(in.StartDate.AsTime()), nil
}

func setGenerateRouteFunc(
	t *testing.T,
	generateRouteFunc func(
		ctx context.Context,
		in *route_generator_api_v1.GenerateRouteRequest,
		opts ...grpc.CallOption,
	) (*route_generator_api_v1.GenerateRouteResponse, error),
) {
	t.Helper()

	testDeps.routeGeneratorClient.generateRouteFunc = generateRouteFunc

	t.Cleanup(func() {
		testDeps.routeGeneratorClient.generateRouteFunc = nil
	})
}

func registerCleanup(t *testing.T, cleanup func(context.Context) error) {
	t.Helper()

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()

		if err := cleanup(cleanupCtx); err != nil {
			t.Errorf("integration test cleanup failed: %v", err)
		}
	})
}

func createTestUser(t *testing.T, ctx context.Context) *domain.User {
	t.Helper()

	request := testutils.UserCreateTest()
	sequence := testUserSequence.Add(1)
	suffix := fmt.Sprintf("%d_%d", time.Now().UnixNano(), sequence)
	request.Username = "routes_user_" + suffix
	request.Email = "routes_" + suffix + "@example.com"
	request.PhoneNumber = fmt.Sprintf("+7999%07d", sequence%10_000_000)

	user, err := testDeps.usersService.CreateUser(
		ctx,
		&request,
		testutils.TestPassword,
	)
	require.NoError(t, err)
	require.NotNil(t, user)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(
			ctx,
			`DELETE FROM app.users WHERE id = $1;`,
			user.ID,
		)
		return err
	})

	return user
}

func createTestTrip(
	t *testing.T,
	ctx context.Context,
	ownerID int,
	name string,
) *domain.Trip {
	t.Helper()

	trip, err := testDeps.tripsService.CreateTrip(
		ctx,
		ownerID,
		testutils.TripCreateTest(ownerID, name),
	)
	require.NoError(t, err)
	require.NotNil(t, trip)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(
			ctx,
			`DELETE FROM app.trips WHERE owner_id = $1 AND id = $2;`,
			ownerID,
			trip.ID,
		)
		return err
	})

	return trip
}

func createTestRoute(
	t *testing.T,
	ctx context.Context,
	trip *domain.Trip,
) *domain.Route {
	t.Helper()

	route, err := testDeps.routesService.CreateRoute(ctx, trip.OwnerID, trip.ID)
	require.NoError(t, err)
	require.NotNil(t, route)

	return route
}

func assertRouteMatchesGenerated(
	t *testing.T,
	generated *route_generator_api_v1.GenerateRouteResponse,
	actual *domain.Route,
	trip *domain.Trip,
) {
	t.Helper()

	require.NotNil(t, actual)

	assert.NotZero(t, actual.ID)
	assert.Equal(t, trip.ID, actual.TripID)
	assert.Equal(t, trip.Version, actual.TripVersion)
	assert.Equal(t, 1, actual.Version)
	assert.Equal(t, domain.RouteStatus(domain.RouteStatusGenerated), actual.Status)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())

	require.Len(t, actual.Days, len(generated.Days))

	for _, generatedDay := range generated.Days {
		actualDay, found := routeDayByNumber(actual.Days, int(generatedDay.DayNumber))
		require.Truef(t, found, "route day %d is not found", generatedDay.DayNumber)

		assert.NotZero(t, actualDay.ID)
		assertSameDate(t, generatedDay.Date.AsTime(), actualDay.Date)
		assert.Equal(t, int(generatedDay.DayNumber), actualDay.DayNumber)
		require.Len(t, actualDay.Stops, len(generatedDay.Stops))

		stopIDsMapping := make(map[int64]int, len(generatedDay.Stops))

		for _, generatedStop := range generatedDay.Stops {
			actualStop, stopFound := routeStopByPosition(actualDay, int(generatedStop.Position))
			require.Truef(
				t,
				stopFound,
				"route stop with position %d is not found in day %d",
				generatedStop.Position,
				generatedDay.DayNumber,
			)

			assert.NotZero(t, actualStop.ID)
			assert.Equal(t, int(generatedStop.PlaceId), actualStop.PlaceID)
			assert.Equal(t, int(generatedStop.Position), actualStop.Position)
			assert.Equal(t, generatedStop.Name, actualStop.Name)
			assert.Equal(t, generatedStop.Description, actualStop.Description)
			assert.Equal(t, generatedStop.Coordinates.Latitude, actualStop.Coordinates.Latitude)
			assert.Equal(t, generatedStop.Coordinates.Longitude, actualStop.Coordinates.Longitude)
			assertSameInstant(t, generatedStop.ArrivalTime.AsTime(), actualStop.ArrivalTime)
			assertSameInstant(t, generatedStop.DepartureTime.AsTime(), actualStop.DepartureTime)
			assert.Equal(t, generatedStop.VisitDuration.AsDuration(), actualStop.VisitDuration)

			stopIDsMapping[generatedStop.Id] = actualStop.ID
		}

		require.Len(t, actualDay.Legs, len(generatedDay.Legs))

		for _, generatedLeg := range generatedDay.Legs {
			expectedFromStopID := stopIDsMapping[int64(generatedLeg.FromStopId)]
			expectedToStopID := stopIDsMapping[int64(generatedLeg.ToStopId)]

			actualLeg, legFound := routeLegByStopIDs(actualDay, expectedFromStopID, expectedToStopID)
			require.Truef(
				t,
				legFound,
				"route leg %d -> %d is not found in day %d",
				expectedFromStopID,
				expectedToStopID,
				generatedDay.DayNumber,
			)

			assert.NotZero(t, actualLeg.ID)
			assert.Equal(t, expectedFromStopID, actualLeg.FromStopID)
			assert.Equal(t, expectedToStopID, actualLeg.ToStopID)
			assert.Equal(t, domain.TravelModeWalking, actualLeg.TravelMode)
			assert.Equal(t, int(generatedLeg.Distance), actualLeg.Distance)
			assert.Equal(t, generatedLeg.Duration.AsDuration(), actualLeg.Duration)
		}
	}
}

func routeDayByNumber(days []domain.RouteDay, dayNumber int) (domain.RouteDay, bool) {
	for _, day := range days {
		if day.DayNumber == dayNumber {
			return day, true
		}
	}

	return domain.RouteDay{}, false
}

func routeStopByPosition(day domain.RouteDay, position int) (domain.RouteStop, bool) {
	for _, stop := range day.Stops {
		if stop.Position == position {
			return stop, true
		}
	}

	return domain.RouteStop{}, false
}

func routeLegByStopIDs(
	day domain.RouteDay,
	fromStopID int,
	toStopID int,
) (domain.RouteLeg, bool) {
	for _, leg := range day.Legs {
		if leg.FromStopID == fromStopID && leg.ToStopID == toStopID {
			return leg, true
		}
	}

	return domain.RouteLeg{}, false
}

func assertSameDate(t *testing.T, expected time.Time, actual time.Time) {
	t.Helper()

	assert.Equal(t, expected.Format("2006-01-02"), actual.Format("2006-01-02"))
}

func assertSameInstant(t *testing.T, expected time.Time, actual time.Time) {
	t.Helper()

	assert.Truef(
		t,
		expected.Equal(actual),
		"expected time %s, got %s",
		expected.Format(time.RFC3339Nano),
		actual.Format(time.RFC3339Nano),
	)
}
