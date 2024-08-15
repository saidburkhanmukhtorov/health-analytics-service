package service

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"github.com/health-analytics-service/health-analytics-service/storage"
)

// HealthRecommendationService implements the health.HealthRecommendationServiceServer interface.
type HealthRecommendationService struct {
	storage     storage.StorageI
	redisClient *redis.Client
	health.UnimplementedHealthRecommendationServiceServer
}

// NewHealthRecommendationService creates a new HealthRecommendationService instance.
func NewHealthRecommendationService(storage storage.StorageI) *HealthRecommendationService {
	return &HealthRecommendationService{
		storage: storage,
	}
}

// CreateHealthRecommendation creates a new health recommendation.
func (s *HealthRecommendationService) CreateHealthRecommendation(ctx context.Context, req *health.HealthRecommendation) (*health.Empty, error) {
	createdID, err := s.storage.HealthRecommendation().CreateHealthRecommendation(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create health recommendation: %w", err)
	}

	fmt.Println("Created Health Recommendation with ID:", createdID) // Optional: Print the created ID

	return &health.Empty{}, nil
}

// GetHealthRecommendation retrieves a health recommendation by its ID.
func (s *HealthRecommendationService) GetHealthRecommendation(ctx context.Context, req *health.ByIdRequest) (*health.HealthRecommendation, error) {
	recommendation, err := s.storage.HealthRecommendation().GetHealthRecommendation(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get health recommendation: %w", err)
	}

	return recommendation, nil
}

// UpdateHealthRecommendation updates an existing health recommendation.
func (s *HealthRecommendationService) UpdateHealthRecommendation(ctx context.Context, req *health.HealthRecommendation) (*health.Empty, error) {
	err := s.storage.HealthRecommendation().UpdateHealthRecommendation(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update health recommendation: %w", err)
	}

	return &health.Empty{}, nil
}

// DeleteHealthRecommendation deletes a health recommendation by its ID.
func (s *HealthRecommendationService) DeleteHealthRecommendation(ctx context.Context, req *health.ByIdRequest) (*health.Empty, error) {
	err := s.storage.HealthRecommendation().DeleteHealthRecommendation(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete health recommendation: %w", err)
	}

	return &health.Empty{}, nil
}

// ListHealthRecommendations retrieves a list of health recommendations based on the provided request.
func (s *HealthRecommendationService) ListHealthRecommendations(ctx context.Context, req *health.ListHealthRecommendationsRequest) (*health.ListHealthRecommendationsResponse, error) {
	recommendations, err := s.storage.HealthRecommendation().ListHealthRecommendations(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list health recommendations: %w", err)
	}

	return &health.ListHealthRecommendationsResponse{
		HealthRecommendations: recommendations,
	}, nil
}
