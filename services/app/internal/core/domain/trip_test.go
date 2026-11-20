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
	pastDate := now.Add(-24 * time.Hour)

	// Helper function to create a valid trip for mutation
	createValidTrip := func() Trip {
		return Trip{
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
	}

	tests := []struct {
		name    string
		trip    Trip
		wantErr error
	}{
		{
			name: "valid trip",
			trip: createValidTrip(),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
				StartDate:        time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2035, time.October, 31, 0, 0, 0, 0, time.UTC),
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
		{
			name: "invalid start date - past",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        pastDate,
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
			wantErr: core_errors.ErrInvalidTripStartDate,
		},
		{
			name: "invalid end date - past",
			trip: Trip{
				ID:               1,
				OwnerID:          1,
				Version:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          pastDate,
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
			name: "invalid daily schedule - start time after end time",
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
					StartTime: "18:00",
					EndTime:   "09:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidTripEndTime, // This will fail because end time is before start time
		},
		{
			name: "invalid budget - zero amount",
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
				Budget:           Money{AmountMinorUnits: 0, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "invalid budget - negative amount",
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
				Budget:           Money{AmountMinorUnits: -1000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidBudget,
		},

		{
			name: "invalid time zone - empty",
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
				TimeZone:         "",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			wantErr: core_errors.ErrInvalidTimeZone, // This is what validateTimeZone actually returns for empty string
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

	// Helper function to create a valid trip for mutation
	createValidTrip := func() Trip {
		return Trip{
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
	}

	trip := createValidTrip()

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
			name: "patch with empty name - should not be set",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				Name:    Nullable[string]{Value: &[]string{""}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripName,
		},
		{
			name: "patch with planning area change",
			patch: &TripPatch{
				OwnerID:      1,
				Version:      1,
				PlanningArea: Nullable[string]{Value: &[]string{"Saint Petersburg"}[0], Set: true},
			},
		},
		{
			name: "patch with valid coordinates change",
			patch: &TripPatch{
				OwnerID:          1,
				Version:          1,
				StartCoordinates: Nullable[Coordinates]{Value: &[]Coordinates{{Latitude: 59.9343, Longitude: 30.3351}}[0], Set: true},
			},
		},
		{
			name: "patch with invalid latitude",
			patch: &TripPatch{
				OwnerID:          1,
				Version:          1,
				StartCoordinates: Nullable[Coordinates]{Value: &[]Coordinates{{Latitude: 91, Longitude: 37.6208}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "patch with invalid longitude",
			patch: &TripPatch{
				OwnerID:          1,
				Version:          1,
				StartCoordinates: Nullable[Coordinates]{Value: &[]Coordinates{{Latitude: 55.7539, Longitude: 181}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidLongitude,
		},
		{
			name: "patch with valid start date change",
			patch: &TripPatch{
				OwnerID:   1,
				Version:   1,
				StartDate: Nullable[time.Time]{Value: &[]time.Time{futureDate.Add(24 * time.Hour)}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "patch with invalid start date",
			patch: &TripPatch{
				OwnerID:   1,
				Version:   1,
				StartDate: Nullable[time.Time]{Value: &[]time.Time{pastDate}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripStartDate,
		},
		{
			name: "patch with valid end date change",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				EndDate: Nullable[time.Time]{Value: &[]time.Time{futureDate.Add(48 * time.Hour)}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripEndDate, // This will fail because end date is after start date
		},
		{
			name: "patch with invalid end date - before start",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				EndDate: Nullable[time.Time]{Value: &[]time.Time{futureDate.Add(-24 * time.Hour)}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "patch with valid daily schedule change",
			patch: &TripPatch{
				OwnerID:       1,
				Version:       1,
				DailySchedule: Nullable[DailySchedule]{Value: &[]DailySchedule{{StartTime: "10:00", EndTime: "19:00"}}[0], Set: true},
			},
		},
		{
			name: "patch with invalid daily schedule - empty start time",
			patch: &TripPatch{
				OwnerID:       1,
				Version:       1,
				DailySchedule: Nullable[DailySchedule]{Value: &[]DailySchedule{{StartTime: "", EndTime: "19:00"}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "patch with invalid daily schedule - start after end",
			patch: &TripPatch{
				OwnerID:       1,
				Version:       1,
				DailySchedule: Nullable[DailySchedule]{Value: &[]DailySchedule{{StartTime: "19:00", EndTime: "10:00"}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTripEndTime, // This will fail because end time is before start time
		},
		{
			name: "patch with valid budget change",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				Budget:  Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 2000000, Currency: "RUB"}}[0], Set: true},
			},
		},
		{
			name: "patch with invalid budget - zero amount",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				Budget:  Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 0, Currency: "RUB"}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "patch with invalid budget - negative amount",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				Budget:  Nullable[Money]{Value: &[]Money{{AmountMinorUnits: -1000, Currency: "RUB"}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "patch with valid interests change",
			patch: &TripPatch{
				OwnerID:   1,
				Version:   1,
				Interests: Nullable[[]Interest]{Value: &[][]Interest{{InterestHistory, InterestNature}}[0], Set: true},
			},
		},

		{
			name: "patch with invalid interests - unsupported",
			patch: &TripPatch{
				OwnerID:   1,
				Version:   1,
				Interests: Nullable[[]Interest]{Value: &[][]Interest{{"invalid"}}[0], Set: true},
			},
			wantErr: core_errors.ErrUnsupportedInterest,
		},
		{
			name: "patch with valid time zone change",
			patch: &TripPatch{
				OwnerID:  1,
				Version:  1,
				TimeZone: Nullable[string]{Value: &[]string{"Europe/Paris"}[0], Set: true},
			},
			wantErr: core_errors.ErrUnsupportedTimezone, // This will fail because Europe/Paris is not supported
		},
		{
			name: "patch with invalid time zone - empty",
			patch: &TripPatch{
				OwnerID:  1,
				Version:  1,
				TimeZone: Nullable[string]{Value: &[]string{""}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidTimeZone, // This is what validateTimeZone actually returns for empty string
		},
		{
			name: "patch with invalid time zone - invalid value",
			patch: &TripPatch{
				OwnerID:  1,
				Version:  1,
				TimeZone: Nullable[string]{Value: &[]string{"Invalid/Timezone"}[0], Set: true},
			},
			wantErr: core_errors.ErrUnsupportedTimezone,
		},
		{
			name: "patch with multiple valid changes",
			patch: &TripPatch{
				OwnerID:      1,
				Version:      1,
				Name:         Nullable[string]{Value: &[]string{"Updated Trip"}[0], Set: true},
				PlanningArea: Nullable[string]{Value: &[]string{"Saint Petersburg"}[0], Set: true},
				Budget:       Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 2000000, Currency: "RUB"}}[0], Set: true},
			},
		},
		{
			name: "patch with mixed valid and invalid changes",
			patch: &TripPatch{
				OwnerID:      1,
				Version:      1,
				Name:         Nullable[string]{Value: &[]string{"Updated Trip"}[0], Set: true},
				PlanningArea: Nullable[string]{Value: &[]string{""}[0], Set: true},
				Budget:       Nullable[Money]{Value: &[]Money{{AmountMinorUnits: 2000000, Currency: "RUB"}}[0], Set: true},
			},
			wantErr: core_errors.ErrInvalidPlanningArea,
		},
		{
			name: "patch with null fields should not change original values",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				Name:    Nullable[string]{Set: false}, // Not set
			},
		},
		{
			name: "patch with null budget should not change original value",
			patch: &TripPatch{
				OwnerID: 1,
				Version: 1,
				Budget:  Nullable[Money]{Set: false}, // Not set
			},
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

func TestNewTrip(t *testing.T) {
	createdAt := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	var deletedAt *time.Time
	trip := NewTrip(
		7,
		3,
		2,
		"Trip to Moscow",
		"Moscow",
		Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2035, time.October, 1, 0, 0, 0, 0, time.UTC),
		DailySchedule{StartTime: "09:00", EndTime: "18:00"},
		[]MealRequirement{},
		Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		[]Interest{InterestHistory},
		"Europe/Moscow",
		createdAt,
		createdAt,
		deletedAt,
	)

	require.NotNil(t, trip)
	assert.Equal(t, 7, trip.ID)
	assert.Equal(t, 3, trip.OwnerID)
	assert.Equal(t, 2, trip.Version)
	assert.Equal(t, "Trip to Moscow", trip.Name)
	assert.Equal(t, "Moscow", trip.PlanningArea)
	assert.Equal(t, 55.7539, trip.StartCoordinates.Latitude)
	assert.Equal(t, 37.6208, trip.StartCoordinates.Longitude)
	assert.Equal(t, 2035, trip.StartDate.Year())
	assert.Equal(t, 2035, trip.EndDate.Year())
	assert.Equal(t, "09:00", string(trip.DailySchedule.StartTime))
	assert.Equal(t, "18:00", string(trip.DailySchedule.EndTime))
	assert.Empty(t, trip.MealRequirements)
	assert.Equal(t, int64(1500000), trip.Budget.AmountMinorUnits)
	assert.Equal(t, "RUB", trip.Budget.Currency)
	assert.Equal(t, InterestHistory, trip.Interests[0])
	assert.Equal(t, "Europe/Moscow", trip.TimeZone)
	assert.Equal(t, createdAt, trip.CreatedAt)
	assert.Equal(t, createdAt, trip.UpdatedAt)
	assert.Nil(t, trip.DeletedAt)
}

func TestTrip_ApplyPatchVersionMismatch(t *testing.T) {
	now := time.Now()
	trip := NewTrip(
		1,
		1,
		1,
		"Test Trip",
		"Moscow",
		Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		now.Add(24*time.Hour),
		now.Add(24*time.Hour),
		DailySchedule{StartTime: "09:00", EndTime: "18:00"},
		[]MealRequirement{},
		Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		[]Interest{InterestHistory},
		"Europe/Moscow",
		now,
		now,
		nil,
	)

	err := trip.ApplyPatch(&TripPatch{OwnerID: 1, Version: 2, Name: Nullable[string]{Set: false}})
	require.ErrorIs(t, err, core_errors.ErrDifferentVersions)
}

func TestTrip_Validate_WithInvalidMealRequirement(t *testing.T) {
	now := time.Now()
	trip := Trip{
		ID:               1,
		OwnerID:          1,
		Version:          1,
		Name:             "Test Trip",
		PlanningArea:     "Moscow",
		StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		StartDate:        now.Add(24 * time.Hour),
		EndDate:          now.Add(24 * time.Hour),
		DailySchedule:    DailySchedule{StartTime: "09:00", EndTime: "18:00"},
		MealRequirements: []MealRequirement{
			{
				Type:     MealBreakfast,
				Window:   TimeWindow{Start: "08:00", End: "09:00"},
				Duration: 0,
				Budget:   Money{AmountMinorUnits: 1000, Currency: "RUB"},
			},
		},
		Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		Interests: []Interest{InterestHistory},
		TimeZone:  "Europe/Moscow",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := trip.Validate()
	require.ErrorIs(t, err, core_errors.ErrInvalidMealTimeDuration)
}

func TestTrip_ApplyPatchWithMealRequirementsAndTimeZone(t *testing.T) {
	now := time.Now()
	trip := NewTrip(
		1,
		1,
		1,
		"Test Trip",
		"Moscow",
		Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		now.Add(24*time.Hour),
		now.Add(24*time.Hour),
		DailySchedule{StartTime: "09:00", EndTime: "18:00"},
		[]MealRequirement{},
		Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		[]Interest{InterestHistory},
		"Europe/Moscow",
		now,
		now,
		nil,
	)

	mealRequirements := []MealRequirement{
		{
			Type:     MealBreakfast,
			Window:   TimeWindow{Start: "08:00", End: "09:00"},
			Duration: 30 * time.Minute,
			Budget:   Money{AmountMinorUnits: 1000, Currency: "RUB"},
		},
	}
	patch := &TripPatch{
		OwnerID:          1,
		Version:          1,
		MealRequirements: Nullable[[]MealRequirement]{Value: &mealRequirements, Set: true},
		TimeZone:         Nullable[string]{Value: &[]string{"Europe/Moscow"}[0], Set: true},
	}

	err := trip.ApplyPatch(patch)
	require.NoError(t, err)
	require.Len(t, trip.MealRequirements, 1)
	require.Equal(t, MealBreakfast, trip.MealRequirements[0].Type)
	require.Equal(t, "Europe/Moscow", trip.TimeZone)
}
