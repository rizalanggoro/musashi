import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listUnitPricesOptions,
  createUnitPriceMutation,
  updateUnitPriceMutation,
  deleteUnitPriceMutation,
} from '@/api/@tanstack/react-query.gen'
import type { UnitPriceDto } from '@/api/types.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FieldGroup, Field, FieldLabel } from '@/components/ui/field'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ToggleGroup, ToggleGroupItem } from '@/components/ui/toggle-group'
import { Badge } from '@/components/ui/badge'
import { formatRupiah } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import { MoreHorizontal, Pencil, Trash2 } from 'lucide-react'

export const Route = createFileRoute('/unit-prices/')({ component: UnitPricesPage })

const categories = ['ALL', 'LABOR', 'MATERIALS', 'EQUIPMENT', 'OTHER'] as const

function UnitPricesPage() {
  const queryClient = useQueryClient()
  const { isAdmin, user, openLoginModal } = useAuth()

  const [search, setSearch] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<string>('ALL')
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingItem, setEditingItem] = useState<UnitPriceDto | null>(null)
  const [errorMessage, setErrorMessage] = useState('')

  // Form
  const [code, setCode] = useState('')
  const [category, setCategory] = useState<'LABOR' | 'MATERIALS' | 'EQUIPMENT' | 'OTHER'>('MATERIALS')
  const [workItem, setWorkItem] = useState('')
  const [unit, setUnit] = useState('')
  const [unitPrice, setUnitPrice] = useState('')
  const [note, setNote] = useState('')

  const { data: items = [], isLoading } = useQuery({
    ...listUnitPricesOptions({
      query: {
        q: search || undefined,
        category: selectedCategory === 'ALL' ? undefined : (selectedCategory as any),
      },
    }),
  })

  const createMut = useMutation(createUnitPriceMutation())
  const updateMut = useMutation(updateUnitPriceMutation())
  const deleteMut = useMutation(deleteUnitPriceMutation())

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal()
      return
    }
    setEditingItem(null)
    setCode('')
    setCategory('MATERIALS')
    setWorkItem('')
    setUnit('')
    setUnitPrice('')
    setNote('')
    setErrorMessage('')
    setIsDialogOpen(true)
  }

  const handleOpenEdit = (item: UnitPriceDto) => {
    if (!user) {
      openLoginModal()
      return
    }
    setEditingItem(item)
    setCode(item.code)
    setCategory((item.category as any) || 'MATERIALS')
    setWorkItem(item.work_item)
    setUnit(item.unit)
    setUnitPrice(item.unit_price)
    setNote(item.note || '')
    setErrorMessage('')
    setIsDialogOpen(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorMessage('')
    if (!code.trim() || !workItem.trim() || !unit.trim()) {
      setErrorMessage('Code, item, dan unit wajib diisi.')
      return
    }

    try {
      if (editingItem) {
        await updateMut.mutateAsync({
          path: { id: editingItem.id },
          body: {
            code: code.trim(),
            category,
            work_item: workItem.trim(),
            unit: unit.trim(),
            unit_price: unitPrice.replace(/\D/g, '') || '0',
            note: note.trim(),
          },
        })
      } else {
        await createMut.mutateAsync({
          body: {
            code: code.trim(),
            category,
            work_item: workItem.trim(),
            unit: unit.trim(),
            unit_price: unitPrice.replace(/\D/g, '') || '0',
            note: note.trim(),
          },
        })
      }
      await queryClient.invalidateQueries()
      setIsDialogOpen(false)
      setEditingItem(null)
    } catch (err: any) {
      setErrorMessage(
        err?.message || (editingItem ? 'Gagal mengubah harga satuan.' : 'Gagal menyimpan harga satuan.')
      )
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
    <div className="flex flex-col gap-4">
      <PageHeader title="Unit Prices" onSearch={setSearch} onAdd={handleOpenAdd} addLabel="Unit Price" />

      {/* Category Filter using ToggleGroup */}
      <ToggleGroup
        value={[selectedCategory]}
        onValueChange={(val) => {
          if (val && val.length > 0) {
            setSelectedCategory(val[0])
          }
        }}
        spacing={2}
        className="flex flex-wrap"
      >
        {categories.map((cat) => (
          <ToggleGroupItem key={cat} value={cat} className="h-8 px-3 text-xs">
            {cat}
          </ToggleGroupItem>
        ))}
      </ToggleGroup>

      <div className="rounded-md border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12 text-center">No</TableHead>
              <TableHead className="w-28">Code</TableHead>
              <TableHead className="w-28">Category</TableHead>
              <TableHead>Work Item</TableHead>
              <TableHead className="w-20 text-center">Unit</TableHead>
              <TableHead className="w-36 text-right">Unit Price</TableHead>
              <TableHead>Note</TableHead>
              {isAdmin && <TableHead className="w-20 text-center">Action</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 8 : 7} className="h-32 text-center text-muted-foreground">
                  Loading data...
                </TableCell>
              </TableRow>
            ) : items.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 8 : 7} className="h-48 text-center">
                  <Empty>
                    <EmptyHeader>
                      <EmptyTitle>Tidak ada harga satuan</EmptyTitle>
                      <EmptyDescription>
                        {search ? `Tidak ada item yang cocok dengan "${search}".` : 'Mulai tambahkan daftar harga satuan baru.'}
                      </EmptyDescription>
                    </EmptyHeader>
                  </Empty>
                </TableCell>
              </TableRow>
            ) : (
              items.map((item, idx) => (
                <TableRow key={item.id}>
                  <TableCell className="text-center text-muted-foreground">{idx + 1}</TableCell>
                  <TableCell className="font-mono text-xs font-semibold">{item.code}</TableCell>
                  <TableCell>
                    <Badge variant="outline" className="text-[10px] font-normal uppercase">
                      {item.category}
                    </Badge>
                  </TableCell>
                  <TableCell className="font-medium text-foreground">{item.work_item}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{item.unit}</TableCell>
                  <TableCell className="text-right font-medium">{formatRupiah(item.unit_price)}</TableCell>
                  <TableCell className="text-muted-foreground">{item.note || '-'}</TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <DropdownMenu>
                        <DropdownMenuTrigger
                          render={
                            <Button
                              variant="ghost"
                              size="icon"
                              className="size-8 text-muted-foreground hover:text-foreground"
                            />
                          }
                        >
                          <MoreHorizontal className="size-4" />
                          <span className="sr-only">Menu</span>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-32">
                          <DropdownMenuItem onClick={() => handleOpenEdit(item)}>
                            <Pencil className="size-4 mr-2" />
                            <span>Ubah</span>
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            variant="destructive"
                            onClick={() => handleDelete(item.id, item.code)}
                          >
                            <Trash2 className="size-4 mr-2" />
                            <span>Hapus</span>
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  )}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Modal Add / Edit Unit Price */}
      <Dialog
        open={isDialogOpen}
        onOpenChange={(open) => {
          setIsDialogOpen(open)
          if (!open) setEditingItem(null)
        }}
      >
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{editingItem ? 'Edit Unit Price' : 'Add Unit Price'}</DialogTitle>
            <DialogDescription>
              {editingItem
                ? 'Perbarui data harga satuan.'
                : 'Tambahkan harga satuan baru ke master database.'}
            </DialogDescription>
          </DialogHeader>

          {errorMessage && (
            <Alert variant="destructive">
              <AlertDescription>{errorMessage}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <FieldGroup>
              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="up-code">Code</FieldLabel>
                  <Input
                    id="up-code"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    placeholder="MTR-001"
                    required
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="up-category">Category</FieldLabel>
                  <Select
                    value={category}
                    onValueChange={(val) => setCategory(val as any)}
                  >
                    <SelectTrigger id="up-category">
                      <SelectValue placeholder="Pilih Kategori" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="LABOR">LABOR</SelectItem>
                        <SelectItem value="MATERIALS">MATERIALS</SelectItem>
                        <SelectItem value="EQUIPMENT">EQUIPMENT</SelectItem>
                        <SelectItem value="OTHER">OTHER</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="up-workitem">Work Item</FieldLabel>
                <Input
                  id="up-workitem"
                  value={workItem}
                  onChange={(e) => setWorkItem(e.target.value)}
                  placeholder="Nama item pekerjaan / material..."
                  required
                />
              </Field>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="up-unit">Unit</FieldLabel>
                  <Input
                    id="up-unit"
                    value={unit}
                    onChange={(e) => setUnit(e.target.value)}
                    placeholder="m3, kg, OH, dll."
                    required
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="up-price">Unit Price (Rp)</FieldLabel>
                  <Input
                    id="up-price"
                    value={unitPrice}
                    onChange={(e) => setUnitPrice(e.target.value)}
                    placeholder="150000"
                    required
                  />
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="up-note">Note</FieldLabel>
                <Input
                  id="up-note"
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                  placeholder="Catatan spesifikasi..."
                />
              </Field>
            </FieldGroup>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setIsDialogOpen(false)
                  setEditingItem(null)
                }}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={editingItem ? updateMut.isPending : createMut.isPending}
              >
                {editingItem
                  ? updateMut.isPending
                    ? 'Saving...'
                    : 'Save'
                  : createMut.isPending
                    ? 'Saving...'
                    : 'Save'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
