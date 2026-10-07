package users_postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testhelpers"
)

func TestCreateUserRepository_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	createUserRequest := testhelpers.UserCreateTest()
	expectedCreatedUser := testhelpers.UserTest()

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
		*scanDestination[0].(*int) = expectedCreatedUser.ID
		*scanDestination[1].(*int) = expectedCreatedUser.Version
		*scanDestination[2].(*string) = expectedCreatedUser.Name
		*scanDestination[3].(*string) = expectedCreatedUser.Surname
		*scanDestination[4].(*string) = expectedCreatedUser.Username
		*scanDestination[5].(*time.Time) = expectedCreatedUser.BirthDate
		*scanDestination[6].(*string) = expectedCreatedUser.Description
		*scanDestination[7].(*string) = expectedCreatedUser.Email
		*scanDestination[8].(*string) = expectedCreatedUser.PhoneNumber
		*scanDestination[9].(*time.Time) = expectedCreatedUser.CreatedAt
		*scanDestination[10].(*time.Time) = expectedCreatedUser.UpdatedAt
		*scanDestination[11].(**time.Time) = expectedCreatedUser.DeletedAt

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
			query,
			createUserRequest.Name,
			createUserRequest.Surname,
			createUserRequest.Username,
			createUserRequest.BirthDate,
			createUserRequest.Description,
			createUserRequest.Email,
			createUserRequest.PhoneNumber,
			createUserRequest.PasswordHash,
		).
		Return(mockRow)

	actualUser, err := usersRepository.CreateUser(context.Background(), &createUserRequest)
	require.NoError(t, err)

	assert.Equal(t, expectedCreatedUser, *actualUser)
}

func TestCreateUserRepository_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	usersRepository := NewUsersRepository(mockPool)

	createUserRequest := testhelpers.UserCreateTest()

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
		Return(core_postgres_pool_errors.ErrNoRows)

	mockPool.
		EXPECT().
		OpTimeout().
		Return(5 * time.Second)

	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			query,
			createUserRequest.Name,
			createUserRequest.Surname,
			createUserRequest.Username,
			createUserRequest.BirthDate,
			createUserRequest.Description,
			createUserRequest.Email,
			createUserRequest.PhoneNumber,
			createUserRequest.PasswordHash,
		).
		Return(mockRow)

	actualUser, err := usersRepository.CreateUser(context.Background(), &createUserRequest)
	require.ErrorIs(t, err, core_postgres_pool_errors.ErrNoRows)

	assert.Nil(t, actualUser)
}
