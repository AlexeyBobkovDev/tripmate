package routes_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/service/mocks"
	trips_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/service"
	trips_mocks "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/service/mocks"
)

type fakeRouteGeneratorClient struct {
	generateRouteRequest  *route_generator_api_v1.GenerateRouteRequest
	generateRouteResponse *route_generator_api_v1.GenerateRouteResponse
	generateRouteError    error
}

func (c *fakeRouteGeneratorClient) GenerateRoute(
	ctx context.Context,
	in *route_generator_api_v1.GenerateRouteRequest,
	opts ...grpc.CallOption,
) (*route_generator_api_v1.GenerateRouteResponse, error) {
	c.generateRouteRequest = in
	return c.generateRouteResponse, c.generateRouteError
}

func TestRoutesService_CreateRoute_InvalidIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	mockTripsRepository := trips_mocks.NewMockTripsRepository(ctrl)
	mockGenerateRouteClient := &fakeRouteGeneratorClient{}
	tripsService := trips_service.NewTripsService(mockTripsRepository)
	routesService := NewRoutesService(mockRoutesRepository, tripsService, mockGenerateRouteClient)

	tests := []struct {
		name    string
		ownerID int
		tripID  int
	}{
		{
			name:    "zero ownerID",
			ownerID: 0,
			tripID:  1,
		},
		{
			name:    "negative ownerID",
			ownerID: -1,
			tripID:  1,
		},
		{
			name:    "zero tripID",
			ownerID: 1,
			tripID:  0,
		},
		{
			name:    "negative tripID",
			ownerID: 1,
			tripID:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoutesRepository.
				EXPECT().
				CreateRoute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			returnedRoute, err := routesService.CreateRoute(context.Background(), tt.ownerID, tt.tripID)
			require.ErrorIs(t, err, core_errors.ErrBadRequest)
			assert.Nil(t, returnedRoute)
		})
	}
}

func TestRoutesService_CreateRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	mockTripsRepository := trips_mocks.NewMockTripsRepository(ctrl)
	mockGenerateRouteClient := &fakeRouteGeneratorClient{
		generateRouteResponse: &route_generator_api_v1.GenerateRouteResponse{
			Status: route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED,
		},
	}
	tripsService := trips_service.NewTripsService(mockTripsRepository)
	routesService := NewRoutesService(mockRoutesRepository, tripsService, mockGenerateRouteClient)

	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate

	trip := &domain.Trip{
		ID:           1,
		OwnerID:      1,
		Version:      1,
		Name:         "Trip to Moscow",
		PlanningArea: "Moscow",
		StartCoordinates: domain.Coordinates{
			Latitude:  55.7539,
			Longitude: 37.6208,
		},
		StartDate: startDate,
		EndDate:   endDate,
		Budget: domain.Money{
			AmountMinorUnits: 1500000,
			Currency:         "RUB",
		},
		Interests: []domain.Interest{"history"},
		TimeZone:  "Europe/Moscow",
	}

	expectedRoute := &domain.Route{
		ID:          1,
		TripID:      1,
		TripVersion: 1,
		Version:     1,
		Status:      domain.RouteStatusGenerated,
		Days:        []domain.RouteDay{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockTripsRepository.
		EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(trip, nil)

	mockRoutesRepository.
		EXPECT().
		CreateRoute(gomock.Any(), 1, trip, gomock.Any()).
		Return(expectedRoute, nil)

	returnedRoute, err := routesService.CreateRoute(context.Background(), 1, 1)
	require.NoError(t, err)
	assert.Equal(t, expectedRoute, returnedRoute)

	require.NotNil(t, mockGenerateRouteClient.generateRouteRequest)
	assert.Equal(t, trip.Name, mockGenerateRouteClient.generateRouteRequest.Name)
	assert.Equal(t, trip.PlanningArea, mockGenerateRouteClient.generateRouteRequest.PlanningArea)
	assert.Equal(t, trip.StartCoordinates.Latitude, mockGenerateRouteClient.generateRouteRequest.StartCoordinates.Latitude)
	assert.Equal(t, trip.StartCoordinates.Longitude, mockGenerateRouteClient.generateRouteRequest.StartCoordinates.Longitude)
	assert.True(t, mockGenerateRouteClient.generateRouteRequest.StartDate.AsTime().Equal(trip.StartDate))
	assert.True(t, mockGenerateRouteClient.generateRouteRequest.EndDate.AsTime().Equal(trip.EndDate))
	assert.Equal(t, trip.Budget.AmountMinorUnits, mockGenerateRouteClient.generateRouteRequest.Budget.AmountMinorUnits)
	assert.Equal(t, trip.Budget.Currency, mockGenerateRouteClient.generateRouteRequest.Budget.Currency)
	require.Len(t, mockGenerateRouteClient.generateRouteRequest.Interests, len(trip.Interests))
	for ind := range trip.Interests {
		assert.Equal(t, route_generator_api_v1.Interest(ind+1), mockGenerateRouteClient.generateRouteRequest.Interests[ind])
	}
	assert.Equal(t, trip.TimeZone, mockGenerateRouteClient.generateRouteRequest.TimeZone)
}

func TestRoutesService_CreateRoute_GetTripError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	mockTripsRepository := trips_mocks.NewMockTripsRepository(ctrl)
	mockGenerateRouteClient := &fakeRouteGeneratorClient{}
	tripsService := trips_service.NewTripsService(mockTripsRepository)
	routesService := NewRoutesService(mockRoutesRepository, tripsService, mockGenerateRouteClient)

	getTripError := errors.New("failed to retrieve trip")
	mockTripsRepository.
		EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(nil, getTripError)

	returnedRoute, err := routesService.CreateRoute(context.Background(), 1, 1)
	require.Error(t, err)
	assert.Nil(t, returnedRoute)
}

func TestRoutesService_CreateRoute_GenerateRouteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	mockTripsRepository := trips_mocks.NewMockTripsRepository(ctrl)
	mockGenerateRouteClient := &fakeRouteGeneratorClient{generateRouteError: errors.New("failed to generate route")}
	tripsService := trips_service.NewTripsService(mockTripsRepository)
	routesService := NewRoutesService(mockRoutesRepository, tripsService, mockGenerateRouteClient)

	mockTripsRepository.
		EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(&domain.Trip{ID: 1, OwnerID: 1, Version: 1}, nil)

	returnedRoute, err := routesService.CreateRoute(context.Background(), 1, 1)
	require.Error(t, err)
	assert.Nil(t, returnedRoute)
}

func TestRoutesService_CreateRoute_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRoutesRepository := mocks.NewMockRoutesRepository(ctrl)
	mockTripsRepository := trips_mocks.NewMockTripsRepository(ctrl)
	mockGenerateRouteClient := &fakeRouteGeneratorClient{
		generateRouteResponse: &route_generator_api_v1.GenerateRouteResponse{
			Status: route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED,
		},
	}
	tripsService := trips_service.NewTripsService(mockTripsRepository)
	routesService := NewRoutesService(mockRoutesRepository, tripsService, mockGenerateRouteClient)

	mockTripsRepository.
		EXPECT().
		GetTrip(context.Background(), 1, 1).
		Return(&domain.Trip{ID: 1, OwnerID: 1, Version: 1}, nil)

	repoError := errors.New("database error")
	mockRoutesRepository.
		EXPECT().
		CreateRoute(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, repoError)

	returnedRoute, err := routesService.CreateRoute(context.Background(), 1, 1)
	require.ErrorIs(t, err, repoError)
	assert.Nil(t, returnedRoute)
}
