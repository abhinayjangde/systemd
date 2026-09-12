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

type NotificationListResponse struct {
	Notifications []NotificationDTO `json:"notifications"`
	NextCursor    string            `json:"next_cursor,omitempty"`
}

type NotificationDTO struct {
	ID          string          `json:"id"`
	RecipientID string          `json:"recipient_id"`
	Type        string          `json:"type"`
	Title       string          `json:"title"`
	Body        string          `json:"body"`
	Data        json.RawMessage `json:"data"`
	ReadAt      *time.Time      `json:"read_at"`
	CreatedAt   time.Time       `json:"created_at"`
}
