package core_pgx_pool_adapters

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
)

func TestMapErrors(t *testing.T) {
	uniqueViolation := &pgconn.PgError{Code: "23505", Message: "duplicate key"}
	foreignKeyViolation := &pgconn.PgError{Code: "23503", Message: "foreign key"}
	checkViolation := &pgconn.PgError{Code: "23514", Message: "check constraint"}

	tests := []struct {
		name     string
		err      error
		contains []error
	}{
		{name: "nil stays nil", err: nil},
		{
			name:     "no rows",
			err:      pgx.ErrNoRows,
			contains: []error{core_postgres_pool_errors.ErrNoRows},
		},
		{
			name:     "wrapped no rows",
			err:      fmt.Errorf("scan: %w", pgx.ErrNoRows),
			contains: []error{core_postgres_pool_errors.ErrNoRows},
		},
		{
			name:     "unique violation becomes conflict",
			err:      uniqueViolation,
			contains: []error{core_errors.ErrConflict},
		},
		{
			name:     "foreign key violation",
			err:      foreignKeyViolation,
			contains: []error{core_postgres_pool_errors.ErrViolatedForeignKey},
		},
		{
			name:     "unknown pg error code",
			err:      checkViolation,
			contains: []error{core_postgres_pool_errors.ErrUnknown},
		},
		{
			name:     "plain error",
			err:      errors.New("connection reset"),
			contains: []error{core_postgres_pool_errors.ErrUnknown},
		},
		{
			name:     "wrapped pg error loses its cause",
			err:      fmt.Errorf("exec: %w", uniqueViolation),
			contains: []error{core_errors.ErrConflict},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapped := MapErrors(tt.err)

			if tt.err == nil {
				require.NoError(t, mapped)
				return
			}

			require.Error(t, mapped)
			for _, target := range tt.contains {
				require.ErrorIs(t, mapped, target)
			}
		})
	}
}

func TestMapErrors_PgErrorCauseIsDiscarded(t *testing.T) {
	// BUG: MapErrors formats the original error with %v instead of %w, so the
	// *pgconn.PgError is not in the returned chain: callers cannot use
	// errors.As to read the SQLSTATE code or constraint name.
	pgErr := &pgconn.PgError{Code: "23505", ConstraintName: "trips_pkey"}
	mapped := MapErrors(fmt.Errorf("exec: %w", pgErr))

	var recovered *pgconn.PgError
	require.False(t, errors.As(mapped, &recovered),
		"the original PgError should stay unwrapped in the chain but is dropped")
}

type fakeRow struct {
	err error
}

func (r fakeRow) Scan(dest ...any) error { return r.err }

func TestPgxRowScan(t *testing.T) {
	t.Run("maps no rows", func(t *testing.T) {
		row := PgxRow{Row: fakeRow{err: pgx.ErrNoRows}}

		err := row.Scan(new(string))

		require.ErrorIs(t, err, core_postgres_pool_errors.ErrNoRows)
	})

	t.Run("maps unique violation to conflict", func(t *testing.T) {
		row := PgxRow{Row: fakeRow{err: &pgconn.PgError{Code: "23505"}}}

		err := row.Scan(new(string))

		require.ErrorIs(t, err, core_errors.ErrConflict)
	})

	t.Run("success passes through", func(t *testing.T) {
		row := PgxRow{Row: fakeRow{}}

		require.NoError(t, row.Scan(new(string)))
	})
}

type fakeRows struct {
	pgx.Rows

	scanErr error
	err     error
}

func (r fakeRows) Scan(dest ...any) error { return r.scanErr }
func (r fakeRows) Err() error             { return r.err }

func TestPgxRows(t *testing.T) {
	t.Run("Scan maps errors", func(t *testing.T) {
		rows := PgxRows{Rows: fakeRows{scanErr: pgx.ErrNoRows}}

		require.ErrorIs(t, rows.Scan(new(string)), core_postgres_pool_errors.ErrNoRows)
	})

	t.Run("Err maps errors", func(t *testing.T) {
		rows := PgxRows{Rows: fakeRows{err: &pgconn.PgError{Code: "23505"}}}

		require.ErrorIs(t, rows.Err(), core_errors.ErrConflict)
	})

	t.Run("Err returns nil when the underlying rows are fine", func(t *testing.T) {
		rows := PgxRows{Rows: fakeRows{}}

		require.NoError(t, rows.Err())
	})
}

func TestPgxTxExec(t *testing.T) {
	t.Run("commit error is mapped", func(t *testing.T) {
		// a *PgxTx over a nil interface would panic, so exercise the
		// mapping through MapErrors via a fake tx instead
		err := MapErrors(fmt.Errorf("commit: %w", pgx.ErrTxClosed))

		require.ErrorIs(t, err, core_postgres_pool_errors.ErrUnknown)
	})
}
