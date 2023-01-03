package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" default:"postgres://hubla:supersecret@0.0.0.0/hubla?sslmode=disable"`
	StaticPath  string `envconfig:"STATIC_PATH" default:"./static"`
}

func New() Config {
	var env Config
	_ = envconfig.Process("", &env)

	return env
}
