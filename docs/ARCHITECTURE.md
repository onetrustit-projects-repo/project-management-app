# Project Management App - Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENTS                                  │
│   Web Browser (Next.js)  │  Mobile App  │  Third-party APIs   │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NGINX REVERSE PROXY                         │
│              (SSL Termination, Rate Limiting)                    │
└─────────────────────────────────────────────────────────────────┘
                                │
            ┌───────────────────┼───────────────────┐
            ▼                   ▼                   ▼
    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
    │   Frontend   │    │   Backend   │    │  WebSocket   │
    │  (Next.js)   │    │    (Go)     │    │   Server     │
    │   :3000      │    │    :8080    │    │    :8081     │
    └──────────────┘    └──────────────┘    └──────────────┘
            │                   │                   │
            └───────────────────┼───────────────────┘
                                │
            ┌───────────────────┼───────────────────┐
            ▼                   ▼                   ▼
    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
    │  PostgreSQL  │    │    Redis    │    │   MinIO/S3  │
    │    :5432     │    │    :6379    │    │    :9000    │
    └──────────────┘    └──────────────┘    └──────────────┘
                                │
                    ┌───────────┴───────────┐
                    ▼                       ▼
            ┌──────────────┐        ┌──────────────┐
            │ Background   │        │   Cache &    │
            │    Jobs      │        │   Sessions   │
            │  (Workers)   │        │              │
            └──────────────┘        └──────────────┘
```

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Frontend | Next.js 14 (App Router) | React UI framework |
| UI Components | Tailwind CSS + shadcn/ui | Modern component library |
| Drag & Drop | dnd-kit | Kanban drag-and-drop |
| State | Zustand + React Query | Client state management |
| Backend | Go 1.21+ (Gin) | REST API server |
| Database | PostgreSQL 16 | Primary data store |
| Cache | Redis 7 | Sessions, caching, real-time |
| WebSocket | gorilla/websocket | Real-time updates |
| Auth | JWT + Refresh tokens | Authentication |
| File Storage | MinIO (S3-compatible) | File attachments |
| Container | Docker + Docker Compose | Containerization |
| Monitoring | Prometheus + Grafana | Observability |

## Database Schema

### Core Entities

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    Workspace    │────<│     Project     │────<│      Task       │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ id (UUID)       │     │ id (UUID)       │     │ id (UUID)       │
│ name            │     │ workspace_id FK  │     │ project_id FK   │
│ slug            │     │ name            │     │ title           │
│ description     │     │ description      │     │ description     │
│ owner_id FK     │     │ key (PM-1)      │     │ status          │
│ settings JSONB   │     │ status          │     │ priority        │
│ created_at      │     │ start_date       │     │ assignee_id FK  │
│ updated_at      │     │ end_date         │     │ reporter_id FK  │
└─────────────────┘     │ created_at       │     │ due_date        │
        │               │ updated_at       │     │ estimated_hours │
        │               └─────────────────┘     │ actual_hours    │
        │                       │               │ created_at      │
        ▼                       │               │ updated_at      │
┌─────────────────┐             │               └─────────────────┘
│      User       │             │                       │
├─────────────────┤             │                       │
│ id (UUID)       │             │                       ▼
│ email           │             │               ┌─────────────────┐
│ password_hash   │             │               │   TaskDependency│
│ name            │             │               ├─────────────────┤
│ avatar_url      │             │               │ id (UUID)       │
│ role            │             │               │ blocking_id FK  │
│ created_at      │             │               │ blocked_id FK   │
│ updated_at      │             │               └─────────────────┘
└─────────────────┘             │
        │                       ▼
        │               ┌─────────────────┐
        ▼               │   TaskComment   │
┌─────────────────┐     ├─────────────────┤
│  ProjectMember  │     │ id (UUID)       │
├─────────────────┤     │ task_id FK      │
│ id (UUID)       │     │ user_id FK      │
│ project_id FK   │     │ content         │
│ user_id FK      │     │ created_at      │
│ role            │     │ updated_at      │
│ created_at      │     └─────────────────┘
└─────────────────┘
        │
        ▼
┌─────────────────┐
│  TaskAttachment │
├─────────────────┤
│ id (UUID)       │
│ task_id FK      │
│ filename        │
│ file_url        │
│ file_size       │
│ uploaded_by FK  │
│ created_at      │
└─────────────────┘
```

