'use client'

import { useState } from 'react'
import { KanbanBoard } from '@/components/kanban/KanbanBoard'
import { TaskList } from '@/components/kanban/TaskList'
import { CalendarView } from '@/components/calendar/CalendarView'
import { TimelineView } from '@/components/timeline/TimelineView'
import { Dashboard } from '@/components/dashboard/Dashboard'
import { Sidebar } from '@/components/common/Sidebar'
import { Header } from '@/components/common/Header'
import { useAuthStore } from '@/lib/store'

type View = 'kanban' | 'list' | 'calendar' | 'timeline' | 'dashboard'

export default function Home() {
  const [currentView, setCurrentView] = useState<View>('kanban')
  const [selectedProject, setSelectedProject] = useState<string>('PM')
  const { user } = useAuthStore()

  return (
    <div className="flex h-screen bg-background">
      <Sidebar 
        selectedProject={selectedProject} 
        onSelectProject={setSelectedProject} 
      />
      
      <div className="flex-1 flex flex-col overflow-hidden">
        <Header 
          currentView={currentView} 
          onViewChange={setCurrentView}
          projectKey={selectedProject}
        />
        
        <main className="flex-1 overflow-auto p-6">
          {currentView === 'kanban' && (
            <KanbanBoard projectKey={selectedProject} />
          )}
          {currentView === 'list' && (
            <TaskList projectKey={selectedProject} />
          )}
          {currentView === 'calendar' && (
            <CalendarView projectKey={selectedProject} />
          )}
          {currentView === 'timeline' && (
            <TimelineView projectKey={selectedProject} />
          )}
          {currentView === 'dashboard' && (
            <Dashboard projectKey={selectedProject} />
          )}
        </main>
      </div>
    </div>
  )
}
