import { useState, useMemo } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listProjectsOptions,
  getProjectRabOptions,
  createRabItemMutation,
  deleteRabItemMutation,
  getProjectBoqOptions,
  listAhspsOptions,
} from '@/api/@tanstack/react-query.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { formatRupiah, formatDecimal } from '@/lib/utils'
import { useAuth } from '@/lib/auth'

export const Route = createFileRoute('/rab')({
  component: RabPage,
  validateSearch: (search: Record<string, unknown>): { projectId?: string } => ({
    projectId: (search.projectId as string) || undefined,
  }),
})

function RabPage() {
  const queryClient = useQueryClient()
  const searchParams = Route.useSearch()
  const queryProjectId = searchParams.projectId || ''
  const { isAdmin, user, openLoginModal } = useAuth()

  const [selectedProjectId, setSelectedProjectId] = useState<string>(queryProjectId)
  const [search, setSearch] = useState('')
  const [isAddOpen, setIsAddOpen] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')

  // Form state
  const [selectedBoqId, setSelectedBoqId] = useState('')
  const [workItem, setWorkItem] = useState('')
  const [unit, setUnit] = useState('')
  const [volume, setVolume] = useState('0')
  const [unitPrice, setUnitPrice] = useState('0')

  // Fetch projects
  const { data: projectsData } = useQuery({
    ...listProjectsOptions({ query: { per_page: 100 } }),
  })
  const projects = projectsData?.data ?? []
  const activeProjectId = selectedProjectId || (projects.length > 0 ? projects[0].id : '')

  // Fetch RAB items & total
  const { data: rabData, isLoading } = useQuery({
    ...getProjectRabOptions({
      path: { id: activeProjectId },
      query: { q: search || undefined },
    }),
    enabled: !!activeProjectId,
  })

  // Fetch BOQ items for project to populate Add options
  const { data: boqItems = [] } = useQuery({
    ...getProjectBoqOptions({ path: { id: activeProjectId } }),
    enabled: isAddOpen && !!activeProjectId,
  })

  // Fetch AHSPs for price lookups
  const { data: ahsps = [] } = useQuery({
    ...listAhspsOptions({}),
    enabled: isAddOpen,
  })

  const createMut = useMutation(createRabItemMutation())
  const deleteMut = useMutation(deleteRabItemMutation())

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal()
      return
    }
    if (!activeProjectId) {
      alert('Pilih proyek terlebih dahulu.')
      return
    }
    setSelectedBoqId('')
    setWorkItem('')
    setUnit('')
    setVolume('1')
    setUnitPrice('0')
    setErrorMessage('')
    setIsAddOpen(true)
  }

  const handleSelectBoqItem = (boqId: string) => {
    setSelectedBoqId(boqId)
    const boq = boqItems.find((b) => b.id === boqId)
    if (boq) {
      setWorkItem(boq.work_item)
      setUnit(boq.unit)
      setVolume(boq.planned_volume || '1')

      // Look up AHSP price if linked
      if (boq.ahsp_id) {
        const foundAhsp = ahsps.find((a) => a.id === boq.ahsp_id)
        if (foundAhsp) {
          setUnitPrice(foundAhsp.unit_price_of_work)
        }
      }
    }
  }

  const handleSelectAhspDirect = (ahspId: string) => {
    const foundAhsp = ahsps.find((a) => a.id === ahspId)
    if (foundAhsp) {
      setWorkItem(foundAhsp.work_item)
      setUnit(foundAhsp.unit)
      setUnitPrice(foundAhsp.unit_price_of_work)
    }
  }

  const previewTotal = useMemo(() => {
    const v = Number.parseFloat(volume) || 0
    const p = Number.parseInt(unitPrice.replace(/\D/g, ''), 10) || 0
    return Math.round(v * p)
  }, [volume, unitPrice])

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
          boq_item_id: selectedBoqId || undefined,
          work_item: workItem.trim(),
          unit: unit.trim(),
          volume: volume || '0',
          unit_price: unitPrice.replace(/\D/g, '') || '0',
        },
      })
      await queryClient.invalidateQueries()
      setIsAddOpen(false)
    } catch (err: any) {
      setErrorMessage(err?.message || 'Gagal menambahkan baris RAB.')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus baris RAB ini?')) return
    try {
      await deleteMut.mutateAsync({ path: { id } })
      await queryClient.invalidateQueries()
    } catch (err: any) {
      alert(err?.message || 'Gagal menghapus baris RAB.')
    }
  }

  const items = rabData?.items ?? []
  const hasZeroPrice = items.some((it) => it.unit_price === '0')

  return (
    <div>
      <PageHeader title="RAB" onSearch={setSearch} onAdd={handleOpenAdd} />

      {/* Project Selector (Mockup Hal. 10) */}
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

      {/* Warning banner if any item lacks price */}
      {hasZeroPrice && (
        <div className="mb-4 rounded-xl border border-amber-300 bg-amber-50 p-3 text-xs text-amber-800">
          ⚠️ Ada item pekerjaan yang belum memiliki harga satuan (Unit Price = Rp0).
        </div>
      )}

      <div className="overflow-x-auto border bg-white shadow-sm">
        <Table className="text-xs">
          <TableHeader className="bg-white">
            <TableRow className="border-b hover:bg-transparent">
              <TableHead className="w-12 border-r text-center font-bold">NO</TableHead>
              <TableHead className="min-w-64 border-r text-center font-bold">WORK ITEM</TableHead>
              <TableHead className="w-28 border-r text-center font-bold">VOLUME</TableHead>
              <TableHead className="w-20 border-r text-center font-bold">UNIT</TableHead>
              <TableHead className="w-40 border-r text-center font-bold">UNIT PRICE</TableHead>
              <TableHead className="w-44 border-r text-center font-bold">TOTAL COST</TableHead>
              {isAdmin && <TableHead className="w-16 text-center font-bold">AKSI</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 7 : 6} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  Memuat data...
                </TableCell>
              </TableRow>
            ) : items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 7 : 6} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  {search ? `Tidak ada item yang cocok dengan "${search}".` : 'Belum ada data RAB untuk proyek ini. Klik ADD untuk membuat.'}
                </TableCell>
              </TableRow>
            ) : (
              items.map((item, idx) => (
                <TableRow key={item.id} className="border-b bg-muted/40 hover:bg-muted/80">
                  <TableCell className="border-r text-center font-medium">{idx + 1}</TableCell>
                  <TableCell className="border-r font-medium text-foreground">{item.work_item}</TableCell>
                  <TableCell className="border-r text-right">{formatDecimal(item.volume)}</TableCell>
                  <TableCell className="border-r text-center">{item.unit}</TableCell>
                  <TableCell className="border-r text-right font-medium">
                    {item.unit_price === '0' ? (
                      <span className="text-amber-600 font-normal">⚠ belum ada harga</span>
                    ) : (
                      formatRupiah(item.unit_price)
                    )}
                  </TableCell>
                  <TableCell className="border-r text-right font-bold text-foreground">
                    {formatRupiah(item.total_cost)}
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

        {/* Footer: TOTAL RAB */}
        {items.length > 0 && (
          <div className="flex items-center justify-between border-t bg-primary px-6 py-3 text-sm font-extrabold text-primary-foreground">
            <span>TOTAL RAB</span>
            <span>{formatRupiah(rabData?.total_rab ?? '0')}</span>
          </div>
        )}
      </div>

      {/* Modal Add RAB (Mockup Hal. 10) */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-md rounded-2xl">
          <DialogHeader>
            <DialogTitle className="text-lg font-bold tracking-wide">TAMBAH ITEM RAB</DialogTitle>
          </DialogHeader>

          {errorMessage && (
            <div className="rounded-md border border-red-200 bg-red-50 p-3 text-xs text-red-600">
              {errorMessage}
            </div>
          )}

          <form onSubmit={handleCreate} className="space-y-4 text-xs font-semibold">
            {boqItems.length > 0 && (
              <div>
                <Label className="text-xs">PILIH DARI ITEM BOQ / VOLUME</Label>
                <select
                  className="mt-1 h-9 w-full rounded-md border border-gray-300 bg-white px-3 py-1 text-xs"
                  value={selectedBoqId}
                  onChange={(e) => handleSelectBoqItem(e.target.value)}
                >
                  <option value="">-- Pilih Item BOQ Proyek --</option>
                  {boqItems.map((b) => (
                    <option key={b.id} value={b.id}>
                      {b.work_item} ({b.planned_volume} {b.unit})
                    </option>
                  ))}
                </select>
              </div>
            )}

            {ahsps.length > 0 && (
              <div>
                <Label className="text-xs">ATAU PILIH DARI MASTER AHSP (HARGA SATUAN)</Label>
                <select
                  className="mt-1 h-9 w-full rounded-md border border-gray-300 bg-white px-3 py-1 text-xs"
                  onChange={(e) => handleSelectAhspDirect(e.target.value)}
                  defaultValue=""
                >
                  <option value="" disabled>-- Pilih Master AHSP --</option>
                  {ahsps.map((a) => (
                    <option key={a.id} value={a.id}>
                      {a.code} - {a.work_item} ({formatRupiah(a.unit_price_of_work)}/{a.unit})
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

            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label className="text-xs">VOLUME</Label>
                <Input
                  value={volume}
                  onChange={(e) => setVolume(e.target.value)}
                  placeholder="0.00"
                  className="mt-1"
                  required
                />
              </div>
              <div>
                <Label className="text-xs">UNIT</Label>
                <Input
                  value={unit}
                  onChange={(e) => setUnit(e.target.value)}
                  placeholder="m3, ls, unit..."
                  className="mt-1"
                  required
                />
              </div>
            </div>

            <div>
              <Label className="text-xs">UNIT PRICE (HARGA SATUAN)</Label>
              <Input
                value={unitPrice}
                onChange={(e) => setUnitPrice(e.target.value)}
                placeholder="Nominal Rupiah..."
                className="mt-1"
                required
              />
            </div>

            <div className="rounded-xl border bg-muted/40 p-3">
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground">Perkiraan Biaya (Total Cost):</span>
                <span className="text-base font-extrabold text-primary">{formatRupiah(previewTotal)}</span>
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
