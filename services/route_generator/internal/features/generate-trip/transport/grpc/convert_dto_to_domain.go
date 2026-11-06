package transport_grpc_generate_route

import (
	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

func convertTripDTOToDomain(trip *route_generator_api_v1.GenerateRouteRequest) *domain.Trip {
	return &domain.Trip{
		Name:             trip.Name,
		PlanningArea:     trip.PlanningArea,
		StartCoordinates: convertStartCoordinatesDTOToDomain(trip.StartCoordinates),
		StartDate:        trip.StartDate.AsTime(),
		EndDate:          trip.EndDate.AsTime(),
		Budget:           convertBudgetDTOToDomain(trip.Budget),
		Interests:        convertInterestsDTOToDomain(trip.Interests),
		TimeZone:         trip.TimeZone,
	}
}

func convertStartCoordinatesDTOToDomain(startCoordinates *route_generator_common_v1.Coordinates) domain.Coordinates {
	return domain.Coordinates{
		Latitude:  startCoordinates.Latitude,
		Longitude: startCoordinates.Longitude,
	}
}

func convertBudgetDTOToDomain(budget *route_generator_common_v1.Money) domain.Money {
	return domain.Money{
		AmountMinorUnits: budget.AmountMinorUnits,
		Currency:         budget.Currency,
	}
}

func convertInterestsDTOToDomain(interests []route_generator_api_v1.Interest) []domain.Interest {
	convertedInterests := make([]domain.Interest, len(interests))
	for ind, interest := range interests {
		convertedInterests[ind] = domain.Interest(interest)
	}

	return convertedInterests
}
