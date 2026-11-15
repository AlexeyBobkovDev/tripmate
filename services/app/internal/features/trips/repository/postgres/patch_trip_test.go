package trips_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestTripsRepository_PatchTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
	tripPatch := testTripFixture()
	tripPatch.ID = tripID
	tripPatch.OwnerID = ownerID
	tripPatch.Version = 2
	tripPatch.Name = "Updated trip"

	expectedTrip := tripPatch
	expectedTrip.Version = 3
	expectedTrip.DailySchedule = tripPatch.DailySchedule

	// Expect transaction creation
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	// Expect query on transaction
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			patchTripQuery,
			ownerID,
			tripID,
			tripPatch.Version,
			tripPatch.Name,
			tripPatch.PlanningArea,
			tripPatch.StartCoordinates.Latitude,
			tripPatch.StartCoordinates.Longitude,
			tripPatch.StartDate,
			tripPatch.EndDate,
			gomock.Any(), // startTime (will be parsed)
			gomock.Any(), // endTime (will be parsed)
			tripPatch.Budget.AmountMinorUnits,
			tripPatch.Budget.Currency,
			tripPatch.Interests,
			tripPatch.TimeZone,
		).
		Return(mockRow)

	// Expect scan
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		})

	// Expect deletion of existing meal requirements
	mockTx.EXPECT().Exec(gomock.Any(), deleteTripMealRequirementsQuery, tripID).Return(pgconn.CommandTag{}, nil)

	// Expect insertion of new meal requirements (if any)
	if len(tripPatch.MealRequirements) > 0 {
		mockTx.EXPECT().
			QueryRow(
				gomock.Any(),
				insertTripMealRequirementsQuery,
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
				gomock.Any(),
			).
			Return(mockRow)
	}

	// Expect transaction commit (success case)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)
	// In the production code, there's a defer that calls Rollback,
	// but it's a no-op if Commit succeeds
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.PatchTrip(
		context.Background(),
		ownerID,
		tripID,
		&tripPatch,
	)

	require.NoError(t, err)
	require.NotNil(t, actualTrip)
	assert.Equal(t, expectedTrip, *actualTrip)
}

func TestTripsRepository_PatchTrip_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripPatch := testTripFixture()
	updateError := errors.New("update error")

	// Expect transaction creation
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	// Expect query on transaction
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			patchTripQuery,
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		Return(updateError)

	// Expect transaction rollback (error case)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.PatchTrip(
		context.Background(),
		1,
		1,
		&tripPatch,
	)

	require.ErrorIs(t, err, updateError)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_PatchTrip_VersionConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockExistsRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripPatch := testTripFixture()

	// Expect transaction creation
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	// Expect query on transaction - returns no rows
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			patchTripQuery,
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		Return(core_postgres_pool_errors.ErrNoRows)

	// Expect query on pool (not transaction) for exists check
	mockPool.EXPECT().
		QueryRow(gomock.Any(), tripExistsQuery, 1, 1).
		Return(mockExistsRow)
	mockExistsRow.EXPECT().
		Scan(gomock.Any()).
		DoAndReturn(func(dest ...any) error {
			*dest[0].(*bool) = true
			return nil
		})

	// Expect transaction rollback (error case)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.PatchTrip(
		context.Background(),
		1,
		1,
		&tripPatch,
	)

	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.ErrorIs(t, err, core_postgres_pool_errors.ErrNoRows)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_PatchTrip_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockExistsRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripPatch := testTripFixture()

	// Expect transaction creation
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	// Expect query on transaction - returns no rows
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			patchTripQuery,
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		Return(core_postgres_pool_errors.ErrNoRows)

	// Expect query on pool (not transaction) for exists check
	mockPool.EXPECT().
		QueryRow(gomock.Any(), tripExistsQuery, 1, 1).
		Return(mockExistsRow)
	mockExistsRow.EXPECT().
		Scan(gomock.Any()).
		DoAndReturn(func(dest ...any) error {
			*dest[0].(*bool) = false
			return nil
		})

	// Expect transaction rollback (error case)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.PatchTrip(
		context.Background(),
		1,
		1,
		&tripPatch,
	)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_PatchTrip_ExistsCheckError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockExistsRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripPatch := testTripFixture()
	existsCheckError := errors.New("exists check error")

	// Expect transaction creation
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	// Expect query on transaction - returns no rows
	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			patchTripQuery,
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		Return(core_postgres_pool_errors.ErrNoRows)

	// Expect query on pool (not transaction) for exists check
	mockPool.EXPECT().
		QueryRow(gomock.Any(), tripExistsQuery, 1, 1).
		Return(mockExistsRow)
	mockExistsRow.EXPECT().
		Scan(gomock.Any()).
		Return(existsCheckError)

	// Expect transaction rollback (error case)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.PatchTrip(
		context.Background(),
		1,
		1,
		&tripPatch,
	)

	require.ErrorIs(t, err, existsCheckError)
	assert.Nil(t, actualTrip)
}
