import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Try to refresh token
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        try {
          const response = await axios.post(`${API_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          })
          const { access_token, refresh_token } = response.data
          localStorage.setItem('access_token', access_token)
          localStorage.setItem('refresh_token', refresh_token)
          
          // Retry original request
          error.config.headers.Authorization = `Bearer ${access_token}`
          return axios(error.config)
        } catch {
          // Refresh failed, logout
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          window.location.href = '/login'
        }
      }
    }
    return Promise.reject(error)
  }
)

export const authAPI = {
  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),
  register: (email: string, password: string, name: string) =>
    api.post('/auth/register', { email, password, name }),
  refresh: (refreshToken: string) =>
    api.post('/auth/refresh', { refresh_token: refreshToken }),
  logout: (refreshToken: string) =>
    api.post('/auth/logout', { refresh_token: refreshToken }),
}

export const workspaceAPI = {
  list: () => api.get('/workspaces'),
  create: (data: { name: string; description?: string }) =>
    api.post('/workspaces', data),
  get: (id: string) => api.get(`/workspaces/${id}`),
  update: (id: string, data: Partial<{ name: string; description: string }>) =>
    api.put(`/workspaces/${id}`, data),
  delete: (id: string) => api.delete(`/workspaces/${id}`),
  getMembers: (id: string) => api.get(`/workspaces/${id}/members`),
}

export const projectAPI = {
  list: (workspaceId: string) => api.get(`/projects?workspace_id=${workspaceId}`),
  create: (workspaceId: string, data: { name: string; key: string; description?: string }) =>
    api.post(`/projects?workspace_id=${workspaceId}`, data),
  get: (key: string) => api.get(`/projects/${key}`),
  update: (key: string, data: Partial<{ name: string; description: string; status: string }>) =>
    api.put(`/projects/${key}`, data),
  delete: (key: string) => api.delete(`/projects/${key}`),
  getMembers: (key: string) => api.get(`/projects/${key}/members`),
  addMember: (key: string, userId: string, role?: string) =>
    api.post(`/projects/${key}/members`, { user_id: userId, role }),
}

export const taskAPI = {
  listByProject: (projectKey: string) => api.get(`/projects/${projectKey}/tasks`),
  create: (projectKey: string, data: {
    title: string
    description?: string
    status?: string
    priority?: string
    assignee_id?: string
    due_date?: string
  }) => api.post(`/projects/${projectKey}/tasks`, data),
  get: (id: string) => api.get(`/tasks/${id}`),
  update: (id: string, data: Partial<{
    title: string
    description: string
    status: string
    priority: string
    assignee_id: string
    due_date: string
  }>) => api.put(`/tasks/${id}`, data),
  delete: (id: string) => api.delete(`/tasks/${id}`),
  move: (id: string, status: string, position: number) =>
    api.post(`/tasks/${id}/move`, { status, position }),
  getComments: (id: string) => api.get(`/tasks/${id}/comments`),
  addComment: (id: string, content: string) =>
    api.post(`/tasks/${id}/comments`, { content }),
  getActivity: (id: string) => api.get(`/tasks/${id}/activity`),
}

export const dashboardAPI = {
  getStats: (workspaceId: string) => api.get(`/dashboard/stats?workspace_id=${workspaceId}`),
  getBurndown: (projectKey: string) => api.get(`/dashboard/projects/${projectKey}/burndown`),
}

export default api
