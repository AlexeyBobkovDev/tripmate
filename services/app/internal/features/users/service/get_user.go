package users_service

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	userID int,
) (*domain.User, error) {
	return s.usersRepository.GetUser(ctx, userID)
}
