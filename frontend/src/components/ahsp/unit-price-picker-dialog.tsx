import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { listUnitPricesOptions } from '@/api/@tanstack/react-query.gen'
import type { UnitPriceDto } from '@/api/types.gen'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { formatRupiah } from '@/lib/utils'
import { Search, Plus } from 'lucide-react'

type Category = 'LABOR' | 'MATERIALS' | 'EQUIPMENT'

const categoryMeta: Record<Category, { title: string; label: string }> = {
  LABOR: { title: 'Pilih Tenaga Kerja (Labor)', label: 'Tenaga Kerja' },
  MATERIALS: { title: 'Pilih Bahan (Materials)', label: 'Bahan Material' },
  EQUIPMENT: { title: 'Pilih Peralatan (Equipment)', label: 'Peralatan' },
}

type Props = {
  open: boolean
  onOpenChange: (open: boolean) => void
  category: Category
  existingUnitPriceIds?: string[]
  onSelectItems: (items: UnitPriceDto[]) => void
}

export function UnitPricePickerDialog({
  open,
  onOpenChange,
  category,
  existingUnitPriceIds = [],
  onSelectItems,
}: Props) {
  const [search, setSearch] = useState('')
  const [selectedIds, setSelectedIds] = useState<string[]>([])

  const { data: allUnitPrices = [], isLoading } = useQuery({
    ...listUnitPricesOptions({
      query: { category, q: search || undefined },
    }),
    enabled: open,
  })

  const meta = categoryMeta[category] || { title: 'Pilih Unit Price', label: 'Item' }

  const toggleSelect = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((it) => it !== id) : [...prev, id]
    )
  }

  const toggleSelectAll = () => {
    if (selectedIds.length === allUnitPrices.length) {
      setSelectedIds([])
    } else {
      setSelectedIds(allUnitPrices.map((p) => p.id))
    }
  }

  const handleConfirmMulti = () => {
    const selected = allUnitPrices.filter((p) => selectedIds.includes(p.id))
    if (selected.length > 0) {
      onSelectItems(selected)
      setSelectedIds([])
      onOpenChange(false)
    }
  }

  const handleQuickSelect = (item: UnitPriceDto) => {
    onSelectItems([item])
    setSelectedIds([])
    onOpenChange(false)
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(nextOpen) => {
        if (!nextOpen) {
          setSelectedIds([])
          setSearch('')
        }
        onOpenChange(nextOpen)
      }}
    >
      <DialogContent className="max-w-3xl max-h-[85vh] flex flex-col p-6 overflow-hidden">
        <DialogHeader>
          <div className="flex items-center justify-between">
            <DialogTitle>{meta.title}</DialogTitle>
            <Badge variant="outline" className="font-mono text-xs">
              {meta.label}
            </Badge>
          </div>
          <DialogDescription>
            Pilih item dari Master Harga Satuan (HSD) untuk ditautkan ke rincian AHSP.
          </DialogDescription>
        </DialogHeader>

        {/* Search Input */}
        <div className="relative my-2">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
          <Input
            placeholder={`Cari kode atau nama ${meta.label.toLowerCase()}...`}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9 bg-card"
            autoFocus
          />
        </div>

        {/* Table List */}
        <div className="flex-1 overflow-y-auto rounded-md border bg-card min-h-[240px]">
          <Table>
            <TableHeader className="sticky top-0 bg-muted/80 backdrop-blur-xs z-10">
              <TableRow>
                <TableHead className="w-10 text-center">
                  <Checkbox
                    checked={
                      allUnitPrices.length > 0 && selectedIds.length === allUnitPrices.length
                    }
                    onCheckedChange={toggleSelectAll}
                    aria-label="Pilih semua"
                  />
                </TableHead>
                <TableHead className="w-28">Kode</TableHead>
                <TableHead>Uraian Pekerjaan / Item</TableHead>
                <TableHead className="w-20 text-center">Satuan</TableHead>
                <TableHead className="w-36 text-right">Harga Satuan</TableHead>
                <TableHead className="w-20 text-center">Aksi</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-40 text-center text-muted-foreground">
                    Memuat data master harga...
                  </TableCell>
                </TableRow>
              ) : allUnitPrices.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-40 text-center">
                    <Empty>
                      <EmptyHeader>
                        <EmptyTitle>Tidak ada item ditemukan</EmptyTitle>
                        <EmptyDescription>
                          {search
                            ? `Tidak ada ${meta.label.toLowerCase()} yang cocok dengan "${search}".`
                            : `Belum ada master ${meta.label.toLowerCase()} terdaftar di Unit Prices.`}
                        </EmptyDescription>
                      </EmptyHeader>
                    </Empty>
                  </TableCell>
                </TableRow>
              ) : (
                allUnitPrices.map((p) => {
                  const isSelected = selectedIds.includes(p.id)
                  const isAlreadyAdded = existingUnitPriceIds.includes(p.id)

                  return (
                    <TableRow
                      key={p.id}
                      className={`cursor-pointer transition-colors ${
                        isSelected ? 'bg-primary/5' : 'hover:bg-muted/40'
                      }`}
                      onClick={() => toggleSelect(p.id)}
                    >
                      <TableCell
                        className="text-center"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <Checkbox
                          checked={isSelected}
                          onCheckedChange={() => toggleSelect(p.id)}
                          aria-label={`Pilih ${p.work_item}`}
                        />
                      </TableCell>
                      <TableCell className="font-mono text-xs font-semibold">
                        {p.code}
                      </TableCell>
                      <TableCell className="font-medium text-foreground">
                        <div className="flex items-center gap-2">
                          <span>{p.work_item}</span>
                          {isAlreadyAdded && (
                            <Badge variant="secondary" className="text-[10px] py-0 px-1.5 font-normal">
                              Sudah ditambahkan
                            </Badge>
                          )}
                        </div>
                      </TableCell>
                      <TableCell className="text-center font-mono text-xs">
                        {p.unit}
                      </TableCell>
                      <TableCell className="text-right font-mono font-semibold text-foreground">
                        {formatRupiah(p.unit_price)}
                      </TableCell>
                      <TableCell
                        className="text-center p-1"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <Button
                          type="button"
                          size="sm"
                          variant={isSelected ? 'default' : 'outline'}
                          className="h-7 px-2 text-xs"
                          onClick={() => handleQuickSelect(p)}
                        >
                          Pilih
                        </Button>
                      </TableCell>
                    </TableRow>
                  )
                })
              )}
            </TableBody>
          </Table>
        </div>

        {/* Footer */}
        <DialogFooter className="mt-4 pt-3 border-t flex items-center justify-between sm:justify-between w-full">
          <div className="text-xs text-muted-foreground">
            {selectedIds.length > 0 ? (
              <span className="font-medium text-foreground">
                {selectedIds.length} item dipilih
              </span>
            ) : (
              <span>Pilih satu atau beberapa item untuk ditambahkan.</span>
            )}
          </div>
          <div className="flex items-center gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                setSelectedIds([])
                onOpenChange(false)
              }}
            >
              Batal
            </Button>
            <Button
              type="button"
              disabled={selectedIds.length === 0}
              onClick={handleConfirmMulti}
            >
              <Plus className="size-4 mr-1.5" />
              Tambahkan ({selectedIds.length})
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
