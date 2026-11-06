package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func TestRouteStop_Validate(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)

	validStop := RouteStop{
		ID:            1,
		PlaceID:       1,
		Position:      1,
		Name:          "Stop",
		Description:   "Stop description",
		Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		ArrivalTime:   futureDate.Add(time.Hour),
		DepartureTime: futureDate.Add(2 * time.Hour),
		VisitDuration: 30 * time.Minute,
	}

	tests := []struct {
		name    string
		mutate  func(*RouteStop)
		wantErr error
	}{
		{name: "valid stop"},
		{name: "invalid stop ID", mutate: func(s *RouteStop) { s.ID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid place ID", mutate: func(s *RouteStop) { s.PlaceID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid position", mutate: func(s *RouteStop) { s.Position = 0 }, wantErr: core_errors.ErrInvalidRouteStopPosition},
		{name: "invalid name", mutate: func(s *RouteStop) { s.Name = "" }, wantErr: core_errors.ErrInvalidRouteStopName},
		{name: "invalid description", mutate: func(s *RouteStop) { s.Description = "" }, wantErr: core_errors.ErrInvalidRouteStopDescription},
		{name: "invalid coordinates", mutate: func(s *RouteStop) { s.Coordinates = Coordinates{Latitude: 91, Longitude: 37.6208} }, wantErr: core_errors.ErrInvalidLatitude},
		{name: "zero arrival time", mutate: func(s *RouteStop) { s.ArrivalTime = time.Time{} }, wantErr: core_errors.ErrInvalidRouteArrivalTime},
		{name: "zero departure time", mutate: func(s *RouteStop) { s.DepartureTime = time.Time{} }, wantErr: core_errors.ErrInvalidRouteDepartureTime},
		{name: "negative visit duration", mutate: func(s *RouteStop) { s.VisitDuration = -time.Minute }, wantErr: core_errors.ErrInvalidRouteStopDatesOrder},
		{name: "arrival after departure", mutate: func(s *RouteStop) { s.DepartureTime = s.ArrivalTime.Add(-time.Hour) }, wantErr: core_errors.ErrInvalidRouteStopDatesOrder},
		{name: "visit duration after departure", mutate: func(s *RouteStop) { s.VisitDuration = 3 * time.Hour }, wantErr: core_errors.ErrInvalidRouteStopDatesOrder},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stop := validStop
			if tt.mutate != nil {
				tt.mutate(&stop)
			}
			err := stop.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
