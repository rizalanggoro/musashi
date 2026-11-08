import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  listProjectsOptions,
  listVolumesOptions,
  deleteVolumeMutation,
} from '@/api/@tanstack/react-query.gen'
import { PageHeader } from '@/components/page-header'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button, buttonVariants } from '@/components/ui/button'
import { Empty, EmptyHeader, EmptyTitle, EmptyDescription } from '@/components/ui/empty'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { formatDecimal } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import { FileSpreadsheet } from 'lucide-react'

export const Route = createFileRoute('/volume/')({
  component: VolumePage,
  validateSearch: (search: Record<string, unknown>): { projectId?: string } => ({
    projectId: (search.projectId as string) || undefined,
  }),
})

export type FormulaType =
  | 'box'          // Balok / Kotak: P × L × T
  | 'cylinder'     // Silinder / Lubang Bor: ¼ × π × D² × T
  | 'area_pl'      // Luas Lantai / Plafon: P × L
  | 'area_pt'      // Luas Dinding / Pagar: P × T
  | 'area_lt'      // Luas Samping: L × T
  | 'area_circle'  // Luas Lingkaran: ¼ × π × D²
  | 'linear_p'     // Meter Lari - Panjang: P
  | 'linear_t'     // Meter Lari - Tinggi/Tiang: T
  | 'weight_rebar' // Berat Besi/Baja: P × Berat Jenis
  | 'unit_only'    // Kuantitas Unit / Lot Saja
  | 'custom'       // Planned Volume × Unit

interface FormulaConfig {
  id: FormulaType
  label: string
  category: string
  description: string
  fields: {
    length?: string
    width?: string
    height?: string
    weight: string
  }
}

export const FORMULA_CONFIGS: Record<FormulaType, FormulaConfig> = {
  box: {
    id: 'box',
    label: 'Kotak / Balok (P × L × T)',
    category: 'm3',
    description: 'Pekerjaan kubikal 3 dimensi: galian, pondasi telapak, balok/kolom beton',
    fields: { length: 'Panjang (m)', width: 'Lebar (m)', height: 'Tinggi (m)', weight: 'Unit / Kali (Pengali)' },
  },
  cylinder: {
    id: 'cylinder',
    label: 'Silinder / Lubang Bor (¼ × π × D² × T)',
    category: 'm3',
    description: 'Bore pile, lubang bor, sumuran, tiang bulat: ¼ × π × D² × Kedalaman × Titik',
    fields: { length: 'Diameter D (m)', height: 'Kedalaman / Tinggi T (m)', weight: 'Jumlah Titik / Lubang (Unit)' },
  },
  area_pl: {
    id: 'area_pl',
    label: 'Luas Lantai / Plafon (P × L)',
    category: 'm2',
    description: 'Luas bidang horizontal: lantai keramik, screed, dak, plafon',
    fields: { length: 'Panjang (m)', width: 'Lebar (m)', weight: 'Unit / Kali (Pengali)' },
  },
  area_pt: {
    id: 'area_pt',
    label: 'Luas Dinding / Pagar (P × T)',
    category: 'm2',
    description: 'Luas bidang vertikal: dinding bata/hebel, plesteran, acian, cat, partisi',
    fields: { length: 'Panjang (m)', height: 'Tinggi (m)', weight: 'Jumlah Sisi / Unit (misal 2 utk 2 sisi)' },
  },
  area_lt: {
    id: 'area_lt',
    label: 'Luas Samping (L × T)',
    category: 'm2',
    description: 'Luas penampang vertikal samping: bordes tangga, sisi kolom',
    fields: { width: 'Lebar (m)', height: 'Tinggi (m)', weight: 'Unit / Kali (Pengali)' },
  },
  area_circle: {
    id: 'area_circle',
    label: 'Luas Lingkaran (¼ × π × D²)',
    category: 'm2',
    description: 'Luas penampang lingkaran: tutup manhole, plat bulat',
    fields: { length: 'Diameter D (m)', weight: 'Jumlah Unit (Buah)' },
  },
  linear_p: {
    id: 'linear_p',
    label: "Meter Lari - Panjang (P)",
    category: "m'",
    description: 'Pengukuran panjang linear horizontal: plin lantai, pipa air, kabel tray',
    fields: { length: 'Panjang P (m)', weight: 'Jumlah Jalur / Unit' },
  },
  linear_t: {
    id: 'linear_t',
    label: "Meter Lari - Tinggi/Tiang (T)",
    category: "m'",
    description: 'Pengukuran linear vertikal: tiang listrik, pipa tegak, kolom profil',
    fields: { height: 'Tinggi T (m)', weight: 'Jumlah Titik / Tiang' },
  },
  weight_rebar: {
    id: 'weight_rebar',
    label: 'Berat Besi / Baja (P × Berat Jenis)',
    category: 'kg',
    description: 'Pembesian tulangan / profil baja: Panjang × Berat Jenis kg/m\' × Batang',
    fields: { length: 'Panjang Batang (m)', width: "Berat Jenis (kg/m')", weight: 'Jumlah Batang (Unit)' },
  },
  unit_only: {
    id: 'unit_only',
    label: 'Kuantitas Unit / Lot Saja',
    category: 'lot/unit',
    description: 'Kuantitas langsung dari jumlah unit/lot: sewa alat, saklar, lampu, sanitair',
    fields: { weight: 'Jumlah Kuantitas / Unit (misal 1 lot, 8 bh)' },
  },
  custom: {
    id: 'custom',
    label: 'Manual / Planned Volume',
    category: 'lainnya',
    description: 'Perhitungan langsung dari planned volume × pengali unit',
    fields: { weight: 'Pengali Unit' },
  },
}

