import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { listBoqProjectsOptions } from '@/api/@tanstack/react-query.gen'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { InputGroup, InputGroupAddon, InputGroupInput } from '@/components/ui/input-group'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { Badge } from '@/components/ui/badge'
import { Download, Search } from 'lucide-react'

export const Route = createFileRoute('/boq/')({ component: BoqPage })

function BoqPage() {
  const navigate = useNavigate()
  const [search, setSearch] = useState('')

  const { data: projects = [], isLoading } = useQuery({
    ...listBoqProjectsOptions({ query: { q: search || undefined } }),
  })

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'complete':
        return <Badge variant="default" className="text-[11px] font-normal">Complete</Badge>
      case 'partial':
        return <Badge variant="secondary" className="text-[11px] font-normal">Partial</Badge>
      default:
        return <Badge variant="outline" className="text-[11px] font-normal text-muted-foreground">Empty</Badge>
    }
  }

  return (
    <div className="flex flex-col gap-4">
      {/* Header bar */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">BOQ (Bill of Quantities)</h1>
        </div>

        <div className="flex items-center gap-3">
          <InputGroup className="w-64">
            <InputGroupAddon>
              <Search data-icon="inline-start" />
            </InputGroupAddon>
            <InputGroupInput
              placeholder="Search projects..."
              onChange={(e) => setSearch(e.target.value)}
            />
          </InputGroup>
          <Button
            variant="outline"
            size="sm"
            onClick={() => navigate({ to: '/volume' })}
          >
            Volume Take-off
          </Button>
          <Button
            size="sm"
            onClick={() => navigate({ to: '/rab' })}
          >
            RAB
          </Button>
        </div>
      </div>

      <div className="rounded-md border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead rowSpan={2} className="w-12 text-center">No</TableHead>
              <TableHead rowSpan={2}>Project</TableHead>
              <TableHead colSpan={2} className="text-center border-b">Volume</TableHead>
              <TableHead colSpan={2} className="text-center border-b">RAB</TableHead>
            </TableRow>
            <TableRow>
              <TableHead className="w-28 text-center text-xs">Status</TableHead>
              <TableHead className="w-32 text-center text-xs">Export</TableHead>
              <TableHead className="w-28 text-center text-xs">Status</TableHead>
              <TableHead className="w-40 text-center text-xs">Export</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={6} className="h-32 text-center text-muted-foreground">
                  Loading data...
                </TableCell>
              </TableRow>
            ) : projects.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="h-48 text-center">
                  <Empty>
                    <EmptyHeader>
                      <EmptyTitle>Belum ada proyek BOQ</EmptyTitle>
                      <EmptyDescription>
                        {search ? `Tidak ada proyek yang cocok dengan "${search}".` : 'Belum ada data proyek.'}
                      </EmptyDescription>
                    </EmptyHeader>
                  </Empty>
                </TableCell>
              </TableRow>
            ) : (
              projects.map((p, idx) => (
                <TableRow key={p.id}>
                  <TableCell className="text-center text-muted-foreground">{idx + 1}</TableCell>
                  <TableCell className="font-medium text-foreground">
                    <button
                      type="button"
                      onClick={() => navigate({ to: `/projects` })}
                      className="text-left hover:underline"
                    >
                      {p.name}
                    </button>
                    <div className="text-xs text-muted-foreground">
                      {p.section || 'General'} · {p.plant || 'Site'}
                    </div>
                  </TableCell>

                  {/* Volume Status & Export */}
                  <TableCell className="text-center">{getStatusBadge(p.volume_status)}</TableCell>
                  <TableCell className="text-center">
                    <Button
                      size="sm"
                      variant="outline"
                      className="h-8 text-xs"
                      render={<a href={`/api/v1/projects/${p.id}/exports/volume.xlsx`} download />}
                      nativeButton={false}
                    >
                      <Download data-icon="inline-start" />
                      Excel
                    </Button>
                  </TableCell>

                  {/* RAB Status & Export */}
                  <TableCell className="text-center">{getStatusBadge(p.rab_status)}</TableCell>
                  <TableCell className="text-center">
                    <div className="flex items-center justify-center gap-1.5">
                      <Button
                        size="sm"
                        variant="outline"
                        className="h-8 text-xs"
                        render={<a href={`/api/v1/projects/${p.id}/exports/rab.xlsx`} download />}
                        nativeButton={false}
                      >
                        <Download data-icon="inline-start" />
                        Excel
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        className="h-8 text-xs"
                        render={<a href={`/api/v1/projects/${p.id}/exports/rab.pdf`} target="_blank" rel="noreferrer" />}
                        nativeButton={false}
                      >
                        <Download data-icon="inline-start" />
                        PDF
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
