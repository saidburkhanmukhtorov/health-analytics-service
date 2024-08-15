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

// GeneticDataConsumer consumes Kafka messages related to genetic data.
type GeneticDataConsumer struct {
	reader  *kafka.Reader
	storage storage.StorageI
	redis   *redis.Client // Add Redis client
}

// NewGeneticDataConsumer creates a new GeneticDataConsumer instance.
func NewGeneticDataConsumer(kafkaBrokers []string, topic string, storage storage.StorageI, redis *redis.Client) *GeneticDataConsumer { // Add Redis client to constructor
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   topic,
		GroupID: "genetic-data-group", // Choose a suitable group ID
	})
	return &GeneticDataConsumer{reader: reader, storage: storage, redis: redis}
}

// Consume starts consuming messages from the Kafka topic.
func (c *GeneticDataConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("error fetching message: %w", err)
		}

		// Determine the message type based on the key
		switch string(msg.Key) {
		case "genetic_data.create":
			var createModel health.GeneticData
			if err := json.Unmarshal(msg.Value, &createModel); err != nil {
				log.Printf("error unmarshalling create genetic data message: %v", err)
				continue
			}
			if _, err := c.storage.GeneticData().CreateGeneticData(ctx, &createModel); err != nil {
				log.Printf("error creating genetic data: %v", err)
			}

			// Send notification for creation
			if err := c.redis.AddNotification(ctx, createModel.UserId, "Your genetic data has been created."); err != nil {
				log.Printf("failed to send notification: %v", err)
				continue
			}

		case "genetic_data.update":
			var updateModel health.GeneticData
			if err := json.Unmarshal(msg.Value, &updateModel); err != nil {
				log.Printf("error unmarshalling update genetic data message: %v", err)
				continue
			}
			if err := c.storage.GeneticData().UpdateGeneticData(ctx, &updateModel); err != nil {
				log.Printf("error updating genetic data: %v", err)
			}

			// Send notification for update
			if err := c.redis.AddNotification(ctx, updateModel.UserId, "Your genetic data has been updated."); err != nil {
				log.Printf("failed to send notification: %v", err)
				continue
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
