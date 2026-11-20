package transport_grpc_generate_route

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

func TestConvertMealRequirementsDTOToDomain(t *testing.T) {
	t.Run("nil meal requirements", func(t *testing.T) {
		converted := convertMealRequirementsDTOToDomain(nil)
		require.NotNil(t, converted)
		assert.Empty(t, converted)
	})

	t.Run("single meal requirement", func(t *testing.T) {
		converted := convertMealRequirementsDTOToDomain([]*route_generator_common_v1.MealRequirement{
			{
				Type:      route_generator_common_v1.MealType_MEAL_TYPE_BREAKFAST,
				StartTime: timestamppb.New(time.Date(2030, time.May, 10, 8, 0, 0, 0, time.UTC)),
				EndTime:   timestamppb.New(time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC)),
				Duration:  durationpb.New(30 * time.Minute),
				Budget: &route_generator_common_v1.Money{
					AmountMinorUnits: 1000,
					Currency:         "RUB",
				},
			},
		})

		require.Len(t, converted, 1)
		convertedMeal := converted[0]
		assert.Equal(t, domain.MealBreakfast, convertedMeal.Type)
		assert.Equal(t, domain.LocalTime("08:00"), convertedMeal.Window.Start)
		assert.Equal(t, domain.LocalTime("09:00"), convertedMeal.Window.End)
		assert.Equal(t, 30*time.Minute, convertedMeal.Duration)
		assert.Equal(t, int64(1000), convertedMeal.Budget.AmountMinorUnits)
		assert.Equal(t, "RUB", convertedMeal.Budget.Currency)
	})

	t.Run("multiple meal requirements", func(t *testing.T) {
		converted := convertMealRequirementsDTOToDomain([]*route_generator_common_v1.MealRequirement{
			{
				Type:      route_generator_common_v1.MealType_MEAL_TYPE_LUNCH,
				StartTime: timestamppb.New(time.Date(2030, time.May, 10, 13, 0, 0, 0, time.UTC)),
				EndTime:   timestamppb.New(time.Date(2030, time.May, 10, 14, 0, 0, 0, time.UTC)),
				Duration:  durationpb.New(45 * time.Minute),
				Budget: &route_generator_common_v1.Money{
					AmountMinorUnits: 2000,
					Currency:         "RUB",
				},
			},
			{
				Type:      route_generator_common_v1.MealType_MEAL_TYPE_DINNER,
				StartTime: timestamppb.New(time.Date(2030, time.May, 10, 19, 0, 0, 0, time.UTC)),
				EndTime:   timestamppb.New(time.Date(2030, time.May, 10, 20, 0, 0, 0, time.UTC)),
				Duration:  durationpb.New(60 * time.Minute),
				Budget: &route_generator_common_v1.Money{
					AmountMinorUnits: 3000,
					Currency:         "RUB",
				},
			},
		})

		require.Len(t, converted, 2)
		assert.Equal(t, domain.MealLunch, converted[0].Type)
		assert.Equal(t, domain.LocalTime("13:00"), converted[0].Window.Start)
		assert.Equal(t, domain.LocalTime("14:00"), converted[0].Window.End)
		assert.Equal(t, 45*time.Minute, converted[0].Duration)
		assert.Equal(t, int64(2000), converted[0].Budget.AmountMinorUnits)

		assert.Equal(t, domain.MealDinner, converted[1].Type)
		assert.Equal(t, domain.LocalTime("19:00"), converted[1].Window.Start)
		assert.Equal(t, domain.LocalTime("20:00"), converted[1].Window.End)
		assert.Equal(t, 60*time.Minute, converted[1].Duration)
		assert.Equal(t, int64(3000), converted[1].Budget.AmountMinorUnits)
	})
}

func TestConvertDailyScheduleDTOToDomain(t *testing.T) {
	schedule := convertDailyScheduleDTOToDomain(&route_generator_common_v1.DailySchedule{
		StartTime: timestamppb.New(time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC)),
		EndTime:   timestamppb.New(time.Date(2030, time.May, 10, 18, 0, 0, 0, time.UTC)),
	})

	require.NotNil(t, schedule)
	assert.Equal(t, domain.LocalTime("09:00"), schedule.StartTime)
	assert.Equal(t, domain.LocalTime("18:00"), schedule.EndTime)
}

