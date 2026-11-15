package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestInterest_Validation(t *testing.T) {
	tests := []struct {
		name     string
		interest Interest
		wantErr  error
	}{
		{
			name:     "valid interest - history",
			interest: InterestHistory,
		},
		{
			name:     "valid interest - nature",
			interest: InterestNature,
		},
		{
			name:     "valid interest - food",
			interest: InterestFood,
		},
		{
			name:     "valid interest - architecture",
			interest: InterestArchitecture,
		},
		{
			name:     "invalid interest - unknown",
			interest: "unknown",
			wantErr:  core_errors.ErrUnsupportedInterest,
		},
		{
			name:     "invalid interest - empty",
			interest: "",
			wantErr:  core_errors.ErrInvalidInterest,
		},

		{
			name:     "invalid interest - uppercase",
			interest: "HISTORY",
			wantErr:  core_errors.ErrUnsupportedInterest,
		},
		{
			name:     "invalid interest - mixed case",
			interest: "History",
			wantErr:  core_errors.ErrUnsupportedInterest,
		},
		{
			name:     "invalid interest - with spaces",
			interest: "history stuff",
			wantErr:  core_errors.ErrUnsupportedInterest,
		},
		{
			name:     "invalid interest - with numbers",
			interest: "history123",
			wantErr:  core_errors.ErrUnsupportedInterest,
		},
		{
			name:     "invalid interest - with special characters",
			interest: "history!",
			wantErr:  core_errors.ErrUnsupportedInterest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInterests([]Interest{tt.interest})

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v but got nil", tt.wantErr)
				} else {
					t.Logf("Got error: %v", err)
					assert.ErrorIs(t, err, tt.wantErr)
				}
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestInterest_Constants(t *testing.T) {
	assert.Equal(t, "history", string(InterestHistory))
	assert.Equal(t, "nature", string(InterestNature))
	assert.Equal(t, "food", string(InterestFood))
	assert.Equal(t, "architecture", string(InterestArchitecture))
}

func TestInterest_Equality(t *testing.T) {
	assert.True(t, InterestHistory == Interest("history"), "Same interest should be equal")
	assert.False(t, InterestHistory == Interest("nature"), "Different interests should not be equal")
}

func TestInterest_Hashability(t *testing.T) {
	// Test that Interest can be used as a map key
	interestMap := make(map[Interest]bool)

	interestMap[InterestHistory] = true
	interestMap[InterestNature] = true
	interestMap[InterestFood] = true
	interestMap[InterestArchitecture] = true

	assert.True(t, interestMap[InterestHistory], "History interest should be in map")
	assert.True(t, interestMap[InterestNature], "Nature interest should be in map")
	assert.True(t, interestMap[InterestFood], "Food interest should be in map")
	assert.True(t, interestMap[InterestArchitecture], "Architecture interest should be in map")
	assert.False(t, interestMap["unknown"], "Unknown interest should not be in map")
}

func TestAllowedInterests_Map(t *testing.T) {
	assert.Contains(t, AllowedInterests, InterestHistory)
	assert.Contains(t, AllowedInterests, InterestNature)
	assert.Contains(t, AllowedInterests, InterestFood)
	assert.Contains(t, AllowedInterests, InterestArchitecture)
	assert.NotContains(t, AllowedInterests, "unknown")
	assert.NotContains(t, AllowedInterests, "")
}
