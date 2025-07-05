import { apiClient } from './api'
import { Article } from '@/types/article'

export interface DashboardStats {
  total_articles: number
  published_articles: number
  draft_articles: number
}

export interface DashboardData {
  stats: DashboardStats
  recent_articles: Article[]
}

export interface CreateArticleData {
  title: string
  content: string
  author: string
  published?: boolean
}

export interface UpdateArticleData {
  title?: string
  content?: string
  author?: string
  published?: boolean
}

export const adminApi = {
  // Dashboard
  getDashboard: async (): Promise<DashboardData> => {
    const response = await apiClient.get('/admin/dashboard')
    return response.data
  },

  // Articles
  getArticles: async (): Promise<Article[]> => {
    const response = await apiClient.get('/admin/articles')
    return response.data
  },

  getArticle: async (id: number): Promise<Article> => {
    const response = await apiClient.get(`/admin/articles/${id}`)
    return response.data
  },

  createArticle: async (data: CreateArticleData): Promise<Article> => {
    const response = await apiClient.post('/admin/articles', data)
    return response.data
  },

  updateArticle: async (id: number, data: UpdateArticleData): Promise<Article> => {
    const response = await apiClient.put(`/admin/articles/${id}`, data)
    return response.data
  },

  deleteArticle: async (id: number): Promise<void> => {
    await apiClient.delete(`/admin/articles/${id}`)
  },
}