package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
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
		&Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208},
		time.Date(2027, time.January, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2027, time.January, 15, 0, 0, 0, 0, time.UTC),
		&Money{AmountMinorUnits: 100000, Currency: "RUB"},
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
		{name: "optional start location omitted", mutate: func(trip *Trip) { trip.StartLocation = nil }},
		{name: "optional budget omitted", mutate: func(trip *Trip) { trip.Budget = nil }},
		{name: "empty interests allowed", mutate: func(trip *Trip) { trip.Interests = nil }},
		{name: "invalid trip ID", mutate: func(trip *Trip) { trip.ID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid owner ID", mutate: func(trip *Trip) { trip.OwnerID = 0 }, wantErr: core_errors.ErrInvalidID},
		{name: "invalid version", mutate: func(trip *Trip) { trip.Version = 0 }, wantErr: core_errors.ErrInvalidVersion},
		{name: "blank name", mutate: func(trip *Trip) { trip.Name = " \t " }, wantErr: core_errors.ErrInvalidTripName},
		{name: "blank planning area", mutate: func(trip *Trip) { trip.PlanningArea = "  " }, wantErr: core_errors.ErrInvalidPlanningArea},
		{name: "invalid start location", mutate: func(trip *Trip) { trip.StartLocation = &Location{Name: "", Latitude: 0, Longitude: 0} }, wantErr: core_errors.ErrInvalidLocationName},
		{name: "zero start date", mutate: func(trip *Trip) { trip.StartDate = time.Time{} }, wantErr: core_errors.ErrInvalidTripStartDate},
		{name: "zero end date", mutate: func(trip *Trip) { trip.EndDate = time.Time{} }, wantErr: core_errors.ErrInvalidTripEndDate},
		{name: "end before start", mutate: func(trip *Trip) { trip.EndDate = trip.StartDate.Add(-time.Hour) }, wantErr: core_errors.ErrInvalidTripEndDate},
		{name: "same start and end", mutate: func(trip *Trip) { trip.EndDate = trip.StartDate }, wantErr: core_errors.ErrInvalidTripEndDate},
		{name: "invalid budget", mutate: func(trip *Trip) { trip.Budget = &Money{AmountMinorUnits: 0, Currency: "RUB"} }, wantErr: core_errors.ErrInvalidBudget},
		{name: "unsupported interest", mutate: func(trip *Trip) { trip.Interests = []Interest{"shopping"} }, wantErr: core_errors.ErrUnsupportedInterest},
		{name: "blank interest", mutate: func(trip *Trip) { trip.Interests = []Interest{"  "} }, wantErr: core_errors.ErrInvalidInterest},
		{name: "blank time zone", mutate: func(trip *Trip) { trip.TimeZone = " " }, wantErr: core_errors.ErrInvalidTimeZone},
		{name: "unsupported time zone", mutate: func(trip *Trip) { trip.TimeZone = "America/New_York" }, wantErr: core_errors.ErrUnsupportedTimezone},
		{name: "zero created at", mutate: func(trip *Trip) { trip.CreatedAt = time.Time{} }, wantErr: core_errors.ErrInvalidCreatedAt},
		{name: "updated before created", mutate: func(trip *Trip) { trip.UpdatedAt = trip.CreatedAt.Add(-time.Second) }, wantErr: core_errors.ErrInvalidUpdatedAt},
		{name: "start date before creation", mutate: func(trip *Trip) { trip.StartDate = trip.CreatedAt }, wantErr: core_errors.ErrInvalidTripStartDate},
		{name: "deleted before updated", mutate: func(trip *Trip) { deleted := trip.UpdatedAt.Add(-time.Second); trip.DeletedAt = &deleted }, wantErr: core_errors.ErrInvalidDeletedAt},
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

func TestTrip_ApplyPatch(t *testing.T) {
	t.Run("updates supplied fields and preserves omitted fields", func(t *testing.T) {
		trip := tripTestValidTrip()
		originalStartDate := trip.StartDate
		newName := "A week in Moscow"
		newArea := "Central Moscow"
		newInterests := []Interest{InterestNature, InterestArchitecture}
		newTimeZone := "Europe/Moscow"
		patch := &TripPatch{
			OwnerID:      trip.OwnerID,
			Version:      trip.Version,
			Name:         Nullable[string]{Value: &newName, Set: true},
			PlanningArea: Nullable[string]{Value: &newArea, Set: true},
			Interests:    Nullable[[]Interest]{Value: &newInterests, Set: true},
			TimeZone:     Nullable[string]{Value: &newTimeZone, Set: true},
		}

		err := trip.ApplyPatch(patch)

		require.NoError(t, err)
		assert.Equal(t, newName, trip.Name)
		assert.Equal(t, newArea, trip.PlanningArea)
		assert.Equal(t, newInterests, trip.Interests)
		assert.Equal(t, newTimeZone, trip.TimeZone)
		assert.True(t, trip.StartDate.Equal(originalStartDate), "omitted fields should remain unchanged")
	})

	t.Run("rejects invalid patch and wraps validation error", func(t *testing.T) {
		trip := tripTestValidTrip()
		before := *trip
		patch := &TripPatch{OwnerID: trip.OwnerID, Version: trip.Version, Name: Nullable[string]{Set: true}}

		err := trip.ApplyPatch(patch)

		require.Error(t, err)
		assert.ErrorIs(t, err, core_errors.ErrInvalidArgument)
		assert.Contains(t, err.Error(), "validate trip patch")
		assert.Equal(t, before, *trip, "trip should not change when patch validation fails")
	})

	t.Run("rejects stale version", func(t *testing.T) {
		trip := tripTestValidTrip()
		before := *trip
		newName := "New name"
		patch := &TripPatch{
			OwnerID: trip.OwnerID,
			Version: trip.Version + 1,
			Name:    Nullable[string]{Value: &newName, Set: true},
		}

		err := trip.ApplyPatch(patch)

		assert.ErrorIs(t, err, core_errors.ErrDifferentVersions)
		assert.Equal(t, before, *trip, "trip should not change on version conflict")
	})

	t.Run("rejects patch that makes date range invalid", func(t *testing.T) {
		trip := tripTestValidTrip()
		before := *trip
		newEndDate := trip.StartDate
		patch := &TripPatch{
			OwnerID: trip.OwnerID,
			Version: trip.Version,
			EndDate: Nullable[time.Time]{Value: &newEndDate, Set: true},
		}

		err := trip.ApplyPatch(patch)

		assert.ErrorIs(t, err, core_errors.ErrInvalidTripEndDate)
		assert.Contains(t, err.Error(), "validate patched trip")
		assert.Equal(t, before, *trip, "trip should remain unchanged when resulting state is invalid")
	})

	t.Run("can clear optional start location", func(t *testing.T) {
		trip := tripTestValidTrip()
		var newLocation *Location
		patch := &TripPatch{
			OwnerID:       trip.OwnerID,
			Version:       trip.Version,
			StartLocation: Nullable[*Location]{Value: &newLocation, Set: true},
		}

		require.NoError(t, trip.ApplyPatch(patch))
		assert.Nil(t, trip.StartLocation)
	})

	t.Run("can clear optional budget", func(t *testing.T) {
		trip := tripTestValidTrip()
		var newBudget *Money
		patch := &TripPatch{
			OwnerID: trip.OwnerID,
			Version: trip.Version,
			Budget:  Nullable[*Money]{Value: &newBudget, Set: true},
		}

		require.NoError(t, trip.ApplyPatch(patch))
		assert.Nil(t, trip.Budget)
	})
}

func TestNewTrip(t *testing.T) {
	createdAt := time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	deletedAt := updatedAt.Add(time.Hour)
	location := &Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208}
	budget := &Money{AmountMinorUnits: 250000, Currency: "RUB"}
	interests := []Interest{InterestHistory, InterestFood}
	startDate := time.Date(2027, time.March, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 3)

	got := NewTrip(10, 20, 3, "Trip", "Moscow", location, startDate, endDate, budget, interests, "Europe/Moscow", createdAt, updatedAt, &deletedAt)

	require.NotNil(t, got)
	assert.Equal(t, 10, got.ID)
	assert.Equal(t, 20, got.OwnerID)
	assert.Equal(t, 3, got.Version)
	assert.Equal(t, "Trip", got.Name)
	assert.Equal(t, "Moscow", got.PlanningArea)
	assert.Same(t, location, got.StartLocation)
	assert.True(t, startDate.Equal(got.StartDate))
	assert.True(t, endDate.Equal(got.EndDate))
	assert.Same(t, budget, got.Budget)
	assert.Equal(t, interests, got.Interests)
	assert.Equal(t, "Europe/Moscow", got.TimeZone)
	assert.True(t, createdAt.Equal(got.CreatedAt))
	assert.True(t, updatedAt.Equal(got.UpdatedAt))
	require.NotNil(t, got.DeletedAt)
	assert.True(t, deletedAt.Equal(*got.DeletedAt))
}

func TestTrip_ApplyPatch_ErrorIsWrapped(t *testing.T) {
	trip := tripTestValidTrip()
	patch := &TripPatch{OwnerID: trip.OwnerID, Version: trip.Version, Name: Nullable[string]{Set: true}}
	err := trip.ApplyPatch(patch)
	require.Error(t, err)
	assert.True(t, errors.Is(err, core_errors.ErrInvalidArgument))
}
