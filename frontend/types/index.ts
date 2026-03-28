export type TaskStatus = 'backlog' | 'todo' | 'in_progress' | 'in_review' | 'done' | 'cancelled'
export type TaskPriority = 'low' | 'medium' | 'high' | 'urgent'
export type MemberRole = 'admin' | 'manager' | 'member' | 'viewer'

export interface User {
  id: string
  email: string
  name: string
  avatar_url?: string
  role: string
}

export interface Workspace {
  id: string
  name: string
  slug: string
  description: string
  owner_id: string
  settings: Record<string, any>
  created_at: string
  updated_at: string
}

export interface Project {
  id: string
  workspace_id: string
  name: string
  description: string
  key: string
  status: string
  start_date?: string
  end_date?: string
  owner_id: string
  created_at: string
  updated_at: string
}

export interface Task {
  id: string
  project_id: string
  title: string
  description: string
  status: TaskStatus
  priority: TaskPriority
  assignee_id?: string
  reporter_id?: string
  parent_id?: string
  due_date?: string
  estimated_hours?: number
  actual_hours: number
  position: number
  task_number: number
  created_at: string
  updated_at: string
  assignee?: User
  reporter?: User
  subtasks?: Task[]
}

export interface TaskComment {
  id: string
  task_id: string
  user_id: string
  content: string
  created_at: string
  updated_at: string
  user?: User
}

export interface ActivityLog {
  id: string
  workspace_id: string
  project_id: string
  task_id: string
  user_id: string
  action: string
  entity_type: string
  entity_id: string
  changes?: Record<string, any>
  created_at: string
}

export interface ProjectMember {
  id: string
  project_id: string
  user_id: string
  role: MemberRole
  created_at: string
  user?: User
}

export interface DashboardStats {
  total_tasks: number
  completed_tasks: number
  in_progress_tasks: number
  overdue_tasks: number
  total_projects: number
  active_projects: number
  total_members: number
  tasks_by_priority: Record<string, number>
  tasks_by_status: Record<string, number>
}
