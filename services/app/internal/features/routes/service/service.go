package routes_service

import (
	"context"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	trips_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/service"
)

//go:generate mockgen -source=service.go -destination=mocks/route_repository.go -package mocks
type RoutesRepository interface {
	CreateRoute(
		ctx context.Context,
		ownerID int,
		trip *domain.Trip,
		generatedRoute *domain.RouteCreate,
	) (*domain.Route, error)
	GetRoute(
		ctx context.Context,
		tripID int,
		routeID int,
	) (*domain.Route, error)
	GetRoutes(
		ctx context.Context,
		tripID int,
	) ([]*domain.Route, error)
	DeleteRoute(
		ctx context.Context,
		tripID int,
		routeID int,
	) error
	// PatchRoute(
	// 	ctx context.Context,
	// 	ownerID int,
	// 	routeID int,
	// 	route *domain.Route,
	// ) (*domain.Route, error)
}

type RoutesService struct {
	routesRepository    RoutesRepository
	tripsService        *trips_service.TripsService
	generateRouteClient route_generator_api_v1.RouteGeneratorServiceClient
}

func NewRoutesService(
	routesRepository RoutesRepository,
	tripsService *trips_service.TripsService,
	generateRouteClient route_generator_api_v1.RouteGeneratorServiceClient,
) *RoutesService {
	return &RoutesService{
		routesRepository:    routesRepository,
		tripsService:        tripsService,
		generateRouteClient: generateRouteClient,
	}
}
