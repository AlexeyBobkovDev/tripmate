package trips_postgres_repository

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

func TestTripsRepository_DeleteTrip_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1

	mockPool.
		EXPECT().
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
			deleteTripQuery,
			ownerID,
			tripID,
		).
		Return(mockCommandTag, nil)
	mockTx.
		EXPECT().
		Rollback(gomock.Any()).
		Return(nil)
	mockTx.
		EXPECT().
		Commit(gomock.Any()).
		Return(nil)

	mockPool.
		EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	err := tripsRepository.DeleteTrip(context.Background(), ownerID, tripID)
	require.NoError(t, err)
}

func TestTripsRepository_DeleteTrip_BeginTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
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

	err := tripsRepository.DeleteTrip(context.Background(), ownerID, tripID)
	require.ErrorIs(t, err, beginTxError)
}

func TestTripsRepository_DeleteTrip_ExecError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
	txExecError := errors.New("Tx exec error")

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockTx.
		EXPECT().
		Exec(
			gomock.Any(),
			deleteTripQuery,
			ownerID,
			tripID,
		).
		Return(mockCommandTag, txExecError)
	mockTx.
		EXPECT().
		Rollback(gomock.Any()).
		Return(nil)

	mockPool.
		EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	err := tripsRepository.DeleteTrip(context.Background(), ownerID, tripID)
	require.ErrorIs(t, err, txExecError)
}

func TestTripsRepository_DeleteTrip_CommitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1
	txCommitError := errors.New("Tx commit error")

	mockPool.
		EXPECT().
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
			deleteTripQuery,
			ownerID,
			tripID,
		).
		Return(mockCommandTag, nil)
	mockTx.
		EXPECT().
		Rollback(gomock.Any()).
		Return(nil)
	mockTx.
		EXPECT().
		Commit(gomock.Any()).
		Return(txCommitError)

	mockPool.
		EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	err := tripsRepository.DeleteTrip(context.Background(), ownerID, tripID)
	require.ErrorIs(t, err, txCommitError)
}

func TestTripsRepository_DeleteTrip_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	tripsRepository := NewTripsRepository(mockPool)

	ownerID := 1
	tripID := 1

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
			deleteTripQuery,
			ownerID,
			tripID,
		).
		Return(mockCommandTag, nil)
	mockTx.
		EXPECT().
		Rollback(gomock.Any()).
		Return(nil)

	mockPool.
		EXPECT().
		Begin(
			gomock.Any(),
		).
		Return(mockTx, nil)

	err := tripsRepository.DeleteTrip(context.Background(), ownerID, tripID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
}
