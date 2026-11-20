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
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestRoutesRepository_GetRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	routeID := 1
	expectedRoute := &domain.Route{
		ID:          1,
		TripID:      tripID,
		TripVersion: 1,
		Version:     1,
		Status:      domain.RouteStatusGenerated,
		Days:        nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock route scan
	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoute.ID
			*destinations[1].(*int) = expectedRoute.TripID
			*destinations[2].(*int) = expectedRoute.TripVersion
			*destinations[3].(*int) = expectedRoute.Version
			*destinations[4].(*domain.RouteStatus) = expectedRoute.Status
			*destinations[5].(*time.Time) = expectedRoute.CreatedAt
			*destinations[6].(*time.Time) = expectedRoute.UpdatedAt
			return nil
		})

	// Mock route day rows (empty)
	mockRouteDayRows.EXPECT().Next().Return(false)

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			getRouteQuery,
			tripID,
			routeID,
		).
		Return(mockRow)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRouteDayQuery,
			expectedRoute.ID,
		).
		Return(mockRouteDayRows, nil)

	mockTx.EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)
	require.NoError(t, err)
	assert.Equal(t, expectedRoute, route)
}

func TestRoutesRepository_GetRoute_BeginTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	routeID := 1
	beginTxError := errors.New("begin tx error")

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(nil, beginTxError)

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)
	require.Error(t, err)
	assert.Nil(t, route)
	assert.ErrorIs(t, err, beginTxError)
}

func TestRoutesRepository_GetRoute_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	routeID := 1
	scanError := errors.New("scan error")

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(scanError)

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			getRouteQuery,
			tripID,
			routeID,
		).
		Return(mockRow)

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

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)
	require.Error(t, err)
	assert.Nil(t, route)
	assert.ErrorIs(t, err, scanError)
}

func TestRoutesRepository_GetRoute_RouteDayScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	routeID := 1
	scanError := errors.New("route day scan error")

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(nil)

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			getRouteQuery,
			tripID,
			routeID,
		).
		Return(mockRow)

	mockTx.EXPECT().
		Query(
			gomock.Any(),
			getRouteDayQuery,
			gomock.Any(),
		).
		Return(mockRouteDayRows, nil)

	mockRouteDayRows.EXPECT().Next().Return(true)
	mockRouteDayRows.EXPECT().
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

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)
	require.Error(t, err)
	assert.Nil(t, route)
	assert.ErrorIs(t, err, scanError)
}

func TestRoutesRepository_GetRoute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	routeID := 1
	noRowsError := core_postgres_pool_errors.ErrNoRows

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(noRowsError)

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	mockTx.EXPECT().
		QueryRow(
			gomock.Any(),
			getRouteQuery,
			tripID,
			routeID,
		).
		Return(mockRow)

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

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)
	require.Error(t, err)
	assert.Nil(t, route)
	assert.ErrorIs(t, err, errors.Unwrap(errors.Unwrap(err)), errors.Unwrap(errors.Unwrap(errors.Unwrap(errors.Unwrap(err)))))
}

func TestRoutesRepository_GetRoute_RouteDayQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID, routeID := 1, 1
	dayQueryError := errors.New("route day query error")

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(nil)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().QueryRow(gomock.Any(), getRouteQuery, tripID, routeID).Return(mockRow)
	mockTx.EXPECT().Query(gomock.Any(), getRouteDayQuery, gomock.Any()).Return(nil, dayQueryError)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)

	require.ErrorIs(t, err, dayQueryError)
	assert.Nil(t, route)
}

func TestRoutesRepository_GetRoute_RouteStopQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID, routeID := 1, 1
	stopQueryError := errors.New("route stop query error")

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(nil)
	mockRouteDayRows.EXPECT().Next().Return(true)
	mockRouteDayRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRouteDayRows.EXPECT().Next().Return(false)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().QueryRow(gomock.Any(), getRouteQuery, tripID, routeID).Return(mockRow)
	mockTx.EXPECT().Query(gomock.Any(), getRouteDayQuery, gomock.Any()).Return(mockRouteDayRows, nil)
	mockTx.EXPECT().Query(gomock.Any(), getRouteStopQuery, gomock.Any()).Return(nil, stopQueryError)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)

	require.ErrorIs(t, err, stopQueryError)
	assert.Nil(t, route)
}

func TestRoutesRepository_GetRoute_RouteStopScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := mocks.NewMockRows(ctrl)
	mockRouteStopRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID, routeID := 1, 1
	stopScanError := errors.New("route stop scan error")

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(nil)
	mockRouteDayRows.EXPECT().Next().Return(true)
	mockRouteDayRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRouteDayRows.EXPECT().Next().Return(false)
	mockRouteStopRows.EXPECT().Next().Return(true)
	mockRouteStopRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(stopScanError)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().QueryRow(gomock.Any(), getRouteQuery, tripID, routeID).Return(mockRow)
	mockTx.EXPECT().Query(gomock.Any(), getRouteDayQuery, gomock.Any()).Return(mockRouteDayRows, nil)
	mockTx.EXPECT().Query(gomock.Any(), getRouteStopQuery, gomock.Any()).Return(mockRouteStopRows, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)

	require.ErrorIs(t, err, stopScanError)
	assert.Nil(t, route)
}

