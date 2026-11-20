package generate_route_service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

type fakeGenerateRouteRepository struct {
	called bool
}

func (f *fakeGenerateRouteRepository) FindCandidates(
	ctx context.Context,
	criteria FindCandidatesCriteria,
) ([]domain.Route, error) {
	f.called = true
	return []domain.Route{}, nil
}

func newValidTrip() *domain.Trip {
	startDate := time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC)
	return &domain.Trip{
		Name:             "Trip to Moscow",
		PlanningArea:     "Moscow",
		StartCoordinates: domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		StartDate:        startDate,
		EndDate:          startDate,
		DailySchedule:    domain.DailySchedule{StartTime: "09:00", EndTime: "18:00"},
		MealRequirements: []domain.MealRequirement{},
		Budget:           domain.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		Interests:        []domain.Interest{domain.InterestHistory},
		TimeZone:         "Europe/Moscow",
	}
}

func TestGenerateRouteService_GenerateRoute(t *testing.T) {
	service := NewGenerateRouteService(&fakeGenerateRouteRepository{})

	tests := []struct {
		name     string
		trip     *domain.Trip
		wantErr  error
		wantCode codes.Code
	}{
		{
			name:     "valid trip",
			trip:     newValidTrip(),
			wantErr:  nil,
			wantCode: codes.OK,
		},
		{
			name:     "empty trip name",
			trip:     func() *domain.Trip { trip := newValidTrip(); trip.Name = ""; return trip }(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "blank planning area",
			trip:     func() *domain.Trip { trip := newValidTrip(); trip.PlanningArea = "   "; return trip }(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid start coordinates",
			trip: func() *domain.Trip {
				trip := newValidTrip()
				trip.StartCoordinates = domain.Coordinates{Latitude: 91, Longitude: 37.62}
				return trip
			}(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "zero trip start date",
			trip:     func() *domain.Trip { trip := newValidTrip(); trip.StartDate = time.Time{}; return trip }(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
		{
			name: "past trip start date",
			trip: func() *domain.Trip {
				trip := newValidTrip()
				trip.StartDate = time.Now().Add(-24 * time.Hour)
				trip.EndDate = trip.StartDate
				return trip
			}(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
		{
			name: "trip end date before start date",
			trip: func() *domain.Trip {
				trip := newValidTrip()
				trip.StartDate = time.Now().Add(48 * time.Hour)
				trip.EndDate = time.Now().Add(24 * time.Hour)
				return trip
			}(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid daily schedule",
			trip: func() *domain.Trip {
				trip := newValidTrip()
				trip.DailySchedule = domain.DailySchedule{StartTime: "18:00", EndTime: "09:00"}
				return trip
			}(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid meal requirement duration",
			trip: func() *domain.Trip {
				trip := newValidTrip()
				trip.MealRequirements = []domain.MealRequirement{{Type: domain.MealBreakfast, Window: domain.TimeWindow{Start: "08:00", End: "09:00"}, Duration: 0, Budget: domain.Money{AmountMinorUnits: 1000, Currency: "RUB"}}}
				return trip
			}(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid trip budget",
			trip: func() *domain.Trip {
				trip := newValidTrip()
				trip.Budget = domain.Money{AmountMinorUnits: 0, Currency: "RUB"}
				return trip
			}(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "unsupported time zone",
			trip:     func() *domain.Trip { trip := newValidTrip(); trip.TimeZone = "Mars/Olympus"; return trip }(),
			wantErr:  status.Error(codes.InvalidArgument, "invalid trip data"),
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route, err := service.GenerateRoute(context.Background(), tt.trip)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Nil(t, route)
				require.Equal(t, tt.wantCode, status.Code(err))
			} else {
				require.NoError(t, err)
				require.NotNil(t, route)
			}
		})
	}
}

func TestGenerateRouteService_ValidTrip_ReturnsGeneratedRoute(t *testing.T) {
	service := NewGenerateRouteService(&fakeGenerateRouteRepository{})

	route, err := service.GenerateRoute(context.Background(), newValidTrip())
	require.NoError(t, err)
	require.NotNil(t, route)

	assert.Equal(t, domain.RouteStatus(domain.RouteStatusGenerated), route.Status)
	require.Len(t, route.Days, 1)
	require.Len(t, route.Days[0].Stops, 2)
	require.Len(t, route.Days[0].Legs, 1)

	day := route.Days[0]
	assert.Equal(t, 1, day.ID)
	assert.Equal(t, 1, day.DayNumber)
	assert.NotZero(t, day.Date)
}

func TestGenerateRouteService_GeneratedRouteUsesTripStartDate(t *testing.T) {
	service := NewGenerateRouteService(&fakeGenerateRouteRepository{})
	trip := newValidTrip()
	trip.StartDate = time.Date(2031, time.June, 15, 9, 0, 0, 0, time.UTC)
	trip.EndDate = trip.StartDate

	route, err := service.GenerateRoute(context.Background(), trip)
	require.NoError(t, err)
	require.NotNil(t, route)

	// The first generated day must be placed on the trip's start date.
	require.Len(t, route.Days, 1)
	assert.Equal(t, trip.StartDate.Format(core_consts.Layout), route.Days[0].Date.Format(core_consts.Layout))
}
