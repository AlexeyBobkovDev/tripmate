package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func TestValidateTripStartDate(t *testing.T) {
	require.ErrorIs(t, validateTripStartDate(time.Time{}), core_errors.ErrInvalidTripStartDate)
	require.NoError(t, validateTripStartDate(time.Now()))
	require.NoError(t, validateTripStartDate(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)))
}

func TestValidateTripEndDate(t *testing.T) {
	require.ErrorIs(t, validateTripEndDate(time.Time{}), core_errors.ErrInvalidTripEndDate)
	require.NoError(t, validateTripEndDate(time.Now()))
	require.NoError(t, validateTripEndDate(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)))
}

func TestValidateTripDatesOrder(t *testing.T) {
	created := time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name                             string
		start, end, createdAt, updatedAt time.Time
		deleted                          *time.Time
		want                             error
	}{
		{name: "valid future trip", start: created.Add(time.Hour), end: created.Add(time.Hour), createdAt: created, updatedAt: created},
		{name: "start before end", start: created.Add(1 * time.Hour), end: created.Add(2 * time.Hour), createdAt: created, updatedAt: created, want: core_errors.ErrInvalidTripEndDate},
		{name: "start after end", start: created.Add(2 * time.Hour), end: created.Add(1 * time.Hour), createdAt: created, updatedAt: created, want: core_errors.ErrInvalidTripEndDate},
		{name: "start before creation", start: created.Add(-time.Hour), end: created.Add(-time.Hour), createdAt: created, updatedAt: created, want: core_errors.ErrInvalidTripStartDate},
		{name: "start equals creation", start: created, end: created, createdAt: created, updatedAt: created, want: core_errors.ErrInvalidTripStartDate},
		{name: "system dates invalid", start: created.Add(time.Hour), end: created.Add(time.Hour), createdAt: created, updatedAt: created.Add(-time.Second), deleted: nil, want: core_errors.ErrInvalidUpdatedAt},
		{name: "deleted before updated", start: created.Add(2 * time.Hour), end: created.Add(2 * time.Hour), createdAt: created, updatedAt: created.Add(time.Hour), deleted: timePtr(created.Add(30 * time.Minute)), want: core_errors.ErrInvalidDeletedAt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTripDatesOrder(tt.start, tt.end, tt.createdAt, tt.updatedAt, tt.deleted)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateBudget(t *testing.T) {
	require.ErrorIs(t, validateBudget(Money{}), core_errors.ErrInvalidBudget)
	require.NoError(t, validateBudget(Money{AmountMinorUnits: 1, Currency: "RUB"}))
	require.ErrorIs(t, validateBudget(Money{AmountMinorUnits: 0, Currency: "RUB"}), core_errors.ErrInvalidBudget)
	require.ErrorIs(t, validateBudget(Money{AmountMinorUnits: 100, Currency: "USD"}), core_errors.ErrUnsupportedCurrency)
}

func TestValidateTripName(t *testing.T) {
	for _, value := range []string{"Trip", "Поездка", "Trip 2026", "   Trip   "} {
		require.NoError(t, validateTripName(value), "value: %q", value)
	}
	for _, value := range []string{"", " ", "\t\n"} {
		require.ErrorIs(t, validateTripName(value), core_errors.ErrInvalidTripName, "value: %q", value)
	}
}

func TestValidateTripPlanningArea(t *testing.T) {
	for _, value := range []string{"Nizhny Novgorod", "Paris, France", "Москва"} {
		require.NoError(t, validatePlanningArea(value), "value: %q", value)
	}
	for _, value := range []string{"", " ", "\t\n"} {
		require.ErrorIs(t, validatePlanningArea(value), core_errors.ErrInvalidPlanningArea, "value: %q", value)
	}
}

func TestValidateCoordinates(t *testing.T) {
	require.ErrorIs(t, validateCoordinates(Coordinates{}), nil)
	require.NoError(t, validateCoordinates(Coordinates{Latitude: 55.75, Longitude: 37.62}))
	require.NoError(t, validateCoordinates(Coordinates{Latitude: 55.75, Longitude: 37.62}))
	require.ErrorIs(t, validateCoordinates(Coordinates{Latitude: 91, Longitude: 37.62}), core_errors.ErrInvalidLatitude)
	require.ErrorIs(t, validateCoordinates(Coordinates{Latitude: 55.75, Longitude: 181}), core_errors.ErrInvalidLongitude)
}

func TestValidateInterests(t *testing.T) {
	tests := []struct {
		name   string
		values []Interest
		want   error
	}{
		{name: "nil list"},
		{name: "empty list", values: []Interest{}},
		{name: "single supported", values: []Interest{InterestHistory}},
		{name: "multiple supported", values: []Interest{InterestHistory, InterestNature, InterestFood, InterestArchitecture}},
		{name: "blank interest", values: []Interest{" "}, want: core_errors.ErrInvalidInterest},
		{name: "unsupported interest", values: []Interest{"shopping"}, want: core_errors.ErrUnsupportedInterest},
		{name: "supported then unsupported", values: []Interest{InterestFood, "shopping"}, want: core_errors.ErrUnsupportedInterest},
		{name: "blank precedes unsupported", values: []Interest{" ", "shopping"}, want: core_errors.ErrInvalidInterest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInterests(tt.values)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateTimeZone(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  error
	}{
		{name: "supported Moscow", value: "Europe/Moscow"},
		{name: "empty", value: "", want: core_errors.ErrInvalidTimeZone},
		{name: "whitespace", value: " \t ", want: core_errors.ErrInvalidTimeZone},
		{name: "invalid IANA zone", value: "Mars/Olympus", want: core_errors.ErrUnsupportedTimezone},
		{name: "valid UTC but unsupported by product policy", value: "UTC", want: core_errors.ErrUnsupportedTimezone},
		{name: "valid New York but unsupported by product policy", value: "America/New_York", want: core_errors.ErrUnsupportedTimezone},
		{name: "Moscow with whitespace currently rejected", value: " Europe/Moscow ", want: core_errors.ErrUnsupportedTimezone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTimeZone(tt.value)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
