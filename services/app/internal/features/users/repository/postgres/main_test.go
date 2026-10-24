//go:build integration

package users_postgres_repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
	core_pgx_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/pgx"
)

const (
	pathToMigrationsFiles = "file://../../../../../migrations/"
)

var poolForIntegrationTesting core_postgres_pool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	postgresCtr := setupPostgres(ctx)

	defer func() {
		poolForIntegrationTesting.Close()

		if err := testcontainers.TerminateContainer(postgresCtr); err != nil {
			fmt.Printf("failed to terminate postgres container: %v\n", err)
		}
	}()

	dsn := getDSN(ctx, postgresCtr)
	applyMigrations(pathToMigrationsFiles, dsn)
	poolForIntegrationTesting = setupPostgresPool(ctx, dsn)

	os.Exit(m.Run())
}

func setupPostgres(ctx context.Context) *postgres.PostgresContainer {
	ctr, err := postgres.Run(
		ctx,
		"postgres:18",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("test-user"),
		postgres.WithPassword("test-password"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to start postgres container: %w", err))
	}

	return ctr
}

func getDSN(ctx context.Context, ctr *postgres.PostgresContainer) string {
	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("failed to get postgres container's conn string: %w", err))
	}
	return dsn
}

func setupPostgresPool(ctx context.Context, dsn string) core_postgres_pool.Pool {
	pool, err := core_pgx_pool.NewPool(ctx, dsn, core_pgx_pool.WithOperationTimeout(15*time.Second))
	if err != nil {
		panic(fmt.Errorf("failed to start postgres pool: %w", err))
	}

	return pool
}

func applyMigrations(sourceURL string, databaseURL string) {
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(fmt.Errorf("failed to apply migrations: %w", err))
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Errorf("failed to up migrations: %w", err))
	}
}
