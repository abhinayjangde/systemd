package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	"github.com/abhinayjangde/notification-system/internal/config"
	"github.com/abhinayjangde/notification-system/internal/db"
	"github.com/abhinayjangde/notification-system/internal/handlers"
	"github.com/abhinayjangde/notification-system/internal/kafka"
)

func main() {
	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Kafka producer
	producer := kafka.NewProducer(
		cfg.KafkaBrokerURL,
		cfg.KafkaTopic,
	)
	defer producer.Close()

	// mux is the HTTP request multiplexer that matches incoming requests to their respective handler functions.
	mux := http.NewServeMux()

	// Notification handlers
	nh := handlers.NewNotificationHandler(db)
	mux.HandleFunc("POST /notifications", nh.Create)
	mux.HandleFunc("GET /notifications", nh.List)
	mux.HandleFunc("PATCH /notifications/{id}/read", nh.Patch)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	// this keeps main thread alive and allows the server to run in the background
	go func() {
		log.Println("starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server forced to close: %v\n", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdown
	log.Printf("Received signal: %v. Initiating graceful shutdown...\n", sig)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v\n", err)
	}

	log.Println("Server exited cleanly.")
}
