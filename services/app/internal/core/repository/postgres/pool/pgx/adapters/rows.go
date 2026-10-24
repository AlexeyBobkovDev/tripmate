package core_pgx_pool_adapters

import (
	"github.com/jackc/pgx/v5"
)

type PgxRows struct {
	pgx.Rows
}

func (r PgxRows) Scan(dest ...any) error {
	if err := r.Rows.Scan(dest...); err != nil {
		return MapErrors(err)
	}

	return nil
}

func (r PgxRows) Next() bool {
	return r.Rows.Next()
}

func (r PgxRows) Close() {
	r.Rows.Close()
}

func (r PgxRows) Err() error {
	return MapErrors(r.Rows.Err())
}

func (r PgxRows) CommandTag() PgxCommandTag {
	return PgxCommandTag{r.Rows.CommandTag()}
}
