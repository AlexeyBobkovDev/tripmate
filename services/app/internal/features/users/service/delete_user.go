package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUser(
	ctx context.Context,
	userID int,
) error {
	if userID <= 0 {
		return fmt.Errorf("userID must be greater than 0")
	}
	return s.usersRepository.DeleteUser(ctx, userID)
}
