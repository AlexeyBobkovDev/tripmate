package core_postgres_pool

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	User      string    `envconfig:"USER" required:"true"`
	Password  string    `envconfig:"PASSWORD" required:"true"`
	Host      string    `envconfig:"HOST" required:"true"`
	Port      string    `envconfig:"PORT" default:"5432"`
	Name      string    `envconfig:"NAME" required:"true"`
	SSLMode   string    `envconfig:"SSLMODE" required:"true"`
	OpTimeout time.Time `envconfig:"OPTIMEOUT" required:"true"`
}

func (cfg Config) BuildDSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   "/" + cfg.Name,
	}

	if cfg.SSLMode != "" {
		q := u.Query()
		q.Set("sslmode", cfg.SSLMode)
		u.RawQuery = q.Encode()

	}

	return u.String()
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return Config{}, fmt.Errorf("failed to load postgres config environment variables")
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
