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
	Env        `yaml:"env" env-required:"true"`
	HTTPConfig `yaml:"http" env-required:"true"`
	S3Config   `yaml:"s3Config"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type S3Config struct {
	BucketName string `yaml:"bucketName"`
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
