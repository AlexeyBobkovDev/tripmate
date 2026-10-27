package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func tripPatchTestValid() *TripPatch {
	return &TripPatch{
		OwnerID: 1,
		Version: 1,
	}
}

func TestTripPatchValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*TripPatch)
		wantErr error
	}{
		{
			name: "valid patch with no fields set",
		},
		{
			name: "invalid owner ID",
			mutate: func(p *TripPatch) {
				p.OwnerID = 0
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid version",
			mutate: func(p *TripPatch) {
				p.Version = 0
			},
			wantErr: core_errors.ErrInvalidVersion,
		},
		{
			name: "name value is nil when set",
			mutate: func(p *TripPatch) {
				p.Name = Nullable[string]{Set: true}
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "empty name",
			mutate: func(p *TripPatch) {
				value := " \t "
				p.Name = Nullable[string]{Set: true, Value: &value}
			},
			wantErr: core_errors.ErrInvalidTripName,
		},
		{
			name: "valid name",
			mutate: func(p *TripPatch) {
				value := "Weekend in Moscow"
				p.Name = Nullable[string]{Set: true, Value: &value}
			},
		},
		{
			name: "planning area value is nil when set",
			mutate: func(p *TripPatch) {
				p.PlanningArea = Nullable[string]{Set: true}
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "empty planning area",
			mutate: func(p *TripPatch) {
				value := ""
				p.PlanningArea = Nullable[string]{Set: true, Value: &value}
			},
			wantErr: core_errors.ErrInvalidPlanningArea,
		},
		{
			name: "start location value is nil when set",
			mutate: func(p *TripPatch) {
				p.StartLocation = Nullable[*Location]{Set: true}
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "invalid start location",
			mutate: func(p *TripPatch) {
				location := &Location{Name: " ", Latitude: 0, Longitude: 0}
				p.StartLocation = Nullable[*Location]{Set: true, Value: &location}
			},
			wantErr: core_errors.ErrInvalidLocationName,
		},
		{
			name: "start date value is nil when set",
			mutate: func(p *TripPatch) {
				p.StartDate = Nullable[time.Time]{Set: true}
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "zero start date",
			mutate: func(p *TripPatch) {
				value := time.Time{}
				p.StartDate = Nullable[time.Time]{Set: true, Value: &value}
			},
			wantErr: core_errors.ErrInvalidTripStartDate,
		},
		{
			name: "end date value is nil when set",
			mutate: func(p *TripPatch) {
				p.EndDate = Nullable[time.Time]{Set: true}
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "zero end date",
			mutate: func(p *TripPatch) {
				value := time.Time{}
				p.EndDate = Nullable[time.Time]{Set: true, Value: &value}
			},
			wantErr: core_errors.ErrInvalidTripEndDate,
		},
		{
			name: "budget value is nil when set",
			mutate: func(p *TripPatch) {
				p.Budget = Nullable[*Money]{Set: true}
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "invalid budget",
			mutate: func(p *TripPatch) {
				budget := &Money{AmountMinorUnits: 0, Currency: "RUB"}
				p.Budget = Nullable[*Money]{Set: true, Value: &budget}
			},
			wantErr: core_errors.ErrInvalidBudget,
		},
		{
			name: "interests value is nil when set",
			mutate: func(p *TripPatch) {
				p.Interests = Nullable[[]Interest]{Set: true}
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "unsupported interest",
			mutate: func(p *TripPatch) {
				value := []Interest{"not-supported"}
				p.Interests = Nullable[[]Interest]{Set: true, Value: &value}
			},
			wantErr: core_errors.ErrUnsupportedInterest,
		},
		{
			name: "time zone value is nil when set",
			mutate: func(p *TripPatch) {
				p.TimeZone = Nullable[string]{Set: true}
			},
			wantErr: core_errors.ErrInvalidArgument,
		},
		{
			name: "empty time zone",
			mutate: func(p *TripPatch) {
				value := ""
				p.TimeZone = Nullable[string]{Set: true, Value: &value}
			},
			wantErr: core_errors.ErrInvalidTimeZone,
		},
		{
			name: "unsupported time zone",
			mutate: func(p *TripPatch) {
				value := "Mars/Olympus"
				p.TimeZone = Nullable[string]{Set: true, Value: &value}
			},
			wantErr: core_errors.ErrUnsupportedTimezone,
		},
		{
			name: "valid values for all fields",
			mutate: func(p *TripPatch) {
				name := "Weekend trip"
				area := "Moscow"
				location := &Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208}
				start := time.Now().Add(48 * time.Hour)
				end := start.Add(24 * time.Hour)
				budget := &Money{AmountMinorUnits: 100000, Currency: "RUB"}
				interests := []Interest{InterestHistory, InterestFood}
				timeZone := "Europe/Moscow"

				p.Name = Nullable[string]{Set: true, Value: &name}
				p.PlanningArea = Nullable[string]{Set: true, Value: &area}
				p.StartLocation = Nullable[*Location]{Set: true, Value: &location}
				p.StartDate = Nullable[time.Time]{Set: true, Value: &start}
				p.EndDate = Nullable[time.Time]{Set: true, Value: &end}
				p.Budget = Nullable[*Money]{Set: true, Value: &budget}
				p.Interests = Nullable[[]Interest]{Set: true, Value: &interests}
				p.TimeZone = Nullable[string]{Set: true, Value: &timeZone}
			},
		},
		{
			name: "optional location can be cleared",
			mutate: func(p *TripPatch) {
				var location *Location
				p.StartLocation = Nullable[*Location]{Set: true, Value: &location}
			},
		},
		{
			name: "optional budget can be cleared",
			mutate: func(p *TripPatch) {
				var budget *Money
				p.Budget = Nullable[*Money]{Set: true, Value: &budget}
			},
		},
		{
			name: "empty interests list is valid",
			mutate: func(p *TripPatch) {
				value := []Interest{}
				p.Interests = Nullable[[]Interest]{Set: true, Value: &value}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patch := tripPatchTestValid()
			if tt.mutate != nil {
				tt.mutate(patch)
			}

			err := patch.Validate()
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewTripPatch(t *testing.T) {
	startDate := time.Date(2027, time.June, 1, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2027, time.June, 3, 12, 0, 0, 0, time.UTC)
	name := "Summer trip"
	area := "Moscow"
	location := &Location{Name: "Red Square", Latitude: 55.7539, Longitude: 37.6208}
	budget := &Money{AmountMinorUnits: 250000, Currency: "RUB"}
	interests := []Interest{InterestHistory, InterestArchitecture}
	timeZone := "Europe/Moscow"

	nameField := Nullable[string]{Set: true, Value: &name}
	areaField := Nullable[string]{Set: true, Value: &area}
	locationField := Nullable[*Location]{Set: true, Value: &location}
	startDateField := Nullable[time.Time]{Set: true, Value: &startDate}
	endDateField := Nullable[time.Time]{Set: true, Value: &endDate}
	budgetField := Nullable[*Money]{Set: true, Value: &budget}
	interestsField := Nullable[[]Interest]{Set: true, Value: &interests}
	timeZoneField := Nullable[string]{Set: true, Value: &timeZone}

	got := NewTripPatch(
		42,
		7,
		nameField,
		areaField,
		locationField,
		startDateField,
		endDateField,
		budgetField,
		interestsField,
		timeZoneField,
	)

	require.NotNil(t, got)
	assert.Equal(t, 42, got.OwnerID)
	assert.Equal(t, 7, got.Version)
	assert.Equal(t, nameField, got.Name)
	assert.Equal(t, areaField, got.PlanningArea)
	assert.Equal(t, locationField, got.StartLocation)
	assert.Equal(t, startDateField, got.StartDate)
	assert.Equal(t, endDateField, got.EndDate)
	assert.Equal(t, budgetField, got.Budget)
	assert.Equal(t, interestsField, got.Interests)
	assert.Equal(t, timeZoneField, got.TimeZone)
}
