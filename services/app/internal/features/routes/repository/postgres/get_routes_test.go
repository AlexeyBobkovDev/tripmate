package routes_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestRoutesRepository_GetRoutes_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteRow := mocks.NewMockRow(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	mockRouteStopsRows := mocks.NewMockRows(ctrl)
	mockOpeningHoursRows := mocks.NewMockRows(ctrl)
	mockRouteLegRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	expectedRoutes := []*domain.Route{
		{
			ID:          1,
			TripID:      1,
			TripVersion: 1,
			Version:     1,
			Status:      domain.RouteStatusGenerated,
			Days: []domain.RouteDay{
				{
					ID:        1,
					Date:      time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC),
					DayNumber: 1,
					Stops: []domain.RouteStop{
						{
							ID:                   1,
							PlaceID:              1,
							Position:             1,
							Name:                 "Test Stop",
							Description:          "Test Description",
							Address:              "Test Address",
							WebsiteURL:           nil,
							Phone:                nil,
							VisitDurationMinutes: 60,
							Coordinates:          domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
							ArrivalTime:          time.Date(2026, time.November, 14, 9, 0, 0, 0, time.UTC),
							DepartureTime:        time.Date(2026, time.November, 14, 10, 0, 0, 0, time.UTC),
						},
					},
					Legs: []domain.RouteLeg{},
				},
			},
			CreatedAt: time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, time.November, 14, 0, 0, 0, 0, time.UTC),
		},
	}

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)

	mockPool.EXPECT().Query(gomock.Any(), getRoutesQuery, tripID).Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			route := expectedRoutes[0]
			*destinations[0].(*int) = route.ID
			*destinations[1].(*int) = route.TripID
			*destinations[2].(*int) = route.TripVersion
			*destinations[3].(*int) = route.Version
			*destinations[4].(*domain.RouteStatus) = route.Status
			*destinations[5].(*time.Time) = route.CreatedAt
			*destinations[6].(*time.Time) = route.UpdatedAt
			return nil
		})
	mockRows.EXPECT().Next().Return(false)
	mockRows.EXPECT().Err().Return(nil)
	mockRows.EXPECT().Close()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), getRouteQuery, tripID, expectedRoutes[0].ID).
		Return(mockRouteRow)
	mockRouteRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			route := expectedRoutes[0]
			*destinations[0].(*int) = route.ID
			*destinations[1].(*int) = route.TripID
			*destinations[2].(*int) = route.TripVersion
			*destinations[3].(*int) = route.Version
			*destinations[4].(*domain.RouteStatus) = route.Status
			*destinations[5].(*time.Time) = route.CreatedAt
			*destinations[6].(*time.Time) = route.UpdatedAt
			return nil
		})

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteDayQuery, expectedRoutes[0].ID).
		Return(mockRouteDaysRows, nil)
	mockRouteDaysRows.EXPECT().Next().Return(true)
	mockRouteDaysRows.EXPECT().
		Scan(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(destinations ...any) error {
			day := expectedRoutes[0].Days[0]
			*destinations[0].(*int) = day.ID
			*destinations[1].(*time.Time) = day.Date
			*destinations[2].(*int) = day.DayNumber
			return nil
		})
	mockRouteDaysRows.EXPECT().Next().Return(false)
	mockRouteDaysRows.EXPECT().Err().Return(nil)
	mockRouteDaysRows.EXPECT().Close()

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteStopQuery, expectedRoutes[0].Days[0].ID).
		Return(mockRouteStopsRows, nil)
	mockRouteStopsRows.EXPECT().Next().Return(true)
	mockRouteStopsRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			stop := expectedRoutes[0].Days[0].Stops[0]
			*destinations[0].(*int) = stop.ID
			*destinations[1].(*int) = stop.PlaceID
			*destinations[2].(*int) = stop.Position
			*destinations[3].(*string) = stop.Name
			*destinations[4].(*string) = stop.Description
			*destinations[5].(*[]domain.Category) = stop.Categories
			*destinations[6].(*[]domain.Interest) = stop.Interests
			*destinations[7].(*string) = stop.Address
			*destinations[8].(**string) = stop.WebsiteURL
			*destinations[9].(**string) = stop.Phone
			*destinations[10].(*int) = stop.VisitDurationMinutes
			*destinations[11].(*float64) = stop.Coordinates.Latitude
			*destinations[12].(*float64) = stop.Coordinates.Longitude
			*destinations[13].(*time.Time) = stop.ArrivalTime
			*destinations[14].(*time.Time) = stop.DepartureTime
			return nil
		})

	mockPool.EXPECT().
		Query(gomock.Any(), getOpeningHoursQuery, expectedRoutes[0].Days[0].Stops[0].ID).
		Return(mockOpeningHoursRows, nil)
	mockOpeningHoursRows.EXPECT().Next().Return(false)
	mockOpeningHoursRows.EXPECT().Err().Return(nil)
	mockOpeningHoursRows.EXPECT().Close()

	mockRouteStopsRows.EXPECT().Next().Return(false)
	mockRouteStopsRows.EXPECT().Err().Return(nil)
	mockRouteStopsRows.EXPECT().Close()

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteLegQuery, expectedRoutes[0].Days[0].ID).
		Return(mockRouteLegRows, nil)
	mockRouteLegRows.EXPECT().Next().Return(false)
	mockRouteLegRows.EXPECT().Err().Return(nil)
	mockRouteLegRows.EXPECT().Close()

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)

	require.NoError(t, err)
	require.NotNil(t, routes)
	require.Len(t, routes, 1)
	require.Equal(t, expectedRoutes[0].ID, routes[0].ID)
	require.Equal(t, expectedRoutes[0].TripID, routes[0].TripID)
	require.Equal(t, expectedRoutes[0].TripVersion, routes[0].TripVersion)
	require.Equal(t, expectedRoutes[0].Version, routes[0].Version)
	require.Equal(t, expectedRoutes[0].Status, routes[0].Status)
	require.Len(t, routes[0].Days, 1)
	require.Len(t, routes[0].Days[0].Stops, 1)
	require.Len(t, routes[0].Days[0].Legs, 0)
}

