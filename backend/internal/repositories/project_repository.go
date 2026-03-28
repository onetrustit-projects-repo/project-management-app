package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pm-app/backend/internal/models"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, p *models.Project) error {
	query := `
		INSERT INTO projects (id, workspace_id, name, description, key, status, start_date, end_date, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		p.ID, p.WorkspaceID, p.Name, p.Description, p.Key, p.Status, p.StartDate, p.EndDate, p.OwnerID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	query := `
		SELECT id, workspace_id, name, description, key, status, start_date, end_date, owner_id, created_at, updated_at
		FROM projects WHERE id = $1
	`
	p := &models.Project{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.Key, &p.Status,
		&p.StartDate, &p.EndDate, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return p, err
}

func (r *ProjectRepository) GetByKey(ctx context.Context, key string) (*models.Project, error) {
	query := `
		SELECT id, workspace_id, name, description, key, status, start_date, end_date, owner_id, created_at, updated_at
		FROM projects WHERE key = $1
	`
	p := &models.Project{}
	err := r.db.QueryRow(ctx, query, key).Scan(
		&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.Key, &p.Status,
		&p.StartDate, &p.EndDate, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return p, err
}

func (r *ProjectRepository) Update(ctx context.Context, p *models.Project) error {
	query := `
		UPDATE projects SET name = $2, description = $3, status = $4, start_date = $5, end_date = $6, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, p.ID, p.Name, p.Description, p.Status, p.StartDate, p.EndDate)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *ProjectRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*models.Project, error) {
	query := `
		SELECT id, workspace_id, name, description, key, status, start_date, end_date, owner_id, created_at, updated_at
		FROM projects WHERE workspace_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		p := &models.Project{}
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.Key, &p.Status,
			&p.StartDate, &p.EndDate, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepository) AddMember(ctx context.Context, member *models.ProjectMember) error {
	query := `
		INSERT INTO project_members (id, project_id, user_id, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (project_id, user_id) DO UPDATE SET role = $4
	`
	_, err := r.db.Exec(ctx, query, member.ID, member.ProjectID, member.UserID, member.Role)
	return err
}

func (r *ProjectRepository) GetMembers(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectMember, error) {
	query := `
		SELECT pm.id, pm.project_id, pm.user_id, pm.role, pm.created_at,
		       u.id, u.email, u.name, u.avatar_url, u.role
		FROM project_members pm
		JOIN users u ON pm.user_id = u.id
		WHERE pm.project_id = $1
		ORDER BY pm.created_at
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.ProjectMember
	for rows.Next() {
		m := &models.ProjectMember{User: &models.User{}}
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.CreatedAt,
			&m.User.ID, &m.User.Email, &m.User.Name, &m.User.AvatarURL, &m.User.Role); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *ProjectRepository) RemoveMember(ctx context.Context, projectID, userID uuid.UUID) error {
	query := `DELETE FROM project_members WHERE project_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, projectID, userID)
	return err
}

func (r *ProjectRepository) CountByWorkspace(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM projects WHERE workspace_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, workspaceID).Scan(&count)
	return count, err
}
