package core_postgres_pool

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	User      string    `envconfig:"User" required:"true"`
	Password  string    `envconfig:"Password" required:"true"`
	Host      string    `envconfig:"Host" required:"true"`
	Port      string    `envconfig:"Port" default:"5432"`
	Name      string    `envconfig:"Name" required:"true"`
	SSLMode   string    `envconfig:"SSLMode" required:"true"`
	OpTimeout time.Time `envconfig:"OpTimeout" required:"true"`
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
