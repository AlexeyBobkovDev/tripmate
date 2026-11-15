package core_errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommonSentinelCategories(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		category Category
	}{
		{name: "invalid argument", err: ErrInvalidArgument, category: CategoryBadRequest},
		{name: "bad request", err: ErrBadRequest, category: CategoryBadRequest},
		{name: "not found", err: ErrNotFound, category: CategoryNotFound},
		{name: "conflict", err: ErrConflict, category: CategoryConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var codedError *CodedError

			require.ErrorAs(t, tt.err, &codedError)
			assert.Equal(t, tt.category, codedError.Category)
			assert.NotEmpty(t, tt.err.Error())
			require.ErrorIs(t, tt.err, codedError.Err)
		})
	}
}

func TestValidationSentinelsAreBadRequests(t *testing.T) {
	sentinels := []error{
		ErrInvalidBirthDateFormat,

		ErrInvalidID,
		ErrInvalidVersion,
		ErrDifferentVersions,
		ErrInvalidCreatedAt,
		ErrInvalidUpdatedAt,
		ErrInvalidDeletedAt,

		ErrInvalidUserName,
		ErrInvalidUserSurname,
		ErrInvalidUserUsername,
		ErrInvalidUserBirthdate,
		ErrInvalidUserDescription,
		ErrInvalidUserEmail,
		ErrInvalidUserPhoneNumber,
		ErrInvalidUserPassword,

		ErrInvalidTripStartDate,
		ErrInvalidTripEndDate,
		ErrInvalidTripStartTime,
		ErrInvalidTripEndTime,
		ErrInvalidBudget,
		ErrUnsupportedCurrency,
		ErrInvalidLatitude,
		ErrInvalidLongitude,
		ErrInvalidTripName,
		ErrInvalidPlanningArea,
		ErrInvalidTimeZone,
		ErrUnsupportedTimezone,
		ErrInvalidInterest,
		ErrUnsupportedInterest,

		ErrUnsupportedRouteTravelMode,
		ErrInvalidRouteDistance,
		ErrInvalidRouteDuration,
		ErrInvalidRouteDaysAmount,
		ErrInvalidRouteDayStopsAmount,
		ErrInvalidRouteDayLegsAmount,
		ErrUnsupportedRouteStatus,
		ErrInvalidRouteStopPosition,
		ErrInvalidRouteStopName,
		ErrInvalidRouteStopDescription,
		ErrInvalidRouteArrivalTime,
		ErrInvalidRouteDepartureTime,
		ErrInvalidRouteStopDatesOrder,
		ErrInvalidRouteDayDate,
		ErrInvalidRouteDayDayNumber,
		ErrInvalidRouteLegStopReference,
		ErrDuplicateRouteStopID,

		ErrInvalidLocalTimeHours,
		ErrInvalidLocalTimeMinutes,

		ErrUnsupportedMealType,
		ErrInvalidMealTimeDuration,
		ErrInvalidMealTimeWindow,
		ErrInvalidMealTimeEndTime,
		ErrInvalidLocalTime,
		ErrInvalidDailyScheduleStartTime,
		ErrInvalidDailyScheduleEndTime,
	}

	seen := make(map[string]struct{}, len(sentinels))

	for _, sentinel := range sentinels {
		t.Run(sentinel.Error(), func(t *testing.T) {
			_, duplicate := seen[sentinel.Error()]
			assert.False(t, duplicate, "sentinel messages must be distinguishable")
			seen[sentinel.Error()] = struct{}{}

			var codedError *CodedError

			require.ErrorAs(t, sentinel, &codedError)
			assert.Equal(t, CategoryBadRequest, codedError.Category)
			assert.NotEmpty(t, sentinel.Error())
		})
	}
}

func TestNewError_UnwrapsAndPreservesCause(t *testing.T) {
	cause := errors.New("cause")

	err := NewError(cause, CategoryDatabase)

	require.ErrorIs(t, err, cause)

	var codedError *CodedError
	require.ErrorAs(t, err, &codedError)
	assert.Equal(t, CategoryDatabase, codedError.Category)
	assert.Equal(t, "cause", err.Error())
}
