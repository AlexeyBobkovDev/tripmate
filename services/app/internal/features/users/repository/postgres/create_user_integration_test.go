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

func TestUsersRepository_CreateUser_Integration_Success(t *testing.T) {
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

	assert.NotZero(t, actualUser.ID)
	assert.Equal(t, 1, actualUser.Version)
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

	var storedPasswordHash string
	require.NoError(t, usersRepository.pool.QueryRow(
		ctx,
		`SELECT password_hash FROM app.users WHERE id=$1;`,
		actualUser.ID,
	).Scan(&storedPasswordHash))
	assert.Equal(t, createUserRequest.PasswordHash, storedPasswordHash)
}

func TestUsersRepository_CreateUser_Integration_DuplicateUniqueField_ConflictError(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	tests := []struct {
		name      string
		duplicate func(duplicate domain.UserCreate, existing domain.UserCreate) domain.UserCreate
	}{
		{
			name: "duplicate username",
			duplicate: func(duplicate domain.UserCreate, existing domain.UserCreate) domain.UserCreate {
				duplicate.Username = existing.Username
				return duplicate
			},
		},
		{
			name: "duplicate email",
			duplicate: func(duplicate domain.UserCreate, existing domain.UserCreate) domain.UserCreate {
				duplicate.Email = existing.Email
				return duplicate
			},
		},
		{
			name: "duplicate phone number",
			duplicate: func(duplicate domain.UserCreate, existing domain.UserCreate) domain.UserCreate {
				duplicate.PhoneNumber = existing.PhoneNumber
				return duplicate
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			checkIfDatabaseIsDirty(t)

			existingUserCreate := uniqueUserCreate(t)
			existingUser, err := usersRepository.CreateUser(ctx, &existingUserCreate)
			require.NoError(t, err)
			require.NotNil(t, existingUser)

			registerCleanup(t, func(ctx context.Context) error {
				_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, existingUser.ID)
				return err
			})

			duplicateUserCreate := test.duplicate(uniqueUserCreate(t), existingUserCreate)

			duplicateUser, err := usersRepository.CreateUser(ctx, &duplicateUserCreate)
			require.ErrorIs(t, err, core_errors.ErrConflict)
			assert.Nil(t, duplicateUser)
		})
	}
}

// TestUsersRepository_CreateUser_Integration_ReusesUsernameAfterDeletion checks
// the intent of the partial unique indexes: unique values are only reserved by
// users that are not deleted.
func TestUsersRepository_CreateUser_Integration_ReusesUsernameAfterDeletion(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	createUserRequest := uniqueUserCreate(t)

	deletedUser, err := usersRepository.CreateUser(ctx, &createUserRequest)
	require.NoError(t, err)
	require.NotNil(t, deletedUser)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, deletedUser.ID)
		return err
	})

	require.NoError(t, usersRepository.DeleteUser(ctx, deletedUser.ID))

	recreatedUser, err := usersRepository.CreateUser(ctx, &createUserRequest)
	require.NoError(t, err)
	require.NotNil(t, recreatedUser)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := usersRepository.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, recreatedUser.ID)
		return err
	})

	assert.NotEqual(t, deletedUser.ID, recreatedUser.ID)
	assert.Equal(t, createUserRequest.Username, recreatedUser.Username)
	assert.Equal(t, createUserRequest.Email, recreatedUser.Email)
	assert.Equal(t, createUserRequest.PhoneNumber, recreatedUser.PhoneNumber)
}
