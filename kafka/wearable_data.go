package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"github.com/health-analytics-service/health-analytics-service/storage"
	"github.com/segmentio/kafka-go"
)

// WearableDataConsumer consumes Kafka messages related to wearable data.
type WearableDataConsumer struct {
	reader  *kafka.Reader
	storage storage.StorageI
}

// NewWearableDataConsumer creates a new WearableDataConsumer instance.
func NewWearableDataConsumer(kafkaBrokers []string, topic string, storage storage.StorageI) *WearableDataConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   topic,
		GroupID: "wearable-data-group", // Choose a suitable group ID
	})
	return &WearableDataConsumer{reader: reader, storage: storage}
}

// Consume starts consuming messages from the Kafka topic.
func (c *WearableDataConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("error fetching message: %w", err)
		}

		// Determine the message type based on the key
		switch string(msg.Key) {
		case "wearable_data.create":
			var createModel health.WearableData
			if err := json.Unmarshal(msg.Value, &createModel); err != nil {
				log.Printf("error unmarshalling create wearable data message: %v", err)
				continue
			}
			if _, err := c.storage.WearableData().CreateWearableData(ctx, &createModel); err != nil {
				log.Printf("error creating wearable data: %v", err)
			}

		case "wearable_data.update":
			var updateModel health.WearableData
			if err := json.Unmarshal(msg.Value, &updateModel); err != nil {
				log.Printf("error unmarshalling update wearable data message: %v", err)
				continue
			}
			if err := c.storage.WearableData().UpdateWearableData(ctx, &updateModel); err != nil {
				log.Printf("error updating wearable data: %v", err)
			}
		default:
			log.Printf("unknown message key: %s", msg.Key)
		}

		// Commit the message
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("error committing message: %w", err)
		}
	}
}
