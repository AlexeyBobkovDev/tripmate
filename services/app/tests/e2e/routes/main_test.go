//go:build e2e

package routes_e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"google.golang.org/grpc"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	core_config "github.com/AlexeyBobkovDev/tripmate/services/app/config"
	core_app "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/app"
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
	core_pgx_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/pgx"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

const (
	pathToMigrationsFiles = "file://../../../migrations/"
	pathToDotEnvFile      = "../../../../../.env.test"
)

type routeTestDeps struct {
	pool core_postgres_pool.Pool
}

var (
	testDeps                 routeTestDeps
	testApp                  *core_app.Application
	testRouteGeneratorServer *fakeRouteGeneratorServer
	freePort                 = fmt.Sprintf(":%d", getFreePort())
)

type fakeRouteGeneratorServer struct {
	route_generator_api_v1.UnimplementedRouteGeneratorServiceServer

	mu                sync.RWMutex
	generateRouteFunc func(
		ctx context.Context,
		request *route_generator_api_v1.GenerateRouteRequest,
	) (*route_generator_api_v1.GenerateRouteResponse, error)
}

func (s *fakeRouteGeneratorServer) GenerateRoute(
	ctx context.Context,
	request *route_generator_api_v1.GenerateRouteRequest,
) (*route_generator_api_v1.GenerateRouteResponse, error) {
	s.mu.RLock()
	generateRouteFunc := s.generateRouteFunc
	s.mu.RUnlock()

	if generateRouteFunc != nil {
		return generateRouteFunc(ctx, request)
	}

	return testutils.RouteGeneratorResponseTest(request.StartDate.AsTime()), nil
}

func (s *fakeRouteGeneratorServer) setGenerateRouteFunc(
	generateRouteFunc func(
		ctx context.Context,
		request *route_generator_api_v1.GenerateRouteRequest,
	) (*route_generator_api_v1.GenerateRouteResponse, error),
) func() {
	s.mu.Lock()
	previousGenerateRouteFunc := s.generateRouteFunc
	s.generateRouteFunc = generateRouteFunc
	s.mu.Unlock()

	return func() {
		s.mu.Lock()
		s.generateRouteFunc = previousGenerateRouteFunc
		s.mu.Unlock()
	}
}

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

	if err := os.Setenv("APP_POSTGRES_HOST", host); err != nil {
		panic(fmt.Errorf("failed to set APP_POSTGRES_HOST: %w", err))
	}
	if err := os.Setenv("APP_POSTGRES_PORT", port.Port()); err != nil {
		panic(fmt.Errorf("failed to set APP_POSTGRES_PORT: %w", err))
	}

	dsn := getDSN(ctx, postgresCtr)
	applyMigrations(pathToMigrationsFiles, dsn)
	testDeps.pool = setupPostgresPool(ctx, dsn)

	routeGeneratorServer, routeGeneratorListener := setupRouteGeneratorServer()
	defer routeGeneratorServer.GracefulStop()

	if err := os.Setenv(
		"ROUTE_GENERATOR_SERVER_TO_CONNECT",
		routeGeneratorListener.Addr().String(),
	); err != nil {
		panic(fmt.Errorf("failed to set ROUTE_GENERATOR_SERVER_TO_CONNECT: %w", err))
	}

	if err := os.Setenv("APP_SERVER_ADDR", freePort); err != nil {
		panic(fmt.Errorf("failed to set APP_SERVER_ADDR: %w", err))
	}

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
		http.MethodGet,
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
	var healthCheckErr error
	for {
		if counter == 10 {
			break
		}

		response, requestErr := client.Do(healthCheckRequest)
		healthCheckErr = requestErr
		if requestErr == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				break
			}
		}

		counter++
		time.Sleep(100 * time.Millisecond)
	}
	if counter >= 10 {
		panic(fmt.Errorf("failed to check if the app is healthy: %w", healthCheckErr))
	}

	return m.Run()
}

func setupRouteGeneratorServer() (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Errorf("failed to listen for route generator: %w", err))
	}

	server := grpc.NewServer()
	testRouteGeneratorServer = &fakeRouteGeneratorServer{}
	route_generator_api_v1.RegisterRouteGeneratorServiceServer(
		server,
		testRouteGeneratorServer,
	)

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			panic(fmt.Errorf("failed to run fake route generator server: %w", err))
		}
	}()

	return server, listener
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
