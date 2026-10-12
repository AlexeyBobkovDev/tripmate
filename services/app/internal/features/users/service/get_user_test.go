package users_service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
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
			usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

			mockUsersRepository.EXPECT().
				GetUser(
					gomock.Any(),
					gomock.Any(),
				).
				Times(0)

			retrievedUser, err := usersService.GetUser(context.Background(), tt.id)
			assert.Nil(t, retrievedUser)

			require.ErrorIs(t, err, core_errors.ErrBadRequest)
		},
		)
	}
}

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	userID := 1

	testUser := testhelpers.UserTest()

	mockUsersRepository.
		EXPECT().
		GetUser(context.Background(), userID).
		Return(&testUser, nil)

	retrievedUser, err := usersService.GetUser(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, testUser, *retrievedUser)
}

func TestGetUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	userID := 1

	mockUsersRepository.
		EXPECT().
		GetUser(context.Background(), userID).
		Return(nil, core_postgres_pool_errors.ErrViolatedForeignKey)

	retrievedUser, err := usersService.GetUser(context.Background(), userID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, core_postgres_pool_errors.ErrViolatedForeignKey)
	assert.Nil(t, retrievedUser)
}
