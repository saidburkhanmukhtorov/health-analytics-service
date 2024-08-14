package service

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"github.com/health-analytics-service/health-analytics-service/storage"
)

// MedicalRecordService implements the health.MedicalRecordServiceServer interface.
type MedicalRecordService struct {
	storage     storage.StorageI
	redisClient *redis.Client
	health.UnimplementedMedicalRecordServiceServer
}

// NewMedicalRecordService creates a new MedicalRecordService instance.
func NewMedicalRecordService(storage storage.StorageI) *MedicalRecordService {
	return &MedicalRecordService{
		storage: storage,
	}
}

// // CreateMedicalRecord creates a new medical record.
// func (s *MedicalRecordService) CreateMedicalRecord(ctx context.Context, req *health.MedicalRecord) (*health.Empty, error) {
// 	createdID, err := s.storage.MedicalRecord().CreateMedicalRecord(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create medical record: %w", err)
// 	}

// 	fmt.Println("Created Medical Record with ID:", createdID) // Optional: Print the created ID

// 	return &health.Empty{}, nil
// }

// GetMedicalRecord retrieves a medical record by its ID.
func (s *MedicalRecordService) GetMedicalRecord(ctx context.Context, req *health.ByIdRequest) (*health.MedicalRecord, error) {
	record, err := s.storage.MedicalRecord().GetMedicalRecord(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical record: %w", err)
	}

	return record, nil
}

// // UpdateMedicalRecord updates an existing medical record.
// func (s *MedicalRecordService) UpdateMedicalRecord(ctx context.Context, req *health.MedicalRecord) (*health.Empty, error) {
// 	err := s.storage.MedicalRecord().UpdateMedicalRecord(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update medical record: %w", err)
// 	}

// 	return &health.Empty{}, nil
// }

// DeleteMedicalRecord deletes a medical record by its ID.
func (s *MedicalRecordService) DeleteMedicalRecord(ctx context.Context, req *health.ByIdRequest) (*health.Empty, error) {
	err := s.storage.MedicalRecord().DeleteMedicalRecord(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete medical record: %w", err)
	}

	return &health.Empty{}, nil
}

// ListMedicalRecords retrieves a list of medical records based on the provided request.
func (s *MedicalRecordService) ListMedicalRecords(ctx context.Context, req *health.ListMedicalRecordsRequest) (*health.ListMedicalRecordsResponse, error) {
	records, err := s.storage.MedicalRecord().ListMedicalRecords(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list medical records: %w", err)
	}

	return &health.ListMedicalRecordsResponse{
		MedicalRecords: records,
	}, nil
}
