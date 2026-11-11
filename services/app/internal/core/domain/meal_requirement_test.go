package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func validMealBudget() Money {
	return Money{
		AmountMinorUnits: 50000,
		Currency:         "RUB",
	}
}

func validMealRequirement() MealRequirement {
	return MealRequirement{
		Type: MealBreakfast,
		Window: TimeWindow{
			Start: "08:00",
			End:   "09:00",
		},
		Duration: 30 * time.Minute,
		Budget:   validMealBudget(),
	}
}

func TestMealRequirement_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(*MealRequirement)
		wantErr error
	}{
		{name: "valid breakfast"},
		{
			name: "valid lunch",
			mutate: func(requirement *MealRequirement) {
				requirement.Type = MealLunch
			},
		},
		{
			name: "valid dinner",
			mutate: func(requirement *MealRequirement) {
				requirement.Type = MealDinner
			},
		},
		{
			name: "duration equals window length",
			mutate: func(requirement *MealRequirement) {
				requirement.Duration = time.Hour
			},
		},
		{
			name: "minimum positive duration",
			mutate: func(requirement *MealRequirement) {
				requirement.Duration = time.Nanosecond
			},
		},
		{
			name: "sub-minute duration fits window",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.End = "08:02"
				requirement.Duration = 90 * time.Second
			},
		},
		{
			name: "duration equals one minute window",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.End = "08:01"
				requirement.Duration = time.Minute
			},
		},
		{
			name: "single digit window hours",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.Start = "8:00"
				requirement.Window.End = "9:00"
				requirement.Duration = time.Hour
			},
		},
		{
			name: "almost full day window",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.Start = "00:00"
				requirement.Window.End = "23:59"
				requirement.Duration = 23*time.Hour + 59*time.Minute
			},
		},
		{
			name: "large budget amount",
			mutate: func(requirement *MealRequirement) {
				requirement.Budget = Money{
					AmountMinorUnits: 1_000_000_000,
					Currency:         "RUB",
				}
			},
		},
		{
			name: "empty meal type",
			mutate: func(requirement *MealRequirement) {
				requirement.Type = ""
			},
			wantErr: core_errors.ErrUnsupportedMealType,
		},
		{
			name: "unknown meal type",
			mutate: func(requirement *MealRequirement) {
				requirement.Type = "brunch"
			},
			wantErr: core_errors.ErrUnsupportedMealType,
		},
		{
			name: "uppercase meal type",
			mutate: func(requirement *MealRequirement) {
				requirement.Type = "Breakfast"
			},
			wantErr: core_errors.ErrUnsupportedMealType,
		},
		{
			name: "meal type with whitespace",
			mutate: func(requirement *MealRequirement) {
				requirement.Type = " breakfast"
			},
			wantErr: core_errors.ErrUnsupportedMealType,
		},
		{
			name: "invalid window start",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.Start = "25:00"
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "empty window start",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.Start = ""
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid window end",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.End = ""
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid window end minutes",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.End = "09:60"
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "window start after end",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.Start = "10:00"
				requirement.Window.End = "09:00"
			},
			wantErr: core_errors.ErrInvalidTripEndTime,
		},
		{
			name: "zero duration",
			mutate: func(requirement *MealRequirement) {
				requirement.Duration = 0
			},
			wantErr: core_errors.ErrInvalidMealTimeDuration,
		},
		{
			name: "negative duration",
			mutate: func(requirement *MealRequirement) {
				requirement.Duration = -time.Minute
			},
			wantErr: core_errors.ErrInvalidMealTimeDuration,
		},
		{
			name: "negative sub-minute duration",
			mutate: func(requirement *MealRequirement) {
				requirement.Duration = -time.Nanosecond
			},
			wantErr: core_errors.ErrInvalidMealTimeDuration,
		},
		{
			name: "duration exceeds window by one nanosecond",
			mutate: func(requirement *MealRequirement) {
				requirement.Duration = time.Hour + time.Nanosecond
			},
			wantErr: core_errors.ErrInvalidMealTimeEndTime,
		},
		{
			name: "duration exceeds window",
			mutate: func(requirement *MealRequirement) {
				requirement.Duration = 61 * time.Minute
			},
			wantErr: core_errors.ErrInvalidMealTimeEndTime,
		},
		{
			name: "duration with seconds exceeds whole-minute window",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.End = "08:01"
				requirement.Duration = 60*time.Second + time.Nanosecond
			},
			wantErr: core_errors.ErrInvalidMealTimeEndTime,
		},
		{
			name: "maximum duration exceeds window",
			mutate: func(requirement *MealRequirement) {
				requirement.Duration = time.Duration(1<<63 - 1)
			},
			wantErr: core_errors.ErrInvalidMealTimeEndTime,
		},
		{
			name: "positive duration in zero-length window",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.End = requirement.Window.Start
			},
			wantErr: core_errors.ErrInvalidMealTimeEndTime,
		},
		{
			name: "zero amount budget",
			mutate: func(requirement *MealRequirement) {
				requirement.Budget = Money{AmountMinorUnits: 0, Currency: "RUB"}
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "negative budget",
			mutate: func(requirement *MealRequirement) {
				requirement.Budget = Money{AmountMinorUnits: -1, Currency: "RUB"}
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "zero value budget",
			mutate: func(requirement *MealRequirement) {
				requirement.Budget = Money{}
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "empty budget currency",
			mutate: func(requirement *MealRequirement) {
				requirement.Budget = Money{AmountMinorUnits: 50000, Currency: ""}
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "lowercase budget currency",
			mutate: func(requirement *MealRequirement) {
				requirement.Budget = Money{AmountMinorUnits: 50000, Currency: "rub"}
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "unsupported budget currency",
			mutate: func(requirement *MealRequirement) {
				requirement.Budget = Money{AmountMinorUnits: 50000, Currency: "USD"}
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "unknown type is validated before window",
			mutate: func(requirement *MealRequirement) {
				requirement.Type = "brunch"
				requirement.Window.Start = "25:00"
			},
			wantErr: core_errors.ErrUnsupportedMealType,
		},
		{
			name: "window is validated before duration",
			mutate: func(requirement *MealRequirement) {
				requirement.Window.Start = "25:00"
				requirement.Duration = 0
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "duration is validated before budget",
			mutate: func(requirement *MealRequirement) {
				requirement.Duration = 0
				requirement.Budget = Money{AmountMinorUnits: 0, Currency: "USD"}
			},
			wantErr: core_errors.ErrInvalidMealTimeDuration,
		},
		{
			name: "budget amount is validated before currency",
			mutate: func(requirement *MealRequirement) {
				requirement.Budget = Money{AmountMinorUnits: 0, Currency: "USD"}
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			requirement := validMealRequirement()
			if tt.mutate != nil {
				tt.mutate(&requirement)
			}

			err := requirement.Validate()

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMealRequirement_Validate_ZeroValue(t *testing.T) {
	t.Parallel()

	var requirement MealRequirement

	require.ErrorIs(t, requirement.Validate(), core_errors.ErrUnsupportedMealType)
}

func TestMealRequirement_Validate_DoesNotMutateRequirement(t *testing.T) {
	t.Parallel()

	requirement := validMealRequirement()
	before := requirement

	require.NoError(t, requirement.Validate())
	require.Equal(t, before, requirement)
}

func TestTimeWindow_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		start   LocalTime
		end     LocalTime
		wantErr error
	}{
		{name: "typical window", start: "09:00", end: "18:00"},
		{name: "same start and end", start: "12:00", end: "12:00"},
		{name: "full day", start: "00:00", end: "23:59"},
		{name: "one minute window", start: "23:58", end: "23:59"},
		{name: "single digit start hour", start: "9:00", end: "18:00"},
		{name: "single digit end hour", start: "09:00", end: "9:30"},
		{name: "invalid start", start: "", end: "18:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "invalid end", start: "09:00", end: "24:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "invalid start minute", start: "09:5", end: "18:00", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "invalid end minute", start: "09:00", end: "18:5", wantErr: core_errors.ErrInvalidLocalTime},
		{name: "start after end", start: "18:00", end: "09:00", wantErr: core_errors.ErrInvalidTripEndTime},
		{name: "start one minute after end", start: "09:01", end: "09:00", wantErr: core_errors.ErrInvalidTripEndTime},
		{name: "window crossing midnight", start: "23:59", end: "00:00", wantErr: core_errors.ErrInvalidTripEndTime},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			window := TimeWindow{
				Start: tt.start,
				End:   tt.end,
			}

			err := window.Validate()

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestTimeWindow_Validate_ZeroValue(t *testing.T) {
	t.Parallel()

	var window TimeWindow

	require.ErrorIs(t, window.Validate(), core_errors.ErrInvalidLocalTime)
}

func TestTimeWindow_Validate_DoesNotMutateWindow(t *testing.T) {
	t.Parallel()

	window := TimeWindow{Start: "09:00", End: "18:00"}
	before := window

	require.NoError(t, window.Validate())
	require.Equal(t, before, window)
}
