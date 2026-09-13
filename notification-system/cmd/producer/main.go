package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/abhinayjangde/notification-system/internal/kafka"
	"github.com/google/uuid"
)

func main() {
	producer := kafka.NewProducer(
		"localhost:9092",
		"notifications",
	)

	defer producer.Close()

	event := kafka.NotificationEvent{
		EventID:     uuid.NewString(),
		EventType:   "ORDER_SHIPPED",
		RecipientID: uuid.NewString(),
		Data: json.RawMessage(`{
			"order_id": "order-458",
			"tracking_id": "TRACK127"
		}`),
		CreatedAt: time.Now(),
	}

	ctx := context.Background()

	if err := producer.Publish(ctx, event); err != nil {
		log.Fatal(err)
	}

	log.Println("event published successfully")
}
