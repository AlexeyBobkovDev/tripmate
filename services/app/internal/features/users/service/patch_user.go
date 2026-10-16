package users_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	userID int,
	userPatch *domain.UserPatch,
) (*domain.User, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("userID must be greater than zero: %w", core_errors.ErrBadRequest)
	}
	user, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.Version != userPatch.Version {
		return nil, fmt.Errorf(
			"patch user %d: expected version %d, current version %d: %w",
			userID,
			userPatch.Version,
			user.Version,
			core_errors.ErrConflict,
		)
	}
	if err := user.ApplyPatch(userPatch); err != nil {
		return nil, fmt.Errorf("apply user patch: %w", errors.Join(err, core_errors.ErrBadRequest))
	}

	user, err = s.usersRepository.PatchUser(ctx, userID, user)
	if err != nil {
		return nil, fmt.Errorf("patch user: %w", err)
	}

	return user, nil
}
