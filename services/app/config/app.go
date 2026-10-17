package core_config

import (
	"fmt"
	"os"
	"time"
)

type AppConfig struct {
	TimeZone *time.Location
}

func NewAppConfig() (AppConfig, error) {
	tz := os.Getenv("TIME_ZONE")
	if tz == "" {
		tz = "UTC"
	}
	zone, err := time.LoadLocation(tz)
	if err != nil {
		return AppConfig{}, fmt.Errorf("load time zone: %s, %w", tz, err)
	}

	return AppConfig{
		TimeZone: zone,
	}, nil
}

func NewAppConfigMust() AppConfig {
	config, err := NewAppConfig()
	if err != nil {
		panic(err)
	}

	return config
}
