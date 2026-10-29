package domain

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func validTripCreate() *TripCreate {
	now := time.Now()
	startDate := now.AddDate(0, 0, 30).UTC()
	endDate := startDate.AddDate(0, 0, 3)

	return NewTripCreate(
		"Weekend in Moscow",
		"Moscow",
		Location{
			Name:      "Red Square",
			Latitude:  55.7539,
			Longitude: 37.6208,
		},
		startDate,
		endDate,
		Money{
			AmountMinorUnits: 100000,
			Currency:         "RUB",
		},
		[]Interest{InterestHistory, InterestFood},
		"Europe/Moscow",
	)
}

func TestTripCreateValidate_Valid(t *testing.T) {
	t.Parallel()

	t.Run("all fields provided", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		require.NoError(t, got.Validate())
	})

	t.Run("optional location and budget omitted", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.StartLocation = Location{}
		got.Budget = Money{}

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidLocationName)
	})

	t.Run("empty interests", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.Interests = nil

		require.NoError(t, got.Validate())
	})
}

func TestTripCreateValidate_InvalidName(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.Name = ""

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripName)
}

func TestTripCreateValidate_InvalidPlanningArea(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.PlanningArea = ""

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidPlanningArea)
}

func TestTripCreateValidate_InvalidStartLocation(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.StartLocation = Location{
		Name:      "Invalid location",
		Latitude:  math.NaN(),
		Longitude: 37.6208,
	}

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidLatitude)
}

func TestTripCreateValidate_InvalidStartDate(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.StartDate = time.Now().AddDate(0, 0, -1)
	got.EndDate = got.StartDate.AddDate(0, 0, 3)

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripStartDate)
}

func TestTripCreateValidate_InvalidEndDate(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.EndDate = time.Now().AddDate(0, 0, -1)

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripEndDate)
}

func TestTripCreateValidate_StartDateMustPrecedeEndDate(t *testing.T) {
	t.Parallel()

	t.Run("start date after end date", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.StartDate = time.Now().AddDate(0, 0, 40)
		got.EndDate = time.Now().AddDate(0, 0, 39)

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripEndDate)
	})

	t.Run("start date equals end date", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.EndDate = got.StartDate

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripEndDate)
	})
}

func TestTripCreateValidate_InvalidBudget(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.Budget = Money{
		AmountMinorUnits: 0,
		Currency:         "RUB",
	}

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidBudget)
}

func TestTripCreateValidate_InvalidInterests(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.Interests = []Interest{"unsupported-interest"}

	require.ErrorIs(t, got.Validate(), core_errors.ErrUnsupportedInterest)
}

func TestTripCreateValidate_InvalidTimeZone(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.TimeZone = "Not/A_Real_Time_Zone"

	require.ErrorIs(t, got.Validate(), core_errors.ErrUnsupportedTimezone)
}

func TestNewTripCreate(t *testing.T) {
	t.Parallel()

	startLocation := Location{
		Name:      "Red Square",
		Latitude:  55.7539,
		Longitude: 37.6208,
	}
	startDate := time.Now().AddDate(0, 0, 60).UTC()
	endDate := startDate.AddDate(0, 0, 3)
	budget := Money{AmountMinorUnits: 100000, Currency: "RUB"}
	interests := []Interest{InterestHistory, InterestNature}

	got := NewTripCreate(
		"Winter trip",
		"Moscow",
		startLocation,
		startDate,
		endDate,
		budget,
		interests,
		"Europe/Moscow",
	)

	require.NotNil(t, got)
	assert.Equal(t, "Winter trip", got.Name)
	assert.Equal(t, "Moscow", got.PlanningArea)
	assert.Equal(t, startLocation, got.StartLocation)
	assert.True(t, startDate.Equal(got.StartDate))
	assert.True(t, endDate.Equal(got.EndDate))
	assert.Equal(t, budget, got.Budget)
	assert.Equal(t, interests, got.Interests)
	assert.Equal(t, "Europe/Moscow", got.TimeZone)
}
