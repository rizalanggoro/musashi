import { useState } from 'react'
import { createRootRoute, Link, Outlet } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { getHealthOptions } from '@/api/@tanstack/react-query.gen'
import { cn } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import { Button } from '@/components/ui/button'

const nav = [
  {
    group: 'MONITORING',
    items: [
      { to: '/projects', label: 'PROJECT' },
      { to: '/schedule', label: 'SCHEDULE' },
    ],
  },
  {
    group: 'COST ESTIMATE',
    items: [
      { to: '/unit-prices', label: 'UNIT PRICE' },
      { to: '/ahsp', label: 'AHSP' },
      { to: '/boq', label: 'BOQ' },
      { to: '/volume', label: 'VOLUME' },
      { to: '/rab', label: 'RAB' },
    ],
  },
] as const

export const Route = createRootRoute({
  component: Shell,
})

function Shell() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const { user, openLoginModal, logout } = useAuth()

  return (
    <div className="flex min-h-screen flex-col bg-background font-sans">
      <header className="flex h-14 items-center gap-4 bg-secondary px-4">
        <button
          type="button"
          aria-label="Toggle sidebar"
          onClick={() => setSidebarOpen((prev) => !prev)}
          className="flex size-9 items-center justify-center rounded-md hover:bg-black/5"
        >
          <span className="flex flex-col gap-1">
            <span className="block h-0.5 w-6 bg-foreground" />
            <span className="block h-0.5 w-6 bg-foreground" />
            <span className="block h-0.5 w-6 bg-foreground" />
          </span>
        </button>
        <Link to="/projects" className="flex items-center gap-2">
          <Logo />
        </Link>
        <div className="ml-auto flex items-center gap-3">
          <HealthBadge />
          {user ? (
            <div className="flex items-center gap-2 text-xs">
              <span className="rounded bg-black/10 px-2 py-1 font-medium">
                {user.name} ({user.role})
              </span>
              <Button size="sm" variant="ghost" className="h-7 text-xs" onClick={() => logout()}>
                Keluar
              </Button>
            </div>
          ) : (
            <Button size="sm" variant="outline" className="h-7 text-xs" onClick={openLoginModal}>
              Masuk
            </Button>
          )}
        </div>
      </header>

      <div className="flex flex-1">
        {sidebarOpen && (
          <aside className="w-64 shrink-0 border-r bg-white">
            {nav.map((section) => (
              <div key={section.group}>
                <div className="bg-primary px-4 py-2.5 text-sm font-semibold tracking-wider text-primary-foreground">
                  {section.group}
                </div>
                {section.items.map((item) => (
                  <Link
                    key={item.to}
                    to={item.to}
                    className="block border-b px-4 py-3 text-sm font-medium tracking-wide text-foreground transition-colors hover:bg-muted/70"
                    activeProps={{
                      className: 'bg-muted font-bold text-foreground border-l-4 border-primary pl-3',
                    }}
                  >
                    {item.label}
                  </Link>
                ))}
              </div>
            ))}
          </aside>
        )}

        <main className="min-w-0 flex-1 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

function Logo() {
  return (
    <span className="flex items-center gap-2">
      <svg viewBox="0 0 40 28" className="h-6 w-9 text-primary" aria-hidden>
        <ellipse cx="20" cy="14" rx="18" ry="6.5" fill="none" stroke="currentColor" strokeWidth="2.5" />
        <circle cx="20" cy="14" r="8" fill="none" stroke="currentColor" strokeWidth="2.5" />
      </svg>
      <span className="text-xl font-extrabold tracking-tight">MUSASHI</span>
    </span>
  )
}

function HealthBadge() {
  const { data, isError } = useQuery({ ...getHealthOptions(), refetchInterval: 30_000 })
  const ok = !isError && data?.status === 'ok'

  return (
    <span
      className={cn(
        'flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium',
        ok ? 'border-emerald-300 bg-emerald-50 text-emerald-700' : 'border-red-300 bg-red-50 text-red-700',
      )}
    >
      <span className={cn('size-2 rounded-full', ok ? 'bg-emerald-500' : 'bg-red-500')} />
      {ok ? 'API OK' : 'API OFF'}
    </span>
  )
}
