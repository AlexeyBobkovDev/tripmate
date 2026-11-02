package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func TestRouteDay_Validate(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name    string
		mutate  func(*RouteDay)
		wantErr error
	}{
		{name: "valid day"},
		{name: "invalid day ID", mutate: func(d *RouteDay) { d.ID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid date", mutate: func(d *RouteDay) { d.Date = time.Time{} }, wantErr: core_errors.ErrInvalidRouteDayDate},
		{name: "past date", mutate: func(d *RouteDay) { d.Date = time.Now().Add(-24 * time.Hour) }, wantErr: core_errors.ErrInvalidRouteDayDate},
		{name: "invalid day number", mutate: func(d *RouteDay) { d.DayNumber = 0 }, wantErr: core_errors.ErrInvalidRouteDayDayNumber},
		{name: "no stops", mutate: func(d *RouteDay) { d.Stops = nil }, wantErr: core_errors.ErrInvalidRouteDayStopsAmount},
		{name: "invalid stop", mutate: func(d *RouteDay) { d.Stops = []RouteStop{{ID: 1, PlaceID: 1, Position: 0}} }, wantErr: core_errors.ErrInvalidRouteStopPosition},
		{name: "no legs", mutate: func(d *RouteDay) { d.Legs = nil }, wantErr: core_errors.ErrInvalidRouteDayLegsAmount},
		{name: "invalid leg", mutate: func(d *RouteDay) {
			d.Legs = []RouteLeg{{ID: 1, FromStopID: 1, ToStopID: 2, TravelMode: "unknown", Distance: 5, Duration: time.Minute}}
		}, wantErr: core_errors.ErrUnsupportedRouteTravelMode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			day := validRouteDay(futureDate)
			if tt.mutate != nil {
				tt.mutate(&day)
			}
			err := day.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
