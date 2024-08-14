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

// GeneticDataRepo implements the storage.GeneticDataRepoI interface for MongoDB.
type GeneticDataRepo struct {
	db *mongo.Database
}

// NewGeneticDataRepo creates a new GeneticDataRepo instance.
func NewGeneticDataRepo(db *mongo.Database) *GeneticDataRepo {
	return &GeneticDataRepo{
		db: db,
	}
}

// CreateGeneticData creates a new genetic data record in the database.
func (r *GeneticDataRepo) CreateGeneticData(ctx context.Context, data *health.GeneticData) (string, error) {
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
	// Convert the Any proto message to a BSON document
	dataVal, err := protojson.Marshal(data.DataValue)
	if err != nil {
		return "", err
	}

	// Convert the model to a BSON document
	bsonData := bson.M{
		"_id":           objectID,
		"user_id":       data.UserId,
		"data_type":     data.DataType,
		"data_value":    string(dataVal),
		"analysis_date": data.AnalysisDate,
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	}

	// Insert the document into the collection
	result, err := r.db.Collection("genetic_data").InsertOne(ctx, bsonData)
	if err != nil {
		return "", fmt.Errorf("failed to create genetic data: %w", err)
	}

	// Get the inserted ID as a string
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert inserted ID to string")
	}
	return insertedID.Hex(), nil
}

// GetGeneticData retrieves a genetic data record by its ID.
func (r *GeneticDataRepo) GetGeneticData(ctx context.Context, id string) (*health.GeneticData, error) {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid genetic data ID: %w", err)
	}

	// Find the document by ID
	var bsonData bson.M
	err = r.db.Collection("genetic_data").FindOne(ctx, bson.M{"_id": objID}).Decode(&bsonData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "genetic data not found")
		}
		return nil, fmt.Errorf("failed to get genetic data by ID: %w", err)
	}

	// Convert the BSON document to a proto message
	dataModel, err := bsonToGeneticData(bsonData)
	if err != nil {
		return nil, err
	}
	return dataModel, nil
}

// UpdateGeneticData updates an existing genetic data record in the database.
func (r *GeneticDataRepo) UpdateGeneticData(ctx context.Context, data *health.GeneticData) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return fmt.Errorf("invalid genetic data ID: %w", err)
	}
	dataVal, err := protojson.Marshal(data.DataValue)
	if err != nil {
		return err
	}
	// Convert the model to a BSON document
	bsonData := bson.M{
		"user_id":       data.UserId,
		"data_type":     data.DataType,
		"data_value":    string(dataVal),
		"analysis_date": data.AnalysisDate,
		"updated_at":    time.Now(),
	}

	// Update the document in the collection
	result, err := r.db.Collection("genetic_data").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bsonData})
	if err != nil {
		return fmt.Errorf("failed to update genetic data: %w", err)
	}

	if result.ModifiedCount == 0 {
		return status.Errorf(codes.NotFound, "genetic data not found")
	}

	return nil
}

// DeleteGeneticData deletes a genetic data record from the database.
func (r *GeneticDataRepo) DeleteGeneticData(ctx context.Context, id string) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid genetic data ID: %w", err)
	}

	// Delete the document from the collection
	result, err := r.db.Collection("genetic_data").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete genetic data: %w", err)
	}

	if result.DeletedCount == 0 {
		return status.Errorf(codes.NotFound, "genetic data not found")
	}

	return nil
}

// ListGeneticData retrieves all genetic data records for a given user ID, applying filters if provided.
func (r *GeneticDataRepo) ListGeneticData(ctx context.Context, req *health.ListGeneticDataRequest) ([]*health.GeneticData, error) {
	// Build the filter query based on the request parameters
	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.DataType != "" {
		filter["data_type"] = req.DataType
	}
	if req.AnalysisDate != "" {
		filter["analysis_date"] = req.AnalysisDate
	}

	// Find the documents based on the filter
	cursor, err := r.db.Collection("genetic_data").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list genetic data: %w", err)
	}
	defer cursor.Close(ctx)

	var geneticDataRecords []*health.GeneticData

	// Iterate through the cursor and convert each document to a proto message
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

// bsonToGeneticData converts a BSON document to a health.GeneticData proto message.
func bsonToGeneticData(bsonData bson.M) (*health.GeneticData, error) {
	dataModel := &health.GeneticData{}

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
	if val, ok := bsonData["analysis_date"].(string); ok {
		dataModel.AnalysisDate = val
	}
	// Convert data_value from BSON to Any proto message
	if val, ok := bsonData["data_value"].(string); ok { // Correct type assertion here
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
