package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func TestRouteCreate_Validate(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour).UTC()

	tests := []struct {
		name    string
		create  RouteCreate
		wantErr error
	}{
		{
			name: "valid route create",
			create: RouteCreate{
				Status: RouteStatusGenerated,
				Days: []RouteDay{
					{
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
			},
		},
		{
			name: "invalid status",
			create: RouteCreate{
				Status: "invalid",
				Days: []RouteDay{
					{
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
					},
				},
			},
			wantErr: core_errors.ErrUnsupportedRouteStatus,
		},
		{
			name: "invalid days - empty",
			create: RouteCreate{
				Status: RouteStatusGenerated,
				Days:   []RouteDay{},
			},
			wantErr: core_errors.ErrInvalidRouteDaysAmount,
		},
		{
			name: "invalid days - first day invalid",
			create: RouteCreate{
				Status: RouteStatusGenerated,
				Days: []RouteDay{
					{
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
				},
			},
			wantErr: core_errors.ErrInvalidID,
		},
		{
			name: "invalid days - invalid date",
			create: RouteCreate{
				Status: RouteStatusGenerated,
				Days: []RouteDay{
					{
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
				},
			},
			wantErr: core_errors.ErrInvalidRouteDayDate,
		},
		{
			name: "invalid days - invalid day number",
			create: RouteCreate{
				Status: RouteStatusGenerated,
				Days: []RouteDay{
					{
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
				},
			},
			wantErr: core_errors.ErrInvalidRouteDayDayNumber,
		},
		{
			name: "invalid days - invalid stops amount",
			create: RouteCreate{
				Status: RouteStatusGenerated,
				Days: []RouteDay{
					{
						ID:        1,
						Date:      futureDate,
						DayNumber: 1,
						Stops:     []RouteStop{},
					},
				},
			},
			wantErr: core_errors.ErrInvalidRouteDayStopsAmount,
		},
		{
			name: "invalid days - invalid leg stop reference",
			create: RouteCreate{
				Status: RouteStatusGenerated,
				Days: []RouteDay{
					{
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
				},
			},
			wantErr: core_errors.ErrInvalidRouteLegStopReference,
		},
		{
			name: "invalid days - duplicate stop ID",
			create: RouteCreate{
				Status: RouteStatusGenerated,
				Days: []RouteDay{
					{
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
								ID:            1,
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
					},
				},
			},
			wantErr: core_errors.ErrDuplicateRouteStopID,
		},
		{
			name: "valid route create with edited status",
			create: RouteCreate{
				Status: RouteStatusEdited,
				Days: []RouteDay{
					{
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
					},
				},
			},
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
