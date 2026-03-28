'use client'

import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Kanban, List, Calendar, GitCompare, LayoutDashboard, Plus, Bell, Search } from 'lucide-react'
import { useAuthStore } from '@/lib/store'

type View = 'kanban' | 'list' | 'calendar' | 'timeline' | 'dashboard'

interface HeaderProps {
  currentView: View
  onViewChange: (view: View) => void
  projectKey: string
}

const viewConfig: { id: View; label: string; icon: any }[] = [
  { id: 'kanban', label: 'Board', icon: Kanban },
  { id: 'list', label: 'List', icon: List },
  { id: 'calendar', label: 'Calendar', icon: Calendar },
  { id: 'timeline', label: 'Timeline', icon: GitCompare },
  { id: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
]

export function Header({ currentView, onViewChange, projectKey }: HeaderProps) {
  const { user } = useAuthStore()

  return (
    <header className="border-b bg-card">
      <div className="flex items-center justify-between px-6 py-3">
        <div className="flex items-center gap-4">
          <h1 className="text-lg font-semibold">Project: {projectKey}</h1>
          
          {/* View tabs */}
          <div className="flex items-center gap-1 ml-8">
            {viewConfig.map(view => (
              <Button
                key={view.id}
                variant={currentView === view.id ? 'secondary' : 'ghost'}
                size="sm"
                onClick={() => onViewChange(view.id)}
                className="gap-2"
              >
                <view.icon className="h-4 w-4" />
                {view.label}
              </Button>
            ))}
          </div>
        </div>

        <div className="flex items-center gap-3">
          <Button variant="outline" size="sm" className="gap-2">
            <Plus className="h-4 w-4" />
            New Task
          </Button>
          
          <Button variant="ghost" size="icon">
            <Search className="h-4 w-4" />
          </Button>
          
          <Button variant="ghost" size="icon">
            <Bell className="h-4 w-4" />
          </Button>

          <Avatar className="h-8 w-8" fallback={user?.name?.[0] || 'U'}>
            {user?.avatar_url && <img src={user.avatar_url} alt={user?.name || 'User'} />}
          </Avatar>
        </div>
      </div>
    </header>
  )
}
