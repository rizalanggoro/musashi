import { PageHeader } from '@/components/page-header'

export function ModulePlaceholder({ title }: { title: string }) {
  return (
    <div>
      <PageHeader title={title} />
      <div className="border p-16 text-center text-sm text-muted-foreground">
        Modul {title} belum diimplementasi.
      </div>
    </div>
  )
}
