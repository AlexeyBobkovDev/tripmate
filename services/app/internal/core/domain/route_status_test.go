package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestRouteStatus_Validate(t *testing.T) {
	tests := []struct {
		name    string
		status  RouteStatus
		wantErr error
	}{
		{
			name:   "valid status - generated",
			status: RouteStatusGenerated,
		},
		{
			name:   "valid status - edited",
			status: RouteStatusEdited,
		},
		{
			name:    "invalid status - unknown",
			status:  "unknown",
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name:    "invalid status - empty",
			status:  "",
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name:    "invalid status - uppercase",
			status:  "GENERATED",
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name:    "invalid status - mixed case",
			status:  "Generated",
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name:    "invalid status - with spaces",
			status:  "generated status",
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name:    "invalid status - with numbers",
			status:  "generated123",
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name:    "invalid status - with special characters",
			status:  "generated!",
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRouteStatus(tt.status)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestRouteStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   RouteStatus
		expected string
	}{
		{
			name:     "RouteStatusGenerated constant",
			status:   RouteStatusGenerated,
			expected: "generated",
		},
		{
			name:     "RouteStatusEdited constant",
			status:   RouteStatusEdited,
			expected: "edited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestRouteStatus_StringRepresentation(t *testing.T) {
	tests := []struct {
		name     string
		status   RouteStatus
		expected string
	}{
		{
			name:     "RouteStatusGenerated string representation",
			status:   RouteStatusGenerated,
			expected: "generated",
		},
		{
			name:     "RouteStatusEdited string representation",
			status:   RouteStatusEdited,
			expected: "edited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}
