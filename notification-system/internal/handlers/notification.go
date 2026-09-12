package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		fmt.Println("encode response error:", err)
	}

}
