import { useState, useMemo } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  createAhspMutation,
  updateAhspMutation,
  getAhspOptions,
  listUnitPricesOptions,
} from '@/api/@tanstack/react-query.gen'
import type { AhspDetailDto, CreateAhspItemInput } from '@/api/types.gen'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { FieldGroup, Field, FieldLabel } from '@/components/ui/field'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'
import { Badge } from '@/components/ui/badge'
import { formatRupiah } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import { ArrowLeft, Plus, Trash2, Save, X, HardHat, Package, Wrench } from 'lucide-react'

type LocalItem = CreateAhspItemInput & { id: string }

type AhspFormProps = {
  id?: string
}

export function AhspForm({ id }: AhspFormProps) {
  const isEdit = Boolean(id)

  const { data: ahspDetail, isLoading: isLoadingDetail } = useQuery({
    ...getAhspOptions({ path: { id: id! } }),
    enabled: isEdit,
  })

  if (isEdit && isLoadingDetail) {
    return (
      <div className="flex h-64 flex-col items-center justify-center gap-3">
        <div className="size-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        <p className="text-sm text-muted-foreground">Memuat data AHSP...</p>
      </div>
    )
  }

  return <AhspFormContent id={id} initialData={ahspDetail} />
}

function AhspFormContent({ id, initialData }: { id?: string; initialData?: AhspDetailDto }) {
  const isEdit = Boolean(id)
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user, openLoginModal } = useAuth()

  const [errorMessage, setErrorMessage] = useState('')

  // Header State
  const [code, setCode] = useState(initialData?.code || '')
  const [workItem, setWorkItem] = useState(initialData?.work_item || '')
  const [unit, setUnit] = useState(initialData?.unit || '')
  const [overheadPct, setOverheadPct] = useState(initialData?.overhead_profit_pct || '10.00')
  const [note, setNote] = useState(initialData?.note || '')
  const [status, setStatus] = useState<'draft' | 'final'>((initialData?.status as any) || 'draft')
  const [items, setItems] = useState<LocalItem[]>(() => {
    if (initialData?.items && initialData.items.length > 0) {
      return initialData.items.map((it) => ({
        id: it.id || crypto.randomUUID(),
        category: (it.category as any) || 'MATERIALS',
        sort: it.sort || 1,
        description: it.description || '',
        code: it.code || '',
        unit: it.unit || '',
        coefficient: it.coefficient || '1.0000',
        unit_price: it.unit_price || '0',
        unit_price_id: it.unit_price_id || undefined,
      }))
    }
    if (isEdit) return []
    return [
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
    ]
  })

  // Master unit prices for auto-complete/selection
  const { data: masterPrices = [] } = useQuery({
    ...listUnitPricesOptions({}),
  })

  const createMut = useMutation(createAhspMutation())
  const updateMut = useMutation(updateAhspMutation())
  const isPending = isEdit ? updateMut.isPending : createMut.isPending

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

  const updateItem = (itemId: string, field: keyof LocalItem, value: any) => {
    setItems((prev) =>
      prev.map((it) => {
        if (it.id !== itemId) return it
        const updated = { ...it, [field]: value }
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

  const removeItem = (itemId: string) => {
    setItems((prev) => prev.filter((it) => it.id !== itemId))
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

    if (!user) {
      openLoginModal()
      return
    }

    if (!code.trim() || !workItem.trim() || !unit.trim()) {
      setErrorMessage('Kode AHSP, uraian pekerjaan, dan satuan unit wajib diisi.')
      return
    }

    try {
      const payload = {
        code: code.trim(),
        work_item: workItem.trim(),
        unit: unit.trim(),
        overhead_profit_pct: overheadPct,
        status,
        note: note.trim(),
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
      }

      if (isEdit && id) {
        await updateMut.mutateAsync({
          path: { id },
          body: payload,
        })
      } else {
        await createMut.mutateAsync({
          body: payload,
        })
      }

      await queryClient.invalidateQueries()
      navigate({ to: '/ahsp' })
    } catch (err: any) {
      setErrorMessage(
        err?.message || (isEdit ? 'Gagal memperbarui data AHSP.' : 'Gagal menyimpan data AHSP baru.')
      )
    }
  }

  return (
    <div className="flex flex-col gap-6 pb-16">
      {/* Top Navigation & Action Header */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => navigate({ to: '/ahsp' })}
            className="gap-1.5"
          >
            <ArrowLeft className="size-4" />
            <span>Kembali ke AHSP</span>
          </Button>
          <div>
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              {isEdit ? 'Edit Analisa Harga Satuan (AHSP)' : 'Tambah AHSP Baru'}
            </h1>
            <p className="text-xs text-muted-foreground">
              {isEdit
                ? `Perbarui rincian item koefisien dan harga untuk ${code || 'AHSP'}`
                : 'Susun rincian koefisien komponen Tenaga Kerja (A), Bahan (B), dan Peralatan (C)'}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate({ to: '/ahsp' })}
            disabled={isPending}
          >
            <X className="size-4 mr-1.5" />
            Batal
          </Button>
          <Button type="button" onClick={handleSave} disabled={isPending}>
            <Save className="size-4 mr-1.5" />
            {isPending ? 'Menyimpan...' : isEdit ? 'Simpan Perubahan' : 'Simpan AHSP'}
          </Button>
        </div>
      </div>

      {errorMessage && (
        <Alert variant="destructive">
          <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
      )}

      <form onSubmit={handleSave} className="flex flex-col gap-6">
        {/* Section 1: Header AHSP Info */}
        <Card>
          <CardHeader>
            <CardTitle>Informasi Pekerjaan (Header)</CardTitle>
            <CardDescription>
              Tentukan kode klasifikasi, uraian pekerjaan, satuan pengukuran, dan catatan spesifikasi.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <FieldGroup className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <Field>
                <FieldLabel htmlFor="ahsp-code">
                  Kode AHSP <span className="text-destructive">*</span>
                </FieldLabel>
                <Input
                  id="ahsp-code"
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  placeholder="e.g. AHSP-BETON-01"
                  className="font-mono bg-card"
                  required
                />
              </Field>

              <div className="md:col-span-2">
                <Field>
                  <FieldLabel htmlFor="ahsp-work-item">
                    Uraian Pekerjaan <span className="text-destructive">*</span>
                  </FieldLabel>
                  <Input
                    id="ahsp-work-item"
                    value={workItem}
                    onChange={(e) => setWorkItem(e.target.value)}
                    placeholder="e.g. Pembuatan 1 m3 Beton Mutu K-225"
                    className="bg-card"
                    required
                  />
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="ahsp-unit">
                  Satuan <span className="text-destructive">*</span>
                </FieldLabel>
                <Input
                  id="ahsp-unit"
                  value={unit}
                  onChange={(e) => setUnit(e.target.value)}
                  placeholder="e.g. m3, m2, kg, OH"
                  className="bg-card"
                  required
                />
              </Field>

              <div className="md:col-span-3">
                <Field>
                  <FieldLabel htmlFor="ahsp-note">Catatan / Spesifikasi Teknis</FieldLabel>
                  <Textarea
                    id="ahsp-note"
                    value={note}
                    onChange={(e) => setNote(e.target.value)}
                    placeholder="Catatan metode kerja, spesifikasi standar, atau referensi SNI..."
                    rows={2}
                    className="bg-card"
                  />
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="ahsp-status">Status</FieldLabel>
                <NativeSelect
                  id="ahsp-status"
                  value={status}
                  onChange={(e) => setStatus(e.target.value as any)}
                  className="w-full bg-card"
                >
                  <NativeSelectOption value="draft">Draft</NativeSelectOption>
                  <NativeSelectOption value="final">Final / Disetujui</NativeSelectOption>
                </NativeSelect>
              </Field>
            </FieldGroup>
          </CardContent>
        </Card>

        {/* Section 2: Components Table (A, B, C) */}
        <Card>
          <CardHeader>
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div>
                <CardTitle>Rincian Komponen Koefisien & Harga Satuan</CardTitle>
                <CardDescription>
                  Masukkan koefisien dan harga satuan untuk setiap rincian tenaga kerja, bahan, dan peralatan.
                </CardDescription>
              </div>
              <div className="flex items-center gap-2 text-xs">
                <Badge variant="outline" className="gap-1 font-mono">
                  Total Items: {items.length}
                </Badge>
              </div>
            </div>
          </CardHeader>
          <CardContent className="p-0 sm:p-4">
            <div className="rounded-lg border bg-card overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow className="bg-muted/50">
                    <TableHead className="w-12 text-center">No</TableHead>
                    <TableHead className="min-w-[280px]">Uraian Komponen & Master Harga</TableHead>
                    <TableHead className="w-28">Kode</TableHead>
                    <TableHead className="w-20 text-center">Satuan</TableHead>
                    <TableHead className="w-28 text-right">Koefisien</TableHead>
                    <TableHead className="w-36 text-right">Harga Satuan (Rp)</TableHead>
                    <TableHead className="w-36 text-right">Jumlah Biaya (Rp)</TableHead>
                    <TableHead className="w-12 text-center"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {/* Category A: LABOR */}
                  <TableRow className="bg-muted/70 font-semibold border-t">
                    <TableCell className="text-center font-bold text-primary">A</TableCell>
                    <TableCell colSpan={6} className="text-foreground">
                      <div className="flex items-center gap-2">
                        <HardHat className="size-4 text-primary" />
                        <span>TENAGA KERJA (LABOR)</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-center p-1">
                      <Button
                        type="button"
                        size="sm"
                        variant="ghost"
                        onClick={() => addItem('LABOR')}
                        title="Tambah Tenaga Kerja"
                        className="h-8 w-8 p-0 text-primary hover:text-primary hover:bg-primary/10"
                      >
                        <Plus className="size-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                  {items.filter((it) => it.category === 'LABOR').length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={8} className="text-center py-3 text-xs text-muted-foreground italic">
                        Belum ada item tenaga kerja.{' '}
                        <button
                          type="button"
                          onClick={() => addItem('LABOR')}
                          className="text-primary underline hover:text-primary/80 font-medium ml-1"
                        >
                          + Tambah Tenaga Kerja
                        </button>
                      </TableCell>
                    </TableRow>
                  ) : (
                    items
                      .filter((it) => it.category === 'LABOR')
                      .map((item, idx) => (
                        <ComponentRow
                          key={item.id}
                          item={item}
                          no={idx + 1}
                          masterPrices={masterPrices}
                          onUpdate={updateItem}
                          onRemove={removeItem}
                        />
                      ))
                  )}
                  <TableRow className="border-t bg-muted/20 font-medium text-xs">
                    <TableCell colSpan={6} className="text-right text-muted-foreground font-semibold">
                      Jumlah Tenaga Kerja (A):
                    </TableCell>
                    <TableCell className="text-right font-bold text-foreground">
                      {formatRupiah(calc.labor)}
                    </TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* Category B: MATERIALS */}
                  <TableRow className="bg-muted/70 font-semibold border-t-2">
                    <TableCell className="text-center font-bold text-primary">B</TableCell>
                    <TableCell colSpan={6} className="text-foreground">
                      <div className="flex items-center gap-2">
                        <Package className="size-4 text-primary" />
                        <span>BAHAN (MATERIALS)</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-center p-1">
                      <Button
                        type="button"
                        size="sm"
                        variant="ghost"
                        onClick={() => addItem('MATERIALS')}
                        title="Tambah Bahan"
                        className="h-8 w-8 p-0 text-primary hover:text-primary hover:bg-primary/10"
                      >
                        <Plus className="size-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                  {items.filter((it) => it.category === 'MATERIALS').length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={8} className="text-center py-3 text-xs text-muted-foreground italic">
                        Belum ada item bahan.{' '}
                        <button
                          type="button"
                          onClick={() => addItem('MATERIALS')}
                          className="text-primary underline hover:text-primary/80 font-medium ml-1"
                        >
                          + Tambah Bahan
                        </button>
                      </TableCell>
                    </TableRow>
                  ) : (
                    items
                      .filter((it) => it.category === 'MATERIALS')
                      .map((item, idx) => (
                        <ComponentRow
                          key={item.id}
                          item={item}
                          no={idx + 1}
                          masterPrices={masterPrices}
                          onUpdate={updateItem}
                          onRemove={removeItem}
                        />
                      ))
                  )}
                  <TableRow className="border-t bg-muted/20 font-medium text-xs">
                    <TableCell colSpan={6} className="text-right text-muted-foreground font-semibold">
                      Jumlah Bahan (B):
                    </TableCell>
                    <TableCell className="text-right font-bold text-foreground">
                      {formatRupiah(calc.material)}
                    </TableCell>
                    <TableCell></TableCell>
                  </TableRow>

                  {/* Category C: EQUIPMENT */}
                  <TableRow className="bg-muted/70 font-semibold border-t-2">
                    <TableCell className="text-center font-bold text-primary">C</TableCell>
                    <TableCell colSpan={6} className="text-foreground">
                      <div className="flex items-center gap-2">
                        <Wrench className="size-4 text-primary" />
                        <span>PERALATAN (EQUIPMENT)</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-center p-1">
                      <Button
                        type="button"
                        size="sm"
                        variant="ghost"
                        onClick={() => addItem('EQUIPMENT')}
                        title="Tambah Peralatan"
                        className="h-8 w-8 p-0 text-primary hover:text-primary hover:bg-primary/10"
                      >
                        <Plus className="size-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                  {items.filter((it) => it.category === 'EQUIPMENT').length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={8} className="text-center py-3 text-xs text-muted-foreground italic">
                        Belum ada item peralatan.{' '}
                        <button
                          type="button"
                          onClick={() => addItem('EQUIPMENT')}
                          className="text-primary underline hover:text-primary/80 font-medium ml-1"
                        >
                          + Tambah Peralatan
                        </button>
                      </TableCell>
                    </TableRow>
                  ) : (
                    items
                      .filter((it) => it.category === 'EQUIPMENT')
                      .map((item, idx) => (
                        <ComponentRow
                          key={item.id}
                          item={item}
                          no={idx + 1}
                          masterPrices={masterPrices}
                          onUpdate={updateItem}
                          onRemove={removeItem}
                        />
                      ))
                  )}
                  <TableRow className="border-t bg-muted/20 font-medium text-xs">
                    <TableCell colSpan={6} className="text-right text-muted-foreground font-semibold">
                      Jumlah Peralatan (C):
                    </TableCell>
                    <TableCell className="text-right font-bold text-foreground">
                      {formatRupiah(calc.equipment)}
                    </TableCell>
                    <TableCell></TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>

        {/* Section 3: Summary & Overhead (D, E, F) */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card className="md:col-span-1 border-muted">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-semibold text-muted-foreground">
                D. Total Biaya Langsung
              </CardTitle>
              <CardDescription className="text-xs">Komponen Tenaga (A) + Bahan (B) + Alat (C)</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold font-mono text-foreground">
                {formatRupiah(calc.d)}
              </div>
            </CardContent>
          </Card>

          <Card className="md:col-span-1 border-muted">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-semibold text-muted-foreground">
                  E. Overhead & Keuntungan
                </CardTitle>
                <div className="flex items-center gap-1.5">
                  <Input
                    type="number"
                    step="0.01"
                    min="0"
                    value={overheadPct}
                    onChange={(e) => setOverheadPct(e.target.value)}
                    className="h-7 w-16 text-center font-bold text-xs bg-card"
                  />
                  <span className="text-xs font-semibold">%</span>
                </div>
              </div>
              <CardDescription className="text-xs">{overheadPct}% dari Biaya Langsung (D)</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold font-mono text-foreground">
                {formatRupiah(calc.e)}
              </div>
            </CardContent>
          </Card>

          <Card className="md:col-span-1 border-primary/40 bg-primary/5">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-bold text-primary">
                  F. Harga Satuan Pekerjaan (AHSP)
                </CardTitle>
                <Badge variant="default" className="text-[10px] uppercase font-semibold">
                  D + E
                </Badge>
              </div>
              <CardDescription className="text-xs">
                Harga per 1 {unit ? <span className="font-semibold text-foreground">{unit}</span> : 'satuan'}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-extrabold font-mono text-primary">
                {formatRupiah(calc.f)}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Bottom Actions */}
        <div className="flex items-center justify-end gap-3 pt-4 border-t">
          <Button
            type="button"
            variant="outline"
            size="lg"
            onClick={() => navigate({ to: '/ahsp' })}
            disabled={isPending}
          >
            Batal
          </Button>
          <Button type="submit" size="lg" disabled={isPending} className="min-w-36">
            <Save className="size-4 mr-2" />
            {isPending ? 'Menyimpan...' : isEdit ? 'Simpan Perubahan' : 'Simpan AHSP'}
          </Button>
        </div>
      </form>
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
    <TableRow className="text-xs hover:bg-muted/30">
      <TableCell className="text-center text-muted-foreground font-mono">{no}</TableCell>
      <TableCell className="p-2">
        <div className="flex flex-col sm:flex-row gap-2">
          {filteredPrices.length > 0 && (
            <NativeSelect
              value={item.unit_price_id || ''}
              onChange={(e) => onUpdate(item.id, 'unit_price_id', e.target.value)}
              className="h-8 max-w-[180px] shrink-0 text-xs"
            >
              <NativeSelectOption value="">-- Pilih Master --</NativeSelectOption>
              {filteredPrices.map((m) => (
                <NativeSelectOption key={m.id} value={m.id}>
                  {m.code} - {m.work_item}
                </NativeSelectOption>
              ))}
            </NativeSelect>
          )}
          <Input
            value={item.description}
            onChange={(e) => onUpdate(item.id, 'description', e.target.value)}
            placeholder="Deskripsi uraian komponen..."
            className="h-8 text-xs bg-card flex-1"
            required
          />
        </div>
      </TableCell>
      <TableCell className="p-2">
        <Input
          value={item.code || ''}
          onChange={(e) => onUpdate(item.id, 'code', e.target.value)}
          placeholder="Kode..."
          className="h-8 text-xs font-mono bg-card"
        />
      </TableCell>
      <TableCell className="p-2 text-center">
        <Input
          value={item.unit}
          onChange={(e) => onUpdate(item.id, 'unit', e.target.value)}
          placeholder="OH/m3"
          className="h-8 text-xs text-center bg-card"
          required
        />
      </TableCell>
      <TableCell className="p-2 text-right">
        <Input
          type="number"
          step="any"
          value={item.coefficient}
          onChange={(e) => onUpdate(item.id, 'coefficient', e.target.value)}
          className="h-8 text-xs text-right font-mono bg-card"
          required
        />
      </TableCell>
      <TableCell className="p-2 text-right">
        <Input
          type="number"
          value={item.unit_price}
          onChange={(e) => onUpdate(item.id, 'unit_price', e.target.value)}
          className="h-8 text-xs text-right font-mono bg-card"
          required
        />
      </TableCell>
      <TableCell className="text-right font-mono font-semibold text-foreground p-2">
        {formatRupiah(cost)}
      </TableCell>
      <TableCell className="p-2 text-center">
        <Button
          type="button"
          size="icon"
          variant="ghost"
          className="size-8 text-muted-foreground hover:text-destructive"
          onClick={() => onRemove(item.id)}
          title="Hapus baris"
        >
          <Trash2 className="size-4" />
        </Button>
      </TableCell>
    </TableRow>
  )
}
