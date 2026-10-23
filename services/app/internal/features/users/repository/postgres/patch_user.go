package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
)

const patchUserQuery = `
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
WHERE id=$1 AND version=$2 AND deleted_at IS NULL
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

const userExistsQuery = `
SELECT EXISTS(
	SELECT 1
	FROM app.users
	WHERE id=$1 and deleted_at IS NULL
);
`

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	userID int,
	user *domain.User,
) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var response domain.User
	err := r.pool.QueryRow(
		ctx,
		patchUserQuery,
		userID,
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
		if !errors.Is(err, core_postgres_pool_errors.ErrNoRows) {
			return nil, fmt.Errorf("update user: %w", err)
		}

		var exists bool

		checkErr := r.pool.QueryRow(ctx, userExistsQuery, userID).Scan(&exists)
		if checkErr != nil {
			return nil, fmt.Errorf(
				"check user %d after patch matched no rows: %w",
				userID,
				checkErr,
			)
		}
		if !exists {
			return nil, fmt.Errorf(
				"patch user %d: %w",
				userID,
				core_errors.ErrNotFound,
			)
		}

		return nil, fmt.Errorf(
			"patch user %d: version conflict: %w",
			userID,
			errors.Join(err, core_errors.ErrConflict),
		)
	}

	return &response, nil
}
