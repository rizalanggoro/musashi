import { useMemo } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  listProjectsOptions,
  createVolumeMutation,
  listAhspsOptions,
  getProjectBoqOptions,
} from '@/api/@tanstack/react-query.gen'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FieldGroup, Field, FieldLabel, FieldError } from '@/components/ui/field'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { formatDecimal } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import {
  FORMULA_CONFIGS,
  detectFormula,
  computeVolume,
  type FormulaType,
} from '@/lib/volume'
import { ArrowLeft, Calculator, Save, Building2, Ruler } from 'lucide-react'

export const Route = createFileRoute('/volume/new/')({
  component: VolumeNewPage,
  validateSearch: (search: Record<string, unknown>): { projectId?: string } => ({
    projectId: (search.projectId as string) || undefined,
  }),
})

const volumeFormSchema = z.object({
  projectId: z.string().min(1, 'Pilih proyek terlebih dahulu'),
  boqItemId: z.string(),
  ahspId: z.string(),
  workItem: z.string().min(1, 'Nama pekerjaan (Work Item) wajib diisi'),
  unit: z.string().min(1, 'Satuan (Unit) wajib diisi'),
  formula: z.enum([
    'box',
    'cylinder',
    'area_pl',
    'area_pt',
    'area_lt',
    'area_circle',
    'linear_p',
    'linear_t',
    'weight_rebar',
    'unit_only',
    'custom',
  ]),
  plannedVolume: z.string(),
  date: z.string(),
  length: z.string(),
  width: z.string(),
  height: z.string(),
  unitWeight: z.string().min(1, 'Unit / pengali wajib diisi'),
  description: z.string(),
})

type VolumeFormValues = z.infer<typeof volumeFormSchema>

