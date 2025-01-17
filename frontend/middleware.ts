import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const PUBLIC_PATHS = ['/login', '/registration']
const STATIC_PATHS = /^\/_next\/static\//

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  // Исключаем статические файлы из обработки
  if (STATIC_PATHS.test(pathname)) {
    return NextResponse.next()
  }

  // Разрешаем публичные страницы
  if (PUBLIC_PATHS.includes(pathname)) {
    return NextResponse.next()
  }

  // Проверяем наличие токена
  const token = request.cookies.get('accToken')
  if (!token) {
    const loginUrl = new URL('/login', request.url)
    return NextResponse.redirect(loginUrl)
  }

  // Если токен есть, разрешаем доступ
  return NextResponse.next()
}

export const config = {
  matcher: '/:path*',
}
