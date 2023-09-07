package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

var (
	PathNotSetErr = errors.New("CONFIG_PATH is not set")
)

type Env string

const (
	Dev  Env = "dev"
	Prod Env = "prod"
)

type Config struct {
	Env          `yaml:"env" env-required:"true"`
	ServerConfig `yaml:"server" env-required:"true"`
}

type ServerConfig struct {
	Port         int           `yaml:"port" env-required:"true"`
	Timeout      time.Duration `yaml:"timeout" env-required:"true"`
	S3BucketName string        `yaml:"s3bucketName"`
}

func Load() (*Config, error) {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, PathNotSetErr
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, err
	}

	switch cfg.Env {
	case Dev:
	case Prod:
	default:
		return nil, fmt.Errorf("unexpected env value %v", cfg.Env)
	}

	return &cfg, nil
}