function VolumeNewPage() {
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const searchParams = Route.useSearch()
  const { user, openLoginModal } = useAuth()

  // Fetch projects list
  const { data: projectsData } = useQuery({
    ...listProjectsOptions({ query: { per_page: 100 } }),
  })
  const projects = projectsData?.data ?? []

  // Default projectId from search query or first project
  const initialProjectId = searchParams.projectId || (projects.length > 0 ? projects[0].id : '')

  // React Hook Form with Zod validation
  const {
    register,
    control,
    handleSubmit,
    watch,
    setValue,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<VolumeFormValues>({
    resolver: zodResolver(volumeFormSchema) as any,
    defaultValues: {
      projectId: initialProjectId,
      boqItemId: 'custom',
      ahspId: '',
      workItem: '',
      unit: '',
      formula: 'box',
      plannedVolume: '0',
      date: new Date().toISOString().slice(0, 10),
      length: '',
      width: '',
      height: '',
      unitWeight: '1',
      description: '',
    },
  })

  // Watched form values
  const activeProjectId = watch('projectId') || initialProjectId
  const selectedBoqId = watch('boqItemId')
  const selectedAhspId = watch('ahspId')
  const formula = watch('formula')
  const unit = watch('unit')
  const length = watch('length')
  const width = watch('width')
  const height = watch('height')
  const unitWeight = watch('unitWeight')
  const plannedVolume = watch('plannedVolume')

  const activeProject = projects.find((p) => p.id === activeProjectId)

  // Fetch existing BoQ items of this project
  const { data: projectBoqItems = [] } = useQuery({
    ...getProjectBoqOptions({ path: { id: activeProjectId } }),
    enabled: !!activeProjectId,
  })

  // Fetch AHSPs for master reference
  const { data: ahsps = [] } = useQuery({
    ...listAhspsOptions({}),
  })

  const projectItems = useMemo(
    () =>
      projects.map((p) => ({
        value: p.id,
        label: `${p.name} (${p.plant} - ${p.section})`,
      })),
    [projects]
  )

  const boqSelectItems = useMemo(
    () => [
      { value: 'custom', label: '+ Buat Item Baru (Input Bebas / Master AHSP)' },
      ...projectBoqItems.map((b) => ({
        value: b.id,
        label: `${b.work_item} (${b.unit}) — Plan: ${formatDecimal(b.planned_volume)}`,
      })),
    ],
    [projectBoqItems]
  )

  const ahspSelectItems = useMemo(
    () =>
      ahsps.map((a) => ({
        value: a.id,
        label: `${a.code} - ${a.work_item} (${a.unit})`,
      })),
    [ahsps]
  )

  const formulaSelectItems = useMemo(
    () =>
      Object.values(FORMULA_CONFIGS).map((cfg) => ({
        value: cfg.id,
        label: cfg.label,
      })),
    []
  )

  const createMut = useMutation(createVolumeMutation())

  // Handle selecting an existing BoQ item
  const handleSelectBoqItem = (boqId: string) => {
    setValue('boqItemId', boqId)
    if (boqId === 'custom') {
      setValue('workItem', '')
      setValue('unit', '')
      setValue('plannedVolume', '0')
      return
    }

    const item = projectBoqItems.find((b) => b.id === boqId)
    if (item) {
      setValue('workItem', item.work_item)
      setValue('unit', item.unit)
      setValue('plannedVolume', item.planned_volume || '0')
      setValue('ahspId', item.ahsp_id || '')
      const detected = detectFormula(item.unit)
      setValue('formula', detected)
    }
  }

  // Handle selecting AHSP master
  const handleSelectAhsp = (ahspId: string) => {
    setValue('ahspId', ahspId)
    const ahsp = ahsps.find((a) => a.id === ahspId)
    if (ahsp) {
      setValue('workItem', ahsp.work_item)
      setValue('unit', ahsp.unit)
      const detected = detectFormula(ahsp.unit)
      setValue('formula', detected)
    }
  }

  // Auto-detect formula on unit change
  const handleUnitChange = (newUnit: string) => {
    setValue('unit', newUnit)
    if (newUnit.trim()) {
      setValue('formula', detectFormula(newUnit))
    }
  }

  const currentFormulaConfig = FORMULA_CONFIGS[formula] || FORMULA_CONFIGS.box

  // Live calculation preview
  const preview = useMemo(() => {
    return computeVolume(formula, length, width, height, unitWeight, plannedVolume)
  }, [formula, length, width, height, unitWeight, plannedVolume])

  const onSubmit = async (values: VolumeFormValues) => {
    if (!user) {
      openLoginModal()
      return
    }

    try {
      await createMut.mutateAsync({
        path: { id: values.projectId },
        body: {
          boq_item_id: values.boqItemId !== 'custom' ? values.boqItemId : undefined,
          ahsp_id: values.ahspId || undefined,
          work_item: values.workItem.trim(),
          unit: values.unit.trim(),
          planned_volume: values.plannedVolume || '0',
          date: values.date || undefined,
          length: currentFormulaConfig.fields.length ? values.length || undefined : undefined,
          width: currentFormulaConfig.fields.width ? values.width || undefined : undefined,
          height: currentFormulaConfig.fields.height ? values.height || undefined : undefined,
          formula: values.formula,
          quantity: preview.qty.toFixed(4),
          unit_weight: values.unitWeight || '1',
          description: values.description.trim(),
        },
      })
      await queryClient.invalidateQueries()
      navigate({
        to: '/volume',
        search: { projectId: values.projectId },
      })
    } catch (err: any) {
      setError('root', {
        message: err?.message || 'Gagal menyimpan data volume take-off.',
      })
    }
  }

  return (
    <div className="flex flex-col gap-6 max-w-6xl mx-auto pb-12">
      {/* Top Navigation & Breadcrumb */}
      <div className="flex items-center justify-between border-b pb-4">
        <div className="flex items-center gap-3">
          <Button
            variant="ghost"
            size="sm"
            className="gap-2 text-muted-foreground hover:text-foreground"
            onClick={() => navigate({ to: '/volume', search: { projectId: activeProjectId } })}
          >
            <ArrowLeft className="size-4" />
            <span>Kembali ke Daftar Volume</span>
          </Button>
          <div className="h-4 w-px bg-border" />
          <div>
            <h1 className="text-xl font-bold tracking-tight text-foreground flex items-center gap-2">
              <Calculator className="size-5 text-primary" />
              <span>Tambah Volume Take-off</span>
            </h1>
            <p className="text-xs text-muted-foreground mt-0.5">
              Halaman input perhitungan volume terhubung langsung dengan BoQ dan spesifikasi pengukuran proyek.
            </p>
          </div>
        </div>

        {activeProject && (
          <div className="flex items-center gap-2">
            <Badge variant="secondary" className="text-xs">
              {activeProject.plant}
            </Badge>
            <Badge variant="outline" className="text-xs">
              {activeProject.section}
            </Badge>
          </div>
        )}
      </div>

      {errors.root?.message && (
        <Alert variant="destructive">
          <AlertDescription>{errors.root.message}</AlertDescription>
        </Alert>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        {/* Left Column: Form Inputs (8 cols) */}
        <div className="lg:col-span-8 flex flex-col gap-6">
          <FieldGroup>
            {/* Card 1: Informasi Proyek & Item Pekerjaan */}
            <Card className="shadow-xs">
              <CardHeader className="pb-4">
                <CardTitle className="text-base flex items-center gap-2">
                  <Building2 className="size-4 text-primary" />
                  <span>Pilih Proyek & Item Pekerjaan</span>
                </CardTitle>
                <CardDescription className="text-xs">
                  Hubungkan pengukuran ini dengan item BoQ proyek atau buat pekerjaan baru.
                </CardDescription>
              </CardHeader>
              <CardContent className="flex flex-col gap-4">
                {/* Proyek Selector */}
                <Field data-invalid={!!errors.projectId}>
                  <FieldLabel htmlFor="new-vol-project">Proyek Terkait</FieldLabel>
                  <Controller
                    name="projectId"
                    control={control}
                    render={({ field }) => (
                      <Select
                        items={projectItems}
                        value={field.value}
                        onValueChange={(val) => {
                          if (val) {
                            field.onChange(val)
                            setValue('boqItemId', 'custom')
                            navigate({
                              search: ((prev: any) => ({ ...prev, projectId: val })) as any,
                              replace: true,
                            })
                          }
                        }}
                      >
                        <SelectTrigger id="new-vol-project" className="w-full">
                          <SelectValue placeholder="Pilih Proyek" />
                        </SelectTrigger>
                        <SelectContent className="w-full min-w-72">
                          <SelectGroup>
                            {projectItems.map((p) => (
                              <SelectItem key={p.value} value={p.value}>
                                {p.label}
                              </SelectItem>
                            ))}
                          </SelectGroup>
                        </SelectContent>
                      </Select>
                    )}
                  />
                  {errors.projectId && <FieldError>{errors.projectId.message}</FieldError>}
                </Field>

                {/* BoQ Selector */}
                <Field>
                  <FieldLabel htmlFor="new-vol-boq">Pilih dari Item BoQ Proyek</FieldLabel>
                  <Select items={boqSelectItems} value={selectedBoqId} onValueChange={(val) => { if (val) handleSelectBoqItem(val) }}>
                    <SelectTrigger id="new-vol-boq" className="w-full">
                      <SelectValue placeholder="-- Pilih Item BoQ Proyek --" />
                    </SelectTrigger>
                    <SelectContent className="w-full min-w-72">
                      <SelectGroup>
                        {boqSelectItems.map((b) => (
                          <SelectItem key={b.value} value={b.value}>
                            {b.label}
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </Field>

                {/* Master AHSP picker when custom */}
                {selectedBoqId === 'custom' && ahsps.length > 0 && (
                  <Field>
                    <FieldLabel htmlFor="new-vol-ahsp">Pilih dari Master AHSP (Opsional)</FieldLabel>
                    <Select items={ahspSelectItems} value={selectedAhspId} onValueChange={(val) => { if (val) handleSelectAhsp(val) }}>
                      <SelectTrigger id="new-vol-ahsp" className="w-full">
                        <SelectValue placeholder="-- Salin nama pekerjaan dari AHSP --" />
                      </SelectTrigger>
                      <SelectContent className="w-full min-w-72">
                        <SelectGroup>
                          {ahspSelectItems.map((a) => (
                            <SelectItem key={a.value} value={a.value}>
                              {a.label}
                            </SelectItem>
                          ))}
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                  </Field>
                )}

                {/* Work Item & Unit */}
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  <div className="sm:col-span-2">
                    <Field data-invalid={!!errors.workItem}>
                      <FieldLabel htmlFor="new-vol-workitem">Nama Pekerjaan (Work Item)</FieldLabel>
                      <Input
                        id="new-vol-workitem"
                        {...register('workItem')}
                        placeholder="Contoh: Pengecoran Bore Pile D40 cm..."
                        aria-invalid={!!errors.workItem}
                      />
                      {errors.workItem && <FieldError>{errors.workItem.message}</FieldError>}
                    </Field>
                  </div>
                  <Field data-invalid={!!errors.unit}>
                    <FieldLabel htmlFor="new-vol-unit">Satuan (Unit)</FieldLabel>
                    <Input
                      id="new-vol-unit"
                      value={unit}
                      onChange={(e) => handleUnitChange(e.target.value)}
                      placeholder="m3, m2, m', kg, lot..."
                      aria-invalid={!!errors.unit}
                    />
                    {errors.unit && <FieldError>{errors.unit.message}</FieldError>}
                  </Field>
                </div>

                {/* Planned Volume & Date */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <Field>
                    <FieldLabel htmlFor="new-vol-planned">Target Volume (Planned Volume)</FieldLabel>
                    <Input
                      id="new-vol-planned"
                      type="number"
                      step="0.0001"
                      {...register('plannedVolume')}
                    />
                  </Field>
                  <Field>
                    <FieldLabel htmlFor="new-vol-date">Tanggal Pengukuran</FieldLabel>
                    <Input
                      id="new-vol-date"
                      type="date"
                      {...register('date')}
                    />
                  </Field>
                </div>
              </CardContent>
            </Card>

            {/* Card 2: Metode Rumus & Dimensi Pengukuran */}
            <Card className="shadow-xs">
              <CardHeader className="pb-4">
                <CardTitle className="text-base flex items-center gap-2">
                  <Ruler className="size-4 text-primary" />
                  <span>Metode Rumus & Dimensi Pengukuran</span>
                </CardTitle>
                <CardDescription className="text-xs">
                  Pilih rumus yang sesuai dengan karakteristik fisik pekerjaan dan masukkan dimensi lapangannya.
                </CardDescription>
              </CardHeader>
              <CardContent className="flex flex-col gap-4">
                {/* Formula Select */}
                <Field data-invalid={!!errors.formula}>
                  <div className="flex items-center justify-between">
                    <FieldLabel htmlFor="new-vol-formula">Pilihan Rumus Perhitungan</FieldLabel>
                    <Badge variant="outline" className="text-[11px] font-mono">
                      Kategori: {currentFormulaConfig.category}
                    </Badge>
                  </div>
                  <Controller
                    name="formula"
                    control={control}
                    render={({ field }) => (
                      <Select
                        items={formulaSelectItems}
                        value={field.value}
                        onValueChange={(val) => {
                          if (val) field.onChange(val as FormulaType)
                        }}
                      >
                        <SelectTrigger id="new-vol-formula" className="w-full">
                          <SelectValue placeholder="Pilih Rumus" />
                        </SelectTrigger>
                        <SelectContent className="w-full min-w-72">
                          <SelectGroup>
                            {formulaSelectItems.map((cfg) => (
                              <SelectItem key={cfg.value} value={cfg.value}>
                                {cfg.label}
                              </SelectItem>
                            ))}
                          </SelectGroup>
                        </SelectContent>
                      </Select>
                    )}
                  />
                  <p className="text-xs text-muted-foreground mt-1">
                    {currentFormulaConfig.description}
                  </p>
                </Field>

                {/* Dynamic Dimension Inputs */}
                <div className="rounded-lg border bg-muted/20 p-4">
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                    {currentFormulaConfig.fields.length && (
                      <Field>
                        <FieldLabel htmlFor="new-vol-l">{currentFormulaConfig.fields.length}</FieldLabel>
                        <Input
                          id="new-vol-l"
                          type="number"
                          step="0.0001"
                          {...register('length')}
                          placeholder="0.0000"
                        />
                      </Field>
                    )}

                    {currentFormulaConfig.fields.width && (
                      <Field>
                        <FieldLabel htmlFor="new-vol-w">{currentFormulaConfig.fields.width}</FieldLabel>
                        <Input
                          id="new-vol-w"
                          type="number"
                          step="0.0001"
                          {...register('width')}
                          placeholder="0.0000"
                        />
                      </Field>
                    )}

                    {currentFormulaConfig.fields.height && (
                      <Field>
                        <FieldLabel htmlFor="new-vol-h">{currentFormulaConfig.fields.height}</FieldLabel>
                        <Input
                          id="new-vol-h"
                          type="number"
                          step="0.0001"
                          {...register('height')}
                          placeholder="0.0000"
                        />
                      </Field>
                    )}

                    <Field data-invalid={!!errors.unitWeight}>
                      <FieldLabel htmlFor="new-vol-uw">{currentFormulaConfig.fields.weight}</FieldLabel>
                      <Input
                        id="new-vol-uw"
                        type="number"
                        step="0.0001"
                        {...register('unitWeight')}
                        placeholder="1"
                        aria-invalid={!!errors.unitWeight}
                      />
                      {errors.unitWeight && <FieldError>{errors.unitWeight.message}</FieldError>}
                    </Field>
                  </div>
                </div>

                {/* Description / Note */}
                <Field>
                  <FieldLabel htmlFor="new-vol-desc">Keterangan / Catatan Lokasi Pengukuran</FieldLabel>
                  <Input
                    id="new-vol-desc"
                    {...register('description')}
                    placeholder="Contoh: Titik Bore Pile P-01 s/d P-12 Sisi Timur Bak IPAL..."
                  />
                </Field>
              </CardContent>
            </Card>
          </FieldGroup>
        </div>

        {/* Right Column: Sticky Summary & Calculation Card (4 cols) */}
        <div className="lg:col-span-4 sticky top-6 flex flex-col gap-4">
          <Card className="border-primary/30 shadow-md bg-card">
            <CardHeader className="pb-3 border-b bg-primary/5 rounded-t-xl">
              <CardTitle className="text-sm font-semibold flex items-center justify-between text-primary">
                <span>Ringkasan Kalkulasi</span>
                <Calculator className="size-4" />
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-4 flex flex-col gap-4">
              {/* Formula Expression */}
              <div className="flex flex-col gap-1 p-3 rounded-md bg-muted/50 border text-xs">
                <span className="text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">
                  Persamaan Rumus:
                </span>
                <span className="font-mono text-foreground font-medium text-xs leading-relaxed break-words">
                  {preview.expression || '-'}
                </span>
              </div>

              {/* Volume KPI */}
              <div className="flex flex-col gap-3 py-2">
                <div className="flex items-center justify-between text-xs text-muted-foreground">
                  <span>Kuantitas Satuan:</span>
                  <strong className="text-foreground font-mono text-sm">
                    {formatDecimal(preview.qty)} {unit || '-'}
                  </strong>
                </div>

                <div className="flex items-center justify-between text-xs text-muted-foreground">
                  <span>Pengali / Unit:</span>
                  <span className="font-mono font-medium text-foreground">
                    × {formatDecimal(unitWeight || 1)}
                  </span>
                </div>

                <div className="h-px bg-border my-1" />

                <div className="flex flex-col gap-1">
                  <span className="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
                    Total Volume Take-off:
                  </span>
                  <div className="flex items-baseline gap-2">
                    <span className="text-2xl font-bold font-mono text-emerald-600 dark:text-emerald-400">
                      {formatDecimal(preview.totalQty)}
                    </span>
                    <span className="text-sm font-semibold text-muted-foreground font-mono">
                      {unit || ''}
                    </span>
                  </div>
                </div>
              </div>

              {/* Target vs Total Progress Info */}
              {parseFloat(plannedVolume || '0') > 0 && (
                <div className="rounded-lg bg-muted/40 p-3 text-xs flex flex-col gap-1.5 border">
                  <div className="flex justify-between text-[11px]">
                    <span className="text-muted-foreground">Volume Rencana (BoQ):</span>
                    <span className="font-mono font-medium">{formatDecimal(plannedVolume)} {unit}</span>
                  </div>
                  <div className="flex justify-between text-[11px]">
                    <span className="text-muted-foreground">Persentase terhadap BoQ:</span>
                    <span className="font-mono font-bold text-primary">
                      {((preview.totalQty / parseFloat(plannedVolume)) * 100).toFixed(1)}%
                    </span>
                  </div>
                </div>
              )}

              {/* Actions */}
              <div className="flex flex-col gap-2 pt-2 border-t">
                <Button
                  type="submit"
                  size="default"
                  className="w-full gap-2 font-semibold shadow-xs"
                  disabled={isSubmitting || createMut.isPending}
                >
                  <Save className="size-4" />
                  <span>
                    {isSubmitting || createMut.isPending ? 'Menyimpan Take-off...' : 'Simpan Take-off'}
                  </span>
                </Button>

                <Button
                  type="button"
                  variant="outline"
                  size="default"
                  className="w-full"
                  onClick={() => navigate({ to: '/volume', search: { projectId: activeProjectId } })}
                >
                  Batal
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </form>
    </div>
  )
}
