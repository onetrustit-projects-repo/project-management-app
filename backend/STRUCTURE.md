# Backend Folder Structure (Go)

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── config/
│   │   ├── config.go              # Configuration loading
│   │   ├── config_test.go
│   │   └── .env.example
│   │
│   ├── domain/
│   │   ├── user.go                # User entity
│   │   ├── project.go             # Project entity
│   │   ├── task.go                # Task entity
│   │   ├── workspace.go           # Workspace entity
│   │   ├── comment.go             # Comment entity
│   │   ├── notification.go        # Notification entity
│   │   └── webhook.go             # Webhook entity
│   │
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── user_repo.go
│   │   │   ├── project_repo.go
│   │   │   ├── task_repo.go
│   │   │   ├── workspace_repo.go
│   │   │   └── ...
│   │   ├── redis/
│   │   │   ├── cache.go           # Redis caching
│   │   │   └── queue.go           # Redis queues
│   │   └── interfaces.go          # Repository interfaces
│   │
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── project_service.go
│   │   ├── task_service.go
│   │   ├── notification_service.go
│   │   ├── ai_service.go           # AI/LLM integration
│   │   └── automation_service.go  # Workflow automation
│   │
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── project_handler.go
│   │   ├── task_handler.go
│   │   ├── workspace_handler.go
│   │   ├── comment_handler.go
│   │   ├── notification_handler.go
│   │   ├── file_handler.go
│   │   ├── webhook_handler.go
│   │   └── ai_handler.go
│   │
│   ├── middleware/
│   │   ├── auth.go                # JWT authentication
│   │   ├── rbac.go                # Role-based access control
│   │   ├── tenant.go              # Multi-tenant isolation
│   │   ├── ratelimit.go           # Rate limiting
│   │   ├── logging.go             # Request logging
│   │   ├── cors.go                # CORS handling
│   │   └── validation.go          # Request validation
│   │
│   ├── websocket/
│   │   ├── hub.go                 # WebSocket hub
│   │   ├── client.go              # WebSocket client
│   │   └── handler.go             # WebSocket handler
│   │
│   ├── worker/
│   │   ├── notification_worker.go
│   │   ├── email_worker.go
│   │   ├── ai_worker.go           # AI task processing
│   │   ├── webhook_worker.go
│   │   └── worker.go              # Worker pool manager
│   │
│   └── pkg/
│       ├── jwt/
│       │   ├── jwt.go             # JWT utilities
│       │   └── refresh.go          # Refresh token handling
│       ├── password/
│       │   └── password.go        # Password hashing (bcrypt)
│       ├── response/
│       │   └── response.go        # Standardized responses
│       ├── validator/
│       │   └── validator.go       # Custom validators
│       └── minio/
│           └── client.go          # S3/MinIO client
│
├── migrations/
│   ├── 001_initial_schema.sql
│   ├── 002_add_ai_tables.sql
│   └── ...
│
├── api/
│   └── v1/
│       ├── router.go              # API routes
│       ├── middleware.go          # API middleware
│       └── routes/                 # Route definitions
│           ├── auth.go
│           ├── users.go
│           ├── projects.go
│           └── tasks.go
│
├── scripts/
│   ├── generate_token.go          # Utility for generating tokens
│   └── seed_data.go              # Database seeding
│
├── Dockerfile
├── Dockerfile.worker
├── go.mod
├── go.sum
├── .env.example
└── docker-compose.yml
```

## Key Files Description

### main.go
Application bootstrap, server initialization, graceful shutdown.

### config.go
Environment-based configuration with sensible defaults:
- Database connection
- Redis connection
- JWT secrets
- MinIO settings
- Server port

### Domain Entities
Pure Go structs representing business entities with JSON tags.

### Repository Layer
Data access abstraction supporting PostgreSQL and Redis.

### Service Layer
Business logic, validation, and orchestration.

### Handler Layer
HTTP handlers, request/response transformation.

### Middleware
Cross-cutting concerns (auth, logging, rate limiting).

### WebSocket Hub
Manages real-time connections and broadcasting.

## API Patterns

### Handler Pattern
```go
func (h *TaskHandler) Create(c *fiber.Ctx) error {
    var req CreateTaskRequest
    if err := c.BodyParser(&req); err != nil {
        return response.Error(c, 400, "invalid request")
    }
    
    task, err := h.taskService.Create(c.Context(), &req)
    if err != nil {
        return response.Error(c, 500, err.Error())
    }
    
    return response.Created(c, task)
}
```

### Service Pattern
```go
func (s *TaskService) Create(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
    if err := s.validateTask(req); err != nil {
        return nil, err
    }
    
    task := &Task{
        // ... mapping
    }
    
    return s.repo.Create(ctx, task)
}
```

## Testing Strategy
- Unit tests for services
- Handler tests with fiber test utils
- Integration tests with test database
