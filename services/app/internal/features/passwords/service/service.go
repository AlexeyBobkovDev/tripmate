package passwords_service

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

// Password Repository

type PasswordRepository interface {
	CreatePassword(
		ctx context.Context,
		userID int,
		hash []byte,
		salt []byte,
		time uint32,
		memory uint32,
		threads uint8,
		keyLen uint32,
	) (*domain.Password, error)
	GetPassword(
		ctx context.Context,
		userID int,
	) (*domain.Password, error)
}

type PasswordService struct {
	passwordRepository PasswordRepository
	passwordHasher     *Argon2IDHasher
}

func NewPasswordService(
	passwordRepository PasswordRepository,
	passwordHasher *Argon2IDHasher,
) *PasswordService {
	return &PasswordService{
		passwordRepository: passwordRepository,
		passwordHasher:     passwordHasher,
	}
}

// Password Hasher

type Argon2IDHasher struct {
	Times   uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

func NewArgon2IDHasher(
	Times uint32,
	Memory uint32,
	Threads uint8,
	KeyLen uint32,
) *Argon2IDHasher {
	return &Argon2IDHasher{
		Times:   Times,
		Memory:  Memory,
		Threads: Threads,
		KeyLen:  KeyLen,
	}
}
