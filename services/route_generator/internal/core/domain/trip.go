package domain

import (
	"time"
)

type Trip struct {
	Name             string
	PlanningArea     string
	StartCoordinates Coordinates

	StartDate time.Time
	EndDate   time.Time

	Budget    Money
	Interests []Interest
	TimeZone  string
}

func (t *Trip) Validate() error {
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
	if err := validateTripEndDate(t.EndDate); err != nil {
		return err
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
	if err := validateTripDatesOrder(t.StartDate, t.EndDate); err != nil {
		return err
	}

	return nil
}

func NewTrip(
	name string,
	planningArea string,
	startCoordinates Coordinates,
	startDate time.Time,
	endDate time.Time,
	budget Money,
	interests []Interest,
	timeZone string,
) *Trip {
	return &Trip{
		Name:             name,
		PlanningArea:     planningArea,
		StartCoordinates: startCoordinates,
		StartDate:        startDate,
		EndDate:          endDate,
		Budget:           budget,
		Interests:        interests,
		TimeZone:         timeZone,
	}
}
