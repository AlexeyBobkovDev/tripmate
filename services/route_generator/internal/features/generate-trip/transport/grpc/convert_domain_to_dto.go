package transport_grpc_generate_route

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

func convertRouteDomainToDTO(route *domain.Route) *route_generator_api_v1.GenerateRouteResponse {
	routeStatusMap := map[domain.RouteStatus]route_generator_api_v1.RouteStatus{
		domain.RouteStatusGenerated: route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED,
		domain.RouteStatusEdited:    route_generator_api_v1.RouteStatus_ROUTE_STATUS_EDITED,
	}
	routeStatus, _ := routeStatusMap[route.Status]

	return &route_generator_api_v1.GenerateRouteResponse{
		Status: routeStatus,
		Days:   convertDaysFromDomainToDTO(route.Days),
	}
}

func convertDaysFromDomainToDTO(days []domain.RouteDay) []*route_generator_api_v1.RouteDay {
	convertedDays := make([]*route_generator_api_v1.RouteDay, len(days))
	for ind, day := range days {
		convertedDays[ind] = convertDayFromDomainToDTO(day)
	}
	return convertedDays
}

func convertDayFromDomainToDTO(day domain.RouteDay) *route_generator_api_v1.RouteDay {
	return &route_generator_api_v1.RouteDay{
		Id:        int64(day.ID),
		Date:      timestamppb.New(day.Date),
		DayNumber: int32(day.DayNumber),
		Stops:     convertStopsFromDomainToDTO(day.Stops),
		Legs:      convertLegsFromDomainToDTO(day.Legs),
	}
}

func convertStopsFromDomainToDTO(stops []domain.RouteStop) []*route_generator_api_v1.RouteStop {
	convertedStop := make([]*route_generator_api_v1.RouteStop, len(stops))
	for ind, stop := range stops {
		convertedStop[ind] = convertStopFromDomainToDTO(stop)
	}
	return convertedStop
}

func convertStopFromDomainToDTO(stop domain.RouteStop) *route_generator_api_v1.RouteStop {
	return &route_generator_api_v1.RouteStop{
		Id:            int64(stop.ID),
		PlaceId:       int64(stop.PlaceID),
		Position:      int32(stop.Position),
		Name:          stop.Name,
		Description:   stop.Description,
		Coordinates:   convertStartCoordinatesDomainToDTO(stop.Coordinates),
		ArrivalTime:   timestamppb.New(stop.ArrivalTime),
		DepartureTime: timestamppb.New(stop.DepartureTime),
		VisitDuration: durationpb.New(stop.VisitDuration),
	}
}

func convertLegsFromDomainToDTO(legs []domain.RouteLeg) []*route_generator_api_v1.RouteLeg {
	convertedLegs := make([]*route_generator_api_v1.RouteLeg, len(legs))
	for ind, leg := range legs {
		convertedLegs[ind] = convertLegFromDomainToDTO(leg)
	}
	return convertedLegs
}

func convertLegFromDomainToDTO(leg domain.RouteLeg) *route_generator_api_v1.RouteLeg {
	travelModeMap := map[domain.TravelMode]route_generator_api_v1.TravelMode{
		domain.TravelModeCycling: route_generator_api_v1.TravelMode_TRAVEL_MODE_CYCLING,
		domain.TravelModeWalking: route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING,
		domain.TravelModeDriving: route_generator_api_v1.TravelMode_TRAVEL_MODE_DRIVING,
		domain.TravelModeTransit: route_generator_api_v1.TravelMode_TRAVEL_MODE_TRANSIT,
	}
	travelMode, _ := travelModeMap[leg.TravelMode]

	return &route_generator_api_v1.RouteLeg{
		Id:         int64(leg.ID),
		FromStopId: int32(leg.FromStopID),
		ToStopId:   int32(leg.ToStopID),
		TravelMode: travelMode,
		Distance:   int32(leg.Distance),
		Duration:   durationpb.New(leg.Duration),
	}
}

func convertStartCoordinatesDomainToDTO(startCoordinates domain.Coordinates) *route_generator_common_v1.Coordinates {
	return &route_generator_common_v1.Coordinates{
		Latitude:  startCoordinates.Latitude,
		Longitude: startCoordinates.Longitude,
	}
}
