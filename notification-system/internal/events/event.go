package events

import (
	"encoding/json"
	"time"
)

// NotificationEvent is the transport-independent event consumed by the
// notification service and published to Kafka.
type NotificationEvent struct {
	EventID     string          `json:"event_id"`
	EventType   string          `json:"event_type"`
	RecipientID string          `json:"recipient_id"`
	Data        json.RawMessage `json:"data"`
	CreatedAt   time.Time       `json:"created_at"`
}
