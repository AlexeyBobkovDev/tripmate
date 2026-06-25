package main

import (
	"fmt"
	"os"
	"time"

	core_config "github.com/AlexeyBobkovDev/tripmate/services/app/config"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()
}
