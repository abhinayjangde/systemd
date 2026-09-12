package models

import (
	"encoding/json"
	"time"
	"uuid"
)

type Notification struct {
	ID          uuid.UUID       `json:"id"`
	RecipientID string          `json:"recipient_id"`
	Type        string          `json:"type"`
	Title       string          `json:"title"`
	Body        string          `json:"body"`
	Data        json.RawMessage `json:"data"`
	ReadAt      *time.Time      `json:"read_at"`
	CreatedAt   time.Time       `json:"created_at"`
}
