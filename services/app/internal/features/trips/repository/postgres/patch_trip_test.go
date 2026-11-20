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

	domain "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
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

func TestTripsRepository_PatchTrip_BeginError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripPatch := testTripFixture()

	beginError := errors.New("begin error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, beginError)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), 1, 1, &tripPatch)

	require.ErrorIs(t, err, beginError)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_PatchTrip_InvalidDailyScheduleTime(t *testing.T) {
	tests := []struct {
		name          string
		startTime     string
		endTime       string
		expectedError error
	}{
		{
			name:          "invalid start time",
			startTime:     "not-a-time",
			endTime:       "18:00",
			expectedError: core_errors.ErrInvalidTripStartTime,
		},
		{
			name:          "invalid end time",
			startTime:     "09:00",
			endTime:       "not-a-time",
			expectedError: core_errors.ErrInvalidTripEndTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPool := mocks.NewMockPool(ctrl)
			mockTx := mocks.NewMockTx(ctrl)
			tripsRepository := NewTripsRepository(mockPool)
			mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
			mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
			mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
			tripPatch := testTripFixture()
			tripPatch.DailySchedule = domain.DailySchedule{
				StartTime: domain.LocalTime(tt.startTime),
				EndTime:   domain.LocalTime(tt.endTime),
			}

			actualTrip, err := tripsRepository.PatchTrip(context.Background(), 1, 1, &tripPatch)

			require.ErrorIs(t, err, tt.expectedError)
			assert.Nil(t, actualTrip)
		})
	}
}

func TestTripsRepository_PatchTrip_MealRequirementsDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripPatch := testTripFixture()
	tripPatch.ID = 1
	tripPatch.OwnerID = 1
	tripPatch.Version = 2

	deleteError := errors.New("delete meal requirements error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

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
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, tripPatch)(dest...)
		},
	)

	mockTx.EXPECT().
		Exec(gomock.Any(), deleteTripMealRequirementsQuery, gomock.Any()).
		Return(mockCommandTag, deleteError)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), 1, 1, &tripPatch)

	require.ErrorIs(t, err, deleteError)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_PatchTrip_MealRequirements_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockMealRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripPatch := testTripFixture()
	tripPatch.ID = 1
	tripPatch.OwnerID = 1
	tripPatch.Version = 2
	tripPatch.MealRequirements = []domain.MealRequirement{
		{
			Type:     domain.MealBreakfast,
			Window:   domain.TimeWindow{Start: domain.LocalTime("08:00"), End: domain.LocalTime("09:00")},
			Duration: 30 * time.Minute,
			Budget:   domain.Money{AmountMinorUnits: 1000, Currency: "RUB"},
		},
	}

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

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
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, tripPatch)(dest...)
		},
	)

	mockTx.EXPECT().
		Exec(gomock.Any(), deleteTripMealRequirementsQuery, gomock.Any()).
		Return(mockCommandTag, nil)
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
		Return(mockMealRow)
	mealStartTime, err := time.Parse("15:04", "08:00")
	require.NoError(t, err)
	mealEndTime, err := time.Parse("15:04", "09:00")
	require.NoError(t, err)
	mockMealRow.EXPECT().Scan(anyArgs(6)...).DoAndReturn(
		func(dest ...any) error {
			*dest[0].(*domain.MealType) = domain.MealBreakfast
			*dest[1].(*time.Time) = mealStartTime
			*dest[2].(*time.Time) = mealEndTime
			*dest[3].(*time.Duration) = 30
			*dest[4].(*int64) = int64(1000)
			*dest[5].(*string) = "RUB"
			return nil
		},
	)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), 1, 1, &tripPatch)

	require.NoError(t, err)
	require.NotNil(t, actualTrip)
	require.Len(t, actualTrip.MealRequirements, 1)
	assert.Equal(t, domain.MealBreakfast, actualTrip.MealRequirements[0].Type)
	assert.Equal(t, 30*time.Minute, actualTrip.MealRequirements[0].Duration)
}

