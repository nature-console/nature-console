import { apiClient } from './api'
import { Article } from '@/types/article'

export const articlesApi = {
  getPublishedArticles: async (): Promise<Article[]> => {
    const response = await apiClient.get('/articles')
    return response.data
  },

  getArticle: async (id: number): Promise<Article> => {
    const response = await apiClient.get(`/articles/${id}`)
    return response.data
  },
}