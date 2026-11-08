import { Badge } from '@/components/ui/badge'

interface FormulaBadgeProps {
  formula?: string
}

export function FormulaBadge({ formula }: FormulaBadgeProps) {
  switch (formula) {
    case 'cylinder':
      return (
        <Badge variant="outline" className="bg-blue-50 text-blue-700 dark:bg-blue-950 dark:text-blue-300 border-blue-200">
          Silinder (D² × T)
        </Badge>
      )
    case 'area_circle':
      return (
        <Badge variant="outline" className="bg-sky-50 text-sky-700 dark:bg-sky-950 dark:text-sky-300 border-sky-200">
          Lingkaran (¼πD²)
        </Badge>
      )
    case 'area_pl':
      return (
        <Badge variant="outline" className="bg-emerald-50 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300 border-emerald-200">
          Luas (P × L)
        </Badge>
      )
    case 'area_pt':
      return (
        <Badge variant="outline" className="bg-teal-50 text-teal-700 dark:bg-teal-950 dark:text-teal-300 border-teal-200">
          Dinding (P × T)
        </Badge>
      )
    case 'area_lt':
      return (
        <Badge variant="outline" className="bg-cyan-50 text-cyan-700 dark:bg-cyan-950 dark:text-cyan-300 border-cyan-200">
          Samping (L × T)
        </Badge>
      )
    case 'linear_p':
    case 'linear_t':
      return (
        <Badge variant="outline" className="bg-indigo-50 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300 border-indigo-200">
          Meter Lari (m')
        </Badge>
      )
    case 'weight_rebar':
      return (
        <Badge variant="outline" className="bg-amber-50 text-amber-700 dark:bg-amber-950 dark:text-amber-300 border-amber-200">
          Besi (P × BJ)
        </Badge>
      )
    case 'unit_only':
      return (
        <Badge variant="outline" className="bg-purple-50 text-purple-700 dark:bg-purple-950 dark:text-purple-300 border-purple-200">
          Unit / Lot
        </Badge>
      )
    case 'box':
      return (
        <Badge variant="outline" className="bg-slate-50 text-slate-700 dark:bg-slate-900 dark:text-slate-300 border-slate-200">
          Balok (P × L × T)
        </Badge>
      )
    default:
      return <Badge variant="secondary">Standar</Badge>
  }
}

