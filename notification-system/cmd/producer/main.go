package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/abhinayjangde/notification-system/internal/kafka"
)

func main() {
	producer := kafka.NewProducer(
		"localhost:9092",
		"notifications",
	)

	defer producer.Close()

	event := kafka.NotificationEvent{
		EventID:     "evt-002",
		EventType:   "ORDER_SHIPPED",
		RecipientID: "user-123",
		Data: json.RawMessage(`{
			"order_id": "order-457",
			"tracking_id": "TRACK127"
		}`),
		CreatedAt: time.Now().Local().String(),
	}

	ctx := context.Background()

	if err := producer.Publish(ctx, event); err != nil {
		log.Fatal(err)
	}

	log.Println("event published successfully")
}
