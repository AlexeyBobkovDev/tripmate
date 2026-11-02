package domain

import (
	"time"
)

type Trip struct {
	ID      int
	OwnerID int
	Version int

	Name             string
	PlanningArea     string
	StartCoordinates Coordinates

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

func NewTrip(
	id int,
	ownerID int,
	version int,
	name string,
	planningArea string,
	startCoordinates Coordinates,
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
		ID:               id,
		OwnerID:          ownerID,
		Version:          version,
		Name:             name,
		PlanningArea:     planningArea,
		StartCoordinates: startCoordinates,
		StartDate:        startDate,
		EndDate:          endDate,
		Budget:           budget,
		Interests:        interests,
		TimeZone:         timeZone,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		DeletedAt:        deletedAt,
	}
}
