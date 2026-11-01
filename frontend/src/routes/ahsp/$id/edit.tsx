import { createFileRoute } from '@tanstack/react-router'
import { AhspForm } from '@/components/ahsp'

export const Route = createFileRoute('/ahsp/$id/edit')({
  component: AhspEditPage,
})

function AhspEditPage() {
  const { id } = Route.useParams()
  return <AhspForm id={id} />
}
