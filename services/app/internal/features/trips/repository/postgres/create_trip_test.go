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
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestTripsRepository_CreateTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripCreate := testTripCreateFixture(42)

	expectedTrip := testTripFixture()
	expectedTrip.ID = 10
	expectedTrip.OwnerID = tripCreate.OwnerID
	expectedTrip.Name = tripCreate.Name
	expectedTrip.PlanningArea = tripCreate.PlanningArea
	expectedTrip.StartCoordinates = tripCreate.StartCoordinates
	expectedTrip.StartDate = tripCreate.StartDate
	expectedTrip.EndDate = tripCreate.EndDate
	expectedTrip.Budget = tripCreate.Budget
	expectedTrip.Interests = tripCreate.Interests
	expectedTrip.TimeZone = tripCreate.TimeZone
	expectedTrip.DailySchedule = tripCreate.DailySchedule

	// Expect transaction creation
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	// Expect query on transaction
	mockTx.EXPECT().QueryRow(
		gomock.Any(),
		createTripQuery,
		tripCreate.OwnerID,
		tripCreate.Name,
		tripCreate.PlanningArea,
		tripCreate.StartCoordinates.Latitude,
		tripCreate.StartCoordinates.Longitude,
		tripCreate.StartDate,
		tripCreate.EndDate,
		gomock.Any(), // startTime (will be parsed)
		gomock.Any(), // endTime (will be parsed)
		tripCreate.Budget.AmountMinorUnits,
		tripCreate.Budget.Currency,
		tripCreate.Interests,
		tripCreate.TimeZone,
	).Return(mockRow)

	// Expect scan
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		},
	)

	// Expect transaction commit (success case)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)
	// The defer statement will also call Rollback, but it's a no-op if Commit succeeds
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.CreateTrip(
		context.Background(),
		tripCreate.OwnerID,
		tripCreate,
	)

	require.NoError(t, err)
	require.NotNil(t, actualTrip)
	assert.Equal(t, expectedTrip, *actualTrip)
}

func TestTripsRepository_CreateTrip_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripCreate := testTripCreateFixture(1)
	scanError := errors.New("scan error")

	// Expect transaction creation
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	// Expect query on transaction
	mockTx.EXPECT().QueryRow(
		gomock.Any(),
		createTripQuery,
		tripCreate.OwnerID,
		tripCreate.Name,
		tripCreate.PlanningArea,
		tripCreate.StartCoordinates.Latitude,
		tripCreate.StartCoordinates.Longitude,
		tripCreate.StartDate,
		tripCreate.EndDate,
		gomock.Any(), // startTime (will be parsed)
		gomock.Any(), // endTime (will be parsed)
		tripCreate.Budget.AmountMinorUnits,
		tripCreate.Budget.Currency,
		tripCreate.Interests,
		tripCreate.TimeZone,
	).Return(mockRow)

	// Expect scan error
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).Return(scanError)

	// Expect transaction rollback (error case)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.CreateTrip(
		context.Background(),
		tripCreate.OwnerID,
		tripCreate,
	)

	require.ErrorIs(t, err, scanError)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_CreateTrip_ScheduleRoundTrip(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripCreate := testTripCreateFixture(1)

	// The database returns different time-of-day values than the request on
	// purpose: the repository must trust the stored values.
	expectedTrip := testTripFixture()
	expectedTrip.DailySchedule = domain.DailySchedule{
		StartTime: "07:30",
		EndTime:   "21:45",
	}

	// Expect transaction creation
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	// Expect query on transaction
	mockTx.EXPECT().QueryRow(
		gomock.Any(),
		createTripQuery,
		tripCreate.OwnerID,
		tripCreate.Name,
		tripCreate.PlanningArea,
		tripCreate.StartCoordinates.Latitude,
		tripCreate.StartCoordinates.Longitude,
		tripCreate.StartDate,
		tripCreate.EndDate,
		gomock.Any(), // startTime (will be parsed)
		gomock.Any(), // endTime (will be parsed)
		tripCreate.Budget.AmountMinorUnits,
		tripCreate.Budget.Currency,
		tripCreate.Interests,
		tripCreate.TimeZone,
	).Return(mockRow)

	// Expect scan
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, expectedTrip)(dest...)
		},
	)

	// Expect transaction commit (success case)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)
	// The defer statement will also call Rollback, but it's a no-op if Commit succeeds
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	actualTrip, err := tripsRepository.CreateTrip(
		context.Background(),
		tripCreate.OwnerID,
		tripCreate,
	)

	require.NoError(t, err)
	require.NotNil(t, actualTrip)
	assert.Equal(t, expectedTrip.DailySchedule, actualTrip.DailySchedule)
}
