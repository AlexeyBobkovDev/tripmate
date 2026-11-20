package trips_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestTripsRepository_GetTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1

	expectedTrip := testTripFixture()
	expectedTrip.ID = tripID
	expectedTrip.OwnerID = ownerID

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getTripQuery, ownerID, tripID).
		Return(mockRow)

	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		})

	mockMealRows := mocks.NewMockRows(ctrl)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, expectedTrip.ID).
		Return(mockMealRows, nil)
	mockMealRows.EXPECT().Next().Return(false).AnyTimes()

	actualTrip, err := tripsRepository.GetTrip(context.Background(), ownerID, tripID)

	require.NoError(t, err)
	require.NotNil(t, actualTrip)
	assert.Equal(t, expectedTrip, *actualTrip)
}

func TestTripsRepository_GetTrip_ScheduleRoundTrip(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	expectedTrip := testTripFixture()
	expectedTrip.DailySchedule = domain.DailySchedule{
		StartTime: "06:15",
		EndTime:   "23:40",
	}

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getTripQuery, gomock.Any(), gomock.Any()).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		})

	mockMealRows := mocks.NewMockRows(ctrl)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, gomock.Any()).
		Return(mockMealRows, nil)
	mockMealRows.EXPECT().Next().Return(false).AnyTimes()

	actualTrip, err := tripsRepository.GetTrip(context.Background(), 1, 1)

	require.NoError(t, err)
	require.NotNil(t, actualTrip)
	assert.Equal(t, expectedTrip.DailySchedule, actualTrip.DailySchedule)
}

func TestTripsRepository_GetTrip_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	scanError := errors.New("query row error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, 1).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		Return(scanError)

	actualTrip, err := tripsRepository.GetTrip(context.Background(), 1, 1)

	require.ErrorIs(t, err, scanError)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_GetTrip_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, 1).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		Return(core_postgres_pool_errors.ErrNoRows)

	actualTrip, err := tripsRepository.GetTrip(context.Background(), 1, 1)

	require.ErrorIs(t, err, core_postgres_pool_errors.ErrNoRows)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_GetTrip_NotFoundMappedFromPgx(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, 1).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		Return(errors.Join(pgx.ErrNoRows, core_errors.ErrNotFound))

	actualTrip, err := tripsRepository.GetTrip(context.Background(), 1, 1)

	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_GetTrip_MealRequirementsQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, 1).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, testTripFixture())(dest...)
		})

	mealRequirementsError := errors.New("meal requirements query error")
	mockPool.
		EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, gomock.Any()).
		Return(nil, mealRequirementsError)

	actualTrip, err := tripsRepository.GetTrip(context.Background(), 1, 1)

	require.ErrorIs(t, err, mealRequirementsError)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_GetTrip_MealRequirementsScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, 1).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, testTripFixture())(dest...)
		})

	mockMealRows := mocks.NewMockRows(ctrl)
	mockPool.
		EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, gomock.Any()).
		Return(mockMealRows, nil)
	mockMealRows.EXPECT().Next().Return(true)
	mockMealRows.EXPECT().Scan(anyArgs(6)...).Return(errors.New("meal requirements scan error"))

	actualTrip, err := tripsRepository.GetTrip(context.Background(), 1, 1)

	require.Error(t, err)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_GetTrip_MealRequirements_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	expectedTrip := testTripFixture()
	expectedTrip.ID = 1
	expectedTrip.OwnerID = 1

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, 1).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		})

	mockMealRows := mocks.NewMockRows(ctrl)
	mockPool.
		EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, expectedTrip.ID).
		Return(mockMealRows, nil)
	gomock.InOrder(
		mockMealRows.EXPECT().Next().Return(true),
		mockMealRows.EXPECT().Scan(anyArgs(6)...).DoAndReturn(
			func(dest ...any) error {
				mealStartTime, err := time.Parse("15:04", "08:00")
				require.NoError(t, err)
				mealEndTime, err := time.Parse("15:04", "09:00")
				require.NoError(t, err)

				*dest[0].(*domain.MealType) = domain.MealBreakfast
				*dest[1].(*time.Time) = mealStartTime
				*dest[2].(*time.Time) = mealEndTime
				*dest[3].(*time.Duration) = 30
				*dest[4].(*int64) = int64(1000)
				*dest[5].(*string) = "RUB"
				return nil
			},
		),
		mockMealRows.EXPECT().Next().Return(false),
	)

	actualTrip, err := tripsRepository.GetTrip(context.Background(), 1, 1)

	require.NoError(t, err)
	require.NotNil(t, actualTrip)
	require.Len(t, actualTrip.MealRequirements, 1)
	assert.Equal(t, domain.MealBreakfast, actualTrip.MealRequirements[0].Type)
	assert.Equal(t, 30*time.Minute, actualTrip.MealRequirements[0].Duration)
}
