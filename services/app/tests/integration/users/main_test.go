//go:build integration

package users_integration_test

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

	core_crypto_argon2id "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/crypto/argon2id"
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
	core_pgx_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/pgx"
	users_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/repository/postgres"
	users_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

const (
	pathToMigrationsFiles = "file://../../../migrations/"
)

type userTestDeps struct {
	pool            core_postgres_pool.Pool
	usersRepository *users_postgres_repository.UsersRepository
	usersService    *users_service.UsersService
}

var testDeps = userTestDeps{}

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	ctx := context.Background()

	postgresCtr := setupPostgres(ctx)

	defer func() {
		testDeps.pool.Close()

		if err := testcontainers.TerminateContainer(postgresCtr); err != nil {
			fmt.Printf("failed to terminate postgres container: %v\n", err)
		}
	}()

	dsn := getDSN(ctx, postgresCtr)
	applyMigrations(pathToMigrationsFiles, dsn)
	testDeps.pool = setupPostgresPool(ctx, dsn)

	testDeps.usersRepository = users_postgres_repository.NewUsersRepository(testDeps.pool)
	usersPasswordManager := core_crypto_argon2id.NewArgon2IDHash(
		testutils.Memory,
		testutils.Iterations,
		testutils.Parallelism,
		testutils.SaltLength,
		testutils.KeyLength,
	)
	testDeps.usersService = users_service.NewUsersService(testDeps.usersRepository, usersPasswordManager)

	return m.Run()
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
