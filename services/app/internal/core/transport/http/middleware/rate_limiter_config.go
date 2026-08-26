package core_middleware

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MaxReqAmount int           `envconfig:"MAX_REQUESTS_AMOUNT" required:"true" validate:"gt=0"`
	WindowSize   time.Duration `envconfig:"WINDOW_SIZE" required:"true" validate:"gt=0"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
