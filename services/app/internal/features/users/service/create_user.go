package users_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user *domain.User,
	passwordToHash []byte,
) (*domain.User, error) {
	userResponse, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	_, err = s.passwordHasher.CreatePassword(ctx, userResponse.ID, passwordToHash)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	return userResponse, nil
}
