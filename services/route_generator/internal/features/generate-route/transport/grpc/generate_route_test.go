package transport_grpc_generate_route

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	route_generator_common_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/common/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

type fakeGenerateRouteService struct {
	generateRouteFunc func(ctx context.Context, trip *domain.Trip) (*domain.Route, error)
	called            bool
}

func (f *fakeGenerateRouteService) GenerateRoute(
	ctx context.Context,
	trip *domain.Trip,
) (*domain.Route, error) {
	f.called = true
	return f.generateRouteFunc(ctx, trip)
}

func newTripRequest() *route_generator_api_v1.GenerateRouteRequest {
	startDate := time.Date(2030, time.May, 10, 9, 0, 0, 0, time.UTC)
	endDate := startDate
	return &route_generator_api_v1.GenerateRouteRequest{
		Name:             "Trip to Moscow",
		PlanningArea:     "Moscow",
		StartCoordinates: &route_generator_common_v1.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
		StartDate:        timestamppb.New(startDate),
		EndDate:          timestamppb.New(endDate),
		DailySchedule: &route_generator_common_v1.DailySchedule{
			StartTime: timestamppb.New(startDate),
			EndTime:   timestamppb.New(endDate),
		},
		Budget:   &route_generator_common_v1.Money{AmountMinorUnits: 1500000, Currency: "RUB"},
		TimeZone: "Europe/Moscow",
	}
}

func TestGenerateRouteGRPCTransport_GenerateRoute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fakeService := &fakeGenerateRouteService{
			generateRouteFunc: func(ctx context.Context, trip *domain.Trip) (*domain.Route, error) {
				require.NotNil(t, trip)
				require.NotZero(t, trip.StartDate)
				return &domain.Route{
					Status: domain.RouteStatusGenerated,
					Days: []domain.RouteDay{
						{
							ID:        1,
							Date:      trip.StartDate,
							DayNumber: 1,
							Stops: []domain.RouteStop{
								{
									PlaceID:       1,
									Position:      1,
									Name:          "Kremlin",
									Description:   "Kremlin",
									Coordinates:   domain.Coordinates{Latitude: 56.07981, Longitude: 46.192347},
									ArrivalTime:   trip.StartDate.Add(2 * time.Hour),
									DepartureTime: trip.StartDate.Add(3 * time.Hour),
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
				}, nil
			},
		}
		transport := NewGenerateRouteServer(fakeService)

		response, err := transport.GenerateRoute(context.Background(), newTripRequest())
		require.NoError(t, err)
		require.NotNil(t, response)

		assert.Equal(t, route_generator_api_v1.RouteStatus_ROUTE_STATUS_GENERATED, response.Status)
		require.Len(t, response.Days, 1)
		require.Len(t, response.Days[0].Stops, 1)
		require.Len(t, response.Days[0].Legs, 1)
	})

	t.Run("service invalid argument error is propagated", func(t *testing.T) {
		fakeService := &fakeGenerateRouteService{
			generateRouteFunc: func(ctx context.Context, trip *domain.Trip) (*domain.Route, error) {
				return nil, status.Error(codes.InvalidArgument, "invalid trip data")
			},
		}
		transport := NewGenerateRouteServer(fakeService)

		response, err := transport.GenerateRoute(context.Background(), newTripRequest())
		require.Error(t, err)
		require.Nil(t, response)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("service internal error is propagated", func(t *testing.T) {
		fakeService := &fakeGenerateRouteService{
			generateRouteFunc: func(ctx context.Context, trip *domain.Trip) (*domain.Route, error) {
				return nil, status.Error(codes.Internal, "failed to generate route")
			},
		}
		transport := NewGenerateRouteServer(fakeService)

		response, err := transport.GenerateRoute(context.Background(), newTripRequest())
		require.Error(t, err)
		require.Nil(t, response)
		require.Equal(t, codes.Internal, status.Code(err))
	})
}
