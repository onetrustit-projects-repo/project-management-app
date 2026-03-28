package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string
type TaskStatus string
type TaskPriority string
type MemberRole string

const (
	RoleOwner   UserRole = "owner"
	RoleAdmin   UserRole = "admin"
	RoleManager UserRole = "manager"
	RoleMember  UserRole = "member"
	RoleViewer  UserRole = "viewer"
)

const (
	StatusBacklog    TaskStatus = "backlog"
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusInReview   TaskStatus = "in_review"
	StatusDone       TaskStatus = "done"
	StatusCancelled  TaskStatus = "cancelled"
)

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
	PriorityUrgent TaskPriority = "urgent"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Name         string     `json:"name"`
	AvatarURL    *string    `json:"avatar_url"`
	Role         UserRole   `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at"`
}

type Workspace struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Settings    JSONMap   `json:"settings"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WorkspaceMember struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Role        MemberRole `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	User        *User      `json:"user,omitempty"`
}

type Project struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Key         string     `json:"key"`
	Status      string     `json:"status"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	OwnerID     uuid.UUID  `json:"owner_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ProjectMember struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"project_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Role      MemberRole `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	User      *User      `json:"user,omitempty"`
}

type Task struct {
	ID              uuid.UUID     `json:"id"`
	ProjectID       uuid.UUID     `json:"project_id"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Status          TaskStatus    `json:"status"`
	Priority        TaskPriority  `json:"priority"`
	AssigneeID      *uuid.UUID    `json:"assignee_id"`
	ReporterID      *uuid.UUID    `json:"reporter_id"`
	ParentID        *uuid.UUID    `json:"parent_id"`
	DueDate         *time.Time    `json:"due_date"`
	EstimatedHours  *float64      `json:"estimated_hours"`
	ActualHours     float64       `json:"actual_hours"`
	Position        int           `json:"position"`
	TaskNumber      int           `json:"task_number"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	Assignee        *User         `json:"assignee,omitempty"`
	Reporter        *User         `json:"reporter,omitempty"`
	Subtasks        []Task        `json:"subtasks,omitempty"`
	Comments        []TaskComment `json:"comments,omitempty"`
	Dependencies    []TaskDependency `json:"dependencies,omitempty"`
}

type TaskDependency struct {
	ID             uuid.UUID `json:"id"`
	BlockingTaskID uuid.UUID `json:"blocking_task_id"`
	BlockedTaskID  uuid.UUID `json:"blocked_task_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type TaskComment struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user,omitempty"`
}

type TaskAttachment struct {
	ID         uuid.UUID `json:"id"`
	TaskID     uuid.UUID `json:"task_id"`
	Filename   string    `json:"filename"`
	FileURL    string    `json:"file_url"`
	FileSize   int64     `json:"file_size"`
	MimeType   string    `json:"mime_type"`
	UploadedBy uuid.UUID `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type ActivityLog struct {
	ID         uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	ProjectID  uuid.UUID `json:"project_id"`
	TaskID     uuid.UUID `json:"task_id"`
	UserID     uuid.UUID `json:"user_id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	Changes    JSONMap   `json:"changes"`
	CreatedAt  time.Time `json:"created_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	IPAddress    string    `json:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type JSONMap map[string]interface{}
