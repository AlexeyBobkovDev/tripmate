//go:build integration

package trips_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestTripsRepository_DeleteTrip_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreateFixture(ownerID))
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	require.NotEmpty(t, createdTrip.MealRequirements, "create fixture must contain a meal requirement")

	require.NoError(t, tripsRepository.DeleteTrip(ctx, ownerID, createdTrip.ID))

	var tripCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM app.trips WHERE id=$1;`,
		createdTrip.ID,
	).Scan(&tripCounter))
	assert.Zero(t, tripCounter, "DeleteTrip must remove the trip row")

	var mealRequirementsCounter int
	require.NoError(t, poolForIntegrationTesting.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM app.meal_requirements WHERE trip_id=$1;`,
		createdTrip.ID,
	).Scan(&mealRequirementsCounter))
	assert.Zero(t, mealRequirementsCounter, "meal requirements must be removed when the trip is deleted")

	retrievedTrip, err := tripsRepository.GetTrip(ctx, ownerID, createdTrip.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, retrievedTrip)

	trips, err := tripsRepository.GetTrips(ctx, ownerID)
	require.NoError(t, err)
	require.NotNil(t, trips)
	assert.Empty(t, *trips)

	require.ErrorIs(
		t,
		tripsRepository.DeleteTrip(ctx, ownerID, createdTrip.ID),
		core_errors.ErrNotFound,
	)
}

func TestTripsRepository_DeleteTrip_Integration_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	err := tripsRepository.DeleteTrip(ctx, 1, 999999)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestTripsRepository_DeleteTrip_Integration_WrongOwner_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)
	otherOwnerID := insertOwner(t, ctx)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreateFixture(ownerID))
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	err = tripsRepository.DeleteTrip(ctx, otherOwnerID, createdTrip.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)

	retrievedTrip, err := tripsRepository.GetTrip(ctx, ownerID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)
	assert.Equal(t, createdTrip.ID, retrievedTrip.ID)
}
