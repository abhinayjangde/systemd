package main

import (
	"net/http"
	"os"
	"time"

	"log"

	"github.com/abhinayjangde/notification-system/internal/config"
	"github.com/abhinayjangde/notification-system/internal/db"
	"github.com/abhinayjangde/notification-system/internal/handlers"
)

func main() {
	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	// Notification handlers
	nh := handlers.NewNotificationHandler(db)
	mux.HandleFunc("POST /notifications", nh.Create)

	log.Println("starting server", "port", cfg.Port)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("server stopped", "err", err)
		os.Exit(1)
	}
}
