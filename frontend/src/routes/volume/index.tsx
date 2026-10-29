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
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FieldGroup, Field, FieldLabel } from '@/components/ui/field'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Card, CardContent } from '@/components/ui/card'
import { formatDecimal } from '@/lib/utils'
import { useAuth } from '@/lib/auth'

export const Route = createFileRoute('/volume/')({
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
    const ahsp = ahsps.find((a) => a.id === ahspId)
    if (ahsp) {
      setWorkItem(ahsp.work_item)
      setUnit(ahsp.unit)
    }
  }

  // Live calculation preview
  const preview = useMemo(() => {
    const l = parseFloat(length || '0')
    const w = parseFloat(width || '0')
    const h = parseFloat(height || '0')
    const uw = parseFloat(unitWeight || '1') || 1
    const planned = parseFloat(plannedVolume || '0')

    let qty = 0
    if (l > 0 && w > 0 && h > 0) {
      qty = l * w * h
    } else {
      qty = planned
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
      setErrorMessage(err?.message || 'Gagal menyimpan data volume.')
    }
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
      <PageHeader title="Volume Take-off" onSearch={setSearch} onAdd={handleOpenAdd} addLabel="Take-off" />

      {/* Project Selector Bar */}
      <div className="flex items-center gap-4 rounded-lg border bg-card p-3">
        <span className="text-sm font-medium text-muted-foreground whitespace-nowrap">Proyek:</span>
        <div className="w-72">
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
                    {p.name}
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
              <TableHead rowSpan={2} className="w-10 text-center">No</TableHead>
              <TableHead rowSpan={2}>Work Item</TableHead>
              <TableHead rowSpan={2} className="w-24 text-right">Planned Vol</TableHead>
              <TableHead rowSpan={2} className="w-16 text-center">Unit</TableHead>
              <TableHead rowSpan={2} className="w-24 text-center">Date</TableHead>
              <TableHead colSpan={3} className="text-center border-b">Dimensions (m)</TableHead>
              <TableHead rowSpan={2} className="w-24 text-right">Quantity</TableHead>
              <TableHead rowSpan={2} className="w-24 text-right">Unit Wt</TableHead>
              <TableHead rowSpan={2} className="w-28 text-right">Total Qty</TableHead>
              <TableHead rowSpan={2}>Description</TableHead>
              {isAdmin && <TableHead rowSpan={2} className="w-16 text-center">Action</TableHead>}
            </TableRow>
            <TableRow>
              <TableHead className="w-20 text-center text-xs">L</TableHead>
              <TableHead className="w-20 text-center text-xs">W</TableHead>
              <TableHead className="w-20 text-center text-xs">H</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 13 : 12} className="h-32 text-center text-muted-foreground">
                  Loading data...
                </TableCell>
              </TableRow>
            ) : volumeRows.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 13 : 12} className="h-48 text-center">
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
                  <TableCell className="text-right text-muted-foreground">{formatDecimal(r.planned_volume)}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{r.unit}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{r.date || '-'}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{r.length ? formatDecimal(r.length) : '-'}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{r.width ? formatDecimal(r.width) : '-'}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{r.height ? formatDecimal(r.height) : '-'}</TableCell>
                  <TableCell className="text-right text-foreground font-medium">{formatDecimal(r.quantity)}</TableCell>
                  <TableCell className="text-right text-muted-foreground">{formatDecimal(r.unit_weight)}</TableCell>
                  <TableCell className="text-right font-bold text-foreground">{formatDecimal(r.total_quantity)}</TableCell>
                  <TableCell className="text-muted-foreground text-xs">{r.description || '-'}</TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="text-xs text-destructive hover:bg-destructive/10"
                        onClick={() => handleDelete(r.id, r.work_item)}
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

      {/* Modal Add Volume Take-off */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-xl">
          <DialogHeader>
            <DialogTitle>Add Volume Take-off</DialogTitle>
            <DialogDescription>Input dimensi pengukuran untuk menghitung volume secara otomatis.</DialogDescription>
          </DialogHeader>

          {errorMessage && (
            <Alert variant="destructive">
              <AlertDescription>{errorMessage}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleCreate} className="flex flex-col gap-4">
            <FieldGroup>
              {ahsps.length > 0 && (
                <Field>
                  <FieldLabel htmlFor="vol-ahsp">Master AHSP (Opsional)</FieldLabel>
                  <Select value={selectedAhspId} onValueChange={(val) => { if (val) handleSelectAhsp(val) }}>
                    <SelectTrigger id="vol-ahsp">
                      <SelectValue placeholder="-- Pilih dari Master AHSP --" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        {ahsps.map((a) => (
                          <SelectItem key={a.id} value={a.id}>
                            {a.code} - {a.work_item}
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </Field>
              )}

              <div className="grid grid-cols-3 gap-4">
                <div className="col-span-2">
                  <Field>
                    <FieldLabel htmlFor="vol-workitem">Work Item</FieldLabel>
                    <Input
                      id="vol-workitem"
                      value={workItem}
                      onChange={(e) => setWorkItem(e.target.value)}
                      placeholder="Nama Pekerjaan..."
                      required
                    />
                  </Field>
                </div>
                <Field>
                  <FieldLabel htmlFor="vol-unit">Unit</FieldLabel>
                  <Input
                    id="vol-unit"
                    value={unit}
                    onChange={(e) => setUnit(e.target.value)}
                    placeholder="m3, m2, kg..."
                    required
                  />
                </Field>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="vol-planned">Planned Volume</FieldLabel>
                  <Input
                    id="vol-planned"
                    type="number"
                    step="0.0001"
                    value={plannedVolume}
                    onChange={(e) => setPlannedVolume(e.target.value)}
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="vol-date">Date</FieldLabel>
                  <Input
                    id="vol-date"
                    type="date"
                    value={date}
                    onChange={(e) => setDate(e.target.value)}
                  />
                </Field>
              </div>

              <div className="grid grid-cols-4 gap-3">
                <Field>
                  <FieldLabel htmlFor="vol-l">Length (m)</FieldLabel>
                  <Input
                    id="vol-l"
                    type="number"
                    step="0.0001"
                    value={length}
                    onChange={(e) => setLength(e.target.value)}
                    placeholder="L"
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="vol-w">Width (m)</FieldLabel>
                  <Input
                    id="vol-w"
                    type="number"
                    step="0.0001"
                    value={width}
                    onChange={(e) => setWidth(e.target.value)}
                    placeholder="W"
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="vol-h">Height (m)</FieldLabel>
                  <Input
                    id="vol-h"
                    type="number"
                    step="0.0001"
                    value={height}
                    onChange={(e) => setHeight(e.target.value)}
                    placeholder="H"
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="vol-uw">Unit Wt</FieldLabel>
                  <Input
                    id="vol-uw"
                    type="number"
                    step="0.0001"
                    value={unitWeight}
                    onChange={(e) => setUnitWeight(e.target.value)}
                  />
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="vol-desc">Description / Note</FieldLabel>
                <Input
                  id="vol-desc"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Catatan pengukuran..."
                />
              </Field>

              {/* Calculated preview card */}
              <Card className="bg-muted/40">
                <CardContent className="p-3">
                  <div className="flex items-center justify-between text-xs text-muted-foreground">
                    <span>Formula: {length && width && height ? `${length} × ${width} × ${height}` : 'Planned Volume'} × {unitWeight || 1}</span>
                    <div className="flex items-center gap-3">
                      <span>Qty: <strong className="text-foreground">{formatDecimal(preview.qty)}</strong></span>
                      <span>Total Qty: <strong className="text-primary font-bold">{formatDecimal(preview.totalQty)}</strong></span>
                    </div>
                  </div>
                </CardContent>
              </Card>
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
