package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestTripPatch_Validate(t *testing.T) {
	now := time.Now()
	futureDate := now.Add(24 * time.Hour)
	laterDate := now.Add(48 * time.Hour)

	tests := []struct {
		name    string
		patch   TripPatch
		wantErr error
	}{
		{
			name: "valid patch - all fields set",
			patch: TripPatch{
				OwnerID:          1,
				Version:          1,
				Name:             Nullable[string]{Value: &[]string{"Updated Trip"}[0], Set: true},
				PlanningArea:     Nullable[string]{Value: &[]string{"Updated Area"}[0], Set: true},
				StartCoordinates: Nullable[Coordinates]{Value: &[]Coordinates{{Latitude: 55.7539, Longitude: 37.6208}}[0], Set: true},
				StartDate:        Nullable[time.Time]{Value: &[]time.Time{futureDate}[0], Set: true},
				EndDate:          Nullable[time.Time]{Value: &[]time.Time{futureDate}[0], Set: true},
				DailySchedule:    Nullable[DailySchedule]{Value: &[]DailySchedule{{StartTime: "10:00", EndTime: "19:00"}}[0], Set: true},
				MealRequirements: Nullable[[]MealRequirement]{Value: &[][]MealRequirement{{}}[0], Set: true},
				Budget:           Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 2000000, Currency: "RUB"}}[0], Set: true},
				Interests:        Nullable[[]Interest]{Value: &[][]Interest{{InterestNature}}[0], Set: true},
				TimeZone:         Nullable[string]{Value: &[]string{"Europe/Moscow"}[0], Set: true},
			},
		},
		{
			name: "valid patch - no fields set",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
			},
		},
		{
			name: "valid patch - only name set",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				Name:    Nullable[string]{Value: &[]string{"Updated Trip"}[0], Set: true},
			},
		},
		{
			name: "valid patch - only budget set",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				Budget:  Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 2000000, Currency: "RUB"}}[0], Set: true},
			},
		},
		{
			name: "invalid owner ID",
			patch: TripPatch{
				OwnerID: -1,
				Version: 1,
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid version",
			patch: TripPatch{
				OwnerID: 1,
				Version: 0,
			},
			wantErr: core_errors.ErrInvalidVersion,
		},
		{
			name: "invalid name - empty",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				Name:    Nullable[string]{Value: &[]string{""}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripName,
		},
		{
			name: "long name is allowed",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				Name:    Nullable[string]{Value: &[]string{"This is a very long name that exceeds the maximum allowed length for trip names"}[0], Set: true},
			},
		},
		{
			name: "invalid planning area",
			patch: TripPatch{
				OwnerID:      1,
				Version:      1,
				PlanningArea: Nullable[string]{Value: &[]string{""}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidPlanningArea,
		},
		{
			name: "invalid coordinates",
			patch: TripPatch{
				OwnerID:          1,
				Version:          1,
				StartCoordinates: Nullable[Coordinates]{Value: &[]Coordinates{{Latitude: 200, Longitude: 200}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "invalid start date - zero value",
			patch: TripPatch{
				OwnerID:   1,
				Version:   1,
				StartDate: Nullable[time.Time]{Value: &[]time.Time{{}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripStartDate,
		},
		{
			name: "invalid end date - zero value",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				EndDate: Nullable[time.Time]{Value: &[]time.Time{{}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "invalid start date after end date",
			patch: TripPatch{
				OwnerID:   1,
				Version:   1,
				StartDate: Nullable[time.Time]{Value: &[]time.Time{laterDate}[0], Set: true},
				EndDate:   Nullable[time.Time]{Value: &[]time.Time{futureDate}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "invalid daily schedule",
			patch: TripPatch{
				OwnerID:       1,
				Version:       1,
				DailySchedule: Nullable[DailySchedule]{Value: &[]DailySchedule{{}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid budget",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				Budget:  Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 0, Currency: "RUB"}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "invalid budget - wrong currency",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				Budget:  Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 1500000, Currency: "USD"}}[0], Set: true},
			},
			wantErr: core_errors.ErrUnsupportedCurrency,
		},
		{
			name: "invalid interests",
			patch: TripPatch{
				OwnerID:   1,
				Version:   1,
				Interests: Nullable[[]Interest]{Value: &[][]Interest{{"invalid"}}[0], Set: true},
			},
			wantErr: core_errors.ErrUnsupportedInterest,
		},
		{
			name: "invalid time zone",
			patch: TripPatch{
				OwnerID:  1,
				Version:  1,
				TimeZone: Nullable[string]{Value: &[]string{"Invalid"}[0], Set: true},
			},
			wantErr: core_errors.ErrUnsupportedTimezone,
		},
		{
			name: "name set to nil",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				Name:    Nullable[string]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "planning area set to nil",
			patch: TripPatch{
				OwnerID:      1,
				Version:      1,
				PlanningArea: Nullable[string]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "coordinates set to nil",
			patch: TripPatch{
				OwnerID:          1,
				Version:          1,
				StartCoordinates: Nullable[Coordinates]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "start date set to nil",
			patch: TripPatch{
				OwnerID:   1,
				Version:   1,
				StartDate: Nullable[time.Time]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "end date set to nil",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				EndDate: Nullable[time.Time]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "daily schedule set to nil",
			patch: TripPatch{
				OwnerID:       1,
				Version:       1,
				DailySchedule: Nullable[DailySchedule]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "meal requirements set to nil",
			patch: TripPatch{
				OwnerID:          1,
				Version:          1,
				MealRequirements: Nullable[[]MealRequirement]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "budget set to nil",
			patch: TripPatch{
				OwnerID: 1,
				Version: 1,
				Budget:  Nullable[Money]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "interests set to nil",
			patch: TripPatch{
				OwnerID:   1,
				Version:   1,
				Interests: Nullable[[]Interest]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "time zone set to nil",
			patch: TripPatch{
				OwnerID:  1,
				Version:  1,
				TimeZone: Nullable[string]{Value: nil, Set: true},
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.patch.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewTripPatch(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour)

	patch := NewTripPatch(
		1,
		1,
		Nullable[string]{Value: &[]string{"Updated Trip"}[0], Set: true},
		Nullable[string]{Value: &[]string{"Updated Area"}[0], Set: true},
		Nullable[Coordinates]{Value: &[]Coordinates{{Latitude: 55.7539, Longitude: 37.6208}}[0], Set: true},
		Nullable[time.Time]{Value: &[]time.Time{futureDate}[0], Set: true},
		Nullable[time.Time]{Value: &[]time.Time{futureDate}[0], Set: true},
		Nullable[DailySchedule]{Value: &[]DailySchedule{{StartTime: "10:00", EndTime: "19:00"}}[0], Set: true},
		Nullable[[]MealRequirement]{Value: &[][]MealRequirement{{}}[0], Set: true},
		Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 2000000, Currency: "RUB"}}[0], Set: true},
		Nullable[[]Interest]{Value: &[][]Interest{{InterestNature}}[0], Set: true},
		Nullable[string]{Value: &[]string{"Europe/Moscow"}[0], Set: true},
	)

	// Verify all fields are set
	assert.Equal(t, 1, patch.OwnerID)
	assert.Equal(t, 1, patch.Version)
	assert.True(t, patch.Name.Set)
	assert.True(t, patch.PlanningArea.Set)
	assert.True(t, patch.StartCoordinates.Set)
	assert.True(t, patch.StartDate.Set)
	assert.True(t, patch.EndDate.Set)
	assert.True(t, patch.DailySchedule.Set)
	assert.True(t, patch.MealRequirements.Set)
	assert.True(t, patch.Budget.Set)
	assert.True(t, patch.Interests.Set)
	assert.True(t, patch.TimeZone.Set)

	// Verify values are correct
	assert.Equal(t, "Updated Trip", *patch.Name.Value)
	assert.Equal(t, "Updated Area", *patch.PlanningArea.Value)
	assert.Equal(t, Coordinates{Latitude: 55.7539, Longitude: 37.6208}, *patch.StartCoordinates.Value)
	assert.Equal(t, futureDate, *patch.StartDate.Value)
	assert.Equal(t, futureDate, *patch.EndDate.Value)
	assert.Equal(t, DailySchedule{StartTime: "10:00", EndTime: "19:00"}, *patch.DailySchedule.Value)
	assert.Equal(t, Money{AmountMinorUnits: 2000000, Currency: "RUB"}, *patch.Budget.Value)
	assert.Equal(t, []Interest{InterestNature}, *patch.Interests.Value)
	assert.Equal(t, "Europe/Moscow", *patch.TimeZone.Value)

	// Verify validation passes
	require.NoError(t, patch.Validate())
}

func TestTripPatch_Validate_WithInvalidMealRequirement(t *testing.T) {
	mealRequirements := []MealRequirement{
		{
			Type:     MealBreakfast,
			Window:   TimeWindow{Start: "08:00", End: "09:00"},
			Duration: 0,
			Budget:   Money{AmountMinorUnits: 1000, Currency: "RUB"},
		},
	}
	patch := &TripPatch{
		OwnerID:          1,
		Version:          1,
		MealRequirements: Nullable[[]MealRequirement]{Value: &mealRequirements, Set: true},
	}

	err := patch.Validate()
	require.ErrorIs(t, err, core_errors.ErrInvalidMealTimeDuration)
}
