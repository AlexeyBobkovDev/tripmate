package transport_grpc_generate_route

import (
	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

func convertMealRequirementsDTOToDomain(mealRequirements []*route_generator_common_v1.MealRequirement) []domain.MealRequirement {
	mealRequirementTypes := map[route_generator_common_v1.MealType]domain.MealType{
		route_generator_common_v1.MealType_MEAL_TYPE_BREAKFAST: domain.MealBreakfast,
		route_generator_common_v1.MealType_MEAL_TYPE_LUNCH:     domain.MealLunch,
		route_generator_common_v1.MealType_MEAL_TYPE_DINNER:    domain.MealDinner,
	}
	convertedMealRequirements := make([]domain.MealRequirement, len(mealRequirements))
	for ind, mealRequirement := range mealRequirements {
		mealRequirementType, _ := mealRequirementTypes[mealRequirement.Type]
		convertedMealRequirements[ind] = domain.MealRequirement{
			Type: mealRequirementType,
			Window: domain.TimeWindow{
				Start: domain.LocalTime(mealRequirement.StartTime.AsTime().Format("15:04")),
				End:   domain.LocalTime(mealRequirement.EndTime.AsTime().Format("15:04")),
			},
			Duration: mealRequirement.Duration.AsDuration(),
			Budget:   convertBudgetDTOToDomain(mealRequirement.Budget),
		}
	}
	return convertedMealRequirements
}

func convertDailyScheduleDTOToDomain(dailySchedule *route_generator_common_v1.DailySchedule) domain.DailySchedule {
	return domain.DailySchedule{
		StartTime: domain.LocalTime(dailySchedule.StartTime.AsTime().Format("15:04")),
		EndTime:   domain.LocalTime(dailySchedule.EndTime.AsTime().Format("15:04")),
	}
}

func convertTripDTOToDomain(trip *route_generator_api_v1.GenerateRouteRequest) *domain.Trip {
	return &domain.Trip{
		Name:             trip.Name,
		PlanningArea:     trip.PlanningArea,
		StartCoordinates: convertStartCoordinatesDTOToDomain(trip.StartCoordinates),
		StartDate:        trip.StartDate.AsTime(),
		EndDate:          trip.EndDate.AsTime(),
		DailySchedule:    convertDailyScheduleDTOToDomain(trip.DailySchedule),
		MealRequirements: convertMealRequirementsDTOToDomain(trip.MealRequirements),
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
	mapInterestDTOToDomain := map[route_generator_api_v1.Interest]domain.Interest{
		route_generator_api_v1.Interest_INTEREST_HISTORY:      domain.InterestHistory,
		route_generator_api_v1.Interest_INTEREST_NATURE:       domain.InterestNature,
		route_generator_api_v1.Interest_INTEREST_FOOD:         domain.InterestFood,
		route_generator_api_v1.Interest_INTEREST_ARCHITECTURE: domain.InterestArchitecture,
	}
	convertedInterests := make([]domain.Interest, len(interests))
	for ind, interest := range interests {
		mappedInterest := mapInterestDTOToDomain[interest]
		convertedInterests[ind] = mappedInterest
	}

	return convertedInterests
}
