package service

import (
	"context"
	"fmt"

	"github.com/health-analytics-service/health-analytics-service/genproto/health"
	"github.com/health-analytics-service/health-analytics-service/storage"
)

// HealthMonitoringService implements the health.HealthMonitoringServiceServer interface.
type HealthMonitoringService struct {
	storage storage.StorageI

	health.UnimplementedHealthMonitoringServiceServer
}

// NewHealthMonitoringService creates a new HealthMonitoringService instance.
func NewHealthMonitoringService(storage storage.StorageI) *HealthMonitoringService {
	return &HealthMonitoringService{
		storage: storage,
	}
}

// GetDailySummary retrieves a daily summary of health data for a given user ID and date.
func (s *HealthMonitoringService) GetDailySummary(ctx context.Context, req *health.DailySummaryRequest) (*health.SummaryResponse, error) {
	summary, err := s.storage.HealthMonitoring().GetDailySummary(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily summary: %w", err)
	}

	return summary, nil
}

// GetWeeklySummary retrieves a weekly summary of health data for a given user ID and date range.
func (s *HealthMonitoringService) GetWeeklySummary(ctx context.Context, req *health.WeeklySummaryRequest) (*health.SummaryResponse, error) {
	summary, err := s.storage.HealthMonitoring().GetWeeklySummary(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly summary: %w", err)
	}

	return summary, nil
}
