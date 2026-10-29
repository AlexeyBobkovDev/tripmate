package trips_service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/service/mocks"
)

func TestTripsService_CreateTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTripsRepository := mocks.NewMockTripsRepository(ctrl)
	tripsService := NewTripsService(mockTripsRepository)

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

	expectedTrip := &domain.Trip{
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

	mockTripsRepository.
		EXPECT().
		CreateTrip(context.Background(), tripCreate.OwnerID, &tripCreate).
		Return(expectedTrip, nil)

	returnedTrip, err := tripsService.CreateTrip(context.Background(), tripCreate.OwnerID, &tripCreate)
	require.NoError(t, err)
	assert.Equal(t, expectedTrip, returnedTrip)
}

func TestTripsService_CreateTrip_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTripsRepository := mocks.NewMockTripsRepository(ctrl)
	tripsService := NewTripsService(mockTripsRepository)

	tests := []struct {
		name    string
		ownerID int
		trip    domain.TripCreate
	}{
		{
			name:    "invalid owner ID",
			ownerID: 0,
			trip: domain.TripCreate{
				OwnerID: 0,
			},
		},
		{
			name:    "empty name",
			ownerID: 1,
			trip: domain.TripCreate{
				OwnerID: 1,
				Name:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTripsRepository.EXPECT().CreateTrip(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

			returnedTrip, err := tripsService.CreateTrip(context.Background(), tt.ownerID, &tt.trip)
			require.Nil(t, returnedTrip)
			assert.Error(t, err)
		})
	}
}

func TestTripsService_CreateTrip_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTripsRepository := mocks.NewMockTripsRepository(ctrl)
	tripsService := NewTripsService(mockTripsRepository)

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

	repoError := fmt.Errorf("database error")

	mockTripsRepository.
		EXPECT().
		CreateTrip(context.Background(), tripCreate.OwnerID, &tripCreate).
		Return(nil, repoError)

	returnedTrip, err := tripsService.CreateTrip(context.Background(), tripCreate.OwnerID, &tripCreate)
	require.Nil(t, returnedTrip)
	assert.ErrorIs(t, err, repoError)
}
