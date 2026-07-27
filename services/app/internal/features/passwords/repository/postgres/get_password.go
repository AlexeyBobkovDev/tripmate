package passwords_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (r *PasswordsPostgresRepository) GetPassword(
	ctx context.Context,
	userID int,
) (*domain.Password, error) {
	query := `
	SELECT user_id, version, password_hash, salt, times, memory, threads, key_len
	FROM app.passwords
	WHERE user_id=$1;
	`

	var password domain.Password

	row := r.pool.QueryRow(ctx, query, userID)
	err := row.Scan(
		&password.UserID,
		&password.Version,
		&password.Hash,
		&password.Salt,
		&password.Times,
		&password.Memory,
		&password.Threads,
		&password.KeyLen,
	)
	if err != nil {
		return nil, fmt.Errorf("get password data: %w", err)
	}

	return &password, nil
}
