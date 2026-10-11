package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testhelpers"
)

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	expectedUser := testhelpers.UserTest()

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			getUserQuery,
			userID,
		).
		Return(mockRow)

	mockRow.
		EXPECT().
		Scan(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(scanDestination ...any) error {
			*scanDestination[0].(*int) = expectedUser.ID
			*scanDestination[1].(*int) = expectedUser.Version
			*scanDestination[2].(*string) = expectedUser.Name
			*scanDestination[3].(*string) = expectedUser.Surname
			*scanDestination[4].(*string) = expectedUser.Username
			*scanDestination[5].(*time.Time) = expectedUser.BirthDate
			*scanDestination[6].(*string) = expectedUser.Description
			*scanDestination[7].(*string) = expectedUser.Email
			*scanDestination[8].(*string) = expectedUser.PhoneNumber
			*scanDestination[9].(*time.Time) = expectedUser.CreatedAt
			*scanDestination[10].(*time.Time) = expectedUser.UpdatedAt
			*scanDestination[11].(**time.Time) = expectedUser.DeletedAt

			return nil
		})

	actualUser, err := usersRepository.GetUser(context.Background(), userID)
	require.NoError(t, err)
	assert.NotNil(t, actualUser)
	assert.Equal(t, expectedUser, *actualUser)
}

func TestGetUser_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	scanError := errors.New("query row error")

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			getUserQuery,
			userID,
		).
		Return(mockRow)

	mockRow.
		EXPECT().
		Scan(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(scanError)

	actualUser, err := usersRepository.GetUser(context.Background(), userID)
	require.ErrorIs(t, err, scanError)
	assert.Nil(t, actualUser)
}

func TestGetUser_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			getUserQuery,
			userID,
		).
		Return(mockRow)

	mockRow.
		EXPECT().
		Scan(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		Return(pgx.ErrNoRows)

	actualUser, err := usersRepository.GetUser(context.Background(), userID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	assert.Nil(t, actualUser)
	assert.EqualError(
		t,
		err,
		fmt.Errorf("get userID %d: %w", userID, pgx.ErrNoRows).Error(),
	)
}
