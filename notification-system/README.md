# Notification System

building a notification system in Go, PostgreSQL

## packages used

1. github.com/joho/godotenv - Loading environment variables
2. github.com/caarlos0/env/v11 - 
3. github.com/golang-migrate/migrate/v4 - Database migrations
4. github.com/jackc/pgx/v5 - Postgres driver
5. github.com/segmentio/kafka-go - Kafka client
6. go get github.com/google/uuid - UUID generation

## Configuration

Copy `.env.example` to `.env` and set the runtime values before starting the
API, producer, consumer, or migration command.

Required Kafka variables:

- `KAFKA_BROKER_URL`
- `KAFKA_TOPIC`
- `KAFKA_CONSUMER_GROUP`