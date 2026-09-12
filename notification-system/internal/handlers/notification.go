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
	ctx := r.Context()

	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Decode error:", err.Error())
		http.Error(w, "invalid json data", http.StatusBadRequest)
		return
	}
	// basic validation
	if req.RecipientID == "" ||
		req.Type == "" ||
		req.Title == "" ||
		req.Body == "" {

		http.Error(
			w,
			"missing required fields",
			http.StatusBadRequest,
		)
		return
	}

	// inserting in db
	var out NotificationResponse
	err := nh.db.QueryRowContext(ctx, `
		INSERT INTO notifications (recipient_id, type, title, body, data)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, req.RecipientID, req.Type, req.Title, req.Body, req.Data).Scan(&out.ID, &out.CreatedAt)

	if err != nil {
		fmt.Println("Insert notification error:", err.Error())
		http.Error(w, "inserting notifications error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// response
	if err := json.NewEncoder(w).Encode(out); err != nil {
		fmt.Println("Encode error:", err.Error())
		http.Error(w, "invalid json res data", http.StatusInternalServerError)
		return
	}
}
