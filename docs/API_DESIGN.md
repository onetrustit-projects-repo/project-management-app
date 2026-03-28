# API Design Specification

## Base URL

```
Production: https://api.yourdomain.com/api/v1
Development: http://localhost:8080/api/v1
```

## Authentication

All authenticated requests require:
```
Authorization: Bearer <access_token>
```

### Endpoints

#### Register
```
POST /auth/register
Content-Type: application/json

Request:
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}

Response (201):
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "550e8400-e29b-41d4-a716-...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "member"
  }
}
```

#### Login
```
POST /auth/login
Content-Type: application/json

Request:
{
  "email": "user@example.com",
  "password": "securepassword"
}

Response (200):
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "550e8400-e29b-41d4-a716-...",
  "user": { ... }
}
```

#### Refresh Token
```
POST /auth/refresh
Content-Type: application/json

Request:
{
  "refresh_token": "550e8400-e29b-41d4-a716-..."
}

Response (200):
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "new-refresh-token...",
  "user": { ... }
}
```

## Workspaces

### List Workspaces
```
GET /workspaces

Response (200):
[
  {
    "id": "uuid",
    "name": "My Workspace",
    "slug": "my-workspace",
    "description": "Workspace description",
    "owner_id": "uuid",
    "settings": {},
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### Create Workspace
```
POST /workspaces
Content-Type: application/json

Request:
{
  "name": "New Workspace",
  "description": "Optional description"
}

Response (201):
{
  "id": "uuid",
  "name": "New Workspace",
  "slug": "new-workspace-240101",
  ...
}
```

### Get Workspace
```
GET /workspaces/:id
```

### Update Workspace
```
PUT /workspaces/:id
Content-Type: application/json

Request:
{
  "name": "Updated Name",
  "description": "Updated description"
}
```

### Delete Workspace
```
DELETE /workspaces/:id
```

### Get Workspace Members
```
GET /workspaces/:id/members

Response (200):
[
  {
    "id": "uuid",
    "workspace_id": "uuid",
    "user_id": "uuid",
    "role": "admin",
    "created_at": "...",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "avatar_url": "..."
    }
  }
]
```

## Projects

### List Projects
```
GET /projects?workspace_id=<uuid>

Response (200):
[
  {
    "id": "uuid",
    "workspace_id": "uuid",
    "name": "Project Name",
    "key": "PN",
    "description": "...",
    "status": "active",
    "created_at": "...",
    "updated_at": "..."
  }
]
```

### Create Project
```
POST /projects?workspace_id=<uuid>
Content-Type: application/json

Request:
{
  "name": "Marketing Campaign",
  "key": "MC",
  "description": "Q1 marketing initiatives"
}

Response (201):
{
  "id": "uuid",
  "key": "MC",
  "name": "Marketing Campaign",
  ...
}
```

### Get Project
```
GET /projects/:key
```

### Update Project
```
PUT /projects/:key
Content-Type: application/json

Request:
{
  "name": "Updated Name",
  "status": "archived"
}
```

### Delete Project
```
DELETE /projects/:key
```

## Tasks

### List Tasks
```
GET /projects/:key/tasks

Response (200):
[
  {
    "id": "uuid",
    "project_id": "uuid",
    "task_number": 1,
    "title": "Setup landing page",
    "description": "Create initial landing page",
    "status": "in_progress",
    "priority": "high",
    "position": 0,
    "assignee_id": "uuid",
    "due_date": "2024-01-15T00:00:00Z",
    "estimated_hours": 8,
    "actual_hours": 3,
    "created_at": "...",
    "updated_at": "...",
    "assignee": {
      "id": "uuid",
      "name": "Jane Doe",
      "avatar_url": "..."
    }
  }
]
```

### Create Task
```
POST /projects/:key/tasks
Content-Type: application/json

