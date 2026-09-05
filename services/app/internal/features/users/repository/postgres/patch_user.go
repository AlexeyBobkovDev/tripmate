package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	userID int,
	user *domain.User,
) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE app.users
	SET
		version=version + 1,
		name=$3,
		surname=$4,
		username=$5,
		birth_date=$6,
		description=$7,
		email=$8,
		phone_number=$9
	WHERE id=$1 AND version=$2
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
		deleted_at;
	`

	var response domain.User
	err := r.pool.QueryRow(
		ctx,
		query,
		user.ID,
		user.Version,
		user.Name,
		user.Surname,
		user.Username,
		user.BirthDate,
		user.Description,
		user.Email,
		user.PhoneNumber,
	).Scan(
		&response.ID,
		&response.Version,
		&response.Name,
		&response.Surname,
		&response.Username,
		&response.BirthDate,
		&response.Description,
		&response.Email,
		&response.PhoneNumber,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.DeletedAt,
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

	return &response, nil
}
