import { apiClient } from './api'
import { AdminUser, LoginCredentials } from '@/types/auth'

export const authApi = {
  login: async (credentials: LoginCredentials) => {
    const response = await apiClient.post('/auth/login', credentials)
    return response.data
  },

  logout: async () => {
    const response = await apiClient.post('/auth/logout')
    return response.data
  },

  me: async (): Promise<{ user: AdminUser }> => {
    const response = await apiClient.get('/auth/me')
    return response.data
  },
}