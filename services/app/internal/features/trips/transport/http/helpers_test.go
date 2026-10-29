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
	endDate := startDate.Add(48 * time.Hour)
	return domain.Trip{
		ID:            1,
		OwnerID:       1,
		Version:       1,
		Name:          "Walk around Moscow",
		PlanningArea:  "Moscow",
		StartLocation: domain.Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		Interests:     []domain.Interest{"history", "nature"},
		TimeZone:      "Europe/Moscow",
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
		DeletedAt:     nil,
	}
}

func testExpectedTripResponse() TripResponse {
	expected := testTripResponse()
	return TripResponse{
		ID:            expected.ID,
		OwnerID:       expected.OwnerID,
		Version:       expected.Version,
		Name:          expected.Name,
		PlanningArea:  expected.PlanningArea,
		StartLocation: expected.StartLocation,
		StartDate:     expected.StartDate,
		EndDate:       expected.EndDate,
		Budget:        expected.Budget,
		Interests:     expected.Interests,
		TimeZone:      expected.TimeZone,
		CreatedAt:     expected.CreatedAt,
		UpdatedAt:     expected.UpdatedAt,
		DeletedAt:     expected.DeletedAt,
	}
}

func testTripRequestBodyWithDefaults() *strings.Reader {
	startDate := time.Now().Add(24 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)
	body := fmt.Sprintf(`
		{
			"owner_id": 1,
			"name": "Walk around Moscow",
			"planning_area": "Moscow",
			"start_location": {"name": "Red Square", "latitude": 55.7539, "longitude": 37.6208},
			"start_date": "%s",
			"end_date": "%s",
			"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
			"interests": ["history", "nature"],
			"time_zone": "Europe/Moscow"
		}`,
		startDate.Format(core_consts.Layout),
		endDate.Format(core_consts.Layout),
	)
	return strings.NewReader(body)
}
