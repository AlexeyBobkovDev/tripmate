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

func TestUsersService_CreateUser_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createUserRequest := testutils.UserCreateTest()

	createdUser := createUser(t, ctx, &createUserRequest, testutils.Password)

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
}

func TestUsersService_CreateUser_Integration_UserAlreadyExists_ConflictError(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createUserRequest := testutils.UserCreateTest()

	_ = createUser(t, ctx, &createUserRequest, testutils.Password)

	createdUserCreatedAgain, err := testDeps.usersService.CreateUser(ctx, &createUserRequest, testutils.Password)

	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.Nil(t, createdUserCreatedAgain)
}
