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
)

func TestDeleteUser_InvalidID(t *testing.T) {
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
				DeleteUser(
					gomock.Any(),
					gomock.Any(),
				).
				Times(0)

			err := usersService.DeleteUser(context.Background(), tt.id)

			require.ErrorAs(t, err, &core_errors.ErrBadRequest)
		},
		)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	userID := 1

	mockUsersRepository.
		EXPECT().
		DeleteUser(context.Background(), userID).
		Return(nil)

	err := usersService.DeleteUser(context.Background(), userID)
	require.NoError(t, err)
}

func TestDeleteUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockPasswordManager := mocks.NewMockPasswordManager(ctrl)
	mockUsersRepository := mocks.NewMockUsersRepository(ctrl)
	usersService := NewUsersService(mockUsersRepository, mockPasswordManager)

	userID := 1

	mockUsersRepository.
		EXPECT().
		DeleteUser(context.Background(), userID).
		Return(core_postgres_pool_errors.ErrViolatedForeignKey)

	err := usersService.DeleteUser(context.Background(), userID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, core_postgres_pool_errors.ErrViolatedForeignKey)
}
