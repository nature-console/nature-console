'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Article } from '@/types/article'
import { articlesApi } from '@/lib/articles-api'

export default function ArticlesPage() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchArticles = async () => {
      try {
        const data = await articlesApi.getPublishedArticles()
        setArticles(data)
      } catch (err) {
        setError('記事の取得に失敗しました')
        console.error('Error fetching articles:', err)
      } finally {
        setLoading(false)
      }
    }

    fetchArticles()
  }, [])

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
    <div className="min-h-screen">
      {/* Header */}
      <header className="border-b border-nature-beige">
        <div className="container mx-auto px-4 py-6">
          <div className="flex justify-between items-center">
            <Link href="/" className="text-2xl font-bold hover:text-nature-green transition-colors">
              Nature Console
            </Link>
            <nav className="space-x-6">
              <Link href="/" className="hover:text-nature-green transition-colors">
                ホーム
              </Link>
              <Link href="/admin/login" className="hover:text-nature-green transition-colors">
                管理画面
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {/* Articles Section */}
      <section className="py-12">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl font-bold text-center mb-12">記事一覧</h1>
          
          {articles.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-xl text-gray-600">まだ記事が投稿されていません。</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {articles.map((article) => (
                <article
                  key={article.id}
                  className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
                >
                  <div className="p-6">
                    <h2 className="text-xl font-semibold mb-3 line-clamp-2">
                      <Link 
                        href={`/articles/${article.id}`}
                        className="hover:text-nature-green transition-colors"
                      >
                        {article.title}
                      </Link>
                    </h2>
                    <p className="text-gray-600 mb-4 line-clamp-3">
                      {article.content.substring(0, 150)}...
                    </p>
                    <div className="flex justify-between items-center text-sm text-gray-500">
                      <span>著者: {article.author}</span>
                      <span>
                        {new Date(article.created_at).toLocaleDateString('ja-JP')}
                      </span>
                    </div>
                  </div>
                </article>
              ))}
            </div>
          )}
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-nature-beige py-8">
        <div className="container mx-auto px-4 text-center text-gray-600">
          <p>&copy; 2024 Nature Console. All rights reserved.</p>
        </div>
      </footer>
    </div>
  )
}