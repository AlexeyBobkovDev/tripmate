//go:build integration

package trips_postgres_repository

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	users_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/repository/postgres"
)

func TestTripsRepository_GetTrips_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)
	otherOwnerID := insertOwner(t, ctx)

	firstTripCreate := tripCreateFixture(ownerID)
	firstTripCreate.Name = "First trip"
	firstCreatedTrip, err := tripsRepository.CreateTrip(ctx, ownerID, firstTripCreate)
	require.NoError(t, err)
	require.NotNil(t, firstCreatedTrip)

	secondTripCreate := tripCreateFixture(ownerID)
	secondTripCreate.Name = "Second trip"
	secondCreatedTrip, err := tripsRepository.CreateTrip(ctx, ownerID, secondTripCreate)
	require.NoError(t, err)
	require.NotNil(t, secondCreatedTrip)

	otherTrip, err := tripsRepository.CreateTrip(ctx, otherOwnerID, tripCreateFixture(otherOwnerID))
	require.NoError(t, err)
	require.NotNil(t, otherTrip)

	trips, err := tripsRepository.GetTrips(ctx, ownerID)
	require.NoError(t, err)
	require.NotNil(t, trips)
	require.Len(t, *trips, 2)

	sort.Slice(*trips, func(i int, j int) bool {
		return (*trips)[i].ID < (*trips)[j].ID
	})

	expectedTrips := map[int]*domain.TripCreate{
		firstCreatedTrip.ID:  firstTripCreate,
		secondCreatedTrip.ID: secondTripCreate,
	}

	for ind := range *trips {
		tripCreate, ok := expectedTrips[(*trips)[ind].ID]
		require.Truef(t, ok, "unexpected trip %d in the result", (*trips)[ind].ID)

		assertTripMatchesCreate(t, tripCreate, &(*trips)[ind])
	}

	for _, trip := range *trips {
		assert.NotEqualf(t, otherTrip.ID, trip.ID, "trips of other owners must not be returned")
	}
}

func TestTripsRepository_GetTrips_Integration_Empty(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)

	trips, err := tripsRepository.GetTrips(ctx, ownerID)
	require.NoError(t, err)
	require.NotNil(t, trips)
	assert.Empty(t, *trips)
}

func TestTripsRepository_GetTrips_Integration_ExcludesSoftDeletedTrips(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tripsRepository := NewTripsRepository(poolForIntegrationTesting)
	usersRepository := users_postgres_repository.NewUsersRepository(poolForIntegrationTesting)

	ownerID := insertOwner(t, ctx)

	createdTrip, err := tripsRepository.CreateTrip(ctx, ownerID, tripCreateFixture(ownerID))
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	// DeleteUser soft-deletes the user's trips instead of removing the rows.
	require.NoError(t, usersRepository.DeleteUser(ctx, ownerID))

	trips, err := tripsRepository.GetTrips(ctx, ownerID)
	require.NoError(t, err)
	require.NotNil(t, trips)
	assert.Empty(t, *trips)
}
