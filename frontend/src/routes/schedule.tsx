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
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/lib/auth'

export const Route = createFileRoute('/schedule')({
  component: SchedulePage,
  validateSearch: (search: Record<string, unknown>) => ({
    projectId: (search.projectId as string) || '',
  }),
})

function SchedulePage() {
  const queryClient = useQueryClient()
  const { projectId: queryProjectId } = Route.useSearch()
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
      setErrorMessage(err?.message || 'Gagal menyimpan item jadwal.')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus item jadwal ini?')) return
    try {
      await deleteMut.mutateAsync({ path: { id } })
      await queryClient.invalidateQueries()
    } catch (err: any) {
      alert(err?.message || 'Gagal menghapus item jadwal.')
    }
  }

  const filtered = schedules.filter((s) => {
    if (!search) return true
    const q = search.toLowerCase()
    return s.work_item.toLowerCase().includes(q) || (s.code ? s.code.toLowerCase().includes(q) : false)
  })

  return (
    <div>
      <PageHeader title="SCHEDULE" onSearch={setSearch} onAdd={handleOpenAdd} />

      {/* Project Selector */}
      <div className="mb-4 flex items-center gap-3 rounded-xl bg-white p-3 border shadow-xs">
        <span className="text-xs font-bold text-foreground">PROJECT:</span>
        <div className="w-80">
          <Select
            value={activeProjectId}
            onValueChange={(val) => {
              if (val) setSelectedProjectId(val)
            }}
          >
            <SelectTrigger className="h-9 text-xs font-semibold">
              <SelectValue placeholder="Pilih Proyek..." />
            </SelectTrigger>
            <SelectContent>
              {projects.map((p) => (
                <SelectItem key={p.id} value={p.id}>
                  {p.name} {p.section ? `(${p.section})` : ''}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="overflow-x-auto border bg-white shadow-sm">
        <Table className="text-xs">
          <TableHeader className="bg-white">
            <TableRow className="border-b hover:bg-transparent">
              <TableHead className="w-12 border-r text-center font-bold">NO</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">CODE</TableHead>
              <TableHead className="min-w-64 border-r text-center font-bold">WORK ITEM</TableHead>
              <TableHead className="w-28 border-r text-center font-bold">START DATE</TableHead>
              <TableHead className="w-28 border-r text-center font-bold">END DATE</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">DURATION</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">WEIGHT (%)</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">PROGRESS (%)</TableHead>
              <TableHead className="w-28 border-r text-center font-bold">STATUS</TableHead>
              {isAdmin && <TableHead className="w-16 text-center font-bold">AKSI</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 10 : 9} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  Memuat data...
                </TableCell>
              </TableRow>
            ) : filtered.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 10 : 9} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  {search ? `Tidak ada jadwal yang cocok dengan "${search}".` : 'Belum ada jadwal untuk proyek ini. Klik ADD untuk membuat.'}
                </TableCell>
              </TableRow>
            ) : (
              filtered.map((item, idx) => (
                <TableRow key={item.id} className="border-b bg-muted/40 hover:bg-muted/80">
                  <TableCell className="border-r text-center font-medium">{idx + 1}</TableCell>
                  <TableCell className="border-r font-mono font-medium">{item.code || '-'}</TableCell>
                  <TableCell className="border-r font-medium text-foreground">{item.work_item}</TableCell>
                  <TableCell className="border-r text-center">{item.start_date}</TableCell>
                  <TableCell className="border-r text-center">{item.end_date}</TableCell>
                  <TableCell className="border-r text-center font-medium">{item.duration} hari</TableCell>
                  <TableCell className="border-r text-right">{item.weight}%</TableCell>
                  <TableCell className="border-r text-right font-semibold">{item.progress}%</TableCell>
                  <TableCell className="border-r text-center">
                    <Badge variant={item.status === 'done' ? 'default' : 'secondary'} className="text-[10px]">
                      {item.status}
                    </Badge>
                  </TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-7 text-xs text-red-600 hover:bg-red-50 hover:text-red-700"
                        onClick={() => handleDelete(item.id)}
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

      {/* Modal Add Schedule */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-md rounded-2xl">
          <DialogHeader>
            <DialogTitle className="text-lg font-bold tracking-wide">TAMBAH ITEM JADWAL</DialogTitle>
          </DialogHeader>

          {errorMessage && (
            <div className="rounded-md border border-red-200 bg-red-50 p-3 text-xs text-red-600">
              {errorMessage}
            </div>
          )}

          <form onSubmit={handleCreate} className="space-y-4 text-xs font-semibold">
            <div>
              <Label className="text-xs">CODE</Label>
              <Input
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="misal WBS-01"
                className="mt-1 font-mono"
              />
            </div>

            <div>
              <Label className="text-xs">WORK ITEM</Label>
              <Input
                value={workItem}
                onChange={(e) => setWorkItem(e.target.value)}
                placeholder="Uraian pekerjaan jadwal..."
                className="mt-1"
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label className="text-xs">START DATE</Label>
                <Input
                  type="date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                  className="mt-1"
                  required
                />
              </div>
              <div>
                <Label className="text-xs">END DATE</Label>
                <Input
                  type="date"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                  className="mt-1"
                  required
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label className="text-xs">WEIGHT (BOBOT %)</Label>
                <Input
                  value={weight}
                  onChange={(e) => setWeight(e.target.value)}
                  placeholder="0.00"
                  className="mt-1"
                />
              </div>
              <div>
                <Label className="text-xs">PROGRESS (%)</Label>
                <Input
                  value={progress}
                  onChange={(e) => setProgress(e.target.value)}
                  placeholder="0.00"
                  className="mt-1"
                />
              </div>
            </div>

            <div>
              <Label className="text-xs">STATUS</Label>
              <Select value={status} onValueChange={(val: any) => setStatus(val)}>
                <SelectTrigger className="mt-1">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="planning">Planning (Rencana)</SelectItem>
                  <SelectItem value="on_going">On Going (Berjalan)</SelectItem>
                  <SelectItem value="done">Done (Selesai)</SelectItem>
                  <SelectItem value="late">Late (Terlambat)</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="flex justify-end gap-2 pt-2">
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
    </div>
  )
}
