//go:build e2e

package trips_e2e_test

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
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
	core_pgx_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/pgx"
)

const (
	pathToMigrationsFiles = "file://../../../migrations/"
	pathToDotEnvFile      = "../../../.env.test"
)

type tripTestDeps struct {
	pool core_postgres_pool.Pool
}

var (
	testDeps tripTestDeps
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

	if err := godotenv.Load(pathToDotEnvFile); err != nil {
		panic(fmt.Errorf("failed to load .env file: %w", err))
	}

	ctx := context.Background()
	postgresCtr := setupPostgres(ctx)

	defer func() {
		if testDeps.pool != nil {
			testDeps.pool.Close()
		}

		if err := testcontainers.TerminateContainer(postgresCtr); err != nil {
			fmt.Printf("failed to terminate PostgreSQL container: %v\n", err)
		}
	}()

	host, err := postgresCtr.Host(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get PostgreSQL host: %w", err))
	}

	port, err := postgresCtr.MappedPort(ctx, "5432/tcp")
	if err != nil {
		panic(fmt.Errorf("failed to get PostgreSQL port: %w", err))
	}

	if err := os.Setenv("POSTGRES_HOST", host); err != nil {
		panic(fmt.Errorf("failed to set POSTGRES_HOST: %w", err))
	}
	if err := os.Setenv("POSTGRES_PORT", port.Port()); err != nil {
		panic(fmt.Errorf("failed to set POSTGRES_PORT: %w", err))
	}

	dsn := getDSN(ctx, postgresCtr)
	applyMigrations(pathToMigrationsFiles, dsn)
	testDeps.pool = setupPostgresPool(ctx, dsn)

	os.Setenv("SERVER_ADDR", freePort)

	app, err := core_app.NewApplication(ctx, false)
	if err != nil {
		panic(fmt.Errorf("failed to initialize application: %w", err))
	}
	testApp = app

	defer func() {
		if err := testApp.Close(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("failed to close application: %v\n", err)
		}
	}()

	go func() {
		if err := testApp.Run(ctx); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to run test application: %w", err))
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
		panic(fmt.Errorf("failed to start PostgreSQL container: %w", err))
	}

	return ctr
}

func getDSN(ctx context.Context, ctr *postgres.PostgresContainer) string {
	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("failed to get PostgreSQL connection string: %w", err))
	}

	return dsn
}

func setupPostgresPool(
	ctx context.Context,
	dsn string,
) core_postgres_pool.Pool {
	pool, err := core_pgx_pool.NewPool(
		ctx,
		dsn,
		core_pgx_pool.WithOperationTimeout(15*time.Second),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create PostgreSQL pool: %w", err))
	}

	return pool
}

func applyMigrations(sourceURL, databaseURL string) {
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(fmt.Errorf("failed to initialize migrations: %w", err))
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Errorf("failed to apply migrations: %w", err))
	}
}
