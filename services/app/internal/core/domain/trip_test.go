package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestTrip_Validate(t *testing.T) {
	now := time.Now()
	futureDate := now.Add(24 * time.Hour)

	tests := []struct {
		name    string
		trip    Trip
		wantErr error
	}{
		{
			name: "valid trip",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
		},
		{
			name: "invalid ID",
			trip: Trip{
				ID:               -1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid owner ID",
			trip: Trip{
				ID:               1,
				OwnerID:          -1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid version",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          0,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidVersion,
		},
		{
			name: "invalid name",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidTripName,
		},
		{
			name: "invalid planning area",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidPlanningArea,
		},
		{
			name: "invalid start coordinates",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 200, Longitude: 200},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "invalid start date",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Time{},
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidTripStartDate,
		},
		{
			name: "invalid end date",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Time{},
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "invalid end date - before start",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "invalid budget",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 0, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "invalid interests",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{"invalid"},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrUnsupportedInterest,
		},
		{
			name: "invalid time zone",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Invalid/Timezone",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrUnsupportedTimezone,
		},
		{
			name: "invalid created at",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        time.Time{},
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidCreatedAt,
		},
		{
			name: "invalid updated at",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        time.Time{},
			},
			wantErr: core_errors.ErrInvalidUpdatedAt,
		},
		{
			name: "invalid deleted at - future",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
				DeletedAt:        &[]time.Time{now.AddDate(1, 0, 0)}[0],
			},
			wantErr: core_errors.ErrInvalidDeletedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.trip.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestTrip_ApplyPatch(t *testing.T) {
	now := time.Now()
	futureDate := now.Add(24 * time.Hour)
	pastDate := now.Add(-24 * time.Hour)
	trip := Trip{
		ID:               1,
		OwnerID:          1,
		Version:          1,
		Name:             "Test Trip",
		PlanningArea:     "Moscow",
		StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		StartDate:        futureDate,
		EndDate:          futureDate,
		DailySchedule: DailySchedule{
			StartTime: "09:00",
			EndTime:   "18:00",
		},
		MealRequirements: []MealRequirement{},
		Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		Interests:        []Interest{InterestHistory},
		TimeZone:         "Europe/Moscow",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	tests := []struct {
		name    string
		patch   *TripPatch
		wantErr error
	}{
		{
			name: "valid patch",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				Name:    Nullable[string]{Value: &[]string{"Updated Trip"}[0], Set: true},
			},
		},
		{
			name: "invalid version",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 2,
			},
			wantErr: core_errors.ErrDifferentVersions,
		},
		{
			name: "invalid name",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				Name:    Nullable[string]{Value: &[]string{""}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripName,
		},
		{
			name: "patch with invalid planning area",
			patch: &TripPatch{
				OwnerID:      1,
				Version:      1,
				PlanningArea: Nullable[string]{Value: &[]string{""}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidPlanningArea,
		},
		{
			name: "patch with invalid coordinates",
			patch: &TripPatch{
				OwnerID:          1,
				Version:          1,
				StartCoordinates: Nullable[Coordinates]{Value: &[]Coordinates{{Latitude: 200, Longitude: 200}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "patch with invalid start date",
			patch: &TripPatch{
				OwnerID:   1,
				Version:   1,
				StartDate: Nullable[time.Time]{Value: &[]time.Time{pastDate}[0], Set: true},
				EndDate:   Nullable[time.Time]{Value: &[]time.Time{pastDate}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripStartDate,
		},
		{
			name: "patch with invalid end date",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				EndDate: Nullable[time.Time]{Value: &[]time.Time{now}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "patch with invalid daily schedule",
			patch: &TripPatch{
				OwnerID:       1,
				Version:       1,
				DailySchedule: Nullable[DailySchedule]{Value: &[]DailySchedule{{}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "patch with invalid budget",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				Budget:  Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 0}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "patch with invalid interests",
			patch: &TripPatch{
				OwnerID:   1,
				Version:   1,
				Interests: Nullable[[]Interest]{Value: &[][]Interest{{"invalid"}}[0], Set: true},
			},
			wantErr: core_errors.ErrUnsupportedInterest,
		},
		{
			name: "patch with invalid time zone",
			patch: &TripPatch{
				OwnerID:  1,
				Version:  1,
				TimeZone: Nullable[string]{Value: &[]string{"Invalid"}[0], Set: true},
			},
			wantErr: core_errors.ErrUnsupportedTimezone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current := trip
			err := current.ApplyPatch(tt.patch)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
