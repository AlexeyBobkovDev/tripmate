package trips_transport_http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func TestMealRequirementRequestDTO_ToDomain(t *testing.T) {
	tests := []struct {
		name             string
		mealType         domain.MealType
		durationMinutes  int
		expectedDuration time.Duration
	}{
		{name: "breakfast 30 minutes", mealType: domain.MealBreakfast, durationMinutes: 30, expectedDuration: 30 * time.Minute},
		{name: "lunch one minute", mealType: domain.MealLunch, durationMinutes: 1, expectedDuration: time.Minute},
		{name: "dinner full day", mealType: domain.MealDinner, durationMinutes: 1440, expectedDuration: 24 * time.Hour},
		{name: "zero minutes", mealType: domain.MealLunch, durationMinutes: 0, expectedDuration: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := MealRequirementRequestDTO{
				Type:     tt.mealType,
				Window:   domain.TimeWindow{Start: "12:00", End: "14:00"},
				Duration: tt.durationMinutes,
				Budget:   domain.Money{AmountMinorUnits: 150000, Currency: "RUB"},
			}

			meal := dto.ToDomain()

			assert.Equal(t, tt.mealType, meal.Type)
			assert.Equal(t, dto.Window, meal.Window)
			assert.Equal(t, tt.expectedDuration, meal.Duration)
			assert.Equal(t, dto.Budget, meal.Budget)
		})
	}
}

func TestMealRequirementRequestDTO_ToDomain_DoesNotMutateDTO(t *testing.T) {
	dto := MealRequirementRequestDTO{Type: domain.MealBreakfast, Duration: 45}

	_ = dto.ToDomain()
	_ = dto.ToDomain()

	assert.Equal(t, 45, dto.Duration)
}
