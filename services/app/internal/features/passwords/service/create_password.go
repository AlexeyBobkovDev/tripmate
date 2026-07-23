package passwords_service

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (s *PasswordService) CreatePassword(
	ctx context.Context,
	userID int,
	password []byte,
) (*domain.Password, error) {
	hash, salt, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return s.passwordRepository.CreatePassword(ctx, userID, hash, salt, s.passwordHasher.Times, s.passwordHasher.Memory, s.passwordHasher.Threads, s.passwordHasher.KeyLen)
}
