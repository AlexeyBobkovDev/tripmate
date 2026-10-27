package domain

import (
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type TripPatch struct {
	OwnerID int
	Version int

	Name          Nullable[string]
	PlanningArea  Nullable[string]
	StartLocation Nullable[*Location]

	StartDate Nullable[time.Time]
	EndDate   Nullable[time.Time]

	Budget    Nullable[*Money]
	Interests Nullable[[]Interest]
	TimeZone  Nullable[string]
}

func (t *TripPatch) Validate() error {
	if err := validateID(t.OwnerID); err != nil {
		return err
	}
	if err := validateVersion(t.Version); err != nil {
		return err
	}
	if t.Name.Set {
		if t.Name.Value == nil {
			return core_errors.ErrInvalidArgument
		}
		if err := validateTripName(*t.Name.Value); err != nil {
			return err
		}
	}
	if t.PlanningArea.Set {
		if t.PlanningArea.Value == nil {
			return core_errors.ErrInvalidArgument
		}
		if err := validatePlanningArea(*t.PlanningArea.Value); err != nil {
			return err
		}
	}
	if t.StartLocation.Set {
		if t.StartLocation.Value == nil {
			return core_errors.ErrInvalidArgument
		}
		if err := validateStartLocation(*t.StartLocation.Value); err != nil {
			return err
		}
	}
	if t.StartDate.Set {
		if t.StartDate.Value == nil {
			return core_errors.ErrInvalidArgument
		}
		if err := validateStartDate(*t.StartDate.Value); err != nil {
			return err
		}
	}
	if t.EndDate.Set {
		if t.EndDate.Value == nil {
			return core_errors.ErrInvalidArgument
		}
		if err := validateEndDate(*t.EndDate.Value); err != nil {
			return err
		}
	}
	if t.Budget.Set {
		if t.Budget.Value == nil {
			return core_errors.ErrInvalidArgument
		}
		if err := validateBudget(*t.Budget.Value); err != nil {
			return err
		}
	}
	if t.Interests.Set {
		if t.Interests.Value == nil {
			return core_errors.ErrInvalidArgument
		}
		if err := validateInterests(*t.Interests.Value); err != nil {
			return err
		}
	}
	if t.TimeZone.Set {
		if t.TimeZone.Value == nil {
			return core_errors.ErrInvalidArgument
		}
		if err := validateTimeZone(*t.TimeZone.Value); err != nil {
			return err
		}
	}

	return nil
}

func NewTripPatch(
	ownerID int,
	version int,
	name Nullable[string],
	planningArea Nullable[string],
	startLocation Nullable[*Location],
	startDate Nullable[time.Time],
	endDate Nullable[time.Time],
	budget Nullable[*Money],
	interests Nullable[[]Interest],
	timeZone Nullable[string],
) *TripPatch {
	return &TripPatch{
		OwnerID:       ownerID,
		Version:       version,
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
