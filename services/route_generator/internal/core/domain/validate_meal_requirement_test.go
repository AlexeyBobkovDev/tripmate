package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func TestValidateMealType(t *testing.T) {
	tests := []struct {
		name     string
		mealType MealType
		wantErr  error
	}{
		{name: "breakfast", mealType: MealBreakfast},
		{name: "lunch", mealType: MealLunch},
		{name: "dinner", mealType: MealDinner},
		{name: "empty", mealType: "", wantErr: core_errors.ErrUnsupportedMealType},
		{name: "unknown", mealType: "brunch", wantErr: core_errors.ErrUnsupportedMealType},
		{name: "uppercase", mealType: "Breakfast", wantErr: core_errors.ErrUnsupportedMealType},
		{name: "leading whitespace", mealType: " breakfast", wantErr: core_errors.ErrUnsupportedMealType},
		{name: "trailing whitespace", mealType: "dinner ", wantErr: core_errors.ErrUnsupportedMealType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMealType(tt.mealType)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
