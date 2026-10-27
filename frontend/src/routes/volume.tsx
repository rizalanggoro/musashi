import { useState, useMemo } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listProjectsOptions,
  listVolumesOptions,
  createVolumeMutation,
  deleteVolumeMutation,
  listAhspsOptions,
} from '@/api/@tanstack/react-query.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { formatDecimal } from '@/lib/utils'
import { useAuth } from '@/lib/auth'

export const Route = createFileRoute('/volume')({
  component: VolumePage,
  validateSearch: (search: Record<string, unknown>): { projectId?: string } => ({
    projectId: (search.projectId as string) || undefined,
  }),
})

function VolumePage() {
  const queryClient = useQueryClient()
  const searchParams = Route.useSearch()
  const queryProjectId = searchParams.projectId || ''
  const { isAdmin, user, openLoginModal } = useAuth()

  const [selectedProjectId, setSelectedProjectId] = useState<string>(queryProjectId)
  const [search, setSearch] = useState('')
  const [isAddOpen, setIsAddOpen] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')

  // Form state
  const [selectedAhspId, setSelectedAhspId] = useState('')
  const [workItem, setWorkItem] = useState('')
  const [unit, setUnit] = useState('')
  const [plannedVolume, setPlannedVolume] = useState('0')
  const [date, setDate] = useState('')
  const [length, setLength] = useState('')
  const [width, setWidth] = useState('')
  const [height, setHeight] = useState('')
  const [unitWeight, setUnitWeight] = useState('1')
  const [description, setDescription] = useState('')

  // Fetch projects for selector
  const { data: projectsData } = useQuery({
    ...listProjectsOptions({ query: { per_page: 100 } }),
  })
  const projects = projectsData?.data ?? []

  // Default project selection if not set
  const activeProjectId = selectedProjectId || (projects.length > 0 ? projects[0].id : '')

  // Fetch volume rows
  const { data: volumeRows = [], isLoading } = useQuery({
    ...listVolumesOptions({
      path: { id: activeProjectId },
      query: { q: search || undefined },
    }),
    enabled: !!activeProjectId,
  })

  // Fetch AHSPs for master auto-complete
  const { data: ahsps = [] } = useQuery({
    ...listAhspsOptions({}),
    enabled: isAddOpen,
  })

  const createMut = useMutation(createVolumeMutation())
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
    setSelectedAhspId('')
    setWorkItem('')
    setUnit('')
    setPlannedVolume('0')
    setDate(new Date().toISOString().slice(0, 10))
    setLength('')
    setWidth('')
    setHeight('')
    setUnitWeight('1')
    setDescription('')
    setErrorMessage('')
    setIsAddOpen(true)
  }

  const handleSelectAhsp = (ahspId: string) => {
    setSelectedAhspId(ahspId)
    const found = ahsps.find((a) => a.id === ahspId)
    if (found) {
      setWorkItem(found.work_item)
      setUnit(found.unit)
    }
  }

  // Live preview of quantity & total quantity
  const preview = useMemo(() => {
    const l = Number.parseFloat(length) || 0
    const w = Number.parseFloat(width) || 0
    const h = Number.parseFloat(height) || 0
    const uw = Number.parseFloat(unitWeight) || 1
    const pVol = Number.parseFloat(plannedVolume) || 0

    let qty = 0
    if (l > 0 && w > 0 && h > 0) {
      qty = l * w * h
    } else {
      qty = pVol
    }
    const totalQty = qty * uw
    return { qty, totalQty }
  }, [length, width, height, unitWeight, plannedVolume])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage('')
    if (!workItem.trim() || !unit.trim()) {
      setErrorMessage('Work item dan unit wajib diisi.')
      return
    }

    try {
      await createMut.mutateAsync({
        path: { id: activeProjectId },
        body: {
          ahsp_id: selectedAhspId || undefined,
          work_item: workItem.trim(),
          unit: unit.trim(),
          planned_volume: plannedVolume || '0',
          date: date || undefined,
          length: length || undefined,
          width: width || undefined,
          height: height || undefined,
          unit_weight: unitWeight || '1',
          description: description.trim(),
        },
      })
      await queryClient.invalidateQueries()
      setIsAddOpen(false)
    } catch (err: any) {
      setErrorMessage(err?.message || 'Gagal menyimpan take-off volume.')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus baris volume take-off ini?')) return
    try {
      await deleteMut.mutateAsync({ path: { id } })
      await queryClient.invalidateQueries()
    } catch (err: any) {
      alert(err?.message || 'Gagal menghapus baris volume.')
    }
  }

  return (
    <div>
      <PageHeader title="VOLUME" onSearch={setSearch} onAdd={handleOpenAdd} />

      {/* Project Selector (Mockup Hal. 9) */}
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
              <TableHead className="w-10 border-r text-center font-bold">NO</TableHead>
              <TableHead className="min-w-48 border-r text-center font-bold">WORK ITEM</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">VOLUME</TableHead>
              <TableHead className="w-16 border-r text-center font-bold">UNIT</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">LENGTH (m)</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">WIDTH (m)</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">HEIGHT (m)</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">QUANTITY</TableHead>
              <TableHead className="min-w-36 border-r text-center font-bold">UNIT WEIGHT (kg/m'; kg/m³)</TableHead>
              <TableHead className="w-28 border-r text-center font-bold">TOTAL QUANTITY</TableHead>
              {isAdmin && <TableHead className="w-16 text-center font-bold">AKSI</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 11 : 10} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  Memuat data...
                </TableCell>
              </TableRow>
            ) : volumeRows.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 11 : 10} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  {search ? `Tidak ada volume yang cocok dengan "${search}".` : 'Belum ada data take-off volume untuk proyek ini. Klik ADD untuk menambah.'}
                </TableCell>
              </TableRow>
            ) : (
              volumeRows.map((row, idx) => (
                <TableRow key={row.id} className="border-b bg-muted/40 hover:bg-muted/80">
                  <TableCell className="border-r text-center font-medium">{idx + 1}</TableCell>
                  <TableCell className="border-r font-medium text-foreground">
                    {row.work_item}
                    {row.description && <span className="block text-[10px] text-muted-foreground">{row.description}</span>}
                  </TableCell>
                  <TableCell className="border-r text-right">{formatDecimal(row.planned_volume)}</TableCell>
                  <TableCell className="border-r text-center">{row.unit}</TableCell>
                  <TableCell className="border-r text-right">{row.length ? formatDecimal(row.length) : '-'}</TableCell>
                  <TableCell className="border-r text-right">{row.width ? formatDecimal(row.width) : '-'}</TableCell>
                  <TableCell className="border-r text-right">{row.height ? formatDecimal(row.height) : '-'}</TableCell>
                  <TableCell className="border-r text-right font-medium">{formatDecimal(row.quantity)}</TableCell>
                  <TableCell className="border-r text-right">{formatDecimal(row.unit_weight)}</TableCell>
                  <TableCell className="border-r text-right font-bold text-foreground">{formatDecimal(row.total_quantity)}</TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-7 text-xs text-red-600 hover:bg-red-50 hover:text-red-700"
                        onClick={() => handleDelete(row.id)}
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

      {/* Modal Add Volume (Mockup Hal. 9) */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-xl rounded-2xl">
          <DialogHeader>
            <DialogTitle className="text-lg font-bold tracking-wide">TAMBAH VOLUME TAKE-OFF</DialogTitle>
          </DialogHeader>

          {errorMessage && (
            <div className="rounded-md border border-red-200 bg-red-50 p-3 text-xs text-red-600">
              {errorMessage}
            </div>
          )}

          <form onSubmit={handleCreate} className="space-y-4 text-xs font-semibold">
            {ahsps.length > 0 && (
              <div>
                <Label className="text-xs">PILIH DARI MASTER AHSP (OPSIONAL)</Label>
                <select
                  className="mt-1 h-9 w-full rounded-md border border-gray-300 bg-white px-3 py-1 text-xs"
                  value={selectedAhspId}
                  onChange={(e) => handleSelectAhsp(e.target.value)}
                >
                  <option value="">-- Input Manual atau Pilih AHSP --</option>
                  {ahsps.map((a) => (
                    <option key={a.id} value={a.id}>
                      {a.code} - {a.work_item} ({a.unit})
                    </option>
                  ))}
                </select>
              </div>
            )}

            <div>
              <Label className="text-xs">WORK ITEM</Label>
              <Input
                value={workItem}
                onChange={(e) => setWorkItem(e.target.value)}
                placeholder="Uraian pekerjaan..."
                className="mt-1"
                required
              />
            </div>

            <div className="grid grid-cols-3 gap-3">
              <div>
                <Label className="text-xs">UNIT</Label>
                <Input
                  value={unit}
                  onChange={(e) => setUnit(e.target.value)}
                  placeholder="m3, kg, m'..."
                  className="mt-1"
                  required
                />
              </div>
              <div>
                <Label className="text-xs">VOLUME RENCANA</Label>
                <Input
                  value={plannedVolume}
                  onChange={(e) => setPlannedVolume(e.target.value)}
                  placeholder="0.00"
                  className="mt-1"
                />
              </div>
              <div>
                <Label className="text-xs">DATE</Label>
                <Input
                  type="date"
                  value={date}
                  onChange={(e) => setDate(e.target.value)}
                  className="mt-1"
                />
              </div>
            </div>

            {/* Dimensi Take-off */}
            <div className="rounded-xl border bg-muted/30 p-3 space-y-3">
              <span className="text-xs font-bold text-muted-foreground uppercase">Dimensi Take-off (Jika Ada)</span>
              <div className="grid grid-cols-3 gap-3">
                <div>
                  <Label className="text-xs">LENGTH (m)</Label>
                  <Input
                    value={length}
                    onChange={(e) => setLength(e.target.value)}
                    placeholder="Panjang"
                    className="mt-1 bg-white"
                  />
                </div>
                <div>
                  <Label className="text-xs">WIDTH (m)</Label>
                  <Input
                    value={width}
                    onChange={(e) => setWidth(e.target.value)}
                    placeholder="Lebar"
                    className="mt-1 bg-white"
                  />
                </div>
                <div>
                  <Label className="text-xs">HEIGHT (m)</Label>
                  <Input
                    value={height}
                    onChange={(e) => setHeight(e.target.value)}
                    placeholder="Tinggi"
                    className="mt-1 bg-white"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <Label className="text-xs">UNIT WEIGHT (default 1)</Label>
                  <Input
                    value={unitWeight}
                    onChange={(e) => setUnitWeight(e.target.value)}
                    placeholder="1"
                    className="mt-1 bg-white"
                  />
                </div>
                <div>
                  <Label className="text-xs">DESCRIPTION / ZONA</Label>
                  <Input
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    placeholder="misal Kolom Lantai 2"
                    className="mt-1 bg-white"
                  />
                </div>
              </div>

              <div className="flex items-center justify-between rounded-lg bg-white p-2 text-xs border">
                <div>
                  <span className="text-muted-foreground">Kalkulasi Quantity: </span>
                  <span className="font-bold">{formatDecimal(preview.qty)}</span>
                </div>
                <div>
                  <span className="text-muted-foreground">Total Quantity: </span>
                  <span className="font-extrabold text-primary">{formatDecimal(preview.totalQty)}</span>
                </div>
              </div>
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
