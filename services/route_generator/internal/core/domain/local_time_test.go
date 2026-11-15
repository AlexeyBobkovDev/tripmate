package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func TestLocalTime_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   LocalTime
		wantErr error
	}{
		{name: "midnight", value: "00:00"},
		{name: "single digit midnight", value: "0:00"},
		{name: "typical morning time", value: "09:00"},
		{name: "single digit hour", value: "9:05"},
		{name: "noon", value: "12:00"},
		{name: "last minute of the day", value: "23:59"},
		{name: "empty value", value: "", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "hour out of range", value: "24:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "hour far out of range", value: "99:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "minute out of range", value: "23:60", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "minute far out of range", value: "23:99", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "missing minute digits", value: "09:0", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "single digit minute", value: "9:5", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "seconds included", value: "09:00:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "fractional seconds included", value: "09:00.0", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "time zone suffix included", value: "09:00Z", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "AM/PM suffix included", value: "9:05PM", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "leading whitespace", value: " 09:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "trailing whitespace", value: "09:00 ", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "trailing newline", value: "09:00\n", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "negative hour", value: "-1:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "explicit positive hour", value: "+9:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "non-numeric value", value: "ab:cd", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "fullwidth digits", value: "９:００", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "too many hour digits", value: "009:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "missing separator", value: "0900", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "trailing digit", value: "09:000", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "date instead of time", value: "2026-11-14", wantErr: core_errors.ErrInvalidLocalTime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := tt.value

			err := value.Validate()

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestLocalTime_Validate_AllMinutesOfDay(t *testing.T) {
	t.Parallel()

	for hour := range 24 {
		for minute := range 60 {
			value := LocalTime(fmt.Sprintf("%02d:%02d", hour, minute))
			require.NoErrorf(t, value.Validate(), "value %q must be valid", value)

			if hour >= 10 {
				continue
			}

			singleDigitHour := LocalTime(fmt.Sprintf("%d:%02d", hour, minute))
			require.NoErrorf(t, singleDigitHour.Validate(), "value %q must be valid", singleDigitHour)
		}
	}
}

// time.Parse("15:04", ...) already rejects out-of-range hours and minutes, so
// the dedicated ErrInvalidLocalTimeHours and ErrInvalidLocalTimeMinutes
// sentinels are not reachable through LocalTime.Validate. Out-of-range values
// must still be reported as ErrInvalidLocalTime.
func TestLocalTime_Validate_OutOfRangeValuesUseGenericError(t *testing.T) {
	t.Parallel()

	values := []LocalTime{"24:00", "99:00", "23:60", "23:99"}

	for _, value := range values {
		err := value.Validate()

		require.ErrorIs(t, err, core_errors.ErrInvalidLocalTime)
		require.NotErrorIs(t, err, core_errors.ErrInvalidLocalTimeHours)
		require.NotErrorIs(t, err, core_errors.ErrInvalidLocalTimeMinutes)
	}
}

func TestLocalTime_Validate_DoesNotMutateValue(t *testing.T) {
	t.Parallel()

	value := LocalTime("9:05")

	require.NoError(t, value.Validate())
	require.Equal(t, LocalTime("9:05"), value)
}

func TestLocalTime_MinutesSinceMidnight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   LocalTime
		want    int
		wantErr error
	}{
		{name: "midnight", value: "00:00", want: 0},
		{name: "single digit midnight", value: "0:00", want: 0},
		{name: "start of working day", value: "09:00", want: 540},
		{name: "half past nine", value: "09:30", want: 570},
		{name: "single digit hour", value: "9:05", want: 545},
		{name: "noon", value: "12:00", want: 720},
		{name: "last minute of the day", value: "23:59", want: 1439},
		{name: "empty value", value: "", want: 0, wantErr: core_errors.ErrInvalidLocalTime},
		{name: "hour out of range", value: "24:00", want: 0, wantErr: core_errors.ErrInvalidLocalTime},
		{name: "minute out of range", value: "12:60", want: 0, wantErr: core_errors.ErrInvalidLocalTime},
		{name: "single digit minute", value: "12:5", want: 0, wantErr: core_errors.ErrInvalidLocalTime},
		{name: "trailing whitespace", value: "12:00 ", want: 0, wantErr: core_errors.ErrInvalidLocalTime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value := tt.value

			got, err := value.MinutesSinceMidnight()

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Zero(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestLocalTime_MinutesSinceMidnight_AllMinutesOfDay(t *testing.T) {
	t.Parallel()

	for hour := range 24 {
		for minute := range 60 {
			value := LocalTime(fmt.Sprintf("%02d:%02d", hour, minute))
			want := hour*60 + minute

			got, err := value.MinutesSinceMidnight()

			require.NoErrorf(t, err, "value %q", value)
			require.Equalf(t, want, got, "value %q", value)
		}
	}
}

func TestLocalTime_MinutesSinceMidnight_DoesNotMutateValue(t *testing.T) {
	t.Parallel()

	value := LocalTime("9:05")

	got, err := value.MinutesSinceMidnight()

	require.NoError(t, err)
	require.Equal(t, 545, got)
	require.Equal(t, LocalTime("9:05"), value)
}
