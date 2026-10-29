import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listProjectsOptions,
  listSchedulesOptions,
  createScheduleMutation,
  deleteScheduleMutation,
} from '@/api/@tanstack/react-query.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FieldGroup, Field, FieldLabel } from '@/components/ui/field'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/lib/auth'

export const Route = createFileRoute('/schedule/')({
  component: SchedulePage,
  validateSearch: (search: Record<string, unknown>): { projectId?: string } => ({
    projectId: (search.projectId as string) || undefined,
  }),
})

function SchedulePage() {
  const queryClient = useQueryClient()
  const searchParams = Route.useSearch()
  const queryProjectId = searchParams.projectId || ''
  const { isAdmin, user, openLoginModal } = useAuth()

  const [selectedProjectId, setSelectedProjectId] = useState<string>(queryProjectId)
  const [search, setSearch] = useState('')
  const [isAddOpen, setIsAddOpen] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')

  // Form state
  const [code, setCode] = useState('')
  const [workItem, setWorkItem] = useState('')
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')
  const [weight, setWeight] = useState('0')
  const [progress, setProgress] = useState('0')
  const [status, setStatus] = useState<'planning' | 'on_going' | 'done' | 'late'>('planning')

  const { data: projectsData } = useQuery({
    ...listProjectsOptions({ query: { per_page: 100 } }),
  })
  const projects = projectsData?.data ?? []
  const activeProjectId = selectedProjectId || (projects.length > 0 ? projects[0].id : '')

  const { data: schedules = [], isLoading } = useQuery({
    ...listSchedulesOptions({ path: { id: activeProjectId } }),
    enabled: !!activeProjectId,
  })

  const createMut = useMutation(createScheduleMutation())
  const deleteMut = useMutation(deleteScheduleMutation())

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal()
      return
    }
    if (!activeProjectId) {
      alert('Pilih proyek terlebih dahulu.')
      return
    }
    setCode('')
    setWorkItem('')
    setStartDate(new Date().toISOString().slice(0, 10))
    setEndDate(new Date().toISOString().slice(0, 10))
    setWeight('0')
    setProgress('0')
    setStatus('planning')
    setErrorMessage('')
    setIsAddOpen(true)
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage('')
    if (!workItem.trim() || !startDate || !endDate) {
      setErrorMessage('Work item, start date, dan end date wajib diisi.')
      return
    }

    try {
      await createMut.mutateAsync({
        path: { id: activeProjectId },
        body: {
          code: code.trim(),
          work_item: workItem.trim(),
          start_date: startDate,
          end_date: endDate,
          weight,
          progress,
          status,
        },
      })
      await queryClient.invalidateQueries()
      setIsAddOpen(false)
    } catch (err: any) {
      setErrorMessage(err?.message || 'Gagal menyimpan schedule.')
    }
  }

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(`Hapus jadwal "${name}"?`)) return
    try {
      await deleteMut.mutateAsync({ path: { id } })
      await queryClient.invalidateQueries()
    } catch (err: any) {
      alert(err?.message || 'Gagal menghapus schedule.')
    }
  }

  const filtered = schedules.filter(
    (s) =>
      s.work_item.toLowerCase().includes(search.toLowerCase()) ||
      (s.code ? s.code.toLowerCase().includes(search.toLowerCase()) : false)
  )

  const getStatusBadgeVariant = (st: string) => {
    switch (st) {
      case 'done':
        return 'default'
      case 'on_going':
        return 'secondary'
      case 'late':
        return 'destructive'
      default:
        return 'outline'
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <PageHeader title="Schedule" onSearch={setSearch} onAdd={handleOpenAdd} addLabel="Schedule" />

      {/* Project Selector Bar */}
      <div className="flex items-center gap-4 rounded-lg border bg-card p-3">
        <span className="text-sm font-medium text-muted-foreground whitespace-nowrap">Proyek:</span>
        <div className="w-80">
          <Select
            value={activeProjectId}
            onValueChange={(val) => {
              if (val) setSelectedProjectId(val)
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder="Pilih Proyek" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                {projects.map((p) => (
                  <SelectItem key={p.id} value={p.id}>
                    {p.name} {p.section ? `(${p.section})` : ''}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="rounded-md border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12 text-center">No</TableHead>
              <TableHead className="w-28">Code</TableHead>
              <TableHead>Work Item</TableHead>
              <TableHead className="w-28 text-center">Start Date</TableHead>
              <TableHead className="w-28 text-center">End Date</TableHead>
              <TableHead className="w-24 text-center">Duration</TableHead>
              <TableHead className="w-24 text-right">Weight (%)</TableHead>
              <TableHead className="w-24 text-right">Progress (%)</TableHead>
              <TableHead className="w-28 text-center">Status</TableHead>
              {isAdmin && <TableHead className="w-20 text-center">Action</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 10 : 9} className="h-32 text-center text-muted-foreground">
                  Loading data...
                </TableCell>
              </TableRow>
            ) : filtered.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 10 : 9} className="h-48 text-center">
                  <Empty>
                    <EmptyHeader>
                      <EmptyTitle>Belum ada jadwal pekerjaan</EmptyTitle>
                      <EmptyDescription>
                        {search ? `Tidak ada item cocok dengan "${search}".` : 'Tambahkan tahapan jadwal baru untuk proyek ini.'}
                      </EmptyDescription>
                    </EmptyHeader>
                  </Empty>
                </TableCell>
              </TableRow>
            ) : (
              filtered.map((s, idx) => (
                <TableRow key={s.id}>
                  <TableCell className="text-center text-muted-foreground">{idx + 1}</TableCell>
                  <TableCell className="font-mono text-xs font-semibold">{s.code || '-'}</TableCell>
                  <TableCell className="font-medium text-foreground">{s.work_item}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{s.start_date}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{s.end_date}</TableCell>
                  <TableCell className="text-center">{s.duration} hari</TableCell>
                  <TableCell className="text-right font-medium">{s.weight}%</TableCell>
                  <TableCell className="text-right font-medium">{s.progress}%</TableCell>
                  <TableCell className="text-center">
                    <Badge variant={getStatusBadgeVariant(s.status)} className="capitalize text-xs font-normal">
                      {s.status.replace('_', ' ')}
                    </Badge>
                  </TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="text-xs text-destructive hover:bg-destructive/10"
                        onClick={() => handleDelete(s.id, s.work_item)}
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

      {/* Modal Add Schedule */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Add Schedule Item</DialogTitle>
            <DialogDescription>Tambahkan jadwal tahapan pekerjaan proyek.</DialogDescription>
          </DialogHeader>

          {errorMessage && (
            <Alert variant="destructive">
              <AlertDescription>{errorMessage}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleCreate} className="flex flex-col gap-4">
            <FieldGroup>
              <div className="grid grid-cols-3 gap-3">
                <Field>
                  <FieldLabel htmlFor="sch-code">Code</FieldLabel>
                  <Input
                    id="sch-code"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    placeholder="W-01"
                  />
                </Field>
                <div className="col-span-2">
                  <Field>
                    <FieldLabel htmlFor="sch-workitem">Work Item</FieldLabel>
                    <Input
                      id="sch-workitem"
                      value={workItem}
                      onChange={(e) => setWorkItem(e.target.value)}
                      placeholder="Nama Tahapan..."
                      required
                    />
                  </Field>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="sch-start">Start Date</FieldLabel>
                  <Input
                    id="sch-start"
                    type="date"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                    required
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="sch-end">End Date</FieldLabel>
                  <Input
                    id="sch-end"
                    type="date"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                    required
                  />
                </Field>
              </div>

              <div className="grid grid-cols-3 gap-3">
                <Field>
                  <FieldLabel htmlFor="sch-weight">Weight (%)</FieldLabel>
                  <Input
                    id="sch-weight"
                    type="number"
                    step="0.01"
                    value={weight}
                    onChange={(e) => setWeight(e.target.value)}
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="sch-prog">Progress (%)</FieldLabel>
                  <Input
                    id="sch-prog"
                    type="number"
                    step="0.01"
                    value={progress}
                    onChange={(e) => setProgress(e.target.value)}
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="sch-status">Status</FieldLabel>
                  <Select
                    value={status}
                    onValueChange={(val) => setStatus(val as any)}
                  >
                    <SelectTrigger id="sch-status">
                      <SelectValue placeholder="Status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="planning">Planning</SelectItem>
                        <SelectItem value="on_going">On Going</SelectItem>
                        <SelectItem value="done">Done</SelectItem>
                        <SelectItem value="late">Late</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </Field>
              </div>
            </FieldGroup>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setIsAddOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={createMut.isPending}>
                {createMut.isPending ? 'Saving...' : 'Save'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
