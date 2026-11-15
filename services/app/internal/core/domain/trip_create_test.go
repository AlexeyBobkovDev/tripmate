package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestTripCreate_Validate(t *testing.T) {
	now := time.Now()
	futureDate := now.Add(24 * time.Hour)
	laterDate := now.Add(48 * time.Hour)

	tests := []struct {
		name    string
		create  TripCreate
		wantErr error
	}{
		{
			name: "valid trip create",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2027, time.January, 15, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2027, time.January, 15, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{
					{
						Type:     MealBreakfast,
						Window:   TimeWindow{Start: "09:00", End: "10:00"},
						Duration: time.Hour,
						Budget:   Money{AmountMinorUnits: 1500000, Currency: "RUB"},
					},
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
		},
		{
			name: "invalid multi-day trip - start before end",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2027, time.January, 15, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2027, time.January, 16, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "invalid multi-day trip - start after end",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Date(2027, time.January, 16, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2027, time.January, 15, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				MealRequirements: []MealRequirement{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "owner ID is validated by the service layer",
			create: TripCreate{
				OwnerID:          -1,
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
			},
		},
		{
			name: "invalid name",
			create: TripCreate{
				OwnerID:          1,
				Name:             "",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidTripName,
		},
		{
			name: "invalid planning area",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidPlanningArea,
		},
		{
			name: "invalid start coordinates",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 200, Longitude: 200},
				StartDate:        futureDate,
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidLatitude,
		},
		{
			name: "invalid start date",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        time.Time{},
				EndDate:          time.Date(2026, time.October, 31, 0, 0, 0, 0, time.UTC),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidTripStartDate,
		},
		{
			name: "invalid start date - past",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        now.AddDate(-1, 0, 0),
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidTripStartDate,
		},
		{
			name: "invalid end date",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          time.Time{},
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "invalid end date - past",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          now.AddDate(-1, 0, 0),
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "invalid end date - before start",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        laterDate,
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "valid single-day trip - equal start and end",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
		},
		{
			name: "invalid daily schedule",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          futureDate,
				DailySchedule:    DailySchedule{},
				Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests:        []Interest{InterestHistory},
				TimeZone:         "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidLocalTime,
		},
		{
			name: "invalid budget",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 0, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "invalid interests",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{"invalid"},
				TimeZone:  "Europe/Moscow",
			},
			wantErr: core_errors.ErrUnsupportedInterest,
		},
		{
			name: "invalid time zone",
			create: TripCreate{
				OwnerID:          1,
				Name:             "Test Trip",
				PlanningArea:     "Moscow",
				StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				StartDate:        futureDate,
				EndDate:          futureDate,
				DailySchedule: DailySchedule{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
				Budget:    Money{AmountMinorUnits: 1500000, Currency: "RUB"},
				Interests: []Interest{InterestHistory},
				TimeZone:  "Invalid/Timezone",
			},
			wantErr: core_errors.ErrUnsupportedTimezone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.create.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewTripCreate(t *testing.T) {
	startDate := time.Now().Add(24 * time.Hour)
	create := NewTripCreate(
		"Test Trip",
		"Moscow",
		Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		startDate,
		startDate,
		DailySchedule{StartTime: "09:00", EndTime: "18:00"},
		[]MealRequirement{},
		Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		[]Interest{InterestHistory},
		"Europe/Moscow",
	)

	// Verify all fields are set
	assert.Equal(t, "Test Trip", create.Name)
	assert.Equal(t, "Moscow", create.PlanningArea)
	assert.Equal(t, Coordinates{Latitude: 55.7539, Longitude: 37.6208}, create.StartCoordinates)
	assert.Equal(t, startDate, create.StartDate)
	assert.Equal(t, startDate, create.EndDate)
	assert.Equal(t, DailySchedule{StartTime: "09:00", EndTime: "18:00"}, create.DailySchedule)
	assert.Equal(t, Money{AmountMinorUnits: 1500000, Currency: "RUB"}, create.Budget)
	assert.Equal(t, []Interest{InterestHistory}, create.Interests)
	assert.Equal(t, "Europe/Moscow", create.TimeZone)

	// Verify validation passes
	require.NoError(t, create.Validate())
}