func TestTripsRepository_PatchTrip_MealRequirementInsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockMealRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripPatch := testTripFixture()
	tripPatch.ID = 1
	tripPatch.OwnerID = 1
	tripPatch.Version = 2
	tripPatch.MealRequirements = []domain.MealRequirement{
		{
			Type:     domain.MealBreakfast,
			Window:   domain.TimeWindow{Start: domain.LocalTime("08:00"), End: domain.LocalTime("09:00")},
			Duration: 30 * time.Minute,
			Budget:   domain.Money{AmountMinorUnits: 1000, Currency: "RUB"},
		},
	}
	insertError := errors.New("meal insert error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

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
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, tripPatch)(dest...)
		},
	)

	mockTx.EXPECT().
		Exec(gomock.Any(), deleteTripMealRequirementsQuery, gomock.Any()).
		Return(mockCommandTag, nil)
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
		Return(mockMealRow)
	mockMealRow.EXPECT().Scan(anyArgs(6)...).Return(insertError)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), 1, 1, &tripPatch)

	require.ErrorIs(t, err, insertError)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_PatchTrip_MealRequirementInvalidStartTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	tripPatch := testTripFixture()
	tripPatch.MealRequirements = []domain.MealRequirement{
		{
			Type:     domain.MealBreakfast,
			Window:   domain.TimeWindow{Start: domain.LocalTime("not-a-time"), End: domain.LocalTime("09:00")},
			Duration: 30 * time.Minute,
			Budget:   domain.Money{AmountMinorUnits: 1000, Currency: "RUB"},
		},
	}

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			patchTripQuery,
			1,
			1,
			tripPatch.Version,
			tripPatch.Name,
			tripPatch.PlanningArea,
			tripPatch.StartCoordinates.Latitude,
			tripPatch.StartCoordinates.Longitude,
			tripPatch.StartDate,
			tripPatch.EndDate,
			gomock.Any(),
			gomock.Any(),
			tripPatch.Budget.AmountMinorUnits,
			tripPatch.Budget.Currency,
			tripPatch.Interests,
			tripPatch.TimeZone,
		).
		Return(mockRow)
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, tripPatch)(dest...)
		},
	)
	mockTx.EXPECT().Exec(gomock.Any(), deleteTripMealRequirementsQuery, gomock.Any()).Return(mockCommandTag, nil)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), 1, 1, &tripPatch)

	require.ErrorIs(t, err, core_errors.ErrInvalidTripStartTime)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_PatchTrip_MealRequirementInvalidEndTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	tripPatch := testTripFixture()
	tripPatch.MealRequirements = []domain.MealRequirement{
		{
			Type:     domain.MealBreakfast,
			Window:   domain.TimeWindow{Start: domain.LocalTime("08:00"), End: domain.LocalTime("not-a-time")},
			Duration: 30 * time.Minute,
			Budget:   domain.Money{AmountMinorUnits: 1000, Currency: "RUB"},
		},
	}

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			patchTripQuery,
			1,
			1,
			tripPatch.Version,
			tripPatch.Name,
			tripPatch.PlanningArea,
			tripPatch.StartCoordinates.Latitude,
			tripPatch.StartCoordinates.Longitude,
			tripPatch.StartDate,
			tripPatch.EndDate,
			gomock.Any(),
			gomock.Any(),
			tripPatch.Budget.AmountMinorUnits,
			tripPatch.Budget.Currency,
			tripPatch.Interests,
			tripPatch.TimeZone,
		).
		Return(mockRow)
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, tripPatch)(dest...)
		},
	)
	mockTx.EXPECT().Exec(gomock.Any(), deleteTripMealRequirementsQuery, gomock.Any()).Return(mockCommandTag, nil)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), 1, 1, &tripPatch)

	require.ErrorIs(t, err, core_errors.ErrInvalidTripEndTime)
	assert.Nil(t, actualTrip)
}
