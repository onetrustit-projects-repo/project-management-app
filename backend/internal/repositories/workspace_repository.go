package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pm-app/backend/internal/models"
)

type WorkspaceRepository struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepository(db *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) Create(ctx context.Context, ws *models.Workspace) error {
	query := `
		INSERT INTO workspaces (id, name, slug, description, owner_id, settings)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		ws.ID, ws.Name, ws.Slug, ws.Description, ws.OwnerID, ws.Settings,
	).Scan(&ws.CreatedAt, &ws.UpdatedAt)
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	query := `
		SELECT id, name, slug, description, owner_id, settings, created_at, updated_at
		FROM workspaces WHERE id = $1
	`
	ws := &models.Workspace{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ws.ID, &ws.Name, &ws.Slug, &ws.Description, &ws.OwnerID, &ws.Settings, &ws.CreatedAt, &ws.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return ws, err
}

func (r *WorkspaceRepository) GetBySlug(ctx context.Context, slug string) (*models.Workspace, error) {
	query := `
		SELECT id, name, slug, description, owner_id, settings, created_at, updated_at
		FROM workspaces WHERE slug = $1
	`
	ws := &models.Workspace{}
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&ws.ID, &ws.Name, &ws.Slug, &ws.Description, &ws.OwnerID, &ws.Settings, &ws.CreatedAt, &ws.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return ws, err
}

func (r *WorkspaceRepository) Update(ctx context.Context, ws *models.Workspace) error {
	query := `
		UPDATE workspaces SET name = $2, description = $3, settings = $4, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, ws.ID, ws.Name, ws.Description, ws.Settings)
	return err
}

func (r *WorkspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM workspaces WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *WorkspaceRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*models.Workspace, error) {
	query := `
		SELECT w.id, w.name, w.slug, w.description, w.owner_id, w.settings, w.created_at, w.updated_at
		FROM workspaces w
		LEFT JOIN workspace_members wm ON w.id = wm.workspace_id
		WHERE w.owner_id = $1 OR wm.user_id = $1
		ORDER BY w.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*models.Workspace
	for rows.Next() {
		ws := &models.Workspace{}
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.Slug, &ws.Description, &ws.OwnerID, &ws.Settings, &ws.CreatedAt, &ws.UpdatedAt); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) AddMember(ctx context.Context, member *models.WorkspaceMember) error {
	query := `
		INSERT INTO workspace_members (id, workspace_id, user_id, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (workspace_id, user_id) DO UPDATE SET role = $4
	`
	_, err := r.db.Exec(ctx, query, member.ID, member.WorkspaceID, member.UserID, member.Role)
	return err
}

func (r *WorkspaceRepository) GetMembers(ctx context.Context, workspaceID uuid.UUID) ([]*models.WorkspaceMember, error) {
	query := `
		SELECT wm.id, wm.workspace_id, wm.user_id, wm.role, wm.created_at,
		       u.id, u.email, u.name, u.avatar_url, u.role
		FROM workspace_members wm
		JOIN users u ON wm.user_id = u.id
		WHERE wm.workspace_id = $1
		ORDER BY wm.created_at
	`
	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.WorkspaceMember
	for rows.Next() {
		m := &models.WorkspaceMember{User: &models.User{}}
		if err := rows.Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.CreatedAt,
			&m.User.ID, &m.User.Email, &m.User.Name, &m.User.AvatarURL, &m.User.Role); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *WorkspaceRepository) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	query := `DELETE FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, workspaceID, userID)
	return err
}

func (r *WorkspaceRepository) IsOwner(ctx context.Context, workspaceID, userID uuid.UUID) bool {
	query := `SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = $1 AND owner_id = $2)`
	var exists bool
	r.db.QueryRow(ctx, query, workspaceID, userID).Scan(&exists)
	return exists
}
