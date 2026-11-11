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
		Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		time.Date(2027, time.January, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2027, time.January, 10, 0, 0, 0, 0, time.UTC),
		validDailySchedule(),
		[]MealRequirement{validMealRequirement()},
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
		{name: "empty meal requirements allowed", mutate: func(trip *Trip) { trip.MealRequirements = nil }},
		{name: "invalid daily schedule start time", mutate: func(trip *Trip) { trip.DailySchedule.StartTime = "" }, wantErr: core_errors.ErrInvalidLocalTime},
		{name: "invalid daily schedule end time", mutate: func(trip *Trip) { trip.DailySchedule.EndTime = "25:00" }, wantErr: core_errors.ErrInvalidLocalTime},
		{name: "daily schedule start after end", mutate: func(trip *Trip) {
			trip.DailySchedule = DailySchedule{StartTime: "18:00", EndTime: "09:00"}
		}, wantErr: core_errors.ErrInvalidTripEndTime},
		{name: "invalid meal requirement type", mutate: func(trip *Trip) {
			trip.MealRequirements[0].Type = "brunch"
		}, wantErr: core_errors.ErrUnsupportedMealType},
		{name: "invalid meal requirement window", mutate: func(trip *Trip) {
			trip.MealRequirements[0].Window.Start = ""
		}, wantErr: core_errors.ErrInvalidLocalTime},
		{name: "meal requirement duration exceeds window", mutate: func(trip *Trip) {
			trip.MealRequirements[0].Duration = 2 * time.Hour
		}, wantErr: core_errors.ErrInvalidMealTimeEndTime},
		{name: "invalid meal requirement budget", mutate: func(trip *Trip) {
			trip.MealRequirements[0].Budget = Money{AmountMinorUnits: 0, Currency: "RUB"}
		}, wantErr: core_errors.ErrInvalidBudget},
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
		newEndDate := trip.StartDate.Add(-time.Hour)
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

	t.Run("clear start coordinates", func(t *testing.T) {
		trip := tripTestValidTrip()
		var newCoordinates Coordinates
		patch := &TripPatch{
			OwnerID:          trip.OwnerID,
			Version:          trip.Version,
			StartCoordinates: Nullable[Coordinates]{Value: &newCoordinates, Set: true},
		}

		require.NoError(t, trip.ApplyPatch(patch))
		assert.Equal(t, Coordinates{}, trip.StartCoordinates)
	})

	t.Run("applies valid start and end dates", func(t *testing.T) {
		trip := tripTestValidTrip()
		newStartDate := trip.StartDate.Add(24 * time.Hour)
		newEndDate := newStartDate
		patch := &TripPatch{
			OwnerID:   trip.OwnerID,
			Version:   trip.Version,
			StartDate: Nullable[time.Time]{Value: &newStartDate, Set: true},
			EndDate:   Nullable[time.Time]{Value: &newEndDate, Set: true},
		}

		require.NoError(t, trip.ApplyPatch(patch))
		assert.True(t, trip.StartDate.Equal(newStartDate))
		assert.True(t, trip.EndDate.Equal(newEndDate))
	})

	t.Run("applies valid budget", func(t *testing.T) {
		trip := tripTestValidTrip()
		newBudget := Money{AmountMinorUnits: 250000, Currency: "RUB"}
		patch := &TripPatch{
			OwnerID: trip.OwnerID,
			Version: trip.Version,
			Budget:  Nullable[Money]{Value: &newBudget, Set: true},
		}

		require.NoError(t, trip.ApplyPatch(patch))
		assert.Equal(t, newBudget, trip.Budget)
	})

	t.Run("clear budget", func(t *testing.T) {
		trip := tripTestValidTrip()
		var newBudget Money
		patch := &TripPatch{
			OwnerID: trip.OwnerID,
			Version: trip.Version,
			Budget:  Nullable[Money]{Value: &newBudget, Set: true},
		}

		require.ErrorIs(t, trip.ApplyPatch(patch), core_errors.ErrInvalidBudget)
	})

	t.Run("applies daily schedule and meal requirements", func(t *testing.T) {
		trip := tripTestValidTrip()
		newSchedule := DailySchedule{StartTime: "10:00", EndTime: "20:00"}
		newRequirements := []MealRequirement{
			{
				Type:     MealLunch,
				Window:   TimeWindow{Start: "12:00", End: "14:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 100000, Currency: "RUB"},
			},
			{
				Type:     MealDinner,
				Window:   TimeWindow{Start: "18:00", End: "20:00"},
				Duration: time.Hour,
				Budget:   Money{AmountMinorUnits: 200000, Currency: "RUB"},
			},
		}
		patch := &TripPatch{
			OwnerID:          trip.OwnerID,
			Version:          trip.Version,
			DailySchedule:    Nullable[DailySchedule]{Value: &newSchedule, Set: true},
			MealRequirements: Nullable[[]MealRequirement]{Value: &newRequirements, Set: true},
		}

		require.NoError(t, trip.ApplyPatch(patch))
		assert.Equal(t, newSchedule, trip.DailySchedule)
		assert.Equal(t, newRequirements, trip.MealRequirements)
	})

	t.Run("clears meal requirements", func(t *testing.T) {
		trip := tripTestValidTrip()
		emptyRequirements := []MealRequirement{}
		patch := &TripPatch{
			OwnerID:          trip.OwnerID,
			Version:          trip.Version,
			MealRequirements: Nullable[[]MealRequirement]{Value: &emptyRequirements, Set: true},
		}

		require.NoError(t, trip.ApplyPatch(patch))
		assert.Empty(t, trip.MealRequirements)
	})

	t.Run("rejects patch with invalid daily schedule", func(t *testing.T) {
		trip := tripTestValidTrip()
		before := *trip
		invalidSchedule := DailySchedule{StartTime: "20:00", EndTime: "09:00"}
		patch := &TripPatch{
			OwnerID:       trip.OwnerID,
			Version:       trip.Version,
			DailySchedule: Nullable[DailySchedule]{Value: &invalidSchedule, Set: true},
		}

		err := trip.ApplyPatch(patch)

		require.ErrorIs(t, err, core_errors.ErrInvalidTripEndTime)
		assert.Contains(t, err.Error(), "validate trip patch")
		assert.Equal(t, before, *trip)
	})

	t.Run("rejects patch with invalid meal requirements", func(t *testing.T) {
		trip := tripTestValidTrip()
		before := *trip
		invalidRequirements := []MealRequirement{
			{
				Type:     "",
				Window:   TimeWindow{Start: "12:00", End: "13:00"},
				Duration: time.Hour,
				Budget:   validMealBudget(),
			},
		}
		patch := &TripPatch{
			OwnerID:          trip.OwnerID,
			Version:          trip.Version,
			MealRequirements: Nullable[[]MealRequirement]{Value: &invalidRequirements, Set: true},
		}

		err := trip.ApplyPatch(patch)

		require.ErrorIs(t, err, core_errors.ErrUnsupportedMealType)
		assert.Contains(t, err.Error(), "validate trip patch")
		assert.Equal(t, before, *trip)
	})
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
	dailySchedule := DailySchedule{StartTime: "10:00", EndTime: "19:00"}
	mealRequirements := []MealRequirement{validMealRequirement()}

	got := NewTrip(10, 20, 3, "Trip", "Moscow", location, startDate, endDate, dailySchedule, mealRequirements, budget, interests, "Europe/Moscow", createdAt, updatedAt, &deletedAt)

	require.NotNil(t, got)
	assert.Equal(t, 10, got.ID)
	assert.Equal(t, 20, got.OwnerID)
	assert.Equal(t, 3, got.Version)
	assert.Equal(t, "Trip", got.Name)
	assert.Equal(t, "Moscow", got.PlanningArea)
	assert.Equal(t, location, got.StartCoordinates)
	assert.True(t, startDate.Equal(got.StartDate))
	assert.True(t, endDate.Equal(got.EndDate))
	assert.Equal(t, dailySchedule, got.DailySchedule)
	assert.Equal(t, mealRequirements, got.MealRequirements)
	assert.Equal(t, budget, got.Budget)
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
