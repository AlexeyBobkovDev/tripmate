package core_config

import "github.com/kelseyhightower/envconfig"

type RouteGeneratorConfig struct {
	RouteGeneratorServerToConnect string `envconfig:"SERVER_TO_CONNECT" default:"route-generator:8081"`
}

func NewRouteGeneratorConfig() (RouteGeneratorConfig, error) {
	var config RouteGeneratorConfig
	if err := envconfig.Process("ROUTE_GENERATOR", &config); err != nil {
		return RouteGeneratorConfig{}, err
	}

	return config, nil
}

func NewRouteGeneratorConfigMust() RouteGeneratorConfig {
	config, err := NewRouteGeneratorConfig()
	if err != nil {
		panic(err)
	}

	return config
}
