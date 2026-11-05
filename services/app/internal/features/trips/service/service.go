package trips_service

import (
	"context"

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
	tripsRepository TripsRepository
}

func NewTripsService(
	tripsRepository TripsRepository,
) *TripsService {
	return &TripsService{
		tripsRepository: tripsRepository,
	}
}
