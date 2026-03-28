package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pm-app/backend/internal/models"
	"github.com/pm-app/backend/internal/repositories"
)

type ProjectService struct {
	projectRepo   *repositories.ProjectRepository
	workspaceRepo *repositories.WorkspaceRepository
}

func NewProjectService(projectRepo *repositories.ProjectRepository, workspaceRepo *repositories.WorkspaceRepository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo, workspaceRepo: workspaceRepo}
}

func (s *ProjectService) Create(ctx context.Context, workspaceID, userID uuid.UUID, req *models.CreateProjectRequest) (*models.Project, error) {
	// Verify workspace access
	ws, err := s.workspaceRepo.GetByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if ws == nil {
		return nil, errors.New("workspace not found")
	}

	project := &models.Project{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		Name:        req.Name,
		Description: req.Description,
		Key:         req.Key,
		Status:      "active",
		OwnerID:     userID,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	// Add creator as admin member
	member := &models.ProjectMember{
		ID:        uuid.New(),
		ProjectID: project.ID,
		UserID:    userID,
		Role:      "admin",
	}
	s.projectRepo.AddMember(ctx, member)

	return project, nil
}

func (s *ProjectService) GetByKey(ctx context.Context, key string) (*models.Project, error) {
	return s.projectRepo.GetByKey(ctx, key)
}

func (s *ProjectService) List(ctx context.Context, workspaceID uuid.UUID) ([]*models.Project, error) {
	return s.projectRepo.ListByWorkspace(ctx, workspaceID)
}

func (s *ProjectService) Update(ctx context.Context, key string, req *models.UpdateProjectRequest) (*models.Project, error) {
	project, err := s.projectRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Status != "" {
		project.Status = req.Status
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) Delete(ctx context.Context, key string) error {
	project, err := s.projectRepo.GetByKey(ctx, key)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	return s.projectRepo.Delete(ctx, project.ID)
}

func (s *ProjectService) GetMembers(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectMember, error) {
	return s.projectRepo.GetMembers(ctx, projectID)
}

func (s *ProjectService) AddMember(ctx context.Context, projectID uuid.UUID, req *models.AddMemberRequest) error {
	member := &models.ProjectMember{
		ID:        uuid.New(),
		ProjectID: projectID,
		UserID:    req.UserID,
		Role:      req.Role,
	}
	if member.Role == "" {
		member.Role = "member"
	}
	return s.projectRepo.AddMember(ctx, member)
}
