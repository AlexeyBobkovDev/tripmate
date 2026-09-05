package users_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	userID int,
	userPatch *domain.UserPatch,
) (*domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if err := user.ApplyPatch(userPatch); err != nil {
		return nil, fmt.Errorf("apply user patch: %w", err)
	}

	user, err = s.usersRepository.PatchUser(ctx, userID, user)
	if err != nil {
		return nil, fmt.Errorf("patch user: %w", err)
	}

	return user, nil
}
