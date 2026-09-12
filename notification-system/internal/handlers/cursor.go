package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type notificationCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

func encodeCursor(c notificationCursor) (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("marshal cursor: %w", err)
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

func decodeCursor(value string) (notificationCursor, error) {
	data, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return notificationCursor{}, fmt.Errorf("decode cursor: %w", err)
	}

	var cursor notificationCursor

	if err := json.Unmarshal(data, &cursor); err != nil {
		return notificationCursor{}, fmt.Errorf("unmarshal cursor: %w", err)
	}

	return cursor, nil
}
