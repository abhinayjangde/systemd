package main

import (
	"net/http"
	"os"
	"time"

	"log"

	"github.com/abhinayjangde/notification-system/internal/config"
)

func main() {
	cfg := config.MustLoad()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /notifications", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`notification created`))
	})

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
