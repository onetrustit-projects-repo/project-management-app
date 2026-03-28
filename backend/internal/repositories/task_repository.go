package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pm-app/backend/internal/models"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO tasks (id, project_id, title, description, status, priority, assignee_id, reporter_id, parent_id, due_date, estimated_hours, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING task_number, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		task.ID, task.ProjectID, task.Title, task.Description, task.Status, task.Priority,
		task.AssigneeID, task.ReporterID, task.ParentID, task.DueDate, task.EstimatedHours, task.Position,
	).Scan(&task.TaskNumber, &task.CreatedAt, &task.UpdatedAt)
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	query := `
		SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, 
		       t.assignee_id, t.reporter_id, t.parent_id, t.due_date, t.estimated_hours, 
		       t.actual_hours, t.position, t.task_number, t.created_at, t.updated_at,
		       a.id, a.email, a.name, a.avatar_url,
		       r.id, r.email, r.name, r.avatar_url
		FROM tasks t
		LEFT JOIN users a ON t.assignee_id = a.id
		LEFT JOIN users r ON t.reporter_id = r.id
		WHERE t.id = $1
	`
	task := &models.Task{}
	var assigneeID, reporterID, assigneeEmail, assigneeName, assigneeAvatar *string
	var assigneeUUID, reporterUUID *uuid.UUID

	err := r.db.QueryRow(ctx, query, id).Scan(
		&task.ID, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.Priority,
		&assigneeUUID, &reporterUUID, &task.ParentID, &task.DueDate, &task.EstimatedHours,
		&task.ActualHours, &task.Position, &task.TaskNumber, &task.CreatedAt, &task.UpdatedAt,
		&assigneeUUID, &assigneeEmail, &assigneeName, &assigneeAvatar,
		&reporterUUID, &reporterEmail, &reporterName, &reporterAvatar,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse assignee
	if assigneeUUID != nil {
		task.Assignee = &models.User{ID: *assigneeUUID}
		if assigneeEmail != nil {
			task.Assignee.Email = *assigneeEmail
		}
		if assigneeName != nil {
			task.Assignee.Name = *assigneeName
		}
		task.Assignee.AvatarURL = assigneeAvatar
	}

	return task, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks SET 
			title = $2, description = $3, status = $4, priority = $5, 
			assignee_id = $6, due_date = $7, estimated_hours = $8, actual_hours = $9, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, task.ID, task.Title, task.Description, task.Status, task.Priority,
		task.AssigneeID, task.DueDate, task.EstimatedHours, task.ActualHours)
	return err
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.TaskStatus, position int) error {
	query := `UPDATE tasks SET status = $2, position = $3, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, status, position)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *TaskRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*models.Task, error) {
	query := `
		SELECT t.id, t.project_id, t.title, t.description, t.status, t.priority, 
		       t.assignee_id, t.reporter_id, t.parent_id, t.due_date, t.estimated_hours, 
		       t.actual_hours, t.position, t.task_number, t.created_at, t.updated_at
		FROM tasks t
		WHERE t.project_id = $1
		ORDER BY t.status, t.position
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		t := &models.Task{}
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.AssigneeID, &t.ReporterID, &t.ParentID, &t.DueDate, &t.EstimatedHours,
			&t.ActualHours, &t.Position, &t.TaskNumber, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepository) AddComment(ctx context.Context, comment *models.TaskComment) error {
	query := `
		INSERT INTO task_comments (id, task_id, user_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, comment.ID, comment.TaskID, comment.UserID, comment.Content).
		Scan(&comment.CreatedAt, &comment.UpdatedAt)
}

func (r *TaskRepository) GetComments(ctx context.Context, taskID uuid.UUID) ([]*models.TaskComment, error) {
	query := `
		SELECT c.id, c.task_id, c.user_id, c.content, c.created_at, c.updated_at,
		       u.id, u.email, u.name, u.avatar_url
		FROM task_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.task_id = $1
		ORDER BY c.created_at ASC
	`
	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.TaskComment
	for rows.Next() {
		c := &models.TaskComment{User: &models.User{}}
		if err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt,
			&c.User.ID, &c.User.Email, &c.User.Name, &c.User.AvatarURL); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *TaskRepository) LogActivity(ctx context.Context, log *models.ActivityLog) error {
	query := `
		INSERT INTO activity_log (id, workspace_id, project_id, task_id, user_id, action, entity_type, entity_id, changes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at
	`
	return r.db.QueryRow(ctx, query, log.ID, log.WorkspaceID, log.ProjectID, log.TaskID,
		log.UserID, log.Action, log.EntityType, log.EntityID, log.Changes).Scan(&log.CreatedAt)
}

func (r *TaskRepository) GetActivity(ctx context.Context, taskID uuid.UUID) ([]*models.ActivityLog, error) {
	query := `
		SELECT id, workspace_id, project_id, task_id, user_id, action, entity_type, entity_id, changes, created_at
		FROM activity_log
		WHERE task_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.ActivityLog
	for rows.Next() {
		l := &models.ActivityLog{}
		if err := rows.Scan(&l.ID, &l.WorkspaceID, &l.ProjectID, &l.TaskID, &l.UserID,
			&l.Action, &l.EntityType, &l.EntityID, &l.Changes, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *TaskRepository) GetStats(ctx context.Context, workspaceID uuid.UUID) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{
		TasksByPriority: make(map[string]int),
		TasksByStatus:   make(map[string]int),
	}

	// Total and by status
	query := `
		SELECT t.status, COUNT(*) 
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE p.workspace_id = $1
		GROUP BY t.status
	`
	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		stats.TasksByStatus[status] = count
		stats.TotalTasks += count
		if status == string(models.StatusDone) {
			stats.CompletedTasks = count
		} else if status == string(models.StatusInProgress) {
			stats.InProgressTasks = count
		}
	}

	// Overdue
	overdueQuery := `
		SELECT COUNT(*) 
		FROM tasks t
		JOIN projects p ON t.project_id = p.id
		WHERE p.workspace_id = $1 AND t.due_date < NOW() AND t.status NOT IN ('done', 'cancelled')
	`
	r.db.QueryRow(ctx, overdueQuery, workspaceID).Scan(&stats.OverdueTasks)

	// Projects
	projQuery := `SELECT COUNT(*), COUNT(*) FILTER (WHERE status = 'active') FROM projects WHERE workspace_id = $1`
	r.db.QueryRow(ctx, projQuery, workspaceID).Scan(&stats.TotalProjects, &stats.ActiveProjects)

	// Members
	memQuery := `SELECT COUNT(DISTINCT wm.user_id) FROM workspace_members wm WHERE wm.workspace_id = $1`
	r.db.QueryRow(ctx, memQuery, workspaceID).Scan(&stats.TotalMembers)

	return stats, nil
}

func (r *TaskRepository) GetBurndown(ctx context.Context, projectID uuid.UUID, startDate, endDate time.Time) ([]models.BurndownPoint, error) {
	query := `
		WITH dates AS (
			SELECT generate_series($2::date, $3::date, '1 day'::interval) as day
		),
		daily_completed AS (
			SELECT DATE(t.updated_at) as day, COUNT(*) as completed
			FROM tasks
			WHERE project_id = $1 AND status = 'done' AND updated_at >= $2 AND updated_at <= $3
			GROUP BY DATE(t.updated_at)
		)
		SELECT 
			d.day,
			(SELECT COUNT(*) FROM tasks WHERE project_id = $1 AND status != 'done') - 
			COALESCE(SUM(dc.completed) OVER (ORDER BY d.day), 0) as remaining,
			COALESCE(dc.completed, 0) as completed
		FROM dates d
		LEFT JOIN daily_completed dc ON DATE(d.day) = dc.day
		ORDER BY d.day
	`
	rows, err := r.db.Query(ctx, query, projectID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []models.BurndownPoint
	for rows.Next() {
		var p models.BurndownPoint
		var day time.Time
		if err := rows.Scan(&day, &p.Remaining, &p.Completed); err != nil {
			return nil, err
		}
		p.Date = day.Format("2006-01-02")
		points = append(points, p)
	}
	return points, nil
}

var reporterEmail, reporterName, reporterAvatar *string
