# Project Management App

AI-powered project management platform built with Go backend and Next.js frontend.

## Features

### Core Features
- ✅ Multi-tenant workspaces
- ✅ Project management with Kanban, List, Calendar, and Timeline views
- ✅ Task management with assignments, priorities, dependencies
- ✅ Real-time collaboration via WebSocket
- ✅ Role-based access control (RBAC)
- ✅ Comments and activity tracking

### Views
- **Kanban Board** - Drag-and-drop task management
- **List View** - Tabular task list with sorting/filtering
- **Calendar View** - Tasks plotted by due dates
- **Timeline View** - Gantt-style visualization

### Dashboard
- Project statistics
- Tasks by status/priority charts
- Team metrics

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Next.js 14, React, Tailwind CSS, dnd-kit |
| Backend | Go 1.21+, Gin framework |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Real-time | WebSocket (gorilla/websocket) |
| Auth | JWT |

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Node.js 22+ (for local development)
- Go 1.21+ (for local development)

### Using Docker

```bash
# Clone the repository
git clone https://github.com/your-org/pm-app.git
cd pm-app

# Start all services
docker compose -f docker/docker-compose.yml up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

Access the app at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Docs: http://localhost:8080/swagger

### Local Development

**Backend:**
```bash
cd backend
cp .env.example .env
go mod download
go run cmd/server/main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## Project Structure

```
pm-app/
├── backend/
│   ├── cmd/server/       # Entry point
│   ├── internal/
│   │   ├── models/       # Data models
│   │   ├── handlers/     # HTTP handlers
│   │   ├── middleware/    # Auth, CORS, rate limiting
│   │   ├── repositories/  # Database operations
│   │   ├── services/      # Business logic
│   │   └── websocket/     # WebSocket hub
│   └── migrations/       # SQL migrations
├── frontend/
│   ├── app/              # Next.js App Router
│   ├── components/        # React components
│   │   ├── kanban/       # Kanban board
│   │   ├── calendar/     # Calendar view
│   │   ├── timeline/     # Timeline/Gantt view
│   │   ├── dashboard/    # Dashboard
│   │   └── ui/           # Base UI components
│   ├── lib/              # Utilities, API client
│   └── hooks/            # Custom React hooks
├── docker/               # Docker configurations
├── ci-cd/                # CI/CD pipelines
└── docs/                 # Documentation
```

## API Endpoints

### Authentication
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
```

### Workspaces
```
GET    /api/v1/workspaces
POST   /api/v1/workspaces
GET    /api/v1/workspaces/:id
PUT    /api/v1/workspaces/:id
DELETE /api/v1/workspaces/:id
```

### Projects
```
GET    /api/v1/projects
POST   /api/v1/projects
GET    /api/v1/projects/:key
PUT    /api/v1/projects/:key
DELETE /api/v1/projects/:key
```

### Tasks
```
GET    /api/v1/projects/:key/tasks
POST   /api/v1/projects/:key/tasks
GET    /api/v1/tasks/:id
PUT    /api/v1/tasks/:id
DELETE /api/v1/tasks/:id
POST   /api/v1/tasks/:id/move
POST   /api/v1/tasks/:id/comments
```

## Environment Variables

```env
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/pm_app?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# Auth
JWT_SECRET=your-secret-key

# Server
PORT=8080
GIN_MODE=release
```

## License

MIT
