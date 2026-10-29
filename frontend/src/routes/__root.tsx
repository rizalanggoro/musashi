import { useEffect } from 'react'
import { createRootRoute, Link, Outlet, useLocation, useNavigate } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { getHealthOptions } from '@/api/@tanstack/react-query.gen'
import { useAuth } from '@/lib/auth'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { TooltipProvider } from '@/components/ui/tooltip'
import {
  SidebarProvider,
  Sidebar,
  SidebarHeader,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarTrigger,
  SidebarInset,
} from '@/components/ui/sidebar'
import {
  FolderKanban,
  Calendar,
  Tag,
  Calculator,
  FileSpreadsheet,
  Box,
  Receipt,
  LogOut,
  User as UserIcon,
} from 'lucide-react'

const nav = [
  {
    group: 'Monitoring',
    items: [
      { to: '/projects', label: 'Projects', icon: FolderKanban },
      { to: '/schedule', label: 'Schedule', icon: Calendar },
    ],
  },
  {
    group: 'Cost Estimate',
    items: [
      { to: '/unit-prices', label: 'Unit Prices', icon: Tag },
      { to: '/ahsp', label: 'AHSP', icon: Calculator },
      { to: '/boq', label: 'BOQ', icon: FileSpreadsheet },
      { to: '/volume', label: 'Volume Take-off', icon: Box },
      { to: '/rab', label: 'RAB', icon: Receipt },
    ],
  },
] as const

export const Route = createRootRoute({
  component: Shell,
})

function Shell() {
  const { user, isLoading, logout } = useAuth()
  const location = useLocation()
  const navigate = useNavigate()

  const isAuthPage = location.pathname === '/auth' || location.pathname === '/login'

  useEffect(() => {
    if (!isLoading) {
      if (!user && !isAuthPage) {
        navigate({ to: '/auth', replace: true })
      } else if (user && isAuthPage) {
        navigate({ to: '/projects', replace: true })
      }
    }
  }, [user, isLoading, isAuthPage, navigate])

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="flex flex-col items-center gap-3">
          <div className="size-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
          <p className="text-xs text-muted-foreground">Memuat sesi...</p>
        </div>
      </div>
    )
  }

  if (isAuthPage) {
    return <Outlet />
  }

  if (!user) {
    return null
  }

  return (
    <TooltipProvider>
      <SidebarProvider defaultOpen>
        <Sidebar variant="inset">
          <SidebarHeader className="h-14 justify-center border-b px-4">
            <Link to="/projects" className="flex items-center gap-2 font-bold tracking-tight text-lg text-foreground">
              <span>MUSASHI</span>
            </Link>
          </SidebarHeader>

          <SidebarContent>
            {nav.map((section) => (
              <SidebarGroup key={section.group}>
                <SidebarGroupLabel>{section.group}</SidebarGroupLabel>
                <SidebarGroupContent>
                  <SidebarMenu>
                    {section.items.map((item) => {
                      const Icon = item.icon
                      const isActive = location.pathname === item.to
                      return (
                        <SidebarMenuItem key={item.to}>
                          <SidebarMenuButton
                            isActive={isActive}
                            render={<Link to={item.to} />}
                          >
                            <Icon />
                            <span>{item.label}</span>
                          </SidebarMenuButton>
                        </SidebarMenuItem>
                      )
                    })}
                  </SidebarMenu>
                </SidebarGroupContent>
              </SidebarGroup>
            ))}
          </SidebarContent>

          <SidebarFooter className="border-t p-2">
            <SidebarMenu>
              <SidebarMenuItem>
                <div className="flex items-center gap-2.5 px-2 py-1.5">
                  <div className="flex size-7 shrink-0 items-center justify-center rounded-md bg-muted text-muted-foreground">
                    <UserIcon className="size-4" />
                  </div>
                  <div className="flex flex-col min-w-0 flex-1">
                    <span className="truncate text-xs font-medium text-foreground">{user.name}</span>
                    <span className="truncate text-[10px] uppercase font-semibold text-muted-foreground">{user.role}</span>
                  </div>
                </div>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton
                  onClick={() => logout()}
                  tooltip="Keluar"
                  className="text-destructive hover:bg-destructive/10 hover:text-destructive focus:text-destructive"
                >
                  <LogOut />
                  <span>Keluar</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarFooter>
        </Sidebar>

        <SidebarInset>
          {/* Top Navbar */}
          <header className="flex h-14 shrink-0 items-center justify-between gap-2 border-b bg-background px-4">
            <div className="flex items-center gap-2">
              <SidebarTrigger className="-ml-1" />
              <Separator orientation="vertical" className="mr-2 h-4" />
            </div>

            <div className="flex items-center gap-3">
              <HealthBadge />
            </div>
          </header>

          {/* Main Page Area */}
          <main className="flex-1 p-6">
            <Outlet />
          </main>
        </SidebarInset>
      </SidebarProvider>
    </TooltipProvider>
  )
}

function HealthBadge() {
  const { data, isError } = useQuery({ ...getHealthOptions(), refetchInterval: 30_000 })
  const ok = !isError && data?.status === 'ok'

  return (
    <Badge variant={ok ? 'secondary' : 'destructive'} className="text-[11px] gap-1.5 font-normal">
      <span className={`size-1.5 rounded-full ${ok ? 'bg-emerald-500' : 'bg-red-500'}`} />
      {ok ? 'Online' : 'Offline'}
    </Badge>
  )
}
