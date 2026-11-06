package routes_service

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *RoutesService) CreateRoute(
	ctx context.Context,
	ownerID int,
	tripID int,
) (*domain.Route, error) {
	if ownerID <= 0 {
		return nil, fmt.Errorf("ownerID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	if tripID <= 0 {
		return nil, fmt.Errorf("tripID must be greater than 0: %w", core_errors.ErrBadRequest)
	}

	trip, err := s.tripsService.GetTrip(ctx, ownerID, tripID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve trip: %w", err)
	}

	generatedRoute, err := s.generateRouteClient.GenerateRoute(ctx, tripDomainToDTO(trip))
	if err != nil {
		return nil, fmt.Errorf("failed to generate route: %w", err)
	}

	routeResponse, err := s.routesRepository.CreateRoute(ctx, ownerID, trip, routeDomainFromDTO(generatedRoute))
	if err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}
	return routeResponse, nil
}

func tripDomainToDTO(trip *domain.Trip) *route_generator_api_v1.GenerateRouteRequest {
	return &route_generator_api_v1.GenerateRouteRequest{
		Name:             trip.Name,
		PlanningArea:     trip.PlanningArea,
		StartCoordinates: coordinatesDomainToDTO(trip.StartCoordinates),
		StartDate:        timestamppb.New(trip.StartDate),
		EndDate:          timestamppb.New(trip.EndDate),
		Budget:           budgetDomainToDTO(trip.Budget),
		Interests:        interestsDomainToDTO(trip.Interests),
		TimeZone:         trip.TimeZone,
	}
}

func interestsDomainToDTO(interests []domain.Interest) []route_generator_api_v1.Interest {
	mapInterests := map[domain.Interest]route_generator_api_v1.Interest{
		domain.InterestHistory:      route_generator_api_v1.Interest_INTEREST_HISTORY,
		domain.InterestNature:       route_generator_api_v1.Interest_INTEREST_NATURE,
		domain.InterestFood:         route_generator_api_v1.Interest_INTEREST_FOOD,
		domain.InterestArchitecture: route_generator_api_v1.Interest_INTEREST_ARCHITECTURE,
	}
	convertedInterests := make([]route_generator_api_v1.Interest, len(interests))
	for ind, interest := range interests {
		mappedInterest, _ := mapInterests[interest]
		convertedInterests[ind] = mappedInterest
	}
	return convertedInterests
}

func budgetDomainToDTO(money domain.Money) *route_generator_common_v1.Money {
	return &route_generator_common_v1.Money{
		AmountMinorUnits: money.AmountMinorUnits,
		Currency:         money.Currency,
	}
}

func coordinatesDomainToDTO(coordinates domain.Coordinates) *route_generator_common_v1.Coordinates {
	return &route_generator_common_v1.Coordinates{
		Latitude:  coordinates.Latitude,
		Longitude: coordinates.Longitude,
	}
}

func routeDomainFromDTO(generatedRoute *route_generator_api_v1.GenerateRouteResponse) *domain.RouteCreate {
	return &domain.RouteCreate{
		Status: statusDomainFromDTO(generatedRoute.Status),
		Days:   daysDomainFromDTO(generatedRoute.Days),
	}
}

func statusDomainFromDTO(status route_generator_api_v1.RouteStatus) domain.RouteStatus {
	mapStatus := map[route_generator_api_v1.RouteStatus]domain.RouteStatus{
		route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED: domain.RouteStatusGenerated,
		route_generator_api_v1.RouteStatus_ROUTE_STATUS_EDITED:    domain.RouteStatusEdited,
	}

	mappedStatus, _ := mapStatus[status]

	return mappedStatus
}

func daysDomainFromDTO(days []*route_generator_api_v1.RouteDay) []domain.RouteDay {
	convertedDays := make([]domain.RouteDay, len(days))
	for ind, day := range days {
		convertedDays[ind] = domain.RouteDay{
			ID:        int(day.Id),
			Date:      day.Date.AsTime(),
			DayNumber: int(day.DayNumber),
			Stops:     stopsDomainFromDTO(day.Stops),
			Legs:      legsDomainFromDTO(day.Legs),
		}
	}
	return convertedDays
}

func legsDomainFromDTO(legs []*route_generator_api_v1.RouteLeg) []domain.RouteLeg {
	convertedDays := make([]domain.RouteLeg, len(legs))
	for ind, leg := range legs {
		convertedDays[ind] = domain.RouteLeg{
			ID:         int(leg.Id),
			FromStopID: int(leg.FromStopId),
			ToStopID:   int(leg.ToStopId),
			TravelMode: travelModeDTOToDomain(leg.TravelMode),
			Distance:   int(leg.Distance),
			Duration:   leg.Duration.AsDuration(),
		}
	}
	return convertedDays
}

func travelModeDTOToDomain(travelMode route_generator_api_v1.TravelMode) domain.TravelMode {
	mapTravelModes := map[route_generator_api_v1.TravelMode]domain.TravelMode{
		route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING: domain.TravelModeWalking,
		route_generator_api_v1.TravelMode_TRAVEL_MODE_CYCLING: domain.TravelModeCycling,
		route_generator_api_v1.TravelMode_TRAVEL_MODE_DRIVING: domain.TravelModeDriving,
		route_generator_api_v1.TravelMode_TRAVEL_MODE_TRANSIT: domain.TravelModeTransit,
	}
	mappedTravelMode, _ := mapTravelModes[travelMode]

	return mappedTravelMode
}

func stopsDomainFromDTO(stops []*route_generator_api_v1.RouteStop) []domain.RouteStop {
	convertedDays := make([]domain.RouteStop, len(stops))
	for ind, stop := range stops {
		convertedDays[ind] = domain.RouteStop{
			ID:          int(stop.Id),
			PlaceID:     int(stop.PlaceId),
			Position:    int(stop.Position),
			Name:        stop.Name,
			Description: stop.Description,
			Coordinates: domain.Coordinates{
				Latitude:  stop.Coordinates.Latitude,
				Longitude: stop.Coordinates.Longitude,
			},
			ArrivalTime:   stop.ArrivalTime.AsTime(),
			DepartureTime: stop.DepartureTime.AsTime(),
			VisitDuration: stop.VisitDuration.AsDuration(),
		}
	}
	return convertedDays
}
