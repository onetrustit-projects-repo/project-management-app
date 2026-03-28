# Project Management App - API Specification

## Base URL
```
Production: https://api.pmapp.com/v1
Development: http://localhost:8080/v1
```

## Authentication
All endpoints require JWT Bearer token authentication except:
- `POST /auth/login`
- `POST /auth/register`
- `POST /auth/refresh`

### JWT Structure
```json
{
  "sub": "user-uuid",
  "tenant_id": "tenant-uuid",
  "roles": ["admin", "member"],
  "exp": 1234567890,
  "iat": 1234567890
}
```

---

## Endpoints

### Authentication

#### POST /auth/login
Login with email and password.
```json
// Request
{
  "email": "user@example.com",
  "password": "securepassword"
}

// Response 200
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "avatar_url": "https://...",
    "tenant": { "id": "uuid", "name": "Acme Inc" }
  }
}
```

#### POST /auth/register
Register new user.
```json
// Request
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe",
  "tenant_name": "Acme Inc"
}

// Response 201
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": { ... }
}
```

#### POST /auth/refresh
Refresh access token.
```json
// Request
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}

// Response 200
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900
}
```

---

### Projects

#### GET /projects
List all projects for current tenant.
```json
// Query params: ?page=1&limit=20&status=active
// Response 200
{
  "data": [
    {
      "id": "uuid",
      "name": "Website Redesign",
      "key": "WR",
      "color": "#6366f1",
      "status": "active",
      "member_count": 5,
      "task_count": 42,
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

#### POST /projects
Create new project.
```json
// Request
{
  "name": "Website Redesign",
  "key": "WR",
  "description": "Complete overhaul of company website",
  "color": "#6366f1",
  "start_date": "2024-02-01",
  "end_date": "2024-04-30"
}

// Response 201
{
  "id": "uuid",
  "name": "Website Redesign",
  "key": "WR",
  ...
}
```

#### GET /projects/:id
Get project details.
```json
// Response 200
{
  "id": "uuid",
  "name": "Website Redesign",
  "key": "WR",
  "description": "...",
  "color": "#6366f1",
  "status": "active",
  "start_date": "2024-02-01",
  "end_date": "2024-04-30",
  "members": [...],
  "workspaces": [...],
  "stats": {
    "total_tasks": 42,
    "completed_tasks": 15,
    "in_progress_tasks": 8,
    "overdue_tasks": 3
  }
}
```

#### PATCH /projects/:id
Update project.
```json
// Request
{
  "name": "New Name",
  "status": "archived"
}

// Response 200
{ ... }
```

#### DELETE /projects/:id
Delete project (soft delete).
```json
// Response 204 No Content
```

---

### Workspaces

#### GET /projects/:projectId/workspaces
List workspaces for a project.
```json
// Response 200
{
  "data": [
    {
      "id": "uuid",
      "name": "Main Board",
      "type": "kanban",
      "columns": [
        { "id": "uuid", "name": "To Do", "position": 0, "color": "#ef4444" },
        { "id": "uuid", "name": "In Progress", "position": 1, "color": "#f59e0b" },
        { "id": "uuid", "name": "Done", "position": 2, "color": "#22c55e" }
      ]
    }
  ]
}
```

#### POST /projects/:projectId/workspaces
Create workspace.
```json
// Request
{
  "name": "Main Board",
  "type": "kanban",
  "columns": ["To Do", "In Progress", "Review", "Done"]
}

// Response 201
{ ... }
```

---

### Tasks

#### GET /projects/:projectId/tasks
List tasks with filtering.
```json
// Query params:
// ?status=in_progress&priority=high&assignee=user-uuid&page=1&limit=50

// Response 200
{
  "data": [
    {
      "id": "uuid",
      "task_number": 42,
      "title": "Implement login page",
      "description": "...",
      "priority": "high",
      "status": "in_progress",
      "assignee": { "id": "uuid", "name": "John", "avatar_url": "..." },
      "due_date": "2024-02-15",
      "labels": ["frontend", "auth"],
      "subtask_count": 3,
      "subtask_completed": 1,
      "comment_count": 5,
      "attachment_count": 2,
      "created_at": "2024-01-20T10:00:00Z"
    }
  ],
  "pagination": { ... }
}
```

#### POST /projects/:projectId/tasks
Create task (supports NLP).
```json
// Request (standard)
{
  "title": "Implement login page",
  "description": "Create login form with validation",
  "priority": "high",
  "status": "todo",
  "assignee_id": "user-uuid",
  "column_id": "column-uuid",
  "due_date": "2024-02-15",
  "estimated_hours": 8,
  "labels": ["frontend"]
}

