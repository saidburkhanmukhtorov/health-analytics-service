package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"github.com/health-analytics-service/health-analytics-service/storage"
	"github.com/health-analytics-service/health-analytics-service/storage/redis"
	"github.com/segmentio/kafka-go"
)

// LifestyleDataConsumer consumes Kafka messages related to lifestyle data.
type LifestyleDataConsumer struct {
	reader  *kafka.Reader
	storage storage.StorageI
	redis   *redis.Client // Add Redis client
}

// NewLifestyleDataConsumer creates a new LifestyleDataConsumer instance.
func NewLifestyleDataConsumer(kafkaBrokers []string, topic string, storage storage.StorageI, redis *redis.Client) *LifestyleDataConsumer { // Add Redis client to constructor
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   topic,
		GroupID: "lifestyle-data-group", // Choose a suitable group ID
	})
	return &LifestyleDataConsumer{reader: reader, storage: storage, redis: redis}
}

// Consume starts consuming messages from the Kafka topic.
func (c *LifestyleDataConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("error fetching message: %w", err)
		}

		// Determine the message type based on the key
		switch string(msg.Key) {
		case "lifestyle_data.create":
			var createModel health.LifestyleData
			if err := json.Unmarshal(msg.Value, &createModel); err != nil {
				log.Printf("error unmarshalling create lifestyle data message: %v", err)
				continue
			}
			if _, err := c.storage.LifestyleData().CreateLifestyleData(ctx, &createModel); err != nil {
				log.Printf("error creating lifestyle data: %v", err)
			}

			// Send notification for creation
			if err := c.redis.AddNotification(ctx, createModel.UserId, "Your lifestyle data has been recorded."); err != nil {
				log.Printf("failed to send notification: %v", err)
				// Handle error (e.g., log and continue, retry, etc.)
			}

		case "lifestyle_data.update":
			var updateModel health.LifestyleData
			if err := json.Unmarshal(msg.Value, &updateModel); err != nil {
				log.Printf("error unmarshalling update lifestyle data message: %v", err)
				continue
			}
			if err := c.storage.LifestyleData().UpdateLifestyleData(ctx, &updateModel); err != nil {
				log.Printf("error updating lifestyle data: %v", err)
			}

			// Send notification for update
			if err := c.redis.AddNotification(ctx, updateModel.UserId, "Your lifestyle data has been updated."); err != nil {
				log.Printf("failed to send notification: %v", err)
				// Handle error
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
