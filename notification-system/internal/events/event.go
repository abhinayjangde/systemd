package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// NotificationEvent is the transport-independent event consumed by the
// notification service and published to Kafka.
type NotificationEvent struct {
	EventID     string          `json:"event_id"`
	EventType   string          `json:"event_type"`
	RecipientID string          `json:"recipient_id"`
	Title       string          `json:"title"`
	Body        string          `json:"body"`
	Data        json.RawMessage `json:"data"`
	CreatedAt   time.Time       `json:"created_at"`
}

type PermanentError struct {
	Err error
}

func (e *PermanentError) Error() string {
	return e.Err.Error()
}

func (e *PermanentError) Unwrap() error {
	return e.Err
}

func NewPermanentError(format string, args ...any) error {
	return &PermanentError{Err: fmt.Errorf(format, args...)}
}

func IsPermanent(err error) bool {
	var permanent *PermanentError
	return errors.As(err, &permanent)
}
