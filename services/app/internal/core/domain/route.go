package domain

import "time"

type Route struct {
	ID     int
	TripID int

	TripVersion int
	Version     int

	Status RouteStatus
	Days   []RouteDay

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Route) Validate() error {
	if err := validateID(r.ID); err != nil {
		return err
	}
	if err := validateID(r.TripID); err != nil {
		return err
	}
	if err := validateVersion(r.TripVersion); err != nil {
		return err
	}
	if err := validateVersion(r.Version); err != nil {
		return err
	}
	if err := validateRouteStatus(r.Status); err != nil {
		return err
	}
	if err := validateRouteDays(r.Days); err != nil {
		return err
	}

	if err := validateCreatedAt(r.CreatedAt); err != nil {
		return err
	}
	if err := validateUpdatedAt(r.UpdatedAt); err != nil {
		return err
	}
	if err := validateSystemDatesOrder(r.CreatedAt, r.UpdatedAt, nil); err != nil {
		return err
	}

	return nil
}
