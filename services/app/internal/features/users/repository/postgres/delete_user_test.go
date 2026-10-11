package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1

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
			deleteUserQuery,
			userID,
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

	err := usersRepository.DeleteUser(context.Background(), userID)
	require.NoError(t, err)
}

func TestDeleteUser_BeginTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
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

	err := usersRepository.DeleteUser(context.Background(), userID)
	require.ErrorIs(t, err, beginTxError)
}

func TestDeleteUser_ExecError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	txExecError := errors.New("Tx exec error")

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockTx.
		EXPECT().
		Exec(
			gomock.Any(),
			deleteUserQuery,
			userID,
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

	err := usersRepository.DeleteUser(context.Background(), userID)
	require.ErrorIs(t, err, txExecError)
}

func TestDeleteUser_CommitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
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
			deleteUserQuery,
			userID,
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

	err := usersRepository.DeleteUser(context.Background(), userID)
	require.ErrorIs(t, err, txCommitError)
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCommandTag := mocks.NewMockCommandTag(ctrl)
	mockTx := mocks.NewMockTx(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1

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
			deleteUserQuery,
			userID,
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

	err := usersRepository.DeleteUser(context.Background(), userID)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	assert.EqualError(
		t,
		err,
		fmt.Errorf("the user %d does not exist: %w", userID, core_errors.ErrNotFound).Error(),
	)
}