func TestRoutesRepository_GetRoute_RouteLegQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := mocks.NewMockRows(ctrl)
	mockRouteStopRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID, routeID := 1, 1
	legQueryError := errors.New("route leg query error")

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(nil)
	mockRouteDayRows.EXPECT().Next().Return(true)
	mockRouteDayRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRouteDayRows.EXPECT().Next().Return(false)
	mockRouteStopRows.EXPECT().Next().Return(false)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().QueryRow(gomock.Any(), getRouteQuery, tripID, routeID).Return(mockRow)
	mockTx.EXPECT().Query(gomock.Any(), getRouteDayQuery, gomock.Any()).Return(mockRouteDayRows, nil)
	mockTx.EXPECT().Query(gomock.Any(), getRouteStopQuery, gomock.Any()).Return(mockRouteStopRows, nil)
	mockTx.EXPECT().Query(gomock.Any(), getRouteLegQuery, gomock.Any()).Return(nil, legQueryError)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)

	require.ErrorIs(t, err, legQueryError)
	assert.Nil(t, route)
}

func TestRoutesRepository_GetRoute_RouteLegScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := mocks.NewMockRows(ctrl)
	mockRouteStopRows := mocks.NewMockRows(ctrl)
	mockRouteLegRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID, routeID := 1, 1
	legScanError := errors.New("route leg scan error")

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(nil)
	mockRouteDayRows.EXPECT().Next().Return(true)
	mockRouteDayRows.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRouteDayRows.EXPECT().Next().Return(false)
	mockRouteStopRows.EXPECT().Next().Return(false)
	mockRouteLegRows.EXPECT().Next().Return(true)
	mockRouteLegRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(legScanError)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().QueryRow(gomock.Any(), getRouteQuery, tripID, routeID).Return(mockRow)
	mockTx.EXPECT().Query(gomock.Any(), getRouteDayQuery, gomock.Any()).Return(mockRouteDayRows, nil)
	mockTx.EXPECT().Query(gomock.Any(), getRouteStopQuery, gomock.Any()).Return(mockRouteStopRows, nil)
	mockTx.EXPECT().Query(gomock.Any(), getRouteLegQuery, gomock.Any()).Return(mockRouteLegRows, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)

	require.ErrorIs(t, err, legScanError)
	assert.Nil(t, route)
}

func TestRoutesRepository_GetRoute_Success_WithDayStopAndLeg(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockRouteDayRows := mocks.NewMockRows(ctrl)
	mockRouteStopRows := mocks.NewMockRows(ctrl)
	mockRouteLegRows := mocks.NewMockRows(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	expectedRoute := &domain.Route{
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
	}

	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoute.ID
			*destinations[1].(*int) = expectedRoute.TripID
			*destinations[2].(*int) = expectedRoute.TripVersion
			*destinations[3].(*int) = expectedRoute.Version
			*destinations[4].(*domain.RouteStatus) = expectedRoute.Status
			*destinations[5].(*time.Time) = expectedRoute.CreatedAt
			*destinations[6].(*time.Time) = expectedRoute.UpdatedAt
			return nil
		})

	mockRouteDayRows.EXPECT().Next().Return(true)
	mockRouteDayRows.EXPECT().
		Scan(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoute.Days[0].ID
			*destinations[1].(*time.Time) = expectedRoute.Days[0].Date
			*destinations[2].(*int) = expectedRoute.Days[0].DayNumber
			return nil
		})
	mockRouteDayRows.EXPECT().Next().Return(false)

	mockRouteStopRows.EXPECT().Next().Return(true)
	mockRouteStopRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoute.Days[0].Stops[0].ID
			*destinations[1].(*int) = expectedRoute.Days[0].Stops[0].PlaceID
			*destinations[2].(*int) = expectedRoute.Days[0].Stops[0].Position
			*destinations[3].(*string) = expectedRoute.Days[0].Stops[0].Name
			*destinations[4].(*string) = "Test Description 1"
			*destinations[5].(*float64) = expectedRoute.Days[0].Stops[0].Coordinates.Latitude
			*destinations[6].(*float64) = expectedRoute.Days[0].Stops[0].Coordinates.Longitude
			*destinations[7].(*time.Time) = expectedRoute.Days[0].Stops[0].ArrivalTime
			*destinations[8].(*time.Time) = expectedRoute.Days[0].Stops[0].DepartureTime
			*destinations[9].(*int64) = int64(3600)
			return nil
		})
	mockRouteStopRows.EXPECT().Next().Return(false)

	mockRouteLegRows.EXPECT().Next().Return(true)
	mockRouteLegRows.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedRoute.Days[0].Legs[0].ID
			*destinations[1].(*int) = expectedRoute.Days[0].Legs[0].FromStopID
			*destinations[2].(*int) = expectedRoute.Days[0].Legs[0].ToStopID
			*destinations[3].(*domain.TravelMode) = expectedRoute.Days[0].Legs[0].TravelMode
			*destinations[4].(*int) = expectedRoute.Days[0].Legs[0].Distance
			*destinations[5].(*int64) = int64(60)
			return nil
		})
	mockRouteLegRows.EXPECT().Next().Return(false)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().QueryRow(gomock.Any(), getRouteQuery, gomock.Any(), gomock.Any()).Return(mockRow)
	mockTx.EXPECT().Query(gomock.Any(), getRouteDayQuery, gomock.Any()).Return(mockRouteDayRows, nil)
	mockTx.EXPECT().Query(gomock.Any(), getRouteStopQuery, gomock.Any()).Return(mockRouteStopRows, nil)
	mockTx.EXPECT().Query(gomock.Any(), getRouteLegQuery, gomock.Any()).Return(mockRouteLegRows, nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	route, err := routesRepository.GetRoute(context.Background(), 1, 1)

	require.NoError(t, err)
	require.NotNil(t, route)
	require.Len(t, route.Days, 1)
	require.Len(t, route.Days[0].Stops, 1)
	require.Len(t, route.Days[0].Legs, 1)
}
