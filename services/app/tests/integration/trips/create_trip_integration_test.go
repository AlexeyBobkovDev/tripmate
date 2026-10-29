//go:build integration

package trips_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsService_CreateTrip_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	request := testutils.TripCreateTest(user.ID, "Create trip integration test")

	createdTrip, err := testDeps.tripsService.CreateTrip(ctx, user.ID, request)
	require.NoError(t, err)
	require.NotNil(t, createdTrip)

	assert.NotZero(t, createdTrip.ID)
	assert.Equal(t, user.ID, createdTrip.OwnerID)
	assert.NotZero(t, createdTrip.Version)
	assert.NotZero(t, createdTrip.CreatedAt)
	assert.NotZero(t, createdTrip.UpdatedAt)
	assert.Nil(t, createdTrip.DeletedAt)

	assert.Equal(t, request.Name, createdTrip.Name)
	assert.Equal(t, request.PlanningArea, createdTrip.PlanningArea)
	assert.Equal(t, request.StartLocation, createdTrip.StartLocation)
	assert.True(t, request.StartDate.Equal(createdTrip.StartDate))
	assert.True(t, request.EndDate.Equal(createdTrip.EndDate))
	assert.Equal(t, request.Budget, createdTrip.Budget)
	assert.Equal(t, request.Interests, createdTrip.Interests)
	assert.Equal(t, request.TimeZone, createdTrip.TimeZone)

	createdTripID := createdTrip.ID
	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(
			ctx,
			`DELETE FROM app.trips WHERE owner_id = $1 AND id = $2`,
			user.ID,
			createdTripID,
		)
		return err
	})

	retrievedTrip, err := testDeps.tripsService.GetTrip(ctx, user.ID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)

	assert.Equal(t, createdTrip.ID, retrievedTrip.ID)
	assert.Equal(t, createdTrip.OwnerID, retrievedTrip.OwnerID)
	assert.Equal(t, createdTrip.Version, retrievedTrip.Version)
	assert.Equal(t, createdTrip.Name, retrievedTrip.Name)
	assert.Equal(t, createdTrip.PlanningArea, retrievedTrip.PlanningArea)
	assert.Equal(t, createdTrip.StartLocation, retrievedTrip.StartLocation)
	assert.True(t, createdTrip.StartDate.Equal(retrievedTrip.StartDate))
	assert.True(t, createdTrip.EndDate.Equal(retrievedTrip.EndDate))
	assert.Equal(t, createdTrip.Budget, retrievedTrip.Budget)
	assert.Equal(t, createdTrip.Interests, retrievedTrip.Interests)
	assert.Equal(t, createdTrip.TimeZone, retrievedTrip.TimeZone)
	assert.Equal(t, createdTrip.DeletedAt, retrievedTrip.DeletedAt)
}

func TestTripsService_CreateTrip_Integration_InvalidOwnerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	request := testutils.TripCreateTest(-1, "Invalid owner trip")

	createdTrip, err := testDeps.tripsService.CreateTrip(ctx, -1, request)

	require.Error(t, err)
	require.Nil(t, createdTrip)
}
