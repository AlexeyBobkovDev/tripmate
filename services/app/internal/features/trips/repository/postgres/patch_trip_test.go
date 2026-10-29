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

func TestTripsRepository_PatchTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)
	createdAt := time.Now()
	updatedAt := time.Now()

	patchTripRequest := domain.Trip{
		ID:            tripID,
		OwnerID:       ownerID,
		Version:       1,
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		DeletedAt:     nil,
	}

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).DoAndReturn(func(scanDestination ...any) error {
		*scanDestination[0].(*int) = patchTripRequest.ID
		*scanDestination[1].(*int) = patchTripRequest.OwnerID
		*scanDestination[2].(*int) = patchTripRequest.Version + 1
		*scanDestination[3].(*string) = patchTripRequest.Name
		*scanDestination[4].(*string) = patchTripRequest.PlanningArea
		*scanDestination[5].(*string) = patchTripRequest.StartLocation.Name
		*scanDestination[6].(*float64) = patchTripRequest.StartLocation.Latitude
		*scanDestination[7].(*float64) = patchTripRequest.StartLocation.Longitude
		*scanDestination[8].(*time.Time) = patchTripRequest.StartDate
		*scanDestination[9].(*time.Time) = patchTripRequest.EndDate
		*scanDestination[10].(*int64) = patchTripRequest.Budget.AmountMinorUnits
		*scanDestination[11].(*string) = patchTripRequest.Budget.Currency
		*scanDestination[12].(*[]domain.Interest) = patchTripRequest.Interests
		*scanDestination[13].(*string) = patchTripRequest.TimeZone
		*scanDestination[14].(*time.Time) = patchTripRequest.CreatedAt
		*scanDestination[15].(*time.Time) = patchTripRequest.UpdatedAt
		*scanDestination[16].(**time.Time) = patchTripRequest.DeletedAt

		return nil
	})

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			patchTripQuery,
			ownerID,
			tripID,
			patchTripRequest.Version,
			patchTripRequest.Name,
			patchTripRequest.PlanningArea,
			patchTripRequest.StartLocation.Name,
			patchTripRequest.StartLocation.Latitude,
			patchTripRequest.StartLocation.Longitude,
			patchTripRequest.StartDate,
			patchTripRequest.EndDate,
			patchTripRequest.Budget.AmountMinorUnits,
			patchTripRequest.Budget.Currency,
			patchTripRequest.Interests,
			patchTripRequest.TimeZone,
		).
		Return(mockRow)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), ownerID, tripID, &patchTripRequest)
	require.NoError(t, err)

	assert.Equal(t, tripID, actualTrip.ID)
	assert.Equal(t, ownerID, actualTrip.OwnerID)
	assert.Equal(t, patchTripRequest.Version+1, actualTrip.Version)
	assert.Equal(t, patchTripRequest.Name, actualTrip.Name)
}

func TestTripsRepository_PatchTrip_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
	scanError := errors.New("scan error")
	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)
	patchTripRequest := domain.Trip{
		ID:            tripID,
		OwnerID:       ownerID,
		Version:       1,
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
	}

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).Return(scanError)

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), patchTripQuery, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockRow)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), ownerID, tripID, &patchTripRequest)
	require.ErrorIs(t, err, scanError)
	require.Nil(t, actualTrip)
}

func TestTripsRepository_PatchTrip_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockTripExistsRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)
	patchTripRequest := domain.Trip{
		ID:            tripID,
		OwnerID:       ownerID,
		Version:       2,
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
	}

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).Return(core_postgres_pool_errors.ErrNoRows)

	mockTripExistsRow.EXPECT().
		Scan(
			gomock.Any(),
		).
		DoAndReturn(func(dest ...any) error {
			*dest[0].(*bool) = false
			return nil
		})

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), patchTripQuery, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockRow)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			tripExistsQuery,
			ownerID,
			tripID,
		).
		Return(mockTripExistsRow)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), ownerID, tripID, &patchTripRequest)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.False(t, errors.Is(err, core_errors.ErrConflict))
	require.Nil(t, actualTrip)
}

func TestTripsRepository_PatchTrip_ErrConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockTripExistsRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)
	patchTripRequest := domain.Trip{
		ID:            tripID,
		OwnerID:       ownerID,
		Version:       2,
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
	}

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).Return(core_postgres_pool_errors.ErrNoRows)

	mockTripExistsRow.EXPECT().
		Scan(
			gomock.Any(),
		).
		DoAndReturn(func(dest ...any) error {
			*dest[0].(*bool) = true
			return nil
		})

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(gomock.Any(), patchTripQuery, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockRow)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			tripExistsQuery,
			ownerID,
			tripID,
		).
		Return(mockTripExistsRow)

	actualTrip, err := tripsRepository.PatchTrip(context.Background(), ownerID, tripID, &patchTripRequest)
	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.ErrorIs(t, err, core_postgres_pool_errors.ErrNoRows)
	require.Nil(t, actualTrip)
}
