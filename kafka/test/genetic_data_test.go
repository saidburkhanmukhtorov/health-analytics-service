package test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/health-analytics-service/health-analytics-service/config"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	consumer "github.com/health-analytics-service/health-analytics-service/kafka"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/health-analytics-service/health-analytics-service/storage/test"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestGeneticDataConsumer(t *testing.T) {
	cfg := config.Load()

	// Create a test topic
	topic := "test-genetic-data-topic"
	createTopic(t, cfg.KafkaBrokersTest, topic)
	// defer deleteTopic(t, cfg.KafkaBrokersTest, topic)

	// Initialize MongoDB storage for testing
	storage, err := test.NewMongoStorageTest(cfg)
	if err != nil {
		t.Fatalf("failed to initialize storage: %v", err)
	}

	// Create a test genetic data model
	dataValue, err := anypb.New(&health.MedicalRecord{
		UserId:      uuid.NewString(),
		RecordType:  "Genetic Test",
		RecordDate:  time.Now().Format("2006-01-02"),
		Description: "Sample genetic test data",
	})
	assert.NoError(t, err, "Failed to create Any proto message")

	geneticDataModel := &health.GeneticData{
		Id:           primitive.NewObjectID().Hex(),
		UserId:       uuid.NewString(),
		DataType:     "DNA Sequencing",
		DataValue:    dataValue,
		AnalysisDate: time.Now().Format("2006-01-02"),
	}
	// Create a GeneticDataConsumer with the test storage
	consumer := consumer.NewGeneticDataConsumer(cfg.KafkaBrokersTest, topic, storage)
	// Consume the message
	go func() {
		if err := consumer.Consume(context.Background()); err != nil {
			t.Errorf("Error consuming message: %v", err)
		}
	}()
	time.Sleep(2 * time.Second)
	// Produce a message to the Kafka topic
	produceMessage(t, cfg.KafkaBrokersTest, topic, "genetic_data.create", geneticDataModel)

	// Wait for the message to be consumed (adjust timeout as needed)
	time.Sleep(time.Second * 4)

	// Retrieve the created genetic data from the database
	createdData, err := storage.GeneticData().GetGeneticData(context.Background(), geneticDataModel.Id)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, geneticDataModel.UserId, createdData.UserId)
	assert.Equal(t, geneticDataModel.DataType, createdData.DataType)
	assert.Equal(t, geneticDataModel.AnalysisDate, createdData.AnalysisDate)
	assert.Equal(t, geneticDataModel.DataValue.String(), createdData.DataValue.String())
}

// Helper functions to create, delete, and produce messages to a Kafka topic
func createTopic(t *testing.T, brokers []string, topic string) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokers[0], topic, 0)
	if err != nil {
		t.Fatalf("failed to dial leader: %v", err)
	}
	defer conn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}
}

func deleteTopic(t *testing.T, brokers []string, topic string) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokers[0], topic, 0)
	if err != nil {
		t.Fatalf("failed to dial leader: %v", err)
	}
	defer conn.Close()

	err = conn.DeleteTopics(topic)
	if err != nil {
		t.Fatalf("failed to delete topic: %v", err)
	}
}

func produceMessage(t *testing.T, brokers []string, topic string, key string, message interface{}) {
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}

	value, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		t.Fatalf("failed to write messages: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}
}
