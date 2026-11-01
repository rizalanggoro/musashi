export const UNIT_PRICE_CATEGORIES = [
  "ALL",
  "LABOR",
  "MATERIALS",
  "EQUIPMENT",
  "OTHER",
] as const;

export const UNIT_PRICE_FORM_CATEGORIES = [
  "LABOR",
  "MATERIALS",
  "EQUIPMENT",
  "OTHER",
] as const;

export type UnitPriceCategory = (typeof UNIT_PRICE_CATEGORIES)[number];
export type UnitPriceFormCategory = (typeof UNIT_PRICE_FORM_CATEGORIES)[number];

