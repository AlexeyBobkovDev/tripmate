//go:build integration

package users_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	trips_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/repository/postgres"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersRepository_DeleteUser_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	createUserRequest := uniqueUserCreate(t)

	actualUser, err := usersRepository.CreateUser(ctx, &createUserRequest)
	require.NoError(t, err)
	require.NotNil(t, actualUser)

	userID := actualUser.ID
	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, userID)
		return err
	})

	err = usersRepository.DeleteUser(ctx, actualUser.ID)
	require.NoError(t, err)

	var deletedAt time.Time
	var version int
	err = usersRepository.pool.QueryRow(
		ctx,
		`SELECT deleted_at, version FROM app.users WHERE id=$1;`,
		actualUser.ID,
	).Scan(&deletedAt, &version)
	require.NoError(t, err)
	assert.False(t, deletedAt.IsZero())
	assert.Equal(t, actualUser.Version+1, version)

	retrievedUser, err := usersRepository.GetUser(ctx, actualUser.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, retrievedUser)

	require.ErrorIs(
		t,
		usersRepository.DeleteUser(ctx, actualUser.ID),
		core_errors.ErrNotFound,
	)
}

func TestUsersRepository_DeleteUser_Integration_UserNotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	err := usersRepository.DeleteUser(ctx, -1)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestUsersRepository_DeleteUser_Integration_SoftDeletesUserTrips(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)
	tripsRepository := trips_postgres_repository.NewTripsRepository(poolForIntegrationTesting)

	deletedUserCreate := uniqueUserCreate(t)
	deletedUser, err := usersRepository.CreateUser(ctx, &deletedUserCreate)
	require.NoError(t, err)
	require.NotNil(t, deletedUser)

	otherUserCreate := uniqueUserCreate(t)
	otherUser, err := usersRepository.CreateUser(ctx, &otherUserCreate)
	require.NoError(t, err)
	require.NotNil(t, otherUser)

	deletedUserTrip, err := tripsRepository.CreateTrip(
		ctx,
		deletedUser.ID,
		testutils.TripCreateTest(deletedUser.ID, "Trip of the deleted user"),
	)
	require.NoError(t, err)
	require.NotNil(t, deletedUserTrip)

	otherUserTrip, err := tripsRepository.CreateTrip(
		ctx,
		otherUser.ID,
		testutils.TripCreateTest(otherUser.ID, "Trip of another user"),
	)
	require.NoError(t, err)
	require.NotNil(t, otherUserTrip)

	registerCleanup(t, func(ctx context.Context) error {
		if _, err := usersRepository.pool.Exec(
			ctx,
			`DELETE FROM app.meal_requirements WHERE id=$1 OR id=$2;`,
			deletedUserTrip.ID,
			otherUserTrip.ID,
		); err != nil {
			return err
		}

		_, err := usersRepository.pool.Exec(
			ctx,
			`DELETE FROM app.trips WHERE owner_id=$1 OR owner_id=$2;`,
			deletedUser.ID,
			otherUser.ID,
		)
		if err != nil {
			return err
		}

		_, err = usersRepository.pool.Exec(
			ctx,
			`DELETE FROM app.users WHERE id=$1 OR id=$2;`,
			deletedUser.ID,
			otherUser.ID,
		)
		return err
	})

	err = usersRepository.DeleteUser(ctx, deletedUser.ID)
	require.NoError(t, err)

	var userDeletedAt time.Time
	var userVersion int
	err = usersRepository.pool.QueryRow(
		ctx,
		`SELECT deleted_at, version FROM app.users WHERE id=$1;`,
		deletedUser.ID,
	).Scan(&userDeletedAt, &userVersion)
	require.NoError(t, err)
	assert.False(t, userDeletedAt.IsZero())
	assert.Equal(t, deletedUser.Version+1, userVersion)

	var tripDeletedAt time.Time
	var tripVersion int
	err = usersRepository.pool.QueryRow(
		ctx,
		`SELECT deleted_at, version FROM app.trips WHERE owner_id=$1 AND id=$2;`,
		deletedUser.ID,
		deletedUserTrip.ID,
	).Scan(&tripDeletedAt, &tripVersion)
	require.NoError(t, err)
	assert.False(t, tripDeletedAt.IsZero())
	assert.Equal(t, deletedUserTrip.Version+1, tripVersion)

	retrievedTrip, err := tripsRepository.GetTrip(ctx, deletedUser.ID, deletedUserTrip.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedTrip)

	deletedUserTrips, err := tripsRepository.GetTrips(ctx, deletedUser.ID)
	require.NoError(t, err)
	require.NotNil(t, deletedUserTrips)
	assert.Empty(t, *deletedUserTrips)

	// Another user's data must not be affected.
	otherUserTrips, err := tripsRepository.GetTrips(ctx, otherUser.ID)
	require.NoError(t, err)
	require.Len(t, *otherUserTrips, 1)
	assert.Equal(t, otherUserTrip.ID, (*otherUserTrips)[0].ID)
	assert.Equal(t, otherUserTrip.Version, (*otherUserTrips)[0].Version)
	assert.Nil(t, (*otherUserTrips)[0].DeletedAt)

	otherUserAfterDelete, err := usersRepository.GetUser(ctx, otherUser.ID)
	require.NoError(t, err)
	require.NotNil(t, otherUserAfterDelete)
	assert.Equal(t, otherUser.Version, otherUserAfterDelete.Version)
}
