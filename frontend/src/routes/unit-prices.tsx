import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listUnitPricesOptions, createUnitPriceMutation, deleteUnitPriceMutation } from '@/api/@tanstack/react-query.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { formatRupiah } from '@/lib/utils'
import { useAuth } from '@/lib/auth'

export const Route = createFileRoute('/unit-prices')({ component: UnitPricesPage })

function UnitPricesPage() {
  const queryClient = useQueryClient()
  const { isAdmin, user, openLoginModal } = useAuth()

  const [search, setSearch] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<string>('ALL')
  const [isAddOpen, setIsAddOpen] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')

  // Form
  const [code, setCode] = useState('')
  const [category, setCategory] = useState<'LABOR' | 'MATERIALS' | 'EQUIPMENT' | 'OTHER'>('MATERIALS')
  const [workItem, setWorkItem] = useState('')
  const [unit, setUnit] = useState('')
  const [unitPrice, setUnitPrice] = useState('')
  const [information, setInformation] = useState('')

  const { data: items = [], isLoading } = useQuery({
    ...listUnitPricesOptions({
      query: {
        q: search || undefined,
        category: selectedCategory === 'ALL' ? undefined : (selectedCategory as any),
      },
    }),
  })

  const createMut = useMutation(createUnitPriceMutation())
  const deleteMut = useMutation(deleteUnitPriceMutation())

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal()
      return
    }
    setCode('')
    setCategory('MATERIALS')
    setWorkItem('')
    setUnit('')
    setUnitPrice('')
    setInformation('')
    setErrorMessage('')
    setIsAddOpen(true)
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage('')
    if (!code.trim() || !workItem.trim() || !unit.trim()) {
      setErrorMessage('Code, item, dan unit wajib diisi.')
      return
    }

    try {
      await createMut.mutateAsync({
        body: {
          code: code.trim(),
          category,
          work_item: workItem.trim(),
          unit: unit.trim(),
          unit_price: unitPrice.replace(/\D/g, '') || '0',
          information: information.trim(),
        },
      })
      await queryClient.invalidateQueries()
      setIsAddOpen(false)
    } catch (err: any) {
      setErrorMessage(err?.message || 'Gagal menyimpan harga satuan.')
    }
  }

  const handleDelete = async (id: string, code: string) => {
    if (!confirm(`Hapus harga satuan ${code}?`)) return
    try {
      await deleteMut.mutateAsync({ path: { id } })
      await queryClient.invalidateQueries()
    } catch (err: any) {
      alert(err?.message || 'Gagal menghapus: Item mungkin sedang digunakan pada AHSP.')
    }
  }

  return (
    <div>
      <PageHeader title="UNIT PRICE" onSearch={setSearch} onAdd={handleOpenAdd} />

      {/* Category Filter Pills */}
      <div className="mb-4 flex gap-2">
        {['ALL', 'LABOR', 'MATERIALS', 'EQUIPMENT', 'OTHER'].map((cat) => (
          <button
            key={cat}
            onClick={() => setSelectedCategory(cat)}
            className={`rounded-md px-3 py-1.5 text-xs font-semibold transition-colors ${
              selectedCategory === cat
                ? 'bg-primary text-primary-foreground'
                : 'bg-muted text-muted-foreground hover:bg-muted/80'
            }`}
          >
            {cat}
          </button>
        ))}
      </div>

      <div className="overflow-x-auto border bg-white shadow-sm">
        <Table className="text-xs">
          <TableHeader className="bg-white">
            <TableRow className="border-b hover:bg-transparent">
              <TableHead className="w-12 border-r text-center font-bold">NO</TableHead>
              <TableHead className="w-28 border-r text-center font-bold">CODE</TableHead>
              <TableHead className="w-32 border-r text-center font-bold">CATEGORY</TableHead>
              <TableHead className="min-w-48 border-r text-center font-bold">ITEM</TableHead>
              <TableHead className="w-20 border-r text-center font-bold">UNIT</TableHead>
              <TableHead className="w-36 border-r text-center font-bold">UNIT PRICE</TableHead>
              <TableHead className="min-w-40 border-r text-center font-bold">INFORMATION</TableHead>
              {isAdmin && <TableHead className="w-16 text-center font-bold">AKSI</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 8 : 7} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  Memuat data...
                </TableCell>
              </TableRow>
            ) : items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 8 : 7} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  {search ? `Tidak ada harga satuan yang cocok dengan "${search}".` : 'Belum ada data harga satuan.'}
                </TableCell>
              </TableRow>
            ) : (
              items.map((item, idx) => (
                <TableRow key={item.id} className="border-b bg-muted/40 hover:bg-muted/80">
                  <TableCell className="border-r text-center font-medium">{idx + 1}</TableCell>
                  <TableCell className="border-r font-mono font-semibold">{item.code}</TableCell>
                  <TableCell className="border-r text-center font-medium">{item.category}</TableCell>
                  <TableCell className="border-r font-medium text-foreground">{item.work_item}</TableCell>
                  <TableCell className="border-r text-center">{item.unit}</TableCell>
                  <TableCell className="border-r text-right font-semibold">{formatRupiah(item.unit_price)}</TableCell>
                  <TableCell className="border-r">{item.information || '-'}</TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-7 text-xs text-red-600 hover:bg-red-50 hover:text-red-700"
                        onClick={() => handleDelete(item.id, item.code)}
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

      {/* Modal Add Unit Price (Mockup Hal. 5) */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-md rounded-2xl">
          <DialogHeader>
            <DialogTitle className="text-lg font-bold tracking-wide">TAMBAH UNIT PRICE</DialogTitle>
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
                placeholder="misal MAT-001, LAB-001"
                className="mt-1"
                required
              />
            </div>

            <div>
              <Label className="text-xs">CATEGORY</Label>
              <Select value={category} onValueChange={(val: any) => setCategory(val)}>
                <SelectTrigger className="mt-1">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="LABOR">LABOR (Tenaga Kerja)</SelectItem>
                  <SelectItem value="MATERIALS">MATERIALS (Bahan/Material)</SelectItem>
                  <SelectItem value="EQUIPMENT">EQUIPMENT (Peralatan)</SelectItem>
                  <SelectItem value="OTHER">OTHER (Lain-lain)</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div>
              <Label className="text-xs">WORK ITEM</Label>
              <Input
                value={workItem}
                onChange={(e) => setWorkItem(e.target.value)}
                placeholder="misal Pasir Beton, Mandor, dsb."
                className="mt-1"
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label className="text-xs">UNIT</Label>
                <Input
                  value={unit}
                  onChange={(e) => setUnit(e.target.value)}
                  placeholder="m3, kg, OH, sak..."
                  className="mt-1"
                  required
                />
              </div>
              <div>
                <Label className="text-xs">UNIT PRICE</Label>
                <Input
                  value={unitPrice}
                  onChange={(e) => setUnitPrice(e.target.value)}
                  placeholder="Nominal Rupiah..."
                  className="mt-1"
                  required
                />
              </div>
            </div>

            <div>
              <Label className="text-xs">INFORMATION</Label>
              <Input
                value={information}
                onChange={(e) => setInformation(e.target.value)}
                placeholder="Keterangan..."
                className="mt-1"
              />
            </div>

            <div className="flex justify-end gap-2 pt-4">
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
