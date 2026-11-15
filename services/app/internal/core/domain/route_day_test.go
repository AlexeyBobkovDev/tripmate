package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestRouteDay_Validate(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour).UTC()

	tests := []struct {
		name    string
		day     RouteDay
		wantErr error
	}{
		{
			name: "valid route day",
			day: RouteDay{
				ID:        1,
				Date:      futureDate,
				DayNumber: 1,
				Stops: []RouteStop{
					{
						ID:            1,
						PlaceID:       1,
						Position:      1,
						Name:          "Stop 1",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate,
						DepartureTime: futureDate.Add(2 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
					{
						ID:            2,
						PlaceID:       2,
						Position:      2,
						Name:          "Stop 2",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate.Add(2 * time.Hour),
						DepartureTime: futureDate.Add(4 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
				},
				Legs: []RouteLeg{
					{
						ID:         1,
						FromStopID: 1,
						ToStopID:   2,
						TravelMode: TravelModeWalking,
						Distance:   1000,
						Duration:   15 * time.Minute,
					},
				},
			},
		},
		{
			name: "invalid ID",
			day: RouteDay{
				ID:        -1,
				Date:      futureDate,
				DayNumber: 1,
				Stops: []RouteStop{
					{
						ID:            1,
						PlaceID:       1,
						Position:      1,
						Name:          "Stop 1",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate,
						DepartureTime: futureDate.Add(2 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
				},
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid date",
			day: RouteDay{
				ID:        1,
				Date:      time.Time{},
				DayNumber: 1,
				Stops: []RouteStop{
					{
						ID:            1,
						PlaceID:       1,
						Position:      1,
						Name:          "Stop 1",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate,
						DepartureTime: futureDate.Add(2 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
				},
			},
			wantErr: core_errors.ErrInvalidRouteDayDate,
		},
		{
			name: "invalid day number",
			day: RouteDay{
				ID:        1,
				Date:      futureDate,
				DayNumber: 0,
				Stops: []RouteStop{
					{
						ID:            1,
						PlaceID:       1,
						Position:      1,
						Name:          "Stop 1",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate,
						DepartureTime: futureDate.Add(2 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
				},
			},
			wantErr: core_errors.ErrInvalidRouteDayDayNumber,
		},
		{
			name: "invalid stops amount",
			day: RouteDay{
				ID:        1,
				Date:      futureDate,
				DayNumber: 1,
				Stops:     []RouteStop{},
			},
			wantErr: core_errors.ErrInvalidRouteDayStopsAmount,
		},
		{
			name: "invalid leg stop reference",
			day: RouteDay{
				ID:        1,
				Date:      futureDate,
				DayNumber: 1,
				Stops: []RouteStop{
					{
						ID:            1,
						PlaceID:       1,
						Position:      1,
						Name:          "Stop 1",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate,
						DepartureTime: futureDate.Add(2 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
				},
				Legs: []RouteLeg{
					{
						ID:         1,
						FromStopID: 999,
						ToStopID:   1,
						TravelMode: TravelModeWalking,
						Distance:   1000,
						Duration:   15 * time.Minute,
					},
				},
			},
			wantErr: core_errors.ErrInvalidRouteLegStopReference,
		},
		{
			name: "invalid leg start equals end",
			day: RouteDay{
				ID:        1,
				Date:      futureDate,
				DayNumber: 1,
				Stops: []RouteStop{
					{
						ID:            1,
						PlaceID:       1,
						Position:      1,
						Name:          "Stop 1",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate,
						DepartureTime: futureDate.Add(2 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
					{
						ID:            2,
						PlaceID:       2,
						Position:      2,
						Name:          "Stop 2",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate.Add(2 * time.Hour),
						DepartureTime: futureDate.Add(4 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
				},
				Legs: []RouteLeg{
					{
						ID:         1,
						FromStopID: 1,
						ToStopID:   1,
						TravelMode: TravelModeWalking,
						Distance:   1000,
						Duration:   15 * time.Minute,
					},
				},
			},
			wantErr: core_errors.ErrInvalidRouteLegStopReference,
		},
		{
			name: "invalid leg start not in stops",
			day: RouteDay{
				ID:        1,
				Date:      futureDate,
				DayNumber: 1,
				Stops: []RouteStop{
					{
						ID:            2,
						PlaceID:       1,
						Position:      1,
						Name:          "Stop 1",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate,
						DepartureTime: futureDate.Add(2 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
				},
				Legs: []RouteLeg{
					{
						ID:         1,
						FromStopID: 1,
						ToStopID:   2,
						TravelMode: TravelModeWalking,
						Distance:   1000,
						Duration:   15 * time.Minute,
					},
				},
			},
			wantErr: core_errors.ErrInvalidRouteLegStopReference,
		},
		{
			name: "invalid leg end not in stops",
			day: RouteDay{
				ID:        1,
				Date:      futureDate,
				DayNumber: 1,
				Stops: []RouteStop{
					{
						ID:            1,
						PlaceID:       1,
						Position:      1,
						Name:          "Stop 1",
						Description:   "Description",
						Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
						ArrivalTime:   futureDate,
						DepartureTime: futureDate.Add(2 * time.Hour),
						VisitDuration: 2 * time.Hour,
					},
				},
				Legs: []RouteLeg{
					{
						ID:         1,
						FromStopID: 1,
						ToStopID:   999,
						TravelMode: TravelModeWalking,
						Distance:   1000,
						Duration:   15 * time.Minute,
					},
				},
			},
			wantErr: core_errors.ErrInvalidRouteLegStopReference,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.day.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestRouteDay_ValidationOrder(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour).UTC()

	// Test that validation happens in the correct order
	day := RouteDay{
		ID:        1,
		Date:      futureDate,
		DayNumber: 1,
		Stops: []RouteStop{
			{
				ID:            1,
				PlaceID:       1,
				Position:      1,
				Name:          "Stop 1",
				Description:   "Description",
				Coordinates:   Coordinates{Latitude: 55.7539, Longitude: 37.6208},
				ArrivalTime:   futureDate,
				DepartureTime: futureDate.Add(2 * time.Hour),
				VisitDuration: 2 * time.Hour,
			},
		},
	}

	// Should pass all validations
	require.NoError(t, day.Validate())

	// Now test that invalid ID is caught first
	day.ID = -1
	err := day.Validate()
	assert.ErrorIs(t, err, core_errors.ErrInvalidID)

	// Reset ID and test next validation
	day.ID = 1
	day.Date = time.Time{}
	err = day.Validate()
	assert.ErrorIs(t, err, core_errors.ErrInvalidRouteDayDate)
}
