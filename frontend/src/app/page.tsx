import Link from 'next/link'

export default function Home() {
  return (
    <div className="min-h-screen">
      {/* Header */}
      <header className="border-b border-nature-beige">
        <div className="container mx-auto px-4 py-6">
          <div className="flex justify-between items-center">
            <h1 className="text-2xl font-bold">Nature Console</h1>
            <nav className="space-x-6">
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

      {/* Hero Section */}
      <section className="py-20">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-5xl font-bold mb-6">自然と技術の調和</h2>
          <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
            Nature Consoleは自然と技術の調和を目指すプラットフォームです。<br />
            環境に関する知識や技術的な洞察を共有し、持続可能な未来の構築に貢献します。
          </p>
          <Link 
            href="/articles"
            className="inline-block bg-nature-green text-nature-text px-8 py-3 rounded-lg hover:bg-opacity-80 transition-all font-semibold"
          >
            記事を読む
          </Link>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 bg-nature-beige">
        <div className="container mx-auto px-4">
          <h3 className="text-3xl font-bold text-center mb-12">私たちの使命</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="text-center">
              <div className="w-16 h-16 bg-nature-green rounded-full mx-auto mb-4 flex items-center justify-center">
                <span className="text-2xl">🌱</span>
              </div>
              <h4 className="text-xl font-semibold mb-2">持続可能性</h4>
              <p className="text-gray-600">
                環境に配慮した技術と手法を探求し、持続可能な未来を構築します。
              </p>
            </div>
            <div className="text-center">
              <div className="w-16 h-16 bg-nature-orange rounded-full mx-auto mb-4 flex items-center justify-center">
                <span className="text-2xl">💡</span>
              </div>
              <h4 className="text-xl font-semibold mb-2">技術革新</h4>
              <p className="text-gray-600">
                最新の技術を活用して、環境問題の解決策を見つけ出します。
              </p>
            </div>
            <div className="text-center">
              <div className="w-16 h-16 bg-nature-purple rounded-full mx-auto mb-4 flex items-center justify-center">
                <span className="text-2xl">🤝</span>
              </div>
              <h4 className="text-xl font-semibold mb-2">コミュニティ</h4>
              <p className="text-gray-600">
                知識を共有し、共に学び成長するコミュニティを築きます。
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-16">
        <div className="container mx-auto px-4 text-center">
          <h3 className="text-3xl font-bold mb-6">今すぐ始めましょう</h3>
          <p className="text-xl text-gray-600 mb-8">
            自然と技術の調和に向けた取り組みに参加しませんか？
          </p>
          <Link 
            href="/articles"
            className="inline-block bg-nature-text text-nature-bg px-8 py-3 rounded-lg hover:bg-opacity-80 transition-all font-semibold"
          >
            記事を探索する
          </Link>
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