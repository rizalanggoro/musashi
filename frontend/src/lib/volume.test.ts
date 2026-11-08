import { describe, it, expect } from 'bun:test'
import {
  detectFormula,
  computeVolume,
  renderDimensions,
} from './volume'

describe('detectFormula', () => {
  it('detects unit_only for lot, ls, pcs, unit', () => {
    expect(detectFormula('lot')).toBe('unit_only')
    expect(detectFormula('LS')).toBe('unit_only')
    expect(detectFormula('pcs')).toBe('unit_only')
    expect(detectFormula('unit')).toBe('unit_only')
    expect(detectFormula('titik')).toBe('unit_only')
  })

  it('detects area_pl for m2', () => {
    expect(detectFormula('m2')).toBe('area_pl')
    expect(detectFormula('M2')).toBe('area_pl')
  })

  it('detects box for m3', () => {
    expect(detectFormula('m3')).toBe('box')
  })

  it("detects linear_p for m' and m", () => {
    expect(detectFormula("m'")).toBe('linear_p')
    expect(detectFormula('m')).toBe('linear_p')
  })

  it('detects weight_rebar for kg and ton', () => {
    expect(detectFormula('kg')).toBe('weight_rebar')
    expect(detectFormula('ton')).toBe('weight_rebar')
  })
})

describe('computeVolume calculations', () => {
  it('calculates cylinder / bore pile formula (1/4 * pi * D^2 * T)', () => {
    // D = 0.4m, T = 6m, unit = 5 titik
    const res = computeVolume('cylinder', '0.4', '', '6', '5', '0')
    const expectedQtyPerUnit = 0.25 * Math.PI * (0.4 * 0.4) * 6
    expect(Math.abs(res.qty - expectedQtyPerUnit)).toBeLessThan(0.0001)
    expect(Math.abs(res.totalQty - expectedQtyPerUnit * 5)).toBeLessThan(0.0001)
    expect(res.expression).toContain('0,4')
    expect(res.expression).toContain('6')
  })

  it('calculates area_pt (P x T x sisi) for walls', () => {
    // P = 10, T = 3, sisi = 2
    const res = computeVolume('area_pt', '10', '', '3', '2', '0')
    expect(res.qty).toBe(30)
    expect(res.totalQty).toBe(60)
  })

  it('calculates weight_rebar (P x BJ x batang)', () => {
    // P = 12m, BJ = 1.58 kg/m' (D16), unit = 100 batang
    const res = computeVolume('weight_rebar', '12', '1.58', '', '100', '0')
    expect(Math.abs(res.qty - 12 * 1.58)).toBeLessThan(0.0001)
    expect(Math.abs(res.totalQty - 12 * 1.58 * 100)).toBeLessThan(0.0001)
  })

  it('calculates unit_only (lot / pcs)', () => {
    const res = computeVolume('unit_only', '', '', '', '1', '0')
    expect(res.qty).toBe(1)
    expect(res.totalQty).toBe(1)
  })

  it('calculates box (P x L x T x unit)', () => {
    const res = computeVolume('box', '2', '3', '4', '2', '0')
    expect(res.qty).toBe(24)
    expect(res.totalQty).toBe(48)
  })
})

describe('renderDimensions', () => {
  it('formats cylinder dimensions properly', () => {
    expect(renderDimensions({ formula: 'cylinder', length: '0.4', height: '6' })).toBe('D: 0,4 m, T: 6 m')
  })

  it('formats rebar dimensions properly', () => {
    expect(renderDimensions({ formula: 'weight_rebar', length: '12', width: '1.58' })).toBe("P: 12 m, BJ: 1,58 kg/m'")
  })

  it('returns - for unit_only', () => {
    expect(renderDimensions({ formula: 'unit_only' })).toBe('-')
  })
})

