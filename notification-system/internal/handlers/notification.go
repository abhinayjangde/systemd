package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abhinayjangde/notification-system/internal/events"
	"github.com/abhinayjangde/notification-system/internal/httpx"
	"github.com/google/uuid"
)

type EventPublisher interface {
	Publish(context.Context, events.NotificationEvent) error
}

type NotificationHandler struct {
	db        *sql.DB
	publisher EventPublisher
}

func NewNotificationHandler(db *sql.DB, publisher EventPublisher) *NotificationHandler {
	return &NotificationHandler{
		db:        db,
		publisher: publisher,
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

	event := events.NotificationEvent{
		EventID:     uuid.NewString(),
		EventType:   req.Type,
		RecipientID: req.RecipientID,
		Title:       req.Title,
		Body:        req.Body,
		Data:        req.Data,
		CreatedAt:   time.Now().UTC(),
	}

	if err := nh.publisher.Publish(ctx, event); err != nil {
		fmt.Println("publish notification event error:", err)
		http.Error(w, "publishing notification event error", http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, http.StatusAccepted, NotificationResponse{EventID: event.EventID})
}

func (nh *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	recipientID := r.URL.Query().Get("recipient_id")

	if recipientID == "" {
		http.Error(
			w,
			"recipient_id is required",
			http.StatusBadRequest,
		)
		return
	}

	limit := 20

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		if parsed > 100 {
			parsed = 100
		}
		limit = parsed
	}

	queryLimit := limit + 1

	cursorValue := r.URL.Query().Get("cursor")
	var cursor *notificationCursor

	if cursorValue != "" {
		parsed, err := decodeCursor(cursorValue)
		if err != nil {
			http.Error(w, "invalid cursor", http.StatusBadRequest)
			return
		}

		if parsed.ID == "" || parsed.CreatedAt.IsZero() {
			http.Error(w, "invalid cursor", http.StatusBadRequest)
			return
		}

		cursor = &parsed
	}

	var rows *sql.Rows
	var err error

	if cursor == nil {
		// First page
		rows, err = nh.db.QueryContext(ctx, `
			SELECT
				id,
				recipient_id,
				type,
				title,
				body,
				data,
				read_at,
				created_at
			FROM notifications
			WHERE recipient_id = $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2
		`, recipientID, queryLimit)

	} else {

		// Next page
		rows, err = nh.db.QueryContext(ctx, `
			SELECT
				id,
				recipient_id,
				type,
				title,
				body,
				data,
				read_at,
				created_at
			FROM notifications
			WHERE recipient_id = $1
			  AND (created_at, id) < ($2, $3)
			ORDER BY created_at DESC, id DESC
			LIMIT $4
		`,
			recipientID,
			cursor.CreatedAt,
			cursor.ID,
			queryLimit,
		)

	}

	if err != nil {
		fmt.Println("query notifications error:", err)
		http.Error(w, "querying notifications error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	// scan rows
	notifications := NotificationListResponse{
		Notifications: make([]NotificationDTO, 0, limit),
	}
	for rows.Next() {
		var n NotificationDTO

		if err := rows.Scan(
			&n.ID,
			&n.RecipientID,
			&n.Type,
			&n.Title,
			&n.Body,
			&n.Data,
			&n.ReadAt,
			&n.CreatedAt,
		); err != nil {

			fmt.Println("scan notification error:", err)
			http.Error(w, "scanning notifications error", http.StatusInternalServerError)
			return
		}

		notifications.Notifications = append(
			notifications.Notifications,
			n,
		)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("rows error:", err)
		http.Error(w, "reading notifications error", http.StatusInternalServerError)
		return
	}

	if len(notifications.Notifications) > limit {

		// Remove the extra record.
		notifications.Notifications =
			notifications.Notifications[:limit]

		// The last notification actually returned
		// becomes our next cursor.
		last := notifications.Notifications[len(notifications.Notifications)-1]

		nextCursor, err := encodeCursor(notificationCursor{
			ID:        last.ID,
			CreatedAt: last.CreatedAt,
		})

		if err != nil {
			fmt.Println("encode cursor error:", err)
			http.Error(w, "encoding cursor error", http.StatusInternalServerError)
			return
		}

		notifications.NextCursor = nextCursor
	}

	httpx.WriteJSON(w, http.StatusOK, notifications)

}

func (nh *NotificationHandler) Patch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	// only authenticated user can call this endpoint/can update notification read
	recipientID := r.URL.Query().Get("recipient_id")
	// And later, when we add authentication, we should remove recipient_id from the query parameter and get the user ID from the authenticated request context instead.

	if id == "" {
		http.Error(
			w,
			"notification id is required",
			http.StatusBadRequest,
		)
		return
	}

	if recipientID == "" {
		http.Error(
			w,
			"recipient_id is required",
			http.StatusBadRequest,
		)
		return
	}

	result, err := nh.db.ExecContext(ctx, `
		UPDATE notifications
		SET read_at = NOW()
		WHERE id = $1
		  AND recipient_id = $2
		  AND read_at IS NULL
	`, id, recipientID)

	if err != nil {
		fmt.Println("mark notification as read error:", err)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	if rowsAffected == 0 {
		http.Error(
			w,
			"notification not found",
			http.StatusNotFound,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
