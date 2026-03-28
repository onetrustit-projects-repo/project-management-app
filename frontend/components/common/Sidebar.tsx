'use client'

import { useQuery } from '@tanstack/react-query'
import { workspaceAPI, projectAPI } from '@/lib/api'
import { useProjectStore } from '@/lib/store'
import { Button } from '@/components/ui/button'
import { Plus, Settings, FolderKanban, LogOut, ChevronDown } from 'lucide-react'
import { useState } from 'react'

interface SidebarProps {
  selectedProject: string
  onSelectProject: (key: string) => void
}

export function Sidebar({ selectedProject, onSelectProject }: SidebarProps) {
  const { currentWorkspace, setCurrentWorkspace, setProjects } = useProjectStore()
  const [showWorkspaces, setShowWorkspaces] = useState(false)

  const { data: workspaces = [] } = useQuery({
    queryKey: ['workspaces'],
    queryFn: () => workspaceAPI.list().then(res => res.data),
  })

  const { data: projects = [] } = useQuery({
    queryKey: ['projects', currentWorkspace?.id],
    queryFn: () => currentWorkspace ? projectAPI.list(currentWorkspace.id).then(res => res.data) : Promise.resolve([]),
    enabled: !!currentWorkspace,
  })

  // Set first workspace as current if none selected
  if (!currentWorkspace && workspaces.length > 0) {
    setCurrentWorkspace(workspaces[0])
  }

  return (
    <div className="w-64 border-r bg-card flex flex-col h-full">
      {/* Logo */}
      <div className="p-4 border-b">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
            <span className="text-primary-foreground font-bold">PM</span>
          </div>
          <span className="font-semibold">ProjectFlow</span>
        </div>
      </div>

      {/* Workspace selector */}
      <div className="p-3 border-b">
        <button
          onClick={() => setShowWorkspaces(!showWorkspaces)}
          className="flex items-center justify-between w-full p-2 rounded-lg hover:bg-muted transition-colors"
        >
          <span className="text-sm font-medium truncate">
            {currentWorkspace?.name || 'Select Workspace'}
          </span>
          <ChevronDown className="h-4 w-4 text-muted-foreground" />
        </button>

        {showWorkspaces && (
          <div className="mt-2 space-y-1">
            {workspaces.map((ws: any) => (
              <button
                key={ws.id}
                onClick={() => {
                  setCurrentWorkspace(ws)
                  setShowWorkspaces(false)
                }}
                className={`w-full text-left p-2 rounded-lg text-sm ${
                  currentWorkspace?.id === ws.id ? 'bg-primary text-primary-foreground' : 'hover:bg-muted'
                }`}
              >
                {ws.name}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Projects */}
      <div className="flex-1 overflow-auto p-3">
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs font-semibold text-muted-foreground uppercase">Projects</span>
          <Button variant="ghost" size="icon" className="h-6 w-6">
            <Plus className="h-3 w-3" />
          </Button>
        </div>

        <div className="space-y-1">
          {projects.map((project: any) => (
            <button
              key={project.id}
              onClick={() => onSelectProject(project.key)}
              className={`w-full flex items-center gap-2 p-2 rounded-lg text-sm transition-colors ${
                selectedProject === project.key
                  ? 'bg-primary/10 text-primary'
                  : 'hover:bg-muted'
              }`}
            >
              <FolderKanban className="h-4 w-4" />
              <span className="truncate">{project.name}</span>
              <span className="ml-auto text-xs text-muted-foreground font-mono">
                {project.key}
              </span>
            </button>
          ))}

          {projects.length === 0 && (
            <div className="text-sm text-muted-foreground p-2">
              No projects yet
            </div>
          )}
        </div>
      </div>

      {/* Footer */}
      <div className="p-3 border-t">
        <button className="w-full flex items-center gap-2 p-2 rounded-lg hover:bg-muted text-sm text-muted-foreground">
          <Settings className="h-4 w-4" />
          <span>Settings</span>
        </button>
        <button className="w-full flex items-center gap-2 p-2 rounded-lg hover:bg-muted text-sm text-muted-foreground">
          <LogOut className="h-4 w-4" />
          <span>Log out</span>
        </button>
      </div>
    </div>
  )
}
