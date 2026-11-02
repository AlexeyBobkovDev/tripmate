package domain

import (
	"strings"
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func validateTripStartDate(startDate time.Time) error {
	if startDate.IsZero() {
		return core_errors.ErrInvalidTripStartDate
	}
	return nil
}

func validateTripEndDate(endDate time.Time) error {
	if endDate.IsZero() {
		return core_errors.ErrInvalidTripEndDate
	}
	return nil
}

func validateTripDatesOrder(
	startDate,
	endDate time.Time,
	createdAt,
	updatedAt time.Time,
	deletedAt *time.Time,
) error {
	if !startDate.Equal(endDate) {
		return core_errors.ErrInvalidTripEndDate
	}

	if startDate.Before(createdAt) || startDate.Equal(createdAt) {
		return core_errors.ErrInvalidTripStartDate
	}

	return validateSystemDatesOrder(createdAt, updatedAt, deletedAt)
}

func validateBudget(budget Money) error {
	return budget.Validate()
}

func validateTripName(name string) error {
	if len([]rune(strings.TrimSpace(name))) == 0 {
		return core_errors.ErrInvalidTripName
	}
	return nil
}

func validatePlanningArea(planningArea string) error {
	if len([]rune(strings.TrimSpace(planningArea))) == 0 {
		return core_errors.ErrInvalidPlanningArea
	}
	return nil
}

func validateCoordinates(coordinates Coordinates) error {
	return coordinates.Validate()
}

func validateInterests(interests []Interest) error {
	for _, interest := range interests {
		if len([]rune(strings.TrimSpace(string(interest)))) == 0 {
			return core_errors.ErrInvalidInterest
		}
		if _, ok := AllowedInterests[interest]; !ok {
			return core_errors.ErrUnsupportedInterest
		}
	}
	return nil
}

func validateTimeZone(timeZone string) error {
	if len([]rune(strings.TrimSpace(timeZone))) == 0 {
		return core_errors.ErrInvalidTimeZone
	}
	if _, err := time.LoadLocation(timeZone); err != nil {
		return core_errors.ErrUnsupportedTimezone
	}
	if timeZone != "Europe/Moscow" {
		return core_errors.ErrUnsupportedTimezone
	}
	return nil
}
