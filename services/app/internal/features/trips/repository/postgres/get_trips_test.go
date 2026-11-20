package trips_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestTripsRepository_GetTrips_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRows := mocks.NewMockRows(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	expectedTrip := testTripFixture()
	expectedTrip.OwnerID = ownerID

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, ownerID).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
			func(dest ...any) error {
				return fillTripScan(t, expectedTrip)(dest...)
			},
		),
		mockRows.EXPECT().Next().Return(false),
	)

	mockMealRows := mocks.NewMockRows(ctrl)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, expectedTrip.ID).
		Return(mockMealRows, nil)
	mockMealRows.EXPECT().Next().Return(false).AnyTimes()

	trips, err := tripsRepository.GetTrips(context.Background(), ownerID)

	require.NoError(t, err)
	require.NotNil(t, trips)
	require.Len(t, *trips, 1)
	assert.Equal(t, expectedTrip, (*trips)[0])
}

func TestTripsRepository_GetTrips_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRows := mocks.NewMockRows(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(false)

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, trips)
	assert.Empty(t, *trips)
}

func TestTripsRepository_GetTrips_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	queryError := errors.New("query error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(nil, queryError)

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.ErrorIs(t, err, queryError)
	assert.Nil(t, trips)
}

func TestTripsRepository_GetTrips_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRows := mocks.NewMockRows(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	scanError := errors.New("scan error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().
			Scan(anyArgs(tripScanDestinationsAmount)...).
			Return(scanError),
	)

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.ErrorIs(t, err, scanError)
	assert.Nil(t, trips)
}

func TestTripsRepository_GetTrips_NoRowsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRows := mocks.NewMockRows(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().
			Scan(anyArgs(tripScanDestinationsAmount)...).
			Return(core_postgres_pool_errors.ErrNoRows),
	)

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, trips)
}

func TestTripsRepository_GetTrips_TripsQueryNoRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(nil, core_postgres_pool_errors.ErrNoRows)

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, trips)
}

func TestTripsRepository_GetTrips_MealRequirementsQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRows := mocks.NewMockRows(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)
	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
			func(dest ...any) error {
				return fillTripScan(t, testTripFixture())(dest...)
			},
		),
	)
	mockPool.
		EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, gomock.Any()).
		Return(nil, errors.New("meal requirements query error"))

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.Error(t, err)
	assert.Nil(t, trips)
}

func TestTripsRepository_GetTrips_MealRequirements_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRows := mocks.NewMockRows(ctrl)
	mockMealRows := mocks.NewMockRows(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	expectedTrip := testTripFixture()

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)
	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
			func(dest ...any) error {
				return fillTripScan(t, expectedTrip)(dest...)
			},
		),
		mockPool.EXPECT().Query(gomock.Any(), getTripMealRequirementsQuery, gomock.Any()).
			Return(mockMealRows, nil),
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
		mockRows.EXPECT().Next().Return(false),
	)

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, trips)
	require.Len(t, *trips, 1)
	require.Len(t, (*trips)[0].MealRequirements, 1)
	assert.Equal(t, domain.MealBreakfast, (*trips)[0].MealRequirements[0].Type)
}

func TestTripsRepository_GetTrips_MealRequirementsScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRows := mocks.NewMockRows(ctrl)
	mockMealRows := mocks.NewMockRows(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	expectedTrip := testTripFixture()

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.
		EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)
	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
			func(dest ...any) error {
				return fillTripScan(t, expectedTrip)(dest...)
			},
		),
		mockPool.EXPECT().Query(gomock.Any(), getTripMealRequirementsQuery, gomock.Any()).
			Return(mockMealRows, nil),
		mockMealRows.EXPECT().Next().Return(true),
		mockMealRows.EXPECT().Scan(anyArgs(6)...).Return(errors.New("meal scan error")),
	)

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.Error(t, err)
	assert.Nil(t, trips)
}
