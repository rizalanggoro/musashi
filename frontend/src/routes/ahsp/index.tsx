import { useState, useMemo } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listAhspsOptions,
  createAhspMutation,
  deleteAhspMutation,
  listUnitPricesOptions,
} from '@/api/@tanstack/react-query.gen'
import type { CreateAhspItemInput } from '@/api/types.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FieldGroup, Field, FieldLabel } from '@/components/ui/field'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { Badge } from '@/components/ui/badge'
import { formatRupiah } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import { Plus, Trash2 } from 'lucide-react'

export const Route = createFileRoute('/ahsp/')({ component: AhspPage })

type LocalItem = CreateAhspItemInput & { id: string }

function AhspPage() {
  const queryClient = useQueryClient()
  const { isAdmin, user, openLoginModal } = useAuth()

  const [search, setSearch] = useState('')
  const [isFormOpen, setIsFormOpen] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')

  // Form Header State
  const [code, setCode] = useState('')
  const [workItem, setWorkItem] = useState('')
  const [unit, setUnit] = useState('')
  const [overheadPct, setOverheadPct] = useState('10.00')
  const [status, setStatus] = useState<'draft' | 'final'>('draft')
  const [items, setItems] = useState<LocalItem[]>([])

  const { data: ahsps = [], isLoading } = useQuery({
    ...listAhspsOptions({ query: { q: search || undefined } }),
  })

  // Master unit prices for auto-complete/selection
  const { data: masterPrices = [] } = useQuery({
    ...listUnitPricesOptions({}),
    enabled: isFormOpen,
  })

  const createMut = useMutation(createAhspMutation())
  const deleteMut = useMutation(deleteAhspMutation())

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal()
      return
    }
    setCode('')
    setWorkItem('')
    setUnit('')
    setOverheadPct('10.00')
    setStatus('draft')
    setItems([
      {
        id: crypto.randomUUID(),
        category: 'LABOR',
        sort: 1,
        description: 'Tukang',
        code: '',
        unit: 'OH',
        coefficient: '0.2000',
        unit_price: '120000',
      },
      {
        id: crypto.randomUUID(),
        category: 'MATERIALS',
        sort: 2,
        description: 'Semen Portland',
        code: '',
        unit: 'sak',
        coefficient: '1.0000',
        unit_price: '65000',
      },
    ])
    setErrorMessage('')
    setIsFormOpen(true)
  }

  const addItem = (category: 'LABOR' | 'MATERIALS' | 'EQUIPMENT') => {
    setItems((prev) => [
      ...prev,
      {
        id: crypto.randomUUID(),
        category,
        sort: prev.length + 1,
        description: '',
        code: '',
        unit: '',
        coefficient: '1.0000',
        unit_price: '0',
      },
    ])
  }

  const updateItem = (id: string, field: keyof LocalItem, value: any) => {
    setItems((prev) =>
      prev.map((it) => {
        if (it.id !== id) return it
        const updated = { ...it, [field]: value }
        // Auto-fill from master price if matching code or item selected
        if (field === 'unit_price_id' && value) {
          const master = masterPrices.find((m) => m.id === value)
          if (master) {
            updated.description = master.work_item
            updated.code = master.code
            updated.unit = master.unit
            updated.unit_price = master.unit_price
          }
        }
        return updated
      })
    )
  }

  const removeItem = (id: string) => {
    setItems((prev) => prev.filter((it) => it.id !== id))
  }

  // Calculations for A, B, C, D, E, F
  const calc = useMemo(() => {
    let labor = 0
    let material = 0
    let equipment = 0

    items.forEach((it) => {
      const coef = parseFloat(it.coefficient || '0')
      const price = parseInt(it.unit_price || '0', 10)
      const cost = Math.round(coef * price)
      if (it.category === 'LABOR') labor += cost
      else if (it.category === 'MATERIALS') material += cost
      else if (it.category === 'EQUIPMENT') equipment += cost
    })

    const d = labor + material + equipment
    const pct = parseFloat(overheadPct || '0')
    const e = Math.round((pct / 100.0) * d)
    const f = d + e

    return { labor, material, equipment, d, e, f }
  }, [items, overheadPct])

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage('')
    if (!code.trim() || !workItem.trim() || !unit.trim()) {
      setErrorMessage('Code, item pekerjaan, dan unit wajib diisi.')
      return
    }

    try {
      await createMut.mutateAsync({
        body: {
          code: code.trim(),
          work_item: workItem.trim(),
          unit: unit.trim(),
          overhead_profit_pct: overheadPct,
          status,
          items: items.map((it, idx) => ({
            unit_price_id: it.unit_price_id || undefined,
            category: it.category,
            sort: idx + 1,
            description: it.description,
            code: it.code,
            unit: it.unit,
            coefficient: it.coefficient,
            unit_price: it.unit_price,
          })),
        },
      })
      await queryClient.invalidateQueries()
      setIsFormOpen(false)
    } catch (err: any) {
      setErrorMessage(err?.message || 'Gagal menyimpan AHSP.')
    }
  }

  const handleDelete = async (id: string, code: string) => {
    if (!confirm(`Hapus AHSP ${code}?`)) return
    try {
      await deleteMut.mutateAsync({ path: { id } })
      await queryClient.invalidateQueries()
    } catch (err: any) {
      alert(err?.message || 'Gagal menghapus AHSP.')
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <PageHeader title="AHSP" onSearch={setSearch} onAdd={handleOpenAdd} addLabel="AHSP" />

      <div className="rounded-md border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12 text-center">No</TableHead>
              <TableHead className="w-32">Code</TableHead>
              <TableHead>Work Item</TableHead>
              <TableHead className="w-24 text-center">Unit</TableHead>
              <TableHead className="w-40 text-right">Unit Price of Work</TableHead>
              <TableHead className="w-24 text-center">Status</TableHead>
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
            ) : ahsps.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 7 : 6} className="h-48 text-center">
                  <Empty>
                    <EmptyHeader>
                      <EmptyTitle>Belum ada data AHSP</EmptyTitle>
                      <EmptyDescription>
                        {search ? `Tidak ada AHSP yang cocok dengan "${search}".` : 'Mulai menyusun analisa harga satuan pekerjaan baru.'}
                      </EmptyDescription>
                    </EmptyHeader>
                  </Empty>
                </TableCell>
              </TableRow>
            ) : (
              ahsps.map((a, idx) => (
                <TableRow key={a.id}>
                  <TableCell className="text-center text-muted-foreground">{idx + 1}</TableCell>
                  <TableCell className="font-mono text-xs font-semibold">{a.code}</TableCell>
                  <TableCell className="font-medium text-foreground">{a.work_item}</TableCell>
                  <TableCell className="text-center">{a.unit}</TableCell>
                  <TableCell className="text-right font-medium text-foreground">
                    {formatRupiah(a.unit_price_of_work)}
                  </TableCell>
                  <TableCell className="text-center">
                    <Badge variant={a.status === 'final' ? 'default' : 'secondary'} className="capitalize text-xs font-normal">
                      {a.status}
                    </Badge>
                  </TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="text-xs text-destructive hover:bg-destructive/10"
                        onClick={() => handleDelete(a.id, a.code)}
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

      {/* Modal Form AHSP Komponen A-F */}
      <Dialog open={isFormOpen} onOpenChange={setIsFormOpen}>
        <DialogContent className="max-h-[90vh] max-w-4xl overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Analisa Harga Satuan Pekerjaan (AHSP)</DialogTitle>
            <DialogDescription>Susun rincian koefisien komponen A (Tenaga), B (Bahan), dan C (Peralatan).</DialogDescription>
          </DialogHeader>

          {errorMessage && (
            <Alert variant="destructive">
              <AlertDescription>{errorMessage}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleSave} className="flex flex-col gap-4">
            {/* Header Form */}
            <FieldGroup className="rounded-lg border bg-muted/40 p-4">
              <div className="grid grid-cols-4 gap-4">
                <Field>
                  <FieldLabel htmlFor="ahsp-code">Code</FieldLabel>
                  <Input
                    id="ahsp-code"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    placeholder="e.g. AHSP-BETON-01"
                    className="font-mono bg-card"
                    required
                  />
                </Field>
                <div className="col-span-2">
                  <Field>
                    <FieldLabel htmlFor="ahsp-item">Work Item</FieldLabel>
                    <Input
                      id="ahsp-item"
                      value={workItem}
                      onChange={(e) => setWorkItem(e.target.value)}
                      placeholder="Uraian Pekerjaan..."
                      className="bg-card"
                      required
                    />
                  </Field>
                </div>
                <Field>
                  <FieldLabel htmlFor="ahsp-unit">Unit</FieldLabel>
                  <Input
                    id="ahsp-unit"
                    value={unit}
                    onChange={(e) => setUnit(e.target.value)}
                    placeholder="m3, m2, kg..."
                    className="bg-card"
                    required
                  />
                </Field>
              </div>
            </FieldGroup>

            {/* Components Table A-F */}
            <div className="rounded-md border bg-card overflow-hidden">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-10 text-center">No</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead className="w-24">Code</TableHead>
                    <TableHead className="w-16 text-center">Unit</TableHead>
                    <TableHead className="w-24 text-right">Coefficient</TableHead>
                    <TableHead className="w-28 text-right">Unit Price</TableHead>
                    <TableHead className="w-32 text-right">Cost</TableHead>
                    <TableHead className="w-10 text-center"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {/* Section A: LABOR */}
                  <TableRow className="bg-muted/60 font-semibold">
                    <TableCell className="text-center font-bold">A</TableCell>
                    <TableCell colSpan={6} className="text-foreground">LABOR (Tenaga Kerja)</TableCell>
                    <TableCell className="text-center p-1">
                      <Button type="button" size="icon" variant="ghost" onClick={() => addItem('LABOR')}>
                        <Plus data-icon="inline-start" />
                      </Button>
                    </TableCell>
                  </TableRow>
                  {items.filter((it) => it.category === 'LABOR').map((item, idx) => (
                    <ComponentRow
                      key={item.id}
                      item={item}
                      no={idx + 1}
                      masterPrices={masterPrices}
                      onUpdate={updateItem}
                      onRemove={removeItem}
                    />
                  ))}
                  <TableRow className="border-t font-medium text-xs">
                    <TableCell colSpan={6} className="text-right text-muted-foreground">Jumlah Tenaga Kerja (A):</TableCell>
                    <TableCell className="text-right font-semibold text-foreground">{formatRupiah(calc.labor)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* Section B: MATERIALS */}
                  <TableRow className="bg-muted/60 font-semibold border-t">
                    <TableCell className="text-center font-bold">B</TableCell>
                    <TableCell colSpan={6} className="text-foreground">MATERIALS (Bahan)</TableCell>
                    <TableCell className="text-center p-1">
                      <Button type="button" size="icon" variant="ghost" onClick={() => addItem('MATERIALS')}>
                        <Plus data-icon="inline-start" />
                      </Button>
                    </TableCell>
                  </TableRow>
                  {items.filter((it) => it.category === 'MATERIALS').map((item, idx) => (
                    <ComponentRow
                      key={item.id}
                      item={item}
                      no={idx + 1}
                      masterPrices={masterPrices}
                      onUpdate={updateItem}
                      onRemove={removeItem}
                    />
                  ))}
                  <TableRow className="border-t font-medium text-xs">
                    <TableCell colSpan={6} className="text-right text-muted-foreground">Jumlah Bahan (B):</TableCell>
                    <TableCell className="text-right font-semibold text-foreground">{formatRupiah(calc.material)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* Section C: EQUIPMENT */}
                  <TableRow className="bg-muted/60 font-semibold border-t">
                    <TableCell className="text-center font-bold">C</TableCell>
                    <TableCell colSpan={6} className="text-foreground">EQUIPMENT (Peralatan)</TableCell>
                    <TableCell className="text-center p-1">
                      <Button type="button" size="icon" variant="ghost" onClick={() => addItem('EQUIPMENT')}>
                        <Plus data-icon="inline-start" />
                      </Button>
                    </TableCell>
                  </TableRow>
                  {items.filter((it) => it.category === 'EQUIPMENT').map((item, idx) => (
                    <ComponentRow
                      key={item.id}
                      item={item}
                      no={idx + 1}
                      masterPrices={masterPrices}
                      onUpdate={updateItem}
                      onRemove={removeItem}
                    />
                  ))}
                  <TableRow className="border-t font-medium text-xs">
                    <TableCell colSpan={6} className="text-right text-muted-foreground">Jumlah Peralatan (C):</TableCell>
                    <TableCell className="text-right font-semibold text-foreground">{formatRupiah(calc.equipment)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* Summary D, E, F */}
                  <TableRow className="border-t-2 font-semibold">
                    <TableCell className="text-center">D</TableCell>
                    <TableCell colSpan={5}>Jumlah Biaya Langsung (A + B + C)</TableCell>
                    <TableCell className="text-right font-bold text-foreground">{formatRupiah(calc.d)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  <TableRow className="font-medium text-sm">
                    <TableCell className="text-center">E</TableCell>
                    <TableCell colSpan={3}>
                      <div className="flex items-center gap-2">
                        <span>Biaya Umum & Keuntungan:</span>
                        <Input
                          type="number"
                          step="0.01"
                          value={overheadPct}
                          onChange={(e) => setOverheadPct(e.target.value)}
                          className="h-7 w-20 text-center font-semibold text-xs"
                        />
                        <span>%</span>
                      </div>
                    </TableCell>
                    <TableCell colSpan={2} className="text-right text-muted-foreground">
                      {overheadPct}% × D
                    </TableCell>
                    <TableCell className="text-right font-semibold text-foreground">{formatRupiah(calc.e)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  <TableRow className="border-t-2 bg-muted/40 font-bold text-base">
                    <TableCell className="text-center">F</TableCell>
                    <TableCell colSpan={5} className="text-foreground">
                      Harga Satuan Pekerjaan (D + E)
                    </TableCell>
                    <TableCell className="text-right text-foreground font-extrabold">{formatRupiah(calc.f)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>

            {/* Status and Action Buttons */}
            <div className="flex items-center justify-between pt-2">
              <div className="flex items-center gap-4">
                <label className="flex items-center gap-2 text-sm font-medium">
                  <input
                    type="radio"
                    name="ahsp-status"
                    checked={status === 'draft'}
                    onChange={() => setStatus('draft')}
                    className="size-4"
                  />
                  Draft
                </label>
                <label className="flex items-center gap-2 text-sm font-medium">
                  <input
                    type="radio"
                    name="ahsp-status"
                    checked={status === 'final'}
                    onChange={() => setStatus('final')}
                    className="size-4"
                  />
                  Final
                </label>
              </div>

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setIsFormOpen(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={createMut.isPending}>
                  {createMut.isPending ? 'Saving...' : 'Save AHSP'}
                </Button>
              </DialogFooter>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function ComponentRow({
  item,
  no,
  masterPrices,
  onUpdate,
  onRemove,
}: {
  item: LocalItem
  no: number
  masterPrices: any[]
  onUpdate: (id: string, field: keyof LocalItem, val: any) => void
  onRemove: (id: string) => void
}) {
  const coef = parseFloat(item.coefficient || '0')
  const price = parseInt(item.unit_price || '0', 10)
  const cost = Math.round(coef * price)

  // Filter master prices by category
  const filteredPrices = masterPrices.filter((m) => m.category === item.category)

  return (
    <TableRow className="text-xs">
      <TableCell className="text-center text-muted-foreground">{no}</TableCell>
      <TableCell className="p-1">
        <div className="flex gap-1.5">
          {filteredPrices.length > 0 && (
            <select
              value={item.unit_price_id || ''}
              onChange={(e) => onUpdate(item.id, 'unit_price_id', e.target.value)}
              className="h-8 max-w-[140px] rounded-md border border-input bg-card px-1.5 text-xs text-foreground"
            >
              <option value="">-- Pilih Master --</option>
              {filteredPrices.map((m) => (
                <option key={m.id} value={m.id}>
                  {m.code} - {m.work_item}
                </option>
              ))}
            </select>
          )}
          <Input
            value={item.description}
            onChange={(e) => onUpdate(item.id, 'description', e.target.value)}
            placeholder="Deskripsi..."
            className="h-8 text-xs bg-card"
            required
          />
        </div>
      </TableCell>
      <TableCell className="p-1">
        <Input
          value={item.code || ''}
          onChange={(e) => onUpdate(item.id, 'code', e.target.value)}
          placeholder="Kode..."
          className="h-8 text-xs font-mono bg-card"
        />
      </TableCell>
      <TableCell className="p-1 text-center">
        <Input
          value={item.unit}
          onChange={(e) => onUpdate(item.id, 'unit', e.target.value)}
          placeholder="OH/m3"
          className="h-8 text-xs text-center bg-card"
          required
        />
      </TableCell>
      <TableCell className="p-1 text-right">
        <Input
          type="number"
          step="0.0001"
          value={item.coefficient}
          onChange={(e) => onUpdate(item.id, 'coefficient', e.target.value)}
          className="h-8 text-xs text-right font-mono bg-card"
          required
        />
      </TableCell>
      <TableCell className="p-1 text-right">
        <Input
          type="number"
          value={item.unit_price}
          onChange={(e) => onUpdate(item.id, 'unit_price', e.target.value)}
          className="h-8 text-xs text-right font-mono bg-card"
          required
        />
      </TableCell>
      <TableCell className="text-right font-medium text-foreground">
        {formatRupiah(cost)}
      </TableCell>
      <TableCell className="p-1 text-center">
        <Button
          type="button"
          size="icon"
          variant="ghost"
          className="text-muted-foreground hover:text-destructive"
          onClick={() => onRemove(item.id)}
        >
          <Trash2 data-icon="inline-start" />
        </Button>
      </TableCell>
    </TableRow>
  )
}
