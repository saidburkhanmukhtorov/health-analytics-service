package storage

import (
	"context"

	"github.com/health-analytics-service/health-analytics-service/genproto/health"
)

// StorageI defines the interface for interacting with the MongoDB storage layer.
type StorageI interface {
	MedicalRecord() MedicalRecordRepoI
	GeneticData() GeneticDataRepoI
	LifestyleData() LifestyleDataRepoI
	WearableData() WearableDataRepoI
	HealthRecommendation() HealthRecommendationRepoI
	HealthMonitoring() HealthMonitoringRepoI
}

// MedicalRecordRepoI defines methods for interacting with medical records in MongoDB.
type MedicalRecordRepoI interface {
	CreateMedicalRecord(ctx context.Context, record *health.MedicalRecord) (string, error)
	GetMedicalRecord(ctx context.Context, id string) (*health.MedicalRecord, error)
	UpdateMedicalRecord(ctx context.Context, record *health.MedicalRecord) error
	DeleteMedicalRecord(ctx context.Context, id string) error
	ListMedicalRecords(ctx context.Context, req *health.ListMedicalRecordsRequest) ([]*health.MedicalRecord, error)
}

// GeneticDataRepoI defines methods for interacting with genetic data in MongoDB.
type GeneticDataRepoI interface {
	CreateGeneticData(ctx context.Context, data *health.GeneticData) (string, error)
	GetGeneticData(ctx context.Context, id string) (*health.GeneticData, error)
	UpdateGeneticData(ctx context.Context, data *health.GeneticData) error
	DeleteGeneticData(ctx context.Context, id string) error
	ListGeneticData(ctx context.Context, req *health.ListGeneticDataRequest) ([]*health.GeneticData, error)
}

// LifestyleDataRepoI defines methods for interacting with lifestyle data in MongoDB.
type LifestyleDataRepoI interface {
	CreateLifestyleData(ctx context.Context, data *health.LifestyleData) (string, error)
	GetLifestyleData(ctx context.Context, id string) (*health.LifestyleData, error)
	UpdateLifestyleData(ctx context.Context, data *health.LifestyleData) error
	DeleteLifestyleData(ctx context.Context, id string) error
	ListLifestyleData(ctx context.Context, req *health.ListLifestyleDataRequest) ([]*health.LifestyleData, error)
}

// WearableDataRepoI defines methods for interacting with wearable data in MongoDB.
type WearableDataRepoI interface {
	CreateWearableData(ctx context.Context, data *health.WearableData) (string, error)
	GetWearableData(ctx context.Context, id string) (*health.WearableData, error)
	UpdateWearableData(ctx context.Context, data *health.WearableData) error
	DeleteWearableData(ctx context.Context, id string) error
	ListWearableData(ctx context.Context, req *health.ListWearableDataRequest) ([]*health.WearableData, error)
}

// HealthRecommendationRepoI defines methods for interacting with health recommendations in MongoDB.
type HealthRecommendationRepoI interface {
	CreateHealthRecommendation(ctx context.Context, recommendation *health.HealthRecommendation) (string, error)
	GetHealthRecommendation(ctx context.Context, id string) (*health.HealthRecommendation, error)
	UpdateHealthRecommendation(ctx context.Context, recommendation *health.HealthRecommendation) error
	DeleteHealthRecommendation(ctx context.Context, id string) error
	ListHealthRecommendations(ctx context.Context, req *health.ListHealthRecommendationsRequest) ([]*health.HealthRecommendation, error)
}

type HealthMonitoringRepoI interface {
	GetDailySummary(ctx context.Context, req *health.DailySummaryRequest) (*health.SummaryResponse, error)
	GetWeeklySummary(ctx context.Context, req *health.WeeklySummaryRequest) (*health.SummaryResponse, error)
}
