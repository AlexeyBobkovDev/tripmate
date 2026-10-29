package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/AlexeyBobkovDev/tripmate/services/app/config"
	_ "github.com/AlexeyBobkovDev/tripmate/services/app/docs"
	core_app "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/app"
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

	app, err := core_app.NewApplication(ctx, true)
	if err != nil {
		log.Printf("failed to initialize the application: %v", err)
	}
	defer func() {
		if closeErr := app.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close application: %w", closeErr))
			log.Printf("failed to close the application: %v", err)
		}
	}()

	return app.Run(ctx)
}
