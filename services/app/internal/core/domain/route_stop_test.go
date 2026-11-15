package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestRouteStop_Validate(t *testing.T) {
	tests := []struct {
		name    string
		stop    RouteStop
		wantErr error
	}{
		{
			name: "valid stop",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
				VisitDuration: 2 * time.Hour,
			},
		},
		{
			name: "invalid ID",
			stop: RouteStop{
				ID:            -1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
				VisitDuration: 2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid place ID",
			stop: RouteStop{
				ID:            1,
				PlaceID:       -1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
				VisitDuration: 2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid position",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      0,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
				VisitDuration: 2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidRouteStopPosition,
		},
		{
			name: "invalid name",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
				VisitDuration: 2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidRouteStopName,
		},
		{
			name: "invalid description",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
				VisitDuration: 2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidRouteStopDescription,
		},
		{
			name: "invalid coordinates",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 200, Longitude: 200},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
				VisitDuration: 2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "invalid arrival time",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Time{},
				DepartureTime: time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
				VisitDuration: 2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidRouteArrivalTime,
		},
		{
			name: "invalid departure time",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Time{},
				VisitDuration: 2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidRouteDepartureTime,
		},
		{
			name: "invalid times order - arrival after departure",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				VisitDuration: 2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidRouteStopDatesOrder,
		},
		{
			name: "invalid times order - duration exceeds departure",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 13, 0, 0, 0, time.UTC),
				VisitDuration: 3 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidRouteStopDatesOrder,
		},
		{
			name: "invalid times order - negative duration",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 10, 0, 0, 0, time.UTC),
				VisitDuration: -2 * time.Hour,
			},
			wantErr: core_errors.ErrInvalidRouteStopDatesOrder,
		},
		{
			name: "valid stop with zero duration",
			stop: RouteStop{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				DepartureTime: time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
				VisitDuration: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.stop.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestRouteStop_Positions(t *testing.T) {
	tests := []struct {
		name        string
		position    int
		shouldError bool
	}{
		{
			name:        "position 1 is valid",
			position:    1,
			shouldError: false,
		},
		{
			name:        "position 100 is valid",
			position:    100,
			shouldError: false,
		},
		{
			name:        "position 0 is invalid",
			position:    0,
			shouldError: true,
		},
		{
			name:        "position negative is invalid",
			position:    -1,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteStopPosition(tt.position)

			if tt.shouldError {
				assert.ErrorIs(t, err, core_errors.ErrInvalidRouteStopPosition)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRouteStop_NamesAndDescriptions(t *testing.T) {
	tests := []struct {
		name        string
		stopName    string
		description string
		wantNameErr bool
		wantDescErr bool
	}{
		{
			name:        "valid name and description",
			stopName:    "Stop 1",
			description: "Description",
		},
		{
			name:        "empty name is invalid",
			stopName:    "",
			description: "Description",
			wantNameErr: true,
		},
		{
			name:        "empty description is invalid",
			stopName:    "Stop 1",
			description: "",
			wantDescErr: true,
		},
		{
			name:        "name and description with spaces are valid",
			stopName:    "Stop 1",
			description: "Description with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nameErr := validateRouteStopName(tt.stopName)
			descErr := validateRouteStopDescription(tt.description)

			if tt.wantNameErr {
				assert.Error(t, nameErr)
			} else {
				require.NoError(t, nameErr)
			}

			if tt.wantDescErr {
				assert.Error(t, descErr)
			} else {
				require.NoError(t, descErr)
			}
		})
	}
}
