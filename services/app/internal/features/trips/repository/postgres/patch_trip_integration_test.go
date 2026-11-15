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
)

func TestTripsRepository_PatchTrip_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreateFixture(ownerID))
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	updatedTrip := *createdTrip
	updatedTrip.Name = "Updated Trip Name"
	updatedTrip.PlanningArea = "Saint Petersburg"
	updatedTrip.StartCoordinates = domain.Coordinates{Latitude: 59.9343, Longitude: 30.3351}
	updatedTrip.DailySchedule = domain.DailySchedule{StartTime: "10:00", EndTime: "20:00"}
	updatedTrip.Budget = domain.Money{AmountMinorUnits: 250000, Currency: "RUB"}
	updatedTrip.Interests = []domain.Interest{"food", "architecture"}

	originalVersion := createdTrip.Version

	patchedTrip, err := tripsRepository.PatchTrip(ctx, ownerID, createdTrip.ID, &updatedTrip)
	require.NoError(t, err)
	require.NotNil(t, patchedTrip)

	assert.Equal(t, createdTrip.ID, patchedTrip.ID)
	assert.Equal(t, ownerID, patchedTrip.OwnerID)
	assert.Equal(t, originalVersion+1, patchedTrip.Version)
	assert.Equal(t, updatedTrip.Name, patchedTrip.Name)
	assert.Equal(t, updatedTrip.PlanningArea, patchedTrip.PlanningArea)
	assert.Equal(t, updatedTrip.StartCoordinates, patchedTrip.StartCoordinates)
	assert.Equal(t, updatedTrip.DailySchedule, patchedTrip.DailySchedule)
	assert.Equal(t, updatedTrip.Budget, patchedTrip.Budget)
	assert.Equal(t, updatedTrip.Interests, patchedTrip.Interests)
	assert.Equal(t, updatedTrip.TimeZone, patchedTrip.TimeZone)
	assertSameDate(t, createdTrip.StartDate, patchedTrip.StartDate)
	assertSameDate(t, createdTrip.EndDate, patchedTrip.EndDate)

	retrievedTrip, err := tripsRepository.GetTrip(ctx, ownerID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)

	assertTripsEqual(t, patchedTrip, retrievedTrip)
}

func TestTripsRepository_PatchTrip_Integration_Conflict(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreateFixture(ownerID))
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	conflictingTrip := *createdTrip
	conflictingTrip.Name = "Name that must not be applied"
	conflictingTrip.Version = 999

	patchedTrip, err := tripsRepository.PatchTrip(ctx, ownerID, createdTrip.ID, &conflictingTrip)
	require.ErrorIs(t, err, core_errors.ErrConflict)
	assert.Nil(t, patchedTrip)

	retrievedTrip, err := tripsRepository.GetTrip(ctx, ownerID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)

	assert.Equal(t, createdTrip.Version, retrievedTrip.Version)
	assert.Equal(t, createdTrip.Name, retrievedTrip.Name)
	assert.Equal(t, createdTrip.PlanningArea, retrievedTrip.PlanningArea)
	assert.Equal(t, createdTrip.DailySchedule, retrievedTrip.DailySchedule)
}

func TestTripsRepository_PatchTrip_Integration_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	trip := &domain.Trip{
		ID:      999999,
		OwnerID: 1,
		Version: 1,
		DailySchedule: domain.DailySchedule{
			StartTime: "09:00",
			EndTime:   "18:00",
		},
	}

	patchedTrip, err := tripsRepository.PatchTrip(ctx, 1, 999999, trip)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, patchedTrip)
}

func TestTripsRepository_PatchTrip_Integration_WrongOwner_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)
	otherOwnerID := insertOwner(t, ctx)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreateFixture(ownerID))
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	updatedTrip := *createdTrip
	updatedTrip.Name = "Name that must not be applied"

	patchedTrip, err := tripsRepository.PatchTrip(ctx, otherOwnerID, createdTrip.ID, &updatedTrip)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, patchedTrip)

	retrievedTrip, err := tripsRepository.GetTrip(ctx, ownerID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)
	assert.Equal(t, createdTrip.Name, retrievedTrip.Name)
	assert.Equal(t, createdTrip.Version, retrievedTrip.Version)
}

func TestTripsRepository_PatchTrip_Integration_DeletedTrip_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreateFixture(ownerID))
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	require.NoError(t, tripsRepository.DeleteTrip(ctx, ownerID, createdTrip.ID))

	patchedTrip, err := tripsRepository.PatchTrip(ctx, ownerID, createdTrip.ID, createdTrip)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, patchedTrip)
}

// TestTripsRepository_PatchTrip_Integration_MealRequirementsReplaced checks
// that patching replaces the whole meal requirement set: an existing meal can
// be updated, meals can be added, and meals can be removed. A trip that must
// keep no meals after a patch expects an empty result.
func TestTripsRepository_PatchTrip_Integration_MealRequirementsReplaced(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	tests := []struct {
		name  string
		seed  []domain.MealRequirement
		patch []domain.MealRequirement
	}{
		{
			name:  "updates an existing meal",
			seed:  []domain.MealRequirement{mealRequirement(t, domain.MealBreakfast, "08:00", "09:00", time.Hour)},
			patch: []domain.MealRequirement{mealRequirement(t, domain.MealLunch, "13:00", "14:00", 45*time.Minute)},
		},
		{
			name: "replaces every meal with distinct values",
			seed: []domain.MealRequirement{
				mealRequirement(t, domain.MealBreakfast, "08:00", "09:00", time.Hour),
				mealRequirement(t, domain.MealLunch, "12:00", "13:00", time.Hour),
			},
			patch: []domain.MealRequirement{
				mealRequirement(t, domain.MealDinner, "18:00", "19:00", time.Hour),
				mealRequirement(t, domain.MealLunch, "13:00", "14:30", 45*time.Minute),
			},
		},
		{
			name:  "adds meals when the trip has none",
			seed:  nil,
			patch: []domain.MealRequirement{mealRequirement(t, domain.MealBreakfast, "08:00", "09:00", time.Hour)},
		},
		{
			name:  "removes all meals",
			seed:  []domain.MealRequirement{mealRequirement(t, domain.MealBreakfast, "08:00", "09:00", time.Hour)},
			patch: []domain.MealRequirement{},
		},
		{
			name: "removes a meal",
			seed: []domain.MealRequirement{
				mealRequirement(t, domain.MealBreakfast, "08:00", "09:00", time.Hour),
				mealRequirement(t, domain.MealLunch, "12:00", "13:00", time.Hour),
			},
			patch: []domain.MealRequirement{mealRequirement(t, domain.MealLunch, "12:00", "12:30", 30*time.Minute)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			checkIfDatabaseIsDirty(t)

			ownerID := insertOwner(t, ctx)

			tripCreate := tripCreateFixture(ownerID)
			tripCreate.MealRequirements = nil
			createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreate)
			require.NoError(t, err)
			require.NotNil(t, createdTrip)

			for _, meal := range test.seed {
				insertMealRequirement(t, ctx, createdTrip.ID, meal)
			}

			updatedTrip := *createdTrip
			updatedTrip.MealRequirements = test.patch

			patchedTrip, err := tripsRepository.PatchTrip(ctx, ownerID, createdTrip.ID, &updatedTrip)
			require.NoError(t, err)
			require.NotNil(t, patchedTrip)

			assert.Equal(t, createdTrip.Version+1, patchedTrip.Version)

			assert.ElementsMatch(t, test.patch, patchedTrip.MealRequirements)

			assertMealRequirementsPersisted(t, ctx, createdTrip.ID, test.patch)
		})
	}
}
