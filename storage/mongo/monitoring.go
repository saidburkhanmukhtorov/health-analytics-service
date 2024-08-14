package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// HealthMonitoringRepo implements the storage.HealthMonitoringRepoI interface for MongoDB.
type HealthMonitoringRepo struct {
	db *mongo.Database
}

// NewHealthMonitoringRepo creates a new HealthMonitoringRepo instance.
func NewHealthMonitoringRepo(db *mongo.Database) *HealthMonitoringRepo {
	return &HealthMonitoringRepo{
		db: db,
	}
}

// GetDailySummary retrieves a daily summary of health data for a given user ID and date.
func (r *HealthMonitoringRepo) GetDailySummary(ctx context.Context, req *health.DailySummaryRequest) (*health.SummaryResponse, error) {
	// Build the filter query based on the request parameters
	filter := bson.M{
		"user_id": req.UserId,
		"created_at": bson.M{
			"$gte": time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.UTC),
		},
	}

	// Retrieve medical records
	medicalRecords, err := r.getMedicalRecordsForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve medical records: %w", err)
	}

	// Retrieve genetic data
	geneticData, err := r.getGeneticDataForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve genetic data: %w", err)
	}

	// Retrieve lifestyle data
	lifestyleData, err := r.getLifestyleDataForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve lifestyle data: %w", err)
	}

	// Retrieve wearable data
	wearableData, err := r.getWearableDataForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wearable data: %w", err)
	}

	// Retrieve health recommendations
	healthRecommendations, err := r.getHealthRecommendationsForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve health recommendations: %w", err)
	}

	// Construct the summary response
	summaryResponse := &health.SummaryResponse{
		MedicalRecords:        medicalRecords,
		GeneticData:           geneticData,
		LifestyleData:         lifestyleData,
		WearableData:          wearableData,
		HealthRecommendations: healthRecommendations,
	}

	return summaryResponse, nil
}

// GetWeeklySummary retrieves a weekly summary of health data for a given user ID and date range.
func (r *HealthMonitoringRepo) GetWeeklySummary(ctx context.Context, req *health.WeeklySummaryRequest) (*health.SummaryResponse, error) {
	// Parse start and end dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Build the filter query based on the request parameters
	filter := bson.M{
		"user_id": req.UserId,
		"created_at": bson.M{
			"$gte": startDate,
			"$lt":  endDate.AddDate(0, 0, 1), // Add 1 day to include the end date
		},
	}

	// Retrieve medical records
	medicalRecords, err := r.getMedicalRecordsForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve medical records: %w", err)
	}

	// Retrieve genetic data
	geneticData, err := r.getGeneticDataForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve genetic data: %w", err)
	}

	// Retrieve lifestyle data
	lifestyleData, err := r.getLifestyleDataForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve lifestyle data: %w", err)
	}

	// Retrieve wearable data
	wearableData, err := r.getWearableDataForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wearable data: %w", err)
	}

	// Retrieve health recommendations
	healthRecommendations, err := r.getHealthRecommendationsForSummary(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve health recommendations: %w", err)
	}

	// Construct the summary response
	summaryResponse := &health.SummaryResponse{
		MedicalRecords:        medicalRecords,
		GeneticData:           geneticData,
		LifestyleData:         lifestyleData,
		WearableData:          wearableData,
		HealthRecommendations: healthRecommendations,
	}

	return summaryResponse, nil
}

// Helper functions to retrieve data for summaries

func (r *HealthMonitoringRepo) getMedicalRecordsForSummary(ctx context.Context, filter bson.M) ([]*health.MedicalRecord, error) {
	cursor, err := r.db.Collection("medical_records").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find medical records: %w", err)
	}
	defer cursor.Close(ctx)

	var medicalRecords []*health.MedicalRecord
	for cursor.Next(ctx) {
		var bsonData bson.M
		if err := cursor.Decode(&bsonData); err != nil {
			return nil, fmt.Errorf("failed to decode medical record: %w", err)
		}

		recordModel, err := bsonToMedicalRecord(bsonData)
		if err != nil {
			return nil, err
		}

		medicalRecords = append(medicalRecords, recordModel)
	}

	return medicalRecords, nil
}

func (r *HealthMonitoringRepo) getGeneticDataForSummary(ctx context.Context, filter bson.M) ([]*health.GeneticData, error) {
	cursor, err := r.db.Collection("genetic_data").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find genetic data: %w", err)
	}
	defer cursor.Close(ctx)

	var geneticDataRecords []*health.GeneticData
	for cursor.Next(ctx) {
		var bsonData bson.M
		if err := cursor.Decode(&bsonData); err != nil {
			return nil, fmt.Errorf("failed to decode genetic data: %w", err)
		}

		dataModel, err := bsonToGeneticData(bsonData)
		if err != nil {
			return nil, err
		}

		geneticDataRecords = append(geneticDataRecords, dataModel)
	}

	return geneticDataRecords, nil
}

func (r *HealthMonitoringRepo) getLifestyleDataForSummary(ctx context.Context, filter bson.M) ([]*health.LifestyleData, error) {
	cursor, err := r.db.Collection("lifestyle_data").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find lifestyle data: %w", err)
	}
	defer cursor.Close(ctx)

	var lifestyleDataRecords []*health.LifestyleData
	for cursor.Next(ctx) {
		var bsonData bson.M
		if err := cursor.Decode(&bsonData); err != nil {
			return nil, fmt.Errorf("failed to decode lifestyle data: %w", err)
		}

		dataModel, err := bsonToLifestyleData(bsonData)
		if err != nil {
			return nil, err
		}

		lifestyleDataRecords = append(lifestyleDataRecords, dataModel)
	}

	return lifestyleDataRecords, nil
}

func (r *HealthMonitoringRepo) getWearableDataForSummary(ctx context.Context, filter bson.M) ([]*health.WearableData, error) {
	cursor, err := r.db.Collection("wearable_data").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find wearable data: %w", err)
	}
	defer cursor.Close(ctx)

	var wearableDataRecords []*health.WearableData
	for cursor.Next(ctx) {
		var bsonData bson.M
		if err := cursor.Decode(&bsonData); err != nil {
			return nil, fmt.Errorf("failed to decode wearable data: %w", err)
		}

		dataModel, err := bsonToWearableData(bsonData)
		if err != nil {
			return nil, err
		}

		wearableDataRecords = append(wearableDataRecords, dataModel)
	}

	return wearableDataRecords, nil
}

func (r *HealthMonitoringRepo) getHealthRecommendationsForSummary(ctx context.Context, filter bson.M) ([]*health.HealthRecommendation, error) {
	cursor, err := r.db.Collection("health_recommendations").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find health recommendations: %w", err)
	}
	defer cursor.Close(ctx)

	var healthRecommendations []*health.HealthRecommendation
	for cursor.Next(ctx) {
		var bsonData bson.M
		if err := cursor.Decode(&bsonData); err != nil {
			return nil, fmt.Errorf("failed to decode health recommendation: %w", err)
		}

		recommendationModel, err := bsonToHealthRecommendation(bsonData)
		if err != nil {
			return nil, err
		}

		healthRecommendations = append(healthRecommendations, recommendationModel)
	}

	return healthRecommendations, nil
}
