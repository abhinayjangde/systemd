package handlers

import (
	"encoding/json"
	"time"
)

type NotificationRequest struct {
	RecipientID string          `json:"recipient_id"` // uuid.UUID
	Type        string          `json:"type"`
	Title       string          `json:"title"`
	Body        string          `json:"body"`
	Data        json.RawMessage `json:"data"`
}

type NotificationResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
