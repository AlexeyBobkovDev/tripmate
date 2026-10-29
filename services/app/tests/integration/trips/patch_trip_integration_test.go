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
)

func TestTripsService_PatchTrip_Integration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	createdTrip := createTestTrip(t, ctx, user.ID)

	originalVersion := createdTrip.Version
	originalPlanningArea := createdTrip.PlanningArea
	originalStartLocation := createdTrip.StartLocation
	originalStartDate := createdTrip.StartDate
	originalEndDate := createdTrip.EndDate
	originalBudget := createdTrip.Budget
	originalInterests := createdTrip.Interests
	originalTimeZone := createdTrip.TimeZone

	newName := "Patched integration trip"
	patch := &domain.TripPatch{
		OwnerID: user.ID,
		Version: originalVersion,
		Name: domain.Nullable[string]{
			Value: &newName,
			Set:   true,
		},
	}

	patchedTrip, err := testDeps.tripsService.PatchTrip(
		ctx,
		user.ID,
		createdTrip.ID,
		patch,
	)
	require.NoError(t, err)
	require.NotNil(t, patchedTrip)

	assert.Equal(t, newName, patchedTrip.Name)
	assert.Equal(t, originalVersion+1, patchedTrip.Version)
	assert.Equal(t, createdTrip.ID, patchedTrip.ID)
	assert.Equal(t, createdTrip.OwnerID, patchedTrip.OwnerID)
	assert.Equal(t, originalPlanningArea, patchedTrip.PlanningArea)
	assert.Equal(t, originalStartLocation, patchedTrip.StartLocation)
	assert.True(t, originalStartDate.Equal(patchedTrip.StartDate))
	assert.True(t, originalEndDate.Equal(patchedTrip.EndDate))
	assert.Equal(t, originalBudget, patchedTrip.Budget)
	assert.Equal(t, originalInterests, patchedTrip.Interests)
	assert.Equal(t, originalTimeZone, patchedTrip.TimeZone)

	retrievedTrip, err := testDeps.tripsService.GetTrip(ctx, user.ID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)

	assert.Equal(t, patchedTrip.ID, retrievedTrip.ID)
	assert.Equal(t, newName, retrievedTrip.Name)
	assert.Equal(t, originalVersion+1, retrievedTrip.Version)
	assert.Equal(t, patchedTrip.UpdatedAt, retrievedTrip.UpdatedAt)
}

func TestTripsService_PatchTrip_Integration_Conflict_MismatchingVersions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	createdTrip := createTestTrip(t, ctx, user.ID)

	newName := "Should not be applied"
	patch := &domain.TripPatch{
		OwnerID: user.ID,
		Version: createdTrip.Version + 1,
		Name: domain.Nullable[string]{
			Value: &newName,
			Set:   true,
		},
	}

	patchedTrip, err := testDeps.tripsService.PatchTrip(
		ctx,
		user.ID,
		createdTrip.ID,
		patch,
	)

	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.Nil(t, patchedTrip)

	retrievedTrip, err := testDeps.tripsService.GetTrip(ctx, user.ID, createdTrip.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTrip)

	assert.Equal(t, createdTrip.Version, retrievedTrip.Version)
	assert.Equal(t, createdTrip.Name, retrievedTrip.Name)
}

func TestTripsService_PatchTrip_Integration_InvalidOwnerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID)

	newName := "Invalid owner patch"
	patch := &domain.TripPatch{
		OwnerID: -1,
		Version: trip.Version,
		Name: domain.Nullable[string]{
			Value: &newName,
			Set:   true,
		},
	}

	patchedTrip, err := testDeps.tripsService.PatchTrip(
		ctx,
		-1,
		trip.ID,
		patch,
	)

	require.ErrorIs(t, err, core_errors.ErrBadRequest)
	require.Nil(t, patchedTrip)
}

func TestTripsService_PatchTrip_Integration_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	newName := "Missing trip patch"

	patch := &domain.TripPatch{
		OwnerID: user.ID,
		Version: 1,
		Name: domain.Nullable[string]{
			Value: &newName,
			Set:   true,
		},
	}

	patchedTrip, err := testDeps.tripsService.PatchTrip(
		ctx,
		user.ID,
		2_000_000_000,
		patch,
	)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, patchedTrip)
}

func TestTripsService_PatchTrip_Integration_AlreadyDeleted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTestTrip(t, ctx, user.ID)

	require.NoError(t, testDeps.tripsService.DeleteTrip(ctx, user.ID, trip.ID))

	newName := "Should not update deleted trip"
	patch := &domain.TripPatch{
		OwnerID: user.ID,
		Version: trip.Version,
		Name: domain.Nullable[string]{
			Value: &newName,
			Set:   true,
		},
	}

	patchedTrip, err := testDeps.tripsService.PatchTrip(
		ctx,
		user.ID,
		trip.ID,
		patch,
	)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, patchedTrip)
}
