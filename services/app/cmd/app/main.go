package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/AlexeyBobkovDev/tripmate/services/app/config"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_pgx_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/pgx"
	core_middleware "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/middleware"
	core_server "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/server"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords"
	passwords_api "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords/api"
	passwords_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords/repository/postgres"
	passwords_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords/service"
	users_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/repository/postgres"
	users_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service"
	users_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing new time zone")

	logger.Debug("initializing new postgres pool")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to initialize new postgres pool", zap.Error(err))
	}
	defer pool.Close()
	logger.Debug("successfully created new postgres pool")

	passwordHasherCfg := passwords.NewConfigMust()
	passwordHasher := passwords_service.NewArgon2IDHasher(
		passwordHasherCfg.Times,
		passwordHasherCfg.Memory,
		passwordHasherCfg.Threads,
		passwordHasherCfg.KeyLen,
	)

	passwordsPostgresRepository := passwords_postgres_repository.NewPasswordsPostgresRepository(pool)
	passwordsService := passwords_service.NewPasswordService(
		passwordsPostgresRepository,
		passwordHasher,
	)
	passwordsAPI := passwords_api.NewPasswordAPI(passwordsService)

	usersPostgresRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(
		usersPostgresRepository,
		passwordsAPI,
	)
	usersTransport := users_transport_http.NewUsersHTTPHandler(usersService)

	apiVersionRouterV1 := core_server.NewAPIRouter(
		core_server.APIVersionV1,
	)
	apiVersionRouterV1.RegisterRoutes(usersTransport.Routes()...)

	logger.Debug("initializing new server")
	server := core_server.NewHTTPServer(
		core_server.NewConfigMust(),
		logger,
		core_middleware.LoggerMiddleware(logger),
		core_middleware.RequestIDMiddleware(),
		core_middleware.RecoveryMiddleware(),
	)
	// TODO: rename the method from Health to RegisterHealthMethod or sth like that
	server.Health()
	logger.Debug("successfully initialized server")
	server.RegisterRouters(apiVersionRouterV1)

	server.Run(ctx)
}
