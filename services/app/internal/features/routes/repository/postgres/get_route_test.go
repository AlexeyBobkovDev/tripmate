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
