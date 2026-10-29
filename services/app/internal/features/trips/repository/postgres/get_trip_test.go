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
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestTripsRepository_GetTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)
	createdAt := time.Now()
	updatedAt := time.Now()

	expectedTrip := domain.Trip{
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

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		QueryRow(gomock.Any(), getTripQuery, ownerID, tripID).
		Return(mockRow)

	mockRow.EXPECT().
		Scan(
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
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedTrip.ID
			*destinations[1].(*int) = expectedTrip.OwnerID
			*destinations[2].(*int) = expectedTrip.Version
			*destinations[3].(*string) = expectedTrip.Name
			*destinations[4].(*string) = expectedTrip.PlanningArea
			*destinations[5].(*string) = expectedTrip.StartLocation.Name
			*destinations[6].(*float64) = expectedTrip.StartLocation.Latitude
			*destinations[7].(*float64) = expectedTrip.StartLocation.Longitude
			*destinations[8].(*time.Time) = expectedTrip.StartDate
			*destinations[9].(*time.Time) = expectedTrip.EndDate
			*destinations[10].(*int64) = expectedTrip.Budget.AmountMinorUnits
			*destinations[11].(*string) = expectedTrip.Budget.Currency
			*destinations[12].(*[]domain.Interest) = expectedTrip.Interests
			*destinations[13].(*string) = expectedTrip.TimeZone
			*destinations[14].(*time.Time) = expectedTrip.CreatedAt
			*destinations[15].(*time.Time) = expectedTrip.UpdatedAt
			*destinations[16].(**time.Time) = expectedTrip.DeletedAt

			return nil
		})

	actualTrip, err := tripsRepository.GetTrip(context.Background(), ownerID, tripID)

	require.NoError(t, err)
	require.NotNil(t, actualTrip)
	assert.Equal(t, expectedTrip, *actualTrip)
}

func TestTripsRepository_GetTrip_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
	scanError := errors.New("query row error")

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		QueryRow(gomock.Any(), getTripQuery, ownerID, tripID).
		Return(mockRow)

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		Return(scanError)

	actualTrip, err := tripsRepository.GetTrip(context.Background(), ownerID, tripID)

	require.ErrorIs(t, err, scanError)
	assert.Nil(t, actualTrip)
}

func TestTripsRepository_GetTrip_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		QueryRow(gomock.Any(), getTripQuery, ownerID, tripID).
		Return(mockRow)

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		Return(errors.Join(pgx.ErrNoRows, core_errors.ErrNotFound))

	actualTrip, err := tripsRepository.GetTrip(context.Background(), ownerID, tripID)

	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.Nil(t, actualTrip)
}
