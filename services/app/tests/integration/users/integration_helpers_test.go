//go:build integration

package users_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func registerCleanup(t *testing.T, cleanup func(context.Context) error) {
	t.Helper()

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()

		if err := cleanup(cleanupCtx); err != nil {
			t.Errorf("integration test cleanup failed: %v", err)
		}
	})
}

func createUser(t *testing.T, ctx context.Context, createUserRequest *domain.UserCreate, password string) *domain.User {
	actualUser, err := testDeps.usersService.CreateUser(ctx, createUserRequest, password)
	require.NoError(t, err)
	require.NotNil(t, actualUser)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(ctx, `DELETE FROM app.users WHERE id=$1;`, actualUser.ID)
		return err
	})

	return actualUser
}

func checkIfDatabaseIsDirty(t *testing.T) {
	var counter int

	require.NoError(t, testDeps.pool.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM app.users;
		`,
	).Scan(&counter))
	if counter != 0 {
		panic("dirty database state!!!")
	}
}
