package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestRouteCreateValidate(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)

	validRouteCreate := RouteCreate{
		Status: RouteStatusGenerated,
		Days:   []RouteDay{validRouteDay(futureDate)},
	}

	tests := []struct {
		name    string
		mutate  func(*RouteCreate)
		wantErr error
	}{
		{name: "valid route create", mutate: func(r *RouteCreate) {}, wantErr: nil},
		{name: "valid edited route create", mutate: func(r *RouteCreate) { r.Status = RouteStatusEdited }, wantErr: nil},
		{
			name: "multiple valid days",
			mutate: func(r *RouteCreate) {
				r.Days = []RouteDay{validRouteDay(futureDate), validRouteDay(futureDate)}
			},
			wantErr: nil,
		},
		{
			name:    "empty status",
			mutate:  func(r *RouteCreate) { r.Status = RouteStatus("") },
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name:    "unsupported status",
			mutate:  func(r *RouteCreate) { r.Status = "planned" },
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name: "status validation error takes precedence over days validation error",
			mutate: func(r *RouteCreate) {
				r.Status = "planned"
				r.Days = nil
			},
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name:    "no days",
			mutate:  func(r *RouteCreate) { r.Days = nil },
			wantErr: core_errors.ErrInvalidRouteDaysAmount,
		},
		{
			name:    "empty days",
			mutate:  func(r *RouteCreate) { r.Days = []RouteDay{} },
			wantErr: core_errors.ErrInvalidRouteDaysAmount,
		},
		{
			name:    "day without stops",
			mutate:  func(r *RouteCreate) { r.Days = []RouteDay{{ID: 1, Date: futureDate, DayNumber: 1}} },
			wantErr: core_errors.ErrInvalidRouteDayStopsAmount,
		},
		{
			name: "day with invalid stop",
			mutate: func(r *RouteCreate) {
				r.Days = []RouteDay{{ID: 1, Date: futureDate, DayNumber: 1, Stops: []RouteStop{{ID: 1, PlaceID: 1, Position: 0}}, Legs: []RouteLeg{}}}
			},
			wantErr: core_errors.ErrInvalidRouteStopPosition,
		},
		{
			name: "day with invalid date",
			mutate: func(r *RouteCreate) {
				r.Days = []RouteDay{{ID: 1, Date: time.Time{}, DayNumber: 1, Stops: []RouteStop{validStop(futureDate)}, Legs: []RouteLeg{}}}
			},
			wantErr: core_errors.ErrInvalidRouteDayDate,
		},
		{
			name: "day with invalid day number",
			mutate: func(r *RouteCreate) {
				r.Days = []RouteDay{{ID: 1, Date: futureDate, DayNumber: 0, Stops: []RouteStop{validStop(futureDate)}, Legs: []RouteLeg{}}}
			},
			wantErr: core_errors.ErrInvalidRouteDayDayNumber,
		},
		{
			name: "day with invalid leg",
			mutate: func(r *RouteCreate) {
				r.Days = []RouteDay{{
					ID:        1,
					Date:      futureDate,
					DayNumber: 1,
					Stops:     []RouteStop{validStop(futureDate)},
					Legs: []RouteLeg{
						{
							ID:         1,
							FromStopID: 1,
							ToStopID:   2,
							TravelMode: "unknown",
							Distance:   5,
							Duration:   10 * time.Minute,
						},
					},
				}}
			},
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routeCreate := validRouteCreate
			if tt.mutate != nil {
				tt.mutate(&routeCreate)
			}

			err := routeCreate.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
