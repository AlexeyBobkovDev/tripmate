package routes_transport_http

import (
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func routeDTOFromDomain(route *domain.Route) RouteResponse {
	return RouteResponse{
		ID:          route.ID,
		TripID:      route.TripID,
		TripVersion: route.TripVersion,
		Version:     route.Version,
		Status:      route.Status,
		Days:        route.Days,
		CreatedAt:   route.CreatedAt,
		UpdatedAt:   route.UpdatedAt,
	}
}

func routeDTOsFromDomains(routes []*domain.Route) []RouteResponse {
	responses := make([]RouteResponse, len(routes))
	for ind, route := range routes {
		responses[ind] = routeDTOFromDomain(route)
	}

	return responses
}
