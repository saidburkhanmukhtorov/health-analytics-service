package test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/health-analytics-service/health-analytics-service/config"
	"github.com/health-analytics-service/health-analytics-service/storage"
	mongodb "github.com/health-analytics-service/health-analytics-service/storage/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createMongoDBConnection(t *testing.T) *mongo.Database {
	cfg := config.Load()

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Unable to connect to MongoDB: %v", err)
	}

	// Ping the database to verify the connection
	if err := client.Ping(context.Background(), nil); err != nil {
		t.Fatalf("Unable to ping MongoDB: %v", err)
	}

	return client.Database(cfg.MongoDB)
}

// StorageM implements the storage.StorageI interface for MongoDB.
type StorageM struct {
	db                       *mongo.Database
	medicalRecordRepo        storage.MedicalRecordRepoI
	geneticDataRepo          storage.GeneticDataRepoI
	lifestyleDataRepo        storage.LifestyleDataRepoI
	wearableDataRepo         storage.WearableDataRepoI
	healthRecommendationRepo storage.HealthRecommendationRepoI
	healthMonitoringRepo     storage.HealthMonitoringRepoI
}

// NewMongoStorage creates a new MongoDB storage instance.
func NewMongoStorageTest(cfg config.Config) (storage.StorageI, error) {
	// Construct MongoDB connection URI
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		slog.Warn("Unable to connect to MongoDB:" + err.Error())
		return nil, err
	}

	// Ping the database to verify the connection
	if err := client.Ping(context.Background(), nil); err != nil {
		slog.Warn("Unable to ping MongoDB:" + err.Error())
		return nil, err
	}

	db := client.Database(cfg.MongoDB)

	return &StorageM{
		db:                       db,
		medicalRecordRepo:        mongodb.NewMedicalRecordRepo(db),
		geneticDataRepo:          mongodb.NewGeneticDataRepo(db),
		lifestyleDataRepo:        mongodb.NewLifestyleDataRepo(db),
		wearableDataRepo:         mongodb.NewWearableDataRepo(db),
		healthRecommendationRepo: mongodb.NewHealthRecommendationRepo(db),
		healthMonitoringRepo:     mongodb.NewHealthMonitoringRepo(db),
	}, nil
}

// MedicalRecord returns the MedicalRecordRepoI implementation for MongoDB.
func (s *StorageM) MedicalRecord() storage.MedicalRecordRepoI {
	return s.medicalRecordRepo
}

// GeneticData returns the GeneticDataRepoI implementation for MongoDB.
func (s *StorageM) GeneticData() storage.GeneticDataRepoI {
	return s.geneticDataRepo
}

// LifestyleData returns the LifestyleDataRepoI implementation for MongoDB.
func (s *StorageM) LifestyleData() storage.LifestyleDataRepoI {
	return s.lifestyleDataRepo
}

// WearableData returns the WearableDataRepoI implementation for MongoDB.
func (s *StorageM) WearableData() storage.WearableDataRepoI {
	return s.wearableDataRepo
}

// HealthRecommendation returns the HealthRecommendationRepoI implementation for MongoDB.
func (s *StorageM) HealthRecommendation() storage.HealthRecommendationRepoI {
	return s.healthRecommendationRepo
}

// HealthMonitoring returns the HealthMonitoringRepoI implementation for MongoDB.
func (s *StorageM) HealthMonitoring() storage.HealthMonitoringRepoI {
	return s.healthMonitoringRepo
}
