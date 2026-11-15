package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestTravelMode_Validation(t *testing.T) {
	tests := []struct {
		testName string
		mode     TravelMode
		wantErr  error
	}{
		{
			testName: "valid mode - walking",
			mode:     TravelModeWalking,
		},
		{
			testName: "valid mode - cycling",
			mode:     TravelModeCycling,
		},
		{
			testName: "valid mode - driving",
			mode:     TravelModeDriving,
		},
		{
			testName: "valid mode - transit",
			mode:     TravelModeTransit,
		},
		{
			testName: "invalid mode - unknown",
			mode:     "unknown",
			wantErr:  core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			testName: "invalid mode - empty",
			mode:     "",
			wantErr:  core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			testName: "invalid mode - uppercase",
			mode:     "WALKING",
			wantErr:  core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			testName: "invalid mode - mixed case",
			mode:     "Walking",
			wantErr:  core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			testName: "invalid mode - with spaces",
			mode:     "walking mode",
			wantErr:  core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			testName: "invalid mode - with numbers",
			mode:     "walking123",
			wantErr:  core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			testName: "invalid mode - with special characters",
			mode:     "walking!",
			wantErr:  core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			testName: "invalid mode - with underscore",
			mode:     "walking_mode",
			wantErr:  core_errors.ErrUnsupportedRouteTravelMode,
		},
		{
			testName: "invalid mode - with hyphen",
			mode:     "walking-mode",
			wantErr:  core_errors.ErrUnsupportedRouteTravelMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := validateRouteTravelMode(tt.mode)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestTravelMode_Constants(t *testing.T) {
	tests := []struct {
		testName string
		mode     TravelMode
		expected string
	}{
		{
			testName: "TravelModeWalking constant",
			mode:     TravelModeWalking,
			expected: "walking",
		},
		{
			testName: "TravelModeCycling constant",
			mode:     TravelModeCycling,
			expected: "cycling",
		},
		{
			testName: "TravelModeDriving constant",
			mode:     TravelModeDriving,
			expected: "driving",
		},
		{
			testName: "TravelModeTransit constant",
			mode:     TravelModeTransit,
			expected: "transit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.mode))
		})
	}
}

func TestTravelMode_ValueComparison(t *testing.T) {
	tests := []struct {
		testName    string
		mode1       TravelMode
		mode2       TravelMode
		shouldEqual bool
	}{
		{
			testName:    "same mode",
			mode1:       TravelModeWalking,
			mode2:       TravelModeWalking,
			shouldEqual: true,
		},
		{
			testName:    "different modes",
			mode1:       TravelModeWalking,
			mode2:       TravelModeCycling,
			shouldEqual: false,
		},
		{
			testName:    "same mode string value",
			mode1:       TravelModeWalking,
			mode2:       TravelMode("walking"),
			shouldEqual: true,
		},
		{
			testName:    "different case string value",
			mode1:       TravelModeWalking,
			mode2:       TravelMode("WALKING"),
			shouldEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.shouldEqual {
				assert.Equal(t, tt.mode1, tt.mode2, "Travel modes should be equal")
				assert.Equal(t, string(tt.mode1), string(tt.mode2), "String representations should be equal")
			} else {
				assert.NotEqual(t, tt.mode1, tt.mode2, "Travel modes should not be equal")
				assert.NotEqual(t, string(tt.mode1), string(tt.mode2), "String representations should not be equal")
			}
		})
	}
}

func TestTravelMode_Hashability(t *testing.T) {
	// Test that TravelMode is hashable by using it as a map key
	modeMap := make(map[TravelMode]bool)

	// Add valid modes to the map
	modeMap[TravelModeWalking] = true
	modeMap[TravelModeCycling] = true
	modeMap[TravelModeDriving] = true
	modeMap[TravelModeTransit] = true

	// Test that we can retrieve the modes
	assert.True(t, modeMap[TravelModeWalking], "Walking mode should be in map")
	assert.True(t, modeMap[TravelModeCycling], "Cycling mode should be in map")
	assert.True(t, modeMap[TravelModeDriving], "Driving mode should be in map")
	assert.True(t, modeMap[TravelModeTransit], "Transit mode should be in map")

	// Test that invalid modes are not in the map
	assert.False(t, modeMap["unknown"], "Unknown mode should not be in map")
	assert.False(t, modeMap[""], "Empty mode should not be in map")
	assert.False(t, modeMap["WALKING"], "Uppercase mode should not be in map")
}
