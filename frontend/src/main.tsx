import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'
import { AuthProvider } from './lib/auth'
import { client } from './api/client.gen'
import './index.css'

if (import.meta.env.VITE_API_BASE_URL) {
  const rawBase = import.meta.env.VITE_API_BASE_URL.replace(/\/+$/, '')
  const baseUrl = rawBase.endsWith('/api/v1') ? rawBase : `${rawBase}/api/v1`
  client.setConfig({
    baseUrl,
  })
}

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 30_000, retry: 1 } },
})

const router = createRouter({ routeTree, defaultPreload: 'intent' })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <RouterProvider router={router} />
      </AuthProvider>
    </QueryClientProvider>
  </StrictMode>,
)
