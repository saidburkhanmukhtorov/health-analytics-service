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

// HealthRecommendationConsumer consumes Kafka messages related to health recommendations.
type HealthRecommendationConsumer struct {
	reader  *kafka.Reader
	storage storage.StorageI
	redis   *redis.Client // Add Redis client
}

// NewHealthRecommendationConsumer creates a new HealthRecommendationConsumer instance.
func NewHealthRecommendationConsumer(kafkaBrokers []string, topic string, storage storage.StorageI, redis *redis.Client) *HealthRecommendationConsumer { // Add Redis client to constructor
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   topic,
		GroupID: "health-recommendation-group", // Choose a suitable group ID
	})
	return &HealthRecommendationConsumer{reader: reader, storage: storage, redis: redis}
}

// Consume starts consuming messages from the Kafka topic.
func (c *HealthRecommendationConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("error fetching message: %w", err)
		}

		// Determine the message type based on the key
		switch string(msg.Key) {
		case "health_recommendation.create":
			var createModel health.HealthRecommendation
			if err := json.Unmarshal(msg.Value, &createModel); err != nil {
				log.Printf("error unmarshalling create health recommendation message: %v", err)
				continue
			}
			if _, err := c.storage.HealthRecommendation().CreateHealthRecommendation(ctx, &createModel); err != nil {
				log.Printf("error creating health recommendation: %v", err)
			}

			// Send notification for creation
			if err := c.redis.AddNotification(ctx, createModel.UserId, "You have a new health recommendation."); err != nil {
				log.Printf("failed to send notification: %v", err)
				// Handle error (e.g., log and continue, retry, etc.)
			}

		case "health_recommendation.update":
			var updateModel health.HealthRecommendation
			if err := json.Unmarshal(msg.Value, &updateModel); err != nil {
				log.Printf("error unmarshalling update health recommendation message: %v", err)
				continue
			}
			if err := c.storage.HealthRecommendation().UpdateHealthRecommendation(ctx, &updateModel); err != nil {
				log.Printf("error updating health recommendation: %v", err)
			}

			// Send notification for update
			if err := c.redis.AddNotification(ctx, updateModel.UserId, "A health recommendation has been updated."); err != nil {
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
