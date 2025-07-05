'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useAuth } from '@/contexts/AuthContext'
import { adminApi, DashboardData } from '@/lib/admin-api'

export default function DashboardPage() {
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  
  const { user, logout } = useAuth()

  useEffect(() => {
    const fetchDashboard = async () => {
      try {
        const data = await adminApi.getDashboard()
        setDashboardData(data)
      } catch (err) {
        setError('ダッシュボードデータの取得に失敗しました')
        console.error('Error fetching dashboard:', err)
      } finally {
        setLoading(false)
      }
    }

    fetchDashboard()
  }, [])

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

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl text-red-600">{error}</div>
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
        <h1 className="text-3xl font-bold text-gray-900 mb-8">ダッシュボード</h1>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="p-2 bg-nature-green rounded-lg">
                <span className="text-xl">📝</span>
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">総記事数</p>
                <p className="text-2xl font-bold text-gray-900">
                  {dashboardData?.stats.total_articles || 0}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="p-2 bg-nature-orange rounded-lg">
                <span className="text-xl">🌟</span>
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">公開記事</p>
                <p className="text-2xl font-bold text-gray-900">
                  {dashboardData?.stats.published_articles || 0}
                </p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center">
              <div className="p-2 bg-nature-purple rounded-lg">
                <span className="text-xl">📄</span>
              </div>
              <div className="ml-4">
                <p className="text-sm font-medium text-gray-600">下書き</p>
                <p className="text-2xl font-bold text-gray-900">
                  {dashboardData?.stats.draft_articles || 0}
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">クイックアクション</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Link
              href="/admin/articles/new"
              className="flex items-center p-4 bg-nature-green rounded-lg text-nature-text hover:bg-opacity-80 transition-all"
            >
              <span className="text-2xl mr-3">✏️</span>
              <div>
                <h3 className="font-semibold">新しい記事を書く</h3>
                <p className="text-sm opacity-80">記事を作成して公開しましょう</p>
              </div>
            </Link>

            <Link
              href="/admin/articles"
              className="flex items-center p-4 bg-nature-beige rounded-lg text-nature-text hover:bg-opacity-80 transition-all"
            >
              <span className="text-2xl mr-3">📋</span>
              <div>
                <h3 className="font-semibold">記事を管理</h3>
                <p className="text-sm opacity-80">既存の記事を編集・削除</p>
              </div>
            </Link>
          </div>
        </div>

        {/* Recent Articles */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">最近の記事</h2>
          {dashboardData?.recent_articles && dashboardData.recent_articles.length > 0 ? (
            <div className="space-y-4">
              {dashboardData.recent_articles.map((article) => (
                <div key={article.id} className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                  <div>
                    <h3 className="font-medium text-gray-900">{article.title}</h3>
                    <p className="text-sm text-gray-600">
                      {article.author} • {new Date(article.created_at).toLocaleDateString('ja-JP')}
                    </p>
                  </div>
                  <div className="flex items-center space-x-2">
                    <span
                      className={`px-2 py-1 text-xs rounded-full ${
                        article.published 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-gray-100 text-gray-800'
                      }`}
                    >
                      {article.published ? '公開中' : '下書き'}
                    </span>
                    <Link
                      href={`/admin/articles/${article.id}/edit`}
                      className="text-nature-green hover:text-nature-text transition-colors text-sm"
                    >
                      編集
                    </Link>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-600">まだ記事がありません。</p>
          )}
        </div>
      </div>
    </div>
  )
}