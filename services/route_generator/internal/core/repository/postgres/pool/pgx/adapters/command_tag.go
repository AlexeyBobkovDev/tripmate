package core_pgx_pool_adapters

import "github.com/jackc/pgx/v5/pgconn"

type PgxCommandTag struct {
	pgconn.CommandTag
}