### Indexes

- `tasks`: (project_id, status), (assignee_id), (due_date)
- `projects`: (workspace_id), (owner_id)
- `comments`: (task_id), (user_id)
- `activities`: (task_id), (user_id), (created_at)

## API Design

### REST Endpoints

```
Authentication:
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
GET    /api/v1/auth/me

Workspaces:
GET    /api/v1/workspaces
POST   /api/v1/workspaces
GET    /api/v1/workspaces/:id
PUT    /api/v1/workspaces/:id
DELETE /api/v1/workspaces/:id
GET    /api/v1/workspaces/:id/members

Projects:
GET    /api/v1/projects
POST   /api/v1/projects
GET    /api/v1/projects/:key
PUT    /api/v1/projects/:key
DELETE /api/v1/projects/:key
GET    /api/v1/projects/:key/members

Tasks:
GET    /api/v1/projects/:key/tasks
POST   /api/v1/projects/:key/tasks
GET    /api/v1/tasks/:id
PUT    /api/v1/tasks/:id
DELETE /api/v1/tasks/:id
POST   /api/v1/tasks/:id/move        (Kanban move)
POST   /api/v1/tasks/:id/comments
GET    /api/v1/tasks/:id/comments
GET    /api/v1/tasks/:id/activity

Users:
GET    /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id

Dashboard:
GET    /api/v1/dashboard/stats
GET    /api/v1/dashboard/projects/:key/burndown
```

### WebSocket Events

```
Client → Server:
  subscribe_project    { project_key }
  unsubscribe_project  { project_key }
  task_update          { task_id, changes }
  task_move            { task_id, from_status, to_status, position }

Server → Client:
  task_created         { task }
  task_updated         { task }
  task_moved           { task_id, from_status, to_status, position }
  comment_added        { comment }
  member_joined        { user }
  member_left          { user }
```

## Project Structure

### Backend (Go)

```
backend/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── models/                 # Database models
│   │   ├── user.go
│   │   ├── workspace.go
│   │   ├── project.go
│   │   ├── task.go
│   │   └── comment.go
│   ├── handlers/               # HTTP handlers
│   │   ├── auth.go
│   │   ├── workspace.go
│   │   ├── project.go
│   │   ├── task.go
│   │   └── dashboard.go
│   ├── middleware/             # Auth, logging, cors
│   ├── repositories/          # Database queries
│   ├── services/              # Business logic
│   └── websocket/             # WebSocket hub
├── migrations/                 # SQL migrations
└── pkg/utils/                 # Helpers
```

### Frontend (Next.js)

```
frontend/
├── app/                       # Next.js App Router
│   ├── (auth)/               # Auth pages
│   ├── (dashboard)/           # Main app pages
│   │   ├── projects/
│   │   ├── tasks/
│   │   └── dashboard/
│   └── api/                   # API client
├── components/
│   ├── ui/                    # Base UI components
│   ├── kanban/                # Kanban board
│   ├── calendar/              # Calendar view
│   ├── timeline/              # Gantt chart
│   └── common/                # Shared components
├── hooks/                     # Custom React hooks
├── lib/                       # Utilities
└── types/                     # TypeScript types
```

## Security

### RBAC Roles

| Role | Permissions |
|------|-------------|
| Owner | Full access, delete workspace, manage billing |
| Admin | Manage members, projects, settings |
| Manager | Create/edit projects, assign tasks |
| Member | Create/edit own tasks, comment |
| Viewer | Read-only access |

### Data Isolation

- All queries include `workspace_id` filter
- Row-level security in PostgreSQL
- JWT contains `workspace_id` claim

## Scalability

- Horizontal scaling: Multiple Go instances behind load balancer
- Database: Read replicas for queries
- Redis: Pub/sub for WebSocket scaling
- File storage: MinIO distributed mode
