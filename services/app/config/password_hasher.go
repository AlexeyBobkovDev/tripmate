package core_config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type PasswordHasherConfig struct {
	Memory      uint32 `envconfig:"MEMORY"       required:"true"`
	Iterations  uint32 `envconfig:"ITERATIONS"   required:"true"`
	Parallelism uint8  `envconfig:"PARALLELISM"  required:"true"`
	SaltLength  uint32 `envconfig:"SALT_LENGTH"  required:"true"`
	KeyLength   uint32 `envconfig:"KEY_LENGTH"   required:"true"`
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
