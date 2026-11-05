//go:build integration

package routes_postgres_repository

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	trips_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/repository/postgres"
	users_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/repository/postgres"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

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

func checkIfDatabaseIsDirty(t *testing.T) {
	var counter int

	require.NoError(t, poolForIntegrationTesting.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM app.routes;
		`,
	).Scan(&counter))
	if counter != 0 {
		panic("dirty database state!!!")
	}
}

func newTripForRoutes(t *testing.T, ctx context.Context) *domain.Trip {
	t.Helper()

	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)

	userCreate := testutils.UserCreateTest()
	userCreate.Username = "routes-" + suffix
	userCreate.Email = fmt.Sprintf("routes-%s@example.com", suffix)
	userCreate.PhoneNumber = "+7999" + suffix[len(suffix)-7:]

	usersRepository := users_postgres_repository.NewUsersRepository(poolForIntegrationTesting)

	user, err := usersRepository.CreateUser(ctx, &userCreate)
	require.NoError(t, err)
	require.NotNil(t, user)

	tripsRepository := trips_postgres_repository.NewTripsRepository(poolForIntegrationTesting)

	trip, err := tripsRepository.CreateTrip(
		ctx,
		user.ID,
		testutils.TripCreateTest(user.ID, "Routes integration test "+suffix),
	)
	require.NoError(t, err)
	require.NotNil(t, trip)

	registerCleanup(t, func(ctx context.Context) error {
		if _, err := poolForIntegrationTesting.Exec(
			ctx,
			`DELETE FROM app.trips WHERE owner_id=$1;`,
			user.ID,
		); err != nil {
			return err
		}

		_, err := poolForIntegrationTesting.Exec(
			ctx,
			`DELETE FROM app.users WHERE id=$1;`,
			user.ID,
		)
		return err
	})

	return trip
}

func routeCreateFixture(daysAmount int) *domain.RouteCreate {
	days := make([]domain.RouteDay, daysAmount)

	for ind := range days {
		days[ind] = routeDayFixture(
			ind+1,
			futureRouteDayDate(ind+2),
			100*(ind+1),
		)
	}

	return &domain.RouteCreate{
		Status: domain.RouteStatusGenerated,
		Days:   days,
	}
}

func futureRouteDayDate(offsetDays int) time.Time {
	now := time.Now().UTC()

	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	).AddDate(0, 0, offsetDays)
}

func routeDayFixture(
	dayNumber int,
	date time.Time,
	stopIDBase int,
) domain.RouteDay {
	firstStop := domain.RouteStop{
		ID:          stopIDBase,
		PlaceID:     stopIDBase + 1,
		Position:    1,
		Name:        fmt.Sprintf("Stop %d-%d", dayNumber, 1),
		Description: fmt.Sprintf("Description %d-%d", dayNumber, 1),
		Coordinates: domain.Coordinates{
			Latitude:  55.7539 + float64(dayNumber)/1000,
			Longitude: 37.6208 + float64(dayNumber)/1000,
		},
		ArrivalTime:   date.Add(9 * time.Hour),
		DepartureTime: date.Add(11 * time.Hour),
		VisitDuration: 90 * time.Minute,
	}

	secondStop := domain.RouteStop{
		ID:          stopIDBase + 1,
		PlaceID:     stopIDBase + 2,
		Position:    2,
		Name:        fmt.Sprintf("Stop %d-%d", dayNumber, 2),
		Description: fmt.Sprintf("Description %d-%d", dayNumber, 2),
		Coordinates: domain.Coordinates{
			Latitude:  55.7639 + float64(dayNumber)/1000,
			Longitude: 37.6308 + float64(dayNumber)/1000,
		},
		ArrivalTime:   date.Add(12 * time.Hour),
		DepartureTime: date.Add(13*time.Hour + 30*time.Minute),
		VisitDuration: 45 * time.Minute,
	}

	return domain.RouteDay{
		ID:        dayNumber,
		Date:      date,
		DayNumber: dayNumber,
		Stops:     []domain.RouteStop{firstStop, secondStop},
		Legs: []domain.RouteLeg{
			{
				ID:         dayNumber*10 + 1,
				FromStopID: firstStop.ID,
				ToStopID:   secondStop.ID,
				TravelMode: domain.TravelModeWalking,
				Distance:   1500,
				Duration:   15 * time.Minute,
			},
		},
	}
}

func assertRouteMatchesCreate(
	t *testing.T,
	generated *domain.RouteCreate,
	actual *domain.Route,
	trip *domain.Trip,
) {
	t.Helper()

	require.NotNil(t, actual)

	assert.NotZero(t, actual.ID)
	assert.Equal(t, trip.ID, actual.TripID)
	assert.Equal(t, trip.Version, actual.TripVersion)
	assert.Equal(t, 1, actual.Version)
	assert.Equal(t, generated.Status, actual.Status)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())

	require.Len(t, actual.Days, len(generated.Days))

	for _, generatedDay := range generated.Days {
		actualDay, found := routeDayByNumber(actual.Days, generatedDay.DayNumber)
		require.Truef(t, found, "route day %d is not found", generatedDay.DayNumber)

		assert.NotZero(t, actualDay.ID)
		assertSameDate(t, generatedDay.Date, actualDay.Date)
		assert.Equal(t, generatedDay.DayNumber, actualDay.DayNumber)
		require.Len(t, actualDay.Stops, len(generatedDay.Stops))

		stopIDsMapping := make(map[int]int, len(generatedDay.Stops))

		for _, generatedStop := range generatedDay.Stops {
			actualStop, stopFound := routeStopByPosition(actualDay, generatedStop.Position)
			require.Truef(
				t,
				stopFound,
				"route stop with position %d is not found in day %d",
				generatedStop.Position,
				generatedDay.DayNumber,
			)

			assert.NotZero(t, actualStop.ID)
			assert.Equal(t, generatedStop.PlaceID, actualStop.PlaceID)
			assert.Equal(t, generatedStop.Position, actualStop.Position)
			assert.Equal(t, generatedStop.Name, actualStop.Name)
			assert.Equal(t, generatedStop.Description, actualStop.Description)
			assert.Equal(t, generatedStop.Coordinates, actualStop.Coordinates)
			assertSameInstant(t, generatedStop.ArrivalTime, actualStop.ArrivalTime)
			assertSameInstant(t, generatedStop.DepartureTime, actualStop.DepartureTime)
			assert.Equal(t, generatedStop.VisitDuration, actualStop.VisitDuration)

			stopIDsMapping[generatedStop.ID] = actualStop.ID
		}

		require.Len(t, actualDay.Legs, len(generatedDay.Legs))

		for _, generatedLeg := range generatedDay.Legs {
			expectedFromStopID := stopIDsMapping[generatedLeg.FromStopID]
			expectedToStopID := stopIDsMapping[generatedLeg.ToStopID]

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
			assert.Equal(t, generatedLeg.TravelMode, actualLeg.TravelMode)
			assert.Equal(t, generatedLeg.Distance, actualLeg.Distance)
			assert.Equal(t, generatedLeg.Duration, actualLeg.Duration)
		}
	}
}

func assertRouteRowsPersisted(t *testing.T, route *domain.Route) {
	t.Helper()

	var dayCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM app.route_days WHERE route_id=$1;`,
		route.ID,
	).Scan(&dayCounter))
	assert.Equal(t, len(route.Days), dayCounter)

	var stopCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM app.route_stops AS stops
		JOIN app.route_days AS days ON days.id = stops.route_day_id
		WHERE days.route_id=$1;
		`,
		route.ID,
	).Scan(&stopCounter))

	var expectedStops int
	for _, day := range route.Days {
		expectedStops += len(day.Stops)
	}
	assert.Equal(t, expectedStops, stopCounter)

	var legCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM app.route_legs AS legs
		JOIN app.route_days AS days ON days.id = legs.route_day_id
		WHERE days.route_id=$1;
		`,
		route.ID,
	).Scan(&legCounter))

	var expectedLegs int
	for _, day := range route.Days {
		expectedLegs += len(day.Legs)
	}
	assert.Equal(t, expectedLegs, legCounter)

	for _, day := range route.Days {
		for _, stop := range day.Stops {
			var visitDurationSeconds int64
			require.NoError(t, poolForIntegrationTesting.QueryRow(
				context.Background(),
				`SELECT visit_duration FROM app.route_stops WHERE id=$1;`,
				stop.ID,
			).Scan(&visitDurationSeconds))
			assert.Equal(t, int64(stop.VisitDuration/time.Second), visitDurationSeconds)
		}

		for _, leg := range day.Legs {
			var durationSeconds int64
			require.NoError(t, poolForIntegrationTesting.QueryRow(
				context.Background(),
				`SELECT duration FROM app.route_legs WHERE id=$1;`,
				leg.ID,
			).Scan(&durationSeconds))
			assert.Equal(t, int64(leg.Duration/time.Second), durationSeconds)
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
