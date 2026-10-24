//go:build integration

package users_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersRepository_GetUser_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	createUserRequest := testutils.UserCreateTest()

	actualUser, err := usersRepository.CreateUser(ctx, &createUserRequest)
	require.NoError(t, err)
	require.NotNil(t, actualUser)

	userID := actualUser.ID
	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, userID)
		return err
	})

	retrievedUser, err := usersRepository.GetUser(ctx, actualUser.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedUser)

	assert.NotZero(t, actualUser.ID)
	assert.NotZero(t, actualUser.Version)
	assert.NotZero(t, actualUser.CreatedAt)
	assert.NotZero(t, actualUser.UpdatedAt)

	assert.Nil(t, actualUser.DeletedAt)

	assert.Equal(t, createUserRequest.Name, actualUser.Name)
	assert.Equal(t, createUserRequest.Surname, actualUser.Surname)
	assert.Equal(t, createUserRequest.Username, actualUser.Username)
	assert.Equal(t, createUserRequest.BirthDate, actualUser.BirthDate)
	assert.Equal(t, createUserRequest.Description, actualUser.Description)
	assert.Equal(t, createUserRequest.Email, actualUser.Email)
	assert.Equal(t, createUserRequest.PhoneNumber, actualUser.PhoneNumber)

	assert.Equal(t, actualUser.ID, retrievedUser.ID)
	assert.Equal(t, actualUser.Version, retrievedUser.Version)
	assert.Equal(t, actualUser.Name, retrievedUser.Name)
	assert.Equal(t, actualUser.Surname, retrievedUser.Surname)
	assert.Equal(t, actualUser.Username, retrievedUser.Username)
	assert.Equal(t, actualUser.BirthDate, retrievedUser.BirthDate)
	assert.Equal(t, actualUser.Description, retrievedUser.Description)
	assert.Equal(t, actualUser.Email, retrievedUser.Email)
	assert.Equal(t, actualUser.PhoneNumber, retrievedUser.PhoneNumber)
	assert.Equal(t, actualUser.CreatedAt, retrievedUser.CreatedAt)
	assert.Equal(t, actualUser.UpdatedAt, retrievedUser.UpdatedAt)
	assert.Equal(t, actualUser.DeletedAt, retrievedUser.DeletedAt)
}

func TestUsersRepository_GetUser_Integration_UserNotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userID := -1

	actualUser, err := usersRepository.GetUser(ctx, userID)
	require.Nil(t, actualUser)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestUsersRepository_GetUser_Integration_AlreadyDeleted(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	createUserRequest := testutils.UserCreateTest()

	actualUser, err := usersRepository.CreateUser(ctx, &createUserRequest)
	require.NoError(t, err)
	require.NotNil(t, actualUser)

	userID := actualUser.ID
	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, userID)
		return err
	})

	require.NoError(t, usersRepository.DeleteUser(ctx, actualUser.ID))

	alreadyDeletedUser, err := usersRepository.GetUser(ctx, actualUser.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, alreadyDeletedUser)
}
