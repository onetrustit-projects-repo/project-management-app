# AI Features Implementation

## Overview

This document describes the AI-powered features and their implementation approach.

## AI Service Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         AI Service (Go)                         │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │  NLP Parser  │  │ LLM Client   │  │ Prompt Eng.  │        │
│  │  (task input)│  │ (OpenAI/     │  │ (templates)  │        │
│  │              │  │  Claude)     │  │              │        │
│  └──────┬───────┘  └──────┬───────┘  └──────────────┘        │
│         │                  │                                   │
│         └──────────────────┼───────────────────────────────────┘│
│                            │                                    │
│  ┌─────────────────────────┴────────────────────────────────┐   │
│  │                    AI Worker Pool                         │   │
│  │  - Task breakdown                                        │   │
│  │  - Risk detection                                        │   │
│  │  - Auto-prioritization                                   │   │
│  │  - Report generation                                     │   │
│  └───────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Features

### 1. Natural Language Task Creation

**Input:** "Create 3 tasks for launching marketing campaign by next Friday"

**Output:**
```json
{
  "tasks": [
    { "title": "Prepare launch assets", "priority": "high", "due_date": "2024-03-14" },
    { "title": "Coordinate with design team", "priority": "medium", "due_date": "2024-03-13" },
    { "title": "Send announcement email", "priority": "medium", "due_date": "2024-03-15" }
  ]
}
```

**Implementation:**
```go
// internal/service/ai_service.go
func (s *AIService) ParseTaskCreation(ctx context.Context, input string, projectID uuid.UUID) (*TaskCreationResult, error) {
    prompt := fmt.Sprintf(`
    Parse the following task request and extract task details.
    Return valid JSON with an array of tasks.
    
    Request: "%s"
    Project ID: %s
    
    Return format:
    {
      "tasks": [
        {"title": "...", "description": "...", "priority": "high|medium|low", "due_date": "YYYY-MM-DD"}
      ]
    }
    `, input, projectID)
    
    response, err := s.llmClient.Chat(ctx, &llm.ChatRequest{
        Model: "gpt-4-turbo",
        Messages: []llm.Message{{Role: "user", Content: prompt}},
    })
    
    return parseTaskCreationResponse(response)
}
```

### 2. Smart Task Breakdown

**Trigger:** When a task is marked as "Epic" or has story_points > 5

**Process:**
1. Identify large tasks
2. Query similar past tasks for context
3. Generate subtask suggestions
4. Present to user for confirmation

```go
func (s *AIService) SuggestTaskBreakdown(ctx context.Context, task *Task) ([]SubtaskSuggestion, error) {
    similarTasks, err := s.repo.GetSimilarTasks(ctx, task.ProjectID, task.Title)
    
    prompt := fmt.Sprintf(`
    This task: "%s" seems large. Based on similar tasks:
    %s
    
    Suggest how to break it down into smaller, manageable subtasks.
    Consider dependencies and logical ordering.
    `, task.Title, formatSimilarTasks(similarTasks))
    
    // ... LLM call and parsing
}
```

### 3. Deadline Risk Detection

**Trigger:** Daily worker scan, on task update

**Logic:**
```
IF task.due_date < NOW + 3 days
   AND task.status NOT IN (done, cancelled)
   AND blockers exist AND blockers.status NOT IN (done, cancelled)
THEN flag as "at risk"
```

```go
func (s *AIService) DetectDeadlineRisks(ctx context.Context) ([]RiskAlert, error) {
    // Scan all tasks due within 3 days
    tasks := s.repo.GetTasksDueSoon(ctx, 3*24*time.Hour)
    
    var alerts []RiskAlert
    for _, task := range tasks {
        if s.hasUnresolvedBlockers(ctx, task) {
            alerts = append(alerts, RiskAlert{
                Task: task,
                Risk: "blocked_by_incomplete_dependency",
                Message: fmt.Sprintf("'%s' is blocked by incomplete task(s)", task.Title),
            })
        }
    }
    
    return alerts, nil
}
```

### 4. Auto-Prioritization

**Algorithm:**
1. Calculate priority score based on:
   - Due date urgency (exponential decay)
   - Dependency weight (tasks blocking others score higher)
   - Business impact (label-based)
   - Assignee workload

2. Formula:
```
priority_score = (
    due_date_factor * 0.4 +
    dependency_factor * 0.3 +
    impact_factor * 0.2 +
    workload_factor * 0.1
)
```

