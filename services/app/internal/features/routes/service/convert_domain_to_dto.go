package routes_service

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func tripDomainToDTO(trip *domain.Trip) *route_generator_api_v1.GenerateRouteRequest {
	return &route_generator_api_v1.GenerateRouteRequest{
		Name:             trip.Name,
		PlanningArea:     trip.PlanningArea,
		StartCoordinates: coordinatesDomainToDTO(trip.StartCoordinates),
		StartDate:        timestamppb.New(trip.StartDate),
		EndDate:          timestamppb.New(trip.EndDate),
		DailySchedule:    dailyScheduleDomainToDTO(trip.DailySchedule),
		MealRequirements: mealRequirementsDomainToDTO(trip.MealRequirements),
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

func dailyScheduleDomainToDTO(dailySchedule domain.DailySchedule) *route_generator_common_v1.DailySchedule {
	return &route_generator_common_v1.DailySchedule{
		StartTime: timestamppb.New(parseTimeString(string(dailySchedule.StartTime))),
		EndTime:   timestamppb.New(parseTimeString(string(dailySchedule.EndTime))),
	}
}

func mealRequirementsDomainToDTO(mealRequirements []domain.MealRequirement) []*route_generator_common_v1.MealRequirement {
	converted := make([]*route_generator_common_v1.MealRequirement, len(mealRequirements))
	for ind, mealReq := range mealRequirements {
		converted[ind] = &route_generator_common_v1.MealRequirement{
			Type:      mealTypeDomainToDTO(mealReq.Type),
			StartTime: timestamppb.New(parseTimeString(string(mealReq.Window.Start))),
			EndTime:   timestamppb.New(parseTimeString(string(mealReq.Window.End))),
			Duration:  durationpb.New(mealReq.Duration),
			Budget:    budgetDomainToDTO(mealReq.Budget),
		}
	}
	return converted
}

func mealTypeDomainToDTO(mealType domain.MealType) route_generator_common_v1.MealType {
	switch mealType {
	case domain.MealBreakfast:
		return route_generator_common_v1.MealType_MEAL_TYPE_BREAKFAST
	case domain.MealLunch:
		return route_generator_common_v1.MealType_MEAL_TYPE_LUNCH
	case domain.MealDinner:
		return route_generator_common_v1.MealType_MEAL_TYPE_DINNER
	default:
		return route_generator_common_v1.MealType_MEAL_TYPE_BREAKFAST
	}
}
