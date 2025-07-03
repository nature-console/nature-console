import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname

  // Check if the route is an admin route (except login)
  if (pathname.startsWith('/admin') && !pathname.startsWith('/admin/login')) {
    // Check if user has auth token
    const token = request.cookies.get('token')
    
    if (!token) {
      // Redirect to login if no token
      return NextResponse.redirect(new URL('/admin/login', request.url))
    }
  }

  // If user is logged in and tries to access login page, redirect to dashboard
  if (pathname === '/admin/login') {
    const token = request.cookies.get('token')
    if (token) {
      return NextResponse.redirect(new URL('/admin/dashboard', request.url))
    }
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/admin/:path*']
}