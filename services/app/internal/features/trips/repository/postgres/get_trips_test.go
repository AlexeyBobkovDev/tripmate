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
	expectedTrip.ID = 1

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, ownerID).
		Return(mockRows, nil)

	// Mock the trip ID row
	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).DoAndReturn(
			func(dest ...any) error {
				if len(dest) > 0 {
					if id, ok := dest[0].(*int); ok {
						*id = expectedTrip.ID
					}
				}
				return nil
			},
		),
		mockRows.EXPECT().Next().Return(false),
		mockRows.EXPECT().Err().Return(nil),
	)
	mockRows.EXPECT().Close()

	// Mock the getTrip QueryRow
	mockRow := mocks.NewMockRow(ctrl)
	mockPool.EXPECT().
		QueryRow(gomock.Any(), getTripQuery, ownerID, expectedTrip.ID).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		})

	// Mock the meal requirements query
	mockMealRows := mocks.NewMockRows(ctrl)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, expectedTrip.ID).
		Return(mockMealRows, nil)
	mockMealRows.EXPECT().Next().Return(false).AnyTimes()
	mockMealRows.EXPECT().Err().Return(nil)
	mockMealRows.EXPECT().Close()

	trips, err := tripsRepository.GetTrips(context.Background(), ownerID)

	require.NoError(t, err)
	require.NotNil(t, trips)
	require.Len(t, *trips, 1)
	assert.Equal(t, expectedTrip, (*trips)[0])
}

func TestTripsRepository_GetTrips_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(nil, core_postgres_pool_errors.ErrNoRows)

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, trips)
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
			Scan(gomock.Any()).
			Return(scanError),
	)
	mockRows.EXPECT().Close()

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
			Scan(gomock.Any()).
			Return(core_postgres_pool_errors.ErrNoRows),
	)
	mockRows.EXPECT().Close()

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	var codedErr *core_errors.CodedError
	require.ErrorAs(t, err, &codedErr)
	require.Equal(t, core_errors.CategoryDatabase, codedErr.Category)
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
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).DoAndReturn(
			func(dest ...any) error {
				if len(dest) > 0 {
					if id, ok := dest[0].(*int); ok {
						*id = 1
					}
				}
				return nil
			},
		),
		mockRows.EXPECT().Next().Return(false),
		mockRows.EXPECT().Err().Return(nil),
	)
	mockRows.EXPECT().Close()

	mockPool.EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, gomock.Any()).
		Return(nil, errors.New("meal requirements query error"))

	// Mock the getTrip QueryRow
	mockRow := mocks.NewMockRow(ctrl)
	mockPool.EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, 1).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, testTripFixture())(dest...)
		})

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
	expectedTrip.ID = 1
	expectedTrip.OwnerID = 1

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).DoAndReturn(
			func(dest ...any) error {
				if len(dest) > 0 {
					if id, ok := dest[0].(*int); ok {
						*id = expectedTrip.ID
					}
				}
				return nil
			},
		),
		mockRows.EXPECT().Next().Return(false),
		mockRows.EXPECT().Err().Return(nil),
	)
	mockRows.EXPECT().Close()

	mockRow := mocks.NewMockRow(ctrl)
	mockPool.EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, expectedTrip.ID).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		})

	mockPool.EXPECT().
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
		mockMealRows.EXPECT().Err().Return(nil),
	)
	mockMealRows.EXPECT().Close()

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
	expectedTrip.ID = 1
	expectedTrip.OwnerID = 1

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).DoAndReturn(
			func(dest ...any) error {
				if len(dest) > 0 {
					if id, ok := dest[0].(*int); ok {
						*id = expectedTrip.ID
					}
				}
				return nil
			},
		),
		mockRows.EXPECT().Next().Return(false),
		mockRows.EXPECT().Err().Return(nil),
	)
	mockRows.EXPECT().Close()

	mockRow := mocks.NewMockRow(ctrl)
	mockPool.EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, expectedTrip.ID).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		})

	mockPool.EXPECT().
		Query(gomock.Any(), getTripMealRequirementsQuery, expectedTrip.ID).
		Return(mockMealRows, nil)
	gomock.InOrder(
		mockMealRows.EXPECT().Next().Return(true),
		mockMealRows.EXPECT().Scan(anyArgs(6)...).Return(errors.New("meal scan error")),
		mockMealRows.EXPECT().Close(),
	)

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.Error(t, err)
	assert.Nil(t, trips)
}

func TestTripsRepository_GetTrips_MealRequirementsErrCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRows := mocks.NewMockRows(ctrl)
	mockMealRows := mocks.NewMockRows(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	expectedTrip := testTripFixture()
	expectedTrip.ID = 1
	expectedTrip.OwnerID = 1

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().
		Query(gomock.Any(), getTripsQuery, 1).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).DoAndReturn(
			func(dest ...any) error {
				if len(dest) > 0 {
					if id, ok := dest[0].(*int); ok {
						*id = expectedTrip.ID
					}
				}
				return nil
			},
		),
		mockRows.EXPECT().Next().Return(false),
		mockRows.EXPECT().Err().Return(nil),
	)
	mockRows.EXPECT().Close()

	mockRow := mocks.NewMockRow(ctrl)
	mockPool.EXPECT().
		QueryRow(gomock.Any(), getTripQuery, 1, expectedTrip.ID).
		Return(mockRow)
	mockRow.EXPECT().
		Scan(anyArgs(tripScanDestinationsAmount)...).
		DoAndReturn(func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		})

	mockPool.EXPECT().
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
		mockMealRows.EXPECT().Err().Return(errors.New("meal rows error")),
	)
	mockMealRows.EXPECT().Close()

	trips, err := tripsRepository.GetTrips(context.Background(), 1)

	require.Error(t, err)
	assert.Nil(t, trips)
}
