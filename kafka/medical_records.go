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

// MedicalRecordConsumer consumes Kafka messages related to medical records.
type MedicalRecordConsumer struct {
	reader  *kafka.Reader
	storage storage.StorageI
	redis   *redis.Client // Add Redis client
}

// NewMedicalRecordConsumer creates a new MedicalRecordConsumer instance.
func NewMedicalRecordConsumer(kafkaBrokers []string, topic string, storage storage.StorageI, redis *redis.Client) *MedicalRecordConsumer { // Add Redis client to constructor
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   topic,
		GroupID: "medical-record-group", // Choose a suitable group ID
	})
	return &MedicalRecordConsumer{reader: reader, storage: storage, redis: redis}
}

// Consume starts consuming messages from the Kafka topic.
func (c *MedicalRecordConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("error fetching message: %w", err)
		}

		// Determine the message type based on the key
		switch string(msg.Key) {
		case "medical_record.create":
			var createModel health.MedicalRecord
			if err := json.Unmarshal(msg.Value, &createModel); err != nil {
				log.Printf("error unmarshalling create medical record message: %v", err)
				continue
			}
			if _, err := c.storage.MedicalRecord().CreateMedicalRecord(ctx, &createModel); err != nil {
				log.Printf("error creating medical record: %v", err)
			}

			// Send notification for creation
			if err := c.redis.AddNotification(ctx, createModel.UserId, "Your medical record has been created."); err != nil {
				log.Printf("failed to send notification: %v", err)
				// Handle error (e.g., log and continue, retry, etc.)
			}

		case "medical_record.update":
			var updateModel health.MedicalRecord
			if err := json.Unmarshal(msg.Value, &updateModel); err != nil {
				log.Printf("error unmarshalling update medical record message: %v", err)
				continue
			}
			if err := c.storage.MedicalRecord().UpdateMedicalRecord(ctx, &updateModel); err != nil {
				log.Printf("error updating medical record: %v", err)
			}

			// Send notification for update
			if err := c.redis.AddNotification(ctx, updateModel.UserId, "Your medical record has been updated."); err != nil {
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
