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

func TestUsersService_DeleteUser_Integration_Success(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createUserRequest := testutils.UserCreateTest()

	actualUser, err := testDeps.usersRepository.CreateUser(ctx, &createUserRequest)

	require.NoError(t, err)
	require.NotNil(t, actualUser)

	userID := actualUser.ID
	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, userID)
		return err
	})

	err = testDeps.usersService.DeleteUser(ctx, actualUser.ID)
	require.NoError(t, err)

	var deletedAt time.Time
	err = testDeps.pool.QueryRow(ctx, `SELECT deleted_at FROM app.users WHERE id=$1;`, userID).Scan(&deletedAt)
	require.NoError(t, err)
	assert.False(t, deletedAt.IsZero())
}

func TestUsersService_DeleteUser_Integration_UserNotFound(t *testing.T) {
	checkIfDatabaseIsDirty(t)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userID := -1

	err := testDeps.usersService.DeleteUser(ctx, userID)
	require.ErrorIs(t, err, core_errors.ErrBadRequest)
}