Request:
{
  "title": "New Task",
  "description": "Task description",
  "status": "backlog",
  "priority": "medium",
  "assignee_id": "uuid",
  "due_date": "2024-01-20T00:00:00Z",
  "estimated_hours": 4
}

Response (201):
{
  "id": "uuid",
  "task_number": 5,
  "title": "New Task",
  ...
}
```

### Get Task
```
GET /tasks/:id
```

### Update Task
```
PUT /tasks/:id
Content-Type: application/json

Request:
{
  "title": "Updated title",
  "status": "in_progress",
  "priority": "high",
  "assignee_id": "uuid"
}
```

### Delete Task
```
DELETE /tasks/:id
```

### Move Task (Kanban)
```
POST /tasks/:id/move
Content-Type: application/json

Request:
{
  "status": "done",
  "position": 0
}

Response (200):
{
  "id": "uuid",
  "status": "done",
  "position": 0,
  ...
}
```

### Task Comments

#### Add Comment
```
POST /tasks/:id/comments
Content-Type: application/json

Request:
{
  "content": "This task is blocked by the design review."
}

Response (201):
{
  "id": "uuid",
  "task_id": "uuid",
  "user_id": "uuid",
  "content": "This task is blocked by the design review.",
  "created_at": "...",
  "user": { ... }
}
```

#### Get Comments
```
GET /tasks/:id/comments

Response (200):
[
  {
    "id": "uuid",
    "content": "Comment text",
    "created_at": "...",
    "user": { ... }
  }
]
```

### Task Activity
```
GET /tasks/:id/activity

Response (200):
[
  {
    "id": "uuid",
    "action": "created",
    "entity_type": "task",
    "entity_id": "uuid",
    "user_id": "uuid",
    "changes": {},
    "created_at": "..."
  }
]
```

## Dashboard

### Get Stats
```
GET /dashboard/stats?workspace_id=<uuid>

Response (200):
{
  "total_tasks": 45,
  "completed_tasks": 23,
  "in_progress_tasks": 12,
  "overdue_tasks": 3,
  "total_projects": 5,
  "active_projects": 4,
  "total_members": 8,
  "tasks_by_priority": {
    "low": 10,
    "medium": 25,
    "high": 8,
    "urgent": 2
  },
  "tasks_by_status": {
    "backlog": 5,
    "todo": 15,
    "in_progress": 12,
    "in_review": 5,
    "done": 23,
    "cancelled": 0
  }
}
```

### Get Burndown
```
GET /dashboard/projects/:key/burndown

Response (200):
[
  {
    "date": "2024-01-01",
    "remaining": 45,
    "completed": 0,
    "ideal": 45
  },
  {
    "date": "2024-01-02",
    "remaining": 43,
    "completed": 2,
    "ideal": 43.5
  }
]
```

## WebSocket Events

Connect to: `/api/v1/ws`

### Client → Server

```json
{
  "type": "subscribe_project",
  "project_key": "PM"
}
```

```json
{
  "type": "unsubscribe_project",
  "project_key": "PM"
}
```

### Server → Client

```json
{
  "type": "task_created",
  "data": { ...task object... }
}
```

```json
{
  "type": "task_updated",
  "data": { ...task object... }
}
```

```json
{
  "type": "task_moved",
  "data": {
    "task_id": "uuid",
    "from_status": "todo",
    "to_status": "in_progress",
    "position": 0
  }
}
```

```json
{
  "type": "comment_added",
  "data": { ...comment object... }
}
```

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message description"
}
```

Common status codes:
- `400` - Bad Request (validation error)
- `401` - Unauthorized (invalid/expired token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `429` - Too Many Requests (rate limit)
- `500` - Internal Server Error

## Rate Limiting

- 100 requests per minute per IP for general API
- 5 requests per minute for authentication endpoints
- Rate limit headers returned:
  - `X-RateLimit-Limit`
  - `X-RateLimit-Remaining`
  - `X-RateLimit-Reset`
