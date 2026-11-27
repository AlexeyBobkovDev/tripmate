package domain

import (
	"time"
)

type RouteStop struct {
	ID      int `json:"id"           example:"1"`
	PlaceID int `json:"place_id"     example:"1"`

	Position int `json:"position"         example:"1"`

	Name                 string     `json:"name"            example:"Name"`
	Description          string     `json:"description"     example:"Description"`
	Categories           []Category `json:"categories"`
	Interests            []Interest `json:"interests"       example:"[history]"`
	Address              string     `json:"address"         example:"Russia, Nizhny Novgorod, Upper Volga Embankment, 7"`
	WebsiteURL           *string    `json:"website_url"`
	Phone                *string    `json:"phone"`
	VisitDurationMinutes int        `json:"visit_duration_minutes"   example:"60"`

	OpeningHours []OpeningHours `json:"opening_hours"`

	Coordinates Coordinates `json:"coordinates"`

	ArrivalTime   time.Time `json:"arrival_time"          example:"2026-10-28T12:00:00Z"`
	DepartureTime time.Time `json:"departure_time"        example:"2026-10-28T14:00:00Z"`
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
	if err := validateCategories(s.Categories); err != nil {
		return err
	}
	if err := validateInterests(s.Interests); err != nil {
		return err
	}
	if err := validateRouteStopAddress(s.Address); err != nil {
		return err
	}
	if err := validateRouteStopWebsiteURL(*s.WebsiteURL); err != nil {
		return err
	}
	if s.Phone != nil {
		if err := validatePhoneNumber(*s.Phone); err != nil {
			return err
		}
	}
	if err := validateRouteStopVisitDurationMinutes(s.VisitDurationMinutes); err != nil {
		return err
	}
	if err := validateRouteStopOpeningHours(s.OpeningHours); err != nil {
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
	if err := validateRouteStopTimesOrder(s.ArrivalTime, s.VisitDurationMinutes, s.DepartureTime); err != nil {
		return err
	}

	return nil
}
