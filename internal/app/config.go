package app

import (
	"github.com/caarlos0/env"

	"github.com/fr33dman/go-template/internal/clients/embeddings"
)

type (
	// Config - collects hole application configuration, separated by groups (example: ServerConfig - server settings)
	Config struct {
		Logger     LoggerConfig
		Server     ServerConfig
		Embeddings embeddings.Config
	}
	// LoggerConfig - logger group settings
	LoggerConfig struct {
		Level string `env:"LOGGER_LEVEL" envDefault:"info"`
	}
	// ServerConfig - server group settings
	ServerConfig struct {
		Host            string `env:"SERVER_HOST"             envDefault:"0.0.0.0"`
		APIPort         int    `env:"SERVER_API_PORT"         envDefault:"8000"`
		HealthcheckPort int    `env:"SERVER_HEALTHCHECK_PORT" envDefault:"8081"`
		MetricsPort     int    `env:"SERVER_METRICS_PORT"     envDefault:"9090"`
	}
)

func NewConfig() (Config, error) {
	var (
		loggerCfg LoggerConfig
		serverCfg ServerConfig
	)

	if err := env.Parse(&serverCfg); err != nil {
		return Config{}, err
	}

	if err := env.Parse(&loggerCfg); err != nil {
		return Config{}, err
	}
	embeddingsCfg, err := embeddings.LoadConfigFromEnv()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Server:     serverCfg,
		Logger:     loggerCfg,
		Embeddings: embeddingsCfg,
	}, nil
}
