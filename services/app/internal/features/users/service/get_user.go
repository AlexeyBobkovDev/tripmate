package users_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	userID int,
) (*domain.User, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("userID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	return s.usersRepository.GetUser(ctx, userID)
}
