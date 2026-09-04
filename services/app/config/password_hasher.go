package core_config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type PasswordHasherConfig struct {
	Times   uint32 `envconfig:"TIMES" required:"true"`
	Memory  uint32 `envconfig:"MEMORY" required:"true"`
	Threads uint8  `envconfig:"THREADS" required:"true"`
	KeyLen  uint32 `envconfig:"KEYLEN" required:"true"`
}

func NewPasswordHasherConfig() (PasswordHasherConfig, error) {
	var config PasswordHasherConfig
	if err := envconfig.Process("PASSWORDHASHER", &config); err != nil {
		return PasswordHasherConfig{}, fmt.Errorf("failed to get config variables for password hasher")
	}

	return config, nil
}

func NewPasswordHasherConfigMust() PasswordHasherConfig {
	config, err := NewPasswordHasherConfig()
	if err != nil {
		panic(err)
	}

	return config
}
