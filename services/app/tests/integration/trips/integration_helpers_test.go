//go:build integration

package trips_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func registerCleanup(t *testing.T, cleanup func(context.Context) error) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := cleanup(ctx); err != nil {
			t.Errorf("integration test cleanup failed: %v", err)
		}
	})
}

func createTestUser(t *testing.T, ctx context.Context) *domain.User {
	t.Helper()

	request := testutils.UserCreateTest()

	user, err := testDeps.usersService.CreateUser(
		ctx,
		&request,
		testutils.TestPassword,
	)
	require.NoError(t, err)
	require.NotNil(t, user)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(
			ctx,
			`DELETE FROM app.users WHERE id = $1;`,
			user.ID,
		)
		return err
	})

	return user
}

func createTestTrip(t *testing.T, ctx context.Context, ownerID int) *domain.Trip {
	t.Helper()

	trip, err := testDeps.tripsService.CreateTrip(
		ctx,
		ownerID,
		testutils.TripCreateTest(ownerID, "Integration test trip"),
	)
	require.NoError(t, err)
	require.NotNil(t, trip)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(
			ctx,
			`DELETE FROM app.trips WHERE owner_id = $1 AND id = $2`,
			ownerID,
			trip.ID,
		)
		return err
	})

	return trip
}
