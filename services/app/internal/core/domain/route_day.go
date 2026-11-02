package domain

import (
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

type RouteDay struct {
	ID int

	Date      time.Time
	DayNumber int

	Stops []RouteStop
	Legs  []RouteLeg
}

func (d *RouteDay) Validate() error {
	if err := validateID(d.ID); err != nil {
		return err
	}
	if err := validateRouteDayDate(d.Date); err != nil {
		return err
	}
	if err := validateRouteDayDayNumber(d.DayNumber); err != nil {
		return err
	}
	if len(d.Stops) == 0 {
		return core_errors.ErrInvalidRouteDayStopsAmount
	}
	for _, stop := range d.Stops {
		if err := stop.Validate(); err != nil {
			return err
		}
	}

	if len(d.Legs) == 0 {
		return core_errors.ErrInvalidRouteDayLegsAmount
	}
	for _, leg := range d.Legs {
		if err := leg.Validate(); err != nil {
			return err
		}
	}

	return nil
}