export function detectFormula(unitStr: string): FormulaType {
  const u = unitStr.toLowerCase().trim()
  if (['lot', 'ls', 'lump sum', 'unit', 'pcs', 'titik', 'buah', 'set', 'btg', 'bh'].includes(u)) {
    return 'unit_only'
  }
  if (['m2', 'm²'].includes(u)) {
    return 'area_pl'
  }
  if (['m3', 'm³'].includes(u)) {
    return 'box'
  }
  if (["m'", 'm', 'meter', 'ml'].includes(u)) {
    return 'linear_p'
  }
  if (['kg', 'ton'].includes(u)) {
    return 'weight_rebar'
  }
  return 'box'
}

export function computeVolume(
  formula: FormulaType,
  lStr: string,
  wStr: string,
  hStr: string,
  unitWeightStr: string,
  plannedStr: string
) {
  const l = parseFloat(lStr || '0') || 0
  const w = parseFloat(wStr || '0') || 0
  const h = parseFloat(hStr || '0') || 0
  const uw = parseFloat(unitWeightStr || '1') || 1
  const planned = parseFloat(plannedStr || '0') || 0

  let qty = 0
  let expression = ''

  switch (formula) {
    case 'cylinder': {
      // 1/4 * pi * D^2 * T
      const area = 0.25 * Math.PI * (l * l)
      const depth = h > 0 ? h : 1
      qty = area * depth
      expression = `¼ × π × (${formatDecimal(l)} m)² × ${formatDecimal(h)} m × ${formatDecimal(uw)} titik`
      break
    }
    case 'area_circle': {
      qty = 0.25 * Math.PI * (l * l)
      expression = `¼ × π × (${formatDecimal(l)} m)² × ${formatDecimal(uw)} unit`
      break
    }
    case 'area_pl': {
      qty = l * w
      expression = `${formatDecimal(l)} m × ${formatDecimal(w)} m × ${formatDecimal(uw)} unit`
      break
    }
    case 'area_pt': {
      qty = l * h
      expression = `${formatDecimal(l)} m × ${formatDecimal(h)} m × ${formatDecimal(uw)} sisi`
      break
    }
    case 'area_lt': {
      qty = w * h
      expression = `${formatDecimal(w)} m × ${formatDecimal(h)} m × ${formatDecimal(uw)} unit`
      break
    }
    case 'linear_p': {
      qty = l
      expression = `${formatDecimal(l)} m × ${formatDecimal(uw)} jalur`
      break
    }
    case 'linear_t': {
      qty = h
      expression = `${formatDecimal(h)} m × ${formatDecimal(uw)} titik`
      break
    }
    case 'weight_rebar': {
      const bj = w > 0 ? w : 1
      qty = l * bj
      expression = `${formatDecimal(l)} m × ${formatDecimal(w)} kg/m' × ${formatDecimal(uw)} btg`
      break
    }
    case 'unit_only': {
      qty = 1
      expression = `1.0000 × ${formatDecimal(uw)} unit`
      break
    }
    case 'box': {
      qty = l * w * h
      expression = `${formatDecimal(l)} m × ${formatDecimal(w)} m × ${formatDecimal(h)} m × ${formatDecimal(uw)} unit`
      break
    }
    case 'custom':
    default: {
      qty = planned > 0 ? planned : (l > 0 && w > 0 && h > 0 ? l * w * h : (l > 0 && w > 0 ? l * w : l))
      expression = `${formatDecimal(qty)} × ${formatDecimal(uw)} unit`
      break
    }
  }

  const totalQty = qty * uw
  return { qty, totalQty, expression }
}

