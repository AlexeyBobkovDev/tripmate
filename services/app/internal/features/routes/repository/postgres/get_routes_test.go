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
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	createdAt := time.Now()

	expectedRoutes := []*domain.Route{
		{
			ID:          1,
			TripID:      1,
			TripVersion: 1,
			Version:     1,
			Status:      domain.RouteStatusGenerated,
			Days:        nil,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		},
	}

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRoutesQuery,
			tripID,
		).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoutes[0].ID
			*destinations[1].(*int) = expectedRoutes[0].TripID
			*destinations[2].(*int) = expectedRoutes[0].TripVersion
			*destinations[3].(*int) = expectedRoutes[0].Version
			*destinations[4].(*domain.RouteStatus) = expectedRoutes[0].Status
			*destinations[5].(*time.Time) = expectedRoutes[0].CreatedAt
			*destinations[6].(*time.Time) = expectedRoutes[0].UpdatedAt
			return nil
		}),
		mockRows.EXPECT().Next().Return(false),
	)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRouteDaysQuery,
			expectedRoutes[0].ID,
		).
		Return(mockRouteDaysRows, nil)

	mockRouteDaysRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)
	require.NoError(t, err)
	require.Len(t, routes, 1)
	assert.Equal(t, expectedRoutes, routes)
}

func TestRoutesRepository_GetRoutes_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	queryError := errors.New("query error")

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRoutesQuery,
			tripID,
		).
		Return(nil, queryError)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)
	require.ErrorIs(t, err, queryError)
	assert.Nil(t, routes)
}

func TestRoutesRepository_GetRoutes_BeginTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	beginTxError := errors.New("begin tx error")

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(nil, beginTxError)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)
	require.Error(t, err)
	assert.Nil(t, routes)
	assert.ErrorIs(t, err, beginTxError)
}

func TestRoutesRepository_GetRoutes_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	scanError := errors.New("scan error")

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRoutesQuery,
			tripID,
		).
		Return(mockRows, nil)

	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(scanError)

	mockTx.EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)
	require.Error(t, err)
	assert.Nil(t, routes)
	assert.ErrorIs(t, err, scanError)
}

func TestRoutesRepository_GetRoutes_RouteDayScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	createdAt := time.Now()
	scanError := errors.New("route day scan error")

	expectedRoutes := []*domain.Route{
		{
			ID:          1,
			TripID:      1,
			TripVersion: 1,
			Version:     1,
			Status:      domain.RouteStatusGenerated,
			Days:        nil,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		},
	}

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRoutesQuery,
			tripID,
		).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoutes[0].ID
			*destinations[1].(*int) = expectedRoutes[0].TripID
			*destinations[2].(*int) = expectedRoutes[0].TripVersion
			*destinations[3].(*int) = expectedRoutes[0].Version
			*destinations[4].(*domain.RouteStatus) = expectedRoutes[0].Status
			*destinations[5].(*time.Time) = expectedRoutes[0].CreatedAt
			*destinations[6].(*time.Time) = expectedRoutes[0].UpdatedAt
			return nil
		}),
		mockRows.EXPECT().Next().Return(false),
	)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRouteDaysQuery,
			expectedRoutes[0].ID,
		).
		Return(mockRouteDaysRows, nil)

	mockRouteDaysRows.EXPECT().Next().Return(true)
	mockRouteDaysRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(scanError)

	mockTx.EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)
	require.Error(t, err)
	assert.Nil(t, routes)
	assert.ErrorIs(t, err, scanError)
}

