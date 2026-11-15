package routes_service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func TestTripDomainToDTO(t *testing.T) {
	startDate := time.Date(2027, 3, 10, 9, 0, 0, 0, time.UTC)
	endDate := time.Date(2027, 3, 11, 18, 0, 0, 0, time.UTC)
	trip := domain.NewTrip(
		1,
		1,
		1,
		"Trip to Moscow",
		"Moscow",
		domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		startDate,
		endDate,
		domain.DailySchedule{
			StartTime: domain.LocalTime("09:00"),
			EndTime:   domain.LocalTime("18:00"),
		},
		[]domain.MealRequirement{
			{
				Type:     domain.MealBreakfast,
				Window:   domain.TimeWindow{Start: domain.LocalTime("08:00"), End: domain.LocalTime("09:00")},
				Duration: 30 * time.Minute,
				Budget:   domain.Money{AmountMinorUnits: 50000, Currency: "RUB"},
			},
		},
		domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		[]domain.Interest{domain.InterestHistory},
		"Europe/Moscow",
		time.Time{},
		time.Time{},
		nil,
	)

	dto := tripDomainToDTO(trip)

	require.NotNil(t, dto)
	assert.Equal(t, trip.Name, dto.Name)
	assert.Equal(t, trip.PlanningArea, dto.PlanningArea)

	require.NotNil(t, dto.StartCoordinates)
	assert.Equal(t, trip.StartCoordinates.Latitude, dto.StartCoordinates.Latitude)
	assert.Equal(t, trip.StartCoordinates.Longitude, dto.StartCoordinates.Longitude)

	require.NotNil(t, dto.StartDate)
	assert.True(t, dto.StartDate.AsTime().Equal(trip.StartDate))
	require.NotNil(t, dto.EndDate)
	assert.True(t, dto.EndDate.AsTime().Equal(trip.EndDate))

	require.NotNil(t, dto.DailySchedule)
	parsedScheduleStart, err := time.Parse("15:04", "09:00")
	require.NoError(t, err)
	parsedScheduleEnd, err := time.Parse("15:04", "18:00")
	require.NoError(t, err)
	assert.True(t, dto.DailySchedule.StartTime.AsTime().Equal(parsedScheduleStart))
	assert.True(t, dto.DailySchedule.EndTime.AsTime().Equal(parsedScheduleEnd))

	require.Len(t, dto.MealRequirements, 1)
	require.NotNil(t, dto.MealRequirements[0])
	assert.Equal(t, route_generator_common_v1.MealType_MEAL_TYPE_BREAKFAST, dto.MealRequirements[0].Type)
	require.NotNil(t, dto.MealRequirements[0].StartTime)
	parsedMealStart, err := time.Parse("15:04", "08:00")
	require.NoError(t, err)
	parsedMealEnd, err := time.Parse("15:04", "09:00")
	require.NoError(t, err)
	assert.True(t, dto.MealRequirements[0].StartTime.AsTime().Equal(parsedMealStart))
	assert.True(t, dto.MealRequirements[0].EndTime.AsTime().Equal(parsedMealEnd))
	require.NotNil(t, dto.MealRequirements[0].Duration)
	assert.Equal(t, 30*time.Minute, dto.MealRequirements[0].Duration.AsDuration())
	require.NotNil(t, dto.MealRequirements[0].Budget)
	assert.Equal(t, trip.MealRequirements[0].Budget.AmountMinorUnits, dto.MealRequirements[0].Budget.AmountMinorUnits)
	assert.Equal(t, trip.MealRequirements[0].Budget.Currency, dto.MealRequirements[0].Budget.Currency)

	require.NotNil(t, dto.Budget)
	assert.Equal(t, trip.Budget.AmountMinorUnits, dto.Budget.AmountMinorUnits)
	assert.Equal(t, trip.Budget.Currency, dto.Budget.Currency)

	require.Len(t, dto.Interests, 1)
	assert.Equal(t, route_generator_api_v1.Interest_INTEREST_HISTORY, dto.Interests[0])

	assert.Equal(t, trip.TimeZone, dto.TimeZone)
}

