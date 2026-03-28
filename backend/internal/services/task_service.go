package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/pm-app/backend/internal/models"
	"github.com/pm-app/backend/internal/repositories"
)

type TaskService struct {
	taskRepo    *repositories.TaskRepository
	projectRepo *repositories.ProjectRepository
	redis       *repositories.RedisRepository
}

func NewTaskService(taskRepo *repositories.TaskRepository, projectRepo *repositories.ProjectRepository, redisClient interface{}) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (s *TaskService) Create(ctx context.Context, projectKey string, userID uuid.UUID, req *models.CreateTaskRequest) (*models.Task, error) {
	project, err := s.projectRepo.GetByKey(ctx, projectKey)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	task := &models.Task{
		ID:           uuid.New(),
		ProjectID:    project.ID,
		Title:        req.Title,
		Description:  req.Description,
		Status:       models.StatusBacklog,
		Priority:     models.PriorityMedium,
		ReporterID:   &userID,
		ParentID:     req.ParentID,
		Position:     0,
	}

	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.ReporterID != nil {
		task.ReporterID = req.ReporterID
	}
	if req.DueDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.DueDate)
		task.DueDate = &t
	}
	if req.EstimatedHours != nil {
		task.EstimatedHours = req.EstimatedHours
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	// Log activity
	s.logActivity(ctx, project.WorkspaceID, project.ID, task.ID, userID, "created", "task")

	return task, nil
}

func (s *TaskService) Get(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *TaskService) ListByProject(ctx context.Context, projectKey string) ([]*models.Task, error) {
	project, err := s.projectRepo.GetByKey(ctx, projectKey)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	return s.taskRepo.ListByProject(ctx, project.ID)
}

func (s *TaskService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.UpdateTaskRequest) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.DueDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.DueDate)
		task.DueDate = &t
	}
	if req.EstimatedHours != nil {
		task.EstimatedHours = req.EstimatedHours
	}
	if req.ActualHours != nil {
		task.ActualHours = *req.ActualHours
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	// Log activity
	project, _ := s.projectRepo.GetByID(ctx, task.ProjectID)
	if project != nil {
		s.logActivity(ctx, project.WorkspaceID, project.ID, task.ID, userID, "updated", "task")
	}

	return task, nil
}

func (s *TaskService) Move(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *models.MoveTaskRequest) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	oldStatus := task.Status
	task.Status = req.Status
	task.Position = req.Position

	if err := s.taskRepo.UpdateStatus(ctx, id, req.Status, req.Position); err != nil {
		return nil, err
	}

	// Log activity
	project, _ := s.projectRepo.GetByID(ctx, task.ProjectID)
	if project != nil {
		s.logActivity(ctx, project.WorkspaceID, project.ID, task.ID, userID, "moved", "task")
		_ = oldStatus // Could track status change in activity
	}

	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *TaskService) AddComment(ctx context.Context, taskID, userID uuid.UUID, req *models.AddCommentRequest) (*models.TaskComment, error) {
	comment := &models.TaskComment{
		ID:      uuid.New(),
		TaskID:  taskID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.taskRepo.AddComment(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *TaskService) GetComments(ctx context.Context, taskID uuid.UUID) ([]*models.TaskComment, error) {
	return s.taskRepo.GetComments(ctx, taskID)
}

func (s *TaskService) GetActivity(ctx context.Context, taskID uuid.UUID) ([]*models.ActivityLog, error) {
	return s.taskRepo.GetActivity(ctx, taskID)
}

func (s *TaskService) logActivity(ctx context.Context, workspaceID, projectID, taskID, userID uuid.UUID, action, entityType string) {
	log := &models.ActivityLog{
		ID:         uuid.New(),
		WorkspaceID: workspaceID,
		ProjectID:  projectID,
		TaskID:     taskID,
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   taskID,
	}
	s.taskRepo.LogActivity(ctx, log)
}
