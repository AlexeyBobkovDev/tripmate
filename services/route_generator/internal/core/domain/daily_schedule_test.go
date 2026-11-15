package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func validDailySchedule() DailySchedule {
	return DailySchedule{
		StartTime: "09:00",
		EndTime:   "18:00",
	}
}

func TestDailySchedule_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		startTime LocalTime
		endTime   LocalTime
		wantErr   error
	}{
		{name: "typical schedule", startTime: "09:00", endTime: "18:00"},
		{name: "same start and end", startTime: "12:00", endTime: "12:00"},
		{name: "midnight to almost midnight", startTime: "00:00", endTime: "23:59"},
		{name: "single digit start hour", startTime: "9:00", endTime: "18:00"},
		{name: "single digit end hour", startTime: "09:00", endTime: "9:30"},
		{name: "one minute long schedule", startTime: "09:00", endTime: "09:01"},
		{name: "one minute before midnight", startTime: "23:58", endTime: "23:59"},
		{name: "empty start time", startTime: "", endTime: "18:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "empty end time", startTime: "09:00", endTime: "", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "invalid start time", startTime: "9", endTime: "18:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "invalid end time", startTime: "09:00", endTime: "25:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "invalid start minute", startTime: "09:5", endTime: "18:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "invalid end minute", startTime: "09:00", endTime: "18:5", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "start time with whitespace", startTime: " 09:00", endTime: "18:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "start after end", startTime: "18:00", endTime: "09:00", wantErr: core_errors.ErrInvalidTripEndTime},
		{name: "start one minute after end", startTime: "09:01", endTime: "09:00", wantErr: core_errors.ErrInvalidTripEndTime},
		{name: "typical hours in wrong order", startTime: "10:00", endTime: "09:59", wantErr: core_errors.ErrInvalidTripEndTime},
		{name: "almost full day in wrong order", startTime: "23:59", endTime: "00:00", wantErr: core_errors.ErrInvalidTripEndTime},
		{name: "both times invalid reports start error first", startTime: "", endTime: "25:00", wantErr: core_errors.ErrInvalidLocalTime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			schedule := DailySchedule{
				StartTime: tt.startTime,
				EndTime:   tt.endTime,
			}

			err := schedule.Validate()

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestDailySchedule_Validate_ZeroValue(t *testing.T) {
	t.Parallel()

	var schedule DailySchedule

	require.ErrorIs(t, schedule.Validate(), core_errors.ErrInvalidLocalTime)
}

func TestDailySchedule_Validate_InvalidEndTimeIsReported(t *testing.T) {
	t.Parallel()

	schedule := DailySchedule{
		StartTime: "09:00",
		EndTime:   "24:00",
	}

	require.ErrorIs(t, schedule.Validate(), core_errors.ErrInvalidLocalTime)
}

func TestDailySchedule_Validate_DoesNotMutateSchedule(t *testing.T) {
	t.Parallel()

	schedule := validDailySchedule()
	before := schedule

	require.NoError(t, schedule.Validate())
	require.Equal(t, before, schedule)
}
