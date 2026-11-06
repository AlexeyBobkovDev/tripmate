package domain

import "time"

type RouteStop struct {
	ID      int
	PlaceID int

	Position    int
	Name        string
	Description string

	Coordinates Coordinates

	ArrivalTime   time.Time
	DepartureTime time.Time
	VisitDuration time.Duration
}

func (s *RouteStop) Validate() error {
	if err := validateID(s.ID); err != nil {
		return err
	}
	if err := validateID(s.PlaceID); err != nil {
		return err
	}
	if err := validateRouteStopPosition(s.Position); err != nil {
		return err
	}
	if err := validateRouteStopName(s.Name); err != nil {
		return err
	}
	if err := validateRouteStopDescription(s.Description); err != nil {
		return err
	}
	if err := s.Coordinates.Validate(); err != nil {
		return err
	}
	if err := validateArrivalTime(s.ArrivalTime); err != nil {
		return err
	}
	if err := validateDepartureTime(s.DepartureTime); err != nil {
		return err
	}
	if err := validateRouteStopTimesOrder(s.ArrivalTime, s.VisitDuration, s.DepartureTime); err != nil {
		return err
	}

	return nil
}
