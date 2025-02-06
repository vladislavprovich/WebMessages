package config

import (
	"context"
	"fmt"
	"os"

	"messenger/internal/server"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server server.Config `yaml:"server"`
	Log    LoggerLevel   `yaml:"log"`
}

type LoggerLevel struct {
	Level string `yaml:"level" envconfig:"LOG_LEVEL" default:"development"`
}

func LoadConfig(_ context.Context) (*Config, error) {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file (%s): %w", configPath, err)
	}

	if err = yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err = envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load env config: %w", err)
	}

	return &cfg, nil
}
