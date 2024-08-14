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

// WearableDataRepo implements the storage.WearableDataRepoI interface for MongoDB.
type WearableDataRepo struct {
	db *mongo.Database
}

// NewWearableDataRepo creates a new WearableDataRepo instance.
func NewWearableDataRepo(db *mongo.Database) *WearableDataRepo {
	return &WearableDataRepo{
		db: db,
	}
}

// CreateWearableData creates a new wearable data record in the database.
func (r *WearableDataRepo) CreateWearableData(ctx context.Context, data *health.WearableData) (string, error) {
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
		"_id":                objectID,
		"user_id":            data.UserId,
		"device_type":        data.DeviceType,
		"data_type":          data.DataType,
		"data_value":         string(dataValueJSON), // Store as string
		"recorded_timestamp": data.RecordedTimestamp,
		"created_at":         time.Now(),
		"updated_at":         time.Now(),
	}

	// Insert the document into the collection
	result, err := r.db.Collection("wearable_data").InsertOne(ctx, bsonData)
	if err != nil {
		return "", fmt.Errorf("failed to create wearable data: %w", err)
	}

	// Get the inserted ID as a string
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert inserted ID to string")
	}
	return insertedID.Hex(), nil
}

// GetWearableData retrieves a wearable data record by its ID.
func (r *WearableDataRepo) GetWearableData(ctx context.Context, id string) (*health.WearableData, error) {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid wearable data ID: %w", err)
	}

	// Find the document by ID
	var bsonData bson.M
	err = r.db.Collection("wearable_data").FindOne(ctx, bson.M{"_id": objID}).Decode(&bsonData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "wearable data not found")
		}
		return nil, fmt.Errorf("failed to get wearable data by ID: %w", err)
	}

	// Convert the BSON document to a proto message
	dataModel, err := bsonToWearableData(bsonData)
	if err != nil {
		return nil, err
	}
	return dataModel, nil
}

// UpdateWearableData updates an existing wearable data record in the database.
func (r *WearableDataRepo) UpdateWearableData(ctx context.Context, data *health.WearableData) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return fmt.Errorf("invalid wearable data ID: %w", err)
	}

	// Convert the Any proto message to a JSON string
	dataValueJSON, err := protojson.Marshal(data.DataValue)
	if err != nil {
		return err
	}

	// Convert the model to a BSON document
	bsonData := bson.M{
		"user_id":            data.UserId,
		"device_type":        data.DeviceType,
		"data_type":          data.DataType,
		"data_value":         string(dataValueJSON), // Store as string
		"recorded_timestamp": data.RecordedTimestamp,
		"updated_at":         time.Now(),
	}

	// Update the document in the collection
	result, err := r.db.Collection("wearable_data").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bsonData})
	if err != nil {
		return fmt.Errorf("failed to update wearable data: %w", err)
	}

	if result.ModifiedCount == 0 {
		return status.Errorf(codes.NotFound, "wearable data not found")
	}

	return nil
}

// DeleteWearableData deletes a wearable data record from the database.
func (r *WearableDataRepo) DeleteWearableData(ctx context.Context, id string) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid wearable data ID: %w", err)
	}

	// Delete the document from the collection
	result, err := r.db.Collection("wearable_data").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete wearable data: %w", err)
	}

	if result.DeletedCount == 0 {
		return status.Errorf(codes.NotFound, "wearable data not found")
	}

	return nil
}

// ListWearableData retrieves all wearable data records for a given user ID, applying filters if provided.
func (r *WearableDataRepo) ListWearableData(ctx context.Context, req *health.ListWearableDataRequest) ([]*health.WearableData, error) {
	// Build the filter query based on the request parameters
	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.DeviceType != "" {
		filter["device_type"] = req.DeviceType
	}
	if req.DataType != "" {
		filter["data_type"] = req.DataType
	}
	if req.RecordedTimestamp != "" {
		filter["recorded_timestamp"] = req.RecordedTimestamp
	}

	// Find the documents based on the filter
	cursor, err := r.db.Collection("wearable_data").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list wearable data: %w", err)
	}
	defer cursor.Close(ctx)

	var wearableDataRecords []*health.WearableData

	// Iterate through the cursor and convert each document to a proto message
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

// bsonToWearableData converts a BSON document to a health.WearableData proto message.
func bsonToWearableData(bsonData bson.M) (*health.WearableData, error) {
	dataModel := &health.WearableData{}

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
	if val, ok := bsonData["device_type"].(string); ok {
		dataModel.DeviceType = val
	}
	if val, ok := bsonData["data_type"].(string); ok {
		dataModel.DataType = val
	}
	if val, ok := bsonData["recorded_timestamp"].(string); ok {
		dataModel.RecordedTimestamp = val
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
