package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pm-app/backend/internal/models"
	"github.com/pm-app/backend/internal/repositories"
)

type DashboardService struct {
	taskRepo    *repositories.TaskRepository
	projectRepo *repositories.ProjectRepository
}

func NewDashboardService(taskRepo *repositories.TaskRepository, projectRepo *repositories.ProjectRepository) *DashboardService {
	return &DashboardService{taskRepo: taskRepo, projectRepo: projectRepo}
}

func (s *DashboardService) GetStats(ctx context.Context, workspaceID uuid.UUID) (*models.DashboardStats, error) {
	return s.taskRepo.GetStats(ctx, workspaceID)
}

func (s *DashboardService) GetBurndown(ctx context.Context, projectKey string) ([]models.BurndownPoint, error) {
	project, err := s.projectRepo.GetByKey(ctx, projectKey)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, nil
	}

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	return s.taskRepo.GetBurndown(ctx, project.ID, startDate, endDate)
}
