import { createFileRoute } from '@tanstack/react-router'
import { ModulePlaceholder } from '@/components/module-placeholder'

export const Route = createFileRoute('/unit-prices')({ component: () => <ModulePlaceholder title="UNIT PRICE" /> })
