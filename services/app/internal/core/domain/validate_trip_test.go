package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestValidateTripStartDate(t *testing.T) {
	// Freeze time for predictable testing
	frozenTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return frozenTime }
	defer func() {
		// Restore original time.Now
		timeNow = time.Now
	}()

	tests := []struct {
		name string
		date time.Time
		want error
		desc string
	}{
		{
			name: "valid future date",
			date: frozenTime.Add(24 * time.Hour),
			want: nil,
			desc: "date should be valid when it's in the future",
		},
		{
			name: "valid current date",
			date: frozenTime,
			want: nil,
			desc: "date should be valid when it's the current time",
		},
		{
			name: "zero time",
			date: time.Time{},
			want: core_errors.ErrInvalidTripStartDate,
			desc: "date should be invalid when it's zero time",
		},
		{
			name: "past date",
			date: frozenTime.Add(-24 * time.Hour),
			want: core_errors.ErrInvalidTripStartDate,
			desc: "date should be invalid when it's in the past",
		},
		{
			name: "far future date",
			date: frozenTime.Add(365 * 24 * time.Hour),
			want: nil,
			desc: "date should be valid when it's far in the future",
		},
		{
			name: "just one second in the past",
			date: frozenTime.Add(-1 * time.Second),
			want: core_errors.ErrInvalidTripStartDate,
			desc: "date should be invalid when it's even one second in the past",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTripStartDate(tt.date)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want, "test description: %s", tt.desc)
			} else {
				require.NoError(t, err, "test description: %s", tt.desc)
			}
		})
	}
}

func TestValidateTripEndDate(t *testing.T) {
	// Freeze time for predictable testing
	frozenTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return frozenTime }
	defer func() {
		// Restore original time.Now
		timeNow = time.Now
	}()

	tests := []struct {
		name string
		date time.Time
		want error
		desc string
	}{
		{
			name: "valid future date",
			date: frozenTime.Add(24 * time.Hour),
			want: nil,
			desc: "date should be valid when it's in the future",
		},
		{
			name: "valid current date",
			date: frozenTime,
			want: nil,
			desc: "date should be valid when it's the current time",
		},
		{
			name: "zero time",
			date: time.Time{},
			want: core_errors.ErrInvalidTripEndDate,
			desc: "date should be invalid when it's zero time",
		},
		{
			name: "past date",
			date: frozenTime.Add(-24 * time.Hour),
			want: core_errors.ErrInvalidTripEndDate,
			desc: "date should be invalid when it's in the past",
		},
		{
			name: "far future date",
			date: frozenTime.Add(365 * 24 * time.Hour),
			want: nil,
			desc: "date should be valid when it's far in the future",
		},
		{
			name: "just one second in the past",
			date: frozenTime.Add(-1 * time.Second),
			want: core_errors.ErrInvalidTripEndDate,
			desc: "date should be invalid when it's even one second in the past",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTripEndDate(tt.date)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want, "test description: %s", tt.desc)
			} else {
				require.NoError(t, err, "test description: %s", tt.desc)
			}
		})
	}
}

func TestValidateTripStartTime(t *testing.T) {
	// Freeze time for predictable testing
	frozenTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return frozenTime }
	defer func() {
		// Restore original time.Now
		timeNow = time.Now
	}()

	tests := []struct {
		name string
		time time.Time
		want error
		desc string
	}{
		{
			name: "valid future time",
			time: frozenTime.Add(24 * time.Hour),
			want: nil,
			desc: "time should be valid when it's in the future",
		},
		{
			name: "valid current time",
			time: frozenTime,
			want: nil,
			desc: "time should be valid when it's the current time",
		},
		{
			name: "zero time",
			time: time.Time{},
			want: core_errors.ErrInvalidTripStartTime,
			desc: "time should be invalid when it's zero time",
		},
		{
			name: "past time",
			time: frozenTime.Add(-24 * time.Hour),
			want: core_errors.ErrInvalidTripStartTime,
			desc: "time should be invalid when it's in the past",
		},
		{
			name: "far future time",
			time: frozenTime.Add(365 * 24 * time.Hour),
			want: nil,
			desc: "time should be valid when it's far in the future",
		},
		{
			name: "just one second in the past",
			time: frozenTime.Add(-1 * time.Second),
			want: core_errors.ErrInvalidTripStartTime,
			desc: "time should be invalid when it's even one second in the past",
		},
		{
			name: "time with different timezone",
			time: frozenTime.In(time.FixedZone("UTC+3", 3*3600)),
			want: nil,
			desc: "time should be valid when it's in a different timezone but represents the same moment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTripStartTime(tt.time)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want, "test description: %s", tt.desc)
			} else {
				require.NoError(t, err, "test description: %s", tt.desc)
			}
		})
	}
}

func TestValidateTripEndTime(t *testing.T) {
	// Freeze time for predictable testing
	frozenTime := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return frozenTime }
	defer func() {
		// Restore original time.Now
		timeNow = time.Now
	}()

	tests := []struct {
		name string
		time time.Time
		want error
		desc string
	}{
		{
			name: "valid future time",
			time: frozenTime.Add(24 * time.Hour),
			want: nil,
			desc: "time should be valid when it's in the future",
		},
		{
			name: "valid current time",
			time: frozenTime,
			want: nil,
			desc: "time should be valid when it's the current time",
		},
		{
			name: "zero time",
			time: time.Time{},
			want: core_errors.ErrInvalidTripEndTime,
			desc: "time should be invalid when it's zero time",
		},
		{
			name: "past time",
			time: frozenTime.Add(-24 * time.Hour),
			want: core_errors.ErrInvalidTripEndTime,
			desc: "time should be invalid when it's in the past",
		},
		{
			name: "far future time",
			time: frozenTime.Add(365 * 24 * time.Hour),
			want: nil,
			desc: "time should be valid when it's far in the future",
		},
		{
			name: "just one second in the past",
			time: frozenTime.Add(-1 * time.Second),
			want: core_errors.ErrInvalidTripEndTime,
			desc: "time should be invalid when it's even one second in the past",
		},
		{
			name: "time with different timezone",
			time: frozenTime.In(time.FixedZone("UTC+3", 3*3600)),
			want: nil,
			desc: "time should be valid when it's in a different timezone but represents the same moment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTripEndTime(tt.time)
			if tt.want != nil {
				require.ErrorIs(t, err, tt.want, "test description: %s", tt.desc)
			} else {
				require.NoError(t, err, "test description: %s", tt.desc)
			}
		})
	}
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
	// Zero coordinates are valid (0, 0 is a valid coordinate)
	require.NoError(t, validateCoordinates(Coordinates{}))
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
