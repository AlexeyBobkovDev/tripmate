package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

const deleteUserQuery = `
	UPDATE app.users
	SET deleted_at = NOW(),
		updated_at = NOW(),
		version = version + 1
	WHERE id=$1
		AND deleted_at IS NULL;
	`

func (s *UsersRepository) DeleteUser(
	ctx context.Context,
	userID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OpTimeout())
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create transaction to delete user: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	cmdTag, err := tx.Exec(ctx, deleteUserQuery, userID)
	if err != nil {
		return fmt.Errorf("update deleted_at field: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("the user %d does not exist: %w", userID, core_errors.ErrNotFound)
	}

	return tx.Commit(ctx)
}
