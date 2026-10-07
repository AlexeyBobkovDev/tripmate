package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

const query = `
INSERT INTO app.users (
	name,
	surname,
	username,
	birth_date,
	description,
	email,
	phone_number,
	password_hash
)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
	id,
	version,
	name,
	surname,
	username,
	birth_date,
	description,
	email,
	phone_number,
	created_at,
	updated_at,
	deleted_at
`

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user *domain.UserCreate,
) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var userResponse domain.User
	err := r.pool.QueryRow(
		ctx,
		query,
		user.Name,
		user.Surname,
		user.Username,
		user.BirthDate,
		user.Description,
		user.Email,
		user.PhoneNumber,
		user.PasswordHash,
	).Scan(
		&userResponse.ID,
		&userResponse.Version,
		&userResponse.Name,
		&userResponse.Surname,
		&userResponse.Username,
		&userResponse.BirthDate,
		&userResponse.Description,
		&userResponse.Email,
		&userResponse.PhoneNumber,
		&userResponse.CreatedAt,
		&userResponse.UpdatedAt,
		&userResponse.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return &userResponse, nil
}
