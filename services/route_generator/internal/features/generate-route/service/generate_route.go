package generate_route_service

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

func (s *GenerateRouteService) GenerateRoute(
	ctx context.Context,
	trip *domain.Trip,
) (*domain.Route, error) {
	if err := trip.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid trip data")
	}

	date, err := time.Parse("2006-01-02", "2027-03-10")
	if err != nil {
		return nil, err
	}

	route := &domain.Route{
		Status: domain.RouteStatusGenerated,
		Days: []domain.RouteDay{
			{
				ID:        1,
				Date:      date,
				DayNumber: 1,
				Stops: []domain.RouteStop{
					{
						ID:          1,
						PlaceID:     1,
						Position:    1,
						Name:        "Kremlin",
						Description: "Kremlin",
						Coordinates: domain.Coordinates{
							Latitude:  56.07981,
							Longitude: 46.192347,
						},
						ArrivalTime:   time.Now().Add(2 * time.Hour),
						DepartureTime: time.Now().Add(3 * time.Hour),
						VisitDuration: 3600 * time.Second,
					},
					{
						ID:          2,
						PlaceID:     2,
						Position:    2,
						Name:        "Kremlin",
						Description: "Kremlin",
						Coordinates: domain.Coordinates{
							Latitude:  56.07981,
							Longitude: 46.192347,
						},
						ArrivalTime:   time.Now().Add(5 * time.Hour),
						DepartureTime: time.Now().Add(6 * time.Hour),
						VisitDuration: 3600 * time.Second,
					},
				},
				Legs: []domain.RouteLeg{
					{
						ID:         1,
						FromStopID: 1,
						ToStopID:   2,
						TravelMode: domain.TravelModeWalking,
						Distance:   1000,
						Duration:   7200 * time.Second,
					},
				},
			},
		},
	}

	if err := route.Validate(); err != nil {
		return nil, status.Error(codes.Internal, "invalid returned route data")
	}

	return route, nil
}
