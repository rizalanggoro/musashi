import { formatDecimal } from '@/lib/utils'

export type FormulaType =
  | 'box'          // Balok / Kotak: P × L × T
  | 'cylinder'     // Silinder / Lubang Bor: ¼ × π × D² × T
  | 'area_pl'      // Luas Lantai / Plafon: P × L
  | 'area_pt'      // Luas Dinding / Pagar: P × T
  | 'area_lt'      // Luas Samping: L × T
  | 'area_circle'  // Luas Lingkaran: ¼ × π × D²
  | 'linear_p'     // Meter Lari - Panjang: P
  | 'linear_t'     // Meter Lari - Tinggi/Tiang: T
  | 'weight_rebar' // Berat Besi/Baja: P × Berat Jenis
  | 'unit_only'    // Kuantitas Unit / Lot Saja
  | 'custom'       // Planned Volume × Unit

export interface FormulaConfig {
  id: FormulaType
  label: string
  category: string
  description: string
  fields: {
    length?: string
    width?: string
    height?: string
    weight: string
  }
}

export const FORMULA_CONFIGS: Record<FormulaType, FormulaConfig> = {
  box: {
    id: 'box',
    label: 'Kotak / Balok (P × L × T)',
    category: 'm3',
    description: 'Pekerjaan kubikal 3 dimensi: galian, pondasi telapak, balok/kolom beton',
    fields: { length: 'Panjang (m)', width: 'Lebar (m)', height: 'Tinggi (m)', weight: 'Unit / Kali (Pengali)' },
  },
  cylinder: {
    id: 'cylinder',
    label: 'Silinder / Lubang Bor (¼ × π × D² × T)',
    category: 'm3',
    description: 'Bore pile, lubang bor, sumuran, tiang bulat: ¼ × π × D² × Kedalaman × Titik',
    fields: { length: 'Diameter D (m)', height: 'Kedalaman / Tinggi T (m)', weight: 'Jumlah Titik / Lubang (Unit)' },
  },
  area_pl: {
    id: 'area_pl',
    label: 'Luas Lantai / Plafon (P × L)',
    category: 'm2',
    description: 'Luas bidang horizontal: lantai keramik, screed, dak, plafon',
    fields: { length: 'Panjang (m)', width: 'Lebar (m)', weight: 'Unit / Kali (Pengali)' },
  },
  area_pt: {
    id: 'area_pt',
    label: 'Luas Dinding / Pagar (P × T)',
    category: 'm2',
    description: 'Luas bidang vertikal: dinding bata/hebel, plesteran, acian, cat, partisi',
    fields: { length: 'Panjang (m)', height: 'Tinggi (m)', weight: 'Jumlah Sisi / Unit (misal 2 utk 2 sisi)' },
  },
  area_lt: {
    id: 'area_lt',
    label: 'Luas Samping (L × T)',
    category: 'm2',
    description: 'Luas penampang vertikal samping: bordes tangga, sisi kolom',
    fields: { width: 'Lebar (m)', height: 'Tinggi (m)', weight: 'Unit / Kali (Pengali)' },
  },
  area_circle: {
    id: 'area_circle',
    label: 'Luas Lingkaran (¼ × π × D²)',
    category: 'm2',
    description: 'Luas penampang lingkaran: tutup manhole, plat bulat',
    fields: { length: 'Diameter D (m)', weight: 'Jumlah Unit (Buah)' },
  },
  linear_p: {
    id: 'linear_p',
    label: "Meter Lari - Panjang (P)",
    category: "m'",
    description: 'Pengukuran panjang linear horizontal: plin lantai, pipa air, kabel tray',
    fields: { length: 'Panjang P (m)', weight: 'Jumlah Jalur / Unit' },
  },
  linear_t: {
    id: 'linear_t',
    label: "Meter Lari - Tinggi/Tiang (T)",
    category: "m'",
    description: 'Pengukuran linear vertikal: tiang listrik, pipa tegak, kolom profil',
    fields: { height: 'Tinggi T (m)', weight: 'Jumlah Titik / Tiang' },
  },
  weight_rebar: {
    id: 'weight_rebar',
    label: 'Berat Besi / Baja (P × Berat Jenis)',
    category: 'kg',
    description: 'Pembesian tulangan / profil baja: Panjang × Berat Jenis kg/m\' × Batang',
    fields: { length: 'Panjang Batang (m)', width: "Berat Jenis (kg/m')", weight: 'Jumlah Batang (Unit)' },
  },
  unit_only: {
    id: 'unit_only',
    label: 'Kuantitas Unit / Lot Saja',
    category: 'lot/unit',
    description: 'Kuantitas langsung dari jumlah unit/lot: sewa alat, saklar, lampu, sanitair',
    fields: { weight: 'Jumlah Kuantitas / Unit (misal 1 lot, 8 bh)' },
  },
  custom: {
    id: 'custom',
    label: 'Manual / Planned Volume',
    category: 'lainnya',
    description: 'Perhitungan langsung dari planned volume × pengali unit',
    fields: { weight: 'Pengali Unit' },
  },
}