// Request (NLP - AI powered)
{
  "nlp_input": "Create 3 tasks for launching marketing campaign by next Friday",
  "project_id": "project-uuid"
}

// Response 201
{
  "id": "uuid",
  "task_number": 43,
  "title": "Implement login page",
  ...
}
```

#### GET /tasks/:id
Get task details.
```json
// Response 200
{
  "id": "uuid",
  "task_number": 42,
  "title": "Implement login page",
  "description": "...",
  "priority": "high",
  "status": "in_progress",
  "column": { "id": "uuid", "name": "In Progress" },
  "assignee": { "id": "uuid", "name": "John Doe", "avatar_url": "..." },
  "reporter": { "id": "uuid", "name": "Jane Smith" },
  "due_date": "2024-02-15",
  "estimated_hours": 8,
  "actual_hours": 3.5,
  "story_points": 5,
  "progress": 40,
  "labels": ["frontend", "auth"],
  "custom_fields": { "story_type": "feature" },
  "subtasks": [...],
  "relations": [
    { "id": "uuid", "type": "blocks", "task": { "id": "uuid", "title": "..." } }
  ],
  "attachments": [...],
  "comments": [...],
  "activity": [...],
  "created_at": "...",
  "updated_at": "..."
}
```

#### PATCH /tasks/:id
Update task.
```json
// Request
{
  "status": "review",
  "priority": "urgent",
  "assignee_id": "new-user-uuid"
}

// Response 200
{ ... }
```

#### PATCH /tasks/:id/move
Move task (drag & drop).
```json
// Request
{
  "column_id": "new-column-uuid",
  "position": 0,
  "subtasks": [
    { "id": "subtask-uuid", "position": 0 }
  ]
}

// Response 200
{ ... }
```

#### DELETE /tasks/:id
Delete task.
```json
// Response 204
```

---

### Task Relations

#### POST /tasks/:taskId/relations
Add dependency.
```json
// Request
{
  "related_task_id": "task-uuid",
  "type": "blocks"  // blocks, blocked_by, relates_to, duplicates
}

// Response 201
{ ... }
```

---

### Comments

#### GET /tasks/:taskId/comments
List comments.
```json
// Response 200
{
  "data": [
    {
      "id": "uuid",
      "content": "Looks good! Just one small fix needed.",
      "user": { "id": "uuid", "name": "Jane", "avatar_url": "..." },
      "mentions": ["user-uuid-1", "user-uuid-2"],
      "created_at": "2024-01-20T15:30:00Z",
      "updated_at": "2024-01-20T15:30:00Z"
    }
  ]
}
```

#### POST /tasks/:taskId/comments
Add comment.
```json
// Request
{
  "content": "Fixed the issue. @john please review.",
  "mentions": ["john-user-uuid"]
}

// Response 201
{ ... }
```

---

### Users

#### GET /users
List tenant users.
```json
// Response 200
{
  "data": [
    { "id": "uuid", "name": "John Doe", "email": "john@example.com", "avatar_url": "..." }
  ]
}
```

#### GET /users/:id
Get user profile.
```json
// Response 200
{
  "id": "uuid",
  "name": "John Doe",
  "email": "john@example.com",
  "avatar_url": "...",
  "timezone": "America/New_York",
  "roles": ["admin"],
  "stats": {
    "assigned_tasks": 15,
    "completed_tasks": 42,
    "overdue_tasks": 2
  }
}
```

---

### Notifications

#### GET /notifications
List user notifications.
```json
// Query params: ?unread=true&page=1
// Response 200
{
  "data": [
    {
      "id": "uuid",
      "type": "task_assigned",
      "title": "New task assigned",
      "message": "You were assigned to 'Fix login bug'",
      "data": { "task_id": "uuid" },
      "is_read": false,
      "created_at": "2024-01-20T10:00:00Z"
    }
  ],
  "unread_count": 5
}
```

#### PATCH /notifications/:id/read
Mark as read.
```json
// Response 200
{ "is_read": true, "read_at": "..." }
```

#### POST /notifications/read-all
Mark all as read.
```json
// Response 200
{ "updated": 5 }
```

---

### Time Tracking

#### POST /tasks/:taskId/time-entries
Start timer or log time.
```json
// Start timer
{
  "action": "start"
}

