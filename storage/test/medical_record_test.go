package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	mongodb "github.com/health-analytics-service/health-analytics-service/storage/mongo"
	"github.com/stretchr/testify/assert"
)

func TestMedicalRecordRepo(t *testing.T) {
	db := createMongoDBConnection(t)
	medicalRecordRepo := mongodb.NewMedicalRecordRepo(db)

	t.Run("CreateMedicalRecord", func(t *testing.T) {
		testRecord := &health.MedicalRecord{
			UserId:      uuid.NewString(),
			RecordType:  "Test Record",
			RecordDate:  time.Now().Format("2006-01-02"),
			Description: "This is a test medical record.",
			DoctorId:    uuid.NewString(),
			Attachments: []string{"attachment1.txt", "attachment2.pdf"},
		}
		createdID, err := medicalRecordRepo.CreateMedicalRecord(context.Background(), testRecord)

		assert.NoError(t, err, "CreateMedicalRecord should not return an error")
		assert.NotEmpty(t, createdID, "Created medical record should have a valid ID")
	})

	t.Run("GetMedicalRecord", func(t *testing.T) {
		// 1. Create a record to retrieve
		testRecord := &health.MedicalRecord{
			UserId:      uuid.NewString(),
			RecordType:  "Test Record for Get",
			RecordDate:  time.Now().Format("2006-01-02"),
			Description: "This is a test medical record for Get.",
			DoctorId:    uuid.NewString(),
			Attachments: []string{"attachment1.txt", "attachment2.pdf"},
		}
		createdID, err := medicalRecordRepo.CreateMedicalRecord(context.Background(), testRecord)
		assert.NoError(t, err, "Creating medical record for GetMedicalRecord test failed")
		assert.NotEmpty(t, createdID, "Created medical record should have a valid ID")

		// 2. Get the record
		retrievedRecord, err := medicalRecordRepo.GetMedicalRecord(context.Background(), createdID)

		assert.NoError(t, err, "GetMedicalRecord should not return an error")
		assert.NotNil(t, retrievedRecord, "GetMedicalRecord response should not be nil")
		assert.Equal(t, testRecord.UserId, retrievedRecord.UserId)
		assert.Equal(t, testRecord.RecordType, retrievedRecord.RecordType)
		assert.Equal(t, testRecord.RecordDate, retrievedRecord.RecordDate)
		assert.Equal(t, testRecord.Description, retrievedRecord.Description)
		assert.Equal(t, testRecord.DoctorId, retrievedRecord.DoctorId)
		assert.Equal(t, testRecord.Attachments, retrievedRecord.Attachments)
	})

	t.Run("UpdateMedicalRecord", func(t *testing.T) {
		// 1. Create a record to update
		testRecord := &health.MedicalRecord{
			UserId:      uuid.NewString(),
			RecordType:  "Test Record for Update",
			RecordDate:  time.Now().Format("2006-01-02"),
			Description: "This is a test medical record for Update.",
			DoctorId:    uuid.NewString(),
			Attachments: []string{"attachment1.txt", "attachment2.pdf"},
		}
		createdID, err := medicalRecordRepo.CreateMedicalRecord(context.Background(), testRecord)
		assert.NoError(t, err, "Creating medical record for UpdateMedicalRecord test failed")
		assert.NotEmpty(t, createdID, "Created medical record should have a valid ID")

		// 2. Update the record
		updateRecord := &health.MedicalRecord{
			Id:          createdID,
			UserId:      uuid.NewString(),
			RecordType:  "Updated Record Type",
			Description: "Updated description.",
		}
		err = medicalRecordRepo.UpdateMedicalRecord(context.Background(), updateRecord)
		assert.NoError(t, err, "UpdateMedicalRecord should not return an error")

		// 3. Retrieve the record and verify the update
		retrievedRecord, err := medicalRecordRepo.GetMedicalRecord(context.Background(), createdID)
		assert.NoError(t, err, "GetMedicalRecord after update should not return an error")
		assert.Equal(t, updateRecord.UserId, retrievedRecord.UserId, "UserId should be updated")
		assert.Equal(t, updateRecord.RecordType, retrievedRecord.RecordType, "RecordType should be updated")
		assert.Equal(t, updateRecord.Description, retrievedRecord.Description, "Description should be updated")
		// Other fields should remain unchanged
		assert.Equal(t, testRecord.RecordDate, retrievedRecord.RecordDate, "RecordDate should remain unchanged")
		assert.Equal(t, testRecord.DoctorId, retrievedRecord.DoctorId, "DoctorId should remain unchanged")
		// assert.Equal(t, testRecord.Attachments, retrievedRecord.Attachments, "Attachments should remain unchanged")
	})

	t.Run("DeleteMedicalRecord", func(t *testing.T) {
		// 1. Create a record to delete
		testRecord := &health.MedicalRecord{
			UserId:      uuid.NewString(),
			RecordType:  "Test Record for Delete",
			RecordDate:  time.Now().Format("2006-01-02"),
			Description: "This is a test medical record for Delete.",
			DoctorId:    uuid.NewString(),
			Attachments: []string{"attachment1.txt", "attachment2.pdf"},
		}
		createdID, err := medicalRecordRepo.CreateMedicalRecord(context.Background(), testRecord)
		assert.NoError(t, err, "Creating medical record for DeleteMedicalRecord test failed")
		assert.NotEmpty(t, createdID, "Created medical record should have a valid ID")

		// 2. Delete the record
		err = medicalRecordRepo.DeleteMedicalRecord(context.Background(), createdID)
		assert.NoError(t, err, "DeleteMedicalRecord should not return an error")

		// 3. Attempt to retrieve the deleted record (should fail)
		retrievedRecord, err := medicalRecordRepo.GetMedicalRecord(context.Background(), createdID)
		assert.Error(t, err, "GetMedicalRecord after delete should return an error")
		assert.Nil(t, retrievedRecord, "GetMedicalRecord response should be nil after delete")
	})

	t.Run("ListMedicalRecords", func(t *testing.T) {
		// 1. Create some records for a specific user
		userID := uuid.NewString()
		testRecords := []*health.MedicalRecord{
			{
				UserId:      userID,
				RecordType:  "Test Record 1 for List",
				RecordDate:  time.Now().Format("2006-01-02"),
				Description: "This is test medical record 1 for List.",
				DoctorId:    uuid.NewString(),
				Attachments: []string{"attachment1.txt", "attachment2.pdf"},
			},
			{
				UserId:      userID,
				RecordType:  "Test Record 2 for List",
				RecordDate:  time.Now().Add(time.Hour * 24).Format("2006-01-02"),
				Description: "This is test medical record 2 for List.",
				DoctorId:    uuid.NewString(),
				Attachments: []string{"attachment3.txt", "attachment4.pdf"},
			},
		}

		for _, record := range testRecords {
			_, err := medicalRecordRepo.CreateMedicalRecord(context.Background(), record)
			assert.NoError(t, err, "Creating medical record for ListMedicalRecords test failed")
		}

		// 2. Test listing all records for the user
		req := &health.ListMedicalRecordsRequest{
			UserId: userID,
		}
		retrievedRecords, err := medicalRecordRepo.ListMedicalRecords(context.Background(), req)
		assert.NoError(t, err, "ListMedicalRecords should not return an error")
		assert.NotNil(t, retrievedRecords, "ListMedicalRecords response should not be nil")
		assert.GreaterOrEqual(t, len(retrievedRecords), 2, "Should have at least two medical records for the user")

		// 3. Test filtering by RecordType
		req = &health.ListMedicalRecordsRequest{
			UserId:     userID,
			RecordType: "Test Record 1 for List",
		}
		retrievedRecords, err = medicalRecordRepo.ListMedicalRecords(context.Background(), req)
		assert.NoError(t, err, "ListMedicalRecords should not return an error")
		assert.NotNil(t, retrievedRecords, "ListMedicalRecords response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one medical record matching the filter")
		assert.Equal(t, "Test Record 1 for List", retrievedRecords[0].RecordType, "RecordType should match the filter")

		// 4. Test filtering by RecordDate
		req = &health.ListMedicalRecordsRequest{
			UserId:     userID,
			RecordDate: time.Now().Format("2006-01-02"),
		}
		retrievedRecords, err = medicalRecordRepo.ListMedicalRecords(context.Background(), req)
		assert.NoError(t, err, "ListMedicalRecords should not return an error")
		assert.NotNil(t, retrievedRecords, "ListMedicalRecords response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one medical record matching the filter")
		assert.Equal(t, time.Now().Format("2006-01-02"), retrievedRecords[0].RecordDate, "RecordDate should match the filter")

		// 5. Test filtering by Description
		req = &health.ListMedicalRecordsRequest{
			UserId:      userID,
			Description: "This is test medical record 1 for List.",
		}
		retrievedRecords, err = medicalRecordRepo.ListMedicalRecords(context.Background(), req)
		assert.NoError(t, err, "ListMedicalRecords should not return an error")
		assert.NotNil(t, retrievedRecords, "ListMedicalRecords response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one medical record matching the filter")
		assert.Equal(t, "This is test medical record 1 for List.", retrievedRecords[0].Description, "Description should match the filter")

		// 6. Test filtering by DoctorId
		doctorID := testRecords[0].DoctorId // Use the DoctorId from the first test record
		req = &health.ListMedicalRecordsRequest{
			UserId:   userID,
			DoctorId: doctorID,
		}
		retrievedRecords, err = medicalRecordRepo.ListMedicalRecords(context.Background(), req)
		assert.NoError(t, err, "ListMedicalRecords should not return an error")
		assert.NotNil(t, retrievedRecords, "ListMedicalRecords response should not be nil")
		assert.Equal(t, 1, len(retrievedRecords), "Should have one medical record matching the filter")
		assert.Equal(t, doctorID, retrievedRecords[0].DoctorId, "DoctorId should match the filter")
	})
}
