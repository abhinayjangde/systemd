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
- `KAFKA_DLQ_TOPIC`
- `KAFKA_MAX_RETRIES`

The API accepts notifications asynchronously with `202 Accepted` after the
event is published to Kafka. The consumer persists the notification and then
commits the Kafka offset. Create both the configured main topic and DLQ topic
when Kafka is running, because topic auto-creation is disabled in Docker:

```powershell
docker exec notification-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --if-not-exists --topic notifications --partitions 3 --replication-factor 1
docker exec notification-kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --if-not-exists --topic notifications.dlq --partitions 3 --replication-factor 1
```