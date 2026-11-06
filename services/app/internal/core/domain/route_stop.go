package domain

import "time"

type RouteStop struct {
	ID      int `json:"id"           example:"1"`
	PlaceID int `json:"place_id"     example:"1"`

	Position    int    `json:"position"         example:"1"`
	Name        string `json:"name"             example:"Name"`
	Description string `json:"description"      example:"Description"`

	Coordinates Coordinates

	ArrivalTime   time.Time     `json:"arrival_time"                           example:"2026-10-28T12:00:00Z"`
	DepartureTime time.Time     `json:"departure_time"                         example:"2026-10-28T14:00:00Z"`
	VisitDuration time.Duration `json:"visit_duration"  swaggertype:"integer"  example:"7200"`
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
