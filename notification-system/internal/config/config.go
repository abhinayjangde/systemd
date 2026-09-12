// internal/config/config.go
package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8000"`
	Env         string `env:"ENV" envDefault:"local"`
	DatabaseUrl string `env:"DATABASE_URL,required"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}
