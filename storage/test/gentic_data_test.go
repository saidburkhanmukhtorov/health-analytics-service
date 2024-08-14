package test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	mongodb "github.com/health-analytics-service/health-analytics-service/storage/mongo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestGeneticDataRepo(t *testing.T) {
	db := createMongoDBConnection(t)
	geneticDataRepo := mongodb.NewGeneticDataRepo(db)

	t.Run("CreateGeneticData", func(t *testing.T) {
		// Create a sample Any proto message
		dataValue, err := anypb.New(&health.MedicalRecord{
			UserId:      uuid.NewString(),
			RecordType:  "Genetic Test",
			RecordDate:  time.Now().Format("2006-01-02"),
			Description: "Sample genetic test data",
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testGeneticData := &health.GeneticData{
			UserId:       uuid.NewString(),
			DataType:     "DNA Sequencing",
			DataValue:    dataValue,
			AnalysisDate: time.Now().Format("2006-01-02"),
		}
		createdID, err := geneticDataRepo.CreateGeneticData(context.Background(), testGeneticData)

		assert.NoError(t, err, "CreateGeneticData should not return an error")
		assert.NotEmpty(t, createdID, "Created genetic data should have a valid ID")
	})

	t.Run("GetGeneticData", func(t *testing.T) {
		// 1. Create a record to retrieve
		dataValue, err := anypb.New(&health.MedicalRecord{
			UserId:      uuid.NewString(),
			RecordType:  "Genetic Test",
			RecordDate:  time.Now().Format("2006-01-02"),
			Description: "Sample genetic test data",
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testGeneticData := &health.GeneticData{
			UserId:       uuid.NewString(),
			DataType:     "DNA Sequencing",
			DataValue:    dataValue,
			AnalysisDate: time.Now().Format("2006-01-02"),
		}
		createdID, err := geneticDataRepo.CreateGeneticData(context.Background(), testGeneticData)
		assert.NoError(t, err, "Creating genetic data for GetGeneticData test failed")
		assert.NotEmpty(t, createdID, "Created genetic data should have a valid ID")

		// 2. Get the record
		retrievedData, err := geneticDataRepo.GetGeneticData(context.Background(), createdID)

		assert.NoError(t, err, "GetGeneticData should not return an error")
		assert.NotNil(t, retrievedData, "GetGeneticData response should not be nil")
		assert.Equal(t, testGeneticData.UserId, retrievedData.UserId)
		assert.Equal(t, testGeneticData.DataType, retrievedData.DataType)
		assert.Equal(t, testGeneticData.AnalysisDate, retrievedData.AnalysisDate)
		// Compare the Any proto messages
		assert.Equal(t, testGeneticData.DataValue.String(), retrievedData.DataValue.String())
		log.Println(retrievedData.DataValue)
	})

	t.Run("UpdateGeneticData", func(t *testing.T) {
		// 1. Create a record to update
		dataValue, err := anypb.New(&health.MedicalRecord{
			UserId:      uuid.NewString(),
			RecordType:  "Genetic Test",
			RecordDate:  time.Now().Format("2006-01-02"),
			Description: "Sample genetic test data",
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testGeneticData := &health.GeneticData{
			UserId:       uuid.NewString(),
			DataType:     "DNA Sequencing",
			DataValue:    dataValue,
			AnalysisDate: time.Now().Format("2006-01-02"),
		}
		createdID, err := geneticDataRepo.CreateGeneticData(context.Background(), testGeneticData)
		assert.NoError(t, err, "Creating genetic data for UpdateGeneticData test failed")
		assert.NotEmpty(t, createdID, "Created genetic data should have a valid ID")

		// 2. Update the record
		updatedDataValue, err := anypb.New(&health.MedicalRecord{
			UserId:      uuid.NewString(),
			RecordType:  "Updated Genetic Test",
			RecordDate:  time.Now().Format("2006-01-02"),
			Description: "Updated sample genetic test data",
		})
		assert.NoError(t, err, "Failed to create updated Any proto message")

		updateRecord := &health.GeneticData{
			Id:           createdID,
			UserId:       uuid.NewString(),
			DataType:     "Updated DNA Sequencing",
			DataValue:    updatedDataValue,
			AnalysisDate: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
		}
		err = geneticDataRepo.UpdateGeneticData(context.Background(), updateRecord)
		assert.NoError(t, err, "UpdateGeneticData should not return an error")

		// 3. Retrieve the record and verify the update
		retrievedRecord, err := geneticDataRepo.GetGeneticData(context.Background(), createdID)
		assert.NoError(t, err, "GetGeneticData after update should not return an error")
		assert.Equal(t, updateRecord.UserId, retrievedRecord.UserId, "UserId should be updated")
		assert.Equal(t, updateRecord.DataType, retrievedRecord.DataType, "DataType should be updated")
		assert.Equal(t, updateRecord.AnalysisDate, retrievedRecord.AnalysisDate, "AnalysisDate should be updated")
		// Compare the Any proto messages
		assert.Equal(t, updateRecord.DataValue.String(), retrievedRecord.DataValue.String())
	})

	t.Run("DeleteGeneticData", func(t *testing.T) {
		// 1. Create a record to delete
		dataValue, err := anypb.New(&health.MedicalRecord{
			UserId:      uuid.NewString(),
			RecordType:  "Genetic Test",
			RecordDate:  time.Now().Format("2006-01-02"),
			Description: "Sample genetic test data",
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testGeneticData := &health.GeneticData{
			UserId:       uuid.NewString(),
			DataType:     "DNA Sequencing",
			DataValue:    dataValue,
			AnalysisDate: time.Now().Format("2006-01-02"),
		}
		createdID, err := geneticDataRepo.CreateGeneticData(context.Background(), testGeneticData)
		assert.NoError(t, err, "Creating genetic data for DeleteGeneticData test failed")
		assert.NotEmpty(t, createdID, "Created genetic data should have a valid ID")

		// 2. Delete the record
		err = geneticDataRepo.DeleteGeneticData(context.Background(), createdID)
		assert.NoError(t, err, "DeleteGeneticData should not return an error")

		// 3. Attempt to retrieve the deleted record (should fail)
		retrievedRecord, err := geneticDataRepo.GetGeneticData(context.Background(), createdID)
		assert.Error(t, err, "GetGeneticData after delete should return an error")
		assert.Nil(t, retrievedRecord, "GetGeneticData response should be nil after delete")
	})

	t.Run("ListGeneticData", func(t *testing.T) {
		// 1. Create some records for a specific user
		userID := uuid.NewString()
		testRecords := []*health.GeneticData{
			{
				UserId:       userID,
				DataType:     "DNA Sequencing 1",
				DataValue:    createSampleAny(t),
				AnalysisDate: time.Now().Format("2006-01-02"),
			},
			{
				UserId:       userID,
				DataType:     "DNA Sequencing 2",
				DataValue:    createSampleAny(t),
				AnalysisDate: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
			},
		}

		for _, record := range testRecords {
			_, err := geneticDataRepo.CreateGeneticData(context.Background(), record)
			assert.NoError(t, err, "Creating genetic data for ListGeneticData test failed")
		}

		// 2. Test listing all records for the user
		req := &health.ListGeneticDataRequest{
			UserId: userID,
		}
		retrievedRecords, err := geneticDataRepo.ListGeneticData(context.Background(), req)
		assert.NoError(t, err, "ListGeneticData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListGeneticData response should not be nil")
		assert.GreaterOrEqual(t, len(retrievedRecords), 2, "Should have at least two genetic data records for the user")

		// 3. Test filtering by DataType
		req = &health.ListGeneticDataRequest{
			UserId:   userID,
			DataType: "DNA Sequencing 1",
		}
		retrievedRecords, err = geneticDataRepo.ListGeneticData(context.Background(), req)
		assert.NoError(t, err, "ListGeneticData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListGeneticData response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one genetic data record matching the filter")
		assert.Equal(t, "DNA Sequencing 1", retrievedRecords[0].DataType, "DataType should match the filter")

		// 4. Test filtering by AnalysisDate
		req = &health.ListGeneticDataRequest{
			UserId:       userID,
			AnalysisDate: time.Now().Format("2006-01-02"),
		}
		retrievedRecords, err = geneticDataRepo.ListGeneticData(context.Background(), req)
		assert.NoError(t, err, "ListGeneticData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListGeneticData response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one genetic data record matching the filter")
		assert.Equal(t, time.Now().Format("2006-01-02"), retrievedRecords[0].AnalysisDate, "AnalysisDate should match the filter")
	})
}

func createSampleAny(t *testing.T) *anypb.Any {
	dataValue, err := anypb.New(&health.MedicalRecord{
		UserId:      uuid.NewString(),
		RecordType:  "Genetic Test",
		RecordDate:  time.Now().Format("2006-01-02"),
		Description: "Sample genetic test data",
	})
	assert.NoError(t, err, "Failed to create Any proto message")
	return dataValue
}
