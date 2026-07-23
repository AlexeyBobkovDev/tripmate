package passwords_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (r *PasswordsPostgresRepository) CreatePassword(
	ctx context.Context,
	userID int,
	hash []byte,
	salt []byte,
	times uint32,
	memory uint32,
	threads uint8,
	keyLen uint32,
) (*domain.Password, error) {
	query := `
	INSERT INTO app.passwords (user_id, password_hash, salt, times, memory, threads, key_len)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING user_id, version, password_hash, salt, times, memory, threads, key_len;
	`
	var passwordHash domain.Password
	row := r.pool.QueryRow(ctx, query, userID, hash, salt, times, memory, threads, keyLen)
	err := row.Scan(
		&passwordHash.UserID,
		&passwordHash.Version,
		&passwordHash.Hash,
		&passwordHash.Salt,
		&passwordHash.Times,
		&passwordHash.Memory,
		&passwordHash.Threads,
		&passwordHash.KeyLen,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create password for userID=`%d`: %w", userID, err)
	}
	return &passwordHash, nil
}
