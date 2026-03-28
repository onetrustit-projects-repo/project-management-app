import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User, Project, Workspace } from '@/types'

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  setAuth: (user: User, accessToken: string, refreshToken: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      setAuth: (user, accessToken, refreshToken) => {
        localStorage.setItem('access_token', accessToken)
        localStorage.setItem('refresh_token', refreshToken)
        set({ user, accessToken, refreshToken })
      },
      logout: () => {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        set({ user: null, accessToken: null, refreshToken: null })
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ user: state.user }),
    }
  )
)

interface ProjectState {
  currentWorkspace: Workspace | null
  currentProject: Project | null
  projects: Project[]
  workspaces: Workspace[]
  setCurrentWorkspace: (ws: Workspace) => void
  setCurrentProject: (project: Project) => void
  setProjects: (projects: Project[]) => void
  setWorkspaces: (workspaces: Workspace[]) => void
}

export const useProjectStore = create<ProjectState>((set) => ({
  currentWorkspace: null,
  currentProject: null,
  projects: [],
  workspaces: [],
  setCurrentWorkspace: (ws) => set({ currentWorkspace: ws }),
  setCurrentProject: (project) => set({ currentProject: project }),
  setProjects: (projects) => set({ projects }),
  setWorkspaces: (workspaces) => set({ workspaces }),
}))
