import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listProjectsOptions, createProjectMutation, deleteProjectMutation, getCostEstimateOptions } from '@/api/@tanstack/react-query.gen'
import type { Project } from '@/api/types.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { formatRupiah } from '@/lib/utils'
import { useAuth } from '@/lib/auth'

export const Route = createFileRoute('/projects')({ component: ProjectsPage })

function ProjectsPage() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const { isAdmin, user, openLoginModal } = useAuth()

  const [search, setSearch] = useState('')
  const [page] = useState(1)
  const [isAddOpen, setIsAddOpen] = useState(false)
  const [selectedProject, setSelectedProject] = useState<Project | null>(null)

  // Form state
  const [name, setName] = useState('')
  const [section, setSection] = useState('')
  const [plant, setPlant] = useState('')
  const [planDate, setPlanDate] = useState('')
  const [budget, setBudget] = useState('')
  const [information, setInformation] = useState('')
  const [spb, setSpb] = useState(false)
  const [rab, setRab] = useState(false)
  const [xxx, setXxx] = useState(false)
  const [formError, setFormError] = useState('')

  const { data, isLoading } = useQuery({
    ...listProjectsOptions({
      query: { q: search || undefined, page, per_page: 25 },
    }),
  })

  const createMut = useMutation(createProjectMutation())
  const deleteMut = useMutation(deleteProjectMutation())

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal()
      return
    }
    setName('')
    setSection('')
    setPlant('')
    setPlanDate('')
    setBudget('')
    setInformation('')
    setSpb(false)
    setRab(false)
    setXxx(false)
    setFormError('')
    setIsAddOpen(true)
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setFormError('')
    if (!name.trim()) {
      setFormError('Nama project wajib diisi.')
      return
    }
    try {
      await createMut.mutateAsync({
        body: {
          name: name.trim(),
          section: section.trim(),
          plant: plant.trim(),
          implementation_plan: planDate || undefined,
          budget: budget ? budget.replace(/\D/g, '') : '0',
          information: information.trim(),
          progress_spb: spb ? 'done' : 'pending',
          progress_rab: rab ? 'done' : 'pending',
          progress_xxx: xxx ? 'done' : 'pending',
        },
      })
      await queryClient.invalidateQueries()
      setIsAddOpen(false)
    } catch (err: any) {
      setFormError(err?.message || 'Gagal membuat proyek.')
    }
  }

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation()
    if (!confirm('Yakin ingin menghapus proyek ini?')) return
    await deleteMut.mutateAsync({ path: { id } })
    await queryClient.invalidateQueries()
  }

  const projects = data?.data ?? []

  return (
    <div>
      <PageHeader title="PROJECT" onSearch={setSearch} onAdd={handleOpenAdd} />

      <div className="overflow-x-auto border bg-white shadow-sm">
        <Table className="text-xs">
          <TableHeader className="bg-white">
            <TableRow className="border-b hover:bg-transparent">
              <TableHead rowSpan={2} className="w-12 border-r text-center font-bold">NO</TableHead>
              <TableHead rowSpan={2} className="min-w-48 border-r text-center font-bold">PROJECT</TableHead>
              <TableHead rowSpan={2} className="min-w-28 border-r text-center font-bold">SECTION</TableHead>
              <TableHead rowSpan={2} className="min-w-28 border-r text-center font-bold">PLANT</TableHead>
              <TableHead rowSpan={2} className="min-w-32 border-r text-center font-bold">IMPLEMENTATION PLAN</TableHead>
              <TableHead rowSpan={2} className="min-w-36 border-r text-center font-bold">BUDGET</TableHead>
              <TableHead colSpan={3} className="border-b border-r text-center font-bold">PROGRESS</TableHead>
              <TableHead rowSpan={2} className="min-w-40 border-r text-center font-bold">INFORMATION</TableHead>
              {isAdmin && <TableHead rowSpan={2} className="w-20 text-center font-bold">AKSI</TableHead>}
            </TableRow>
            <TableRow className="border-b hover:bg-transparent">
              <TableHead className="w-16 border-r text-center font-bold">SPB</TableHead>
              <TableHead className="w-16 border-r text-center font-bold">RAB</TableHead>
              <TableHead className="w-16 border-r text-center font-bold">XXX</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 11 : 10} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  Memuat data...
                </TableCell>
              </TableRow>
            ) : projects.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 11 : 10} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  {search ? `Tidak ada proyek yang cocok dengan "${search}".` : 'Belum ada data proyek. Klik ADD untuk membuat.'}
                </TableCell>
              </TableRow>
            ) : (
              projects.map((p, idx) => (
                <TableRow
                  key={p.id}
                  onClick={() => setSelectedProject(p)}
                  className="cursor-pointer border-b bg-muted/40 transition-colors hover:bg-muted/80"
                >
                  <TableCell className="border-r text-center font-medium">{(page - 1) * 25 + idx + 1}</TableCell>
                  <TableCell className="border-r font-semibold text-foreground">{p.name}</TableCell>
                  <TableCell className="border-r">{p.section || '-'}</TableCell>
                  <TableCell className="border-r">{p.plant || '-'}</TableCell>
                  <TableCell className="border-r text-center">{p.implementation_plan || '-'}</TableCell>
                  <TableCell className="border-r text-right font-medium">{formatRupiah(p.budget)}</TableCell>
                  <TableCell className="border-r text-center">
                    {p.progress_spb === 'done' ? <Badge className="bg-emerald-600 text-[10px]">✓</Badge> : '-'}
                  </TableCell>
                  <TableCell className="border-r text-center">
                    {p.progress_rab === 'done' ? <Badge className="bg-emerald-600 text-[10px]">✓</Badge> : '-'}
                  </TableCell>
                  <TableCell className="border-r text-center">
                    {p.progress_xxx === 'done' ? <Badge className="bg-emerald-600 text-[10px]">✓</Badge> : '-'}
                  </TableCell>
                  <TableCell className="border-r">{p.information || '-'}</TableCell>
                  {isAdmin && (
                    <TableCell className="text-center" onClick={(e) => e.stopPropagation()}>
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-7 text-xs text-red-600 hover:bg-red-50 hover:text-red-700"
                        onClick={(e) => handleDelete(p.id, e)}
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

      {/* Modal Add Project (Mockup Hal. 3) */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-xl rounded-2xl">
          <DialogHeader>
            <DialogTitle className="text-lg font-bold tracking-wide">TAMBAH PROYEK</DialogTitle>
          </DialogHeader>

          {formError && (
            <div className="rounded-md border border-red-200 bg-red-50 p-3 text-xs text-red-600">
              {formError}
            </div>
          )}

          <form onSubmit={handleCreate} className="space-y-4 text-xs font-semibold">
            <div>
              <Label className="text-xs">PROJECT</Label>
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Nama Proyek..."
                className="mt-1"
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label className="text-xs">SECTION</Label>
                <Input
                  value={section}
                  onChange={(e) => setSection(e.target.value)}
                  placeholder="Unit / Bagian..."
                  className="mt-1"
                />
              </div>
              <div>
                <Label className="text-xs">PLANT</Label>
                <Input
                  value={plant}
                  onChange={(e) => setPlant(e.target.value)}
                  placeholder="Lokasi Plant / Pabrik..."
                  className="mt-1"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label className="text-xs">DATE</Label>
                <Input
                  type="date"
                  value={planDate}
                  onChange={(e) => setPlanDate(e.target.value)}
                  className="mt-1"
                />
              </div>
              <div>
                <Label className="text-xs">BUDGET</Label>
                <Input
                  value={budget}
                  onChange={(e) => setBudget(e.target.value)}
                  placeholder="Nominal Rupiah..."
                  className="mt-1"
                />
              </div>
            </div>

            <div>
              <Label className="text-xs">PROGRESS</Label>
              <div className="mt-2 flex items-center gap-6">
                <label className="flex items-center gap-2 cursor-pointer font-normal">
                  <input
                    type="checkbox"
                    checked={spb}
                    onChange={(e) => setSpb(e.target.checked)}
                    className="size-4 rounded text-primary"
                  />
                  <span>SPB</span>
                </label>
                <label className="flex items-center gap-2 cursor-pointer font-normal">
                  <input
                    type="checkbox"
                    checked={rab}
                    onChange={(e) => setRab(e.target.checked)}
                    className="size-4 rounded text-primary"
                  />
                  <span>RAB</span>
                </label>
                <label className="flex items-center gap-2 cursor-pointer font-normal">
                  <input
                    type="checkbox"
                    checked={xxx}
                    onChange={(e) => setXxx(e.target.checked)}
                    className="size-4 rounded text-primary"
                  />
                  <span>XXX</span>
                </label>
              </div>
            </div>

            <div>
              <Label className="text-xs">INFORMATION</Label>
              <Input
                value={information}
                onChange={(e) => setInformation(e.target.value)}
                placeholder="Keterangan..."
                className="mt-1"
              />
            </div>

            <div className="flex justify-end gap-2 pt-4">
              <Button type="button" variant="outline" onClick={() => setIsAddOpen(false)}>
                Batal
              </Button>
              <Button type="submit" disabled={createMut.isPending} className="px-8">
                {createMut.isPending ? 'Menyimpan...' : 'ADD'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Project Detail Dialog */}
      {selectedProject && (
        <ProjectDetailDialog
          project={selectedProject}
          open={!!selectedProject}
          onClose={() => setSelectedProject(null)}
          onNavigate={(path) => {
            setSelectedProject(null)
            navigate({ to: path })
          }}
        />
      )}
    </div>
  )
}

function ProjectDetailDialog({
  project,
  open,
  onClose,
  onNavigate,
}: {
  project: Project
  open: boolean
  onClose: () => void
  onNavigate: (to: string) => void
}) {
  const { data: cost } = useQuery({
    ...getCostEstimateOptions({ path: { id: project.id } }),
    enabled: open,
  })

  return (
    <Dialog open={open} onOpenChange={(val) => !val && onClose()}>
      <DialogContent className="max-w-2xl rounded-2xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-3 text-lg font-bold">
            <span className="bg-primary px-3 py-1 text-white">{project.name}</span>
            <span className="text-sm font-normal text-muted-foreground">{project.section} · {project.plant}</span>
          </DialogTitle>
        </DialogHeader>

        {/* Cost Estimate Summary Card (PRD §6.8) */}
        <div className="rounded-xl border bg-muted/40 p-4">
          <h3 className="mb-3 text-xs font-bold uppercase tracking-wider text-muted-foreground">
            Cost Estimate Summary
          </h3>
          <div className="grid grid-cols-3 gap-3 text-xs">
            <div className="rounded-lg bg-white p-3 shadow-xs">
              <div className="text-muted-foreground">BUDGET</div>
              <div className="mt-1 text-base font-bold text-foreground">{formatRupiah(cost?.budget ?? project.budget)}</div>
            </div>
            <div className="rounded-lg bg-white p-3 shadow-xs">
              <div className="text-muted-foreground">RAB (rencana)</div>
              <div className="mt-1 text-base font-bold text-foreground">{formatRupiah(cost?.rab ?? '0')}</div>
            </div>
            <div className="rounded-lg bg-white p-3 shadow-xs">
              <div className="text-muted-foreground">SPB (terbayar)</div>
              <div className="mt-1 text-base font-bold text-foreground">{formatRupiah(cost?.spb ?? '0')}</div>
            </div>
          </div>
          <div className="mt-3 grid grid-cols-2 gap-3 text-xs">
            <div className="rounded-lg bg-white p-3 shadow-xs">
              <div className="text-muted-foreground">VARIANCE (Budget - RAB)</div>
              <div className="mt-1 text-sm font-bold text-foreground">{formatRupiah(cost?.variance ?? '0')}</div>
            </div>
            <div className="rounded-lg bg-white p-3 shadow-xs">
              <div className="text-muted-foreground">SISA (Budget - SPB)</div>
              <div className="mt-1 text-sm font-bold text-foreground">{formatRupiah(cost?.sisa ?? '0')}</div>
            </div>
          </div>
        </div>

        <div className="flex items-center justify-between pt-2">
          <div className="flex gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={() => onNavigate(`/schedule?projectId=${project.id}`)}
            >
              Jadwal
            </Button>
            <Button
              size="sm"
              variant="outline"
              onClick={() => onNavigate(`/volume?projectId=${project.id}`)}
            >
              Volume Take-off
            </Button>
            <Button
              size="sm"
              className="bg-primary text-white"
              onClick={() => onNavigate(`/rab?projectId=${project.id}`)}
            >
              RAB Proyek
            </Button>
          </div>
          <Button variant="ghost" size="sm" onClick={onClose}>
            Tutup
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
