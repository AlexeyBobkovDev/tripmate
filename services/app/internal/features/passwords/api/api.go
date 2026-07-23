package passwords_api

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

type PasswordManager interface {
	CreatePassword(
		ctx context.Context,
		userID int,
		password []byte,
	) (*domain.Password, error)
	GetPassword(
		ctx context.Context,
		userID int,
	) (*domain.Password, error)
}

type PasswordAPI struct {
	passwordService PasswordService
}

func NewPasswordAPI(
	passwordService PasswordService,
) *PasswordAPI {
	return &PasswordAPI{
		passwordService: passwordService,
	}
}

type PasswordService interface {
	CreatePassword(
		ctx context.Context,
		userID int,
		password []byte,
	) (*domain.Password, error)
	GetPassword(
		ctx context.Context,
		userID int,
	) (*domain.Password, error)
}
