package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Node Node
	HTTP HTTP
}

type Node struct {
	Name string `env:"NODE_NAME" env-default:"short01"`
}

type HTTP struct {
	Port            string        `env:"HTTP_PORT" env-default:"8080"`
	DebugPort       string        `env:"HTTP_DEBUG_PORT" env-default:"9000"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"5s"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"3s"`
}

func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
