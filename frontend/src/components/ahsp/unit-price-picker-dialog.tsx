import { useState, useMemo } from 'react'
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
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { InputGroup, InputGroupAddon, InputGroupInput } from '@/components/ui/input-group'
import { formatRupiah } from '@/lib/utils'
import { Search } from 'lucide-react'

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
      query: { category, q: search.trim() || undefined },
    }),
    enabled: open,
  })

  const meta = categoryMeta[category] || {
    title: 'Pilih Unit Price',
    label: 'Item',
  }

  const availablePrices = useMemo(
    () => allUnitPrices.filter((p) => !existingUnitPriceIds.includes(p.id)),
    [allUnitPrices, existingUnitPriceIds]
  )

  const isAllSelected =
    availablePrices.length > 0 &&
    availablePrices.every((p) => selectedIds.includes(p.id))

  const toggleSelect = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((it) => it !== id) : [...prev, id]
    )
  }

  const toggleSelectAll = () => {
    if (isAllSelected) {
      setSelectedIds([])
    } else {
      setSelectedIds(availablePrices.map((p) => p.id))
    }
  }

  const handleConfirm = () => {
    const selected = allUnitPrices.filter((p) => selectedIds.includes(p.id))
    if (selected.length > 0) {
      onSelectItems(selected)
      setSelectedIds([])
      onOpenChange(false)
    }
  }

  const handleOpenChange = (nextOpen: boolean) => {
    if (!nextOpen) {
      setSelectedIds([])
      setSearch('')
    }
    onOpenChange(nextOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-3xl">
        <DialogHeader>
          <DialogTitle>{meta.title}</DialogTitle>
          <DialogDescription>
            Pilih item dari Master Harga Satuan (HSD) untuk dimasukkan ke rincian AHSP.
          </DialogDescription>
        </DialogHeader>

        <InputGroup>
          <InputGroupAddon>
            <Search data-icon="inline-start" />
          </InputGroupAddon>
          <InputGroupInput
            placeholder={`Cari kode atau nama ${meta.label.toLowerCase()}...`}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </InputGroup>

        <div className="rounded-md border max-h-80 overflow-y-auto">
          <Table>
            <TableHeader className="sticky top-0 bg-background">
              <TableRow>
                <TableHead className="w-12 text-center">
                  <Checkbox
                    checked={isAllSelected}
                    onCheckedChange={toggleSelectAll}
                    disabled={availablePrices.length === 0}
                    aria-label="Pilih semua"
                  />
                </TableHead>
                <TableHead className="w-32">Kode</TableHead>
                <TableHead>Uraian Pekerjaan / Item</TableHead>
                <TableHead className="w-20 text-center">Satuan</TableHead>
                <TableHead className="w-36 text-right">Harga Satuan</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-32 text-center text-muted-foreground">
                    Memuat data master harga...
                  </TableCell>
                </TableRow>
              ) : allUnitPrices.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-32 text-center">
                    <Empty>
                      <EmptyHeader>
                        <EmptyTitle>Tidak ada item</EmptyTitle>
                        <EmptyDescription>
                          {search
                            ? `Tidak ada ${meta.label.toLowerCase()} yang cocok dengan "${search}".`
                            : `Belum ada data master ${meta.label.toLowerCase()}.`}
                        </EmptyDescription>
                      </EmptyHeader>
                    </Empty>
                  </TableCell>
                </TableRow>
              ) : (
                allUnitPrices.map((p) => {
                  const isAlreadyInAhsp = existingUnitPriceIds.includes(p.id)
                  const isChecked = selectedIds.includes(p.id)

                  return (
                    <TableRow
                      key={p.id}
                      className={
                        isAlreadyInAhsp
                          ? 'opacity-50 cursor-not-allowed'
                          : 'cursor-pointer'
                      }
                      onClick={() => {
                        if (!isAlreadyInAhsp) toggleSelect(p.id)
                      }}
                    >
                      <TableCell
                        className="text-center"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <Checkbox
                          checked={isChecked || isAlreadyInAhsp}
                          disabled={isAlreadyInAhsp}
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
                          {isAlreadyInAhsp && (
                            <Badge variant="secondary" className="text-[10px]">
                              Sudah ada
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
                    </TableRow>
                  )
                })
              )}
            </TableBody>
          </Table>
        </div>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => handleOpenChange(false)}
          >
            Batal
          </Button>
          <Button
            type="button"
            disabled={selectedIds.length === 0}
            onClick={handleConfirm}
          >
            Tambahkan ({selectedIds.length})
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
