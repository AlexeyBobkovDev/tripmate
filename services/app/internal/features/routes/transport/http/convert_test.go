package routes_transport_http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func TestRouteDTOFromDomain_MapsRouteFields(t *testing.T) {
	route := domain.Route{
		ID:          7,
		TripID:      42,
		TripVersion: 3,
		Version:     5,
		Status:      domain.RouteStatusGenerated,
		Days: []domain.RouteDay{
			{
				DayNumber: 1,
				Legs: []domain.RouteLeg{
					{
						FromStopID: 1,
						ToStopID:   2,
						TravelMode: domain.TravelModeWalking,
						Distance:   800,
						Duration:   600,
					},
				},
			},
		},
	}

	dto := routeDTOFromDomain(&route)

	assert.Equal(t, route.ID, dto.ID)
	assert.Equal(t, route.TripID, dto.TripID)
	assert.Equal(t, route.TripVersion, dto.TripVersion)
	assert.Equal(t, route.Version, dto.Version)
	assert.Equal(t, route.Status, dto.Status)
	assert.Equal(t, route.CreatedAt, dto.CreatedAt)
	assert.Equal(t, route.UpdatedAt, dto.UpdatedAt)
	require.Len(t, dto.Days, 1)
	require.Len(t, dto.Days[0].Legs, 1)
	assert.Equal(t, route.Days[0].Legs[0].FromStopID, dto.Days[0].Legs[0].FromStopID)
	assert.Equal(t, route.Days[0].Legs[0].ToStopID, dto.Days[0].Legs[0].ToStopID)
	assert.Equal(t, route.Days[0].Legs[0].TravelMode, dto.Days[0].Legs[0].TravelMode)
	assert.Equal(t, route.Days[0].Legs[0].Distance, dto.Days[0].Legs[0].Distance)
	assert.Equal(t, route.Days[0].Legs[0].Duration, dto.Days[0].Legs[0].Duration)
}

func TestRouteDTOFromDomain_EmptyDaysStaysEmpty(t *testing.T) {
	route := testRouteResponse()
	route.Days = []domain.RouteDay{}

	dto := routeDTOFromDomain(&route)

	require.NotNil(t, dto.Days)
	require.Empty(t, dto.Days)
}

func TestRouteDTOFromDomain_NilDays(t *testing.T) {
	route := testRouteResponse()
	route.Days = nil

	dto := routeDTOFromDomain(&route)

	assert.Nil(t, dto.Days)
}
