//go:build integration

package users_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func checkIfDatabaseIsDirty(t *testing.T) {
	var counter int

	require.NoError(t, poolForIntegrationTesting.QueryRow(
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
