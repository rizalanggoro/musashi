import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listProjectsOptions,
  listVolumesOptions,
  deleteVolumeMutation,
} from '@/api/@tanstack/react-query.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button, buttonVariants } from '@/components/ui/button'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { FormulaBadge } from '@/components/formula-badge'
import { renderDimensions } from '@/lib/volume'
import { formatDecimal } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import { FileSpreadsheet } from 'lucide-react'

export const Route = createFileRoute('/volume/')({
  component: VolumePage,
  validateSearch: (search: Record<string, unknown>): { projectId?: string } => ({
    projectId: (search.projectId as string) || undefined,
  }),
})

function VolumePage() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const searchParams = Route.useSearch()
  const queryProjectId = searchParams.projectId || ''
  const { isAdmin, user, openLoginModal } = useAuth()

  const [selectedProjectId, setSelectedProjectId] = useState<string>(queryProjectId)
  const [search, setSearch] = useState('')

  // Fetch projects for selector
  const { data: projectsData } = useQuery({
    ...listProjectsOptions({ query: { per_page: 100 } }),
  })
  const projects = projectsData?.data ?? []
  const projectItems = projects.map((p) => ({
    value: p.id,
    label: `${p.name} ${p.section ? `(${p.section})` : ''}`,
  }))

  // Default project selection if not set
  const activeProjectId = selectedProjectId || (projects.length > 0 ? projects[0].id : '')
  const activeProject = projects.find((p) => p.id === activeProjectId)

  // Fetch volume rows for active project
  const { data: volumeRows = [], isLoading } = useQuery({
    ...listVolumesOptions({
      path: { id: activeProjectId },
      query: { q: search || undefined },
    }),
    enabled: !!activeProjectId,
  })

  const deleteMut = useMutation(deleteVolumeMutation())

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal()
      return
    }
    if (!activeProjectId) {
      alert('Pilih proyek terlebih dahulu.')
      return
    }
    navigate({
      to: '/volume/new',
      search: { projectId: activeProjectId },
    })
  }

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(`Hapus baris take-off "${name}"?`)) return
    try {
      await deleteMut.mutateAsync({ path: { id } })
      await queryClient.invalidateQueries()
    } catch (err: any) {
      alert(err?.message || 'Gagal menghapus volume.')
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <PageHeader
        title="Volume Take-off"
        onSearch={setSearch}
        onAdd={handleOpenAdd}
        addLabel="Tambah Take-off"
      />

      {/* Project Selector & Status Bar */}
      <div className="flex flex-wrap items-center justify-between gap-4 rounded-lg border bg-card p-3 shadow-xs">
        <div className="flex items-center gap-3">
          <span className="text-sm font-semibold text-foreground whitespace-nowrap">Proyek:</span>
          <div className="w-72">
            <Select
              items={projectItems}
              value={activeProjectId}
              onValueChange={(val) => {
                if (val) setSelectedProjectId(val)
              }}
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Pilih Proyek" />
              </SelectTrigger>
              <SelectContent className="w-80">
                <SelectGroup>
                  {projectItems.map((item) => (
                    <SelectItem key={item.value} value={item.value}>
                      {item.label}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>

          {activeProject && (
            <div className="hidden sm:flex items-center gap-2">
              <Badge variant="secondary" className="text-xs">
                {activeProject.plant}
              </Badge>
              <Badge variant="outline" className="text-xs">
                {activeProject.section}
              </Badge>
            </div>
          )}
        </div>

        {/* Action Export */}
        {activeProjectId && (
          <a
            href={`/api/v1/projects/${activeProjectId}/exports/volume.xlsx`}
            download
            className={buttonVariants({
              variant: 'outline',
              size: 'sm',
              className: 'gap-2 text-xs',
            })}
          >
            <FileSpreadsheet className="size-3.5 text-emerald-600" />
            <span>Export Excel</span>
          </a>
        )}
      </div>

      {/* Table Data Take-off */}
      <div className="rounded-md border bg-card overflow-x-auto shadow-xs">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-10 text-center">No</TableHead>
              <TableHead className="min-w-44">Work Item</TableHead>
              <TableHead className="w-28 text-center">Metode / Rumus</TableHead>
              <TableHead className="w-24 text-right">Planned Vol</TableHead>
              <TableHead className="w-16 text-center">Unit</TableHead>
              <TableHead className="w-24 text-center">Date</TableHead>
              <TableHead className="w-40 text-center">Dimensi Pengukuran</TableHead>
              <TableHead className="w-24 text-right">Quantity</TableHead>
              <TableHead className="w-24 text-right">Unit / Wt</TableHead>
              <TableHead className="w-28 text-right font-bold">Total Qty</TableHead>
              <TableHead className="min-w-32">Description</TableHead>
              {isAdmin && <TableHead className="w-16 text-center">Action</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 12 : 11} className="h-32 text-center text-muted-foreground">
                  Memuat data volume take-off...
                </TableCell>
              </TableRow>
            ) : volumeRows.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 12 : 11} className="h-48 text-center">
                  <Empty>
                    <EmptyHeader>
                      <EmptyTitle>Tidak ada data take-off volume</EmptyTitle>
                      <EmptyDescription>
                        {search ? `Tidak ada data cocok dengan "${search}".` : 'Mulai tambahkan perhitungan volume take-off untuk proyek ini.'}
                      </EmptyDescription>
                    </EmptyHeader>
                  </Empty>
                </TableCell>
              </TableRow>
            ) : (
              volumeRows.map((r, idx) => (
                <TableRow key={r.id}>
                  <TableCell className="text-center text-muted-foreground">{idx + 1}</TableCell>
                  <TableCell className="font-medium text-foreground">{r.work_item}</TableCell>
                  <TableCell className="text-center">
                    <FormulaBadge formula={r.formula} />
                  </TableCell>
                  <TableCell className="text-right text-muted-foreground">{formatDecimal(r.planned_volume)}</TableCell>
                  <TableCell className="text-center text-muted-foreground font-mono text-xs">{r.unit}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{r.date || '-'}</TableCell>
                  <TableCell className="text-center text-muted-foreground text-xs font-mono">
                    {renderDimensions(r)}
                  </TableCell>
                  <TableCell className="text-right text-foreground font-medium">{formatDecimal(r.quantity)}</TableCell>
                  <TableCell className="text-right text-muted-foreground">{formatDecimal(r.unit_weight)}</TableCell>
                  <TableCell className="text-right font-bold text-foreground text-emerald-600 dark:text-emerald-400">
                    {formatDecimal(r.total_quantity)}
                  </TableCell>
                  <TableCell className="text-muted-foreground text-xs">{r.description || '-'}</TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="text-xs text-destructive hover:bg-destructive/10"
                        onClick={() => handleDelete(r.id, r.work_item)}
                      >
                        Hapus
                      </Button>
                    </TableCell>
                  )}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
