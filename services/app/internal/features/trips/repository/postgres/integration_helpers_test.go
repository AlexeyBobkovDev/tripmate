//go:build integration

package trips_postgres_repository

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

const insertOwnerQuery = `
	INSERT INTO app.users (name, surname, username, birth_date, description, email, phone_number, password_hash)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id;
`

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

// checkIfDatabaseIsDirty guards tests against leaked rows from a previous test.
// Integration tests in this package run sequentially, so the database must be
// empty before each test starts. Checking every table also verifies that
// cascading deletes and cleanups worked.
func checkIfDatabaseIsDirty(t *testing.T) {
	t.Helper()

	tables := []string{
		"app.users",
		"app.trips",
		"app.meal_requirements",
		"app.routes",
		"app.route_days",
		"app.route_stops",
		"app.route_legs",
	}

	for _, table := range tables {
		var counter int

		require.NoError(t, poolForIntegrationTesting.QueryRow(
			context.Background(),
			fmt.Sprintf(`SELECT COUNT(*) FROM %s;`, table),
		).Scan(&counter))
		require.Zero(t, counter, "dirty database state: %s is not empty", table)
	}
}

// insertOwner inserts a user directly (trips repository does not own users) and
// returns its ID. Unique fields are generated per call so a single test can
// create more than one owner.
func insertOwner(t *testing.T, ctx context.Context) int {
	t.Helper()

	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)

	userCreate := testutils.UserCreateTest()
	userCreate.Username = "trips-" + suffix
	userCreate.Email = fmt.Sprintf("trips-%s@example.com", suffix)
	userCreate.PhoneNumber = "+7999" + suffix[len(suffix)-7:]

	var ownerID int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		ctx,
		insertOwnerQuery,
		userCreate.Name,
		userCreate.Surname,
		userCreate.Username,
		userCreate.BirthDate,
		userCreate.Description,
		userCreate.Email,
		userCreate.PhoneNumber,
		userCreate.PasswordHash,
	).Scan(&ownerID))

	registerCleanup(t, func(ctx context.Context) error {
		if _, err := poolForIntegrationTesting.Exec(
			ctx,
			`DELETE FROM app.meal_requirements WHERE trip_id IN (
				SELECT id FROM app.trips WHERE owner_id=$1
			);`,
			ownerID,
		); err != nil {
			return err
		}

		if _, err := poolForIntegrationTesting.Exec(
			ctx,
			`DELETE FROM app.trips WHERE owner_id=$1;`,
			ownerID,
		); err != nil {
			return err
		}

		_, err := poolForIntegrationTesting.Exec(
			ctx,
			`DELETE FROM app.users WHERE id=$1;`,
			ownerID,
		)
		return err
	})

	return ownerID
}

// tripCreateFixture builds a valid future single-day trip. DailySchedule is
// required by the repository, and a missing schedule makes CreateTrip fail
// before it ever reaches PostgreSQL.
func tripCreateFixture(ownerID int) *domain.TripCreate {
	return testutils.TripCreateTest(ownerID, "Trip to Moscow")
}

// assertTripMatchesCreate checks the fields of a trip together with its meal
// requirements. Meal requirements are stored in a separate table, so a
// round-trip test only proves persistence when they are compared too.
func assertTripMatchesCreate(
	t *testing.T,
	expected *domain.TripCreate,
	actual *domain.Trip,
) {
	t.Helper()

	require.NotNil(t, actual)

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.OwnerID, actual.OwnerID)
	assert.Equal(t, 1, actual.Version)

	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.PlanningArea, actual.PlanningArea)
	assert.Equal(t, expected.StartCoordinates, actual.StartCoordinates)
	assertSameDate(t, expected.StartDate, actual.StartDate)
	assertSameDate(t, expected.EndDate, actual.EndDate)
	assert.Equal(t, expected.DailySchedule, actual.DailySchedule)
	assert.Equal(t, expected.Budget, actual.Budget)
	assert.Equal(t, expected.Interests, actual.Interests)
	assert.Equal(t, expected.TimeZone, actual.TimeZone)
	assert.ElementsMatch(t, expected.MealRequirements, actual.MealRequirements)

	assert.NotZero(t, actual.CreatedAt)
	assert.NotZero(t, actual.UpdatedAt)
	assert.Nil(t, actual.DeletedAt)
}

