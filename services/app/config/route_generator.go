package core_config

import "github.com/kelseyhightower/envconfig"

type RouteGeneratorConfig struct {
	RouteGeneratorServerAddr string `envconfig:"SERVER_ADDR" default:"route-generator:8081"`
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