export function renderDimensions(row: {
  formula?: string
  length?: string
  width?: string
  height?: string
  unit_weight?: string
}) {
  const f = row.formula || ''
  const l = row.length ? formatDecimal(row.length) : ''
  const w = row.width ? formatDecimal(row.width) : ''
  const h = row.height ? formatDecimal(row.height) : ''

  switch (f) {
    case 'cylinder':
      return l || h ? `D: ${l || '-'} m, T: ${h || '-'} m` : '-'
    case 'area_circle':
      return l ? `D: ${l} m` : '-'
    case 'area_pl':
      return l || w ? `P: ${l || '-'} m, L: ${w || '-'} m` : '-'
    case 'area_pt':
      return l || h ? `P: ${l || '-'} m, T: ${h || '-'} m` : '-'
    case 'area_lt':
      return w || h ? `L: ${w || '-'} m, T: ${h || '-'} m` : '-'
    case 'linear_p':
      return l ? `P: ${l} m` : '-'
    case 'linear_t':
      return h ? `T: ${h} m` : '-'
    case 'weight_rebar':
      return l || w ? `P: ${l || '-'} m, BJ: ${w || '-'} kg/m'` : '-'
    case 'unit_only':
      return '-'
    case 'box':
    default:
      if (l && w && h) return `${l} × ${w} × ${h} m`
      if (l && w) return `${l} × ${w} m`
      if (l && h) return `${l} × ${h} m`
      if (l) return `P: ${l} m`
      return '-'
  }
}

