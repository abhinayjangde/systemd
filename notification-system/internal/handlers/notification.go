package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotificationHandler struct {
	db *sql.DB
}

func NewNotificationHandler(db *sql.DB) *NotificationHandler {
	return &NotificationHandler{
		db: db,
	}
}

func (nh *NotificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	var req NotificationRequest
	// json.Unmarshal([]byte(req))
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Decode error:", err.Error())
		http.Error(w, "invalid json data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(req); err != nil {
		fmt.Println("Encode error:", err.Error())
		http.Error(w, "invalid json res data", http.StatusInternalServerError)
		return
	}
}
