package core_postgres_pool

import (
	"context"
)

type Pool interface {
	Begin(ctx context.Context) (Tx, error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
}

type Rows interface {
	Scan(dest ...any) error
	Next() bool
	Close()
	Err() error
}

type Row interface {
	Scan(dest ...any) error
}

type CommandTag interface {
	RowsAffected() int64
}

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	Exec(ctx context.Context, sql string, arguments ...any) (commandTag CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
}
