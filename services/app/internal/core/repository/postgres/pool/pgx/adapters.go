package core_pgx_pool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
)

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	if err := r.Row.Scan(dest...); err != nil {
		return mapErrors(err)
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			return fmt.Errorf(
				"%v: %w",
				err,
				core_postgres_pool.ErrViolatedForeignKey,
			)
		case pgerrcode.UniqueViolation:
			return fmt.Errorf(
				"%v: %w",
				err,
				core_errors.ErrConflict,
			)
		}
	}

	return fmt.Errorf(
		"%v: %w",
		err,
		core_postgres_pool.ErrUnknown,
	)
}