func TestMealRequirementsDomainToDTO(t *testing.T) {
	mealRequirements := []domain.MealRequirement{
		{
			Type:     domain.MealBreakfast,
			Window:   domain.TimeWindow{Start: domain.LocalTime("08:00"), End: domain.LocalTime("09:00")},
			Duration: 45 * time.Minute,
			Budget:   domain.Money{AmountMinorUnits: 50000, Currency: "RUB"},
		},
		{
			Type:     domain.MealLunch,
			Window:   domain.TimeWindow{Start: domain.LocalTime("13:00"), End: domain.LocalTime("14:00")},
			Duration: time.Hour,
			Budget:   domain.Money{AmountMinorUnits: 75000, Currency: "RUB"},
		},
	}

	converted := mealRequirementsDomainToDTO(mealRequirements)

	require.Len(t, converted, 2)

	first := converted[0]
	require.NotNil(t, first)
	assert.Equal(t, route_generator_common_v1.MealType_MEAL_TYPE_BREAKFAST, first.Type)
	parsedMealStart, err := time.Parse("15:04", "08:00")
	require.NoError(t, err)
	parsedMealEnd, err := time.Parse("15:04", "09:00")
	require.NoError(t, err)
	require.NotNil(t, first.StartTime)
	assert.True(t, first.StartTime.AsTime().Equal(parsedMealStart))
	require.NotNil(t, first.EndTime)
	assert.True(t, first.EndTime.AsTime().Equal(parsedMealEnd))
	require.NotNil(t, first.Duration)
	assert.Equal(t, 45*time.Minute, first.Duration.AsDuration())
	require.NotNil(t, first.Budget)
	assert.Equal(t, int64(50000), first.Budget.AmountMinorUnits)
	assert.Equal(t, "RUB", first.Budget.Currency)

	second := converted[1]
	require.NotNil(t, second)
	assert.Equal(t, route_generator_common_v1.MealType_MEAL_TYPE_LUNCH, second.Type)
	require.NotNil(t, second.Duration)
	assert.Equal(t, time.Hour, second.Duration.AsDuration())
}

func TestDailyScheduleDomainToDTO(t *testing.T) {
	dailySchedule := domain.DailySchedule{
		StartTime: domain.LocalTime("09:00"),
		EndTime:   domain.LocalTime("18:00"),
	}

	dto := dailyScheduleDomainToDTO(dailySchedule)

	require.NotNil(t, dto)
	parsedStart, err := time.Parse("15:04", "09:00")
	require.NoError(t, err)
	parsedEnd, err := time.Parse("15:04", "18:00")
	require.NoError(t, err)
	assert.True(t, dto.StartTime.AsTime().Equal(parsedStart))
	assert.True(t, dto.EndTime.AsTime().Equal(parsedEnd))
}

func TestBudgetDomainToDTO(t *testing.T) {
	dto := budgetDomainToDTO(domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"})

	require.NotNil(t, dto)
	assert.Equal(t, int64(1500000), dto.AmountMinorUnits)
	assert.Equal(t, "RUB", dto.Currency)
}

func TestCoordinatesDomainToDTO(t *testing.T) {
	dto := coordinatesDomainToDTO(domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208})

	require.NotNil(t, dto)
	assert.Equal(t, 55.7539, dto.Latitude)
	assert.Equal(t, 37.6208, dto.Longitude)
}

func TestInterestsDomainToDTO(t *testing.T) {
	interests := interestsDomainToDTO([]domain.Interest{domain.InterestHistory, domain.InterestNature})

	require.Len(t, interests, 2)
	assert.Equal(t, route_generator_api_v1.Interest_INTEREST_HISTORY, interests[0])
	assert.Equal(t, route_generator_api_v1.Interest_INTEREST_NATURE, interests[1])
}

func TestStatusDomainFromDTO(t *testing.T) {
	assert.Equal(t, domain.RouteStatus("generated"), statusDomainFromDTO(route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED))
	assert.Equal(t, domain.RouteStatus("edited"), statusDomainFromDTO(route_generator_api_v1.RouteStatus_ROUTE_STATUS_EDITED))
}

func TestRouteDomainFromDTO(t *testing.T) {
	routeCreate := routeDomainFromDTO(
		&route_generator_api_v1.GenerateRouteResponse{
			Status: route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED,
			Days: []*route_generator_api_v1.RouteDay{
				{
					Id:        1,
					DayNumber: 1,
					Stops:     []*route_generator_api_v1.RouteStop{},
					Legs:      []*route_generator_api_v1.RouteLeg{},
				},
			},
		},
	)

	require.NotNil(t, routeCreate)
	assert.Equal(t, domain.RouteStatus("generated"), routeCreate.Status)
	require.Len(t, routeCreate.Days, 1)
}

