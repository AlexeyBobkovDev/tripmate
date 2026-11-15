//go:build integration

package users_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestUsersRepository_PatchUser_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userToCreate := uniqueUserCreate(t)
	createdUser, err := usersRepository.CreateUser(ctx, &userToCreate)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, createdUser.ID)
		return err
	})

	userPatch := userPatchFixture(createdUser.Version)
	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	originalVersion := createdUser.Version

	patchedUser, err := usersRepository.PatchUser(ctx, createdUser.ID, createdUser)
	require.NoError(t, err)
	require.NotNil(t, patchedUser)
	assert.Equal(t, originalVersion+1, patchedUser.Version)

	assert.Equal(t, createdUser.Name, patchedUser.Name)
	assert.Equal(t, createdUser.Surname, patchedUser.Surname)
	assert.Equal(t, createdUser.Username, patchedUser.Username)
	assert.Equal(t, createdUser.BirthDate, patchedUser.BirthDate)
	assert.Equal(t, createdUser.Description, patchedUser.Description)
	assert.Equal(t, createdUser.Email, patchedUser.Email)
	assert.Equal(t, createdUser.PhoneNumber, patchedUser.PhoneNumber)

	gotUser, err := usersRepository.GetUser(ctx, createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, gotUser)

	assert.Equal(t, patchedUser, gotUser)
	assert.Equal(t, originalVersion+1, gotUser.Version)
}

func TestUsersRepository_PatchUser_Integration_Conflict_MismatchingUserVersions(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userToCreate := uniqueUserCreate(t)
	createdUser, err := usersRepository.CreateUser(ctx, &userToCreate)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, createdUser.ID)
		return err
	})

	originalUser := *createdUser

	userPatch := userPatchFixture(createdUser.Version)
	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	// An intentional change of user's version
	createdUser.Version++

	patchedUser, err := usersRepository.PatchUser(ctx, createdUser.ID, createdUser)
	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.Nil(t, patchedUser)

	gotUser, err := usersRepository.GetUser(ctx, createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, gotUser)

	assert.Equal(t, originalUser.Version, gotUser.Version)
	assert.Equal(t, originalUser.Name, gotUser.Name)
	assert.Equal(t, originalUser.Surname, gotUser.Surname)
	assert.Equal(t, originalUser.Username, gotUser.Username)
	assert.Equal(t, originalUser.BirthDate, gotUser.BirthDate)
	assert.Equal(t, originalUser.Description, gotUser.Description)
	assert.Equal(t, originalUser.Email, gotUser.Email)
	assert.Equal(t, originalUser.PhoneNumber, gotUser.PhoneNumber)
}

func TestUsersRepository_PatchUser_Integration_NotFound_InvalidUserID(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userToPatch := domain.User{}

	retrievedUser, err := usersRepository.PatchUser(ctx, -1, &userToPatch)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedUser)
}

func TestUsersRepository_PatchUser_Integration_NotFound_AlreadyDeletedUser(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userToCreate := uniqueUserCreate(t)
	createdUser, err := usersRepository.CreateUser(ctx, &userToCreate)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, createdUser.ID)
		return err
	})

	require.NoError(t, usersRepository.DeleteUser(ctx, createdUser.ID))

	userPatch := userPatchFixture(createdUser.Version)
	require.NoError(t, createdUser.ApplyPatch(&userPatch))
	require.NoError(t, createdUser.Validate())

	retrievedUser, err := usersRepository.PatchUser(ctx, createdUser.ID, createdUser)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedUser)
}
