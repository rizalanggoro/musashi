import { createRootRoute, Link, Outlet, useLocation } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { getHealthOptions } from '@/api/@tanstack/react-query.gen'
import { useAuth } from '@/lib/auth'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { TooltipProvider } from '@/components/ui/tooltip'
import {
  SidebarProvider,
  Sidebar,
  SidebarHeader,
  SidebarContent,
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
  LogIn,
  LogOut,
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
  const { user, openLoginModal, logout } = useAuth()
  const location = useLocation()

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

              {user ? (
                <div className="flex items-center gap-2 text-xs">
                  <span className="text-muted-foreground">
                    {user.name} <span className="text-[10px] uppercase font-semibold">({user.role})</span>
                  </span>
                  <Button size="sm" variant="outline" onClick={() => logout()}>
                    <LogOut data-icon="inline-start" />
                    Keluar
                  </Button>
                </div>
              ) : (
                <Button size="sm" variant="default" onClick={openLoginModal}>
                  <LogIn data-icon="inline-start" />
                  Masuk
                </Button>
              )}
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
