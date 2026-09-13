package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type EventHandler interface {
	ProcessEvent(
		ctx context.Context,
		event NotificationEvent,
	) error
}

type Consumer struct {
	reader  *kafka.Reader
	handler EventHandler
}

func NewConsumer(
	brokerURL string,
	topic string,
	groupID string,
	handler EventHandler,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerURL},
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}
func (c *Consumer) Consume(ctx context.Context) error {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("fetch kafka message: %w", err)
		}

		var event NotificationEvent

		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf(
				"invalid kafka message: partition=%d offset=%d error=%v",
				message.Partition,
				message.Offset,
				err,
			)
			continue
		}

		if err := c.handler.ProcessEvent(ctx, event); err != nil {
			return fmt.Errorf(
				"process event %s: %w",
				event.EventID,
				err,
			)
		}

		if err := c.reader.CommitMessages(ctx, message); err != nil {
			return fmt.Errorf(
				"commit kafka message: %w",
				err,
			)
		}

		log.Printf(
			"event processed: event_id=%s partition=%d offset=%d",
			event.EventID,
			message.Partition,
			message.Offset,
		)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

// func (c *Consumer) saveNotification(
// 	ctx context.Context,
// 	event NotificationEvent,
// ) error {
// 	_, err := c.handler.ExecContext(ctx, `
// 		INSERT INTO notifications (
// 			event_id,
// 			recipient_id,
// 			type,
// 			title,
// 			body,
// 			data,
// 			created_at
// 		)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7)
// 		ON CONFLICT (event_id) DO NOTHING
// 	`,
// 		event.EventID,
// 		event.RecipientID,
// 		event.EventType,
// 		"Notification",
// 		"An event occurred",
// 		event.Data,
// 		event.CreatedAt,
// 	)

// 	if err != nil {
// 		return fmt.Errorf("save notification: %w", err)
// 	}

// 	return nil
// }
