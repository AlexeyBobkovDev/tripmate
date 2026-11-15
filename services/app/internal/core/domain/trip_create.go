package domain

import (
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type TripCreate struct {
	OwnerID int

	Name             string
	PlanningArea     string
	StartCoordinates Coordinates

	StartDate time.Time
	EndDate   time.Time

	DailySchedule    DailySchedule
	MealRequirements []MealRequirement

	Budget    Money
	Interests []Interest
	TimeZone  string
}

func (t *TripCreate) Validate() error {
	if err := validateTripName(t.Name); err != nil {
		return err
	}
	if err := validatePlanningArea(t.PlanningArea); err != nil {
		return err
	}
	if err := validateCoordinates(t.StartCoordinates); err != nil {
		return err
	}

	if err := validateTripStartDate(t.StartDate); err != nil {
		return err
	}
	if t.StartDate.Before(time.Now()) {
		return core_errors.ErrInvalidTripStartDate
	}
	if err := validateTripEndDate(t.EndDate); err != nil {
		return err
	}
	if t.EndDate.Before(time.Now()) {
		return core_errors.ErrInvalidTripEndDate
	}
	if !t.StartDate.Equal(t.EndDate) {
		return core_errors.ErrInvalidTripEndDate
	}

	if err := t.DailySchedule.Validate(); err != nil {
		return err
	}
	for _, mealRequirement := range t.MealRequirements {
		if err := mealRequirement.Validate(); err != nil {
			return err
		}
	}

	if err := validateBudget(t.Budget); err != nil {
		return err
	}
	if err := validateInterests(t.Interests); err != nil {
		return err
	}
	if err := validateTimeZone(t.TimeZone); err != nil {
		return err
	}

	return nil
}

func NewTripCreate(
	name string,
	planningArea string,
	startCoordinates Coordinates,
	startDate time.Time,
	endDate time.Time,
	dailySchedule DailySchedule,
	mealRequirements []MealRequirement,
	budget Money,
	interests []Interest,
	timeZone string,
) *TripCreate {
	return &TripCreate{
		Name:             name,
		PlanningArea:     planningArea,
		StartCoordinates: startCoordinates,
		StartDate:        startDate,
		EndDate:          endDate,
		DailySchedule:    dailySchedule,
		MealRequirements: mealRequirements,
		Budget:           budget,
		Interests:        interests,
		TimeZone:         timeZone,
	}
}
