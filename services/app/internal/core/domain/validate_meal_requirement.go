package domain

import core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"

func validateMealType(mealType MealType) error {
	switch mealType {
	case MealBreakfast:
	case MealLunch:
	case MealDinner:
	default:
		return core_errors.ErrUnsupportedMealType
	}
	return nil
}
