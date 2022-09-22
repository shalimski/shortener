package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App   App
	Node  Node
	HTTP  HTTP
	Mongo Mongo
	Redis Redis
}

type App struct {
	ShortURLLength int      `env:"SHORT_URL_LENGTH" env-default:"7"`
	EtcdEndpoints  []string `env:"ETCD_ENDPOINTS" env-default:"http://127.0.0.1:2379"`
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

type Mongo struct {
	Host     string `env:"MONGO_HOST" env-default:"localhost"`
	Port     string `env:"MONGO_PORT" env-default:"27017"`
	User     string `env:"MONGO_USER" env-default:"admin"`
	Password string `env:"MONGO_PASSWORD" env-default:"admin"`
	Database string `env:"MONGO_DATABASE" env-default:"shortener"`
}

type Redis struct {
	DSN      string `env:"REDIS_DSN" env-default:"127.0.0.1:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:"admin"`
}

func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
