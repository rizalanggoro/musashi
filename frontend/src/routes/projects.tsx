import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { PageHeader } from '@/components/page-header'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

export const Route = createFileRoute('/projects')({ component: ProjectsPage })

function ProjectsPage() {
  const [query, setQuery] = useState('')

  return (
    <div>
      <PageHeader title="PROJECT" onSearch={setQuery} />

      <div className="overflow-x-auto border">
        <Table className="text-xs">
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead rowSpan={2} className="text-center">NO</TableHead>
              <TableHead rowSpan={2} className="text-center">PROJECT</TableHead>
              <TableHead rowSpan={2} className="text-center">SECTION</TableHead>
              <TableHead rowSpan={2} className="text-center">PLANT</TableHead>
              <TableHead rowSpan={2} className="text-center">IMPLEMENTATION PLAN</TableHead>
              <TableHead rowSpan={2} className="text-center">BUDGET</TableHead>
              <TableHead colSpan={3} className="border-b text-center">PROGRESS</TableHead>
              <TableHead rowSpan={2} className="text-center">INFORMATION</TableHead>
            </TableRow>
            <TableRow className="hover:bg-transparent">
              <TableHead className="text-center">SPB</TableHead>
              <TableHead className="text-center">RAB</TableHead>
              <TableHead className="text-center">XXX</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow className="hover:bg-muted/60">
              <TableCell colSpan={11} className="h-64 bg-muted text-center text-muted-foreground">
                {query ? `Tidak ada proyek yang cocok dengan "${query}".` : 'Belum ada data proyek.'}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
