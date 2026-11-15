package users_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	users_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service/mocks"
)

func TestUsersService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUsersRepository(ctrl)
	mockHasher := mocks.NewMockPasswordManager(ctrl)

	svc := users_service.NewUsersService(mockRepo, mockHasher)

	ctx := context.Background()
	userCreate := &domain.UserCreate{
		Name:        "Test",
		Surname:     "User",
		Username:    "testuser",
		BirthDate:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Description: "Some kind of description",
		Email:       "test@example.com",
		PhoneNumber: "+79990746988",
	}
	rawPassword := "SecurePass123!@"
	hashedPassword := "hashed_SecurePass123!@"

	t.Run("should successfully hash password and create user", func(t *testing.T) {
		mockHasher.EXPECT().
			Hash(rawPassword).
			Return(hashedPassword, nil)

		mockRepo.EXPECT().
			CreateUser(ctx, gomock.AssignableToTypeOf(&domain.UserCreate{})).
			DoAndReturn(func(_ context.Context, u *domain.UserCreate) (*domain.User, error) {
				assert.Equal(t, hashedPassword, u.PasswordHash)
				return &domain.User{
					ID:       1,
					Username: u.Username,
					Email:    u.Email,
				}, nil
			})

		user, err := svc.CreateUser(ctx, userCreate, rawPassword)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, user.ID)
	})

	t.Run("should return error if password hashing fails", func(t *testing.T) {
		mockHasher.EXPECT().
			Hash(rawPassword).
			Return("", errors.New("hashing error"))

		user, err := svc.CreateUser(ctx, userCreate, rawPassword)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
