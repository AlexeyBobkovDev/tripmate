package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	userID int,
) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT
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
	FROM app.users
	WHERE id=$1
		AND deleted_at IS NULL;
	`

	var userResponse domain.User
	err := r.pool.QueryRow(
		ctx,
		query,
		userID,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(
				"user id %d attempt to concurrently access: %w",
				userID,
				err,
			)
		}
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return &userResponse, nil
}
