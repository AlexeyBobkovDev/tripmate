package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestRouteLeg_Validate(t *testing.T) {
	tests := []struct {
		name    string
		leg     RouteLeg
		wantErr error
	}{
		{
			name: "valid leg - walking",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   1000,
				Duration:   15 * time.Minute,
			},
		},
		{
			name: "valid leg - cycling",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeCycling,
				Distance:   5000,
				Duration:   30 * time.Minute,
			},
		},
		{
			name: "valid leg - driving",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeDriving,
				Distance:   10000,
				Duration:   1 * time.Hour,
			},
		},
		{
			name: "valid leg - transit",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeTransit,
				Distance:   8000,
				Duration:   45 * time.Minute,
			},
		},
		{
			name: "invalid ID",
			leg: RouteLeg{
				ID:         -1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   1000,
				Duration:   15 * time.Minute,
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid from stop ID",
			leg: RouteLeg{
				ID:         1,
				FromStopID: -1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   1000,
				Duration:   15 * time.Minute,
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid to stop ID",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   -1,
				TravelMode: TravelModeWalking,
				Distance:   1000,
				Duration:   15 * time.Minute,
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid travel mode",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: "invalid",
				Distance:   1000,
				Duration:   15 * time.Minute,
			},
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			name: "invalid distance",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   0,
				Duration:   15 * time.Minute,
			},
			wantErr: core_errors.ErrInvalidRouteDistance,
		},
		{
			name: "invalid distance - negative",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   -1000,
				Duration:   15 * time.Minute,
			},
			wantErr: core_errors.ErrInvalidRouteDistance,
		},
		{
			name: "invalid duration",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   1000,
				Duration:   0,
			},
			wantErr: core_errors.ErrInvalidRouteDuration,
		},
		{
			name: "invalid duration - negative",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   1000,
				Duration:   -15 * time.Minute,
			},
			wantErr: core_errors.ErrInvalidRouteDuration,
		},
		{
			// Referential integrity (including from == to) is enforced by
			// validateStopsAndLegs, which has the full stop list available.
			name: "start equals end is not checked at leg level",
			leg: RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   1,
				TravelMode: TravelModeWalking,
				Distance:   1000,
				Duration:   15 * time.Minute,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.leg.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestRouteLeg_TravelModes(t *testing.T) {
	tests := []struct {
		name        string
		travelMode  TravelMode
		shouldError bool
	}{
		{
			name:        "walking is valid",
			travelMode:  TravelModeWalking,
			shouldError: false,
		},
		{
			name:        "cycling is valid",
			travelMode:  TravelModeCycling,
			shouldError: false,
		},
		{
			name:        "driving is valid",
			travelMode:  TravelModeDriving,
			shouldError: false,
		},
		{
			name:        "transit is valid",
			travelMode:  TravelModeTransit,
			shouldError: false,
		},
		{
			name:        "walking uppercase is invalid",
			travelMode:  "Walking",
			shouldError: true,
		},
		{
			name:        "invalid mode is invalid",
			travelMode:  "flying",
			shouldError: true,
		},
		{
			name:        "empty is invalid",
			travelMode:  "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteTravelMode(tt.travelMode)

			if tt.shouldError {
				assert.ErrorIs(t, err, core_errors.ErrUnsupportedRouteTravelMode)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRouteLeg_DistanceAndDuration(t *testing.T) {
	tests := []struct {
		name     string
		distance int
		duration time.Duration
	}{
		{
			name:     "valid distance and duration",
			distance: 1000,
			duration: 15 * time.Minute,
		},
		{
			name:     "large distance and duration",
			distance: 100000,
			duration: 2 * time.Hour,
		},
		{
			name:     "small distance and duration",
			distance: 10,
			duration: 1 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leg := RouteLeg{
				ID:         1,
				FromStopID: 1,
				ToStopID:   2,
				TravelMode: TravelModeWalking,
				Distance:   tt.distance,
				Duration:   tt.duration,
			}

			err := leg.Validate()
			require.NoError(t, err)

			// Test that getters work
			assert.Equal(t, tt.distance, leg.Distance)
			assert.Equal(t, tt.duration, leg.Duration)
		})
	}
}
