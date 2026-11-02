package domain

import "time"

type RouteLeg struct {
	ID int

	FromStopID int
	ToStopID   int

	TravelMode TravelMode
	Distance   int
	Duration   time.Duration
}

func (l *RouteLeg) Validate() error {
	if err := validateID(l.ID); err != nil {
		return err
	}
	if err := validateID(l.FromStopID); err != nil {
		return err
	}
	if err := validateID(l.ToStopID); err != nil {
		return err
	}
	if err := validateRouteTravelMode(l.TravelMode); err != nil {
		return err
	}
	if err := validateRouteDistance(l.Distance); err != nil {
		return err
	}
	if err := validateRouteDuration(l.Duration); err != nil {
		return err
	}

	return nil
}
