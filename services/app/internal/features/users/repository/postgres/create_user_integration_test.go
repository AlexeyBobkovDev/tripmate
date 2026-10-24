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

func TestUsersRepository_CreateUser_Integration_Success(t *testing.T) {
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
}

func TestUsersRepository_CreateUser_Integration_UserAlreadyExists_ConflictError(t *testing.T) {
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

	actualUserCreatedAgain, err := usersRepository.CreateUser(ctx, &createUserRequest)

	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.Nil(t, actualUserCreatedAgain)
}
