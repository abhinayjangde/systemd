package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abhinayjangde/notification-system/internal/config"
	"github.com/abhinayjangde/notification-system/internal/db"
	"github.com/abhinayjangde/notification-system/internal/kafka"
	"github.com/abhinayjangde/notification-system/internal/notification"
)

func main() {
	cfg := config.MustLoad()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repository := notification.NewPostgresRepository(database)

	service := notification.NewService(repository)

	consumer := kafka.NewConsumer(
		cfg.KafkaBrokerURL,
		cfg.KafkaTopic,
		cfg.KafkaConsumerGroup,
		service,
	)

	defer consumer.Close()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	log.Println("notification consumer started")

	if err := consumer.Consume(ctx); err != nil {
		log.Printf("consumer stopped: %v", err)
	}

	log.Println("notification consumer stopped")
}
