'use client'

import { useQuery } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { format, startOfMonth, endOfMonth, eachDayOfInterval, isSameMonth, isToday, addMonths, subMonths } from 'date-fns'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { taskAPI } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import type { Task } from '@/types'

interface CalendarViewProps {
  projectKey: string
}

export function CalendarView({ projectKey }: CalendarViewProps) {
  const [currentDate, setCurrentDate] = useState(new Date())
  
  const { data: tasks = [] } = useQuery({
    queryKey: ['tasks', projectKey],
    queryFn: () => taskAPI.listByProject(projectKey).then(res => res.data),
  })

  const monthStart = startOfMonth(currentDate)
  const monthEnd = endOfMonth(currentDate)
  const days = eachDayOfInterval({ start: monthStart, end: monthEnd })

  // Get first day of month for padding
  const startDay = monthStart.getDay()
  const paddingDays = Array(startDay).fill(null)

  const tasksByDate = useMemo(() => {
    const grouped: Record<string, Task[]> = {}
    tasks.forEach(task => {
      if (task.due_date) {
        const dateKey = format(new Date(task.due_date), 'yyyy-MM-dd')
        if (!grouped[dateKey]) grouped[dateKey] = []
        grouped[dateKey].push(task)
      }
    })
    return grouped
  }, [tasks])

  const previousMonth = () => setCurrentDate(subMonths(currentDate, 1))
  const nextMonth = () => setCurrentDate(addMonths(currentDate, 1))

  return (
    <div className="bg-card rounded-lg border">
      <div className="flex items-center justify-between p-4 border-b">
        <h2 className="font-semibold">{format(currentDate, 'MMMM yyyy')}</h2>
        <div className="flex gap-2">
          <Button variant="outline" size="icon" onClick={previousMonth}>
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button variant="outline" size="icon" onClick={nextMonth}>
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-7">
        {['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'].map(day => (
          <div key={day} className="p-2 text-center text-sm font-medium text-muted-foreground border-b">
            {day}
          </div>
        ))}

        {paddingDays.map((_, i) => (
          <div key={`pad-${i}`} className="min-h-[100px] border-b border-r bg-muted/30" />
        ))}

        {days.map(day => {
          const dateKey = format(day, 'yyyy-MM-dd')
          const dayTasks = tasksByDate[dateKey] || []
          
          return (
            <div
              key={day.toISOString()}
              className={`min-h-[100px] border-b border-r p-2 ${
                !isSameMonth(day, currentDate) ? 'bg-muted/30 text-muted-foreground' : ''
              }`}
            >
              <div className={`text-sm font-medium mb-1 ${
                isToday(day) ? 'bg-primary text-primary-foreground rounded-full w-7 h-7 flex items-center justify-center' : ''
              }`}>
                {format(day, 'd')}
              </div>
              
              <div className="space-y-1">
                {dayTasks.slice(0, 3).map(task => (
                  <div
                    key={task.id}
                    className="text-xs p-1 rounded bg-primary/10 text-primary truncate cursor-pointer hover:bg-primary/20"
                    title={task.title}
                  >
                    {task.title}
                  </div>
                ))}
                {dayTasks.length > 3 && (
                  <div className="text-xs text-muted-foreground">
                    +{dayTasks.length - 3} more
                  </div>
                )}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