func TestDaysDomainFromDTO(t *testing.T) {
	days := daysDomainFromDTO(
		[]*route_generator_api_v1.RouteDay{
			{
				Id:        1,
				DayNumber: 1,
				Stops: []*route_generator_api_v1.RouteStop{
					{
						Id:       1,
						PlaceId:  10,
						Position: 1,
						Name:     "Stop name",
						Coordinates: &route_generator_common_v1.Coordinates{
							Latitude:  55.7539,
							Longitude: 37.6208,
						},
					},
				},
				Legs: []*route_generator_api_v1.RouteLeg{
					{
						Id:         100,
						FromStopId: 1,
						ToStopId:   2,
						TravelMode: route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING,
						Distance:   100,
						Duration:   durationpb.New(800 * time.Second),
					},
				},
			},
		},
	)

	require.Len(t, days, 1)

	day := days[0]
	require.NotNil(t, day)
	assert.Equal(t, 1, day.ID)
	assert.Equal(t, 1, day.DayNumber)
	require.Len(t, day.Stops, 1)
	require.Len(t, day.Legs, 1)
}

func TestLegsDomainFromDTO(t *testing.T) {
	legs := legsDomainFromDTO(
		[]*route_generator_api_v1.RouteLeg{
			{
				Id:         100,
				FromStopId: 1,
				ToStopId:   2,
				TravelMode: route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING,
				Distance:   100,
				Duration:   durationpb.New(800 * time.Second),
			},
		},
	)

	require.Len(t, legs, 1)

	assert.Equal(t, 100, legs[0].ID)
	assert.Equal(t, 1, legs[0].FromStopID)
	assert.Equal(t, 2, legs[0].ToStopID)
	assert.Equal(t, domain.TravelModeWalking, legs[0].TravelMode)
	assert.Equal(t, 100, legs[0].Distance)
	assert.Equal(t, 800*time.Second, legs[0].Duration)
}

func TestTravelModeDTOToDomain(t *testing.T) {
	tests := []struct {
		name       string
		travelMode route_generator_api_v1.TravelMode
		expected   domain.TravelMode
	}{
		{name: "walking", travelMode: route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING, expected: domain.TravelModeWalking},
		{name: "cycling", travelMode: route_generator_api_v1.TravelMode_TRAVEL_MODE_CYCLING, expected: domain.TravelModeCycling},
		{name: "driving", travelMode: route_generator_api_v1.TravelMode_TRAVEL_MODE_DRIVING, expected: domain.TravelModeDriving},
		{name: "transit", travelMode: route_generator_api_v1.TravelMode_TRAVEL_MODE_TRANSIT, expected: domain.TravelModeTransit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, travelModeDTOToDomain(tt.travelMode))
		})
	}
}

func TestStopsDomainFromDTO(t *testing.T) {
	stops := stopsDomainFromDTO(
		[]*route_generator_api_v1.RouteStop{
			{
				Id:       1,
				PlaceId:  10,
				Position: 1,
				Name:     "Stop name",
				Coordinates: &route_generator_common_v1.Coordinates{
					Latitude:  55.7539,
					Longitude: 37.6208,
				},
				ArrivalTime:   timestamppb.New(time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)),
				DepartureTime: timestamppb.New(time.Date(2026, 10, 28, 14, 0, 0, 0, time.UTC)),
				VisitDuration: durationpb.New(7200 * time.Second),
			},
		},
	)

	require.Len(t, stops, 1)

	stop := stops[0]
	require.NotNil(t, stop)
	assert.Equal(t, 1, stop.ID)
	assert.Equal(t, 10, stop.PlaceID)
	assert.Equal(t, 1, stop.Position)
	assert.Equal(t, "Stop name", stop.Name)
	require.NotNil(t, stop.Coordinates)
	assert.Equal(t, 55.7539, stop.Coordinates.Latitude)
	assert.Equal(t, 37.6208, stop.Coordinates.Longitude)
	require.NotNil(t, stop.ArrivalTime)
	assert.True(t, stop.ArrivalTime.Equal(time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC)))
	require.NotNil(t, stop.DepartureTime)
	assert.True(t, stop.DepartureTime.Equal(time.Date(2026, 10, 28, 14, 0, 0, 0, time.UTC)))
	require.NotNil(t, stop.VisitDuration)
	assert.Equal(t, 7200*time.Second, stop.VisitDuration)
}
