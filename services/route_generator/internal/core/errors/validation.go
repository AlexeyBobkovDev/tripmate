package core_errors

import "errors"

var (
	ErrInvalidID         = NewError(errors.New("invalid id"), CategoryBadRequest)
	ErrDifferentVersions = NewError(errors.New("invalid patch version"), CategoryBadRequest)
)

// TODO: refactor Unsupported and Invalid
var (
	ErrInvalidTripStartDate = NewError(errors.New("invalid trip start date"), CategoryBadRequest)
	ErrInvalidTripEndDate   = NewError(errors.New("invalid trip end date"), CategoryBadRequest)
	ErrInvalidBudget        = NewError(errors.New("invalid trip budget"), CategoryBadRequest)
	ErrUnsupportedCurrency  = NewError(errors.New("unsupported trip currency"), CategoryBadRequest)
	ErrInvalidLatitude      = NewError(errors.New("invalid trip latitude"), CategoryBadRequest)
	ErrInvalidLongitude     = NewError(errors.New("invalid trip longitude"), CategoryBadRequest)
	ErrInvalidTripName      = NewError(errors.New("invalid trip name"), CategoryBadRequest)
	ErrInvalidPlanningArea  = NewError(errors.New("invalid planning area"), CategoryBadRequest)
	ErrInvalidTimeZone      = NewError(errors.New("invalid time zone"), CategoryBadRequest)
	ErrUnsupportedTimezone  = NewError(errors.New("unsupported time zone"), CategoryBadRequest)
	ErrInvalidInterest      = NewError(errors.New("invalid interest"), CategoryBadRequest)
	ErrUnsupportedInterest  = NewError(errors.New("unsupported interest"), CategoryBadRequest)
)

var (
	ErrUnsupportedRouteTravelMode  = NewError(errors.New("unsupported route travel mode"), CategoryBadRequest)
	ErrInvalidRouteDistance        = NewError(errors.New("invalid route distance"), CategoryBadRequest)
	ErrInvalidRouteDuration        = NewError(errors.New("invalid route duration"), CategoryBadRequest)
	ErrInvalidRouteDaysAmount      = NewError(errors.New("invalid route days amount"), CategoryBadRequest)
	ErrInvalidRouteDayStopsAmount  = NewError(errors.New("invalid route day stops amount"), CategoryBadRequest)
	ErrInvalidRouteDayLegsAmount   = NewError(errors.New("invalid route day legs amount"), CategoryBadRequest)
	ErrUnsupportedRouteStatus      = NewError(errors.New("unsupported route status"), CategoryBadRequest)
	ErrInvalidRouteStopPosition    = NewError(errors.New("invalid route stop position"), CategoryBadRequest)
	ErrInvalidRouteStopName        = NewError(errors.New("invalid route stop name"), CategoryBadRequest)
	ErrInvalidRouteStopDescription = NewError(errors.New("invalid route stop description"), CategoryBadRequest)
	ErrInvalidRouteArrivalTime     = NewError(errors.New("invalid route stop arrival time"), CategoryBadRequest)
	ErrInvalidRouteDepartureTime   = NewError(errors.New("invalid route stop departure time"), CategoryBadRequest)
	ErrInvalidRouteStopDatesOrder  = NewError(errors.New("invalid route stop dates order"), CategoryBadRequest)
	ErrInvalidRouteDayDate         = NewError(errors.New("invalid route day date"), CategoryBadRequest)
	ErrInvalidRouteDayDayNumber    = NewError(errors.New("invalid route day day number"), CategoryBadRequest)
)