func TestConvertTripDTOToDomain(t *testing.T) {
	trip := convertTripDTOToDomain(&route_generator_api_v1.GenerateRouteRequest{
		Name:             "Trip to Moscow",
		PlanningArea:     "Moscow",
		StartCoordinates: &route_generator_common_v1.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		StartDate:        timestamppb.New(time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC)),
		EndDate:          timestamppb.New(time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC)),
		DailySchedule: &route_generator_common_v1.DailySchedule{
			StartTime: timestamppb.New(time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC)),
			EndTime:   timestamppb.New(time.Date(2030, time.May, 10, 18, 0, 0, 0, time.UTC)),
		},
		MealRequirements: []*route_generator_common_v1.MealRequirement{},
		Budget: &route_generator_common_v1.Money{
			AmountMinorUnits: 1500000,
			Currency:         "RUB",
		},
		Interests: []route_generator_api_v1.Interest{
			route_generator_api_v1.Interest_INTEREST_HISTORY,
			route_generator_api_v1.Interest_INTEREST_FOOD,
		},
		TimeZone: "Europe/Moscow",
	})

	require.NotNil(t, trip)
	assert.Equal(t, "Trip to Moscow", trip.Name)
	assert.Equal(t, "Moscow", trip.PlanningArea)
	require.NotNil(t, trip.StartCoordinates)
	assert.Equal(t, 55.7539, trip.StartCoordinates.Latitude)
	assert.Equal(t, 37.6208, trip.StartCoordinates.Longitude)
	require.NotNil(t, trip.StartDate)
	assert.Equal(t, 2030, trip.StartDate.Year())
	require.NotNil(t, trip.EndDate)
	assert.Equal(t, 2030, trip.EndDate.Year())
	require.NotNil(t, trip.DailySchedule)
	assert.Equal(t, domain.LocalTime("09:00"), trip.DailySchedule.StartTime)
	assert.Equal(t, domain.LocalTime("18:00"), trip.DailySchedule.EndTime)
	require.Len(t, trip.MealRequirements, 0)
	require.NotNil(t, trip.Budget)
	assert.Equal(t, int64(1500000), trip.Budget.AmountMinorUnits)
	assert.Equal(t, "RUB", trip.Budget.Currency)
	require.Len(t, trip.Interests, 2)
	assert.Equal(t, domain.InterestHistory, trip.Interests[0])
	assert.Equal(t, domain.InterestFood, trip.Interests[1])
	assert.Equal(t, "Europe/Moscow", trip.TimeZone)
}

func TestConvertStartCoordinatesDTOToDomain(t *testing.T) {
	coordinates := convertStartCoordinatesDTOToDomain(&route_generator_common_v1.Coordinates{
		Latitude:  55.7539,
		Longitude: 37.6208,
	})

	require.NotNil(t, coordinates)
	assert.Equal(t, 55.7539, coordinates.Latitude)
	assert.Equal(t, 37.6208, coordinates.Longitude)
}

func TestConvertBudgetDTOToDomain(t *testing.T) {
	money := convertBudgetDTOToDomain(&route_generator_common_v1.Money{
		AmountMinorUnits: 1500000,
		Currency:         "RUB",
	})

	require.NotNil(t, money)
	assert.Equal(t, int64(1500000), money.AmountMinorUnits)
	assert.Equal(t, "RUB", money.Currency)
}

func TestConvertInterestsDTOToDomain(t *testing.T) {
	interests := convertInterestsDTOToDomain([]route_generator_api_v1.Interest{
		route_generator_api_v1.Interest_INTEREST_HISTORY,
		route_generator_api_v1.Interest_INTEREST_FOOD,
	})

	require.Len(t, interests, 2)
	assert.Equal(t, domain.InterestHistory, interests[0])
	assert.Equal(t, domain.InterestFood, interests[1])
}
