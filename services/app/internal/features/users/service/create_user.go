package users_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user *domain.UserCreate,
	passwordToHash string,
) (*domain.User, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	if err := domain.ValidatePassword(passwordToHash); err != nil {
		return nil, err
	}

	password, err := s.passwordHasher.Hash(passwordToHash)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = password

	userResponse, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return userResponse, nil
}
