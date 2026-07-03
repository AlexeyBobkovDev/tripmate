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
	core_postgres_pool "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool"
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
	pool, err := core_postgres_pool.NewPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to initialize new postgres pool", zap.Error(err))
	}
	defer pool.Close()
	logger.Debug("successfully created new postgres pool")
}
