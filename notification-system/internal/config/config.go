// internal/config/config.go
package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port               string `env:"PORT" envDefault:"8000"`
	Env                string `env:"ENV" envDefault:"local"`
	DatabaseURL        string `env:"DATABASE_URL,required"`
	KafkaBrokerURL     string `env:"KAFKA_BROKER_URL,required"`
	KafkaTopic         string `env:"KAFKA_TOPIC,required"`
	KafkaConsumerGroup string `env:"KAFKA_CONSUMER_GROUP,required"`
	KafkaDLQTopic      string `env:"KAFKA_DLQ_TOPIC,required"`
	KafkaMaxRetries    int    `env:"KAFKA_MAX_RETRIES" envDefault:"3"`
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
