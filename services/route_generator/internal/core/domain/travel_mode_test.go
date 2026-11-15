package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func TestTravelMode_Constants(t *testing.T) {
	tests := []struct {
		name       string
		travelMode TravelMode
		want       string
	}{
		{name: "walking", travelMode: TravelModeWalking, want: "walking"},
		{name: "cycling", travelMode: TravelModeCycling, want: "cycling"},
		{name: "driving", travelMode: TravelModeDriving, want: "driving"},
		{name: "transit", travelMode: TravelModeTransit, want: "transit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.travelMode))
		})
	}
}

func TestTravelMode_Validation(t *testing.T) {
	tests := []struct {
		name    string
		mode    TravelMode
		wantErr error
	}{
		{
			name: "valid mode - walking",
			mode: TravelModeWalking,
		},
		{
			name: "valid mode - cycling",
			mode: TravelModeCycling,
		},
		{
			name: "valid mode - driving",
			mode: TravelModeDriving,
		},
		{
			name: "valid mode - transit",
			mode: TravelModeTransit,
		},
		{
			name:    "invalid mode - unknown",
			mode:    "unknown",
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			name:    "invalid mode - empty",
			mode:    "",
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			name:    "invalid mode - uppercase",
			mode:    "WALKING",
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			name:    "invalid mode - mixed case",
			mode:    "Walking",
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			name:    "invalid mode - with spaces",
			mode:    "walking mode",
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			name:    "invalid mode - with numbers",
			mode:    "walking123",
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			name:    "invalid mode - with special characters",
			mode:    "walking!",
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			name:    "invalid mode - with underscore",
			mode:    "walking_mode",
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			name:    "invalid mode - with hyphen",
			mode:    "walking-mode",
			wantErr: core_errors.ErrUnsupportedRouteTravelMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteTravelMode(tt.mode)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestTravelMode_Equality(t *testing.T) {
	assert.True(t, TravelModeWalking == TravelMode("walking"), "Same mode should be equal")
	assert.False(t, TravelModeWalking == TravelMode("cycling"), "Different modes should not be equal")
	assert.False(t, TravelModeWalking == TravelMode("WALKING"), "Different case should not be equal")
}

func TestTravelMode_Hashability(t *testing.T) {
	// Test that TravelMode can be used as a map key
	modeMap := make(map[TravelMode]bool)

	modeMap[TravelModeWalking] = true
	modeMap[TravelModeCycling] = true
	modeMap[TravelModeDriving] = true
	modeMap[TravelModeTransit] = true

	assert.True(t, modeMap[TravelModeWalking], "Walking mode should be in map")
	assert.True(t, modeMap[TravelModeCycling], "Cycling mode should be in map")
	assert.True(t, modeMap[TravelModeDriving], "Driving mode should be in map")
	assert.True(t, modeMap[TravelModeTransit], "Transit mode should be in map")
	assert.False(t, modeMap["unknown"], "Unknown mode should not be in map")
	assert.False(t, modeMap[""], "Empty mode should not be in map")
	assert.False(t, modeMap["WALKING"], "Uppercase mode should not be in map")
}
