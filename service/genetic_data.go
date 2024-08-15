package service

import (
	"context"
	"fmt"

	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"github.com/health-analytics-service/health-analytics-service/storage"
	"github.com/health-analytics-service/health-analytics-service/storage/redis"
)

// GeneticDataService implements the health.GeneticDataServiceServer interface.
type GeneticDataService struct {
	storage     storage.StorageI
	redisClient *redis.Client
	health.UnimplementedGeneticDataServiceServer
}

// NewGeneticDataService creates a new GeneticDataService instance.
func NewGeneticDataService(storage storage.StorageI) *GeneticDataService {
	return &GeneticDataService{
		storage: storage,
	}
}

// CreateGeneticData creates a new genetic data record.
func (s *GeneticDataService) CreateGeneticData(ctx context.Context, req *health.GeneticData) (*health.Empty, error) {
	createdID, err := s.storage.GeneticData().CreateGeneticData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create genetic data: %w", err)
	}

	fmt.Println("Created Genetic Data with ID:", createdID) // Optional: Print the created ID

	return &health.Empty{}, nil
}

// GetGeneticData retrieves a genetic data record by its ID.
func (s *GeneticDataService) GetGeneticData(ctx context.Context, req *health.ByIdRequest) (*health.GeneticData, error) {
	data, err := s.storage.GeneticData().GetGeneticData(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get genetic data: %w", err)
	}

	return data, nil
}

// UpdateGeneticData updates an existing genetic data record.
func (s *GeneticDataService) UpdateGeneticData(ctx context.Context, req *health.GeneticData) (*health.Empty, error) {
	err := s.storage.GeneticData().UpdateGeneticData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update genetic data: %w", err)
	}

	return &health.Empty{}, nil
}

// DeleteGeneticData deletes a genetic data record by its ID.
func (s *GeneticDataService) DeleteGeneticData(ctx context.Context, req *health.ByIdRequest) (*health.Empty, error) {
	err := s.storage.GeneticData().DeleteGeneticData(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete genetic data: %w", err)
	}

	return &health.Empty{}, nil
}

// ListGeneticData retrieves a list of genetic data records based on the provided request.
func (s *GeneticDataService) ListGeneticData(ctx context.Context, req *health.ListGeneticDataRequest) (*health.ListGeneticDataResponse, error) {
	data, err := s.storage.GeneticData().ListGeneticData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list genetic data: %w", err)
	}

	return &health.ListGeneticDataResponse{
		GeneticData: data,
	}, nil
}
