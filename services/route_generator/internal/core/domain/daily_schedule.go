package domain

import (
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

type DailySchedule struct {
	StartTime LocalTime `json:"start_time" example:"09:00"`
	EndTime   LocalTime `json:"end_time"   example:"18:00"`
}

func (ds *DailySchedule) Validate() error {
	if err := ds.StartTime.Validate(); err != nil {
		return err
	}
	if err := ds.EndTime.Validate(); err != nil {
		return err
	}

	startMinutes, err := ds.StartTime.MinutesSinceMidnight()
	if err != nil {
		return err
	}

	endMinutes, err := ds.EndTime.MinutesSinceMidnight()
	if err != nil {
		return err
	}

	if startMinutes > endMinutes {
		return core_errors.ErrInvalidTripEndTime
	}

	return nil
}
