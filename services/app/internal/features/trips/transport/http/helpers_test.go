package trips_transport_http

import (
	"fmt"
	"strings"
	"time"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

var timeNow = time.Now()

func testTripResponse() domain.Trip {
	startDate := timeNow.Add(24 * time.Hour)
	endDate := startDate
	return domain.Trip{
		ID:               1,
		OwnerID:          1,
		Version:          1,
		Name:             "Walk around Moscow",
		PlanningArea:     "Moscow",
		StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		StartDate:        startDate,
		EndDate:          endDate,
		Budget:           domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		Interests:        []domain.Interest{"history", "nature"},
		TimeZone:         "Europe/Moscow",
		CreatedAt:        timeNow,
		UpdatedAt:        timeNow,
		DeletedAt:        nil,
	}
}

func testExpectedTripResponse() TripResponse {
	expected := testTripResponse()
	return TripResponse{
		ID:               expected.ID,
		OwnerID:          expected.OwnerID,
		Version:          expected.Version,
		Name:             expected.Name,
		PlanningArea:     expected.PlanningArea,
		StartCoordinates: expected.StartCoordinates,
		StartDate:        expected.StartDate,
		EndDate:          expected.EndDate,
		Budget:           expected.Budget,
		Interests:        expected.Interests,
		TimeZone:         expected.TimeZone,
		CreatedAt:        expected.CreatedAt,
		UpdatedAt:        expected.UpdatedAt,
		DeletedAt:        expected.DeletedAt,
	}
}

func testTripRequestBodyWithDefaults() *strings.Reader {
	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate
	// MealRequirement.Duration is a time.Duration and therefore JSON-encoded as nanoseconds.
	body := fmt.Sprintf(
		`
		{
			"owner_id": 1,
			"name": "Walk around Moscow",
			"planning_area": "Moscow",
			"start_location": {"latitude": 55.7539, "longitude": 37.6208},
			"start_date": "%s",
			"end_date": "%s",
			"daily_schedule": {"start_time": "09:00", "end_time": "18:00"},
			"meal_requirements": [
				{
					"meal_type": "breakfast",
					"window": {"start": "08:00", "end": "09:00"},
					"duration_minutes": 60,
					"budget": {"amount_minor_units": 50000, "currency": "RUB"}
				}
			],
			"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
			"interests": ["history", "nature"],
			"time_zone": "Europe/Moscow"
		}`,
		startDate.Format(core_consts.Layout),
		endDate.Format(core_consts.Layout),
	)
	return strings.NewReader(body)
}
