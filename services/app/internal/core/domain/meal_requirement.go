package domain

import (
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type MealType string

const (
	MealBreakfast MealType = "breakfast"
	MealLunch     MealType = "lunch"
	MealDinner    MealType = "dinner"
)

type MealRequirement struct {
	Type     MealType      `json:"meal_type"                               example:"breakfast"`
	Window   TimeWindow    `json:"window"`
	Duration time.Duration `json:"duration_minutes"  swaggertype:"integer" example:"60"`
	Budget   Money         `json:"budget"`
}

func (ml *MealRequirement) Validate() error {
	if err := validateMealType(ml.Type); err != nil {
		return err
	}

	if err := ml.Window.Validate(); err != nil {
		return err
	}

	if ml.Duration <= 0 {
		return core_errors.ErrInvalidMealTimeDuration
	}

	startMinutes, err := ml.Window.Start.MinutesSinceMidnight()
	if err != nil {
		return err
	}

	endMinutes, err := ml.Window.End.MinutesSinceMidnight()
	if err != nil {
		return err
	}

	windowDuration := time.Duration(endMinutes-startMinutes) * time.Minute
	if ml.Duration > windowDuration {
		return core_errors.ErrInvalidMealTimeEndTime
	}

	if err := ml.Budget.Validate(); err != nil {
		return err
	}

	return nil
}

type TimeWindow struct {
	Start LocalTime `json:"start"    example:"09:00"`
	End   LocalTime `json:"end"      example:"18:00"`
}

func (tw *TimeWindow) Validate() error {
	if err := tw.Start.Validate(); err != nil {
		return err
	}
	if err := tw.End.Validate(); err != nil {
		return err
	}
	startMinutes, err := tw.Start.MinutesSinceMidnight()
	if err != nil {
		return err
	}

	endMinutes, err := tw.End.MinutesSinceMidnight()
	if err != nil {
		return err
	}

	if startMinutes > endMinutes {
		return core_errors.ErrInvalidTripEndTime
	}
	return nil
}
