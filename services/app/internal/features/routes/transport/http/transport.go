package routes_transport_http

import (
	"context"
	"net/http"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_server "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/server"
)

//go:generate mockgen -source=transport.go -destination=mocks/route_service.go -package mocks
type RoutesService interface {
	CreateRoute(
		ctx context.Context,
		ownerID int,
		tripID int,
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
}

type RoutesHTTPHandler struct {
	routesService RoutesService
}

func NewRoutesHTTPHandler(
	routesService RoutesService,
) *RoutesHTTPHandler {
	return &RoutesHTTPHandler{
		routesService: routesService,
	}
}

func (h *RoutesHTTPHandler) Routes() []*core_server.Route {
	return []*core_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/routes/{" + OwnerIDParam + "}" + "/" + "{" + TripIDParam + "}",
			Handler: h.CreateRoute,
		},
		{
			Method:  http.MethodGet,
			Path:    "/routes/{" + TripIDParam + "}" + "/" + "{" + RouteIDParam + "}",
			Handler: h.GetRoute,
		},
		{
			Method:  http.MethodGet,
			Path:    "/routes/{" + TripIDParam + "}",
			Handler: h.GetRoutes,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/routes/{" + TripIDParam + "}" + "/" + "{" + RouteIDParam + "}",
			Handler: h.DeleteRoute,
		},
	}
}
