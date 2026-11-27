package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestDailySchedule_Validate(t *testing.T) {
	tests := []struct {
		name     string
		schedule DailySchedule
		wantErr  error
	}{
		{
			name: "valid daily schedule",
			schedule: DailySchedule{
				StartTime: "09:00",
				EndTime:   "18:00",
			},
		},
		{
			name: "valid daily schedule - midnight to midnight",
			schedule: DailySchedule{
				StartTime: "00:00",
				EndTime:   "00:00",
			},
		},
		{
			name: "valid daily schedule - all day",
			schedule: DailySchedule{
				StartTime: "00:00",
				EndTime:   "23:59",
			},
		},
		{
			name: "invalid start time",
			schedule: DailySchedule{
				StartTime: "invalid",
				EndTime:   "18:00",
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid end time",
			schedule: DailySchedule{
				StartTime: "09:00",
				EndTime:   "invalid",
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "start time after end time",
			schedule: DailySchedule{
				StartTime: "18:00",
				EndTime:   "09:00",
			},
			wantErr: core_errors.ErrInvalidTripEndTime,
		},
		{
			name: "start time equals end time",
			schedule: DailySchedule{
				StartTime: "12:00",
				EndTime:   "12:00",
			},
		},
		{
			name: "invalid start time - hours > 23",
			schedule: DailySchedule{
				StartTime: "24:00",
				EndTime:   "18:00",
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid end time - minutes > 59",
			schedule: DailySchedule{
				StartTime: "09:00",
				EndTime:   "18:60",
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schedule.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDailySchedule_Validate_EdgeCases(t *testing.T) {
	t.Run("valid schedule", func(t *testing.T) {
		err := (&DailySchedule{StartTime: "09:00", EndTime: "18:00"}).Validate()
		require.NoError(t, err)
	})

	t.Run("invalid end time", func(t *testing.T) {
		err := (&DailySchedule{StartTime: "09:00", EndTime: "25:00"}).Validate()
		require.ErrorIs(t, err, core_errors.ErrInvalidLocalTime)
	})

	t.Run("start time after end time", func(t *testing.T) {
		err := (&DailySchedule{StartTime: "18:00", EndTime: "09:00"}).Validate()
		require.ErrorIs(t, err, core_errors.ErrInvalidTripEndTime)
	})
}

func TestDailySchedule_Validate_MinutesSinceMidnightFailures(t *testing.T) {
	// BUG: LocalTime.Validate rejects these formats before
	// MinutesSinceMidnight can fail, so DailySchedule.Validate can never
	// reach its MinutesSinceMidnight error branches (daily_schedule.go:22,
	// 27) — "24:00" parses in "15:04" as 0:00 next day, so even hours=24
	// passes time.Parse; Validate's own range check catches it instead.
	tests := []struct {
		name     string
		schedule DailySchedule
	}{
		{
			name:     "invalid start time fails at MinutesSinceMidnight",
			schedule: DailySchedule{StartTime: "not-a-time", EndTime: "18:00"},
		},
		{
			name:     "invalid end time fails at MinutesSinceMidnight",
			schedule: DailySchedule{StartTime: "09:00", EndTime: "not-a-time"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schedule.Validate()

			require.ErrorIs(t, err, core_errors.ErrInvalidLocalTime)
		})
	}
}
