package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestValidateRouteTravelMode(t *testing.T) {
	tests := []struct {
		name       string
		travelMode TravelMode
		want       error
	}{
		{name: "walking", travelMode: TravelModeWalking},
		{name: "cycling", travelMode: TravelModeCycling},
		{name: "driving", travelMode: TravelModeDriving},
		{name: "transit", travelMode: TravelModeTransit},
		{name: "empty travel mode", travelMode: TravelMode(""), want: core_errors.ErrUnsupportedRouteTravelMode},
		{name: "unsupported travel mode", travelMode: "car", want: core_errors.ErrUnsupportedRouteTravelMode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteTravelMode(tt.travelMode)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteDistance(t *testing.T) {
	tests := []struct {
		name     string
		distance int
		want     error
	}{
		{name: "positive distance", distance: 1},
		{name: "typical leg distance", distance: 1234},
		{name: "zero distance", distance: 0, want: core_errors.ErrInvalidRouteDistance},
		{name: "negative distance", distance: -1, want: core_errors.ErrInvalidRouteDistance},
		{name: "negative typical distance", distance: -1234, want: core_errors.ErrInvalidRouteDistance},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteDistance(tt.distance)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     error
	}{
		{name: "positive duration", duration: time.Minute},
		{name: "typical leg duration", duration: 45 * time.Minute},
		{name: "zero duration", duration: 0, want: core_errors.ErrInvalidRouteDuration},
		{name: "negative duration", duration: -time.Minute, want: core_errors.ErrInvalidRouteDuration},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteDuration(tt.duration)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteStatus(t *testing.T) {
	tests := []struct {
		name       string
		tripStatus RouteStatus
		want       error
	}{
		{name: "generated", tripStatus: RouteStatusGenerated},
		{name: "edited", tripStatus: RouteStatusEdited},
		{name: "empty status", tripStatus: RouteStatus(""), want: core_errors.ErrUnsupportedRouteStatus},
		{name: "unsupported status", tripStatus: "planned", want: core_errors.ErrUnsupportedRouteStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteStatus(tt.tripStatus)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteDays(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)

	validStop := RouteStop{
		ID:            1,
		PlaceID:       1,
		Position:      1,
		Name:          "Stop",
		Description:   "Stop description",
		Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		ArrivalTime:   futureDate.Add(time.Hour),
		VisitDuration: 30 * time.Minute,
	}
	validStop.DepartureTime = futureDate.Add(2 * time.Hour)

	validLeg := RouteLeg{
		ID:         1,
		FromStopID: 1,
		ToStopID:   2,
		TravelMode: TravelModeWalking,
		Distance:   5,
		Duration:   10 * time.Minute,
	}

	validDay := RouteDay{
		ID:        1,
		Date:      futureDate,
		DayNumber: 1,
		Stops:     []RouteStop{validStop},
		Legs:      []RouteLeg{validLeg},
	}

	tests := []struct {
		name string
		days []RouteDay
		want error
	}{
		{name: "no days", days: nil, want: core_errors.ErrInvalidRouteDaysAmount},
		{name: "empty days", days: []RouteDay{}, want: core_errors.ErrInvalidRouteDaysAmount},
		{name: "single valid day", days: []RouteDay{validDay}},
		{name: "multiple valid days", days: []RouteDay{validDay, validDay}},
		{name: "first day has no stops", days: []RouteDay{{ID: 1, Date: futureDate, DayNumber: 1}}, want: core_errors.ErrInvalidRouteDayStopsAmount},
		{name: "first day has invalid stop", days: []RouteDay{{ID: 1, Date: futureDate, DayNumber: 1, Stops: []RouteStop{{ID: 1, PlaceID: 1, Position: 0}}, Legs: []RouteLeg{}}}, want: core_errors.ErrInvalidRouteStopPosition},
		{name: "second day has invalid leg", days: []RouteDay{validDay, {ID: 2, Date: futureDate, DayNumber: 2, Stops: []RouteStop{validStop}, Legs: []RouteLeg{{ID: 1, FromStopID: 1, ToStopID: 2, TravelMode: "unknown", Distance: 5, Duration: time.Minute}}}}, want: core_errors.ErrUnsupportedRouteTravelMode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteDays(tt.days)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteStopPosition(t *testing.T) {
	tests := []struct {
		name     string
		position int
		want     error
	}{
		{name: "first position", position: 1},
		{name: "typical position", position: 3},
		{name: "zero position", position: 0, want: core_errors.ErrInvalidRouteStopPosition},
		{name: "negative position", position: -1, want: core_errors.ErrInvalidRouteStopPosition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteStopPosition(tt.position)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteStopName(t *testing.T) {
	tests := []struct {
		name     string
		stopName string
		want     error
	}{
		{name: "valid name", stopName: "Eiffel Tower"},
		{name: "unicode name", stopName: "Московский Кремль"},
		{name: "single rune", stopName: "A"},
		{name: "empty name", stopName: "", want: core_errors.ErrInvalidRouteStopName},
		{name: "whitespace name is not trimmed and has runes", stopName: " ", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteStopName(tt.stopName)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteStopDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        error
	}{
		{name: "empty description is invalid", description: "", want: core_errors.ErrInvalidRouteStopDescription},
		{name: "short description", description: "Historic monument"},
		{name: "long description", description: "A very long description of the historic monument that visitors should read before arrival at the destination stop of the day"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteStopDescription(tt.description)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateArrivalTime(t *testing.T) {
	tests := []struct {
		name        string
		arrivalTime time.Time
		want        error
	}{
		{name: "future arrival time", arrivalTime: time.Now().Add(time.Hour)},
		{name: "zero arrival time", arrivalTime: time.Time{}, want: core_errors.ErrInvalidRouteArrivalTime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArrivalTime(tt.arrivalTime)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateDepartureTime(t *testing.T) {
	tests := []struct {
		name          string
		departureTime time.Time
		want          error
	}{
		{name: "future departure time", departureTime: time.Now().Add(2 * time.Hour)},
		{name: "zero departure time", departureTime: time.Time{}, want: core_errors.ErrInvalidRouteDepartureTime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDepartureTime(tt.departureTime)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteStopTimesOrder(t *testing.T) {
	tests := []struct {
		name          string
		arrivalTime   time.Time
		visitDuration time.Duration
		departureTime time.Time
		want          error
	}{
		{name: "valid stop", arrivalTime: time.Now().Add(time.Hour), visitDuration: 30 * time.Minute, departureTime: time.Now().Add(2 * time.Hour)},
		{name: "visit duration fits into departure time", arrivalTime: time.Now().Add(time.Hour), visitDuration: 30 * time.Minute, departureTime: time.Now().Add(90 * time.Minute)},
		{name: "visit duration after departure", arrivalTime: time.Now().Add(time.Hour), visitDuration: 2 * time.Hour, departureTime: time.Now().Add(90 * time.Minute), want: core_errors.ErrInvalidRouteStopDatesOrder},
		{name: "arrival after departure", arrivalTime: time.Now().Add(2 * time.Hour), visitDuration: 30 * time.Minute, departureTime: time.Now().Add(time.Hour), want: core_errors.ErrInvalidRouteStopDatesOrder},
		{name: "negative visit duration", arrivalTime: time.Now().Add(time.Hour), visitDuration: -time.Minute, departureTime: time.Now().Add(2 * time.Hour), want: core_errors.ErrInvalidRouteStopDatesOrder},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteStopTimesOrder(tt.arrivalTime, tt.visitDuration, tt.departureTime)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteDayDate(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want error
	}{
		{name: "future date", date: time.Now().Add(24 * time.Hour)},
		{name: "zero date", date: time.Time{}, want: core_errors.ErrInvalidRouteDayDate},
		{name: "past date", date: time.Now().Add(-24 * time.Hour), want: core_errors.ErrInvalidRouteDayDate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteDayDate(tt.date)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateRouteDayDayNumber(t *testing.T) {
	tests := []struct {
		name      string
		dayNumber int
		want      error
	}{
		{name: "first day", dayNumber: 1},
		{name: "typical day number", dayNumber: 3},
		{name: "zero day number", dayNumber: 0, want: core_errors.ErrInvalidRouteDayDayNumber},
		{name: "negative day number", dayNumber: -1, want: core_errors.ErrInvalidRouteDayDayNumber},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteDayDayNumber(tt.dayNumber)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
