package config

import "time"

type DatabaseConfig struct {
	DSN             string        `env:"DATABASE_DSN"`
	MaxConns        int32         `env:"DATABASE_MAX_CONNS"`
	MinConns        int32         `env:"DATABASE_MIN_CONNS"`
	MaxConnLifetime time.Duration `env:"DATABASE_MAX_CONN_LIFETIME"`
	MaxConnIdleTime time.Duration `env:"DATABASE_MAX_CONN_IDLE_TIME"`
}
