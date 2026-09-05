package users_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user *domain.UserCreate,
	passwordToHash []byte,
) (*domain.User, error) {
	userResponse, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return userResponse, nil
}
