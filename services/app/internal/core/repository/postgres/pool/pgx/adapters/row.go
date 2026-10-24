package core_pgx_pool_adapters

import "github.com/jackc/pgx/v5"

type PgxRow struct {
	pgx.Row
}

func (r PgxRow) Scan(dest ...any) error {
	if err := r.Row.Scan(dest...); err != nil {
		return MapErrors(err)
	}

	return nil
}
