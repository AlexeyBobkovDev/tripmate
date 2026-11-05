package routes_transport_http

import (
	"time"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

var timeNow = time.Now()

func testRouteResponse() domain.Route {
	return domain.Route{
		ID:          1,
		TripID:      1,
		TripVersion: 1,
		Version:     1,
		Status:      domain.RouteStatusGenerated,
		Days:        []domain.RouteDay{},
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}
}

func testExpectedRouteResponse() RouteResponse {
	expected := testRouteResponse()
	return RouteResponse{
		ID:          expected.ID,
		TripID:      expected.TripID,
		TripVersion: expected.TripVersion,
		Version:     expected.Version,
		Status:      expected.Status,
		Days:        expected.Days,
		CreatedAt:   expected.CreatedAt,
		UpdatedAt:   expected.UpdatedAt,
	}
}
