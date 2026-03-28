'use client'

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  DndContext,
  DragOverlay,
  closestCorners,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragStartEvent,
  DragEndEvent,
} from '@dnd-kit/core'
import { arrayMove, sortableKeyboardCoordinates } from '@dnd-kit/sortable'
import { useState } from 'react'
import { taskAPI } from '@/lib/api'
import { KanbanColumn } from './KanbanColumn'
import { TaskCard } from './TaskCard'
import { CreateTaskModal } from './CreateTaskModal'
import type { Task, TaskStatus } from '@/types'

const COLUMNS: { id: TaskStatus; title: string; color: string }[] = [
  { id: 'backlog', title: 'Backlog', color: 'bg-gray-500' },
  { id: 'todo', title: 'To Do', color: 'bg-blue-500' },
  { id: 'in_progress', title: 'In Progress', color: 'bg-yellow-500' },
  { id: 'in_review', title: 'In Review', color: 'bg-purple-500' },
  { id: 'done', title: 'Done', color: 'bg-green-500' },
]

interface KanbanBoardProps {
  projectKey: string
}

export function KanbanBoard({ projectKey }: KanbanBoardProps) {
  const queryClient = useQueryClient()
  const [activeTask, setActiveTask] = useState<Task | null>(null)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [createColumn, setCreateColumn] = useState<TaskStatus>('backlog')

  const { data: tasks = [], isLoading } = useQuery({
    queryKey: ['tasks', projectKey],
    queryFn: () => taskAPI.listByProject(projectKey).then(res => res.data),
  })

  const moveMutation = useMutation({
    mutationFn: ({ taskId, status, position }: { taskId: string; status: TaskStatus; position: number }) =>
      taskAPI.move(taskId, status, position),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks', projectKey] })
    },
  })

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
  )

  const getTasksByStatus = (status: TaskStatus) =>
    tasks.filter(t => t.status === status).sort((a, b) => a.position - b.position)

  const handleDragStart = (event: DragStartEvent) => {
    const task = tasks.find(t => t.id === event.active.id)
    if (task) setActiveTask(task)
  }

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    setActiveTask(null)

    if (!over) return

    const taskId = active.id as string
    const task = tasks.find(t => t.id === taskId)
    if (!task) return

    // Determine new status and position
    let newStatus: TaskStatus = task.status
    let newPosition = task.position

    // Check if dropped on a column
    const overColumn = COLUMNS.find(c => c.id === over.id)
    if (overColumn) {
      newStatus = overColumn.id
      const tasksInColumn = getTasksByStatus(newStatus)
      newPosition = tasksInColumn.length
    } else {
      // Dropped on another task
      const overTask = tasks.find(t => t.id === over.id)
      if (overTask) {
        newStatus = overTask.status
        newPosition = overTask.position
      }
    }

    if (newStatus !== task.status || newPosition !== task.position) {
      moveMutation.mutate({ taskId, status: newStatus, position: newPosition })
    }
  }

  const handleCreateTask = (columnId: TaskStatus) => {
    setCreateColumn(columnId)
    setIsCreateModalOpen(true)
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    )
  }

  return (
    <div className="h-full">
      <DndContext
        sensors={sensors}
        collisionDetection={closestCorners}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
      >
        <div className="flex gap-4 h-full overflow-x-auto pb-4">
          {COLUMNS.map(column => (
            <KanbanColumn
              key={column.id}
              id={column.id}
              title={column.title}
              color={column.color}
              tasks={getTasksByStatus(column.id)}
              onAddTask={() => handleCreateTask(column.id)}
            />
          ))}
        </div>

        <DragOverlay>
          {activeTask && (
            <TaskCard task={activeTask} isDragging />
          )}
        </DragOverlay>
      </DndContext>

      <CreateTaskModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        projectKey={projectKey}
        defaultStatus={createColumn}
      />
    </div>
  )
}
