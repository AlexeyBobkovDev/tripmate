//go:build integration

package trips_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestTripsService_DeleteTrip_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID)

	err := testDeps.tripsService.DeleteTrip(ctx, user.ID, trip.ID)
	require.NoError(t, err)

	retrievedTrip, err := testDeps.tripsService.GetTrip(ctx, user.ID, trip.ID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedTrip)
}

func TestTripsService_DeleteTrip_Integration_InvalidOwnerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := testDeps.tripsService.DeleteTrip(ctx, -1, 1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
}

func TestTripsService_DeleteTrip_Integration_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	err := testDeps.tripsService.DeleteTrip(ctx, user.ID, -1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
}

func TestTripsService_DeleteTrip_Integration_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	err := testDeps.tripsService.DeleteTrip(ctx, user.ID, 2_000_000_000)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
}

func TestTripsService_DeleteTrip_Integration_AlreadyDeleted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID)

	require.NoError(t, testDeps.tripsService.DeleteTrip(ctx, user.ID, trip.ID))

	err := testDeps.tripsService.DeleteTrip(ctx, user.ID, trip.ID)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
}
