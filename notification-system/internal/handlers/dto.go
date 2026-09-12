package handlers

import "encoding/json"

type NotificationRequest struct {
	PecipientID string          `json:"recipient_id"`
	Type        string          `json:"type"`
	Title       string          `json:"title"`
	Body        string          `json:"body"`
	Data        json.RawMessage `json:"data"`
}
