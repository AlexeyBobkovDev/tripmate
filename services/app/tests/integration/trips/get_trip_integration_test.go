//go:build integration

package trips_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestTripsService_GetTrip_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	createdTrip := createTestTrip(t, ctx, user.ID)

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
	assert.Equal(t, createdTrip.CreatedAt, retrievedTrip.CreatedAt)
	assert.Equal(t, createdTrip.UpdatedAt, retrievedTrip.UpdatedAt)
	assert.Equal(t, createdTrip.DeletedAt, retrievedTrip.DeletedAt)
}

func TestTripsService_GetTrip_Integration_InvalidOwnerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	trip, err := testDeps.tripsService.GetTrip(ctx, -1, 1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, trip)
}

func TestTripsService_GetTrip_Integration_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	trip, err := testDeps.tripsService.GetTrip(ctx, user.ID, -1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, trip)
}

func TestTripsService_GetTrip_Integration_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	trip, err := testDeps.tripsService.GetTrip(ctx, user.ID, 2_000_000_000)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, trip)
}

func TestTripsService_GetTrip_Integration_WrongOwner(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	owner := createTestUser(t, ctx)

	birthDate := testutils.BirthDateTest()
	anotherUser := domain.User{
		ID:          owner.ID + 1,
		Version:     1,
		Name:        testutils.TestName,
		Surname:     testutils.TestSurname,
		Username:    "new_username",
		BirthDate:   birthDate,
		Description: testutils.TestDescription,
		Email:       "trip-test@email",
		PhoneNumber: "+7980045486",
		CreatedAt:   testutils.CreatedAtTest(),
		UpdatedAt:   testutils.UpdatedAtTest(),
		DeletedAt:   nil,
	}

	createdTrip := createTestTrip(t, ctx, owner.ID)

	trip, err := testDeps.tripsService.GetTrip(ctx, anotherUser.ID, createdTrip.ID)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, trip)
}
