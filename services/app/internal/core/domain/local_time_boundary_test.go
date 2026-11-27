package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

// BUG: LocalTime.Validate's range checks (local_time.go:18, 21) are dead
// code: time.Parse("15:04") already rejects hours > 23 and minutes > 59, so
// the ErrInvalidLocalTimeHours / ErrInvalidLocalTimeMinutes branches can
// never run. Every input below is rejected by time.Parse first. Flip these
// assertions if the range checks are ever made reachable (e.g. by parsing
// manually).
func TestLocalTime_Validate_RangeChecksAreDeadCode(t *testing.T) {
	tests := []struct {
		name string
		time LocalTime
	}{
		{name: "hours beyond 23", time: "24:00"},
		{name: "minutes beyond 59", time: "00:60"},
		{name: "hours way beyond 23", time: "99:99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.time.Validate()

			require.ErrorIs(t, err, core_errors.ErrInvalidLocalTime,
				"time.Parse rejects this before any range check can")
		})
	}
}

func TestMealRequirement_Validate_FailsViaMinutesSinceMidnight(t *testing.T) {
	validBudget := Money{AmountMinorUnits: 1000, Currency: "RUB"}

	t.Run("window start with impossible minutes", func(t *testing.T) {
		meal := MealRequirement{
			Type:     MealBreakfast,
			Window:   TimeWindow{Start: "09:xx", End: "10:00"},
			Duration: 30 * time.Minute,
			Budget:   validBudget,
		}

		require.ErrorIs(t, meal.Validate(), core_errors.ErrInvalidLocalTime)
	})

	t.Run("window end with impossible minutes", func(t *testing.T) {
		meal := MealRequirement{
			Type:     MealBreakfast,
			Window:   TimeWindow{Start: "09:00", End: "10:yy"},
			Duration: 30 * time.Minute,
			Budget:   validBudget,
		}

		require.ErrorIs(t, meal.Validate(), core_errors.ErrInvalidLocalTime)
	})
}

func TestTimeWindow_Validate_FailsViaMinutesSinceMidnight(t *testing.T) {
	t.Run("start unparsable after passing Validate is impossible", func(t *testing.T) {
		window := TimeWindow{Start: "q9:00", End: "10:00"}

		require.ErrorIs(t, window.Validate(), core_errors.ErrInvalidLocalTime)
	})

	t.Run("end unparsable after passing Validate is impossible", func(t *testing.T) {
		window := TimeWindow{Start: "09:00", End: "1!:00"}

		require.ErrorIs(t, window.Validate(), core_errors.ErrInvalidLocalTime)
	})
}

func TestMealRequirement_Validate_DurationJustAboveWindow(t *testing.T) {
	tests := []struct {
		name     string
		start    LocalTime
		end      LocalTime
		duration time.Duration
		wantErr  error
	}{
		{name: "one minute over a 60 minute window", start: "09:00", end: "10:00", duration: 61 * time.Minute, wantErr: core_errors.ErrInvalidMealTimeEndTime},
		{name: "one second over a 60 minute window", start: "09:00", end: "10:00", duration: time.Hour + time.Second, wantErr: core_errors.ErrInvalidMealTimeEndTime},
		{name: "exactly the window passes", start: "09:00", end: "10:00", duration: time.Hour},
		{name: "one second under the window passes", start: "09:00", end: "10:00", duration: time.Hour - time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meal := MealRequirement{
				Type:     MealLunch,
				Window:   TimeWindow{Start: tt.start, End: tt.end},
				Duration: tt.duration,
				Budget:   Money{AmountMinorUnits: 1000, Currency: "RUB"},
			}

			err := meal.Validate()

			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
