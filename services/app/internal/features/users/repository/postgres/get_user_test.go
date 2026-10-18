package users_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	expectedUser := testutils.UserTest()

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(mockRow)

	mockRow.EXPECT().
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
		DoAndReturn(func(destinations ...any) error {
			*destinations[0].(*int) = expectedUser.ID
			*destinations[1].(*int) = expectedUser.Version
			*destinations[2].(*string) = expectedUser.Name
			*destinations[3].(*string) = expectedUser.Surname
			*destinations[4].(*string) = expectedUser.Username
			*destinations[5].(*time.Time) = expectedUser.BirthDate
			*destinations[6].(*string) = expectedUser.Description
			*destinations[7].(*string) = expectedUser.Email
			*destinations[8].(*string) = expectedUser.PhoneNumber
			*destinations[9].(*time.Time) = expectedUser.CreatedAt
			*destinations[10].(*time.Time) = expectedUser.UpdatedAt
			*destinations[11].(**time.Time) = expectedUser.DeletedAt

			return nil
		})

	actualUser, err := usersRepository.GetUser(context.Background(), userID)

	require.NoError(t, err)
	require.NotNil(t, actualUser)
	assert.Equal(t, expectedUser, *actualUser)
}

func TestGetUser_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRow := mocks.NewMockRow(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	scanError := errors.New("query row error")

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(mockRow)

	mockRow.EXPECT().
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

	mockPool.EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.EXPECT().
		QueryRow(gomock.Any(), getUserQuery, userID).
		Return(mockRow)

	mockRow.EXPECT().
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
		Return(errors.Join(pgx.ErrNoRows, core_errors.ErrNotFound))

	actualUser, err := usersRepository.GetUser(context.Background(), userID)

	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, actualUser)
}
