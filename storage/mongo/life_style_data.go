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
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

// LifestyleDataRepo implements the storage.LifestyleDataRepoI interface for MongoDB.
type LifestyleDataRepo struct {
	db *mongo.Database
}

// NewLifestyleDataRepo creates a new LifestyleDataRepo instance.
func NewLifestyleDataRepo(db *mongo.Database) *LifestyleDataRepo {
	return &LifestyleDataRepo{
		db: db,
	}
}

// CreateLifestyleData creates a new lifestyle data record in the database.
func (r *LifestyleDataRepo) CreateLifestyleData(ctx context.Context, data *health.LifestyleData) (string, error) {
	var (
		objectID primitive.ObjectID
		err      error
	)

	if data.Id != "" {
		objectID, err = primitive.ObjectIDFromHex(data.Id)
		if err != nil {
			return "", err
		}
	} else {
		objectID = primitive.NewObjectID()
	}

	// Convert the Any proto message to a JSON string
	dataValueJSON, err := protojson.Marshal(data.DataValue)
	if err != nil {
		return "", err
	}

	// Convert the model to a BSON document
	bsonData := bson.M{
		"_id":           objectID,
		"user_id":       data.UserId,
		"data_type":     data.DataType,
		"data_value":    string(dataValueJSON),
		"recorded_date": data.RecordedDate,
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	}

	// Insert the document into the collection
	result, err := r.db.Collection("lifestyle_data").InsertOne(ctx, bsonData)
	if err != nil {
		return "", fmt.Errorf("failed to create lifestyle data: %w", err)
	}

	// Get the inserted ID as a string
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert inserted ID to string")
	}
	return insertedID.Hex(), nil
}

// GetLifestyleData retrieves a lifestyle data record by its ID.
func (r *LifestyleDataRepo) GetLifestyleData(ctx context.Context, id string) (*health.LifestyleData, error) {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid lifestyle data ID: %w", err)
	}

	// Find the document by ID
	var bsonData bson.M
	err = r.db.Collection("lifestyle_data").FindOne(ctx, bson.M{"_id": objID}).Decode(&bsonData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "lifestyle data not found")
		}
		return nil, fmt.Errorf("failed to get lifestyle data by ID: %w", err)
	}

	// Convert the BSON document to a proto message
	dataModel, err := bsonToLifestyleData(bsonData)
	if err != nil {
		return nil, err
	}
	return dataModel, nil
}

// UpdateLifestyleData updates an existing lifestyle data record in the database.
func (r *LifestyleDataRepo) UpdateLifestyleData(ctx context.Context, data *health.LifestyleData) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return fmt.Errorf("invalid lifestyle data ID: %w", err)
	}

	// Convert the Any proto message to a JSON string
	dataValueJSON, err := protojson.Marshal(data.DataValue)
	if err != nil {
		return err
	}

	// Convert the model to a BSON document
	bsonData := bson.M{
		"user_id":       data.UserId,
		"data_type":     data.DataType,
		"data_value":    string(dataValueJSON), // Store as string
		"recorded_date": data.RecordedDate,
		"updated_at":    time.Now(),
	}

	// Update the document in the collection
	result, err := r.db.Collection("lifestyle_data").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bsonData})
	if err != nil {
		return fmt.Errorf("failed to update lifestyle data: %w", err)
	}

	if result.ModifiedCount == 0 {
		return status.Errorf(codes.NotFound, "lifestyle data not found")
	}

	return nil
}

// DeleteLifestyleData deletes a lifestyle data record from the database.
func (r *LifestyleDataRepo) DeleteLifestyleData(ctx context.Context, id string) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid lifestyle data ID: %w", err)
	}

	// Delete the document from the collection
	result, err := r.db.Collection("lifestyle_data").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete lifestyle data: %w", err)
	}

	if result.DeletedCount == 0 {
		return status.Errorf(codes.NotFound, "lifestyle data not found")
	}

	return nil
}

// ListLifestyleData retrieves all lifestyle data records for a given user ID, applying filters if provided.
func (r *LifestyleDataRepo) ListLifestyleData(ctx context.Context, req *health.ListLifestyleDataRequest) ([]*health.LifestyleData, error) {
	// Build the filter query based on the request parameters
	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.DataType != "" {
		filter["data_type"] = req.DataType
	}
	if req.RecordedDate != "" {
		filter["recorded_date"] = req.RecordedDate
	}

	// Find the documents based on the filter
	cursor, err := r.db.Collection("lifestyle_data").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list lifestyle data: %w", err)
	}
	defer cursor.Close(ctx)

	var lifestyleDataRecords []*health.LifestyleData

	// Iterate through the cursor and convert each document to a proto message
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

// bsonToLifestyleData converts a BSON document to a health.LifestyleData proto message.
func bsonToLifestyleData(bsonData bson.M) (*health.LifestyleData, error) {
	dataModel := &health.LifestyleData{}

	// Handle _id field as string
	if val, ok := bsonData["_id"].(string); ok {
		dataModel.Id = val
	} else if oid, ok := bsonData["_id"].(primitive.ObjectID); ok {
		dataModel.Id = oid.Hex()
	} else {
		return nil, fmt.Errorf("invalid _id type: %T", bsonData["_id"])
	}

	// Handle potentially nil fields
	if val, ok := bsonData["user_id"].(string); ok {
		dataModel.UserId = val
	}
	if val, ok := bsonData["data_type"].(string); ok {
		dataModel.DataType = val
	}
	if val, ok := bsonData["recorded_date"].(string); ok {
		dataModel.RecordedDate = val
	}

	// Convert data_value from string to Any proto message
	if val, ok := bsonData["data_value"].(string); ok {
		dataModel.DataValue = &anypb.Any{}
		err := protojson.Unmarshal([]byte(val), dataModel.DataValue)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data_value from JSON: %w", err)
		}
	}

	// Convert created_at and updated_at fields
	if val, ok := bsonData["created_at"].(primitive.DateTime); ok {
		dataModel.CreatedAt = val.Time().Format(time.RFC3339)
	}
	if val, ok := bsonData["updated_at"].(primitive.DateTime); ok {
		dataModel.UpdatedAt = val.Time().Format(time.RFC3339)
	}

	return dataModel, nil
}
