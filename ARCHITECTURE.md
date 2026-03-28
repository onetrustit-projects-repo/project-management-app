# Project Management App - System Architecture

## Overview

A scalable, production-ready project management platform designed for 10-50 users initially, with architecture supporting growth to 10,000+ users.

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                    │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │   Next.js SPA   │  │   Mobile App    │  │   Third-Party   │              │
│  │   (React/TS)    │  │   (React Native│  │     Integrations │              │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘              │
└───────────┼─────────────────────┼─────────────────────┼──────────────────────┘
            │                     │                     │
            └─────────────────────┼─────────────────────┘
                                  │ HTTPS/WSS
┌─────────────────────────────────┴───────────────────────────────────────────┐
│                              API GATEWAY LAYER                               │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                        Nginx / Traefik                                │   │
│  │                   (SSL Termination, Load Balancing)                  │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                  │
┌─────────────────────────────────┴───────────────────────────────────────────┐
│                              SERVICE LAYER                                   │
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │   API Service   │  │  WebSocket Srv  │  │  AI Service     │              │
│  │   (Go/Fiber)    │  │  (Centrifugo)   │  │  (Go + LLM)     │              │
│  │   :8080         │  │  :8000          │  │  :8081          │              │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘              │
│           │                    │                     │                        │
│  ┌────────┴────────────────────┴─────────────────────┴────────┐           │
│  │                      Message Queue (Redis Streams)            │           │
│  └──────────────────────────────────────────────────────────────┘           │
│           │                    │                     │                        │
│  ┌────────┴────────────────────┴─────────────────────┴────────┐           │
│  │                    Background Workers                        │           │
│  │  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐  │           │
│  │  │ Notification│  │  Email   │  │  AI Task  │  │  Audit    │  │           │
│  │  │  Worker   │  │  Worker  │  │  Worker   │  │  Worker   │  │           │
│  │  └───────────┘  └───────────┘  └───────────┘  └───────────┘  │           │
│  └──────────────────────────────────────────────────────────────┘           │
└─────────────────────────────────────────────────────────────────────────────┘
                                  │
┌─────────────────────────────────┴───────────────────────────────────────────┐
│                              DATA LAYER                                      │
│                                                                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │   PostgreSQL    │  │     Redis       │  │   MinIO/S3      │              │
│  │   (Primary DB)  │  │   (Cache/Queue) │  │   (File Storage)│              │
│  │   :5432         │  │   :6379         │  │   :9000         │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Component Architecture

### Backend (Go)

```
┌────────────────────────────────────────────────────────────────┐
│                        Go API Service                           │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   Handlers   │  │   Services   │  │  Repository  │          │
│  │   (Routes)   │──│   (Business) │──│   (Database) │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│         │                │                  │                   │
│         └────────────────┼──────────────────┘                   │
│                          │                                      │
│  ┌───────────────────────┴────────────────────────────────┐    │
│  │                    Middleware                           │    │
│  │  Auth | RBAC | Logging | Rate Limit | Validation       │    │
│  └─────────────────────────────────────────────────────────┘    │
└────────────────────────────────────────────────────────────────┘
```

### Frontend (Next.js)

```
┌────────────────────────────────────────────────────────────────┐
│                      Next.js Application                        │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │    Pages     │  │  Components  │  │    Hooks     │          │
│  │  /projects   │  │   <Kanban>   │  │  useTasks()  │          │
│  │  /tasks      │  │   <Calendar> │  │  useProjects │          │
│  │  /dashboard  │  │   <Timeline> │  │  useWebSocket│          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                             │   │
│  ┌─────────────────────────────────────────────────────────┴┐  │
│  │                     State Management                      │  │
│  │              (Zustand + React Query + WebSocket)          │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────┘
```

---

## Data Flow

### Task Creation Flow

```
User Input (NLP) → API → Validation → PostgreSQL → Redis Queue → AI Worker → Suggestions
                         ↓
                   WebSocket → Frontend (Real-time Update)
```

### Authentication Flow

```
User Login → JWT Token → Refresh Token Rotation → Secure Cookie
                                ↓
                    RBAC Middleware → Permission Check → Resource Access
```

---

## Deployment Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           Docker Compose (Production)                    │
│                                                                          │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐          │
│  │  nginx  │ │  api    │ │ worker  │ │  redis  │ │ postgres│          │
│  │ (proxy) │ │ (go)    │ │ (go)    │ │ (cache) │ │  (db)   │          │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘          │
│                                  │                                        │
│                         ┌────────┴────────┐                             │
│                         │   minio (s3)    │                             │
│                         │  (file storage) │                             │
│                         └─────────────────┘                             │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Technology Stack Summary

| Layer          | Technology                          |
|----------------|-------------------------------------|
| Frontend       | Next.js 14, React 18, TypeScript   |
| UI Components  | Radix UI, Tailwind CSS             |
| State          | Zustand, TanStack Query             |
| Backend        | Go 1.21+, Fiber v2                 |
| Database       | PostgreSQL 15+                      |
| Cache/Queue    | Redis 7+                            |
| File Storage   | MinIO (S3-compatible)               |
| Real-time      | WebSocket (Gorilla/Fiuba)           |
| AI             | Go + OpenAI/Anthropic API          |
| Auth           | JWT + Refresh Tokens                |
| Container      | Docker, Docker Compose              |
| CI/CD          | GitHub Actions                      |
| Reverse Proxy  | Nginx/Traefik                       |
