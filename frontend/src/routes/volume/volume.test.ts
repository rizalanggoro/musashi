import { describe, it, expect } from 'bun:test'
import {
  detectFormula,
  computeVolume,
  renderDimensions,
} from './index'

describe('detectFormula', () => {
  it('detects unit_only for lot, ls, pcs, unit', () => {
    expect(detectFormula('lot')).toBe('unit_only')
    expect(detectFormula('LS')).toBe('unit_only')
    expect(detectFormula('pcs')).toBe('unit_only')
    expect(detectFormula('unit')).toBe('unit_only')
    expect(detectFormula('titik')).toBe('unit_only')
    expect(detectFormula('bh')).toBe('unit_only')
  })

  it('detects area_pl for m2', () => {
    expect(detectFormula('m2')).toBe('area_pl')
    expect(detectFormula('m²')).toBe('area_pl')
  })

  it('detects box for m3', () => {
    expect(detectFormula('m3')).toBe('box')
    expect(detectFormula('m³')).toBe('box')
  })

  it("detects linear_p for m' and m", () => {
    expect(detectFormula("m'")).toBe('linear_p')
    expect(detectFormula('m')).toBe('linear_p')
    expect(detectFormula('meter')).toBe('linear_p')
  })

  it('detects weight_rebar for kg and ton', () => {
    expect(detectFormula('kg')).toBe('weight_rebar')
    expect(detectFormula('ton')).toBe('weight_rebar')
  })
})

describe('computeVolume calculations', () => {
  it('calculates cylinder / bore pile formula (1/4 * pi * D^2 * T)', () => {
    // D = 0.4 m, Depth = 6 m, 10 titik
    const res = computeVolume('cylinder', '0.4', '', '6', '10', '0')
    const expectedQty = 0.25 * Math.PI * 0.4 * 0.4 * 6
    expect(Math.abs(res.qty - expectedQty)).toBeLessThan(0.001)
    expect(Math.abs(res.totalQty - expectedQty * 10)).toBeLessThan(0.001)
    expect(res.expression).toContain('¼ × π')
  })

  it('calculates area_pt (P x T x sisi) for walls', () => {
    // 60m x 4m, 2 sisi
    const res = computeVolume('area_pt', '60', '', '4', '2', '0')
    expect(res.qty).toBe(240)
    expect(res.totalQty).toBe(480)
    expect(res.expression).toContain('60 m × 4 m × 2 sisi')
  })

  it('calculates weight_rebar (P x BJ x batang)', () => {
    // 12m, BJ 1.042 kg/m', 20 batang
    const res = computeVolume('weight_rebar', '12', '1.042', '', '20', '0')
    expect(res.qty).toBeCloseTo(12.504, 3)
    expect(res.totalQty).toBeCloseTo(250.08, 2)
  })

  it('calculates unit_only (lot / pcs)', () => {
    const res = computeVolume('unit_only', '', '', '', '5', '0')
    expect(res.qty).toBe(1)
    expect(res.totalQty).toBe(5)
  })

  it('calculates box (P x L x T x unit)', () => {
    const res = computeVolume('box', '10', '2', '3', '2.5', '0')
    expect(res.qty).toBe(60)
    expect(res.totalQty).toBe(150)
  })
})

describe('renderDimensions', () => {
  it('formats cylinder dimensions properly', () => {
    const text = renderDimensions({ formula: 'cylinder', length: '0.4000', height: '6.0000' })
    expect(text).toContain('D: 0,4 m, T: 6 m')
  })

  it('formats rebar dimensions properly', () => {
    const text = renderDimensions({ formula: 'weight_rebar', length: '12.0000', width: '1.0420' })
    expect(text).toContain("P: 12 m, BJ: 1,04 kg/m'")
  })

  it('returns - for unit_only', () => {
    const text = renderDimensions({ formula: 'unit_only' })
    expect(text).toBe('-')
  })
})
