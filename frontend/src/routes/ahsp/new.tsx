import { createFileRoute } from '@tanstack/react-router'
import { AhspForm } from '@/components/ahsp'

export const Route = createFileRoute('/ahsp/new')({
  component: AhspNewPage,
})

function AhspNewPage() {
  return <AhspForm />
}
