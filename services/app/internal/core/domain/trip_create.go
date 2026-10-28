package domain

import (
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type TripCreate struct {
	OwnerID int

	Name          string
	PlanningArea  string
	StartLocation Location

	StartDate time.Time
	EndDate   time.Time

	Budget    Money
	Interests []Interest
	TimeZone  string
}

func (t *TripCreate) Validate() error {
	if err := validateID(t.OwnerID); err != nil {
		return err
	}
	if err := validateTripName(t.Name); err != nil {
		return err
	}
	if err := validatePlanningArea(t.PlanningArea); err != nil {
		return err
	}
	if err := validateStartLocation(t.StartLocation); err != nil {
		return err
	}
	if err := validateStartDate(t.StartDate); err != nil {
		return err
	}
	if t.StartDate.Before(time.Now()) {
		return core_errors.ErrInvalidTripStartDate
	}
	if err := validateEndDate(t.EndDate); err != nil {
		return err
	}
	if t.EndDate.Before(time.Now()) {
		return core_errors.ErrInvalidTripEndDate
	}
	if !t.StartDate.Before(t.EndDate) {
		return core_errors.ErrInvalidTripEndDate
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
	ownerID int,
	name string,
	planningArea string,
	startLocation Location,
	startDate time.Time,
	endDate time.Time,
	budget Money,
	interests []Interest,
	timeZone string,
) *TripCreate {
	return &TripCreate{
		OwnerID:       ownerID,
		Name:          name,
		PlanningArea:  planningArea,
		StartLocation: startLocation,
		StartDate:     startDate,
		EndDate:       endDate,
		Budget:        budget,
		Interests:     interests,
		TimeZone:      timeZone,
	}
}
