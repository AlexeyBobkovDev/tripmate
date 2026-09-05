package users_postgres_repository

import (
	"context"
	"fmt"
)

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

	query := `
	UPDATE app.users
	SET deleted_at = NOW(),
		updated_at = NOW(),
		version = version + 1
	WHERE id=$1
		AND deleted_at IS NULL;
	`
	cmdTag, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("update deleted_at field: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("the user %d does not exist", userID)
	}

	return tx.Commit(ctx)
}
