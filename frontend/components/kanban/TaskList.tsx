'use client'

import { useQuery } from '@tanstack/react-query'
import { formatDistanceToNow } from 'date-fns'
import { taskAPI } from '@/lib/api'
import { Avatar } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { MoreHorizontal, Plus } from 'lucide-react'
import type { Task } from '@/types'

interface TaskListProps {
  projectKey: string
}

const statusOrder = ['backlog', 'todo', 'in_progress', 'in_review', 'done']

const statusLabels: Record<string, string> = {
  backlog: 'Backlog',
  todo: 'To Do',
  in_progress: 'In Progress',
  in_review: 'In Review',
  done: 'Done',
}

const priorityColors: Record<string, string> = {
  low: 'bg-gray-100 text-gray-700',
  medium: 'bg-blue-100 text-blue-700',
  high: 'bg-orange-100 text-orange-700',
  urgent: 'bg-red-100 text-red-700',
}

export function TaskList({ projectKey }: TaskListProps) {
  const { data: tasks = [], isLoading } = useQuery({
    queryKey: ['tasks', projectKey],
    queryFn: () => taskAPI.listByProject(projectKey).then(res => res.data),
  })

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    )
  }

  const sortedTasks = [...tasks].sort((a, b) => {
    const statusA = statusOrder.indexOf(a.status)
    const statusB = statusOrder.indexOf(b.status)
    if (statusA !== statusB) return statusA - statusB
    return a.position - b.position
  })

  return (
    <div className="bg-card rounded-lg border">
      <div className="grid grid-cols-12 gap-4 p-4 border-b bg-muted/50 text-sm font-medium">
        <div className="col-span-1">ID</div>
        <div className="col-span-4">Title</div>
        <div className="col-span-2">Status</div>
        <div className="col-span-2">Assignee</div>
        <div className="col-span-2">Due Date</div>
        <div className="col-span-1">Priority</div>
      </div>

      {sortedTasks.map(task => (
        <div
          key={task.id}
          className="grid grid-cols-12 gap-4 p-4 border-b last:border-0 hover:bg-muted/30 transition-colors cursor-pointer"
        >
          <div className="col-span-1 font-mono text-sm text-muted-foreground">
            {task.task_number}
          </div>
          <div className="col-span-4 font-medium">{task.title}</div>
          <div className="col-span-2">
            <Badge variant="outline">{statusLabels[task.status]}</Badge>
          </div>
          <div className="col-span-2">
            {task.assignee && (
              <div className="flex items-center gap-2">
                <Avatar className="h-6 w-6" fallback={task.assignee.name?.[0]}>
                  {task.assignee.avatar_url && (
                    <img src={task.assignee.avatar_url} alt={task.assignee.name} />
                  )}
                </Avatar>
                <span className="text-sm">{task.assignee.name}</span>
              </div>
            )}
          </div>
          <div className="col-span-2 text-sm text-muted-foreground">
            {task.due_date ? formatDistanceToNow(new Date(task.due_date), { addSuffix: true }) : '-'}
          </div>
          <div className="col-span-1">
            <Badge className={priorityColors[task.priority]} variant="secondary">
              {task.priority}
            </Badge>
          </div>
        </div>
      ))}

      {sortedTasks.length === 0 && (
        <div className="p-8 text-center text-muted-foreground">
          No tasks found. Create your first task to get started.
        </div>
      )}
    </div>
  )
}
