package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/errors"
)

func tripTestValidTrip() *Trip {
	createdAt := time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	return NewTrip(
		1,
		2,
		1,
		"Weekend in Moscow",
		"Moscow",
		Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		time.Date(2027, time.January, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2027, time.January, 10, 0, 0, 0, 0, time.UTC),
		Money{AmountMinorUnits: 100000, Currency: "RUB"},
		[]Interest{InterestHistory, InterestFood},
		"Europe/Moscow",
		createdAt,
		updatedAt,
		nil,
	)
}

func TestTrip_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Trip)
		wantErr error
	}{
		{name: "valid trip"},
		{name: "start location omitted", mutate: func(trip *Trip) { trip.StartCoordinates = Coordinates{} }, wantErr: nil},
		{name: "budget omitted", mutate: func(trip *Trip) { trip.Budget = Money{} }, wantErr: core_errors.ErrInvalidBudget},
		{name: "empty interests allowed", mutate: func(trip *Trip) { trip.Interests = nil }},
		{name: "invalid trip ID", mutate: func(trip *Trip) { trip.ID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid owner ID", mutate: func(trip *Trip) { trip.OwnerID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid version", mutate: func(trip *Trip) { trip.Version = 0 }, wantErr: core_errors.ErrInvalidVersion},
		{name: "blank name", mutate: func(trip *Trip) { trip.Name = " \t " }, wantErr: core_errors.ErrInvalidTripName},
		{name: "blank planning area", mutate: func(trip *Trip) { trip.PlanningArea = "  " }, wantErr: core_errors.ErrInvalidPlanningArea},
		{name: "invalid start location", mutate: func(trip *Trip) { trip.StartCoordinates = Coordinates{Latitude: 100, Longitude: 0} }, wantErr: core_errors.ErrInvalidLatitude},
		{name: "zero start date", mutate: func(trip *Trip) { trip.StartDate = time.Time{} }, wantErr: core_errors.ErrInvalidTripStartDate},
		{name: "zero end date", mutate: func(trip *Trip) { trip.EndDate = time.Time{} }, wantErr: core_errors.ErrInvalidTripEndDate},
		{name: "end before start", mutate: func(trip *Trip) { trip.EndDate = trip.StartDate.Add(-time.Hour) }, wantErr: core_errors.ErrInvalidTripEndDate},
		{name: "same start and end", mutate: func(trip *Trip) { trip.EndDate = trip.StartDate }, wantErr: nil},
		{name: "invalid budget", mutate: func(trip *Trip) { trip.Budget = Money{AmountMinorUnits: 0, Currency: "RUB"} }, wantErr: core_errors.ErrInvalidBudget},
		{name: "unsupported interest", mutate: func(trip *Trip) { trip.Interests = []Interest{"shopping"} }, wantErr: core_errors.ErrUnsupportedInterest},
		{name: "blank interest", mutate: func(trip *Trip) { trip.Interests = []Interest{"  "} }, wantErr: core_errors.ErrInvalidInterest},
		{name: "blank time zone", mutate: func(trip *Trip) { trip.TimeZone = " " }, wantErr: core_errors.ErrInvalidTimeZone},
		{name: "unsupported time zone", mutate: func(trip *Trip) { trip.TimeZone = "America/New_York" }, wantErr: core_errors.ErrUnsupportedTimezone},
		{name: "zero created at", mutate: func(trip *Trip) { trip.CreatedAt = time.Time{} }, wantErr: core_errors.ErrInvalidCreatedAt},
		{name: "updated before created", mutate: func(trip *Trip) { trip.UpdatedAt = trip.CreatedAt.Add(-time.Second) }, wantErr: core_errors.ErrInvalidUpdatedAt},
		{name: "start date before creation", mutate: func(trip *Trip) { trip.StartDate = trip.CreatedAt.Add(-time.Hour); trip.EndDate = trip.StartDate }, wantErr: core_errors.ErrInvalidTripStartDate},
		{name: "deleted before updated", mutate: func(trip *Trip) { deleted := trip.UpdatedAt.Add(-time.Second); trip.DeletedAt = &deleted }, wantErr: core_errors.ErrInvalidDeletedAt},
		{name: "deleted at in future", mutate: func(trip *Trip) { deleted := time.Now().Add(time.Hour); trip.DeletedAt = &deleted }, wantErr: core_errors.ErrInvalidDeletedAt},
		{name: "future updated at", mutate: func(trip *Trip) { trip.UpdatedAt = time.Now().Add(time.Hour) }, wantErr: core_errors.ErrInvalidUpdatedAt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trip := tripTestValidTrip()
			if tt.mutate != nil {
				tt.mutate(trip)
			}
			err := trip.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewTrip(t *testing.T) {
	createdAt := time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	deletedAt := updatedAt.Add(time.Hour)
	location := Coordinates{Latitude: 55.7539, Longitude: 37.6208}
	budget := Money{AmountMinorUnits: 250000, Currency: "RUB"}
	interests := []Interest{InterestHistory, InterestFood}
	startDate := time.Date(2027, time.March, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate

	got := NewTrip(10, 20, 3, "Trip", "Moscow", location, startDate, endDate, budget, interests, "Europe/Moscow", createdAt, updatedAt, &deletedAt)

	require.NotNil(t, got)
	assert.Equal(t, 10, got.ID)
	assert.Equal(t, 20, got.OwnerID)
	assert.Equal(t, 3, got.Version)
	assert.Equal(t, "Trip", got.Name)
	assert.Equal(t, "Moscow", got.PlanningArea)
	assert.Equal(t, location, got.StartCoordinates)
	assert.True(t, startDate.Equal(got.StartDate))
	assert.True(t, endDate.Equal(got.EndDate))
	assert.Equal(t, budget, got.Budget)
	assert.Equal(t, interests, got.Interests)
	assert.Equal(t, "Europe/Moscow", got.TimeZone)
	assert.True(t, createdAt.Equal(got.CreatedAt))
	assert.True(t, updatedAt.Equal(got.UpdatedAt))
	require.NotNil(t, got.DeletedAt)
	assert.True(t, deletedAt.Equal(*got.DeletedAt))
}
