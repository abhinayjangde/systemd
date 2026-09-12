package models

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID          string          `json:"id"`
	PecipientID string          `json:"recipient_id"`
	Type        string          `json:"type"`
	Title       string          `json:"title"`
	Body        string          `json:"body"`
	Data        json.RawMessage `json:"data"`
	ReadAt      *time.Time      `json:"read_at"`
	CreatedAt   time.Time       `jsong:"created_at"`
}