func TestRoutesRepository_GetRoutes_RouteStopScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	mockRouteStopsRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	createdAt := time.Now()
	scanError := errors.New("route stop scan error")

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
					Date:      time.Date(2035, 3, 10, 0, 0, 0, 0, time.UTC),
					DayNumber: 1,
					Stops:     []domain.RouteStop{},
					Legs:      []domain.RouteLeg{},
				},
			},
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
	}

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRoutesQuery,
			tripID,
		).
		Return(mockRows, nil)

	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoutes[0].ID
			*destinations[1].(*int) = expectedRoutes[0].TripID
			*destinations[2].(*int) = expectedRoutes[0].TripVersion
			*destinations[3].(*int) = expectedRoutes[0].Version
			*destinations[4].(*domain.RouteStatus) = expectedRoutes[0].Status
			*destinations[5].(*time.Time) = expectedRoutes[0].CreatedAt
			*destinations[6].(*time.Time) = expectedRoutes[0].UpdatedAt
			return nil
		}),
		mockRows.EXPECT().Next().Return(false),
	)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRouteDaysQuery,
			expectedRoutes[0].ID,
		).
		Return(mockRouteDaysRows, nil)

	mockRouteDaysRows.EXPECT().Next().Return(true)
	mockRouteDaysRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoutes[0].Days[0].ID
			*destinations[1].(*time.Time) = expectedRoutes[0].Days[0].Date
			*destinations[2].(*int) = expectedRoutes[0].Days[0].DayNumber
			return nil
		})
	mockRouteDaysRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRouteStopsQuery,
			expectedRoutes[0].Days[0].ID,
		).
		Return(mockRouteStopsRows, nil)

	mockRouteStopsRows.EXPECT().Next().Return(true)
	mockRouteStopsRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(scanError)

	mockTx.EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)
	require.Error(t, err)
	assert.Nil(t, routes)
	assert.ErrorIs(t, err, scanError)
}

func TestRoutesRepository_GetRoutes_RouteDayQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	dayQueryError := errors.New("route day query error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Query(gomock.Any(), getRoutesQuery, tripID).Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().Scan(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil)
	mockRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().Query(gomock.Any(), getRouteDaysQuery, gomock.Any()).Return(nil, dayQueryError)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)

	require.ErrorIs(t, err, dayQueryError)
	assert.Nil(t, routes)
}

func TestRoutesRepository_GetRoutes_RouteStopQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	stopQueryError := errors.New("route stop query error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Query(gomock.Any(), getRoutesQuery, tripID).Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().Scan(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil)
	mockRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().Query(gomock.Any(), getRouteDaysQuery, gomock.Any()).Return(mockRouteDaysRows, nil)
	mockRouteDaysRows.EXPECT().Next().Return(true)
	mockRouteDaysRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRouteDaysRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().Query(gomock.Any(), getRouteStopsQuery, gomock.Any()).Return(nil, stopQueryError)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)

	require.ErrorIs(t, err, stopQueryError)
	assert.Nil(t, routes)
}

func TestRoutesRepository_GetRoutes_RouteLegQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	mockRouteStopsRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	legQueryError := errors.New("route leg query error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Query(gomock.Any(), getRoutesQuery, tripID).Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().Scan(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil)
	mockRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().Query(gomock.Any(), getRouteDaysQuery, gomock.Any()).Return(mockRouteDaysRows, nil)
	mockRouteDaysRows.EXPECT().Next().Return(true)
	mockRouteDaysRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRouteDaysRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().Query(gomock.Any(), getRouteStopsQuery, gomock.Any()).Return(mockRouteStopsRows, nil)
	mockRouteStopsRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().Query(gomock.Any(), getRouteLegsQuery, gomock.Any()).Return(nil, legQueryError)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)

	require.ErrorIs(t, err, legQueryError)
	assert.Nil(t, routes)
}

func TestRoutesRepository_GetRoutes_RouteLegScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	mockRouteStopsRows := mocks.NewMockRows(ctrl)
	mockRouteLegRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	legScanError := errors.New("route leg scan error")

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Query(gomock.Any(), getRoutesQuery, tripID).Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(true)
	mockRows.EXPECT().Scan(
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil)
	mockRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().Query(gomock.Any(), getRouteDaysQuery, gomock.Any()).Return(mockRouteDaysRows, nil)
	mockRouteDaysRows.EXPECT().Next().Return(true)
	mockRouteDaysRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRouteDaysRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().Query(gomock.Any(), getRouteStopsQuery, gomock.Any()).Return(mockRouteStopsRows, nil)
	mockRouteStopsRows.EXPECT().Next().Return(false)

	mockTx.EXPECT().Query(gomock.Any(), getRouteLegsQuery, gomock.Any()).Return(mockRouteLegRows, nil)
	mockRouteLegRows.EXPECT().Next().Return(true)
	mockRouteLegRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(legScanError)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)

	require.ErrorIs(t, err, legScanError)
	assert.Nil(t, routes)
}

