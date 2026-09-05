package users_service

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

type PasswordManager interface {
	Hash(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user *domain.UserCreate,
	) (*domain.User, error)
	GetUser(
		ctx context.Context,
		userID int,
	) (*domain.User, error)
	DeleteUser(
		ctx context.Context,
		userID int,
	) error
	PatchUser(
		ctx context.Context,
		userID int,
		user *domain.User,
	) (*domain.User, error)
}

type UsersService struct {
	usersRepository UsersRepository
	passwordHasher  PasswordManager
}

func NewUsersService(
	usersRepository UsersRepository,
	passwordHasher PasswordManager,
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
		passwordHasher:  passwordHasher,
	}
}