// Log time entry
{
  "action": "log",
  "duration": 3600,  // seconds
  "description": "Working on UI"
}

// Response 201
{
  "id": "uuid",
  "start_time": "2024-01-20T10:00:00Z",
  "duration": null,
  "is_running": true
}
```

#### PATCH /time-entries/:id
Update time entry.
```json
// Request
{
  "action": "stop"
}

// Response 200
{
  "duration": 3600,
  "is_running": false
}
```

---

### AI Features

#### POST /ai/chat
Chat with AI assistant.
```json
// Request
{
  "message": "What's the status of the marketing campaign project?",
  "project_id": "project-uuid",
  "conversation_id": "conv-uuid"  // optional, for context
}

// Response 200
{
  "conversation_id": "conv-uuid",
  "message": "The marketing campaign project is currently...",
  "suggestions": [...]
}
```

#### POST /ai/suggest
Get AI suggestions for task.
```json
// Request
{
  "task_id": "task-uuid",
  "action": "breakdown"  // breakdown, prioritize, deadline_risk
}

// Response 200
{
  "suggestions": [
    {
      "type": "breakdown",
      "tasks": [
        { "title": "Subtask 1", "priority": "high" },
        { "title": "Subtask 2", "priority": "medium" }
      ]
    }
  ]
}
```

---

### Webhooks

#### GET /webhooks
List webhooks.
```json
// Response 200
{
  "data": [
    {
      "id": "uuid",
      "name": "Slack Integration",
      "url": "https://hooks.slack.com/...",
      "events": ["task.created", "task.updated", "comment.created"],
      "is_active": true
    }
  ]
}
```

#### POST /webhooks
Create webhook.
```json
// Request
{
  "name": "Slack Integration",
  "url": "https://hooks.slack.com/...",
  "secret": "whsec_...",
  "events": ["task.created", "task.updated"]
}

// Response 201
{ ... }
```

---

### Files

#### POST /upload
Upload file.
```
Content-Type: multipart/form-data
file: <binary>

// Response 201
{
  "id": "uuid",
  "filename": "screenshot.png",
  "url": "/files/uuid/screenshot.png",
  "file_size": 102400
}
```

#### GET /files/:id
Download file.
```
// Response 200 (binary)
```

---

## Error Responses

```json
// 400 Bad Request
{
  "error": "validation_error",
  "message": "Invalid input data",
  "details": [
    { "field": "email", "message": "Invalid email format" }
  ]
}

// 401 Unauthorized
{
  "error": "unauthorized",
  "message": "Invalid or expired token"
}

// 403 Forbidden
{
  "error": "forbidden",
  "message": "You don't have permission to perform this action"
}

// 404 Not Found
{
  "error": "not_found",
  "message": "Resource not found"
}

// 429 Too Many Requests
{
  "error": "rate_limit_exceeded",
  "message": "Too many requests",
  "retry_after": 60
}

// 500 Internal Server Error
{
  "error": "internal_error",
  "message": "An unexpected error occurred"
}
```

---

## Rate Limiting

- **Authenticated:** 1000 requests/minute
- **Unauthenticated:** 20 requests/minute
- **Upload:** 100 requests/hour

Response headers:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1706200000
```

---

## WebSocket (Real-time)

### Connection
```
wss://api.pmapp.com/ws?token=access_token
```

### Events (Server → Client)

```json
// Task created
{
  "type": "task.created",
  "data": { "task": {...} }
}

// Task updated
{
  "type": "task.updated",
  "data": { "task_id": "uuid", "changes": {...} }
}

// Task moved
{
  "type": "task.moved",
  "data": { "task_id": "uuid", "column_id": "uuid", "position": 0 }
}

// Comment added
{
  "type": "comment.created",
  "data": { "task_id": "uuid", "comment": {...} }
}

// Notification
{
  "type": "notification",
  "data": { "notification": {...} }
}
```

### Events (Client → Server)

```json
// Subscribe to project
{ "action": "subscribe", "channel": "project:uuid" }

// Unsubscribe
{ "action": "unsubscribe", "channel": "project:uuid" }

// Typing indicator
{ "action": "typing", "channel": "task:uuid" }
```
