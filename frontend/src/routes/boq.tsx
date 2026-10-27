import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { listBoqProjectsOptions } from '@/api/@tanstack/react-query.gen'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'

export const Route = createFileRoute('/boq')({ component: BoqPage })

function BoqPage() {
  const navigate = useNavigate()
  const [search, setSearch] = useState('')

  const { data: projects = [], isLoading } = useQuery({
    ...listBoqProjectsOptions({ query: { q: search || undefined } }),
  })

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'complete':
        return <Badge className="bg-emerald-600 text-[10px]">COMPLETE</Badge>
      case 'partial':
        return <Badge className="bg-amber-600 text-[10px]">PARTIAL</Badge>
      default:
        return <Badge variant="secondary" className="text-[10px]">EMPTY</Badge>
    }
  }

  const handleExport = (url: string) => {
    window.open(url, '_blank')
  }

  return (
    <div>
      {/* Header bar matching Mockup Hal. 8 */}
      <div className="mb-6 flex flex-wrap items-center gap-3">
        <h1 className="bg-primary px-6 py-2.5 text-lg font-semibold tracking-wide text-primary-foreground">
          BOQ
        </h1>

        <div className="flex items-center gap-2">
          <Button
            className="bg-primary text-white"
            onClick={() => navigate({ to: '/volume' })}
          >
            VOLUME
          </Button>
          <Button
            className="bg-primary text-white"
            onClick={() => navigate({ to: '/rab' })}
          >
            RAB
          </Button>
        </div>

        <div className="ml-auto">
          <Input
            placeholder="SEARCH"
            className="h-10 w-64 border-0 bg-muted px-4 placeholder:text-muted-foreground"
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>

      <div className="overflow-x-auto border bg-white shadow-sm">
        <Table className="text-xs">
          <TableHeader className="bg-white">
            <TableRow className="border-b hover:bg-transparent">
              <TableHead rowSpan={2} className="w-12 border-r text-center font-bold">NO</TableHead>
              <TableHead rowSpan={2} className="min-w-64 border-r text-center font-bold">PROJECT</TableHead>
              <TableHead colSpan={2} className="border-b border-r text-center font-bold">VOLUME</TableHead>
              <TableHead colSpan={2} className="border-b text-center font-bold">RAB</TableHead>
            </TableRow>
            <TableRow className="border-b hover:bg-transparent">
              <TableHead className="w-28 border-r text-center font-bold">STATUS</TableHead>
              <TableHead className="w-32 border-r text-center font-bold">EXCEL</TableHead>
              <TableHead className="w-28 border-r text-center font-bold">STATUS</TableHead>
              <TableHead className="w-44 text-center font-bold">EXCEL / PDF</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={6} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  Memuat data...
                </TableCell>
              </TableRow>
            ) : projects.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="h-40 bg-muted/30 text-center text-muted-foreground">
                  {search ? `Tidak ada proyek yang cocok dengan "${search}".` : 'Belum ada data proyek.'}
                </TableCell>
              </TableRow>
            ) : (
              projects.map((p, idx) => (
                <TableRow key={p.id} className="border-b bg-muted/40 hover:bg-muted/80">
                  <TableCell className="border-r text-center font-medium">{idx + 1}</TableCell>
                  <TableCell className="border-r font-semibold text-foreground">
                    <button
                      type="button"
                      className="text-left font-semibold text-foreground hover:underline"
                      onClick={() => navigate({ to: `/volume?projectId=${p.id}` })}
                    >
                      {p.name}
                    </button>
                    {p.section && <span className="ml-2 text-xs font-normal text-muted-foreground">({p.section})</span>}
                  </TableCell>
                  <TableCell className="border-r text-center font-medium">
                    {getStatusBadge(p.volume_status)}
                  </TableCell>
                  <TableCell className="border-r text-center">
                    <Button
                      size="sm"
                      variant="outline"
                      className="h-7 text-xs"
                      disabled={p.volume_status === 'empty'}
                      onClick={() => handleExport(`/api/v1/projects/${p.id}/exports/volume.xlsx`)}
                    >
                      Unduh XLSX
                    </Button>
                  </TableCell>
                  <TableCell className="border-r text-center font-medium">
                    {getStatusBadge(p.rab_status)}
                  </TableCell>
                  <TableCell className="text-center">
                    <div className="flex items-center justify-center gap-1.5">
                      <Button
                        size="sm"
                        variant="outline"
                        className="h-7 text-xs"
                        disabled={p.rab_status === 'empty'}
                        onClick={() => handleExport(`/api/v1/projects/${p.id}/exports/rab.xlsx`)}
                      >
                        XLSX
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        className="h-7 text-xs bg-primary/10 text-primary hover:bg-primary/20"
                        disabled={p.rab_status === 'empty'}
                        onClick={() => handleExport(`/api/v1/projects/${p.id}/exports/rab.pdf`)}
                      >
                        PDF
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
