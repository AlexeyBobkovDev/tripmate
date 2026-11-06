package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/AlexeyBobkovDev/tripmate/services/route_generator/config"
	core_app "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/app"
)

func main() {
	if err := run(); err != nil {
		log.Printf("application failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := core_config.NewAppConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app, err := core_app.NewApplication(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize the application: %v", err)
	}
	defer func() {
		if closeErr := app.Close(); closeErr != nil {
			log.Printf("failed to close the application: %v", closeErr)
		}
	}()

	return app.Run(ctx)
}
