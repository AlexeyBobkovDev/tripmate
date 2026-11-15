package trips_transport_http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_http_request "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/request"
)

func createTripRequestBody(mealRequirements string) string {
	return `{
		"name": "Walk around Moscow",
		"planning_area": "Moscow",
		"start_location": {"latitude": 55.7539, "longitude": 37.6208},
		"start_date": "2027-03-10",
		"end_date": "2027-03-10",
		"daily_schedule": {"start_time": "09:00", "end_time": "18:00"},
		"meal_requirements": ` + mealRequirements + `,
		"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
		"interests": ["history"],
		"time_zone": "Europe/Moscow"
	}`
}

func decodeCreateTripRequest(t *testing.T, body string) error {
	t.Helper()

	request := httptest.NewRequest(
		http.MethodPost,
		"/trips/1",
		strings.NewReader(body),
	)

	var decoded CreateTripRequest
	return core_http_request.DecodeAndValidate(request, &decoded)
}

// TestCreateTripRequest_ValidatesNestedMealRequirements_BUG documents that the
// `dive` validator is missing from CreateTripRequest.MealRequirements. Field
// tags inside MealRequirementRequestDTO (oneof, min, max) are never evaluated,
// so invalid meal data reaches the service layer.
func TestCreateTripRequest_ValidatesNestedMealRequirements_BUG(t *testing.T) {
	tests := []struct {
		name             string
		mealRequirements string
	}{
		{
			name:             "zero duration",
			mealRequirements: `[{"meal_type": "breakfast", "window": {"start": "08:00", "end": "09:00"}, "duration_minutes": 0, "budget": {"amount_minor_units": 50000, "currency": "RUB"}}]`,
		},
		{
			name:             "duration exceeds one day",
			mealRequirements: `[{"meal_type": "breakfast", "window": {"start": "08:00", "end": "09:00"}, "duration_minutes": 3600000000000, "budget": {"amount_minor_units": 50000, "currency": "RUB"}}]`,
		},
		{
			name:             "unsupported meal type",
			mealRequirements: `[{"meal_type": "brunch", "window": {"start": "08:00", "end": "09:00"}, "duration_minutes": 60, "budget": {"amount_minor_units": 50000, "currency": "RUB"}}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := decodeCreateTripRequest(
				t,
				createTripRequestBody(tt.mealRequirements),
			)

			require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
		})
	}
}

func TestCreateTripRequest_AcceptsValidMealRequirements(t *testing.T) {
	err := decodeCreateTripRequest(
		t,
		createTripRequestBody(`[{"meal_type": "breakfast", "window": {"start": "08:00", "end": "09:00"}, "duration_minutes": 60, "budget": {"amount_minor_units": 50000, "currency": "RUB"}}]`),
	)

	require.NoError(t, err)
}

func TestCreateTripRequest_RejectsInvalidTopLevelData(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "empty name",
			body: `{
				"name": "",
				"planning_area": "Moscow",
				"start_location": {"latitude": 55.7539, "longitude": 37.6208},
				"start_date": "2027-03-10",
				"end_date": "2027-03-10",
				"daily_schedule": {"start_time": "09:00", "end_time": "18:00"},
				"meal_requirements": [],
				"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
				"interests": ["history"],
				"time_zone": "Europe/Moscow"
			}`,
		},
		{
			name: "empty interests",
			body: `{
				"name": "Walk around Moscow",
				"planning_area": "Moscow",
				"start_location": {"latitude": 55.7539, "longitude": 37.6208},
				"start_date": "2027-03-10",
				"end_date": "2027-03-10",
				"daily_schedule": {"start_time": "09:00", "end_time": "18:00"},
				"meal_requirements": [],
				"budget": {"amount_minor_units": 1500000, "currency": "RUB"},
				"interests": [],
				"time_zone": "Europe/Moscow"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := decodeCreateTripRequest(t, tt.body)

			require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
		})
	}
}
