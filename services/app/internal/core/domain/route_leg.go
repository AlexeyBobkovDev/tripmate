package domain

import "time"

type RouteLeg struct {
	ID int `json:"id"    example:"1"`

	FromStopID int `json:"from_stop_id"   example:"1"`
	ToStopID   int `json:"to_stop_id"     example:"1"`

	TravelMode TravelMode    `json:"travel_mode"                              example:"walking"`
	Distance   int           `json:"distance"                                 example:"100"`
	Duration   time.Duration `json:"duration"      swaggertype:"integer"      example:"800"`
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
