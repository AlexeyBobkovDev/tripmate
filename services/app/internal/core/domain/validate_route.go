package domain

import (
	"time"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func validateRouteTravelMode(travelMode TravelMode) error {
	switch travelMode {
	case TravelModeWalking,
		TravelModeCycling,
		TravelModeDriving,
		TravelModeTransit:
		return nil
	default:
		return core_errors.ErrUnsupportedRouteTravelMode
	}
}

func validateRouteDistance(distance int) error {
	if distance <= 0 {
		return core_errors.ErrInvalidRouteDistance
	}
	return nil
}

func validateRouteDuration(duration time.Duration) error {
	if duration <= 0 {
		return core_errors.ErrInvalidRouteDuration
	}
	return nil
}

func validateRouteStatus(tripStatus RouteStatus) error {
	switch tripStatus {
	case RouteStatusGenerated,
		RouteStatusEdited:
		return nil
	default:
		return core_errors.ErrUnsupportedRouteStatus
	}
}

func validateRouteDays(days []RouteDay) error {
	if len(days) == 0 {
		return core_errors.ErrInvalidRouteDaysAmount
	}
	for _, day := range days {
		if err := day.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func validateRouteStopPosition(position int) error {
	if position <= 0 {
		return core_errors.ErrInvalidRouteStopPosition
	}
	return nil
}

func validateRouteStopName(name string) error {
	if len([]rune(name)) == 0 {
		return core_errors.ErrInvalidRouteStopName
	}
	return nil
}

func validateRouteStopDescription(description string) error {
	if len([]rune(description)) == 0 {
		return core_errors.ErrInvalidRouteStopDescription
	}
	return nil
}

func validateArrivalTime(arrivalTime time.Time) error {
	if arrivalTime.IsZero() {
		return core_errors.ErrInvalidRouteArrivalTime
	}
	return nil
}

func validateDepartureTime(departureTime time.Time) error {
	if departureTime.IsZero() {
		return core_errors.ErrInvalidRouteDepartureTime
	}
	return nil
}

func validateRouteStopTimesOrder(arrivalTime time.Time, visitDuration time.Duration, departureTime time.Time) error {
	if visitDuration < 0 {
		return core_errors.ErrInvalidRouteStopDatesOrder
	}
	if arrivalTime.After(departureTime) {
		return core_errors.ErrInvalidRouteStopDatesOrder
	}
	if arrivalTime.Add(visitDuration).After(departureTime) {
		return core_errors.ErrInvalidRouteStopDatesOrder
	}
	return nil
}

func validateRouteDayDate(date time.Time) error {
	if date.IsZero() {
		return core_errors.ErrInvalidRouteDayDate
	}
	if date.Before(time.Now().UTC()) {
		return core_errors.ErrInvalidRouteDayDate
	}
	return nil
}

func validateRouteDayDayNumber(dayNumber int) error {
	if dayNumber <= 0 {
		return core_errors.ErrInvalidRouteDayDayNumber
	}
	return nil
}

func validateStopsAndLegs(stops []RouteStop, legs []RouteLeg) error {
	existingStops := map[int]struct{}{}

	for _, stop := range stops {
		if _, ok := existingStops[stop.ID]; ok {
			return core_errors.ErrDuplicateRouteStopID
		}
		existingStops[stop.ID] = struct{}{}
	}

	for _, leg := range legs {
		if leg.FromStopID == leg.ToStopID {
			return core_errors.ErrInvalidRouteLegStopReference
		}
		if _, ok := existingStops[leg.FromStopID]; !ok {
			return core_errors.ErrInvalidRouteLegStopReference
		}
		if _, ok := existingStops[leg.ToStopID]; !ok {
			return core_errors.ErrInvalidRouteLegStopReference
		}
	}
	return nil
}
