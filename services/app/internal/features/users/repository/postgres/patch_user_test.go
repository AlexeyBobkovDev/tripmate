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

func TestPatchUserRepository_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	patchUserRequest := testutils.UserTest()

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
		).DoAndReturn(func(scanDestination ...any) error {
		*scanDestination[0].(*int) = userID
		*scanDestination[1].(*int) = patchUserRequest.Version + 1
		*scanDestination[2].(*string) = patchUserRequest.Name
		*scanDestination[3].(*string) = patchUserRequest.Surname
		*scanDestination[4].(*string) = patchUserRequest.Username
		*scanDestination[5].(*time.Time) = patchUserRequest.BirthDate
		*scanDestination[6].(*string) = patchUserRequest.Description
		*scanDestination[7].(*string) = patchUserRequest.Email
		*scanDestination[8].(*string) = patchUserRequest.PhoneNumber
		*scanDestination[9].(*time.Time) = patchUserRequest.CreatedAt
		*scanDestination[10].(*time.Time) = patchUserRequest.UpdatedAt
		*scanDestination[11].(**time.Time) = patchUserRequest.DeletedAt

		return nil
	})

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			patchUserQuery,
			userID,
			patchUserRequest.Version,
			patchUserRequest.Name,
			patchUserRequest.Surname,
			patchUserRequest.Username,
			patchUserRequest.BirthDate,
			patchUserRequest.Description,
			patchUserRequest.Email,
			patchUserRequest.PhoneNumber,
		).
		Return(mockRow)

	actualUser, err := usersRepository.PatchUser(context.Background(), userID, &patchUserRequest)
	require.NoError(t, err)

	assert.Equal(t, userID, actualUser.ID)
	assert.Equal(t, patchUserRequest.Version+1, actualUser.Version)
	assert.Equal(t, patchUserRequest.Name, actualUser.Name)
	assert.Equal(t, patchUserRequest.Surname, actualUser.Surname)
	assert.Equal(t, patchUserRequest.Username, actualUser.Username)
	assert.Equal(t, patchUserRequest.BirthDate, actualUser.BirthDate)
	assert.Equal(t, patchUserRequest.Description, actualUser.Description)
	assert.Equal(t, patchUserRequest.Email, actualUser.Email)
	assert.Equal(t, patchUserRequest.PhoneNumber, actualUser.PhoneNumber)
	assert.Equal(t, patchUserRequest.CreatedAt, actualUser.CreatedAt)
	assert.Equal(t, patchUserRequest.UpdatedAt, actualUser.UpdatedAt)
	assert.Equal(t, patchUserRequest.DeletedAt, actualUser.DeletedAt)
}

func TestPatchUserRepository_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	scanError := errors.New("scan error")
	patchUserRequest := testutils.UserTest()

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
		).Return(scanError)

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			patchUserQuery,
			userID,
			patchUserRequest.Version,
			patchUserRequest.Name,
			patchUserRequest.Surname,
			patchUserRequest.Username,
			patchUserRequest.BirthDate,
			patchUserRequest.Description,
			patchUserRequest.Email,
			patchUserRequest.PhoneNumber,
		).
		Return(mockRow)

	actualUser, err := usersRepository.PatchUser(context.Background(), userID, &patchUserRequest)
	require.ErrorIs(t, err, scanError)
	require.Nil(t, actualUser)
}

func TestPatchUserRepository_ExistsCheckScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockUserExistsRow := mocks.NewMockRow(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	userVersion := 1
	existsCheckError := errors.New("exists check scan error")
	patchUserRequest := testutils.UserTest()
	patchUserRequest.Version = userVersion + 1

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
		).Return(pgx.ErrNoRows)

	mockUserExistsRow.EXPECT().
		Scan(
			gomock.Any(),
		).
		Return(existsCheckError)

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			patchUserQuery,
			userID,
			patchUserRequest.Version,
			patchUserRequest.Name,
			patchUserRequest.Surname,
			patchUserRequest.Username,
			patchUserRequest.BirthDate,
			patchUserRequest.Description,
			patchUserRequest.Email,
			patchUserRequest.PhoneNumber,
		).
		Return(mockRow)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			userExistsQuery,
			userID,
		).
		Return(mockUserExistsRow)

	actualUser, err := usersRepository.PatchUser(context.Background(), userID, &patchUserRequest)
	require.ErrorIs(t, err, existsCheckError)
	require.Nil(t, actualUser)
}

func TestPatchUserRepository_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockUserExistsRow := mocks.NewMockRow(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	userVersion := 1
	patchUserRequest := testutils.UserTest()
	patchUserRequest.Version = userVersion + 1

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
		).Return(pgx.ErrNoRows)

	mockUserExistsRow.EXPECT().
		Scan(
			gomock.Any(),
		).
		DoAndReturn(func(dest ...any) error {
			*dest[0].(*bool) = false
			return nil
		})

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			patchUserQuery,
			userID,
			patchUserRequest.Version,
			patchUserRequest.Name,
			patchUserRequest.Surname,
			patchUserRequest.Username,
			patchUserRequest.BirthDate,
			patchUserRequest.Description,
			patchUserRequest.Email,
			patchUserRequest.PhoneNumber,
		).
		Return(mockRow)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			userExistsQuery,
			userID,
		).
		Return(mockUserExistsRow)

	actualUser, err := usersRepository.PatchUser(context.Background(), userID, &patchUserRequest)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.False(t, errors.Is(err, core_errors.ErrConflict))
	require.Nil(t, actualUser)
}

func TestPatchUserRepository_ErrConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	mockUserExistsRow := mocks.NewMockRow(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	userID := 1
	userVersion := 1
	patchUserRequest := testutils.UserTest()
	patchUserRequest.Version = userVersion + 1

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
		).Return(pgx.ErrNoRows)

	mockUserExistsRow.EXPECT().
		Scan(
			gomock.Any(),
		).
		DoAndReturn(func(dest ...any) error {
			*dest[0].(*bool) = true
			return nil
		})

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			patchUserQuery,
			userID,
			patchUserRequest.Version,
			patchUserRequest.Name,
			patchUserRequest.Surname,
			patchUserRequest.Username,
			patchUserRequest.BirthDate,
			patchUserRequest.Description,
			patchUserRequest.Email,
			patchUserRequest.PhoneNumber,
		).
		Return(mockRow)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			userExistsQuery,
			userID,
		).
		Return(mockUserExistsRow)

	actualUser, err := usersRepository.PatchUser(context.Background(), userID, &patchUserRequest)
	require.ErrorIs(t, err, core_errors.ErrConflict)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.Nil(t, actualUser)
}
