# Frontend Folder Structure (Next.js)

```
frontend/
├── src/
│   ├── app/
│   │   ├── (auth)/
│   │   │   ├── login/
│   │   │   │   └── page.tsx
│   │   │   ├── register/
│   │   │   │   └── page.tsx
│   │   │   └── layout.tsx
│   │   │
│   │   ├── (dashboard)/
│   │   │   ├── layout.tsx          # Dashboard layout with sidebar
│   │   │   │
│   │   │   ├── projects/
│   │   │   │   ├── page.tsx        # Project list
│   │   │   │   └── [projectId]/
│   │   │   │       ├── page.tsx    # Project board
│   │   │   │       ├── settings/
│   │   │   │       │   └── page.tsx
│   │   │   │       └── members/
│   │   │   │           └── page.tsx
│   │   │   │
│   │   │   ├── tasks/
│   │   │   │   ├── page.tsx        # My tasks
│   │   │   │   └── [taskId]/
│   │   │   │       └── page.tsx    # Task detail
│   │   │   │
│   │   │   ├── calendar/
│   │   │   │   └── page.tsx        # Calendar view
│   │   │   │
│   │   │   ├── timeline/
│   │   │   │   └── page.tsx        # Gantt view
│   │   │   │
│   │   │   ├── dashboard/
│   │   │   │   └── page.tsx        # Dashboard home
│   │   │   │
│   │   │   └── settings/
│   │   │       ├── page.tsx
│   │   │       ├── profile/
│   │   │       ├── workspace/
│   │   │       └── billing/
│   │   │
│   │   ├── api/
│   │   │   └── [...trpc]/          # tRPC API routes
│   │   │       └── route.ts
│   │   │
│   │   ├── layout.tsx
│   │   ├── page.tsx                # Landing/redirect
│   │   └── globals.css
│   │
│   ├── components/
│   │   ├── ui/                     # Base UI components
│   │   │   ├── button.tsx
│   │   │   ├── input.tsx
│   │   │   ├── select.tsx
│   │   │   ├── modal.tsx
│   │   │   ├── dropdown-menu.tsx
│   │   │   ├── avatar.tsx
│   │   │   ├── badge.tsx
│   │   │   ├── tooltip.tsx
│   │   │   └── ...
│   │   │
│   │   ├── layout/
│   │   │   ├── sidebar.tsx
│   │   │   ├── header.tsx
│   │   │   ├── breadcrumb.tsx
│   │   │   └── command-menu.tsx    # Cmd+K palette
│   │   │
│   │   ├── projects/
│   │   │   ├── project-card.tsx
│   │   │   ├── project-list.tsx
│   │   │   ├── project-create-dialog.tsx
│   │   │   └── project-settings.tsx
│   │   │
│   │   ├── tasks/
│   │   │   ├── task-card.tsx
│   │   │   ├── task-list.tsx
│   │   │   ├── task-detail.tsx
│   │   │   ├── task-create-dialog.tsx
│   │   │   ├── task-filters.tsx
│   │   │   └── task-nlp-input.tsx   # AI NLP input
│   │   │
│   │   ├── kanban/
│   │   │   ├── board.tsx
│   │   │   ├── column.tsx
│   │   │   ├── draggable-card.tsx
│   │   │   └── add-column-dialog.tsx
│   │   │
│   │   ├── calendar/
│   │   │   ├── calendar-view.tsx
│   │   │   ├── calendar-day.tsx
│   │   │   └── calendar-event.tsx
│   │   │
│   │   ├── timeline/
│   │   │   ├── timeline-view.tsx
│   │   │   ├── timeline-bar.tsx
│   │   │   └── timeline-header.tsx
│   │   │
│   │   ├── comments/
│   │   │   ├── comment-list.tsx
│   │   │   ├── comment-item.tsx
│   │   │   └── comment-input.tsx
│   │   │
│   │   ├── notifications/
│   │   │   ├── notification-bell.tsx
│   │   │   ├── notification-list.tsx
│   │   │   └── notification-item.tsx
│   │   │
│   │   ├── dashboard/
│   │   │   ├── stats-card.tsx
│   │   │   ├── activity-feed.tsx
│   │   │   ├── project-progress.tsx
│   │   │   └── team-workload.tsx
│   │   │
│   │   └── ai/
│   │       ├── chat-panel.tsx
│   │       ├── suggestions.tsx
│   │       └── nlp-input.tsx
│   │
│   ├── hooks/
│   │   ├── use-auth.ts
│   │   ├── use-projects.ts
│   │   ├── use-tasks.ts
│   │   ├── use-websocket.ts
│   │   ├── use-notifications.ts
│   │   └── use-optimistic-update.ts
│   │
│   ├── lib/
│   │   ├── api/
│   │   │   ├── client.ts           # API client setup
│   │   │   ├── endpoints.ts        # API endpoints
│   │   │   └── types.ts            # API types
│   │   │
│   │   ├── utils/
│   │   │   ├── cn.ts              # classname utility
│   │   │   ├── format.ts           # Date, number formatting
│   │   │   └── validators.ts
│   │   │
│   │   └── constants.ts
│   │
│   ├── stores/
│   │   ├── auth-store.ts          # Zustand auth store
│   │   ├── project-store.ts
│   │   ├── task-store.ts
│   │   ├── ui-store.ts            # UI state (sidebar, modals)
│   │   └── notification-store.ts
│   │
│   ├── types/
│   │   ├── user.ts
│   │   ├── project.ts
│   │   ├── task.ts
│   │   └── api.ts
│   │
│   └── styles/
│       └── themes/
│           ├── light.css
│           └── dark.css
│
├── public/
│   ├── favicon.ico
│   └── images/
│
├── components.json               # shadcn/ui component registry
├── tailwind.config.ts
├── next.config.js
├── package.json
├── tsconfig.json
└── Dockerfile
```

