package domain

import (
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

type LocalTime string

func (lt *LocalTime) Validate() error {
	localTime, err := time.Parse("15:04", string(*lt))
	if err != nil {
		return core_errors.ErrInvalidLocalTime
	}

	if localTime.Hour() < 0 || localTime.Hour() > 23 {
		return core_errors.ErrInvalidLocalTimeHours
	}
	if localTime.Minute() < 0 || localTime.Minute() > 59 {
		return core_errors.ErrInvalidLocalTimeMinutes
	}
	return nil
}

func (lt *LocalTime) MinutesSinceMidnight() (int, error) {
	localTime, err := time.Parse("15:04", string(*lt))
	if err != nil {
		return 0, core_errors.ErrInvalidLocalTime
	}
	return localTime.Hour()*60 + localTime.Minute(), nil
}
