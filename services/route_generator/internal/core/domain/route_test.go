package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func validRoute() *Route {
	createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	return &Route{
		ID:          1,
		TripID:      1,
		TripVersion: 1,
		Version:     1,
		Status:      RouteStatusGenerated,
		Days:        []RouteDay{validRouteDay(time.Now().Add(24 * time.Hour))},
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func validRouteDay(futureDate time.Time) RouteDay {
	return RouteDay{
		ID:        1,
		Date:      futureDate,
		DayNumber: 1,
		Stops: []RouteStop{
			{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop",
				Description:   "Stop description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   futureDate.Add(time.Hour),
				DepartureTime: futureDate.Add(2 * time.Hour),
				VisitDuration: 30 * time.Minute,
			},
		},
		Legs: []RouteLeg{
			{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   5,
				Duration:   10 * time.Minute,
			},
		},
	}
}

func TestRoute_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Route)
		wantErr error
	}{
		{name: "valid route"},
		{name: "invalid route ID", mutate: func(r *Route) { r.ID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid trip ID", mutate: func(r *Route) { r.TripID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid trip version", mutate: func(r *Route) { r.TripVersion = 0 }, wantErr: core_errors.ErrInvalidVersion},
		{name: "invalid version", mutate: func(r *Route) { r.Version = 0 }, wantErr: core_errors.ErrInvalidVersion},
		{name: "unsupported status", mutate: func(r *Route) { r.Status = "planned" }, wantErr: core_errors.ErrUnsupportedRouteStatus},
		{name: "no days", mutate: func(r *Route) { r.Days = nil }, wantErr: core_errors.ErrInvalidRouteDaysAmount},
		{name: "day without stops", mutate: func(r *Route) { r.Days = []RouteDay{{ID: 1, Date: time.Now().Add(24 * time.Hour), DayNumber: 1}} }, wantErr: core_errors.ErrInvalidRouteDayStopsAmount},
		{name: "day with invalid stop", mutate: func(r *Route) {
			r.Days = []RouteDay{{ID: 1, Date: time.Now().Add(24 * time.Hour), DayNumber: 1, Stops: []RouteStop{{ID: 1, PlaceID: 1, Position: 0}}, Legs: []RouteLeg{}}}
		}, wantErr: core_errors.ErrInvalidRouteStopPosition},
		{name: "day with invalid date", mutate: func(r *Route) {
			r.Days = []RouteDay{{ID: 1, Date: time.Time{}, DayNumber: 1, Stops: validRouteDay(time.Now().Add(24 * time.Hour)).Stops, Legs: validRouteDay(time.Now().Add(24 * time.Hour)).Legs}}
		}, wantErr: core_errors.ErrInvalidRouteDayDate},
		{name: "day with invalid day number", mutate: func(r *Route) {
			r.Days = []RouteDay{{ID: 1, Date: time.Now().Add(24 * time.Hour), DayNumber: 0, Stops: validRouteDay(time.Now().Add(24 * time.Hour)).Stops, Legs: validRouteDay(time.Now().Add(24 * time.Hour)).Legs}}
		}, wantErr: core_errors.ErrInvalidRouteDayDayNumber},
		{name: "zero created at", mutate: func(r *Route) { r.CreatedAt = time.Time{} }, wantErr: core_errors.ErrInvalidCreatedAt},
		{name: "future created at", mutate: func(r *Route) { r.CreatedAt = time.Now().Add(time.Hour) }, wantErr: core_errors.ErrInvalidCreatedAt},
		{name: "zero updated at", mutate: func(r *Route) { r.UpdatedAt = time.Time{} }, wantErr: core_errors.ErrInvalidUpdatedAt},
		{name: "future updated at", mutate: func(r *Route) { r.UpdatedAt = time.Now().Add(time.Hour) }, wantErr: core_errors.ErrInvalidUpdatedAt},
		{name: "updated before created", mutate: func(r *Route) { r.UpdatedAt = r.CreatedAt.Add(-time.Second) }, wantErr: core_errors.ErrInvalidUpdatedAt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := validRoute()
			if tt.mutate != nil {
				tt.mutate(route)
			}
			err := route.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