```go
func (s *AIService) CalculatePriorityScore(ctx context.Context, task *Task) (float64, error) {
    dueFactor := s.calculateDueDateFactor(task.DueDate)
    depFactor := s.calculateDependencyFactor(ctx, task)
    impactFactor := s.calculateImpactFactor(task.Labels)
    workloadFactor := s.calculateWorkloadFactor(ctx, task.AssigneeID)
    
    return (dueFactor * 0.4) + (depFactor * 0.3) + 
           (impactFactor * 0.2) + (workloadFactor * 0.1), nil
}
```

### 5. Meeting Notes → Tasks

**Input:** Transcript or notes from meeting

**Process:**
1. Extract action items (LLM)
2. Identify assignees (name mentions)
3. Set default due dates
4. Create tasks

```go
func (s *AIService) ExtractTasksFromNotes(ctx context.Context, notes string, projectID uuid.UUID) ([]TaskInput, error) {
    prompt := fmt.Sprintf(`
    Extract actionable tasks from the following meeting notes.
    Identify WHO should do WHAT and by WHEN.
    
    Notes:
    %s
    
    Return JSON array of tasks:
    [{"title": "...", "assignee": "...", "due_date": "YYYY-MM-DD", "priority": "..."}]
    `, notes)
    
    // ... LLM processing
}
```

### 6. AI Status Reports

**Weekly automated report includes:**
- Project progress summary
- Tasks completed vs planned
- At-risk items
- Team workload distribution
- Recommendations

```go
func (s *AIService) GenerateWeeklyReport(ctx context.Context, projectID uuid.UUID) (*Report, error) {
    stats := s.aggregateProjectStats(ctx, projectID)
    
    prompt := fmt.Sprintf(`
    Generate a weekly project status report based on:
    
    - Tasks completed: %d
    - Tasks in progress: %d  
    - Overdue tasks: %d
    - Team members active: %d
    - Progress this week: %s
    
    Provide:
    1. Executive summary (3-4 sentences)
    2. Key accomplishments
    3. Blockers and risks
    4. Recommendations for next week
    
    Format as structured JSON.
    `, stats.Completed, stats.InProgress, stats.Overdue, stats.ActiveMembers, stats.ProgressTrend)
    
    return s.llmClient.Chat(ctx, &llmRequest)
}
```

### 7. Chat Assistant

**Capabilities:**
- Natural language queries
- Task creation and updates
- Status lookups
- Scheduling assistance

```go
type AIChatService struct {
    conversationStore *redis.Store
    projectContext    ProjectContextLoader
    llmClient         LLMClient
}

func (s *AIChatService) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
    // Load conversation history
    history := s.conversationStore.GetMessages(ctx, req.ConversationID)
    
    // Load project context (members, recent tasks, etc.)
    projectCtx := s.projectContext.Load(ctx, req.ProjectID)
    
    prompt := s.buildContextPrompt(projectCtx, history, req.Message)
    
    response, err := s.llmClient.Chat(ctx, &llm.ChatRequest{
        Model: "gpt-4-turbo",
        Messages: append(history, llm.Message{
            Role: "user", 
            Content: prompt,
        }),
    })
    
    // Save to conversation history
    s.conversationStore.AddMessage(ctx, req.ConversationID, req.Message, response)
    
    return &ChatResponse{
        Message:     response.Content,
        Suggestions: s.generateSuggestions(ctx, response),
    }, nil
}
```

## Worker Jobs

```go
// internal/worker/ai_worker.go
func (w *AIWorker) RegisterJobs() {
    // Daily risk detection
    w.cron.Every(1).Day().Do(w.scanDeadlineRisks)
    
    // Weekly reports
    w.cron.Every(1).Week().Monday().At("8:00").Do(w.generateWeeklyReports)
    
    // Real-time task suggestions
    w.realtime.On("task.created", w.suggestBreakdown)
    w.realtime.On("task.updated", w.reassessPriority)
}
```

## LLM Configuration

```yaml
# config.yaml
ai:
  provider: "openai"  # or "anthropic"
  model: "gpt-4-turbo-preview"
  temperature: 0.7
  max_tokens: 2048
  cache_enabled: true
  cache_ttl: 3600  # 1 hour
```

## Safety & Governance

1. **Human-in-the-loop:** AI suggestions require user confirmation before execution
2. **Audit trail:** All AI decisions logged with confidence scores
3. **Fallback:** Clear error handling when AI service unavailable
4. **Rate limiting:** AI endpoints protected with queue-based access
5. **Cost control:** Max tokens per request, caching to reduce API calls
