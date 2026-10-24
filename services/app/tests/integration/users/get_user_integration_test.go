//go:build integration

package users_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersService_GetUser_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createUserRequest := testutils.UserCreateTest()

	createdUser := createUser(t, ctx, &createUserRequest, testutils.Password)

	retrievedUser, err := testDeps.usersService.GetUser(ctx, createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedUser)

	assert.NotZero(t, createdUser.ID)
	assert.NotZero(t, createdUser.Version)
	assert.NotZero(t, createdUser.CreatedAt)
	assert.NotZero(t, createdUser.UpdatedAt)

	assert.Nil(t, createdUser.DeletedAt)

	assert.Equal(t, createUserRequest.Name, createdUser.Name)
	assert.Equal(t, createUserRequest.Surname, createdUser.Surname)
	assert.Equal(t, createUserRequest.Username, createdUser.Username)
	assert.Equal(t, createUserRequest.BirthDate, createdUser.BirthDate)
	assert.Equal(t, createUserRequest.Description, createdUser.Description)
	assert.Equal(t, createUserRequest.Email, createdUser.Email)
	assert.Equal(t, createUserRequest.PhoneNumber, createdUser.PhoneNumber)

	assert.Equal(t, createdUser.ID, retrievedUser.ID)
	assert.Equal(t, createdUser.Version, retrievedUser.Version)
	assert.Equal(t, createdUser.Name, retrievedUser.Name)
	assert.Equal(t, createdUser.Surname, retrievedUser.Surname)
	assert.Equal(t, createdUser.Username, retrievedUser.Username)
	assert.Equal(t, createdUser.BirthDate, retrievedUser.BirthDate)
	assert.Equal(t, createdUser.Description, retrievedUser.Description)
	assert.Equal(t, createdUser.Email, retrievedUser.Email)
	assert.Equal(t, createdUser.PhoneNumber, retrievedUser.PhoneNumber)
	assert.Equal(t, createdUser.CreatedAt, retrievedUser.CreatedAt)
	assert.Equal(t, createdUser.UpdatedAt, retrievedUser.UpdatedAt)
	assert.Equal(t, createdUser.DeletedAt, retrievedUser.DeletedAt)
}

func TestUsersService_GetUser_Integration_UserNotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userID := -1

	createdUser, err := testDeps.usersService.GetUser(ctx, userID)
	require.Nil(t, createdUser)
	require.ErrorIs(t, err, core_errors.ErrBadRequest)
}

func TestUsersService_GetUser_Integration_AlreadyDeleted(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createUserRequest := testutils.UserCreateTest()

	createdUser := createUser(t, ctx, &createUserRequest, testutils.Password)

	require.NoError(t, testDeps.usersRepository.DeleteUser(ctx, createdUser.ID))

	alreadyDeletedUser, err := testDeps.usersRepository.GetUser(ctx, createdUser.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, alreadyDeletedUser)
}
