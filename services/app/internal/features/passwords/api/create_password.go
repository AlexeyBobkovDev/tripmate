package passwords_api

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (s *PasswordAPI) CreatePassword(
	ctx context.Context,
	userID int,
	password []byte,
) (*domain.Password, error) {
	return s.passwordService.CreatePassword(ctx, userID, password)
}