func TestRoutesRepository_GetRoutes_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Query(gomock.Any(), getRoutesQuery, gomock.Any()).Return(nil, errors.New("connection refused"))

	routes, err := routesRepository.GetRoutes(context.Background(), 1)

	require.Error(t, err)
	assert.Nil(t, routes)
	assert.ErrorContains(t, err, "failed to fetch route rows")
}

func TestRoutesRepository_GetRoutes_RowsScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Query(gomock.Any(), getRoutesQuery, gomock.Any()).Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().Scan(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(errors.New("failed to scan route"))
	mockRows.EXPECT().Close()

	routes, err := routesRepository.GetRoutes(context.Background(), 1)

	require.Error(t, err)
	assert.Nil(t, routes)
	assert.ErrorContains(t, err, "failed to fetch route")
}

func TestRoutesRepository_GetRoutes_RowsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Query(gomock.Any(), getRoutesQuery, gomock.Any()).Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(false)
	mockRows.EXPECT().Err().Return(errors.New("rows iteration failed"))
	mockRows.EXPECT().Close()

	routes, err := routesRepository.GetRoutes(context.Background(), 1)

	require.Error(t, err)
	assert.Nil(t, routes)
	assert.ErrorContains(t, err, "error while fetching routes rows")
}

func TestRoutesRepository_GetRoutes_RouteLegScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteRow := mocks.NewMockRow(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	mockRouteStopsRows := mocks.NewMockRows(ctrl)
	mockRouteLegRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	legScanError := errors.New("route leg scan error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)

	mockPool.EXPECT().Query(gomock.Any(), getRoutesQuery, tripID).Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(nil)
	mockRows.EXPECT().Next().Return(false)
	mockRows.EXPECT().Err().Return(nil)
	mockRows.EXPECT().Close()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), getRouteQuery, tripID, gomock.Any()).
		Return(mockRouteRow)
	mockRouteRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(nil)

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteDayQuery, gomock.Any()).
		Return(mockRouteDaysRows, nil)
	mockRouteDaysRows.EXPECT().Next().Return(true)
	mockRouteDaysRows.EXPECT().
		Scan(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	mockRouteDaysRows.EXPECT().Next().Return(false)
	mockRouteDaysRows.EXPECT().Err().Return(nil)
	mockRouteDaysRows.EXPECT().Close()

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteStopQuery, gomock.Any()).
		Return(mockRouteStopsRows, nil)
	mockRouteStopsRows.EXPECT().Next().Return(false)
	mockRouteStopsRows.EXPECT().Err().Return(nil)
	mockRouteStopsRows.EXPECT().Close()

	mockPool.EXPECT().
		Query(gomock.Any(), getRouteLegQuery, gomock.Any()).
		Return(mockRouteLegRows, nil)
	mockRouteLegRows.EXPECT().Next().Return(true)
	mockRouteLegRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(legScanError)
	mockRouteLegRows.EXPECT().Close()

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)

	require.ErrorIs(t, err, legScanError)
	assert.Nil(t, routes)
}
