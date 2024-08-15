package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/health-analytics-service/health-analytics-service/config"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	consumer "github.com/health-analytics-service/health-analytics-service/kafka"
	mongodb "github.com/health-analytics-service/health-analytics-service/storage/mongo"
	"github.com/health-analytics-service/health-analytics-service/storage/redis"

	"github.com/health-analytics-service/health-analytics-service/service"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	// Initialize MongoDB storage
	mongoStorage, err := mongodb.NewMongoStorage(cfg)
	if err != nil {
		log.Fatalf("failed to initialize MongoDB storage: %v", err)
	}
	redisClient, err := redis.Connect(&cfg)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	// Initialize Kafka consumers
	geneticDataConsumer := consumer.NewGeneticDataConsumer(cfg.KafkaBrokers, cfg.KafkaGeneticDataTopic, mongoStorage, redisClient)
	healthRecommendationConsumer := consumer.NewHealthRecommendationConsumer(cfg.KafkaBrokers, cfg.KafkaHealthRecommendationTopic, mongoStorage, redisClient)
	lifestyleDataConsumer := consumer.NewLifestyleDataConsumer(cfg.KafkaBrokers, cfg.KafkaLifestyleDataTopic, mongoStorage, redisClient)
	medicalRecordConsumer := consumer.NewMedicalRecordConsumer(cfg.KafkaBrokers, cfg.KafkaMedicalRecordTopic, mongoStorage, redisClient)
	wearableDataConsumer := consumer.NewWearableDataConsumer(cfg.KafkaBrokers, cfg.KafkaWearableDataTopic, mongoStorage)

	// Start consumers in separate goroutines
	go func() {
		if err := geneticDataConsumer.Consume(context.Background()); err != nil {
			log.Fatalf("genetic data consumer error: %v", err)
		}
	}()

	go func() {
		if err := healthRecommendationConsumer.Consume(context.Background()); err != nil {
			log.Fatalf("health recommendation consumer error: %v", err)
		}
	}()

	go func() {
		if err := lifestyleDataConsumer.Consume(context.Background()); err != nil {
			log.Fatalf("lifestyle data consumer error: %v", err)
		}
	}()

	go func() {
		if err := medicalRecordConsumer.Consume(context.Background()); err != nil {
			log.Fatalf("medical record consumer error: %v", err)
		}
	}()

	go func() {
		if err := wearableDataConsumer.Consume(context.Background()); err != nil {
			log.Fatalf("wearable data consumer error: %v", err)
		}
	}()

	// Initialize gRPC server
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register gRPC services
	health.RegisterGeneticDataServiceServer(s, service.NewGeneticDataService(mongoStorage))
	health.RegisterHealthRecommendationServiceServer(s, service.NewHealthRecommendationService(mongoStorage))
	health.RegisterLifestyleDataServiceServer(s, service.NewLifestyleDataService(mongoStorage))
	health.RegisterMedicalRecordServiceServer(s, service.NewMedicalRecordService(mongoStorage))
	health.RegisterWearableDataServiceServer(s, service.NewWearableDataService(mongoStorage))
	health.RegisterHealthMonitoringServiceServer(s, service.NewHealthMonitoringService(mongoStorage))

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
