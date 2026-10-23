package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
)

const getUserQuery = `
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

func (r *UsersRepository) GetUser(
	ctx context.Context,
	userID int,
) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var userResponse domain.User
	err := r.pool.QueryRow(
		ctx,
		getUserQuery,
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
		if errors.Is(err, core_postgres_pool_errors.ErrNoRows) {
			return nil, fmt.Errorf(
				"get userID %d: %w",
				userID,
				errors.Join(err, core_errors.ErrNotFound),
			)
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &userResponse, nil
}
