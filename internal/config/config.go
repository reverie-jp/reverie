package config

import (
	"github.com/caarlos0/env/v11"
)

type Env string

const (
	EnvDevelopment Env = "development"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

type Config struct {
	Env      Env `env:"ENVIRONMENT"`
	Auth     AuthConfig
	Database DatabaseConfig
	Log      LogConfig
	Server   ServerConfig
}

func New() *Config {
	return &Config{}
}

func (c *Config) LoadFromEnv() error {
	if err := env.Parse(c); err != nil {
		return err
	}
	if err := env.Parse(&c.Auth); err != nil {
		return err
	}
	if err := env.Parse(&c.Database); err != nil {
		return err
	}
	if err := env.Parse(&c.Log); err != nil {
		return err
	}
	if err := env.Parse(&c.Server); err != nil {
		return err
	}
	return nil
}
