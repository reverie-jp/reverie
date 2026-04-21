package config

import "time"

type LiveKitConfig struct {
	URL       string        `env:"LIVEKIT_URL"`
	APIKey    string        `env:"LIVEKIT_API_KEY"`
	APISecret string        `env:"LIVEKIT_API_SECRET"`
	TokenTTL  time.Duration `env:"LIVEKIT_TOKEN_TTL" envDefault:"4h"`
}
