package models

import "github.com/google/uuid"

// Auth DTOs
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Workspace DTOs
type CreateWorkspaceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Settings    JSONMap `json:"settings"`
}

// Project DTOs
type CreateProjectRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Key         string  `json:"key" binding:"required,min=2,max=10"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type UpdateProjectRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type AddMemberRequest struct {
	UserID uuid.UUID  `json:"user_id" binding:"required"`
	Role   MemberRole `json:"role"`
}

// Task DTOs
type CreateTaskRequest struct {
	Title          string        `json:"title" binding:"required"`
	Description    string        `json:"description"`
	Status         TaskStatus    `json:"status"`
	Priority       TaskPriority  `json:"priority"`
	AssigneeID     *uuid.UUID    `json:"assignee_id"`
	ReporterID     *uuid.UUID    `json:"reporter_id"`
	ParentID       *uuid.UUID    `json:"parent_id"`
	DueDate        *string       `json:"due_date"`
	EstimatedHours *float64      `json:"estimated_hours"`
}

type UpdateTaskRequest struct {
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Status         TaskStatus    `json:"status"`
	Priority       TaskPriority  `json:"priority"`
	AssigneeID     *uuid.UUID    `json:"assignee_id"`
	DueDate        *string       `json:"due_date"`
	EstimatedHours *float64      `json:"estimated_hours"`
	ActualHours    *float64      `json:"actual_hours"`
}

type MoveTaskRequest struct {
	Status   TaskStatus `json:"status" binding:"required"`
	Position int        `json:"position"`
}

type AddCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// Dashboard DTOs
type DashboardStats struct {
	TotalTasks       int `json:"total_tasks"`
	CompletedTasks   int `json:"completed_tasks"`
	InProgressTasks  int `json:"in_progress_tasks"`
	OverdueTasks     int `json:"overdue_tasks"`
	TotalProjects    int `json:"total_projects"`
	ActiveProjects   int `json:"active_projects"`
	TotalMembers     int `json:"total_members"`
	TasksByPriority  map[string]int `json:"tasks_by_priority"`
	TasksByStatus    map[string]int `json:"tasks_by_status"`
}

type BurndownPoint struct {
	Date       string  `json:"date"`
	Remaining  int     `json:"remaining"`
	Completed  int     `json:"completed"`
	Ideal      float64 `json:"ideal"`
}

// WebSocket DTOs
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type TaskMoveEvent struct {
	TaskID     uuid.UUID  `json:"task_id"`
	FromStatus TaskStatus `json:"from_status"`
	ToStatus   TaskStatus `json:"to_status"`
	Position   int        `json:"position"`
}
