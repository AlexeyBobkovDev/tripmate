package routes_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestRoutesRepository_DeleteRoute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	routeID := 1

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockCommandTag.
		EXPECT().
		RowsAffected().
		Return(int64(1))

	mockTx.
		EXPECT().
		Exec(
			gomock.Any(),
			deleteRouteQuery,
			tripID,
			routeID,
		).
		Return(mockCommandTag, nil)

	mockTx.
		EXPECT().
		Commit(
			gomock.Any(),
		).
		Return(nil)

	mockTx.
		EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockPool.
		EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	err := routesRepository.DeleteRoute(context.Background(), tripID, routeID)
	require.NoError(t, err)
}

func TestRoutesRepository_DeleteRoute_BeginTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	routeID := 1
	beginTxError := errors.New("begin Tx error")

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, beginTxError)

	err := routesRepository.DeleteRoute(context.Background(), tripID, routeID)
	require.ErrorIs(t, err, beginTxError)
}

func TestRoutesRepository_DeleteRoute_ExecError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	routeID := 1
	txExecError := errors.New("Tx exec error")

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockTx.
		EXPECT().
		Exec(
			gomock.Any(),
			deleteRouteQuery,
			tripID,
			routeID,
		).
		Return(mockCommandTag, txExecError)

	mockTx.
		EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockPool.
		EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	err := routesRepository.DeleteRoute(context.Background(), tripID, routeID)
	require.ErrorIs(t, err, txExecError)
}

func TestRoutesRepository_DeleteRoute_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID := 1
	routeID := 1

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockCommandTag.
		EXPECT().
		RowsAffected().
		Return(int64(0))

	mockTx.
		EXPECT().
		Exec(
			gomock.Any(),
			deleteRouteQuery,
			tripID,
			routeID,
		).
		Return(mockCommandTag, nil)

	mockTx.
		EXPECT().
		Rollback(
			gomock.Any(),
		).
		Return(nil)

	mockPool.
		EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	err := routesRepository.DeleteRoute(context.Background(), tripID, routeID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}
