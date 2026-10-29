import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listProjectsOptions, createProjectMutation, deleteProjectMutation, getCostEstimateOptions } from '@/api/@tanstack/react-query.gen'
import type { Project } from '@/api/types.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FieldGroup, Field, FieldLabel } from '@/components/ui/field'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { formatRupiah } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import { Check } from 'lucide-react'

export const Route = createFileRoute('/projects/')({ component: ProjectsPage })

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
    <div className="flex flex-col gap-4">
      <PageHeader title="Projects" onSearch={setSearch} onAdd={handleOpenAdd} addLabel="Project" />

      <div className="rounded-md border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead rowSpan={2} className="w-12 text-center">No</TableHead>
              <TableHead rowSpan={2}>Project</TableHead>
              <TableHead rowSpan={2}>Section</TableHead>
              <TableHead rowSpan={2}>Plant</TableHead>
              <TableHead rowSpan={2}>Plan Date</TableHead>
              <TableHead rowSpan={2} className="text-right">Budget</TableHead>
              <TableHead colSpan={3} className="text-center border-b">Progress</TableHead>
              <TableHead rowSpan={2}>Information</TableHead>
              {isAdmin && <TableHead rowSpan={2} className="w-20 text-center">Action</TableHead>}
            </TableRow>
            <TableRow>
              <TableHead className="w-16 text-center text-xs">SPB</TableHead>
              <TableHead className="w-16 text-center text-xs">RAB</TableHead>
              <TableHead className="w-16 text-center text-xs">XXX</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 11 : 10} className="h-32 text-center text-muted-foreground">
                  Loading data...
                </TableCell>
              </TableRow>
            ) : projects.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 11 : 10} className="h-48 text-center">
                  <Empty>
                    <EmptyHeader>
                      <EmptyTitle>Belum ada proyek</EmptyTitle>
                      <EmptyDescription>
                        {search ? `Tidak ada proyek yang cocok dengan "${search}".` : 'Mulai dengan menambahkan proyek baru.'}
                      </EmptyDescription>
                    </EmptyHeader>
                  </Empty>
                </TableCell>
              </TableRow>
            ) : (
              projects.map((p, idx) => (
                <TableRow
                  key={p.id}
                  onClick={() => setSelectedProject(p)}
                  className="cursor-pointer"
                >
                  <TableCell className="text-center text-muted-foreground">{(page - 1) * 25 + idx + 1}</TableCell>
                  <TableCell className="font-medium text-foreground">{p.name}</TableCell>
                  <TableCell>{p.section || '-'}</TableCell>
                  <TableCell>{p.plant || '-'}</TableCell>
                  <TableCell>{p.implementation_plan || '-'}</TableCell>
                  <TableCell className="text-right font-medium">{formatRupiah(p.budget)}</TableCell>
                  <TableCell className="text-center">
                    {p.progress_spb === 'done' ? (
                      <Badge variant="secondary" className="size-5 p-0 justify-center">
                        <Check className="size-3" />
                      </Badge>
                    ) : (
                      <span className="text-muted-foreground">-</span>
                    )}
                  </TableCell>
                  <TableCell className="text-center">
                    {p.progress_rab === 'done' ? (
                      <Badge variant="secondary" className="size-5 p-0 justify-center">
                        <Check className="size-3" />
                      </Badge>
                    ) : (
                      <span className="text-muted-foreground">-</span>
                    )}
                  </TableCell>
                  <TableCell className="text-center">
                    {p.progress_xxx === 'done' ? (
                      <Badge variant="secondary" className="size-5 p-0 justify-center">
                        <Check className="size-3" />
                      </Badge>
                    ) : (
                      <span className="text-muted-foreground">-</span>
                    )}
                  </TableCell>
                  <TableCell className="max-w-xs truncate text-muted-foreground">{p.information || '-'}</TableCell>
                  {isAdmin && (
                    <TableCell className="text-center" onClick={(e) => e.stopPropagation()}>
                      <Button
                        size="sm"
                        variant="ghost"
                        className="text-xs text-destructive hover:bg-destructive/10"
                        onClick={(e) => handleDelete(p.id, e)}
                      >
                        Delete
                      </Button>
                    </TableCell>
                  )}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Modal Add Project */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Add New Project</DialogTitle>
            <DialogDescription>Masukkan rincian proyek baru untuk monitoring.</DialogDescription>
          </DialogHeader>

          {formError && (
            <Alert variant="destructive">
              <AlertDescription>{formError}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleCreate} className="flex flex-col gap-4">
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="proj-name">Project Name</FieldLabel>
                <Input
                  id="proj-name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Nama Proyek..."
                  required
                />
              </Field>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="proj-section">Section</FieldLabel>
                  <Input
                    id="proj-section"
                    value={section}
                    onChange={(e) => setSection(e.target.value)}
                    placeholder="Unit / Bagian..."
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="proj-plant">Plant</FieldLabel>
                  <Input
                    id="proj-plant"
                    value={plant}
                    onChange={(e) => setPlant(e.target.value)}
                    placeholder="Lokasi Plant / Pabrik..."
                  />
                </Field>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="proj-plan">Plan Date</FieldLabel>
                  <Input
                    id="proj-plan"
                    type="date"
                    value={planDate}
                    onChange={(e) => setPlanDate(e.target.value)}
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="proj-budget">Budget (Rp)</FieldLabel>
                  <Input
                    id="proj-budget"
                    value={budget}
                    onChange={(e) => setBudget(e.target.value)}
                    placeholder="Contoh: 150000000"
                  />
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="proj-info">Information</FieldLabel>
                <Input
                  id="proj-info"
                  value={information}
                  onChange={(e) => setInformation(e.target.value)}
                  placeholder="Keterangan tambahan..."
                />
              </Field>

              <div className="flex gap-6 pt-2">
                <label className="flex items-center gap-2 text-sm font-medium">
                  <input
                    type="checkbox"
                    checked={spb}
                    onChange={(e) => setSpb(e.target.checked)}
                    className="size-4 rounded border-input"
                  />
                  Progress SPB
                </label>
                <label className="flex items-center gap-2 text-sm font-medium">
                  <input
                    type="checkbox"
                    checked={rab}
                    onChange={(e) => setRab(e.target.checked)}
                    className="size-4 rounded border-input"
                  />
                  Progress RAB
                </label>
                <label className="flex items-center gap-2 text-sm font-medium">
                  <input
                    type="checkbox"
                    checked={xxx}
                    onChange={(e) => setXxx(e.target.checked)}
                    className="size-4 rounded border-input"
                  />
                  Progress XXX
                </label>
              </div>
            </FieldGroup>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setIsAddOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={createMut.isPending}>
                {createMut.isPending ? 'Saving...' : 'Save Project'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Project Detail Modal */}
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
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl">
            <span>{project.name}</span>
            <Badge variant="outline" className="font-normal">
              {project.section || 'General'} · {project.plant || 'Site'}
            </Badge>
          </DialogTitle>
          <DialogDescription>Ringkasan estimasi biaya dan navigasi modul terkait.</DialogDescription>
        </DialogHeader>

        {/* Cost Estimate Summary Card */}
        <div className="flex flex-col gap-4 py-2">
          <div className="grid grid-cols-3 gap-3">
            <Card>
              <CardContent className="p-4">
                <div className="text-xs text-muted-foreground">BUDGET</div>
                <div className="mt-1 text-lg font-bold">{formatRupiah(cost?.budget ?? project.budget)}</div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <div className="text-xs text-muted-foreground">RAB (RENCANA)</div>
                <div className="mt-1 text-lg font-bold">{formatRupiah(cost?.rab ?? '0')}</div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <div className="text-xs text-muted-foreground">SPB (TERBAYAR)</div>
                <div className="mt-1 text-lg font-bold">{formatRupiah(cost?.spb ?? '0')}</div>
              </CardContent>
            </Card>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <Card>
              <CardContent className="p-4">
                <div className="text-xs text-muted-foreground">VARIANCE (Budget - RAB)</div>
                <div className="mt-1 text-base font-semibold">{formatRupiah(cost?.variance ?? '0')}</div>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="p-4">
                <div className="text-xs text-muted-foreground">SISA ANGGARAN (Budget - SPB)</div>
                <div className="mt-1 text-base font-semibold">{formatRupiah(cost?.sisa ?? '0')}</div>
              </CardContent>
            </Card>
          </div>
        </div>

        <DialogFooter className="flex-row items-center justify-between sm:justify-between">
          <div className="flex gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={() => onNavigate(`/schedule?projectId=${project.id}`)}
            >
              Schedule
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
              variant="default"
              onClick={() => onNavigate(`/rab?projectId=${project.id}`)}
            >
              RAB
            </Button>
          </div>
          <Button variant="ghost" size="sm" onClick={onClose}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
