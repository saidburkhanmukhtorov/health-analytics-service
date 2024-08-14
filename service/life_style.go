package service

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"github.com/health-analytics-service/health-analytics-service/storage"
)

// LifestyleDataService implements the health.LifestyleDataServiceServer interface.
type LifestyleDataService struct {
	storage     storage.StorageI
	redisClient *redis.Client
	health.UnimplementedLifestyleDataServiceServer
}

// NewLifestyleDataService creates a new LifestyleDataService instance.
func NewLifestyleDataService(storage storage.StorageI) *LifestyleDataService {
	return &LifestyleDataService{
		storage: storage,
	}
}

// // CreateLifestyleData creates a new lifestyle data record.
// func (s *LifestyleDataService) CreateLifestyleData(ctx context.Context, req *health.LifestyleData) (*health.Empty, error) {
// 	createdID, err := s.storage.LifestyleData().CreateLifestyleData(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create lifestyle data: %w", err)
// 	}

// 	fmt.Println("Created Lifestyle Data with ID:", createdID) // Optional: Print the created ID

// 	return &health.Empty{}, nil
// }

// GetLifestyleData retrieves a lifestyle data record by its ID.
func (s *LifestyleDataService) GetLifestyleData(ctx context.Context, req *health.ByIdRequest) (*health.LifestyleData, error) {
	data, err := s.storage.LifestyleData().GetLifestyleData(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get lifestyle data: %w", err)
	}

	return data, nil
}

// // UpdateLifestyleData updates an existing lifestyle data record.
// func (s *LifestyleDataService) UpdateLifestyleData(ctx context.Context, req *health.LifestyleData) (*health.Empty, error) {
// 	err := s.storage.LifestyleData().UpdateLifestyleData(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update lifestyle data: %w", err)
// 	}

// 	return &health.Empty{}, nil
// }

// DeleteLifestyleData deletes a lifestyle data record by its ID.
func (s *LifestyleDataService) DeleteLifestyleData(ctx context.Context, req *health.ByIdRequest) (*health.Empty, error) {
	err := s.storage.LifestyleData().DeleteLifestyleData(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete lifestyle data: %w", err)
	}

	return &health.Empty{}, nil
}

// ListLifestyleData retrieves a list of lifestyle data records based on the provided request.
func (s *LifestyleDataService) ListLifestyleData(ctx context.Context, req *health.ListLifestyleDataRequest) (*health.ListLifestyleDataResponse, error) {
	data, err := s.storage.LifestyleData().ListLifestyleData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list lifestyle data: %w", err)
	}

	return &health.ListLifestyleDataResponse{
		LifestyleData: data,
	}, nil
}
