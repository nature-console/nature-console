import type { Metadata } from 'next'
import { AuthProvider } from '@/contexts/AuthContext'
import './globals.css'

export const metadata: Metadata = {
  title: 'Nature Console',
  description: 'Nature Consoleは自然と技術の調和を目指すプラットフォームです。環境に関する知識や技術的な洞察を共有し、持続可能な未来の構築に貢献します。',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body className="antialiased bg-nature-bg text-nature-text">
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  )
}