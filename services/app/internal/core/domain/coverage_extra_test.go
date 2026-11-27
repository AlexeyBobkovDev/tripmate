package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestValidateRouteStopAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		address string
		wantErr error
	}{
		{name: "non-empty address", address: "Some address"},
		{name: "empty address", address: "", wantErr: core_errors.ErrInvalidAddress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateRouteStopAddress(tt.address)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateRouteStopWebsiteURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		websiteURL string
		wantErr    error
	}{
		{name: "non-empty website URL", websiteURL: "https://example.com"},
		{name: "empty website URL", websiteURL: "", wantErr: core_errors.ErrInvalidWebsiteURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateRouteStopWebsiteURL(tt.websiteURL)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateRouteStopVisitDurationMinutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		visitDurationMinutes int
		wantErr              error
	}{
		{name: "zero duration", visitDurationMinutes: 0},
		{name: "positive duration", visitDurationMinutes: 60},
		{name: "negative duration", visitDurationMinutes: -1, wantErr: core_errors.ErrInvalidRouteStopVisitDurationMinutes},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateRouteStopVisitDurationMinutes(tt.visitDurationMinutes)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateRouteStopOpeningHours(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		openingHours []OpeningHours
		wantErr      error
	}{
		{name: "no opening hours"},
		{name: "valid opening hours", openingHours: []OpeningHours{{OpensAt: 9, ClosesAt: 18}}},
		{
			name:         "opens at equals closes at",
			openingHours: []OpeningHours{{OpensAt: 9, ClosesAt: 9}},
			wantErr:      core_errors.ErrInvalidRouteStopOpeningHours,
		},
		{
			name:         "opens after closes",
			openingHours: []OpeningHours{{OpensAt: 18, ClosesAt: 9}},
			wantErr:      core_errors.ErrInvalidRouteStopOpeningHours,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateRouteStopOpeningHours(tt.openingHours)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		categories []Category
		wantErr    error
	}{
		{name: "nil list"},
		{name: "empty list", categories: []Category{}},
		{name: "supported category", categories: []Category{CategoryMuseum}},
		{name: "blank category", categories: []Category{" "}, wantErr: core_errors.ErrInvalidCategory},
		{name: "unsupported category", categories: []Category{"unsupported"}, wantErr: core_errors.ErrUnsupportedCategory},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateCategories(tt.categories)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestRouteStop_Validate_NewFields(t *testing.T) {
	t.Parallel()

	websiteURL := "https://stop.example.com"
	emptyWebsiteURL := ""
	validPhone := "+79990000000"
	invalidPhone := "12345"

	validStop := func() RouteStop {
		return RouteStop{
			ID:                   1,
			PlaceID:              1,
			Position:             1,
			Name:                 "Stop 1",
			Description:          "Description",
			Categories:           []Category{CategoryMuseum},
			Interests:            []Interest{InterestHistory},
			Address:              "Some stop address",
			WebsiteURL:           &websiteURL,
			Phone:                &validPhone,
			VisitDurationMinutes: 60,
			OpeningHours:         []OpeningHours{{OpensAt: 9, ClosesAt: 18}},
			Coordinates:          Coordinates{Latitude: 55.7539, Longitude: 37.6208},
			ArrivalTime:          time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
			DepartureTime:        time.Date(2026, time.October, 28, 14, 0, 0, 0, time.UTC),
		}
	}

	tests := []struct {
		name    string
		mutate  func(*RouteStop)
		wantErr error
	}{
		{name: "valid stop"},
		{
			name:    "invalid categories",
			mutate:  func(s *RouteStop) { s.Categories = []Category{"unsupported"} },
			wantErr: core_errors.ErrUnsupportedCategory,
		},
		{
			name:    "invalid interests",
			mutate:  func(s *RouteStop) { s.Interests = []Interest{"unsupported"} },
			wantErr: core_errors.ErrUnsupportedInterest,
		},
		{
			name:    "empty address",
			mutate:  func(s *RouteStop) { s.Address = "" },
			wantErr: core_errors.ErrInvalidAddress,
		},
		{
			name:    "empty website URL",
			mutate:  func(s *RouteStop) { s.WebsiteURL = &emptyWebsiteURL },
			wantErr: core_errors.ErrInvalidWebsiteURL,
		},
		{
			name:    "invalid phone",
			mutate:  func(s *RouteStop) { s.Phone = &invalidPhone },
			wantErr: core_errors.ErrInvalidPhoneNumber,
		},
		{
			name:    "negative visit duration",
			mutate:  func(s *RouteStop) { s.VisitDurationMinutes = -1 },
			wantErr: core_errors.ErrInvalidRouteStopVisitDurationMinutes,
		},
		{
			name:    "invalid opening hours",
			mutate:  func(s *RouteStop) { s.OpeningHours = []OpeningHours{{OpensAt: 18, ClosesAt: 9}} },
			wantErr: core_errors.ErrInvalidRouteStopOpeningHours,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stop := validStop()
			if tt.mutate != nil {
				tt.mutate(&stop)
			}

			err := stop.Validate()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestTripCreate_Validate_RealClockBranches(t *testing.T) {
	realNow := time.Now().UTC()
	frozenNow := realNow.Add(-24 * time.Hour)
	timeNow = func() time.Time { return frozenNow }
	t.Cleanup(func() { timeNow = time.Now })

	validCreate := func() TripCreate {
		return TripCreate{
			OwnerID:          1,
			Name:             "Test Trip",
			PlanningArea:     "Moscow",
			StartCoordinates: Coordinates{Latitude: 55.7539, Longitude: 37.6208},
			StartDate:        frozenNow.Add(time.Hour),
			EndDate:          frozenNow.Add(time.Hour),
			DailySchedule:    DailySchedule{StartTime: "09:00", EndTime: "18:00"},
			Budget:           Money{AmountMinorUnits: 1500000, Currency: "RUB"},
			Interests:        []Interest{InterestHistory},
			TimeZone:         "Europe/Moscow",
		}
	}

	t.Run("start date is before the real clock", func(t *testing.T) {
		create := validCreate()
		err := create.Validate()
		require.ErrorIs(t, err, core_errors.ErrInvalidTripStartDate)
	})

	t.Run("end date is before the real clock", func(t *testing.T) {
		create := validCreate()
		create.StartDate = realNow.Add(time.Hour)
		create.EndDate = frozenNow.Add(time.Hour)
		err := create.Validate()
		require.ErrorIs(t, err, core_errors.ErrInvalidTripEndDate)
	})
}

func TestValidatePlanningArea(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{name: "valid area", value: "Moscow"},
		{name: "empty area", value: "", wantErr: core_errors.ErrInvalidPlanningArea},
		{name: "whitespace area", value: "   ", wantErr: core_errors.ErrInvalidPlanningArea},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validatePlanningArea(tt.value)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
