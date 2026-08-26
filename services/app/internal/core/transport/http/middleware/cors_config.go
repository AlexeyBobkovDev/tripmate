package core_middleware

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type CORSConfigToProcess struct {
	Origins string `envconfig:"ORIGINS"`
	Methods string `envconfig:"METHODS"`
	Headers string `envconfig:"HEADERS"`
}

type CORSConfig struct {
	Origins []string `envconfig:"ORIGINS"`
	Methods []string `envconfig:"METHODS"`
	Headers []string `envconfig:"HEADERS"`
}

func NewCORSConfig() (CORSConfig, error) {
	var configToProcess CORSConfigToProcess
	if err := envconfig.Process("CORS_ALLOWED", &configToProcess); err != nil {
		return CORSConfig{}, fmt.Errorf("retrieve CORS_ALLOWED .env variables: %w", err)
	}
	corsConfig := CORSConfig{
		Origins: strings.Split(configToProcess.Origins, ","),
		Methods: strings.Split(configToProcess.Methods, ","),
		Headers: strings.Split(configToProcess.Headers, ","),
	}
	return corsConfig, nil
}

func NewCORSConfigMust() CORSConfig {
	config, err := NewCORSConfig()
	if err != nil {
		panic(err)
	}

	return config
}
