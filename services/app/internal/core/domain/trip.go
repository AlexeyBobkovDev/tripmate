package domain

import (
	"fmt"
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type Trip struct {
	ID      int
	OwnerID int
	Version int

	Name          string
	PlanningArea  string
	StartLocation Location

	StartDate time.Time
	EndDate   time.Time

	Budget    Money
	Interests []Interest
	TimeZone  string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (t *Trip) Validate() error {
	if err := validateID(t.ID); err != nil {
		return err
	}
	if err := validateID(t.OwnerID); err != nil {
		return err
	}
	if err := validateVersion(t.Version); err != nil {
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
	if err := validateEndDate(t.EndDate); err != nil {
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
	if err := validateCreatedAt(t.CreatedAt); err != nil {
		return err
	}
	if err := validateUpdatedAt(t.UpdatedAt); err != nil {
		return err
	}
	if err := validateDeletedAt(t.DeletedAt); err != nil {
		return err
	}
	if err := validateTripDatesOrder(t.StartDate, t.EndDate, t.CreatedAt, t.UpdatedAt, t.DeletedAt); err != nil {
		return err
	}

	return nil
}

func (t *Trip) ApplyPatch(patch *TripPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate trip patch: %w", err)
	}

	tmp := *t

	if tmp.Version != patch.Version {
		return core_errors.ErrDifferentVersions
	}

	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}
	if patch.PlanningArea.Set {
		tmp.PlanningArea = *patch.PlanningArea.Value
	}
	if patch.StartLocation.Set {
		tmp.StartLocation = *patch.StartLocation.Value
	}
	if patch.StartDate.Set {
		tmp.StartDate = *patch.StartDate.Value
	}
	if patch.EndDate.Set {
		tmp.EndDate = *patch.EndDate.Value
	}
	if patch.Budget.Set {
		tmp.Budget = *patch.Budget.Value
	}
	if patch.Interests.Set {
		tmp.Interests = *patch.Interests.Value
	}
	if patch.TimeZone.Set {
		tmp.TimeZone = *patch.TimeZone.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched trip: %w", err)
	}

	*t = tmp

	return nil
}

func NewTrip(
	id int,
	ownerID int,
	version int,
	name string,
	planningArea string,
	startLocation Location,
	startDate time.Time,
	endDate time.Time,
	budget Money,
	interests []Interest,
	timeZone string,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) *Trip {
	return &Trip{
		ID:            id,
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
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		DeletedAt:     deletedAt,
	}
}
