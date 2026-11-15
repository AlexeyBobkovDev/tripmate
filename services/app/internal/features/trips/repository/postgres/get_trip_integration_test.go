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

func TestTripsRepository_GetTrip_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)
	tripCreate := tripCreateFixture(ownerID)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreate)
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	retrievedTrip, err := tripsRepository.GetTrip(ctx, ownerID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)

	assertTripMatchesCreate(t, tripCreate, retrievedTrip)
	assertTripsEqual(t, createdTrip, retrievedTrip)
}

func TestTripsRepository_GetTrip_Integration_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	retrievedTrip, err := tripsRepository.GetTrip(ctx, 1, 999999)
	require.Nil(t, retrievedTrip)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestTripsRepository_GetTrip_Integration_WrongOwner_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)
	otherOwnerID := insertOwner(t, ctx)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreateFixture(ownerID))
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	retrievedTrip, err := tripsRepository.GetTrip(ctx, otherOwnerID, createdTrip.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, retrievedTrip)

	// The trip must still be readable by its actual owner.
	ownerTrip, err := tripsRepository.GetTrip(ctx, ownerID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, ownerTrip)
	assert.Equal(t, createdTrip.ID, ownerTrip.ID)
}

func TestTripsRepository_GetTrip_Integration_DeletedTrip_NotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreateFixture(ownerID))
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	require.NoError(t, tripsRepository.DeleteTrip(ctx, ownerID, createdTrip.ID))

	retrievedTrip, err := tripsRepository.GetTrip(ctx, ownerID, createdTrip.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, retrievedTrip)
}
