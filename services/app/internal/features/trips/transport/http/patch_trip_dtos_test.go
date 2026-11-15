package trips_transport_http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_http_types "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/types"
)

func TestPatchTripRequest_Validate(t *testing.T) {
	t.Run("returns nil when no field is set", func(t *testing.T) {
		request := PatchTripRequest{}

		require.NoError(t, request.Validate())
	})

	t.Run("returns nil when start and end dates are valid", func(t *testing.T) {
		request := PatchTripRequest{
			StartDate: core_http_types.Nullable[string]{Value: &[]string{"2027-03-10"}[0], Set: true},
			EndDate:   core_http_types.Nullable[string]{Value: &[]string{"2027-03-11"}[0], Set: true},
		}

		require.NoError(t, request.Validate())
	})

	t.Run("returns invalid start date error when start date is invalid", func(t *testing.T) {
		request := PatchTripRequest{
			StartDate: core_http_types.Nullable[string]{Value: &[]string{"not-a-date"}[0], Set: true},
		}

		err := request.Validate()

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidTripStartDate)
	})

	t.Run("returns invalid end date error when end date is invalid", func(t *testing.T) {
		request := PatchTripRequest{
			EndDate: core_http_types.Nullable[string]{Value: &[]string{"not-a-date"}[0], Set: true},
		}

		err := request.Validate()

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidTripEndDate)
	})
}

func TestPatchTripDomainFromDTO(t *testing.T) {
	t.Run("returns unset patch when no field is set", func(t *testing.T) {
		tripPatch, err := patchTripDomainFromDTO(1, PatchTripRequest{})

		require.NoError(t, err)
		require.NotNil(t, tripPatch)
		assert.Equal(t, 1, tripPatch.OwnerID)
		assert.False(t, tripPatch.Name.Set)
		assert.False(t, tripPatch.StartDate.Set)
		assert.False(t, tripPatch.EndDate.Set)
		assert.False(t, tripPatch.MealRequirements.Set)
		assert.False(t, tripPatch.TimeZone.Set)
	})

	t.Run("parses set start and end dates", func(t *testing.T) {
		startDate := "2027-03-10"
		endDate := "2027-03-11"
		tripPatch, err := patchTripDomainFromDTO(1, PatchTripRequest{
			StartDate: core_http_types.Nullable[string]{Value: &startDate, Set: true},
			EndDate:   core_http_types.Nullable[string]{Value: &endDate, Set: true},
		})

		require.NoError(t, err)
		require.NotNil(t, tripPatch)
		require.NotNil(t, tripPatch.StartDate.Value)
		assert.True(t, tripPatch.StartDate.Value.Equal(time.Date(2027, 3, 10, 0, 0, 0, 0, time.UTC)))
		require.NotNil(t, tripPatch.EndDate.Value)
		assert.True(t, tripPatch.EndDate.Value.Equal(time.Date(2027, 3, 11, 0, 0, 0, 0, time.UTC)))
		assert.True(t, tripPatch.StartDate.Set)
		assert.True(t, tripPatch.EndDate.Set)
	})

	t.Run("returns invalid start date error when start date cannot be parsed", func(t *testing.T) {
		startDate := "not-a-date"
		tripPatch, err := patchTripDomainFromDTO(1, PatchTripRequest{
			StartDate: core_http_types.Nullable[string]{Value: &startDate, Set: true},
		})

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidTripStartDate)
		assert.Nil(t, tripPatch)
	})

	t.Run("returns invalid start date error when start date is set to null", func(t *testing.T) {
		tripPatch, err := patchTripDomainFromDTO(1, PatchTripRequest{
			StartDate: core_http_types.Nullable[string]{Value: nil, Set: true},
		})

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidTripStartDate)
		assert.Nil(t, tripPatch)
	})

	t.Run("returns invalid end date error when end date is set to null", func(t *testing.T) {
		tripPatch, err := patchTripDomainFromDTO(1, PatchTripRequest{
			EndDate: core_http_types.Nullable[string]{Value: nil, Set: true},
		})

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidTripEndDate)
		assert.Nil(t, tripPatch)
	})

	t.Run("returns invalid argument error when meal requirements are set to null", func(t *testing.T) {
		tripPatch, err := patchTripDomainFromDTO(1, PatchTripRequest{
			MealRequirements: core_http_types.Nullable[[]MealRequirementRequestDTO]{Value: nil, Set: true},
		})

		require.Error(t, err)
		require.ErrorIs(t, err, core_errors.ErrInvalidArgument)
		assert.Nil(t, tripPatch)
	})

	t.Run("maps meal requirements when they are set", func(t *testing.T) {
		durationMinutes := 60
		budget := domain.Money{AmountMinorUnits: 50000, Currency: "RUB"}
		mealRequirements := []MealRequirementRequestDTO{
			{
				Type:     domain.MealBreakfast,
				Window:   domain.TimeWindow{Start: domain.LocalTime("08:00"), End: domain.LocalTime("09:00")},
				Duration: durationMinutes,
				Budget:   budget,
			},
		}

		tripPatch, err := patchTripDomainFromDTO(1, PatchTripRequest{
			MealRequirements: core_http_types.Nullable[[]MealRequirementRequestDTO]{Value: &mealRequirements, Set: true},
		})

		require.NoError(t, err)
		require.NotNil(t, tripPatch)
		require.NotNil(t, tripPatch.MealRequirements.Value)
		require.Len(t, *tripPatch.MealRequirements.Value, 1)

		converted := (*tripPatch.MealRequirements.Value)[0]
		assert.Equal(t, domain.MealBreakfast, converted.Type)
		assert.Equal(t, 60*time.Minute, converted.Duration)
		assert.Equal(t, budget, converted.Budget)
	})
}
