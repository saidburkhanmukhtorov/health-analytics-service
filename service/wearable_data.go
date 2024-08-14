package service

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"github.com/health-analytics-service/health-analytics-service/storage"
)

// WearableDataService implements the health.WearableDataServiceServer interface.
type WearableDataService struct {
	storage     storage.StorageI
	redisClient *redis.Client
	health.UnimplementedWearableDataServiceServer
}

// NewWearableDataService creates a new WearableDataService instance.
func NewWearableDataService(storage storage.StorageI) *WearableDataService {
	return &WearableDataService{
		storage: storage,
	}
}

// // CreateWearableData creates a new wearable data record.
// func (s *WearableDataService) CreateWearableData(ctx context.Context, req *health.WearableData) (*health.Empty, error) {
// 	createdID, err := s.storage.WearableData().CreateWearableData(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create wearable data: %w", err)
// 	}

// 	fmt.Println("Created Wearable Data with ID:", createdID) // Optional: Print the created ID

// 	return &health.Empty{}, nil
// }

// GetWearableData retrieves a wearable data record by its ID.
func (s *WearableDataService) GetWearableData(ctx context.Context, req *health.ByIdRequest) (*health.WearableData, error) {
	data, err := s.storage.WearableData().GetWearableData(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get wearable data: %w", err)
	}

	return data, nil
}

// // UpdateWearableData updates an existing wearable data record.
// func (s *WearableDataService) UpdateWearableData(ctx context.Context, req *health.WearableData) (*health.Empty, error) {
// 	err := s.storage.WearableData().UpdateWearableData(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update wearable data: %w", err)
// 	}

// 	return &health.Empty{}, nil
// }

// DeleteWearableData deletes a wearable data record by its ID.
func (s *WearableDataService) DeleteWearableData(ctx context.Context, req *health.ByIdRequest) (*health.Empty, error) {
	err := s.storage.WearableData().DeleteWearableData(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete wearable data: %w", err)
	}

	return &health.Empty{}, nil
}

// ListWearableData retrieves a list of wearable data records based on the provided request.
func (s *WearableDataService) ListWearableData(ctx context.Context, req *health.ListWearableDataRequest) (*health.ListWearableDataResponse, error) {
	data, err := s.storage.WearableData().ListWearableData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list wearable data: %w", err)
	}

	return &health.ListWearableDataResponse{
		WearableData: data,
	}, nil
}
