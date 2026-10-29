package trips_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestTripsRepository_CreateTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)

	tripCreate := domain.TripCreate{
		OwnerID:       1,
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
	}

	expectedTrip := domain.Trip{
		ID:            1,
		OwnerID:       tripCreate.OwnerID,
		Version:       1,
		Name:          tripCreate.Name,
		PlanningArea:  tripCreate.PlanningArea,
		StartLocation: tripCreate.StartLocation,
		StartDate:     tripCreate.StartDate,
		EndDate:       tripCreate.EndDate,
		Budget:        tripCreate.Budget,
		Interests:     tripCreate.Interests,
		TimeZone:      tripCreate.TimeZone,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     nil,
	}

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).DoAndReturn(func(destinations ...any) error {
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

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().QueryRow(gomock.Any(), createTripQuery, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRow)

	trip, err := tripsRepository.CreateTrip(context.Background(), tripCreate.OwnerID, &tripCreate)
	require.NoError(t, err)
	assert.Equal(t, &expectedTrip, trip)
}

func TestTripsRepository_CreateTrip_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)

	tripCreate := domain.TripCreate{
		OwnerID:       1,
		Name:          "Trip to Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 100000, Currency: "RUB"},
		Interests:     []domain.Interest{"history"},
		TimeZone:      "Europe/Moscow",
	}

	mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(core_postgres_pool_errors.ErrNoRows)
	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().QueryRow(gomock.Any(), createTripQuery, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockRow)

	trip, err := tripsRepository.CreateTrip(context.Background(), tripCreate.OwnerID, &tripCreate)
	require.Error(t, err)
	assert.Nil(t, trip)
}
