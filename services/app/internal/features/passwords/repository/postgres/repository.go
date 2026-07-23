package passwords_postgres_repository

import core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"

type PasswordsPostgresRepository struct {
	pool core_postgres_pool.Pool
}

func NewPasswordsPostgresRepository(
	pool core_postgres_pool.Pool,
) *PasswordsPostgresRepository {
	return &PasswordsPostgresRepository{
		pool: pool,
	}
}
