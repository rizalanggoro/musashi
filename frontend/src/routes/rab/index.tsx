import { useState } from 'react'
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
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow, TableFooter } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FieldGroup, Field, FieldLabel } from '@/components/ui/field'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Card, CardContent } from '@/components/ui/card'
import { formatRupiah, formatDecimal } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import { AlertCircle, Download } from 'lucide-react'

export const Route = createFileRoute('/rab/')({
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
      setVolume(boq.planned_volume)

      if (boq.ahsp_id) {
        const ahsp = ahsps.find((a) => a.id === boq.ahsp_id)
        if (ahsp) {
          setUnitPrice(ahsp.unit_price_of_work)
        }
      }
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage('')
    if (!workItem.trim() || !unit.trim()) {
      setErrorMessage('Nama item pekerjaan dan unit wajib diisi.')
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
    <div className="flex flex-col gap-4">
      <PageHeader title="RAB (Rencana Anggaran Biaya)" onSearch={setSearch} onAdd={handleOpenAdd} addLabel="RAB Item" />

      {/* Project Selector Bar */}
      <div className="flex flex-wrap items-center justify-between gap-4 rounded-lg border bg-card p-3">
        <div className="flex items-center gap-3">
          <span className="text-sm font-medium text-muted-foreground whitespace-nowrap">Proyek:</span>
          <div className="w-80">
            <Select
              value={activeProjectId}
              onValueChange={(val) => {
                if (val) setSelectedProjectId(val)
              }}
            >
              <SelectTrigger>
                <SelectValue placeholder="Pilih Proyek..." />
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

        {activeProjectId && (
          <div className="flex items-center gap-2">
            <Button
              size="sm"
              variant="outline"
              render={<a href={`/api/v1/projects/${activeProjectId}/exports/rab.xlsx`} download />}
              nativeButton={false}
            >
              <Download data-icon="inline-start" />
              Excel
            </Button>
            <Button
              size="sm"
              variant="outline"
              render={<a href={`/api/v1/projects/${activeProjectId}/exports/rab.pdf`} target="_blank" rel="noreferrer" />}
              nativeButton={false}
            >
              <Download data-icon="inline-start" />
              PDF
            </Button>
          </div>
        )}
      </div>

      {/* Warning banner using Alert */}
      {hasZeroPrice && (
        <Alert>
          <AlertCircle />
          <AlertTitle>Perhatian</AlertTitle>
          <AlertDescription>Terdapat item pekerjaan yang belum memiliki harga satuan (Unit Price = Rp0).</AlertDescription>
        </Alert>
      )}

      <div className="rounded-md border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12 text-center">No</TableHead>
              <TableHead>Work Item</TableHead>
              <TableHead className="w-28 text-right">Volume</TableHead>
              <TableHead className="w-20 text-center">Unit</TableHead>
              <TableHead className="w-40 text-right">Unit Price</TableHead>
              <TableHead className="w-44 text-right">Total Cost</TableHead>
              {isAdmin && <TableHead className="w-20 text-center">Action</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 7 : 6} className="h-32 text-center text-muted-foreground">
                  Loading data...
                </TableCell>
              </TableRow>
            ) : items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 7 : 6} className="h-48 text-center">
                  <Empty>
                    <EmptyHeader>
                      <EmptyTitle>Belum ada data RAB</EmptyTitle>
                      <EmptyDescription>
                        {search ? `Tidak ada item yang cocok dengan "${search}".` : 'Belum ada data RAB untuk proyek ini.'}
                      </EmptyDescription>
                    </EmptyHeader>
                  </Empty>
                </TableCell>
              </TableRow>
            ) : (
              items.map((item, idx) => (
                <TableRow key={item.id}>
                  <TableCell className="text-center text-muted-foreground">{idx + 1}</TableCell>
                  <TableCell className="font-medium text-foreground">{item.work_item}</TableCell>
                  <TableCell className="text-right">{formatDecimal(item.volume)}</TableCell>
                  <TableCell className="text-center">{item.unit}</TableCell>
                  <TableCell className="text-right font-medium">
                    {item.unit_price === '0' ? (
                      <span className="text-amber-600 text-xs font-normal">Belum ada harga</span>
                    ) : (
                      formatRupiah(item.unit_price)
                    )}
                  </TableCell>
                  <TableCell className="text-right font-semibold text-foreground">
                    {formatRupiah(item.total_cost)}
                  </TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="text-xs text-destructive hover:bg-destructive/10"
                        onClick={() => handleDelete(item.id)}
                      >
                        Delete
                      </Button>
                    </TableCell>
                  )}
                </TableRow>
              ))
            )}
          </TableBody>
          {items.length > 0 && (
            <TableFooter>
              <TableRow className="font-bold text-base">
                <TableCell colSpan={5} className="text-right">
                  TOTAL BIAYA KESELURUHAN (RAB):
                </TableCell>
                <TableCell className="text-right font-extrabold text-foreground">
                  {formatRupiah(rabData?.total_rab || '0')}
                </TableCell>
                {isAdmin && <TableCell></TableCell>}
              </TableRow>
            </TableFooter>
          )}
        </Table>
      </div>

      {/* Modal Add RAB Item */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Add RAB Item</DialogTitle>
            <DialogDescription>Tambahkan baris rencana anggaran biaya proyek.</DialogDescription>
          </DialogHeader>

          {errorMessage && (
            <Alert variant="destructive">
              <AlertDescription>{errorMessage}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleCreate} className="flex flex-col gap-4">
            <FieldGroup>
              {boqItems.length > 0 && (
                <Field>
                  <FieldLabel htmlFor="rab-boq">Pilih dari BOQ (Opsional)</FieldLabel>
                  <Select value={selectedBoqId} onValueChange={(val) => { if (val) handleSelectBoqItem(val) }}>
                    <SelectTrigger id="rab-boq">
                      <SelectValue placeholder="-- Pilih Item BOQ --" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        {boqItems.map((b) => (
                          <SelectItem key={b.id} value={b.id}>
                            {b.work_item} ({b.planned_volume} {b.unit})
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </Field>
              )}

              <Field>
                <FieldLabel htmlFor="rab-item">Work Item</FieldLabel>
                <Input
                  id="rab-item"
                  value={workItem}
                  onChange={(e) => setWorkItem(e.target.value)}
                  placeholder="Uraian pekerjaan..."
                  required
                />
              </Field>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="rab-unit">Unit</FieldLabel>
                  <Input
                    id="rab-unit"
                    value={unit}
                    onChange={(e) => setUnit(e.target.value)}
                    placeholder="m3, ls, unit..."
                    required
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="rab-vol">Volume</FieldLabel>
                  <Input
                    id="rab-vol"
                    type="number"
                    step="0.0001"
                    value={volume}
                    onChange={(e) => setVolume(e.target.value)}
                    required
                  />
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="rab-price">Unit Price (Rp)</FieldLabel>
                <Input
                  id="rab-price"
                  value={unitPrice}
                  onChange={(e) => setUnitPrice(e.target.value)}
                  placeholder="Harga satuan..."
                  required
                />
              </Field>

              {/* Live Subtotal Card */}
              <Card className="bg-muted/40">
                <CardContent className="p-3">
                  <div className="flex items-center justify-between text-xs text-muted-foreground">
                    <span>Estimasi Total:</span>
                    <strong className="text-sm font-bold text-foreground">
                      {formatRupiah(Math.round(parseFloat(volume || '0') * parseInt(unitPrice.replace(/\D/g, '') || '0', 10)))}
                    </strong>
                  </div>
                </CardContent>
              </Card>
            </FieldGroup>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setIsAddOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={createMut.isPending}>
                {createMut.isPending ? 'Saving...' : 'Save Item'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
