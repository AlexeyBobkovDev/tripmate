//go:build integration

package trips_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
)

func TestTripsRepository_CreateTrip_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)
	tripCreate := tripCreateFixture(ownerID)

	actualTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreate)
	require.NoError(t, err)

	assertTripMatchesCreate(t, tripCreate, actualTrip)
}

func TestTripsRepository_CreateTrip_Integration_ForeignKeyViolation(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	const nonExistentOwnerID = 999999
	tripCreate := tripCreateFixture(nonExistentOwnerID)

	actualTrip, err := tripsRepository.CreateTrip(ctx, nonExistentOwnerID, tripCreate)
	require.Error(t, err)
	require.ErrorIs(t, err, core_postgres_pool_errors.ErrViolatedForeignKey)
	assert.Nil(t, actualTrip)
}

// TestTripsRepository_CreateTrip_Integration_MealRequirementsPersisted documents
// that meal requirements accepted by TripCreate must survive a round trip.
func TestTripsRepository_CreateTrip_Integration_MealRequirementsPersisted(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)
	tripCreate := tripCreateFixture(ownerID)
	tripCreate.MealRequirements = []domain.MealRequirement{
		mealRequirement(t, domain.MealBreakfast, "08:00", "09:00", time.Hour),
		mealRequirement(t, domain.MealLunch, "13:00", "14:00", 45*time.Minute),
	}
	require.NotEmpty(t, tripCreate.MealRequirements, "fixture must contain meal requirements")

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreate)
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	assert.ElementsMatch(t, tripCreate.MealRequirements, createdTrip.MealRequirements)

	var mealRequirementsCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM app.meal_requirements WHERE trip_id=$1;`,
		createdTrip.ID,
	).Scan(&mealRequirementsCounter))
	assert.Equal(t, len(tripCreate.MealRequirements), mealRequirementsCounter, "meal requirements must be persisted")

	retrievedTrip, err := tripsRepository.GetTrip(ctx, ownerID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)

	assert.ElementsMatch(t, tripCreate.MealRequirements, retrievedTrip.MealRequirements)
}

// TestTripsRepository_CreateTrip_Integration_InvalidMealType_NoPartialTripCreation
// checks that a meal insert failure surfaces as an error and does not leave a
// trip row behind: create must be atomic.
func TestTripsRepository_CreateTrip_Integration_InvalidMealType_NoPartialTripCreation(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)
	tripCreate := tripCreateFixture(ownerID)
	tripCreate.MealRequirements = []domain.MealRequirement{
		mealRequirement(t, domain.MealBreakfast, "08:00", "09:00", time.Hour),
	}

	// The meal type is unknown to the app.meal_type enum and therefore violates
	// a database constraint.
	tripCreate.MealRequirements[0].Type = domain.MealType("brunch")

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreate)

	assert.Error(t, err, "creating a trip with an invalid meal type must fail")
	assert.Nil(t, createdTrip)

	var tripCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM app.trips WHERE owner_id=$1;`,
		ownerID,
	).Scan(&tripCounter))
	assert.Zero(t, tripCounter, "failed create trip must not leave a partially created trip behind")

	var mealRequirementsCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM app.meal_requirements;`,
	).Scan(&mealRequirementsCounter))
	assert.Zero(t, mealRequirementsCounter, "failed create trip must not leave meal requirements behind")
}

func TestTripsRepository_CreateTrip_Integration_InvalidDailySchedule(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	tests := []struct {
		name          string
		dailySchedule domain.DailySchedule
		expectedError error
	}{
		{
			name: "invalid start time",
			dailySchedule: domain.DailySchedule{
				StartTime: "not-a-time",
				EndTime:   "18:00",
			},
			expectedError: core_errors.ErrInvalidTripStartTime,
		},
		{
			name: "invalid end time",
			dailySchedule: domain.DailySchedule{
				StartTime: "09:00",
				EndTime:   "not-a-time",
			},
			expectedError: core_errors.ErrInvalidTripEndTime,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			checkIfDatabaseIsDirty(t)

			ownerID := insertOwner(t, ctx)
			tripCreate := tripCreateFixture(ownerID)
			tripCreate.DailySchedule = test.dailySchedule

			actualTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreate)

			require.ErrorIs(t, err, test.expectedError)
			assert.Nil(t, actualTrip)
		})
	}
}
