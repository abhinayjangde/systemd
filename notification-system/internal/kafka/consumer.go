package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	db     *sql.DB
}

func NewConsumer(
	brokerURL string,
	topic string,
	groupID string,
	db *sql.DB,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerURL},
		Topic:   topic,
		GroupID: groupID,

		// Start consuming from the earliest
		// available message when this group
		// has no committed offset.
		StartOffset: kafka.FirstOffset,
	})

	return &Consumer{
		reader: reader,
		db:     db,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("read kafka message: %w", err)
		}

		var event NotificationEvent

		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf(
				"failed to unmarshal message: partition=%d offset=%d error=%v",
				message.Partition,
				message.Offset,
				err,
			)
			continue
		}

		if err := c.saveNotification(ctx, event); err != nil {
			return err
		}

		log.Printf(
			"notification saved: event_id=%s partition=%d offset=%d",
			event.EventID,
			message.Partition,
			message.Offset,
		)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func (c *Consumer) saveNotification(
	ctx context.Context,
	event NotificationEvent,
) error {
	_, err := c.db.ExecContext(ctx, `
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
		return fmt.Errorf("save notification: %w", err)
	}

	return nil
}
