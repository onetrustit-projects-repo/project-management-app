package services

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pm-app/backend/internal/models"
	"github.com/pm-app/backend/internal/repositories"
)

type WorkspaceService struct {
	workspaceRepo *repositories.WorkspaceRepository
	userRepo      *repositories.UserRepository
}

func NewWorkspaceService(workspaceRepo *repositories.WorkspaceRepository, userRepo *repositories.UserRepository) *WorkspaceService {
	return &WorkspaceService{workspaceRepo: workspaceRepo, userRepo: userRepo}
}

func (s *WorkspaceService) Create(ctx context.Context, userID uuid.UUID, req *models.CreateWorkspaceRequest) (*models.Workspace, error) {
	// Generate slug
	slug := generateSlug(req.Name)
	
	// Ensure slug is unique
	counter := 1
	originalSlug := slug
	for {
		existing, _ := s.workspaceRepo.GetBySlug(ctx, slug)
		if existing == nil {
			break
		}
		slug = originalSlug + "-" + string(rune('0'+counter))
		counter++
	}

	workspace := &models.Workspace{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		OwnerID:     userID,
		Settings:    make(map[string]interface{}),
	}

	if err := s.workspaceRepo.Create(ctx, workspace); err != nil {
		return nil, err
	}

	// Add owner as admin member
	member := &models.WorkspaceMember{
		ID:          uuid.New(),
		WorkspaceID: workspace.ID,
		UserID:      userID,
		Role:        "admin",
	}
	s.workspaceRepo.AddMember(ctx, member)

	return workspace, nil
}

func (s *WorkspaceService) Get(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	return s.workspaceRepo.GetByID(ctx, id)
}

func (s *WorkspaceService) List(ctx context.Context, userID uuid.UUID) ([]*models.Workspace, error) {
	return s.workspaceRepo.ListByUser(ctx, userID)
}

func (s *WorkspaceService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateWorkspaceRequest) (*models.Workspace, error) {
	workspace, err := s.workspaceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, errors.New("workspace not found")
	}

	if req.Name != "" {
		workspace.Name = req.Name
	}
	if req.Description != "" {
		workspace.Description = req.Description
	}
	if req.Settings != nil {
		workspace.Settings = req.Settings
	}

	if err := s.workspaceRepo.Update(ctx, workspace); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (s *WorkspaceService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	if !s.workspaceRepo.IsOwner(ctx, id, userID) {
		return errors.New("only owner can delete workspace")
	}
	return s.workspaceRepo.Delete(ctx, id)
}

func (s *WorkspaceService) GetMembers(ctx context.Context, workspaceID uuid.UUID) ([]*models.WorkspaceMember, error) {
	return s.workspaceRepo.GetMembers(ctx, workspaceID)
}

func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	// Trim hyphens
	slug = strings.Trim(slug, "-")
	// Add timestamp to ensure uniqueness
	slug = slug + "-" + time.Now().Format("060102")
	return slug
}
