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

func TestTripsRepository_CreateTrip_BeginError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripCreate := testTripCreateFixture(1)
	beginError := errors.New("begin error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(nil, beginError)

	actualTrip, err := tripsRepository.CreateTrip(context.Background(), tripCreate.OwnerID, tripCreate)

	require.ErrorIs(t, err, beginError)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_CreateTrip_InvalidDailyScheduleTime(t *testing.T) {
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
			tripCreate := testTripCreateFixture(1)
			tripCreate.DailySchedule = domain.DailySchedule{
				StartTime: domain.LocalTime(tt.startTime),
				EndTime:   domain.LocalTime(tt.endTime),
			}

			actualTrip, err := tripsRepository.CreateTrip(context.Background(), tripCreate.OwnerID, tripCreate)

			require.ErrorIs(t, err, tt.expectedError)
			assert.Nil(t, actualTrip)
		})
	}
}

func TestTripsRepository_CreateTrip_MealRequirementInvalidStartTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	tripCreate := testTripCreateFixture(1)
	tripCreate.MealRequirements = []domain.MealRequirement{
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
			createTripQuery,
			tripCreate.OwnerID,
			tripCreate.Name,
			tripCreate.PlanningArea,
			tripCreate.StartCoordinates.Latitude,
			tripCreate.StartCoordinates.Longitude,
			tripCreate.StartDate,
			tripCreate.EndDate,
			gomock.Any(),
			gomock.Any(),
			tripCreate.Budget.AmountMinorUnits,
			tripCreate.Budget.Currency,
			tripCreate.Interests,
			tripCreate.TimeZone,
		).
		Return(mockRow)
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, testTripFixture())(dest...)
		},
	)

	actualTrip, err := tripsRepository.CreateTrip(context.Background(), tripCreate.OwnerID, tripCreate)

	require.ErrorIs(t, err, core_errors.ErrInvalidTripStartTime)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_CreateTrip_MealRequirementInvalidEndTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	tripCreate := testTripCreateFixture(1)
	tripCreate.MealRequirements = []domain.MealRequirement{
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
			createTripQuery,
			tripCreate.OwnerID,
			tripCreate.Name,
			tripCreate.PlanningArea,
			tripCreate.StartCoordinates.Latitude,
			tripCreate.StartCoordinates.Longitude,
			tripCreate.StartDate,
			tripCreate.EndDate,
			gomock.Any(),
			gomock.Any(),
			tripCreate.Budget.AmountMinorUnits,
			tripCreate.Budget.Currency,
			tripCreate.Interests,
			tripCreate.TimeZone,
		).
		Return(mockRow)
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, testTripFixture())(dest...)
		},
	)

	actualTrip, err := tripsRepository.CreateTrip(context.Background(), tripCreate.OwnerID, tripCreate)

	require.ErrorIs(t, err, core_errors.ErrInvalidTripEndTime)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_CreateTrip_MealRequirements_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockMealRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripCreate := testTripCreateFixture(1)
	tripCreate.MealRequirements = []domain.MealRequirement{
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
			createTripQuery,
			tripCreate.OwnerID,
			tripCreate.Name,
			tripCreate.PlanningArea,
			tripCreate.StartCoordinates.Latitude,
			tripCreate.StartCoordinates.Longitude,
			tripCreate.StartDate,
			tripCreate.EndDate,
			gomock.Any(),
			gomock.Any(),
			tripCreate.Budget.AmountMinorUnits,
			tripCreate.Budget.Currency,
			tripCreate.Interests,
			tripCreate.TimeZone,
		).
		Return(mockRow)
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, testTripFixture())(dest...)
		},
	)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createTripMealRequirementsQuery,
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

	actualTrip, err := tripsRepository.CreateTrip(context.Background(), tripCreate.OwnerID, tripCreate)

	require.NoError(t, err)
	require.NotNil(t, actualTrip)
	require.Len(t, actualTrip.MealRequirements, 1)
	assert.Equal(t, domain.MealBreakfast, actualTrip.MealRequirements[0].Type)
	assert.Equal(t, 30*time.Minute, actualTrip.MealRequirements[0].Duration)
	assert.Equal(t, domain.LocalTime("08:00"), actualTrip.MealRequirements[0].Window.Start)
	assert.Equal(t, domain.LocalTime("09:00"), actualTrip.MealRequirements[0].Window.End)
	assert.Equal(t, int64(1000), actualTrip.MealRequirements[0].Budget.AmountMinorUnits)
}

func TestTripsRepository_CreateTrip_MealRequirementsInsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockMealRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	tripCreate := testTripCreateFixture(1)
	tripCreate.MealRequirements = []domain.MealRequirement{
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
			createTripQuery,
			tripCreate.OwnerID,
			tripCreate.Name,
			tripCreate.PlanningArea,
			tripCreate.StartCoordinates.Latitude,
			tripCreate.StartCoordinates.Longitude,
			tripCreate.StartDate,
			tripCreate.EndDate,
			gomock.Any(),
			gomock.Any(),
			tripCreate.Budget.AmountMinorUnits,
			tripCreate.Budget.Currency,
			tripCreate.Interests,
			tripCreate.TimeZone,
		).
		Return(mockRow)
	mockRow.EXPECT().Scan(anyArgs(tripScanDestinationsAmount)...).DoAndReturn(
		func(dest ...any) error {
			return fillTripScan(t, testTripFixture())(dest...)
		},
	)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			createTripMealRequirementsQuery,
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

	actualTrip, err := tripsRepository.CreateTrip(context.Background(), tripCreate.OwnerID, tripCreate)

	require.ErrorIs(t, err, insertError)
	assert.Nil(t, actualTrip)
}