export function detectFormula(unitStr: string): FormulaType {
  const u = unitStr.toLowerCase().trim()
  if (['lot', 'ls', 'lump sum', 'unit', 'pcs', 'titik', 'buah', 'set', 'btg', 'bh'].includes(u)) {
    return 'unit_only'
  }
  if (['m2', 'm²'].includes(u)) {
    return 'area_pl'
  }
  if (['m3', 'm³'].includes(u)) {
    return 'box'
  }
  if (["m'", 'm', 'meter', 'ml'].includes(u)) {
    return 'linear_p'
  }
  if (['kg', 'ton'].includes(u)) {
    return 'weight_rebar'
  }
  return 'box'
}

export function computeVolume(
  formula: FormulaType,
  lStr: string,
  wStr: string,
  hStr: string,
  unitWeightStr: string,
  plannedStr: string
) {
  const l = parseFloat(lStr || '0') || 0
  const w = parseFloat(wStr || '0') || 0
  const h = parseFloat(hStr || '0') || 0
  const uw = parseFloat(unitWeightStr || '1') || 1
  const planned = parseFloat(plannedStr || '0') || 0

  let qty = 0
  let expression = ''

  switch (formula) {
    case 'cylinder': {
      // 1/4 * pi * D^2 * T
      const area = 0.25 * Math.PI * (l * l)
      const depth = h > 0 ? h : 1
      qty = area * depth
      expression = `¼ × π × (${formatDecimal(l)} m)² × ${formatDecimal(h)} m × ${formatDecimal(uw)} titik`
      break
    }
    case 'area_circle': {
      qty = 0.25 * Math.PI * (l * l)
      expression = `¼ × π × (${formatDecimal(l)} m)² × ${formatDecimal(uw)} unit`
      break
    }
    case 'area_pl': {
      qty = l * w
      expression = `${formatDecimal(l)} m × ${formatDecimal(w)} m × ${formatDecimal(uw)} unit`
      break
    }
    case 'area_pt': {
      qty = l * h
      expression = `${formatDecimal(l)} m × ${formatDecimal(h)} m × ${formatDecimal(uw)} sisi`
      break
    }
    case 'area_lt': {
      qty = w * h
      expression = `${formatDecimal(w)} m × ${formatDecimal(h)} m × ${formatDecimal(uw)} unit`
      break
    }
    case 'linear_p': {
      qty = l
      expression = `${formatDecimal(l)} m × ${formatDecimal(uw)} jalur`
      break
    }
    case 'linear_t': {
      qty = h
      expression = `${formatDecimal(h)} m × ${formatDecimal(uw)} titik`
      break
    }
    case 'weight_rebar': {
      const bj = w > 0 ? w : 1
      qty = l * bj
      expression = `${formatDecimal(l)} m × ${formatDecimal(w)} kg/m' × ${formatDecimal(uw)} btg`
      break
    }
    case 'unit_only': {
      qty = 1
      expression = `1.0000 × ${formatDecimal(uw)} unit`
      break
    }
    case 'box': {
      qty = l * w * h
      expression = `${formatDecimal(l)} m × ${formatDecimal(w)} m × ${formatDecimal(h)} m × ${formatDecimal(uw)} unit`
      break
    }
    case 'custom':
    default: {
      qty = planned > 0 ? planned : (l > 0 && w > 0 && h > 0 ? l * w * h : (l > 0 && w > 0 ? l * w : l))
      expression = `${formatDecimal(qty)} × ${formatDecimal(uw)} unit`
      break
    }
  }

  const totalQty = qty * uw
  return { qty, totalQty, expression }
}

export function renderDimensions(row: {
  formula?: string
  length?: string
  width?: string
  height?: string
  unit_weight?: string
}) {
  const f = row.formula || ''
  const l = row.length ? formatDecimal(row.length) : ''
  const w = row.width ? formatDecimal(row.width) : ''
  const h = row.height ? formatDecimal(row.height) : ''

  switch (f) {
    case 'cylinder':
      return l || h ? `D: ${l || '-'} m, T: ${h || '-'} m` : '-'
    case 'area_circle':
      return l ? `D: ${l} m` : '-'
    case 'area_pl':
      return l || w ? `P: ${l || '-'} m, L: ${w || '-'} m` : '-'
    case 'area_pt':
      return l || h ? `P: ${l || '-'} m, T: ${h || '-'} m` : '-'
    case 'area_lt':
      return w || h ? `L: ${w || '-'} m, T: ${h || '-'} m` : '-'
    case 'linear_p':
      return l ? `P: ${l} m` : '-'
    case 'linear_t':
      return h ? `T: ${h} m` : '-'
    case 'weight_rebar':
      return l || w ? `P: ${l || '-'} m, BJ: ${w || '-'} kg/m'` : '-'
    case 'unit_only':
      return '-'
    case 'box':
    default:
      if (l && w && h) return `${l} × ${w} × ${h} m`
      if (l && w) return `${l} × ${w} m`
      if (l && h) return `${l} × ${h} m`
      if (l) return `P: ${l} m`
      return '-'
  }
}

