package passwords

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Times   uint32 `envconfig:"TIMES" required:"true"`
	Memory  uint32 `envconfig:"MEMORY" required:"true"`
	Threads uint8  `envconfig:"THREADS" required:"true"`
	KeyLen  uint32 `envconfig:"KEYLEN" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("PASSWORDHASHER", &config); err != nil {
		return Config{}, fmt.Errorf("failed to get config variables for password hasher")
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return config
}
