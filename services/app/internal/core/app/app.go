package core_app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_config "github.com/AlexeyBobkovDev/tripmate/services/app/config"
	_ "github.com/AlexeyBobkovDev/tripmate/services/app/docs"
	core_crypto_argon2id "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/crypto/argon2id"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
	core_pgx_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/pgx"
	core_middleware "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/middleware"
	core_server "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/server"
	users_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/repository/postgres"
	users_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service"
	users_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http"
)

type Application struct {
	Logger         *core_logger.Logger
	Pool           core_postgres_pool.Pool
	PasswordHasher users_service.PasswordManager

	server *core_server.HTTPServer
}

// ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
func (app *Application) Run(ctx context.Context) error {
	app.Logger.Debug("initializing new time zone")

	return app.server.Run(ctx)
}

func (app *Application) Handler() http.Handler {
	return app.server.Handler()
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

	passwordHasherCfg := core_config.NewPasswordHasherConfigMust()
	passwordHasher := core_crypto_argon2id.NewArgon2IDHash(
		passwordHasherCfg.Memory,
		passwordHasherCfg.Iterations,
		passwordHasherCfg.Parallelism,
		passwordHasherCfg.SaltLength,
		passwordHasherCfg.KeyLength,
	)

	usersPostgresRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(
		usersPostgresRepository,
		passwordHasher,
	)
	usersTransport := users_transport_http.NewUsersHTTPHandler(usersService)

	apiVersionRouterV1 := core_server.NewAPIRouter(
		core_server.APIVersionV1,
	)
	apiVersionRouterV1.RegisterRoutes(usersTransport.Routes()...)

	logger.Debug("initializing new server")

	middlewareConfig := core_middleware.NewConfigMust()
	// corsConfiguration := core_middleware.NewCORSConfigMust()
	server := core_server.NewHTTPServer(
		core_server.NewConfigMust(),
		logger,
		// core_middleware.CORSMiddleware(
		// 	corsConfiguration,
		// ),
		core_middleware.RateLimiterMiddleware(
			middlewareConfig.MaxReqAmount,
			middlewareConfig.WindowSize,
		),
		core_middleware.LoggerMiddleware(logger),
		core_middleware.RequestIDMiddleware(),
		core_middleware.TraceMiddleware(),
		core_middleware.RecoveryMiddleware(),
	)
	server.RegisterHealth()
	server.RegisterSwagger()
	logger.Debug("successfully initialized server")
	server.RegisterRouters(apiVersionRouterV1)

	return &Application{
		Logger:         logger,
		Pool:           pool,
		PasswordHasher: passwordHasher,
		server:         server,
	}, nil
}
