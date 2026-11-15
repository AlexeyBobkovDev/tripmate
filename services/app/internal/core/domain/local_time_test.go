package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestLocalTime_Validate(t *testing.T) {
	tests := []struct {
		name    string
		time    LocalTime
		wantErr error
	}{
		{
			name: "valid time - 09:00",
			time: "09:00",
		},
		{
			name: "valid time - 00:00",
			time: "00:00",
		},
		{
			name: "valid time - 23:59",
			time: "23:59",
		},
		{
			name: "valid time - 12:34",
			time: "12:34",
		},
		{
			name:    "invalid time - empty",
			time:    "",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - missing hours",
			time:    ":00",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - missing minutes",
			time:    "12:",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - extra characters",
			time:    "12:34:56",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - only hours",
			time:    "12",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - only minutes",
			time:    "34",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - negative hours",
			time:    "-09:00",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - negative minutes",
			time:    "09:-00",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - hours > 23",
			time:    "24:00",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - hours < 0",
			time:    "-01:00",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - minutes > 59",
			time:    "00:60",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - minutes < 0",
			time:    "00:-01",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - non-numeric",
			time:    "ab:cd",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - letters",
			time:    "12:abc",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "invalid time - all letters",
			time:    "abcd:efgh",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "valid time - single digit hour",
			time: "1:00",
		},
		{
			name: "valid time - single digit hour without leading zero",
			time: "9:00",
		},
		{
			name:    "invalid time - single digit minute",
			time:    "09:5",
			wantErr: core_errors.ErrInvalidLocalTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.time.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestLocalTime_MinutesSinceMidnight(t *testing.T) {
	tests := []struct {
		name    string
		time    LocalTime
		want    int
		wantErr error
	}{
		{
			name: "midnight",
			time: "00:00",
			want: 0,
		},
		{
			name: "1 hour",
			time: "01:00",
			want: 60,
		},
		{
			name: "12:00 noon",
			time: "12:00",
			want: 720,
		},
		{
			name: "23:59",
			time: "23:59",
			want: 1439,
		},
		{
			name: "random time",
			time: "14:25",
			want: 865,
		},
		{
			name:    "invalid time",
			time:    "invalid",
			want:    0,
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "empty time",
			time:    "",
			want:    0,
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name:    "time with extra characters",
			time:    "12:34:56",
			want:    0,
			wantErr: core_errors.ErrInvalidLocalTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.time.MinutesSinceMidnight()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, 0, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
