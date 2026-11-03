package trips_service

import (
	"context"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

//go:generate mockgen -source=service.go -destination=mocks/trip_repository.go -package mocks
type TripsRepository interface {
	CreateTrip(
		ctx context.Context,
		ownerID int,
		trip *domain.TripCreate,
	) (*domain.Trip, error)
	GetTrip(
		ctx context.Context,
		ownerID int,
		tripID int,
	) (*domain.Trip, error)
	GetTrips(
		ctx context.Context,
		ownerID int,
	) (*[]domain.Trip, error)
	DeleteTrip(
		ctx context.Context,
		ownerID int,
		tripID int,
	) error
	PatchTrip(
		ctx context.Context,
		ownerID int,
		tripID int,
		trip *domain.Trip,
	) (*domain.Trip, error)
}

type TripsService struct {
	tripsRepository     TripsRepository
	generateRouteClient route_generator_api_v1.RouteGeneratorServiceClient
}

func NewTripsService(
	tripsRepository TripsRepository,
	generateRouteClient route_generator_api_v1.RouteGeneratorServiceClient,
) *TripsService {
	return &TripsService{
		tripsRepository:     tripsRepository,
		generateRouteClient: generateRouteClient,
	}
}
