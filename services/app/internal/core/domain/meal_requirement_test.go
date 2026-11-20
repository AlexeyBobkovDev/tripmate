package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestMealRequirement_Validate(t *testing.T) {
	tests := []struct {
		name        string
		requirement MealRequirement
		wantErr     error
	}{
		{
			name: "valid breakfast requirement",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
		},
		{
			name: "valid lunch requirement",
			requirement: MealRequirement{
				Type:     MealLunch,
				Window:   TimeWindow{Start: "12:00", End: "14:00"},
				Duration: 2 * time.Hour,
				Budget:   Money{AmountMinorUnits: 3000000, Currency: "RUB"},
			},
		},
		{
			name: "valid dinner requirement",
			requirement: MealRequirement{
				Type:     MealDinner,
				Window:   TimeWindow{Start: "19:00", End: "21:00"},
				Duration: 2 * time.Hour,
				Budget:   Money{AmountMinorUnits: 3500000, Currency: "RUB"},
			},
		},
		{
			name: "invalid meal type",
			requirement: MealRequirement{
				Type:     "invalid",
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
			wantErr: core_errors.ErrUnsupportedMealType,
		},
		{
			name: "invalid window start time",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "invalid", End: "10:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid window end time",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "invalid"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid duration - zero",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: 0,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidMealTimeDuration,
		},
		{
			name: "invalid duration - negative",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: -time.Minute,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidMealTimeDuration,
		},
		{
			name: "duration exceeds window",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: 2 * time.Hour,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidMealTimeEndTime,
		},
		{
			name: "duration equals window",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
		},
		{
			name: "duration starts at window end",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: 0,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidMealTimeDuration,
		},
		{
			name: "invalid budget",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 0, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "invalid budget - negative",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: -100, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "invalid budget - wrong currency",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "10:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "USD"},
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "window start equals window end",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "09:00", End: "09:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidMealTimeEndTime,
		},
		{
			name: "window start after window end",
			requirement: MealRequirement{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "10:00", End: "09:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			},
			wantErr: core_errors.ErrInvalidTripEndTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.requirement.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestTimeWindow_Validate(t *testing.T) {
	tests := []struct {
		name    string
		window  TimeWindow
		wantErr error
	}{
		{
			name: "valid time window",
			window: TimeWindow{
				Start: "09:00",
				End:   "18:00",
			},
		},
		{
			name: "valid window - midnight to midnight",
			window: TimeWindow{
				Start: "00:00",
				End:   "00:00",
			},
		},
		{
			name: "valid window - all day",
			window: TimeWindow{
				Start: "00:00",
				End:   "23:59",
			},
		},
		{
			name: "invalid start time",
			window: TimeWindow{
				Start: "invalid",
				End:   "18:00",
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid end time",
			window: TimeWindow{
				Start: "09:00",
				End:   "invalid",
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "start time after end time",
			window: TimeWindow{
				Start: "18:00",
				End:   "09:00",
			},
			wantErr: core_errors.ErrInvalidTripEndTime,
		},
		{
			name: "start time equals end time",
			window: TimeWindow{
				Start: "12:00",
				End:   "12:00",
			},
		},
		{
			name: "invalid start time - hours > 23",
			window: TimeWindow{
				Start: "24:00",
				End:   "18:00",
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid end time - minutes > 59",
			window: TimeWindow{
				Start: "09:00",
				End:   "18:60",
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.window.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestMealRequirement_Validate_EdgeCases(t *testing.T) {
	validMeal := MealRequirement{
		Type:     MealBreakfast,
		Window:   TimeWindow{Start: "08:00", End: "09:00"},
		Duration: 30 * time.Minute,
		Budget:   Money{AmountMinorUnits: 1000, Currency: "RUB"},
	}

	t.Run("zero duration", func(t *testing.T) {
		meal := validMeal
		meal.Duration = 0
		err := meal.Validate()
		require.ErrorIs(t, err, core_errors.ErrInvalidMealTimeDuration)
	})

	t.Run("duration longer than window", func(t *testing.T) {
		meal := validMeal
		meal.Duration = 120 * time.Minute
		err := meal.Validate()
		require.ErrorIs(t, err, core_errors.ErrInvalidMealTimeEndTime)
	})

	t.Run("invalid budget", func(t *testing.T) {
		meal := validMeal
		meal.Budget = Money{AmountMinorUnits: 0, Currency: "RUB"}
		err := meal.Validate()
		require.ErrorIs(t, err, core_errors.ErrInvalidBudget)
	})

	t.Run("unsupported meal type", func(t *testing.T) {
		meal := validMeal
		meal.Type = MealType("unknown")
		err := meal.Validate()
		require.ErrorIs(t, err, core_errors.ErrUnsupportedMealType)
	})

	t.Run("start after end window", func(t *testing.T) {
		meal := validMeal
		meal.Window = TimeWindow{Start: "12:00", End: "11:00"}
		err := meal.Validate()
		require.ErrorIs(t, err, core_errors.ErrInvalidTripEndTime)
	})
}