// assertTripsEqual compares two persisted trips. Timestamps are compared as
// instants because the driver may return a different (but equivalent) time
// representation.
func assertTripsEqual(t *testing.T, expected *domain.Trip, actual *domain.Trip) {
	t.Helper()

	require.NotNil(t, actual)

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.OwnerID, actual.OwnerID)
	assert.Equal(t, expected.Version, actual.Version)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.PlanningArea, actual.PlanningArea)
	assert.Equal(t, expected.StartCoordinates, actual.StartCoordinates)
	assertSameDate(t, expected.StartDate, actual.StartDate)
	assertSameDate(t, expected.EndDate, actual.EndDate)
	assert.Equal(t, expected.DailySchedule, actual.DailySchedule)
	assert.Equal(t, expected.Budget, actual.Budget)
	assert.Equal(t, expected.Interests, actual.Interests)
	assert.Equal(t, expected.TimeZone, actual.TimeZone)
	assert.ElementsMatch(t, expected.MealRequirements, actual.MealRequirements)
	assertSameInstant(t, expected.CreatedAt, actual.CreatedAt)
	assertSameInstant(t, expected.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, expected.DeletedAt, actual.DeletedAt)
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

// mealRequirement creates a meal requirement with the given type and window.
// Distinct windows/durations are used by tests that replace several meals, so
// a bug that updates every row with the same values is observable.
func mealRequirement(
	t *testing.T,
	mealType domain.MealType,
	start, end string,
	duration time.Duration,
) domain.MealRequirement {
	t.Helper()

	return domain.MealRequirement{
		Type:     mealType,
		Window:   domain.TimeWindow{Start: domain.LocalTime(start), End: domain.LocalTime(end)},
		Duration: duration,
		Budget:   domain.Money{AmountMinorUnits: 50000, Currency: "RUB"},
	}
}

// insertMealRequirement inserts a meal requirement directly for the given trip.
func insertMealRequirement(t *testing.T, ctx context.Context, tripID int, meal domain.MealRequirement) {
	t.Helper()

	startTime, err := time.Parse("15:04", string(meal.Window.Start))
	require.NoError(t, err)
	endTime, err := time.Parse("15:04", string(meal.Window.End))
	require.NoError(t, err)

	_, err = poolForIntegrationTesting.Exec(
		ctx,
		`INSERT INTO app.meal_requirements (
			trip_id, meal_type, start_time, end_time, visit_duration_minutes, amount_minor_units, currency
		) VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		tripID,
		meal.Type,
		startTime,
		endTime,
		int(meal.Duration/time.Minute),
		meal.Budget.AmountMinorUnits,
		meal.Budget.Currency,
	)
	require.NoError(t, err)
}

// fetchMealRequirements returns all meal requirements of a trip as persisted.
func fetchMealRequirements(t *testing.T, ctx context.Context, tripID int) []domain.MealRequirement {
	t.Helper()

	rows, err := poolForIntegrationTesting.Query(
		ctx,
		`SELECT meal_type, start_time, end_time, visit_duration_minutes, amount_minor_units, currency
		 FROM app.meal_requirements
		 WHERE trip_id=$1;`,
		tripID,
	)
	require.NoError(t, err)
	defer rows.Close()

	var mealRequirements []domain.MealRequirement
	for rows.Next() {
		var (
			meal            domain.MealRequirement
			startTime       time.Time
			endTime         time.Time
			durationMinutes int
		)

		require.NoError(t, rows.Scan(
			&meal.Type,
			&startTime,
			&endTime,
			&durationMinutes,
			&meal.Budget.AmountMinorUnits,
			&meal.Budget.Currency,
		))

		meal.Window = domain.TimeWindow{
			Start: domain.LocalTime(startTime.Format("15:04")),
			End:   domain.LocalTime(endTime.Format("15:04")),
		}
		meal.Duration = time.Duration(durationMinutes) * time.Minute

		mealRequirements = append(mealRequirements, meal)
	}
	require.NoError(t, rows.Err())

	return mealRequirements
}

// assertMealRequirementsPersisted checks that the persisted meal requirements
// of a trip are exactly the expected set (order does not matter, a trip may have
// several meals).
func assertMealRequirementsPersisted(t *testing.T, ctx context.Context, tripID int, expected []domain.MealRequirement) {
	t.Helper()

	assert.ElementsMatch(
		t,
		expected,
		fetchMealRequirements(t, ctx, tripID),
		"meal requirements for trip %d",
		tripID,
	)
}
