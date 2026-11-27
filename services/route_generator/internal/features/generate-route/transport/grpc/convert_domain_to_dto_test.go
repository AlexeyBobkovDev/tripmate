package transport_grpc_generate_route

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

func TestConvertRouteDomainToDTO(t *testing.T) {
	route := &domain.Route{
		Status: domain.RouteStatusGenerated,
		Days: []domain.RouteDay{
			{
				ID:        1,
				Date:      time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC),
				DayNumber: 1,
				Stops: []domain.RouteStop{
					{
						PlaceID:              1,
						Position:             1,
						Name:                 "Kremlin",
						Description:          "Kremlin",
						Coordinates:          domain.Coordinates{Latitude: 56.07981, Longitude: 46.192347},
						ArrivalTime:          time.Date(2030, time.May, 10, 11, 0, 0, 0, time.UTC),
						DepartureTime:        time.Date(2030, time.May, 10, 12, 0, 0, 0, time.UTC),
						VisitDurationMinutes: 3600,
					},
					{
						PlaceID:              2,
						Position:             2,
						Name:                 "Red Square",
						Description:          "Red Square",
						Coordinates:          domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:          time.Date(2030, time.May, 10, 14, 0, 0, 0, time.UTC),
						DepartureTime:        time.Date(2030, time.May, 10, 15, 0, 0, 0, time.UTC),
						VisitDurationMinutes: 3600,
					},
				},
				Legs: []domain.RouteLeg{
					{
						ID:         1,
						FromStopID: 1,
						ToStopID:   2,
						TravelMode: domain.TravelModeWalking,
						Distance:   1000,
						Duration:   7200 * time.Second,
					},
				},
			},
		},
	}

	response := convertRouteDomainToDTO(route)

	require.NotNil(t, response)
	assert.Equal(t, route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED, response.Status)
	require.Len(t, response.Days, 1)

	day := response.Days[0]
	require.NotNil(t, day)
	assert.Equal(t, int64(1), day.GetId())
	require.NotNil(t, day.GetDate())
	assert.Equal(t, int32(1), day.GetDayNumber())
	require.Len(t, day.GetStops(), 2)
	require.Len(t, day.GetLegs(), 1)

	stop := day.GetStops()[0]
	require.NotNil(t, stop)

	assert.Equal(t, int64(1), stop.GetPlaceId())
	assert.Equal(t, int32(1), stop.GetPosition())
	assert.Equal(t, "Kremlin", stop.GetName())
	assert.Equal(t, "Kremlin", stop.GetDescription())
	require.NotNil(t, stop.GetCoordinates)
	assert.Equal(t, 56.07981, stop.GetCoordinates().GetLatitude())
	assert.Equal(t, 46.192347, stop.GetCoordinates().GetLongitude())
	require.NotNil(t, stop.GetArrivalTime())
	require.NotNil(t, stop.GetDepartureTime())
	assert.Equal(t, int32(3600), stop.GetVisitDurationMinutes())

	leg := day.GetLegs()[0]
	require.NotNil(t, leg)
	assert.Equal(t, int64(1), leg.GetId())
	assert.Equal(t, int32(1), leg.GetFromStopId())
	assert.Equal(t, int32(2), leg.GetToStopId())
	assert.Equal(t, route_generator_api_v1.TravelMode_TRAVEL_MODE_WALKING, leg.GetTravelMode())
	assert.Equal(t, int32(1000), leg.GetDistance())
	require.NotNil(t, leg.GetDuration())
	assert.Equal(t, 7200.0, leg.GetDuration().AsDuration().Seconds())
}

func TestConvertRouteDomainToDTO_EditedStatus(t *testing.T) {
	route := &domain.Route{
		Status: domain.RouteStatusEdited,
		Days:   []domain.RouteDay{},
	}

	response := convertRouteDomainToDTO(route)

	require.NotNil(t, response)
	assert.Equal(t, route_generator_api_v1.RouteStatus_ROUTE_STATUS_EDITED, response.Status)
	assert.Empty(t, response.GetDays())
}

