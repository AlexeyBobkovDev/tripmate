package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	core_config "github.com/AlexeyBobkovDev/tripmate/services/app/config"
	_ "github.com/AlexeyBobkovDev/tripmate/services/app/docs"
	core_crypto_argon2id "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/crypto/argon2id"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_pgx_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/pgx"
	core_middleware "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/middleware"
	core_server "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/server"
	users_postgres_repository "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/repository/postgres"
	users_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/service"
	users_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http"
)

//	@Title			TripMate API
//	@Version		1.0
//	@Description	TripMate is the service that helps you plan your route and just travel
//	@Description	without difficulties
//	@Contact.name	TripMate support
//	@Contact.email	tripmatesupport@gmail.com
//
// Contact.url  https://tripmate./support
//
//	@Host			localhost:8080
//	@BasePath		/api/v1
//	@Accept			json
//	@Produce		json

func main() {
	cfg := core_config.NewAppConfigMust()
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

	poolCfg := core_pgx_pool.NewConfigMust()
	pool, err := core_pgx_pool.NewPool(ctx, poolCfg)
	if err != nil {
		logger.Fatal("failed to initialize new postgres pool", zap.Error(err))
	}
	defer pool.Close()
	logger.Debug("successfully created new postgres pool")

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
	corsConfiguration := core_middleware.NewCORSConfigMust()
	server := core_server.NewHTTPServer(
		core_server.NewConfigMust(),
		logger,
		core_middleware.CORSMiddleware(
			corsConfiguration,
		),
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

	server.Run(ctx)
}
