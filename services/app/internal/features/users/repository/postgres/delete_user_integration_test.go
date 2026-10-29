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

	createUserRequest := testutils.UserCreateTest()

	actualUser, err := usersRepository.CreateUser(ctx, &createUserRequest)

	require.NoError(t, err)

	userID := actualUser.ID
	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, userID)
		return err
	})

	err = usersRepository.DeleteUser(ctx, actualUser.ID)
	require.NoError(t, err)

	var deletedAt time.Time
	err = usersRepository.pool.QueryRow(ctx, `SELECT deleted_at FROM app.users WHERE id=$1;`, actualUser.ID).Scan(&deletedAt)
	require.NoError(t, err)
	assert.False(t, deletedAt.IsZero())
}

func TestUsersRepository_DeleteUser_Integration_UserNotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userID := -1

	err := usersRepository.DeleteUser(ctx, userID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestUsersRepository_DeleteUser_Integration_SoftDeletesUserTrips(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)
	tripsRepository := trips_postgres_repository.NewTripsRepository(poolForIntegrationTesting)

	createUserRequest := testutils.UserCreateTest()

	createdUser, err := usersRepository.CreateUser(ctx, &createUserRequest)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	createdTrip, err := tripsRepository.CreateTrip(
		ctx,
		createdUser.ID,
		testutils.TripCreateTest(createdUser.ID, "Delete user soft delete trips integration test"),
	)
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.trips WHERE owner_id=$1 AND id=$2;`, createdUser.ID, createdTrip.ID)
		if err != nil {
			return err
		}

		_, err = usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, createdUser.ID)
		return err
	})

	err = usersRepository.DeleteUser(ctx, createdUser.ID)
	require.NoError(t, err)

	var userDeletedAt time.Time
	err = usersRepository.pool.QueryRow(
		ctx,
		`SELECT deleted_at FROM app.users WHERE id=$1;`,
		createdUser.ID,
	).Scan(&userDeletedAt)
	require.NoError(t, err)
	assert.False(t, userDeletedAt.IsZero())

	var tripDeletedAt time.Time
	err = usersRepository.pool.QueryRow(
		ctx,
		`SELECT deleted_at FROM app.trips WHERE owner_id=$1 AND id=$2;`,
		createdUser.ID,
		createdTrip.ID,
	).Scan(&tripDeletedAt)
	require.NoError(t, err)
	assert.False(t, tripDeletedAt.IsZero())

	retrievedTrip, err := tripsRepository.GetTrip(ctx, createdUser.ID, createdTrip.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedTrip)
}
