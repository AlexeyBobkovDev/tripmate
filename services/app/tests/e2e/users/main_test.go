//go:build e2e

package users_e2e_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	core_config "github.com/AlexeyBobkovDev/tripmate/services/app/config"
	core_app "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/app"
	core_crypto_argon2id "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/crypto/argon2id"
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
	core_pgx_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/pgx"
	users_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/repository/postgres"
	users_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service"
	users_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

const (
	pathToMigrationsFiles = "file://../../../migrations/"
	pathToDotEnvFile      = "../../../../../.env.test"
)

type userTestDeps struct {
	pool               core_postgres_pool.Pool
	usersRepository    *users_postgres_repository.UsersRepository
	usersService       *users_service.UsersService
	usersTransportHTTP *users_transport_http.UsersHTTPHandler
}

var (
	testDeps = userTestDeps{}
	testApp  *core_app.Application
	freePort = fmt.Sprintf(":%d", getFreePort())
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func getFreePort() int {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func runTests(m *testing.M) int {
	cfg := core_config.NewAppConfigMust()
	time.Local = cfg.TimeZone

	ctx := context.Background()

	if err := godotenv.Load(pathToDotEnvFile); err != nil {
		panic(fmt.Errorf("failed to load .env file: %w", err))
	}

	postgresCtr := setupPostgres(ctx)

	defer func() {
		testDeps.pool.Close()

		if err := testcontainers.TerminateContainer(postgresCtr); err != nil {
			panic(fmt.Errorf("failed to terminate postgres container: %w\n", err))
		}
	}()

	host, err := postgresCtr.Host(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get PostgreSQL container's host: %w", err))
	}
	port, err := postgresCtr.MappedPort(ctx, "5432/tcp")
	if err != nil {
		panic(fmt.Errorf("failed to get PostgreSQL container's port: %w", err))
	}

	os.Setenv("POSTGRES_HOST", host)
	os.Setenv("POSTGRES_PORT", port.Port())

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
	testDeps.usersTransportHTTP = users_transport_http.NewUsersHTTPHandler(testDeps.usersService)

	os.Setenv("SERVER_ADDR", freePort)

	app, err := core_app.NewApplication(ctx, false)
	if err != nil {
		panic(fmt.Errorf("can not initialize the app: %w", err))
	}
	defer func() {
		if err := testApp.Close(); err != nil && err != http.ErrServerClosed {
			fmt.Println("application error: %w", err)
		}
	}()
	testApp = app

	go func() {
		if err := testApp.Run(ctx); err != nil {
			panic("failed to start testApp")
		}
	}()

	healthCheckRequest, err := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost%s/health", freePort),
		nil,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create healthCheckRequest: %w", err))
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	counter := 0
	for {
		if counter == 10 {
			break
		}
		response, err := client.Do(healthCheckRequest)

		if err == nil && response.StatusCode == http.StatusOK {
			break
		}
		counter++
		time.Sleep(100 * time.Millisecond)
	}
	if counter >= 10 {
		panic(fmt.Errorf("failed to check if the app is health: %w", err))
	}

	return m.Run()
}

func setupPostgres(ctx context.Context) *postgres.PostgresContainer {
	ctr, err := postgres.Run(
		ctx,
		"postgres:18",
		postgres.WithDatabase(os.Getenv("POSTGRES_DB")),
		postgres.WithUsername(os.Getenv("POSTGRES_USER")),
		postgres.WithPassword(os.Getenv("POSTGRES_PASSWORD")),
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
