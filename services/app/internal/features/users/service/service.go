package users_service

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	passwords_api "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords/api"
)

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user *domain.User,
	) (*domain.User, error)
}

type UsersService struct {
	usersRepository UsersRepository
	passwordHasher  passwords_api.PasswordManager
}

func NewUsersService(
	usersRepository UsersRepository,
	passwordHasher passwords_api.PasswordManager,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		passwordHasher:  passwordHasher,
	}
}
