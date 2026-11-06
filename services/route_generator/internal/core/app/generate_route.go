package core_app

import (
	"context"
	"errors"
	"fmt"
	"log"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/logger"
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/repository/postgres/pool"
	core_pgx_pool "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/repository/postgres/pool/pgx"
	core_interceptor "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/transport/grpc/interceptors"
	core_server "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/transport/grpc/server"
	generate_route_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/features/generate-route/repository/postgres"
	generate_route_service "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/features/generate-route/service"
	transport_grpc_generate_route "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/features/generate-route/transport/grpc"
)

type Application struct {
	Logger *core_logger.Logger
	Pool   core_postgres_pool.Pool

	server *core_server.GRPCServer
}

// ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
func (app *Application) Run(ctx context.Context) error {
	app.Logger.Debug("initializing new time zone")

	return app.server.Run(ctx)
}

func (app *Application) Close() error {
	if app.Pool != nil {
		app.Pool.Close()
	}

	if app.Logger != nil {
		return app.Logger.Close()
	}
	return nil
}

func NewApplication(ctx context.Context) (*Application, error) {
	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	logger.Debug("initializing new postgres pool")
	poolCfg := core_pgx_pool.NewConfigMust()
	pool, err := core_pgx_pool.NewPool(ctx, poolCfg)
	if err != nil {
		if closeLoggerErr := logger.Close(); closeLoggerErr != nil {
			return nil, errors.Join(
				closeLoggerErr,
				fmt.Errorf("failed to close logger"),
				fmt.Errorf("initialize PostgreSQL pool: %w", err),
			)
		}
		return nil, fmt.Errorf("initialize PostgreSQL pool: %w", err)
	}

	logger.Debug("initializing new gRPC server")

	grpcServer := core_server.NewGRPCServer(
		core_server.NewConfigMust(),
		logger,
		core_interceptor.LoggerInterceptor(logger),
		core_interceptor.RequestIDInterceptor(),
		core_interceptor.TraceInterceptor(),
		core_interceptor.RecoveryInterceptor(),
	)

	// Register health
	grpcServer.RegisterHealth()

	// Register generate-route feature
	generatorRoutesRepository := generate_route_postgres_repository.NewGenerateRouteRepository(pool)
	generatorRoutesService := generate_route_service.NewGenerateRouteService(generatorRoutesRepository)
	generatorRoutesGRPCHandler := transport_grpc_generate_route.NewGenerateRouteServer(generatorRoutesService)
	route_generator_api_v1.RegisterRouteGeneratorServiceServer(grpcServer, generatorRoutesGRPCHandler)

	log.Println("gRPC server is initialized successfully at port: :8081")

	logger.Debug("successfully initialized server")

	return &Application{
		Logger: logger,
		Pool:   pool,
		server: grpcServer,
	}, nil
}
