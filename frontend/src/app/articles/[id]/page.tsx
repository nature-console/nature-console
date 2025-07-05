'use client'

import { useState, useEffect } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { Article } from '@/types/article'
import { articlesApi } from '@/lib/articles-api'

export default function ArticlePage() {
  const params = useParams()
  const [article, setArticle] = useState<Article | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchArticle = async () => {
      try {
        const id = parseInt(params.id as string)
        const data = await articlesApi.getArticle(id)
        setArticle(data)
      } catch (err) {
        setError('記事の取得に失敗しました')
        console.error('Error fetching article:', err)
      } finally {
        setLoading(false)
      }
    }

    if (params.id) {
      fetchArticle()
    }
  }, [params.id])

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">読み込み中...</div>
      </div>
    )
  }

  if (error || !article) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="text-xl text-red-600 mb-4">{error || '記事が見つかりません'}</div>
          <Link 
            href="/articles"
            className="inline-block bg-nature-green text-nature-text px-6 py-2 rounded-lg hover:bg-opacity-80 transition-all"
          >
            記事一覧に戻る
          </Link>
        </div>
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
              <Link href="/articles" className="hover:text-nature-green transition-colors">
                記事一覧
              </Link>
              <Link href="/admin/login" className="hover:text-nature-green transition-colors">
                管理画面
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {/* Article Content */}
      <article className="py-12">
        <div className="container mx-auto px-4 max-w-4xl">
          {/* Back Button */}
          <div className="mb-8">
            <Link 
              href="/articles"
              className="inline-flex items-center text-nature-green hover:text-nature-text transition-colors"
            >
              ← 記事一覧に戻る
            </Link>
          </div>

          {/* Article Header */}
          <header className="mb-8">
            <h1 className="text-4xl font-bold mb-4">{article.title}</h1>
            <div className="flex items-center space-x-4 text-gray-600">
              <span>著者: {article.author}</span>
              <span>•</span>
              <time>
                {new Date(article.created_at).toLocaleDateString('ja-JP', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric'
                })}
              </time>
            </div>
          </header>

          {/* Article Content */}
          <div className="prose prose-lg max-w-none">
            <div className="whitespace-pre-wrap leading-relaxed">
              {article.content}
            </div>
          </div>

          {/* Article Footer */}
          <footer className="mt-12 pt-8 border-t border-nature-beige">
            <div className="flex justify-center">
              <Link 
                href="/articles"
                className="bg-nature-green text-nature-text px-8 py-3 rounded-lg hover:bg-opacity-80 transition-all font-semibold"
              >
                他の記事を読む
              </Link>
            </div>
          </footer>
        </div>
      </article>

      {/* Footer */}
      <footer className="border-t border-nature-beige py-8">
        <div className="container mx-auto px-4 text-center text-gray-600">
          <p>&copy; 2024 Nature Console. All rights reserved.</p>
        </div>
      </footer>
    </div>
  )
}