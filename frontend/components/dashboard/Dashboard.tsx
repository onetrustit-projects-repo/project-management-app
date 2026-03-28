'use client'

import { useQuery } from '@tanstack/react-query'
import { dashboardAPI } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { BarChart3, CheckCircle2, Clock, AlertTriangle, Users, FolderOpen } from 'lucide-react'

interface DashboardProps {
  projectKey: string
}

export function Dashboard({ projectKey }: DashboardProps) {
  const { data: stats, isLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: () => dashboardAPI.getStats('default').then(res => res.data),
  })

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    )
  }

  const statCards = [
    { title: 'Total Tasks', value: stats?.total_tasks || 0, icon: BarChart3, color: 'text-blue-500' },
    { title: 'Completed', value: stats?.completed_tasks || 0, icon: CheckCircle2, color: 'text-green-500' },
    { title: 'In Progress', value: stats?.in_progress_tasks || 0, icon: Clock, color: 'text-yellow-500' },
    { title: 'Overdue', value: stats?.overdue_tasks || 0, icon: AlertTriangle, color: 'text-red-500' },
    { title: 'Projects', value: stats?.total_projects || 0, icon: FolderOpen, color: 'text-purple-500' },
    { title: 'Team Members', value: stats?.total_members || 0, icon: Users, color: 'text-cyan-500' },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">Overview of your project metrics</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {statCards.map(card => (
          <Card key={card.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{card.title}</CardTitle>
              <card.icon className={`h-4 w-4 ${card.color}`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{card.value}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Tasks by Status Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Tasks by Status</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {stats?.tasks_by_status && Object.entries(stats.tasks_by_status).map(([status, count]) => (
              <div key={status} className="flex items-center gap-4">
                <div className="w-24 text-sm capitalize">{status.replace('_', ' ')}</div>
                <div className="flex-1 bg-muted rounded-full h-4 overflow-hidden">
                  <div
                    className="bg-primary h-full transition-all"
                    style={{ width: `${(count / stats.total_tasks) * 100}%` }}
                  />
                </div>
                <div className="w-8 text-sm text-right">{count}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Tasks by Priority */}
      <Card>
        <CardHeader>
          <CardTitle>Tasks by Priority</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-4 gap-4">
            {stats?.tasks_by_priority && Object.entries(stats.tasks_by_priority).map(([priority, count]) => (
              <div key={priority} className="text-center p-4 rounded-lg bg-muted">
                <div className="text-2xl font-bold">{count}</div>
                <div className="text-sm capitalize text-muted-foreground">{priority}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
