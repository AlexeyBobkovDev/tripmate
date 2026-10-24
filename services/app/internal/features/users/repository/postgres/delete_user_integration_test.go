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

func TestUsersRepository_DeleteUser_Integration_Success(t *testing.T) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	usersRepository := NewUsersRepository(poolForIntegrationTesting)

	userID := -1

	err := usersRepository.DeleteUser(ctx, userID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}
