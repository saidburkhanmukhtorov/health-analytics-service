package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MedicalRecordRepo implements the storage.MedicalRecordRepoI interface for MongoDB.
type MedicalRecordRepo struct {
	db *mongo.Database
}

// NewMedicalRecordRepo creates a new MedicalRecordRepo instance.
func NewMedicalRecordRepo(db *mongo.Database) *MedicalRecordRepo {
	return &MedicalRecordRepo{
		db: db,
	}
}

// CreateMedicalRecord creates a new medical record in the database.
func (r *MedicalRecordRepo) CreateMedicalRecord(ctx context.Context, record *health.MedicalRecord) (string, error) {
	var (
		objectID primitive.ObjectID
		err      error
	)

	if record.Id != "" {
		objectID, err = primitive.ObjectIDFromHex(record.Id)
		if err != nil {
			return "", err
		}
	} else {
		objectID = primitive.NewObjectID()
	}

	// Convert the model to a BSON document
	bsonRecord := bson.M{
		"_id":         objectID,
		"user_id":     record.UserId,
		"record_type": record.RecordType,
		"record_date": record.RecordDate,
		"description": record.Description,
		"doctor_id":   record.DoctorId,
		"attachments": record.Attachments,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	// Insert the document into the collection
	result, err := r.db.Collection("medical_records").InsertOne(ctx, bsonRecord)
	if err != nil {
		return "", fmt.Errorf("failed to create medical record: %w", err)
	}

	// Get the inserted ID as a string
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert inserted ID to string")
	}
	return insertedID.Hex(), nil
}

// GetMedicalRecord retrieves a medical record by its ID.
func (r *MedicalRecordRepo) GetMedicalRecord(ctx context.Context, id string) (*health.MedicalRecord, error) {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid medical record ID: %w", err)
	}

	// Find the document by ID
	var bsonRecord bson.M
	err = r.db.Collection("medical_records").FindOne(ctx, bson.M{"_id": objID}).Decode(&bsonRecord)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "medical record not found")
		}
		return nil, fmt.Errorf("failed to get medical record by ID: %w", err)
	}

	// Convert the BSON document to a proto message
	recordModel, err := bsonToMedicalRecord(bsonRecord)
	if err != nil {
		return nil, err
	}

	return recordModel, nil
}

// UpdateMedicalRecord updates an existing medical record in the database.
func (r *MedicalRecordRepo) UpdateMedicalRecord(ctx context.Context, record *health.MedicalRecord) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(record.Id)
	if err != nil {
		return fmt.Errorf("invalid medical record ID: %w", err)
	}

	// Build the update document based on the provided fields
	bsonRecord := bson.M{
		"updated_at": time.Now(),
	}

	if record.UserId != "" {
		bsonRecord["user_id"] = record.UserId
	}
	if record.RecordType != "" {
		bsonRecord["record_type"] = record.RecordType
	}
	if record.RecordDate != "" {
		bsonRecord["record_date"] = record.RecordDate
	}
	if record.Description != "" {
		bsonRecord["description"] = record.Description
	}
	if record.DoctorId != "" {
		bsonRecord["doctor_id"] = record.DoctorId
	}
	if len(record.Attachments) > 0 {
		bsonRecord["attachments"] = record.Attachments
	}

	// Update the document in the collection
	result, err := r.db.Collection("medical_records").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bsonRecord})
	if err != nil {
		return fmt.Errorf("failed to update medical record: %w", err)
	}

	if result.ModifiedCount == 0 {
		return status.Errorf(codes.NotFound, "medical record not found")
	}

	return nil
}

// DeleteMedicalRecord deletes a medical record from the database.
func (r *MedicalRecordRepo) DeleteMedicalRecord(ctx context.Context, id string) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid medical record ID: %w", err)
	}

	// Delete the document from the collection
	result, err := r.db.Collection("medical_records").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete medical record: %w", err)
	}

	if result.DeletedCount == 0 {
		return status.Errorf(codes.NotFound, "medical record not found")
	}

	return nil
}

// ListMedicalRecords retrieves all medical records for a given user ID, applying filters if provided.
func (r *MedicalRecordRepo) ListMedicalRecords(ctx context.Context, req *health.ListMedicalRecordsRequest) ([]*health.MedicalRecord, error) {
	// Build the filter query based on the request parameters
	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.RecordType != "" {
		filter["record_type"] = req.RecordType
	}
	if req.RecordDate != "" {
		filter["record_date"] = req.RecordDate
	}
	if req.Description != "" {
		filter["description"] = req.Description
	}
	if req.DoctorId != "" {
		filter["doctor_id"] = req.DoctorId
	}

	// Find the documents based on the filter
	cursor, err := r.db.Collection("medical_records").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list medical records: %w", err)
	}
	defer cursor.Close(ctx)

	var medicalRecords []*health.MedicalRecord

	// Iterate through the cursor and convert each document to a proto message
	for cursor.Next(ctx) {
		var bsonRecord bson.M
		if err := cursor.Decode(&bsonRecord); err != nil {
			return nil, fmt.Errorf("failed to decode medical record: %w", err)
		}

		recordModel, err := bsonToMedicalRecord(bsonRecord)
		if err != nil {
			return nil, err
		}

		medicalRecords = append(medicalRecords, recordModel)
	}

	return medicalRecords, nil
}

// bsonToMedicalRecord converts a BSON document to a health.MedicalRecord proto message.
func bsonToMedicalRecord(bsonRecord bson.M) (*health.MedicalRecord, error) {
	recordModel := &health.MedicalRecord{}

	// Handle _id field as string
	if val, ok := bsonRecord["_id"].(string); ok {
		recordModel.Id = val
	} else if oid, ok := bsonRecord["_id"].(primitive.ObjectID); ok {
		recordModel.Id = oid.Hex()
	} else {
		return nil, fmt.Errorf("invalid _id type: %T", bsonRecord["_id"])
	}

	// Handle potentially nil fields
	if val, ok := bsonRecord["user_id"].(string); ok {
		recordModel.UserId = val
	}
	if val, ok := bsonRecord["record_type"].(string); ok {
		recordModel.RecordType = val
	}
	if val, ok := bsonRecord["record_date"].(string); ok {
		recordModel.RecordDate = val
	}
	if val, ok := bsonRecord["description"].(string); ok {
		recordModel.Description = val
	}
	if val, ok := bsonRecord["doctor_id"].(string); ok {
		recordModel.DoctorId = val
	}
	if val, ok := bsonRecord["attachments"].(bson.A); ok {
		for _, attachment := range val {
			if strAttachment, ok := attachment.(string); ok {
				recordModel.Attachments = append(recordModel.Attachments, strAttachment)
			}

		}
	}
	// Convert created_at and updated_at fields
	if val, ok := bsonRecord["created_at"].(primitive.DateTime); ok {
		recordModel.CreatedAt = val.Time().Format(time.RFC3339)
	}
	if val, ok := bsonRecord["updated_at"].(primitive.DateTime); ok {
		recordModel.UpdatedAt = val.Time().Format(time.RFC3339)
	}

	return recordModel, nil
}
