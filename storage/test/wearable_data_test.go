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

// ... (Your existing createMongoDBConnection function) ...

func TestWearableDataRepo(t *testing.T) {
	db := createMongoDBConnection(t)
	wearableDataRepo := mongodb.NewWearableDataRepo(db)

	t.Run("CreateWearableData", func(t *testing.T) {
		// Create a sample Any proto message
		dataValue, err := anypb.New(&health.HeartRateData{
			UserId:            uuid.NewString(),
			HeartRate:         80,
			RecordedTimestamp: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testWearableData := &health.WearableData{
			UserId:            uuid.NewString(),
			DeviceType:        "Smartwatch",
			DataType:          "HeartRate",
			DataValue:         dataValue,
			RecordedTimestamp: time.Now().Format(time.RFC3339),
		}
		createdID, err := wearableDataRepo.CreateWearableData(context.Background(), testWearableData)

		assert.NoError(t, err, "CreateWearableData should not return an error")
		assert.NotEmpty(t, createdID, "Created wearable data should have a valid ID")
	})

	t.Run("GetWearableData", func(t *testing.T) {
		// 1. Create a record to retrieve
		dataValue, err := anypb.New(&health.HeartRateData{
			UserId:            uuid.NewString(),
			HeartRate:         80,
			RecordedTimestamp: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testWearableData := &health.WearableData{
			UserId:            uuid.NewString(),
			DeviceType:        "Smartwatch",
			DataType:          "HeartRate",
			DataValue:         dataValue,
			RecordedTimestamp: time.Now().Format(time.RFC3339),
		}
		createdID, err := wearableDataRepo.CreateWearableData(context.Background(), testWearableData)
		assert.NoError(t, err, "Creating wearable data for GetWearableData test failed")
		assert.NotEmpty(t, createdID, "Created wearable data should have a valid ID")

		// 2. Get the record
		retrievedData, err := wearableDataRepo.GetWearableData(context.Background(), createdID)

		assert.NoError(t, err, "GetWearableData should not return an error")
		assert.NotNil(t, retrievedData, "GetWearableData response should not be nil")
		assert.Equal(t, testWearableData.UserId, retrievedData.UserId)
		assert.Equal(t, testWearableData.DeviceType, retrievedData.DeviceType)
		assert.Equal(t, testWearableData.DataType, retrievedData.DataType)
		assert.Equal(t, testWearableData.RecordedTimestamp, retrievedData.RecordedTimestamp)
		// Compare the Any proto messages
		assert.Equal(t, testWearableData.DataValue.String(), retrievedData.DataValue.String())
	})

	t.Run("UpdateWearableData", func(t *testing.T) {
		// 1. Create a record to update
		dataValue, err := anypb.New(&health.HeartRateData{
			UserId:            uuid.NewString(),
			HeartRate:         80,
			RecordedTimestamp: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testWearableData := &health.WearableData{
			UserId:            uuid.NewString(),
			DeviceType:        "Smartwatch",
			DataType:          "HeartRate",
			DataValue:         dataValue,
			RecordedTimestamp: time.Now().Format(time.RFC3339),
		}
		createdID, err := wearableDataRepo.CreateWearableData(context.Background(), testWearableData)
		assert.NoError(t, err, "Creating wearable data for UpdateWearableData test failed")
		assert.NotEmpty(t, createdID, "Created wearable data should have a valid ID")

		// 2. Update the record
		updatedDataValue, err := anypb.New(&health.HeartRateData{
			UserId:            uuid.NewString(),
			HeartRate:         90,
			RecordedTimestamp: time.Now().Add(time.Hour).Format(time.RFC3339),
		})
		assert.NoError(t, err, "Failed to create updated Any proto message")

		updateRecord := &health.WearableData{
			Id:                createdID,
			UserId:            uuid.NewString(),
			DeviceType:        "Fitness Tracker",
			DataType:          "Updated HeartRate",
			DataValue:         updatedDataValue,
			RecordedTimestamp: time.Now().Add(time.Hour).Format(time.RFC3339),
		}
		err = wearableDataRepo.UpdateWearableData(context.Background(), updateRecord)
		assert.NoError(t, err, "UpdateWearableData should not return an error")

		// 3. Retrieve the record and verify the update
		retrievedRecord, err := wearableDataRepo.GetWearableData(context.Background(), createdID)
		assert.NoError(t, err, "GetWearableData after update should not return an error")
		assert.Equal(t, updateRecord.UserId, retrievedRecord.UserId, "UserId should be updated")
		assert.Equal(t, updateRecord.DeviceType, retrievedRecord.DeviceType, "DeviceType should be updated")
		assert.Equal(t, updateRecord.DataType, retrievedRecord.DataType, "DataType should be updated")
		assert.Equal(t, updateRecord.RecordedTimestamp, retrievedRecord.RecordedTimestamp, "RecordedTimestamp should be updated")
		// Compare the Any proto messages
		assert.Equal(t, updateRecord.DataValue.String(), retrievedRecord.DataValue.String())
	})

	t.Run("DeleteWearableData", func(t *testing.T) {
		// 1. Create a record to delete
		dataValue, err := anypb.New(&health.HeartRateData{
			UserId:            uuid.NewString(),
			HeartRate:         80,
			RecordedTimestamp: time.Now().Format(time.RFC3339),
		})
		assert.NoError(t, err, "Failed to create Any proto message")

		testWearableData := &health.WearableData{
			UserId:            uuid.NewString(),
			DeviceType:        "Smartwatch",
			DataType:          "HeartRate",
			DataValue:         dataValue,
			RecordedTimestamp: time.Now().Format(time.RFC3339),
		}
		createdID, err := wearableDataRepo.CreateWearableData(context.Background(), testWearableData)
		assert.NoError(t, err, "Creating wearable data for DeleteWearableData test failed")
		assert.NotEmpty(t, createdID, "Created wearable data should have a valid ID")

		// 2. Delete the record
		err = wearableDataRepo.DeleteWearableData(context.Background(), createdID)
		assert.NoError(t, err, "DeleteWearableData should not return an error")

		// 3. Attempt to retrieve the deleted record (should fail)
		retrievedRecord, err := wearableDataRepo.GetWearableData(context.Background(), createdID)
		assert.Error(t, err, "GetWearableData after delete should return an error")
		assert.Nil(t, retrievedRecord, "GetWearableData response should be nil after delete")
	})

	t.Run("ListWearableData", func(t *testing.T) {
		// 1. Create some records for a specific user
		userID := uuid.NewString()
		testRecords := []*health.WearableData{
			{
				UserId:            userID,
				DeviceType:        "Smartwatch",
				DataType:          "HeartRate 1",
				DataValue:         createSampleAny(t),
				RecordedTimestamp: time.Now().Format(time.RFC3339),
			},
			{
				UserId:            userID,
				DeviceType:        "Fitness Tracker",
				DataType:          "HeartRate 2",
				DataValue:         createSampleAny(t),
				RecordedTimestamp: time.Now().Add(time.Hour).Format(time.RFC3339),
			},
		}

		for _, record := range testRecords {
			_, err := wearableDataRepo.CreateWearableData(context.Background(), record)
			assert.NoError(t, err, "Creating wearable data for ListWearableData test failed")
		}

		// 2. Test listing all records for the user
		req := &health.ListWearableDataRequest{
			UserId: userID,
		}
		retrievedRecords, err := wearableDataRepo.ListWearableData(context.Background(), req)
		assert.NoError(t, err, "ListWearableData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListWearableData response should not be nil")
		assert.GreaterOrEqual(t, len(retrievedRecords), 2, "Should have at least two wearable data records for the user")

		// 3. Test filtering by DeviceType
		req = &health.ListWearableDataRequest{
			UserId:     userID,
			DeviceType: "Smartwatch",
		}
		retrievedRecords, err = wearableDataRepo.ListWearableData(context.Background(), req)
		assert.NoError(t, err, "ListWearableData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListWearableData response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one wearable data record matching the filter")
		assert.Equal(t, "Smartwatch", retrievedRecords[0].DeviceType, "DeviceType should match the filter")

		// 4. Test filtering by DataType
		req = &health.ListWearableDataRequest{
			UserId:   userID,
			DataType: "HeartRate 1",
		}
		retrievedRecords, err = wearableDataRepo.ListWearableData(context.Background(), req)
		assert.NoError(t, err, "ListWearableData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListWearableData response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one wearable data record matching the filter")
		assert.Equal(t, "HeartRate 1", retrievedRecords[0].DataType, "DataType should match the filter")

		// 5. Test filtering by RecordedTimestamp
		recordedTimestamp := testRecords[0].RecordedTimestamp
		req = &health.ListWearableDataRequest{
			UserId:            userID,
			RecordedTimestamp: recordedTimestamp,
		}
		retrievedRecords, err = wearableDataRepo.ListWearableData(context.Background(), req)
		assert.NoError(t, err, "ListWearableData should not return an error")
		assert.NotNil(t, retrievedRecords, "ListWearableData response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one wearable data record matching the filter")
		assert.Equal(t, recordedTimestamp, retrievedRecords[0].RecordedTimestamp, "RecordedTimestamp should match the filter")
	})
}
