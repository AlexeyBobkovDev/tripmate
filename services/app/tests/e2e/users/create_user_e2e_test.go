//go:build e2e

package users_e2e_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersTransportHTTP_CreateUser_E2E_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createdUser, err := createUser(t, ctx, 5*time.Second)
	require.Equal(t, "", err)

	assert.NotZero(t, createdUser.ID)
	assert.NotZero(t, createdUser.Version)
	assert.NotZero(t, createdUser.CreatedAt)
	assert.NotZero(t, createdUser.UpdatedAt)

	assert.Nil(t, createdUser.DeletedAt)

	assert.Equal(t, testutils.TestName, createdUser.Name)
	assert.Equal(t, testutils.TestSurname, createdUser.Surname)
	assert.Equal(t, testutils.TestUsername, createdUser.Username)
	assert.Equal(t, testutils.BirthDateTest(), createdUser.BirthDate)
	assert.Equal(t, testutils.TestDescription, createdUser.Description)
	assert.Equal(t, testutils.TestEmail, createdUser.Email)
	assert.Equal(t, testutils.TestPhoneNumber, createdUser.PhoneNumber)
}

func TestUsersTransportHTTP_CreateUser_E2E_UserAlreadyExists_ConflictError(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createdUser, err := createUser(t, ctx, 5*time.Second)
	require.Equal(t, "", err)
	require.Equal(t, testutils.TestName, createdUser.Name)
	require.Equal(t, testutils.TestSurname, createdUser.Surname)
	require.Equal(t, testutils.TestUsername, createdUser.Username)
	require.Equal(t, testutils.BirthDateTest(), createdUser.BirthDate)
	require.Equal(t, testutils.TestDescription, createdUser.Description)
	require.Equal(t, testutils.TestEmail, createdUser.Email)
	require.Equal(t, testutils.TestPhoneNumber, createdUser.PhoneNumber)

	createdUserAgain, err := createUser(t, ctx, 5*time.Second)

	require.Contains(t, err, core_errors.ErrConflict.Error())
	require.Nil(t, createdUserAgain)
}
