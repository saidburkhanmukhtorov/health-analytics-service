package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	mongodb "github.com/health-analytics-service/health-analytics-service/storage/mongo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestHealthMonitoringRepo(t *testing.T) {
	db := createMongoDBConnection(t)
	healthMonitoringRepo := mongodb.NewHealthMonitoringRepo(db)

	// Create mock data for all services
	userID := uuid.NewString()
	recordDate := time.Now().Format("2006-01-02")
	recordedTimestamp := time.Now().Format(time.RFC3339)

	// Mock Medical Record
	medicalRecordRepo := mongodb.NewMedicalRecordRepo(db)
	_, err := medicalRecordRepo.CreateMedicalRecord(context.Background(), &health.MedicalRecord{
		UserId:      userID,
		RecordType:  "Test Record",
		RecordDate:  recordDate,
		Description: "This is a test medical record.",
		DoctorId:    uuid.NewString(),
		Attachments: []string{"attachment1.txt", "attachment2.pdf"},
	})
	assert.NoError(t, err, "Creating mock medical record failed")

	// Mock Genetic Data
	geneticDataRepo := mongodb.NewGeneticDataRepo(db)
	dataValue, err := anypb.New(&health.MedicalRecord{
		UserId:      uuid.NewString(),
		RecordType:  "Genetic Test",
		RecordDate:  recordDate,
		Description: "Sample genetic test data",
	})
	assert.NoError(t, err, "Failed to create Any proto message")
	_, err = geneticDataRepo.CreateGeneticData(context.Background(), &health.GeneticData{
		UserId:       userID,
		DataType:     "DNA Sequencing",
		DataValue:    dataValue,
		AnalysisDate: recordDate,
	})
	assert.NoError(t, err, "Creating mock genetic data failed")

	// Mock Lifestyle Data
	lifestyleDataRepo := mongodb.NewLifestyleDataRepo(db)
	dataValue, err = anypb.New(&health.SleepData{
		UserId:        uuid.NewString(),
		SleepDuration: int64(8 * time.Hour),
		SleepQuality:  "Good",
		RecordedDate:  recordDate,
	})
	assert.NoError(t, err, "Failed to create Any proto message")
	_, err = lifestyleDataRepo.CreateLifestyleData(context.Background(), &health.LifestyleData{
		UserId:       userID,
		DataType:     "Sleep",
		DataValue:    dataValue,
		RecordedDate: recordDate,
	})
	assert.NoError(t, err, "Creating mock lifestyle data failed")

	// Mock Wearable Data
	wearableDataRepo := mongodb.NewWearableDataRepo(db)
	dataValue, err = anypb.New(&health.HeartRateData{
		UserId:            uuid.NewString(),
		HeartRate:         80,
		RecordedTimestamp: recordedTimestamp,
	})
	assert.NoError(t, err, "Failed to create Any proto message")
	_, err = wearableDataRepo.CreateWearableData(context.Background(), &health.WearableData{
		UserId:            userID,
		DeviceType:        "Smartwatch",
		DataType:          "HeartRate",
		DataValue:         dataValue,
		RecordedTimestamp: recordedTimestamp,
	})
	assert.NoError(t, err, "Creating mock wearable data failed")

	// Mock Health Recommendation
	healthRecommendationRepo := mongodb.NewHealthRecommendationRepo(db)
	_, err = healthRecommendationRepo.CreateHealthRecommendation(context.Background(), &health.HealthRecommendation{
		UserId:             userID,
		RecommendationType: "Exercise",
		Description:        "Engage in at least 30 minutes of moderate-intensity exercise most days of the week.",
		Priority:           2,
	})
	assert.NoError(t, err, "Creating mock health recommendation failed")

	// Test GetDailySummary
	t.Run("GetDailySummary", func(t *testing.T) {
		req := &health.DailySummaryRequest{
			UserId: userID,
			Date:   recordDate,
		}
		summary, err := healthMonitoringRepo.GetDailySummary(context.Background(), req)
		assert.NoError(t, err, "GetDailySummary should not return an error")
		assert.NotNil(t, summary, "GetDailySummary response should not be nil")
	})

	// Test GetWeeklySummary
	t.Run("GetWeeklySummary", func(t *testing.T) {
		startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02") // Start date one week ago
		endDate := recordDate
		req := &health.WeeklySummaryRequest{
			UserId:    userID,
			StartDate: startDate,
			EndDate:   endDate,
		}
		summary, err := healthMonitoringRepo.GetWeeklySummary(context.Background(), req)
		assert.NoError(t, err, "GetWeeklySummary should not return an error")
		assert.NotNil(t, summary, "GetWeeklySummary response should not be nil")
	})
}
