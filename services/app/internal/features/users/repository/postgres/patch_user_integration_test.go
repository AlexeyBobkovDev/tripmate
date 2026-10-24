//go:build integration

package users_postgres_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersRepository_PatchUser_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userPatch := testutils.UserPatchTest()
	// Prepare the data (create user)

	userToCreate := testutils.UserCreateTest()
	createdUser, err := usersRepository.CreateUser(context.Background(), &userToCreate)
	require.NoError(t, err)
	require.NotZero(t, createdUser)

	// Register cleanup
	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, createdUser.ID)
		return err
	})

	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	originalVersion := createdUser.Version

	patchedUser, err := usersRepository.PatchUser(context.Background(), createdUser.ID, createdUser)
	require.NoError(t, err)
	require.NotNil(t, patchedUser)
	assert.Equal(t, originalVersion+1, patchedUser.Version)

	gotUser, err := usersRepository.GetUser(context.Background(), createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, gotUser)

	assert.Equal(t, patchedUser, gotUser)
	assert.Equal(t, originalVersion+1, gotUser.Version)
}

func TestUsersRepository_PatchUser_Integration_Conflict_MismatchingUserVersions(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userPatch := testutils.UserPatchTest()
	// Prepare the data (create user)

	userToCreate := testutils.UserCreateTest()
	createdUser, err := usersRepository.CreateUser(context.Background(), &userToCreate)
	require.NoError(t, err)
	require.NotZero(t, createdUser)

	originalVersion := createdUser.Version

	// Register cleanup
	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, createdUser.ID)
		return err
	})

	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	// An intentional change of user's version
	createdUser.Version++

	retrievedUser, err := usersRepository.PatchUser(context.Background(), createdUser.ID, createdUser)
	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.Nil(t, retrievedUser)

	gotUser, err := usersRepository.GetUser(context.Background(), createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, gotUser)

	assert.Equal(t, originalVersion, gotUser.Version)
}

func TestUsersRepository_PatchUser_Integration_NotFound_InvalidUserID(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userToPatch := domain.User{}

	retrievedUser, err := usersRepository.PatchUser(context.Background(), -1, &userToPatch)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedUser)
}

func TestUsersRepository_PatchUser_Integration_NotFound_AlreadyDeletedUser(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userPatch := testutils.UserPatchTest()
	// Prepare the data (create user)

	userToCreate := testutils.UserCreateTest()
	createdUser, err := usersRepository.CreateUser(context.Background(), &userToCreate)
	require.NoError(t, err)
	require.NotZero(t, createdUser)

	// Register cleanup
	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, createdUser.ID)
		return err
	})

	require.NoError(t, usersRepository.DeleteUser(context.Background(), createdUser.ID))

	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	retrievedUser, err := usersRepository.PatchUser(context.Background(), createdUser.ID, createdUser)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedUser)
}
