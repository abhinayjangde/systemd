package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/abhinayjangde/notification-system/internal/events"
	"github.com/segmentio/kafka-go"
)

type EventHandler interface {
	ProcessEvent(
		ctx context.Context,
		event NotificationEvent,
	) error
}

type Consumer struct {
	reader     *kafka.Reader
	dlqWriter  *kafka.Writer
	dlqTopic   string
	handler    EventHandler
	maxRetries int
}

func NewConsumer(
	brokerURL string,
	topic string,
	groupID string,
	dlqTopic string,
	maxRetries int,
	handler EventHandler,
) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerURL},
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})
	dlqWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    dlqTopic,
		Balancer: &kafka.Hash{},
	}
	if maxRetries < 0 {
		maxRetries = 0
	}

	return &Consumer{
		reader:     reader,
		dlqWriter:  dlqWriter,
		dlqTopic:   dlqTopic,
		handler:    handler,
		maxRetries: maxRetries,
	}
}
func (c *Consumer) Consume(ctx context.Context) error {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return fmt.Errorf("fetch kafka message: %w", err)
		}

		var event NotificationEvent

		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf(
				"invalid kafka message, sending to DLQ: partition=%d offset=%d error=%v",
				message.Partition,
				message.Offset,
				err,
			)
			if err := c.deadLetter(ctx, message, err, 0); err != nil {
				return err
			}
			continue
		}

		var processErr error
		attempts := 0
		for attempt := 0; ; attempt++ {
			attempts++
			processErr = c.handler.ProcessEvent(ctx, event)
			if processErr == nil {
				break
			}
			if events.IsPermanent(processErr) || attempt >= c.maxRetries {
				break
			}

			if err := waitForRetry(ctx, attempt); err != nil {
				return err
			}
		}

		if processErr != nil {
			if err := c.deadLetter(ctx, message, processErr, attempts); err != nil {
				return err
			}
			log.Printf(
				"event sent to DLQ: event_id=%s error=%v",
				event.EventID,
				processErr,
			)
			continue
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
	readerErr := c.reader.Close()
	writerErr := c.dlqWriter.Close()
	if readerErr != nil {
		return readerErr
	}
	return writerErr
}

type deadLetterMessage struct {
	OriginalValue []byte    `json:"original_value"`
	Error         string    `json:"error"`
	Attempts      int       `json:"attempts"`
	Partition     int       `json:"partition"`
	Offset        int64     `json:"offset"`
	FailedAt      time.Time `json:"failed_at"`
}

func (c *Consumer) deadLetter(
	ctx context.Context,
	message kafka.Message,
	reason error,
	attempts int,
) error {
	payload, err := json.Marshal(deadLetterMessage{
		OriginalValue: message.Value,
		Error:         reason.Error(),
		Attempts:      attempts,
		Partition:     message.Partition,
		Offset:        message.Offset,
		FailedAt:      time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("marshal dead-letter message: %w", err)
	}

	if err := c.dlqWriter.WriteMessages(ctx, kafka.Message{
		Key:   message.Key,
		Value: payload,
	}); err != nil {
		return fmt.Errorf("publish message to DLQ %q: %w", c.dlqTopic, err)
	}

	if err := c.reader.CommitMessages(ctx, message); err != nil {
		return fmt.Errorf("commit dead-lettered message: %w", err)
	}
	return nil
}

func waitForRetry(ctx context.Context, attempt int) error {
	delay := 100 * time.Millisecond * time.Duration(1<<min(attempt, 5))
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}