func TestRoutesRepository_GetRoutes_Success_WithDayStopAndLeg(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockRouteDaysRows := mocks.NewMockRows(ctrl)
	mockRouteStopsRows := mocks.NewMockRows(ctrl)
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
							ID:            1,
							PlaceID:       1,
							Position:      1,
							Name:          "Test Stop 1",
							Coordinates:   domain.Coordinates{Latitude: 55.7539, Longitude: 37.6208},
							ArrivalTime:   time.Date(2026, time.November, 14, 9, 0, 0, 0, time.UTC),
							DepartureTime: time.Date(2026, time.November, 14, 10, 0, 0, 0, time.UTC),
							VisitDuration: time.Hour,
						},
					},
					Legs: []domain.RouteLeg{
						{
							ID:         1,
							FromStopID: 1,
							ToStopID:   2,
							Distance:   1000,
						},
					},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Query(gomock.Any(), getRoutesQuery, 1).Return(mockRows, nil)
	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoutes[0].ID
			*destinations[1].(*int) = expectedRoutes[0].TripID
			*destinations[2].(*int) = expectedRoutes[0].TripVersion
			*destinations[3].(*int) = expectedRoutes[0].Version
			*destinations[4].(*domain.RouteStatus) = expectedRoutes[0].Status
			*destinations[5].(*time.Time) = expectedRoutes[0].CreatedAt
			*destinations[6].(*time.Time) = expectedRoutes[0].UpdatedAt
			return nil
		}),
		mockRows.EXPECT().Next().Return(false),
	)

	mockTx.EXPECT().Query(gomock.Any(), getRouteDaysQuery, gomock.Any()).Return(mockRouteDaysRows, nil)
	gomock.InOrder(
		mockRouteDaysRows.EXPECT().Next().Return(true),
		mockRouteDaysRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(destinations ...any) error {
				*destinations[0].(*int) = expectedRoutes[0].Days[0].ID
				*destinations[1].(*time.Time) = expectedRoutes[0].Days[0].Date
				*destinations[2].(*int) = expectedRoutes[0].Days[0].DayNumber
				return nil
			},
		),
		mockRouteDaysRows.EXPECT().Next().Return(false),
	)

	mockTx.EXPECT().Query(gomock.Any(), getRouteStopsQuery, gomock.Any()).Return(mockRouteStopsRows, nil)
	gomock.InOrder(
		mockRouteStopsRows.EXPECT().Next().Return(true),
		mockRouteStopsRows.EXPECT().Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoutes[0].Days[0].Stops[0].ID
			*destinations[1].(*int) = expectedRoutes[0].Days[0].Stops[0].PlaceID
			*destinations[2].(*int) = expectedRoutes[0].Days[0].Stops[0].Position
			*destinations[3].(*string) = expectedRoutes[0].Days[0].Stops[0].Name
			*destinations[4].(*string) = "Test Description 1"
			*destinations[5].(*float64) = expectedRoutes[0].Days[0].Stops[0].Coordinates.Latitude
			*destinations[6].(*float64) = expectedRoutes[0].Days[0].Stops[0].Coordinates.Longitude
			*destinations[7].(*time.Time) = expectedRoutes[0].Days[0].Stops[0].ArrivalTime
			*destinations[8].(*time.Time) = expectedRoutes[0].Days[0].Stops[0].DepartureTime
			*destinations[9].(*int64) = int64(3600)
			return nil
		}),
		mockRouteStopsRows.EXPECT().Next().Return(false),
	)

	mockTx.EXPECT().Query(gomock.Any(), getRouteLegsQuery, gomock.Any()).Return(mockRouteLegRows, nil)
	gomock.InOrder(
		mockRouteLegRows.EXPECT().Next().Return(true),
		mockRouteLegRows.EXPECT().Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoutes[0].Days[0].Legs[0].ID
			*destinations[1].(*int) = expectedRoutes[0].Days[0].Legs[0].FromStopID
			*destinations[2].(*int) = expectedRoutes[0].Days[0].Legs[0].ToStopID
			*destinations[3].(*domain.TravelMode) = expectedRoutes[0].Days[0].Legs[0].TravelMode
			*destinations[4].(*int) = expectedRoutes[0].Days[0].Legs[0].Distance
			*destinations[5].(*int64) = int64(60)
			return nil
		}),
		mockRouteLegRows.EXPECT().Next().Return(false),
	)

	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	routes, err := routesRepository.GetRoutes(context.Background(), tripID)

	require.NoError(t, err)
	require.NotNil(t, routes)
	require.Len(t, routes, 1)
	require.Len(t, routes[0].Days, 1)
	require.Len(t, routes[0].Days[0].Stops, 1)
	require.Len(t, routes[0].Days[0].Legs, 1)
}