## Component Patterns

### Task Card (Kanban)
```tsx
'use client'

import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { Task } from '@/types/task'
import { Avatar, Badge } from '@/components/ui'

interface TaskCardProps {
  task: Task
  onClick?: () => void
}

export function TaskCard({ task, onClick }: TaskCardProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging
  } = useSortable({ id: task.id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      onClick={onClick}
      className="bg-white rounded-lg border border-slate-200 p-3 cursor-grab active:cursor-grabbing hover:shadow-md transition-shadow"
    >
      <div className="flex items-start justify-between gap-2">
        <span className="text-xs text-slate-400">#{task.taskNumber}</span>
        <Badge variant={task.priority}>{task.priority}</Badge>
      </div>
      <h4 className="mt-2 font-medium text-slate-900 line-clamp-2">
        {task.title}
      </h4>
      <div className="mt-3 flex items-center justify-between">
        <div className="flex items-center gap-2">
          {task.labels.slice(0, 2).map(label => (
            <span key={label} className="text-xs bg-slate-100 px-2 py-0.5 rounded">
              {label}
            </span>
          ))}
        </div>
        {task.assignee && <Avatar src={task.assignee.avatarUrl} />}
      </div>
    </div>
  )
}
```

### Drag & Drop Board
```tsx
'use client'

import {
  DndContext,
  DragOverlay,
  closestCorners,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors
} from '@dnd-kit/core'
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy
} from '@dnd-kit/sortable'

export function KanbanBoard({ columns, tasks, onMove }) {
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
  )

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCorners}
      onDragEnd={handleDragEnd}
    >
      <div className="flex gap-4 overflow-x-auto pb-4">
        {columns.map(column => (
          <div key={column.id} className="flex-shrink-0 w-80">
            <Column column={column} tasks={getTasksByColumn(column.id)} />
          </div>
        ))}
      </div>
    </DndContext>
  )
}
```

## State Management

### Zustand Store Example
```ts
// stores/task-store.ts
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface TaskState {
  tasks: Record<string, Task>
  optimisticUpdates: Map<string, Task>
  
  setTasks: (tasks: Task[]) => void
  updateTask: (id: string, changes: Partial<Task>) => void
  moveTask: (taskId: string, columnId: string, position: number) => void
}

export const useTaskStore = create<TaskState>()(
  persist(
    (set, get) => ({
      tasks: {},
      optimisticUpdates: new Map(),
      
      setTasks: (tasks) => set({
        tasks: Object.fromEntries(tasks.map(t => [t.id, t]))
      }),
      
      updateTask: (id, changes) => set((state) => ({
        tasks: {
          ...state.tasks,
          [id]: { ...state.tasks[id], ...changes }
        }
      })),
      
      moveTask: (taskId, columnId, position) => {
        // Optimistic update with WebSocket sync
      }
    }),
    { name: 'task-storage' }
  )
)
```

## Real-time Updates
```tsx
// hooks/use-websocket.ts
import { useEffect } from 'react'
import { useTaskStore } from '@/stores/task-store'
import { useNotificationStore } from '@/stores/notification-store'

export function useWebSocket() {
  const updateTask = useTaskStore(s => s.updateTask)
  const addNotification = useNotificationStore(s => s.add)

  useEffect(() => {
    const ws = new WebSocket(`${WS_URL}?token=${getToken()}`)
    
    ws.onmessage = (event) => {
      const { type, data } = JSON.parse(event.data)
      
      switch (type) {
        case 'task.updated':
          updateTask(data.task_id, data.changes)
          break
        case 'notification':
          addNotification(data.notification)
          break
      }
    }
    
    return () => ws.close()
  }, [])
}
```

## Key Dependencies
- **Next.js 14** - App Router
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **shadcn/ui** - Base components
- **Zustand** - State management
- **TanStack Query** - Server state
- **dnd-kit** - Drag and drop
- **date-fns** - Date handling
- **Lucide** - Icons
- **Zod** - Validation
