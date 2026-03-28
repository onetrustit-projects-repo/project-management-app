package repositories

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(connString string) (*pgxpool.Pool, error) {
	if connString == "" {
		connString = "postgres://postgres:postgres@localhost:5432/pm_app?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func RunMigrations(db *pgxpool.Pool) error {
	ctx := context.Background()

	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,
		
		`CREATE TYPE IF NOT EXISTS user_role AS ENUM ('owner', 'admin', 'manager', 'member', 'viewer')`,
		`CREATE TYPE IF NOT EXISTS task_status AS ENUM ('backlog', 'todo', 'in_progress', 'in_review', 'done', 'cancelled')`,
		`CREATE TYPE IF NOT EXISTS task_priority AS ENUM ('low', 'medium', 'high', 'urgent')`,
		`CREATE TYPE IF NOT EXISTS member_role AS ENUM ('admin', 'manager', 'member', 'viewer')`,

		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			avatar_url TEXT,
			role user_role DEFAULT 'member',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			last_login_at TIMESTAMPTZ
		)`,

		`CREATE TABLE IF NOT EXISTS workspaces (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(100) UNIQUE NOT NULL,
			description TEXT,
			owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
			settings JSONB DEFAULT '{}',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS workspace_members (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role member_role DEFAULT 'member',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(workspace_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS projects (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			key VARCHAR(10) UNIQUE NOT NULL,
			status VARCHAR(50) DEFAULT 'active',
			start_date DATE,
			end_date DATE,
			owner_id UUID NOT NULL REFERENCES users(id),
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS project_members (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role member_role DEFAULT 'member',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(project_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			title VARCHAR(500) NOT NULL,
			description TEXT,
			status task_status DEFAULT 'backlog',
			priority task_priority DEFAULT 'medium',
			assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
			reporter_id UUID REFERENCES users(id) ON DELETE SET NULL,
			parent_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
			due_date TIMESTAMPTZ,
			estimated_hours DECIMAL(10,2),
			actual_hours DECIMAL(10,2) DEFAULT 0,
			position INTEGER DEFAULT 0,
			task_number SERIAL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS task_dependencies (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			blocking_task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			blocked_task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(blocking_task_id, blocked_task_id)
		)`,

		`CREATE TABLE IF NOT EXISTS task_comments (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS activity_log (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
			project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
			task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id) ON DELETE SET NULL,
			action VARCHAR(100) NOT NULL,
			entity_type VARCHAR(50) NOT NULL,
			entity_id UUID,
			changes JSONB,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS sessions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			refresh_token VARCHAR(500) NOT NULL,
			user_agent TEXT,
			ip_address INET,
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
		)`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_assignee ON tasks(assignee_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_workspace ON projects(workspace_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workspace_members_workspace ON workspace_members(workspace_id)`,
		`CREATE INDEX IF NOT EXISTS idx_project_members_project ON project_members(project_id)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(ctx, migration); err != nil {
			return fmt.Errorf("migration failed: %s: %w", migration[:50], err)
		}
	}

	return nil
}
