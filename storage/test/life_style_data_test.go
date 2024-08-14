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

func TestLifestyleDataRepo(t *testing.T) {
	db := createMongoDBConnection(t)
	lifestyleDataRepo := mongodb.NewLifestyleDataRepo(db)

	t.Run("CreateLifestyleData", func(t *testing.T) {
		// Create a sample Any proto message
		dataValue, err := anypb.New(&health.SleepData{
			UserId:        uuid.NewString(),
			SleepDuration: int64(8 * time.Hour),
			SleepQuality:  "Good",
			RecordedDate:  time.Now().Format("2006-01-02"),
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testLifestyleData := &health.LifestyleData{
			UserId:       uuid.NewString(),
			DataType:     "Sleep",
			DataValue:    dataValue,
			RecordedDate: time.Now().Format("2006-01-02"),
		}
		createdID, err := lifestyleDataRepo.CreateLifestyleData(context.Background(), testLifestyleData)

		assert.NoError(t, err, "CreateLifestyleData should not return an error")
		assert.NotEmpty(t, createdID, "Created lifestyle data should have a valid ID")
	})

	t.Run("GetLifestyleData", func(t *testing.T) {
		// 1. Create a record to retrieve
		dataValue, err := anypb.New(&health.SleepData{
			UserId:        uuid.NewString(),
			SleepDuration: int64(8 * time.Hour),
			SleepQuality:  "Good",
			RecordedDate:  time.Now().Format("2006-01-02"),
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testLifestyleData := &health.LifestyleData{
			UserId:       uuid.NewString(),
			DataType:     "Sleep",
			DataValue:    dataValue,
			RecordedDate: time.Now().Format("2006-01-02"),
		}
		createdID, err := lifestyleDataRepo.CreateLifestyleData(context.Background(), testLifestyleData)
		assert.NoError(t, err, "Creating lifestyle data for GetLifestyleData test failed")
		assert.NotEmpty(t, createdID, "Created lifestyle data should have a valid ID")

		// 2. Get the record
		retrievedData, err := lifestyleDataRepo.GetLifestyleData(context.Background(), createdID)

		assert.NoError(t, err, "GetLifestyleData should not return an error")
		assert.NotNil(t, retrievedData, "GetLifestyleData response should not be nil")
		assert.Equal(t, testLifestyleData.UserId, retrievedData.UserId)
		assert.Equal(t, testLifestyleData.DataType, retrievedData.DataType)
		assert.Equal(t, testLifestyleData.RecordedDate, retrievedData.RecordedDate)
		// Compare the Any proto messages
		assert.Equal(t, testLifestyleData.DataValue.String(), retrievedData.DataValue.String())
	})

	t.Run("UpdateLifestyleData", func(t *testing.T) {
		// 1. Create a record to update
		dataValue, err := anypb.New(&health.SleepData{
			UserId:        uuid.NewString(),
			SleepDuration: int64(8 * time.Hour),
			SleepQuality:  "Good",
			RecordedDate:  time.Now().Format("2006-01-02"),
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testLifestyleData := &health.LifestyleData{
			UserId:       uuid.NewString(),
			DataType:     "Sleep",
			DataValue:    dataValue,
			RecordedDate: time.Now().Format("2006-01-02"),
		}
		createdID, err := lifestyleDataRepo.CreateLifestyleData(context.Background(), testLifestyleData)
		assert.NoError(t, err, "Creating lifestyle data for UpdateLifestyleData test failed")
		assert.NotEmpty(t, createdID, "Created lifestyle data should have a valid ID")

		// 2. Update the record
		updatedDataValue, err := anypb.New(&health.SleepData{
			UserId:        uuid.NewString(),
			SleepDuration: int64(8 * time.Hour),
			SleepQuality:  "Average",
			RecordedDate:  time.Now().Add(time.Hour * 24).Format("2006-01-02"),
		})
		assert.NoError(t, err, "Failed to create updated Any proto message")

		updateRecord := &health.LifestyleData{
			Id:           createdID,
			UserId:       uuid.NewString(),
			DataType:     "Updated Sleep",
			DataValue:    updatedDataValue,
			RecordedDate: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
		}
		err = lifestyleDataRepo.UpdateLifestyleData(context.Background(), updateRecord)
		assert.NoError(t, err, "UpdateLifestyleData should not return an error")

		// 3. Retrieve the record and verify the update
		retrievedRecord, err := lifestyleDataRepo.GetLifestyleData(context.Background(), createdID)
		assert.NoError(t, err, "GetLifestyleData after update should not return an error")
		assert.Equal(t, updateRecord.UserId, retrievedRecord.UserId, "UserId should be updated")
		assert.Equal(t, updateRecord.DataType, retrievedRecord.DataType, "DataType should be updated")
		assert.Equal(t, updateRecord.RecordedDate, retrievedRecord.RecordedDate, "RecordedDate should be updated")
		// Compare the Any proto messages
		assert.Equal(t, updateRecord.DataValue.String(), retrievedRecord.DataValue.String())
	})

	t.Run("DeleteLifestyleData", func(t *testing.T) {
		// 1. Create a record to delete
		dataValue, err := anypb.New(&health.SleepData{
			UserId:        uuid.NewString(),
			SleepDuration: int64(8 * time.Hour),
			SleepQuality:  "Good",
			RecordedDate:  time.Now().Format("2006-01-02"),
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testLifestyleData := &health.LifestyleData{
			UserId:       uuid.NewString(),
			DataType:     "Sleep",
			DataValue:    dataValue,
			RecordedDate: time.Now().Format("2006-01-02"),
		}
		createdID, err := lifestyleDataRepo.CreateLifestyleData(context.Background(), testLifestyleData)
		assert.NoError(t, err, "Creating lifestyle data for DeleteLifestyleData test failed")
		assert.NotEmpty(t, createdID, "Created lifestyle data should have a valid ID")

		// 2. Delete the record
		err = lifestyleDataRepo.DeleteLifestyleData(context.Background(), createdID)
		assert.NoError(t, err, "DeleteLifestyleData should not return an error")

		// 3. Attempt to retrieve the deleted record (should fail)
		retrievedRecord, err := lifestyleDataRepo.GetLifestyleData(context.Background(), createdID)
		assert.Error(t, err, "GetLifestyleData after delete should return an error")
		assert.Nil(t, retrievedRecord, "GetLifestyleData response should be nil after delete")
	})

	t.Run("ListLifestyleData", func(t *testing.T) {
		// 1. Create some records for a specific user
		userID := uuid.NewString()
		testRecords := []*health.LifestyleData{
			{
				UserId:       userID,
				DataType:     "Sleep 1",
				DataValue:    createSampleAny(t),
				RecordedDate: time.Now().Format("2006-01-02"),
			},
			{
				UserId:       userID,
				DataType:     "Sleep 2",
				DataValue:    createSampleAny(t),
				RecordedDate: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
			},
		}

		for _, record := range testRecords {
			_, err := lifestyleDataRepo.CreateLifestyleData(context.Background(), record)
			assert.NoError(t, err, "Creating lifestyle data for ListLifestyleData test failed")
		}

		// 2. Test listing all records for the user
		req := &health.ListLifestyleDataRequest{
			UserId: userID,
		}
		retrievedRecords, err := lifestyleDataRepo.ListLifestyleData(context.Background(), req)
		assert.NoError(t, err, "ListLifestyleData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListLifestyleData response should not be nil")
		assert.GreaterOrEqual(t, len(retrievedRecords), 2, "Should have at least two lifestyle data records for the user")

		// 3. Test filtering by DataType
		req = &health.ListLifestyleDataRequest{
			UserId:   userID,
			DataType: "Sleep 1",
		}
		retrievedRecords, err = lifestyleDataRepo.ListLifestyleData(context.Background(), req)
		assert.NoError(t, err, "ListLifestyleData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListLifestyleData response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one lifestyle data record matching the filter")
		assert.Equal(t, "Sleep 1", retrievedRecords[0].DataType, "DataType should match the filter")

		// 4. Test filtering by RecordedDate
		req = &health.ListLifestyleDataRequest{
			UserId:       userID,
			RecordedDate: time.Now().Format("2006-01-02"),
		}
		retrievedRecords, err = lifestyleDataRepo.ListLifestyleData(context.Background(), req)
		assert.NoError(t, err, "ListLifestyleData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListLifestyleData response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one lifestyle data record matching the filter")
		assert.Equal(t, time.Now().Format("2006-01-02"), retrievedRecords[0].RecordedDate, "RecordedDate should match the filter")
	})
}
