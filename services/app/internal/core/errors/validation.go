package core_errors

import "errors"

var (
	ErrInvalidID         = NewError(errors.New("invalid id"), CategoryBadRequest)
	ErrInvalidVersion    = NewError(errors.New("invalid version"), CategoryBadRequest)
	ErrDifferentVersions = NewError(errors.New("invalid patch version"), CategoryBadRequest)
	ErrInvalidCreatedAt  = NewError(errors.New("invalid createdAt"), CategoryBadRequest)
	ErrInvalidUpdatedAt  = NewError(errors.New("invalid updatedAt"), CategoryBadRequest)
	ErrInvalidDeletedAt  = NewError(errors.New("invalid DeletedAt"), CategoryBadRequest)
)

var (
	ErrInvalidUserName        = NewError(errors.New("invalid user name"), CategoryBadRequest)
	ErrInvalidUserSurname     = NewError(errors.New("invalid user surname"), CategoryBadRequest)
	ErrInvalidUserUsername    = NewError(errors.New("invalid user username"), CategoryBadRequest)
	ErrInvalidUserBirthdate   = NewError(errors.New("invalid user birthdate"), CategoryBadRequest)
	ErrInvalidUserDescription = NewError(errors.New("invalid user description"), CategoryBadRequest)
	ErrInvalidUserEmail       = NewError(errors.New("invalid user email"), CategoryBadRequest)
	ErrInvalidUserPhoneNumber = NewError(errors.New("invalid user phone number"), CategoryBadRequest)
	ErrInvalidUserPassword    = NewError(errors.New("invalid user password"), CategoryBadRequest)
)

var (
	ErrInvalidTripStartDate = NewError(errors.New("invalid trip start date"), CategoryBadRequest)
	ErrInvalidTripEndDate   = NewError(errors.New("invalid trip end date"), CategoryBadRequest)
	ErrInvalidBudget        = NewError(errors.New("invalid trip budget"), CategoryBadRequest)
	ErrUnsupportedCurrency  = NewError(errors.New("unsupported trip currency"), CategoryBadRequest)
	ErrInvalidLocationName  = NewError(errors.New("invalid trip location name"), CategoryBadRequest)
	ErrInvalidLatitude      = NewError(errors.New("invalid trip latitude"), CategoryBadRequest)
	ErrInvalidLongitude     = NewError(errors.New("invalid trip longitude"), CategoryBadRequest)
	ErrInvalidTripName      = NewError(errors.New("invalid trip name"), CategoryBadRequest)
	ErrInvalidPlanningArea  = NewError(errors.New("invalid planning area"), CategoryBadRequest)
	ErrInvalidTimeZone      = NewError(errors.New("invalid time zone"), CategoryBadRequest)
	ErrUnsupportedTimezone  = NewError(errors.New("unsupported time zone"), CategoryBadRequest)
	ErrInvalidInterest      = NewError(errors.New("invalid interest"), CategoryBadRequest)
	ErrUnsupportedInterest  = NewError(errors.New("unsupported interest"), CategoryBadRequest)
)
