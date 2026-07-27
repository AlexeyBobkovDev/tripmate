package passwords_service

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (s *PasswordService) GetPassword(
	ctx context.Context,
	userID int,
) (*domain.Password, error) {
	return s.passwordRepository.GetPassword(ctx, userID)
}
