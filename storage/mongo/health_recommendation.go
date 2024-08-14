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

// HealthRecommendationRepo implements the storage.HealthRecommendationRepoI interface for MongoDB.
type HealthRecommendationRepo struct {
	db *mongo.Database
}

// NewHealthRecommendationRepo creates a new HealthRecommendationRepo instance.
func NewHealthRecommendationRepo(db *mongo.Database) *HealthRecommendationRepo {
	return &HealthRecommendationRepo{
		db: db,
	}
}

// CreateHealthRecommendation creates a new health recommendation in the database.
func (r *HealthRecommendationRepo) CreateHealthRecommendation(ctx context.Context, recommendation *health.HealthRecommendation) (string, error) {
	var (
		objectID primitive.ObjectID
		err      error
	)

	if recommendation.Id != "" {
		objectID, err = primitive.ObjectIDFromHex(recommendation.Id)
		if err != nil {
			return "", err
		}
	} else {
		objectID = primitive.NewObjectID()
	}

	// Convert the model to a BSON document
	bsonRecommendation := bson.M{
		"_id":                 objectID,
		"user_id":             recommendation.UserId,
		"recommendation_type": recommendation.RecommendationType,
		"description":         recommendation.Description,
		"priority":            recommendation.Priority,
		"created_at":          time.Now(),
		"updated_at":          time.Now(),
	}

	// Insert the document into the collection
	result, err := r.db.Collection("health_recommendations").InsertOne(ctx, bsonRecommendation)
	if err != nil {
		return "", fmt.Errorf("failed to create health recommendation: %w", err)
	}

	// Get the inserted ID as a string
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert inserted ID to string")
	}
	return insertedID.Hex(), nil
}

// GetHealthRecommendation retrieves a health recommendation by its ID.
func (r *HealthRecommendationRepo) GetHealthRecommendation(ctx context.Context, id string) (*health.HealthRecommendation, error) {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid health recommendation ID: %w", err)
	}

	// Find the document by ID
	var bsonRecommendation bson.M
	err = r.db.Collection("health_recommendations").FindOne(ctx, bson.M{"_id": objID}).Decode(&bsonRecommendation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "health recommendation not found")
		}
		return nil, fmt.Errorf("failed to get health recommendation by ID: %w", err)
	}

	// Convert the BSON document to a proto message
	recommendationModel, err := bsonToHealthRecommendation(bsonRecommendation)
	if err != nil {
		return nil, err
	}

	return recommendationModel, nil
}

// UpdateHealthRecommendation updates an existing health recommendation in the database.
func (r *HealthRecommendationRepo) UpdateHealthRecommendation(ctx context.Context, recommendation *health.HealthRecommendation) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(recommendation.Id)
	if err != nil {
		return fmt.Errorf("invalid health recommendation ID: %w", err)
	}

	// Convert the model to a BSON document
	bsonRecommendation := bson.M{
		"user_id":             recommendation.UserId,
		"recommendation_type": recommendation.RecommendationType,
		"description":         recommendation.Description,
		"priority":            recommendation.Priority,
		"updated_at":          time.Now(),
	}

	// Update the document in the collection
	result, err := r.db.Collection("health_recommendations").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bsonRecommendation})
	if err != nil {
		return fmt.Errorf("failed to update health recommendation: %w", err)
	}

	if result.ModifiedCount == 0 {
		return status.Errorf(codes.NotFound, "health recommendation not found")
	}

	return nil
}

// DeleteHealthRecommendation deletes a health recommendation from the database.
func (r *HealthRecommendationRepo) DeleteHealthRecommendation(ctx context.Context, id string) error {
	// Convert the string ID to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid health recommendation ID: %w", err)
	}

	// Delete the document from the collection
	result, err := r.db.Collection("health_recommendations").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete health recommendation: %w", err)
	}

	if result.DeletedCount == 0 {
		return status.Errorf(codes.NotFound, "health recommendation not found")
	}

	return nil
}

// ListHealthRecommendations retrieves all health recommendations for a given user ID, applying filters if provided.
func (r *HealthRecommendationRepo) ListHealthRecommendations(ctx context.Context, req *health.ListHealthRecommendationsRequest) ([]*health.HealthRecommendation, error) {
	// Build the filter query based on the request parameters
	filter := bson.M{}
	if req.UserId != "" {
		filter["user_id"] = req.UserId
	}
	if req.RecommendationType != "" {
		filter["recommendation_type"] = req.RecommendationType
	}
	if req.Priority != 0 {
		filter["priority"] = req.Priority
	}

	// Find the documents based on the filter
	cursor, err := r.db.Collection("health_recommendations").Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list health recommendations: %w", err)
	}
	defer cursor.Close(ctx)

	var healthRecommendations []*health.HealthRecommendation

	// Iterate through the cursor and convert each document to a proto message
	for cursor.Next(ctx) {
		var bsonRecommendation bson.M
		if err := cursor.Decode(&bsonRecommendation); err != nil {
			return nil, fmt.Errorf("failed to decode health recommendation: %w", err)
		}

		recommendationModel, err := bsonToHealthRecommendation(bsonRecommendation)
		if err != nil {
			return nil, err
		}

		healthRecommendations = append(healthRecommendations, recommendationModel)
	}

	return healthRecommendations, nil
}

// bsonToHealthRecommendation converts a BSON document to a health.HealthRecommendation proto message.
func bsonToHealthRecommendation(bsonRecommendation bson.M) (*health.HealthRecommendation, error) {
	recommendationModel := &health.HealthRecommendation{}

	// Handle _id field as string
	if val, ok := bsonRecommendation["_id"].(string); ok {
		recommendationModel.Id = val
	} else if oid, ok := bsonRecommendation["_id"].(primitive.ObjectID); ok {
		recommendationModel.Id = oid.Hex()
	} else {
		return nil, fmt.Errorf("invalid _id type: %T", bsonRecommendation["_id"])
	}

	// Handle potentially nil fields
	if val, ok := bsonRecommendation["user_id"].(string); ok {
		recommendationModel.UserId = val
	}
	if val, ok := bsonRecommendation["recommendation_type"].(string); ok {
		recommendationModel.RecommendationType = val
	}
	if val, ok := bsonRecommendation["description"].(string); ok {
		recommendationModel.Description = val
	}
	if val, ok := bsonRecommendation["priority"].(int32); ok {
		recommendationModel.Priority = val
	}

	// Convert created_at and updated_at fields
	if val, ok := bsonRecommendation["created_at"].(primitive.DateTime); ok {
		recommendationModel.CreatedAt = val.Time().Format(time.RFC3339)
	}
	if val, ok := bsonRecommendation["updated_at"].(primitive.DateTime); ok {
		recommendationModel.UpdatedAt = val.Time().Format(time.RFC3339)
	}

	return recommendationModel, nil
}
