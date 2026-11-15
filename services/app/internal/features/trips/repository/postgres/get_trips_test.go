package trips_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

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