func TestConvertDaysFromDomainToDTO(t *testing.T) {
	days := []domain.RouteDay{
		{
			ID:        7,
			Date:      time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC),
			DayNumber: 3,
			Stops:     []domain.RouteStop{},
			Legs:      []domain.RouteLeg{},
		},
	}

	converted := convertDaysFromDomainToDTO(days)
	require.Len(t, converted, 1)
	require.NotNil(t, converted[0])
	assert.Equal(t, int64(7), converted[0].GetId())
	assert.Equal(t, int32(3), converted[0].GetDayNumber())
	assert.Empty(t, converted[0].GetStops())
	assert.Empty(t, converted[0].GetLegs())
}

func TestConvertLegsFromDomainToDTO(t *testing.T) {
	legs := []domain.RouteLeg{
		{
			ID:         5,
			FromStopID: 1,
			ToStopID:   2,
			TravelMode: domain.TravelModeCycling,
			Distance:   1500,
			Duration:   1800 * time.Second,
		},
	}

	converted := convertLegsFromDomainToDTO(legs)
	require.Len(t, converted, 1)
	require.NotNil(t, converted[0])
	assert.Equal(t, int64(5), converted[0].GetId())
	assert.Equal(t, int32(1), converted[0].GetFromStopId())
	assert.Equal(t, int32(2), converted[0].GetToStopId())
	assert.Equal(t, route_generator_api_v1.TravelMode_TRAVEL_MODE_CYCLING, converted[0].GetTravelMode())
	assert.Equal(t, int32(1500), converted[0].GetDistance())
	require.NotNil(t, converted[0].GetDuration())
	assert.Equal(t, 1800.0, converted[0].GetDuration().AsDuration().Seconds())
}

func TestConvertStopsFromDomainToDTO(t *testing.T) {
	stops := []domain.RouteStop{
		{
			PlaceID:              42,
			Position:             1,
			Name:                 "Moscow Kremlin",
			Description:          "Kremlin",
			Coordinates:          domain.Coordinates{Latitude: 56.07981, Longitude: 46.192347},
			ArrivalTime:          time.Date(2030, time.May, 10, 11, 0, 0, 0, time.UTC),
			DepartureTime:        time.Date(2030, time.May, 10, 12, 0, 0, 0, time.UTC),
			VisitDurationMinutes: 3600,
		},
	}

	converted := convertStopsFromDomainToDTO(stops)
	require.Len(t, converted, 1)
	require.NotNil(t, converted[0])
	stop := converted[0]

	assert.Equal(t, int64(42), stop.GetPlaceId())
	assert.Equal(t, int32(1), stop.GetPosition())
	assert.Equal(t, "Moscow Kremlin", stop.GetName())
	assert.Equal(t, "Kremlin", stop.GetDescription())
	require.NotNil(t, stop.GetCoordinates)
	assert.Equal(t, 56.07981, stop.GetCoordinates().GetLatitude())
	assert.Equal(t, 46.192347, stop.GetCoordinates().GetLongitude())
	require.NotNil(t, stop.GetArrivalTime())
	require.NotNil(t, stop.GetDepartureTime())
	assert.Equal(t, int32(3600), stop.GetVisitDurationMinutes())
}

func TestConvertStartCoordinatesDomainToDTO(t *testing.T) {
	coordinates := convertStartCoordinatesDomainToDTO(domain.Coordinates{Latitude: 55.75, Longitude: 37.62})

	require.NotNil(t, coordinates)
	assert.Equal(t, 55.75, coordinates.GetLatitude())
	assert.Equal(t, 37.62, coordinates.GetLongitude())
}

func TestConvertDayFromDomainToDTO(t *testing.T) {
	day := convertDayFromDomainToDTO(domain.RouteDay{
		ID:        3,
		Date:      time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC),
		DayNumber: 2,
		Stops:     []domain.RouteStop{},
		Legs:      []domain.RouteLeg{},
	})

	require.NotNil(t, day)
	assert.Equal(t, int64(3), day.GetId())
	require.NotNil(t, day.GetDate())
	assert.Equal(t, int32(2), day.GetDayNumber())
	assert.Empty(t, day.GetStops())
	assert.Empty(t, day.GetLegs())
}
