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

func TestTripsService_GetTrips_Integration_Empty(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	trips, err := testDeps.tripsService.GetTrips(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, trips)
	assert.Empty(t, *trips)
}

func TestTripsService_GetTrips_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	firstTrip := createTestTrip(t, ctx, user.ID)
	secondTrip, err := testDeps.tripsService.CreateTrip(
		ctx,
		user.ID,
		testutils.TripCreateTest(user.ID, "Second integration trip"),
	)
	require.NoError(t, err)
	require.NotNil(t, secondTrip)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(
			ctx,
			`DELETE FROM app.trips WHERE owner_id = $1 AND id = $2`,
			user.ID,
			secondTrip.ID,
		)
		return err
	})

	trips, err := testDeps.tripsService.GetTrips(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, trips)
	require.Len(t, *trips, 2)

	gotByID := make(map[int]string, len(*trips))
	for _, trip := range *trips {
		assert.Equal(t, user.ID, trip.OwnerID)
		assert.Nil(t, trip.DeletedAt)
		gotByID[trip.ID] = trip.Name
	}

	assert.Equal(t, firstTrip.Name, gotByID[firstTrip.ID])
	assert.Equal(t, secondTrip.Name, gotByID[secondTrip.ID])
}

func TestTripsService_GetTrips_Integration_OnlyReturnsOwnerTrips(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	firstOwner := createTestUser(t, ctx)

	birthDate := testutils.BirthDateTest()
	request := domain.UserCreate{
		Name:        testutils.TestName,
		Surname:     testutils.TestSurname,
		Username:    "new_username",
		BirthDate:   birthDate,
		Description: testutils.TestDescription,
		Email:       "trip-test-user-2@email.com",
		PhoneNumber: "+79800454867",
	}
	secondOwner, err := testDeps.usersService.CreateUser(
		ctx,
		&request,
		testutils.TestPassword,
	)
	require.NoError(t, err)
	require.NotNil(t, secondOwner)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(
			ctx,
			`DELETE FROM app.users WHERE id = $1;`,
			secondOwner.ID,
		)
		return err
	})

	ownedTrip := createTestTrip(t, ctx, firstOwner.ID)
	_ = createTestTrip(t, ctx, secondOwner.ID)

	trips, err := testDeps.tripsService.GetTrips(ctx, firstOwner.ID)
	require.NoError(t, err)
	require.NotNil(t, trips)
	require.Len(t, *trips, 1)

	assert.Equal(t, ownedTrip.ID, (*trips)[0].ID)
	assert.Equal(t, firstOwner.ID, (*trips)[0].OwnerID)
	assert.Equal(t, ownedTrip.Name, (*trips)[0].Name)
}

func TestTripsService_GetTrips_Integration_InvalidOwnerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	trips, err := testDeps.tripsService.GetTrips(ctx, -1)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, trips)
}
