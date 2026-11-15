//go:build integration

package trips_integration_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

var testUserSequence atomic.Int64

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
	sequence := testUserSequence.Add(1)
	suffix := fmt.Sprintf("%d_%d", time.Now().UnixNano(), sequence)
	// app.users has unique indexes on username, email and phone_number, so every
	// test user must be unique even when a test creates more than one user.
	request.Username = "integration_user_" + suffix
	request.Email = "integration_" + suffix + "@example.com"
	request.PhoneNumber = fmt.Sprintf("+7999%07d", sequence%10_000_000)

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

// deleteTripRows removes a trip and everything that references it. Meal
// requirements reference app.trips without ON DELETE CASCADE, so they must be
// deleted first.
func deleteTripRows(ctx context.Context, ownerID int, tripID int) error {
	if _, err := testDeps.pool.Exec(
		ctx,
		`DELETE FROM app.meal_requirements WHERE id = $1;`,
		tripID,
	); err != nil {
		return err
	}

	_, err := testDeps.pool.Exec(
		ctx,
		`DELETE FROM app.trips WHERE owner_id = $1 AND id = $2;`,
		ownerID,
		tripID,
	)
	return err
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
		return deleteTripRows(ctx, ownerID, trip.ID)
	})

	return trip
}
