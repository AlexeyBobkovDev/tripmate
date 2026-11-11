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
	startDate := now.AddDate(0, 0, 10).UTC()
	endDate := startDate

	return NewTripCreate(
		"Weekend in Moscow",
		"Moscow",
		Coordinates{
			Latitude:  55.7539,
			Longitude: 37.6208,
		},
		startDate,
		endDate,
		validDailySchedule(),
		[]MealRequirement{validMealRequirement()},
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

	t.Run("optional location omitted", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.StartCoordinates = Coordinates{}

		require.NoError(t, got.Validate())
	})

	t.Run("empty interests", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.Interests = nil

		require.NoError(t, got.Validate())
	})

	t.Run("empty meal requirements", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.MealRequirement = nil

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

func TestTripCreateValidate_InvalidStartCoordinates(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.StartCoordinates = Coordinates{
		Latitude:  math.NaN(),
		Longitude: 37.6208,
	}

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidLatitude)
}

func TestTripCreateValidate_InvalidStartDate(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.StartDate = time.Now().AddDate(0, 0, -1)
	got.EndDate = got.StartDate

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripStartDate)
}

func TestTripCreateValidate_ZeroStartDate(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.StartDate = time.Time{}
	got.EndDate = got.StartDate

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripStartDate)
}

func TestTripCreateValidate_InvalidEndDate(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.EndDate = time.Now().AddDate(0, 0, -1)

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripEndDate)
}

func TestTripCreateValidate_ZeroEndDate(t *testing.T) {
	t.Parallel()

	got := validTripCreate()
	got.EndDate = time.Time{}

	require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripEndDate)
}

func TestTripCreateValidate_StartDateMustEqualEndDate(t *testing.T) {
	t.Parallel()

	t.Run("start date before end date", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.EndDate = got.StartDate.AddDate(0, 0, 1)

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripEndDate)
	})

	t.Run("start date after end date", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.EndDate = got.StartDate.AddDate(0, 0, -1)

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripEndDate)
	})

	t.Run("start date equals end date", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.EndDate = got.StartDate

		require.NoError(t, got.Validate())
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

func TestTripCreateValidate_InvalidDailySchedule(t *testing.T) {
	t.Parallel()

	t.Run("invalid start time", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.DailySchedule.StartTime = ""

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidLocalTime)
	})

	t.Run("start after end", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.DailySchedule = DailySchedule{StartTime: "18:00", EndTime: "09:00"}

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidTripEndTime)
	})
}

func TestTripCreateValidate_InvalidMealRequirement(t *testing.T) {
	t.Parallel()

	t.Run("unsupported meal type", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.MealRequirement = []MealRequirement{validMealRequirement()}
		got.MealRequirement[0].Type = "brunch"

		require.ErrorIs(t, got.Validate(), core_errors.ErrUnsupportedMealType)
	})

	t.Run("invalid meal window", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.MealRequirement = []MealRequirement{validMealRequirement()}
		got.MealRequirement[0].Window.Start = ""

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidLocalTime)
	})

	t.Run("invalid meal duration", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.MealRequirement = []MealRequirement{validMealRequirement()}
		got.MealRequirement[0].Duration = 0

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidMealTimeDuration)
	})

	t.Run("invalid meal budget", func(t *testing.T) {
		t.Parallel()

		got := validTripCreate()
		got.MealRequirement = []MealRequirement{validMealRequirement()}
		got.MealRequirement[0].Budget = Money{AmountMinorUnits: 0, Currency: "RUB"}

		require.ErrorIs(t, got.Validate(), core_errors.ErrInvalidBudget)
	})
}

func TestNewTripCreate(t *testing.T) {
	t.Parallel()

	startCoordinates := Coordinates{
		Latitude:  55.7539,
		Longitude: 37.6208,
	}
	startDate := time.Now().AddDate(0, 0, 60).UTC()
	endDate := startDate
	budget := Money{AmountMinorUnits: 100000, Currency: "RUB"}
	interests := []Interest{InterestHistory, InterestNature}
	dailySchedule := DailySchedule{StartTime: "08:00", EndTime: "22:00"}
	mealRequirements := []MealRequirement{validMealRequirement()}

	got := NewTripCreate(
		"Winter trip",
		"Moscow",
		startCoordinates,
		startDate,
		endDate,
		dailySchedule,
		mealRequirements,
		budget,
		interests,
		"Europe/Moscow",
	)

	require.NotNil(t, got)
	assert.Equal(t, "Winter trip", got.Name)
	assert.Equal(t, "Moscow", got.PlanningArea)
	assert.Equal(t, startCoordinates, got.StartCoordinates)
	assert.True(t, startDate.Equal(got.StartDate))
	assert.True(t, endDate.Equal(got.EndDate))
	assert.Equal(t, dailySchedule, got.DailySchedule)
	assert.Equal(t, mealRequirements, got.MealRequirement)
	assert.Equal(t, budget, got.Budget)
	assert.Equal(t, interests, got.Interests)
	assert.Equal(t, "Europe/Moscow", got.TimeZone)
}