function getFormulaBadge(f?: string) {
  switch (f) {
    case 'cylinder':
      return <Badge variant="outline" className="bg-blue-50 text-blue-700 dark:bg-blue-950 dark:text-blue-300 border-blue-200">Silinder (D² × T)</Badge>
    case 'area_circle':
      return <Badge variant="outline" className="bg-sky-50 text-sky-700 dark:bg-sky-950 dark:text-sky-300 border-sky-200">Lingkaran (¼πD²)</Badge>
    case 'area_pl':
      return <Badge variant="outline" className="bg-emerald-50 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300 border-emerald-200">Luas (P × L)</Badge>
    case 'area_pt':
      return <Badge variant="outline" className="bg-teal-50 text-teal-700 dark:bg-teal-950 dark:text-teal-300 border-teal-200">Dinding (P × T)</Badge>
    case 'area_lt':
      return <Badge variant="outline" className="bg-cyan-50 text-cyan-700 dark:bg-cyan-950 dark:text-cyan-300 border-cyan-200">Samping (L × T)</Badge>
    case 'linear_p':
    case 'linear_t':
      return <Badge variant="outline" className="bg-indigo-50 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300 border-indigo-200">Meter Lari (m')</Badge>
    case 'weight_rebar':
      return <Badge variant="outline" className="bg-amber-50 text-amber-700 dark:bg-amber-950 dark:text-amber-300 border-amber-200">Besi (P × BJ)</Badge>
    case 'unit_only':
      return <Badge variant="outline" className="bg-purple-50 text-purple-700 dark:bg-purple-950 dark:text-purple-300 border-purple-200">Unit / Lot</Badge>
    case 'box':
      return <Badge variant="outline" className="bg-slate-50 text-slate-700 dark:bg-slate-900 dark:text-slate-300 border-slate-200">Balok (P × L × T)</Badge>
    default:
      return <Badge variant="secondary">Standar</Badge>
  }
}

function VolumePage() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const searchParams = Route.useSearch()
  const queryProjectId = searchParams.projectId || ''
  const { isAdmin, user, openLoginModal } = useAuth()

  const [selectedProjectId, setSelectedProjectId] = useState<string>(queryProjectId)
  const [search, setSearch] = useState('')

  // Fetch projects for selector
  const { data: projectsData } = useQuery({
    ...listProjectsOptions({ query: { per_page: 100 } }),
  })
  const projects = projectsData?.data ?? []

  // Default project selection if not set
  const activeProjectId = selectedProjectId || (projects.length > 0 ? projects[0].id : '')
  const activeProject = projects.find((p) => p.id === activeProjectId)

  // Fetch volume rows for active project
  const { data: volumeRows = [], isLoading } = useQuery({
    ...listVolumesOptions({
      path: { id: activeProjectId },
      query: { q: search || undefined },
    }),
    enabled: !!activeProjectId,
  })

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
    navigate({
      to: '/volume/new',
      search: { projectId: activeProjectId },
    })
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
      <PageHeader
        title="Volume Take-off"
        onSearch={setSearch}
        onAdd={handleOpenAdd}
        addLabel="Tambah Take-off"
      />

      {/* Project Selector & Status Bar */}
      <div className="flex flex-wrap items-center justify-between gap-4 rounded-lg border bg-card p-3 shadow-xs">
        <div className="flex items-center gap-3">
          <span className="text-sm font-semibold text-foreground whitespace-nowrap">Proyek:</span>
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

          {activeProject && (
            <div className="hidden sm:flex items-center gap-2">
              <Badge variant="secondary" className="text-xs">
                {activeProject.plant}
              </Badge>
              <Badge variant="outline" className="text-xs">
                {activeProject.section}
              </Badge>
            </div>
          )}
        </div>

        {/* Action Export */}
        {activeProjectId && (
          <a
            href={`/api/v1/projects/${activeProjectId}/exports/volume.xlsx`}
            download
            className={buttonVariants({
              variant: 'outline',
              size: 'sm',
              className: 'gap-2 text-xs',
            })}
          >
            <FileSpreadsheet className="size-3.5 text-emerald-600" />
            <span>Export Excel</span>
          </a>
        )}
      </div>

      {/* Table Data Take-off */}
      <div className="rounded-md border bg-card overflow-x-auto shadow-xs">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-10 text-center">No</TableHead>
              <TableHead className="min-w-44">Work Item</TableHead>
              <TableHead className="w-28 text-center">Metode / Rumus</TableHead>
              <TableHead className="w-24 text-right">Planned Vol</TableHead>
              <TableHead className="w-16 text-center">Unit</TableHead>
              <TableHead className="w-24 text-center">Date</TableHead>
              <TableHead className="w-40 text-center">Dimensi Pengukuran</TableHead>
              <TableHead className="w-24 text-right">Quantity</TableHead>
              <TableHead className="w-24 text-right">Unit / Wt</TableHead>
              <TableHead className="w-28 text-right font-bold">Total Qty</TableHead>
              <TableHead className="min-w-32">Description</TableHead>
              {isAdmin && <TableHead className="w-16 text-center">Action</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 12 : 11} className="h-32 text-center text-muted-foreground">
                  Memuat data volume take-off...
                </TableCell>
              </TableRow>
            ) : volumeRows.length === 0 ? (
              <TableRow>
                <TableCell colSpan={isAdmin ? 12 : 11} className="h-48 text-center">
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
                  <TableCell className="text-center">{getFormulaBadge(r.formula)}</TableCell>
                  <TableCell className="text-right text-muted-foreground">{formatDecimal(r.planned_volume)}</TableCell>
                  <TableCell className="text-center text-muted-foreground font-mono text-xs">{r.unit}</TableCell>
                  <TableCell className="text-center text-muted-foreground">{r.date || '-'}</TableCell>
                  <TableCell className="text-center text-muted-foreground text-xs font-mono">
                    {renderDimensions(r)}
                  </TableCell>
                  <TableCell className="text-right text-foreground font-medium">{formatDecimal(r.quantity)}</TableCell>
                  <TableCell className="text-right text-muted-foreground">{formatDecimal(r.unit_weight)}</TableCell>
                  <TableCell className="text-right font-bold text-foreground text-emerald-600 dark:text-emerald-400">
                    {formatDecimal(r.total_quantity)}
                  </TableCell>
                  <TableCell className="text-muted-foreground text-xs">{r.description || '-'}</TableCell>
                  {isAdmin && (
                    <TableCell className="text-center">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="text-xs text-destructive hover:bg-destructive/10"
                        onClick={() => handleDelete(r.id, r.work_item)}
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
    </div>
  )
}
