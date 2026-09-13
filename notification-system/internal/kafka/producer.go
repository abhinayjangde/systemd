package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type NotificationEvent struct {
	EventID     string          `json:"event_id"`
	EventType   string          `json:"event_type"`
	RecipientID string          `json:"recipient_id"`
	Data        json.RawMessage `json:"data"`
	CreatedAt   string          `json:"created_at"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokerURL string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) Publish(
	ctx context.Context,
	event NotificationEvent,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.RecipientID),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
