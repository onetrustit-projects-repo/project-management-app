'use client'

import { useQuery } from '@tanstack/react-query'
import { format, differenceInDays, addDays } from 'date-fns'
import { taskAPI } from '@/lib/api'
import type { Task } from '@/types'

interface TimelineViewProps {
  projectKey: string
}

export function TimelineView({ projectKey }: TimelineViewProps) {
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

  const tasksWithDates = tasks.filter(t => t.due_date)
  const today = new Date()

  // Calculate timeline range
  const startDates = tasksWithDates.map(t => t.due_date ? new Date(t.due_date) : today)
  const earliestDate = tasksWithDates.length > 0 
    ? new Date(Math.min(...startDates.map(d => d.getTime())))
    : addDays(today, -7)
  const latestDate = addDays(today, 60)
  const totalDays = differenceInDays(latestDate, earliestDate)

  const getPosition = (date: string) => {
    const d = new Date(date)
    return (differenceInDays(d, earliestDate) / totalDays) * 100
  }

  const getWidth = (estimatedHours?: number) => {
    if (!estimatedHours) return 2 // minimum width
    return Math.min(estimatedHours / 8 * 10, 30) // max 30%
  }

  const statusColors: Record<string, string> = {
    backlog: 'bg-gray-400',
    todo: 'bg-blue-500',
    in_progress: 'bg-yellow-500',
    in_review: 'bg-purple-500',
    done: 'bg-green-500',
  }

  return (
    <div className="bg-card rounded-lg border">
      <div className="p-4 border-b">
        <h2 className="font-semibold">Timeline View</h2>
        <p className="text-sm text-muted-foreground">Tasks plotted by their due dates</p>
      </div>

      {/* Timeline header */}
      <div className="flex border-b pl-24">
        {Array.from({ length: Math.ceil(totalDays / 7) }).map((_, i) => {
          const date = addDays(earliestDate, i * 7)
          return (
            <div key={i} className="flex-1 text-xs text-muted-foreground border-l pl-1">
              {format(date, 'MMM d')}
            </div>
          )
        })}
      </div>

      {/* Today marker */}
      <div className="relative h-8 bg-primary/10">
        <div 
          className="absolute top-0 bottom-0 w-0.5 bg-primary"
          style={{ left: `${getPosition(today.toISOString())}%` }}
        />
        <div 
          className="absolute -top-1 text-xs text-primary font-medium"
          style={{ left: `${getPosition(today.toISOString())}%` }}
        >
          Today
        </div>
      </div>

      {/* Task bars */}
      <div className="relative">
        {tasksWithDates.map((task, index) => (
          <div
            key={task.id}
            className="flex items-center border-b h-12 pl-4 hover:bg-muted/50"
          >
            <div className="w-20 text-sm truncate font-mono">{task.task_number}</div>
            <div className="flex-1 relative h-6">
              <div
                className={`absolute top-0 h-full rounded ${statusColors[task.status]} opacity-80 hover:opacity-100`}
                style={{
                  left: `${getPosition(task.due_date!)}%`,
                  width: `${getWidth(task.estimated_hours)}%`,
                }}
                title={`${task.title} - ${task.due_date ? format(new Date(task.due_date), 'MMM d') : ''}`}
              />
            </div>
          </div>
        ))}

        {tasksWithDates.length === 0 && (
          <div className="p-8 text-center text-muted-foreground">
            No tasks with due dates to display on timeline.
          </div>
        )}
      </div>
    </div>
  )
}
