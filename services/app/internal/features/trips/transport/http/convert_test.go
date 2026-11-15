package trips_transport_http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

// TestTripDTOFromDomain_MealDurationConversionDoesNotMutateSource_BUG
// documents that tripDTOFromDomain divides the duration in place. The domain
// slice is shared, so every extra conversion divides the duration again
// (90 minutes -> 90 nanoseconds -> 0). Converting the same trip twice must be
// idempotent and must not change the domain object.
func TestTripDTOFromDomain_MealDurationConversionDoesNotMutateSource_BUG(t *testing.T) {
	trip := testTripResponse()
	trip.MealRequirements = []domain.MealRequirement{
		{
			Type:     domain.MealBreakfast,
			Window:   domain.TimeWindow{Start: "08:00", End: "10:00"},
			Duration: 90 * time.Minute,
			Budget:   domain.Money{AmountMinorUnits: 50000, Currency: "RUB"},
		},
	}

	originalDuration := trip.MealRequirements[0].Duration

	first := tripDTOFromDomain(&trip)
	second := tripDTOFromDomain(&trip)

	require.Len(t, first.MealRequirements, 1)
	require.Len(t, second.MealRequirements, 1)

	// The duration should be preserved (90 minutes = 1 hour 30 minutes)
	expectedDuration := 90 * time.Minute
	assert.Equal(t, expectedDuration, first.MealRequirements[0].Duration)
	assert.Equal(t, expectedDuration, second.MealRequirements[0].Duration)
	assert.Equal(t, originalDuration, trip.MealRequirements[0].Duration)
}

func TestTripDTOFromDomain_MapsTripFields(t *testing.T) {
	trip := testTripResponse()

	response := tripDTOFromDomain(&trip)

	assert.Equal(t, trip.ID, response.ID)
	assert.Equal(t, trip.OwnerID, response.OwnerID)
	assert.Equal(t, trip.Version, response.Version)
	assert.Equal(t, trip.Name, response.Name)
	assert.Equal(t, trip.PlanningArea, response.PlanningArea)
	assert.Equal(t, trip.StartCoordinates, response.StartCoordinates)
	assert.Equal(t, trip.StartDate, response.StartDate)
	assert.Equal(t, trip.EndDate, response.EndDate)
	assert.Equal(t, trip.DailySchedule, response.DailySchedule)
	assert.Equal(t, trip.Budget, response.Budget)
	assert.Equal(t, trip.Interests, response.Interests)
	assert.Equal(t, trip.TimeZone, response.TimeZone)
	assert.Equal(t, trip.CreatedAt, response.CreatedAt)
	assert.Equal(t, trip.UpdatedAt, response.UpdatedAt)
	assert.Equal(t, trip.DeletedAt, response.DeletedAt)
}

func TestTripDTOsFromDomains_MapsEachTrip(t *testing.T) {
	first := testTripResponse()
	second := testTripResponse()
	second.ID = 2
	second.OwnerID = 2

	trips := []domain.Trip{first, second}

	responses := tripDTOsFromDomains(&trips)

	require.Len(t, responses, 2)
	assert.Equal(t, first.ID, responses[0].ID)
	assert.Equal(t, second.ID, responses[1].ID)
}
