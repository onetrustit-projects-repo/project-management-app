'use client'

import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { formatDistanceToNow } from 'date-fns'
import { MessageSquare, Paperclip, AlertCircle } from 'lucide-react'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import type { Task } from '@/types'

interface TaskCardProps {
  task: Task
  isDragging?: boolean
}

const priorityColors = {
  low: 'bg-gray-100 text-gray-700',
  medium: 'bg-blue-100 text-blue-700',
  high: 'bg-orange-100 text-orange-700',
  urgent: 'bg-red-100 text-red-700',
}

export function TaskCard({ task, isDragging }: TaskCardProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging: isSortableDragging,
  } = useSortable({ id: task.id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={`bg-card rounded-lg border p-3 cursor-grab active:cursor-grabbing shadow-sm ${
        isSortableDragging ? 'opacity-50 ring-2 ring-primary' : ''
      } ${isDragging ? 'shadow-lg rotate-3' : ''}`}
    >
      <div className="flex items-start justify-between gap-2 mb-2">
        <span className="text-xs text-muted-foreground font-mono">
          {task.task_number}
        </span>
        <Badge className={priorityColors[task.priority]} variant="secondary">
          {task.priority}
        </Badge>
      </div>

      <h4 className="font-medium text-sm mb-2 line-clamp-2">{task.title}</h4>

      {task.due_date && (
        <div className={`flex items-center gap-1 text-xs mb-2 ${
          new Date(task.due_date) < new Date() ? 'text-red-600' : 'text-muted-foreground'
        }`}>
          <AlertCircle className="h-3 w-3" />
          {formatDistanceToNow(new Date(task.due_date), { addSuffix: true })}
        </div>
      )}

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {task.assignee && (
            <Avatar className="h-6 w-6" fallback={task.assignee.name?.[0] || '?'}>
              {task.assignee.avatar_url && (
                <img src={task.assignee.avatar_url} alt={task.assignee.name} />
              )}
            </Avatar>
          )}
        </div>

        <div className="flex items-center gap-2 text-muted-foreground">
          <span className="flex items-center gap-1 text-xs">
            <MessageSquare className="h-3 w-3" />
          </span>
          <span className="flex items-center gap-1 text-xs">
            <Paperclip className="h-3 w-3" />
          </span>
        </div>
      </div>
    </div>
  )
}
