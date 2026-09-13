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
)

func main() {
	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	consumer := kafka.NewConsumer(
		"localhost:9092",
		"notifications",
		"notification-service",
		db,
	)

	defer consumer.Close()

	log.Println("notification consumer started")

	if err := consumer.Consume(ctx); err != nil {
		log.Printf("consumer stopped: %v", err)
	}

	log.Println("notification consumer stopped")
}
