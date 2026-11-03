package trips_transport_http

import (
	"context"
	"net/http"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_server "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/server"
)

//go:generate mockgen -source=transport.go -destination=mocks/trip_service.go -package mocks
type TripsService interface {
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
		tripPatch *domain.TripPatch,
	) (*domain.Trip, error)
}

type TripsHTTPHandler struct {
	tripsService TripsService
}

func NewTripsHTTPHandler(
	tripsService TripsService,
) *TripsHTTPHandler {
	return &TripsHTTPHandler{
		tripsService: tripsService,
	}
}

func (h *TripsHTTPHandler) Routes() []*core_server.Route {
	return []*core_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/trips/{" + OwnerIDParam + "}",
			Handler: h.CreateTrip,
		},
		{
			Method:  http.MethodGet,
			Path:    "/trips/{" + OwnerIDParam + "}" + "/" + "{" + TripIDParam + "}",
			Handler: h.GetTrip,
		},
		{
			Method:  http.MethodGet,
			Path:    "/trips/{" + OwnerIDParam + "}",
			Handler: h.GetTrips,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/trips/{" + OwnerIDParam + "}" + "/" + "{" + TripIDParam + "}",
			Handler: h.DeleteTrip,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/trips/{" + OwnerIDParam + "}" + "/" + "{" + TripIDParam + "}",
			Handler: h.PatchTrip,
		},
	}
}
