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
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { formatRupiah } from '@/lib/utils'
import { useAuth } from '@/lib/auth'

export const Route = createFileRoute('/ahsp')({ component: AhspPage })

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

  const removeItem = (id: string) => {
    setItems((prev) => prev.filter((item) => item.id !== id))
  }

  const updateItem = (id: string, field: keyof LocalItem, value: any) => {
    setItems((prev) =>
      prev.map((item) => {
        if (item.id !== id) return item
        return { ...item, [field]: value }
      }),
    )
  }

  const handleSelectMaster = (id: string, masterId: string) => {
    const found = masterPrices.find((m) => m.id === masterId)
    if (!found) return
    setItems((prev) =>
      prev.map((item) => {
        if (item.id !== id) return item
        return {
          ...item,
          unit_price_id: found.id,
          code: found.code,
          description: found.work_item,
          unit: found.unit,
          unit_price: found.unit_price,
        }
      }),
    )
  }

  // Live calculation of A, B, C, D, E, F
  const calc = useMemo(() => {
    let labor = 0
    let material = 0
    let equipment = 0

    items.forEach((item) => {
      const coef = Number.parseFloat(item.coefficient) || 0
      const price = Number.parseInt(item.unit_price, 10) || 0
      const cost = Math.round(coef * price)

      if (item.category === 'LABOR') labor += cost
      else if (item.category === 'MATERIALS') material += cost
      else if (item.category === 'EQUIPMENT') equipment += cost
    })

    const totalCostD = labor + material + equipment
    const pct = Number.parseFloat(overheadPct) || 0
    const overheadE = Math.round((pct / 100) * totalCostD)
    const unitPriceF = totalCostD + overheadE

    return { labor, material, equipment, totalCostD, overheadE, unitPriceF }
  }, [items, overheadPct])

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage('')
    if (!code.trim() || !workItem.trim() || !unit.trim()) {
      setErrorMessage('Code, Work Item, dan Unit wajib diisi.')
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
          items: items.map(({ id: _id, ...rest }) => ({
            ...rest,
            coefficient: rest.coefficient || '0',
            unit_price: rest.unit_price || '0',
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
    <div>
      <PageHeader title="AHSP" onSearch={setSearch} onAdd={handleOpenAdd} />

      <div className="overflow-x-auto border bg-white shadow-sm">
        <Table className="text-xs">
          <TableHeader className="bg-white">
            <TableRow className="border-b hover:bg-transparent">
              <TableHead className="w-12 border-r text-center font-bold">NO</TableHead>
              <TableHead className="w-28 border-r text-center font-bold">CODE</TableHead>
              <TableHead className="min-w-64 border-r text-center font-bold">ITEM</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">UNIT</TableHead>
              <TableHead className="w-40 border-r text-center font-bold">COST (F)</TableHead>
              <TableHead className="w-24 border-r text-center font-bold">STATUS</TableHead>
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
            ) : ahsps.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 7 : 6} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  {search ? `Tidak ada AHSP yang cocok dengan "${search}".` : 'Belum ada data AHSP.'}
                </TableCell>
              </TableRow>
            ) : (
              ahsps.map((a, idx) => (
                <TableRow key={a.id} className="border-b bg-muted/40 hover:bg-muted/80">
                  <TableCell className="border-r text-center font-medium">{idx + 1}</TableCell>
                  <TableCell className="border-r font-mono font-semibold">{a.code}</TableCell>
                  <TableCell className="border-r font-medium text-foreground">{a.work_item}</TableCell>
                  <TableCell className="border-r text-center">{a.unit}</TableCell>
                  <TableCell className="border-r text-right font-bold text-foreground">
                    {formatRupiah(a.unit_price_of_work)}
                  </TableCell>
                  <TableCell className="border-r text-center">
                    <Badge variant={a.status === 'final' ? 'default' : 'secondary'} className="text-[10px]">
                      {a.status}
                    </Badge>
                  </TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-7 text-xs text-red-600 hover:bg-red-50 hover:text-red-700"
                        onClick={() => handleDelete(a.id, a.code)}
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

      {/* Modal Form AHSP Komponen A-F (Mockup Hal. 7) */}
      <Dialog open={isFormOpen} onOpenChange={setIsFormOpen}>
        <DialogContent className="max-h-[90vh] max-w-4xl overflow-y-auto rounded-2xl">
          <DialogHeader>
            <DialogTitle className="text-lg font-bold tracking-wide">ANALISA HARGA SATUAN PEKERJAAN (AHSP)</DialogTitle>
          </DialogHeader>

          {errorMessage && (
            <div className="rounded-md border border-red-200 bg-red-50 p-3 text-xs text-red-600">
              {errorMessage}
            </div>
          )}

          <form onSubmit={handleSave} className="space-y-4 text-xs">
            {/* Header Form */}
            <div className="grid grid-cols-4 gap-4 rounded-xl bg-muted/30 p-3">
              <div>
                <Label className="text-xs font-bold">CODE</Label>
                <Input
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  placeholder="misal AHSP-BETON-01"
                  className="mt-1 bg-white font-mono"
                  required
                />
              </div>
              <div className="col-span-2">
                <Label className="text-xs font-bold">WORK ITEM</Label>
                <Input
                  value={workItem}
                  onChange={(e) => setWorkItem(e.target.value)}
                  placeholder="Uraian Pekerjaan..."
                  className="mt-1 bg-white font-medium"
                  required
                />
              </div>
              <div>
                <Label className="text-xs font-bold">UNIT</Label>
                <Input
                  value={unit}
                  onChange={(e) => setUnit(e.target.value)}
                  placeholder="m3, m2, m', kg..."
                  className="mt-1 bg-white"
                  required
                />
              </div>
            </div>

            {/* Components Table A-F */}
            <div className="overflow-x-auto border">
              <Table className="text-xs">
                <TableHeader className="bg-white">
                  <TableRow className="border-b">
                    <TableHead className="w-10 border-r text-center font-bold">NO</TableHead>
                    <TableHead className="min-w-44 border-r text-center font-bold">DESCRIPTION</TableHead>
                    <TableHead className="w-24 border-r text-center font-bold">CODE</TableHead>
                    <TableHead className="w-16 border-r text-center font-bold">UNIT</TableHead>
                    <TableHead className="w-24 border-r text-center font-bold">COEFFICIENT</TableHead>
                    <TableHead className="w-28 border-r text-center font-bold">UNIT PRICE</TableHead>
                    <TableHead className="w-28 border-r text-center font-bold">COST</TableHead>
                    <TableHead className="w-12 text-center font-bold"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {/* Section A: LABOR */}
                  <TableRow className="bg-primary/10 font-bold">
                    <TableCell className="border-r text-center">A</TableCell>
                    <TableCell colSpan={6}>LABOR (TENAGA KERJA)</TableCell>
                    <TableCell className="text-center">
                      <Button type="button" size="sm" variant="outline" className="h-6 text-[10px]" onClick={() => addItem('LABOR')}>
                        + Item
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
                      onSelectMaster={handleSelectMaster}
                    />
                  ))}
                  <TableRow className="bg-muted/70 font-semibold">
                    <TableCell colSpan={6} className="border-r text-right">TOTAL LABOR COST</TableCell>
                    <TableCell className="border-r text-right">{formatRupiah(calc.labor)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* Section B: MATERIALS */}
                  <TableRow className="bg-primary/10 font-bold">
                    <TableCell className="border-r text-center">B</TableCell>
                    <TableCell colSpan={6}>MATERIALS (BAHAN / MATERIAL)</TableCell>
                    <TableCell className="text-center">
                      <Button type="button" size="sm" variant="outline" className="h-6 text-[10px]" onClick={() => addItem('MATERIALS')}>
                        + Item
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
                      onSelectMaster={handleSelectMaster}
                    />
                  ))}
                  <TableRow className="bg-muted/70 font-semibold">
                    <TableCell colSpan={6} className="border-r text-right">TOTAL MATERIAL COST</TableCell>
                    <TableCell className="border-r text-right">{formatRupiah(calc.material)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* Section C: EQUIPMENT */}
                  <TableRow className="bg-primary/10 font-bold">
                    <TableCell className="border-r text-center">C</TableCell>
                    <TableCell colSpan={6}>EQUIPMENT (PERALATAN)</TableCell>
                    <TableCell className="text-center">
                      <Button type="button" size="sm" variant="outline" className="h-6 text-[10px]" onClick={() => addItem('EQUIPMENT')}>
                        + Item
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
                      onSelectMaster={handleSelectMaster}
                    />
                  ))}
                  <TableRow className="bg-muted/70 font-semibold">
                    <TableCell colSpan={6} className="border-r text-right">TOTAL EQUIPMENT COST</TableCell>
                    <TableCell className="border-r text-right">{formatRupiah(calc.equipment)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* D: TOTAL COST */}
                  <TableRow className="bg-muted font-bold text-foreground">
                    <TableCell className="border-r text-center">D</TableCell>
                    <TableCell colSpan={5} className="border-r">TOTAL COST (A + B + C)</TableCell>
                    <TableCell className="border-r text-right">{formatRupiah(calc.totalCostD)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* E: OVERHEAD + PROFIT */}
                  <TableRow className="bg-white">
                    <TableCell className="border-r text-center font-bold">E</TableCell>
                    <TableCell colSpan={4} className="border-r">
                      <div className="flex items-center gap-2">
                        <span>OVERHEAD + PROFIT</span>
                        <Input
                          type="number"
                          step="0.1"
                          value={overheadPct}
                          onChange={(e) => setOverheadPct(e.target.value)}
                          className="h-6 w-16 text-center font-bold text-xs"
                        />
                        <span>%</span>
                      </div>
                    </TableCell>
                    <TableCell className="border-r text-center">{overheadPct}%</TableCell>
                    <TableCell className="border-r text-right font-medium">{formatRupiah(calc.overheadE)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* F: UNIT PRICE OF WORK */}
                  <TableRow className="bg-primary text-primary-foreground font-extrabold text-sm">
                    <TableCell className="border-r text-center">F</TableCell>
                    <TableCell colSpan={5} className="border-r">UNIT PRICE OF WORK (D + E)</TableCell>
                    <TableCell className="border-r text-right">{formatRupiah(calc.unitPriceF)}</TableCell>
                    <TableCell></TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>

            {/* Total Display & Actions */}
            <div className="flex items-center justify-between rounded-xl bg-muted/40 p-3">
              <div className="flex items-center gap-4">
                <span className="font-bold">TOTAL COST:</span>
                <span className="text-base font-extrabold text-primary">{formatRupiah(calc.unitPriceF)}</span>
              </div>
              <div className="flex gap-2">
                <Button type="button" variant="outline" onClick={() => setIsFormOpen(false)}>
                  Batal
                </Button>
                <Button type="submit" disabled={createMut.isPending} className="px-8">
                  {createMut.isPending ? 'Menyimpan...' : 'ADD'}
                </Button>
              </div>
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
  onSelectMaster,
}: {
  item: LocalItem
  no: number
  masterPrices: any[]
  onUpdate: (id: string, field: keyof LocalItem, val: any) => void
  onRemove: (id: string) => void
  onSelectMaster: (id: string, masterId: string) => void
}) {
  const coef = Number.parseFloat(item.coefficient) || 0
  const price = Number.parseInt(item.unit_price, 10) || 0
  const cost = Math.round(coef * price)

  const relevantMasters = masterPrices.filter((m) => m.category === item.category)

  return (
    <TableRow className="border-b bg-white hover:bg-muted/30">
      <TableCell className="border-r text-center font-medium">{no}</TableCell>
      <TableCell className="border-r p-1">
        <div className="space-y-1">
          <Input
            value={item.description}
            onChange={(e) => onUpdate(item.id, 'description', e.target.value)}
            placeholder="Deskripsi komponen..."
            className="h-7 text-xs"
          />
          {relevantMasters.length > 0 && (
            <select
              className="h-6 w-full rounded border border-gray-200 bg-muted/30 text-[10px] text-muted-foreground"
              onChange={(e) => onSelectMaster(item.id, e.target.value)}
              defaultValue=""
            >
              <option value="" disabled>-- Pilih dari Master Harga --</option>
              {relevantMasters.map((m) => (
                <option key={m.id} value={m.id}>
                  {m.code} - {m.work_item} ({formatRupiah(m.unit_price)}/{m.unit})
                </option>
              ))}
            </select>
          )}
        </div>
      </TableCell>
      <TableCell className="border-r p-1">
        <Input
          value={item.code || ''}
          onChange={(e) => onUpdate(item.id, 'code', e.target.value)}
          placeholder="Code"
          className="h-7 font-mono text-xs"
        />
      </TableCell>
      <TableCell className="border-r p-1">
        <Input
          value={item.unit}
          onChange={(e) => onUpdate(item.id, 'unit', e.target.value)}
          placeholder="Unit"
          className="h-7 text-xs"
        />
      </TableCell>
      <TableCell className="border-r p-1">
        <Input
          value={item.coefficient}
          onChange={(e) => onUpdate(item.id, 'coefficient', e.target.value)}
          placeholder="0.00"
          className="h-7 text-right text-xs"
        />
      </TableCell>
      <TableCell className="border-r p-1">
        <Input
          value={item.unit_price}
          onChange={(e) => onUpdate(item.id, 'unit_price', e.target.value)}
          placeholder="Harga"
          className="h-7 text-right text-xs"
        />
      </TableCell>
      <TableCell className="border-r text-right font-medium">{formatRupiah(cost)}</TableCell>
      <TableCell className="text-center p-1">
        <Button
          type="button"
          size="sm"
          variant="ghost"
          className="h-6 w-6 p-0 text-red-500 hover:text-red-700"
          onClick={() => onRemove(item.id)}
        >
          ×
        </Button>
      </TableCell>
    </TableRow>
  )
}
