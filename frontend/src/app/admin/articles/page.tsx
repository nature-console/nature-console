'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useAuth } from '@/contexts/AuthContext'
import { adminApi } from '@/lib/admin-api'
import { Article } from '@/types/article'

export default function ArticleManagementPage() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  
  const { user, logout } = useAuth()

  const fetchArticles = async () => {
    try {
      const data = await adminApi.getArticles()
      setArticles(data)
    } catch (err) {
      setError('記事の取得に失敗しました')
      console.error('Error fetching articles:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchArticles()
  }, [])

  const handleDelete = async (id: number) => {
    if (!confirm('この記事を削除してもよろしいですか？')) {
      return
    }

    try {
      await adminApi.deleteArticle(id)
      setArticles(articles.filter(article => article.id !== id))
    } catch (err) {
      console.error('Error deleting article:', err)
      alert('記事の削除に失敗しました')
    }
  }

  const handleLogout = async () => {
    try {
      await logout()
    } catch (err) {
      console.error('Logout error:', err)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">読み込み中...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200">
        <div className="container mx-auto px-4 py-4">
          <div className="flex justify-between items-center">
            <div className="flex items-center space-x-4">
              <Link href="/" className="text-xl font-bold text-nature-text hover:text-nature-green transition-colors">
                Nature Console
              </Link>
              <span className="text-gray-500">管理画面</span>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-600">こんにちは、{user?.name}さん</span>
              <button
                onClick={handleLogout}
                className="text-sm text-red-600 hover:text-red-800 transition-colors"
              >
                ログアウト
              </button>
            </div>
          </div>
        </div>
      </header>

      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <div>
            <Link 
              href="/admin/dashboard"
              className="inline-flex items-center text-nature-green hover:text-nature-text transition-colors mb-2"
            >
              ← ダッシュボードに戻る
            </Link>
            <h1 className="text-3xl font-bold text-gray-900">記事管理</h1>
          </div>
          <Link
            href="/admin/articles/new"
            className="bg-nature-green text-nature-text px-6 py-2 rounded-lg hover:bg-opacity-80 transition-all"
          >
            新しい記事を作成
          </Link>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg mb-6">
            {error}
          </div>
        )}

        <div className="bg-white rounded-lg shadow">
          {articles.length === 0 ? (
            <div className="p-8 text-center">
              <p className="text-gray-600 mb-4">まだ記事がありません。</p>
              <Link
                href="/admin/articles/new"
                className="inline-block bg-nature-green text-nature-text px-6 py-2 rounded-lg hover:bg-opacity-80 transition-all"
              >
                最初の記事を作成
              </Link>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      タイトル
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      著者
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      ステータス
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      作成日
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      操作
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {articles.map((article) => (
                    <tr key={article.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4">
                        <div className="text-sm font-medium text-gray-900">
                          {article.title}
                        </div>
                        <div className="text-sm text-gray-500 truncate max-w-xs">
                          {article.content.substring(0, 100)}...
                        </div>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-900">
                        {article.author}
                      </td>
                      <td className="px-6 py-4">
                        <span
                          className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                            article.published 
                              ? 'bg-green-100 text-green-800' 
                              : 'bg-gray-100 text-gray-800'
                          }`}
                        >
                          {article.published ? '公開中' : '下書き'}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-500">
                        {new Date(article.created_at).toLocaleDateString('ja-JP')}
                      </td>
                      <td className="px-6 py-4 text-sm space-x-2">
                        {article.published && (
                          <Link
                            href={`/articles/${article.id}`}
                            target="_blank"
                            className="text-nature-green hover:text-nature-text transition-colors"
                          >
                            表示
                          </Link>
                        )}
                        <Link
                          href={`/admin/articles/${article.id}/edit`}
                          className="text-blue-600 hover:text-blue-800 transition-colors"
                        >
                          編集
                        </Link>
                        <button
                          onClick={() => handleDelete(article.id)}
                          className="text-red-600 hover:text-red-800 transition-colors"
                        >
                          削除
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}