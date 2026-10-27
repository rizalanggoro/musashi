import { cn } from 'cn'

export { cn }

export function formatRupiah(value: string | number | null | undefined): string {
  if (value === null || value === undefined || value === '') return 'Rp0'
  const num = typeof value === 'string' ? Number.parseInt(value, 10) : value
  if (Number.isNaN(num)) return 'Rp0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(num)
}

export function formatDecimal(value: string | number | null | undefined, digits = 2): string {
  if (value === null || value === undefined || value === '') return '0'
  const num = typeof value === 'string' ? Number.parseFloat(value) : value
  if (Number.isNaN(num)) return '0'
  return new Intl.NumberFormat('id-ID', {
    minimumFractionDigits: 0,
    maximumFractionDigits: digits,
  }).format(num)
}
