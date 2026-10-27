import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

type Props = {
  title: string
  onSearch?: (value: string) => void
  onAdd?: () => void
}

export function PageHeader({ title, onSearch, onAdd }: Props) {
  return (
    <div className="mb-6 flex flex-wrap items-center gap-3">
      <h1 className="bg-primary px-6 py-2.5 text-lg font-semibold tracking-wide text-primary-foreground">
        {title}
      </h1>

      <div className="ml-auto flex items-center gap-3">
        <Input
          placeholder="SEARCH"
          className="h-10 w-64 border-0 bg-muted px-4 placeholder:text-muted-foreground"
          onChange={(e) => onSearch?.(e.target.value)}
        />
        <Button className="h-10 px-6" onClick={onAdd}>
          ADD
        </Button>
      </div>
    </div>
  )
}
