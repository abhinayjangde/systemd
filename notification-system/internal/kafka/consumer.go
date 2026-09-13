package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(
	brokerURL string,
	topic string,
	groupID string,
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

		log.Printf(
			"received event: event_id=%s event_type=%s recipient_id=%s partition=%d offset=%d",
			event.EventID,
			event.EventType,
			event.RecipientID,
			message.Partition,
			message.Offset,
		)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
