package users_service

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testhelpers"
)

func TestGetUser_InvalidID(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{
			name: "negative ID",
			id:   -1,
		},
		{
			name: "zero ID",
			id:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
			mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
			usersService := NewUsersService(
				mockUsersRepository,
				mockPasswordManager,
			)

			mockUsersRepository.EXPECT().
				GetUser(gomock.Any(), gomock.Any()).
				Times(0)

			retrievedUser, err := usersService.GetUser(
				context.Background(),
				tt.id,
			)

			require.ErrorIs(t, err, core_errors.ErrBadRequest)
			assert.Nil(t, retrievedUser)
		})
	}
}

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(
		mockUsersRepository,
		mockPasswordManager,
	)

	userID := 1
	expectedUser := testhelpers.UserTest()

	mockUsersRepository.EXPECT().
		GetUser(context.Background(), userID).
		Return(&expectedUser, nil)

	retrievedUser, err := usersService.GetUser(
		context.Background(),
		userID,
	)

	require.NoError(t, err)
	require.NotNil(t, retrievedUser)
	assert.Equal(t, expectedUser, *retrievedUser)
}

func TestGetUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(
		mockUsersRepository,
		mockPasswordManager,
	)

	userID := 1

	mockUsersRepository.EXPECT().
		GetUser(context.Background(), userID).
		Return(nil, errors.Join(pgx.ErrNoRows, core_errors.ErrNotFound))

	retrievedUser, err := usersService.GetUser(
		context.Background(),
		userID,
	)

	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.ErrorIs(t, err, core_errors.ErrNotFound)
	require.Nil(t, retrievedUser)
}
