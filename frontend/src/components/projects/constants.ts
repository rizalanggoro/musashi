export const PROJECT_SECTIONS = [
  "GA Building",
  "EHS",
  "Wellfare",
] as const;

export type ProjectSection = (typeof PROJECT_SECTIONS)[number];

export const PROJECT_PLANTS = [
  "Cikarang",
  "Karawang",
] as const;

export type ProjectPlant = (typeof PROJECT_PLANTS)[number];

export type ProgressKey =
  | "progress_ppt"
  | "progress_drawing"
  | "progress_rab"
  | "progress_vendor"
  | "progress_spb"
  | "progress_pr"
  | "progress_ringi"
  | "progress_po"
  | "progress_construction"
  | "progress_bast";

export interface StageOption {
  key: ProgressKey;
  value: string;
  label: string;
}

export const PROGRESS_STAGES: StageOption[] = [
  { key: "progress_ppt", value: "progress_ppt", label: "PPT" },
  { key: "progress_drawing", value: "progress_drawing", label: "Drawing" },
  { key: "progress_rab", value: "progress_rab", label: "RAB" },
  { key: "progress_vendor", value: "progress_vendor", label: "Vendor" },
  { key: "progress_spb", value: "progress_spb", label: "SPB" },
  { key: "progress_pr", value: "progress_pr", label: "PR" },
  { key: "progress_ringi", value: "progress_ringi", label: "RINGI" },
  { key: "progress_po", value: "progress_po", label: "PO" },
  { key: "progress_construction", value: "progress_construction", label: "Construction" },
  { key: "progress_bast", value: "progress_bast", label: "BAST" },
];

