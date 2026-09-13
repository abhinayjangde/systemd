package notification

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateFromEvent(
	ctx context.Context,
	event Event,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO notifications (
			event_id,
			recipient_id,
			type,
			title,
			body,
			data,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (event_id) DO NOTHING
	`,
		event.EventID,
		event.RecipientID,
		event.EventType,
		"Notification",
		"An event occurred",
		event.Data,
		event.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}

	return nil
}
