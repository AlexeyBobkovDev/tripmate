package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {
	query := `
	INSERT INTO app.users (name, surname, username, birth_date, description, email, phone_number)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, version, name, surname, username, birth_date, description, email, phone_number, created_at, updated_at, deleted_at
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.Name,
		user.Surname,
		user.Username,
		user.BirthDate,
		user.Description,
		user.Email,
		user.PhoneNumber,
	)

	var userResponse domain.User

	err := row.Scan(
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
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &userResponse, nil
}
