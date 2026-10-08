package users_service

import (
	"context"
	"fmt"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *UsersService) DeleteUser(
	ctx context.Context,
	userID int,
) error {
	if userID <= 0 {
		return fmt.Errorf("userID must be greater than 0: %w", core_errors.ErrBadRequest)
	}
	return s.usersRepository.DeleteUser(ctx, userID)
}
