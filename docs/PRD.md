# PRD: MUSASHI — Project Monitoring & Cost Estimate

**Versi:** 2.0 (revisi mengikuti mockup `docs/musashi.pdf`)
**Sumber kebenaran desain:** 10 halaman mockup `docs/musashi.pdf`
**Status:** Draft untuk implementasi MVP

---

## 1. Ringkasan

MUSASHI adalah aplikasi web internal untuk memonitor proyek konstruksi dan menyusun
estimasi biaya. Cakupan mengikuti mockup: data proyek, jadwal, harga satuan, AHSP,
BOQ, volume (take-off), dan RAB, dengan keluaran Excel/PDF untuk Volume dan RAB.

Bukan ERP penuh. Di luar cakupan: accounting, procurement, HR, inventory, payroll.

### Tujuan

1. Satu tempat untuk data proyek + budget + progress.
2. Harga satuan dan AHSP menjadi satu sumber harga untuk RAB.
3. Volume pekerjaan dihitung dari dimensi (take-off), bukan diketik sebagai angka jadi.
4. RAB dan Volume bisa diekspor (Excel / PDF) untuk kebutuhan internal.

### Prinsip

- **Mockup adalah kontrak UI.** Label, kolom, dan susunan form mengikuti PDF apa adanya.
  Penyimpangan hanya boleh karena alasan teknis, dan dicatat di §14.
- Perhitungan uang terjadi di backend, frontend hanya menampilkan.
- Satu pola halaman untuk semua modul (judul + search + add + tabel).

---

## 2. Inventaris Layar (dari PDF)

| Hal. | Layar | Isi |
|---|---|---|
| 1 | Shell / navigasi | Topbar (hamburger + logo MUSASHI), sidebar 2 grup |
| 2 | Project (list) | Tabel proyek + Search + Add |
| 3 | Project (form) | Modal form proyek |
| 4 | Unit Price (list) | Tabel harga satuan + Search + Add |
| 5 | Unit Price (form) | Modal form harga satuan |
| 6 | AHSP (list) | Tabel AHSP + Search + Add |
| 7 | AHSP (form) | Form AHSP + tabel komponen A–F + Total Cost |
| 8 | BOQ | Tabel proyek + status/export Volume & RAB + tombol Volume/RAB |
| 9 | Volume | Pilih proyek + tabel take-off + Add |
| 10 | RAB | Pilih proyek + tabel RAB + Add |

Tidak ada layar Dashboard, Login, Schedule detail, atau Cost Estimate detail di mockup.
Schedule hanya muncul sebagai menu; Cost Estimate hanya sebagai judul grup menu.

---

## 3. Navigasi & IA

```text
Topbar: [hamburger] [logo MUSASHI]

Sidebar
├── MONITORING            (label grup, bukan link)
│   ├── PROJECT
│   └── SCHEDULE
└── COST ESTIMATE         (label grup, bukan link)
    ├── UNIT PRICE
    ├── AHSP
    └── BOQ ──┬── VOLUME   (tombol di halaman BOQ)
              └── RAB      (tombol di halaman BOQ)
```

Aturan:

- Item aktif: background pink (`--accent`) dengan bar merah di kiri. Label grup (`MONITORING`,
  `COST ESTIMATE`) selalu merah penuh, teks putih. Item non-aktif: teks hitam di atas putih,
  garis pemisah abu muda.
- Sidebar bisa di-collapse lewat hamburger (mockup hanya menampilkan ikonnya).
- `VOLUME` dan `RAB` **wajib** dipilih proyek dulu (lewat halaman BOQ atau dropdown proyek),
  karena keduanya selalu ter-scope ke satu proyek.

---

## 4. Design System

Diambil dari mockup; dipetakan ke shadcn + Tailwind theme.

| Token | Nilai | Pemakaian |
|---|---|---|
| `--primary` | `#E4000F` (merah MUSASHI) | tombol Add, chip judul halaman, item sidebar aktif |
| `--primary-foreground` | `#FFFFFF` | teks di atas merah |
| `--muted` / `--accent` | `#FADCDC` (pink) | isi sel tabel, input search |
| `--secondary` | `#D9D9D9` (abu) | topbar |
| `--border` | `#000000` / `#E5E7EB` | border tabel (hitam tipis), garis sidebar (abu) |
| Radius | `rounded-2xl` untuk card/modal, `rounded-none` untuk tabel | — |
| Font | sans-serif geometric (mockup mirip Poppins) | seluruh UI |

### Pola halaman (semua modul)

```text
┌───────────────────────────────────────────────────────────────┐
│ [ JUDUL MODUL ]                    [ 🔍 SEARCH ] [   ADD   ] │
├───────────────────────────────────────────────────────────────┤
│  tabel berborder, header rata tengah, body pink               │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

- Judul modul: chip merah, teks putih, huruf kapital, lebar mengikuti isi.
- Search: input pink dengan ikon kaca pembesar, placeholder `SEARCH`, debounce 300 ms,
  mencari di kolom teks utama modul.
- Add: tombol merah, membuka modal form (bukan halaman baru).
- Form/modal: card putih rounded, label kapital di atas input, tombol `ADD` di kanan bawah.
- Tabel: header kapital rata tengah, baris body pink `#FADCDC`, border hitam 1px,
  kolom `NO` sebagai penomoran urut (bukan ID database).
- Angka uang rata kanan, format `Rp1.234.567` (id-ID). Desimal volume/koefisien pakai koma.

---

## 5. Peran Pengguna

Mockup tidak menampilkan login maupun role. Untuk MVP tetap dipakai 2 role (kebutuhan operasional,
dicatat sebagai tambahan di §14).

| Aksi | Admin | Staff |
|---|---|---|
| Lihat semua modul & export | ✅ | ✅ |
| CRUD Project | ✅ | ❌ |
| CRUD Unit Price | ✅ | ❌ |
| CRUD AHSP | ✅ | ❌ |
| Input Volume (take-off) | ✅ | ✅ |
| Input/edit RAB | ✅ | ❌ |
| CRUD Schedule | ✅ | ❌ |

Autentikasi: session cookie (httpOnly) + tabel `users`. Tanpa registrasi publik; user dibuat via seed/CLI.

---

## 6. Modul

### 6.1 Project

**List (hal. 2)**

| Kolom | Tipe | Catatan |
|---|---|---|
| NO | int | urut tampil |
| PROJECT | text | nama proyek |
| SECTION | text | mis. struktur/arsitektur/unit kerja |
| PLANT | text | lokasi/pabrik/site |
| IMPLEMENTATION PLAN | date | tanggal rencana pelaksanaan (lihat §15) |
| BUDGET | money | pagu anggaran |
| PROGRESS — SPB | status | status tahap SPB |
| PROGRESS — RAB | status | status tahap RAB |
| PROGRESS — XXX | status | kolom ketiga belum didefinisikan (lihat §15) |
| INFORMATION | text | keterangan bebas |

**Form (hal. 3)**

```text
PROJECT           [ input ]
SECTION           [ input ]        PLANT  [ input ]
DATE              [ date  ]        BUDGET [ money ]
PROGRESS          SPB [ ] RAB [ ] XXX [ ]
                                                    [ ADD ]
```

**Aturan**

- `PROJECT` wajib, unik (case-insensitive).
- `BUDGET` ≥ 0.
- Checkbox PROGRESS menyimpan status per tahap (`done` bila dicentang, `pending` bila tidak).
- Klik baris membuka detail proyek: ringkasan + tautan ke Volume/RAB proyek tersebut.
- Hapus proyek = soft delete (`deleted_at`), agar riwayat RAB tidak hilang.

### 6.2 Schedule

Menu ada di mockup, layarnya tidak. Implementasi minimal, mengikuti pola halaman yang sama.

| Kolom | Tipe |
|---|---|
| NO | int |
| CODE | text |
| WORK ITEM | text |
| START DATE | date |
| END DATE | date |
| DURATION | int (hari, dihitung) |
| WEIGHT | decimal (%) |
| PROGRESS | decimal (%) |
| STATUS | enum: `planning`, `on_going`, `done`, `late` |

Aturan: `DURATION = end_date - start_date + 1`. `PROGRESS` manual (belum ada sumber volume
per item jadwal). Σ `WEIGHT` per proyek disarankan 100% (divalidasi sebagai warning, bukan blokir).

### 6.3 Unit Price

**List (hal. 4)**

| Kolom | Tipe |
|---|---|
| NO | int |
| CODE | text, unik |
| CATEGORY | enum: `LABOR`, `MATERIALS`, `EQUIPMENT`, `OTHER` |
| ITEM | text (label mockup: ITEM; form: WORK ITEM) |
| UNIT | text (m³, kg, OH, ls, …) |
| UNIT PRICE | money |
| INFORMATION | text |

**Form (hal. 5):** `CODE`, `CATEGORY`, `WORK ITEM`, `UNIT`, `UNIT PRICE`, `INFORMATION`.

Aturan:

- `CODE` unik, format bebas tapi disarankan `MAT-001` / `LAB-001` / `EQP-001`.
- Perubahan `UNIT PRICE` **tidak** mengubah RAB yang sudah ada (RAB menyimpan harga saat dipakai).
- Hapus ditolak bila harga dipakai komponen AHSP (409 + pesan).

### 6.4 AHSP (Analisa Harga Satuan Pekerjaan)

**List (hal. 6):** `CODE`, `ITEM`, `UNIT`, `COST`, `STATUS`.
`COST` = `UNIT PRICE OF WORK` (F) hasil perhitungan; `STATUS` = `draft` | `final`.

**Form (hal. 7)**

```text
CODE [ ]   WORK ITEM [ ]   UNIT [ ]

NO  DESCRIPTION                    CODE  UNIT  COEFFICIENT  UNIT PRICE  COST
A   LABOR
      <baris komponen>
    TOTAL LABOR COST
B   MATERIALS
      <baris komponen>
    TOTAL MATERIAL COST
C   EQUIPMENT
      <baris komponen>
    TOTAL EQUIPMENT COST
D   TOTAL COST
E   OVERHEAD + PROFIT        (persen, mis. 10%)
F   UNIT PRICE OF WORK

TOTAL COST [ ]
```

Aturan & formula:

```text
A. TOTAL LABOR COST       = Σ (coefficient × unit_price) grup LABOR
B. TOTAL MATERIAL COST    = Σ (coefficient × unit_price) grup MATERIALS
C. TOTAL EQUIPMENT COST   = Σ (coefficient × unit_price) grup EQUIPMENT
D. TOTAL COST (D)         = A + B + C
E. OVERHEAD + PROFIT      = overhead_profit_pct × D          (default 10%)
F. UNIT PRICE OF WORK     = D + E
```

- `COST` per baris = `coefficient × unit_price`, dibulatkan ke rupiah (half-up).
- Baris komponen memilih `Unit Price` (master) atau input manual (`code`, `unit`, `unit_price`).
  Bila memilih master, `code`/`unit`/harga terisi otomatis dan tetap bisa di-override.
- `TOTAL COST` di luar tabel = `F` (harga satuan pekerjaan), ditampilkan read-only.
- Perubahan koefisien/harga langsung menghitung ulang (optimistic di UI, angka final dari backend).
- AHSP tidak terikat proyek (master), dipakai ulang lintas proyek.

### 6.5 BOQ (halaman hub per proyek)

**Layar (hal. 8)**

```text
[ BOQ ]                        [ VOLUME ] [ RAB ]  [ 🔍 SEARCH ]

NO | PROJECT | VOLUME: STATUS | EXCEL | RAB: STATUS | EXCEL/PDF
```

Aturan:

- BOQ **tidak** menyimpan item baru di layar ini. Isinya daftar proyek + status kelengkapan:
  - `VOLUME / STATUS`: `empty` (belum ada item volume), `partial`, `complete`
  - `RAB / STATUS`: `empty`, `partial`, `complete`
- `VOLUME / EXCEL` dan `RAB / EXCEL|PDF` = tombol unduh. Nonaktif bila status `empty`.
- Tombol `VOLUME` / `RAB` mengarah ke layar modul dengan proyek terpilih
  (`/boq/{projectId}/volume`, `/boq/{projectId}/rab`).
- Sumber item BOQ = item pekerjaan yang dibuat di layar Volume (work item + volume rencana).

### 6.6 Volume (take-off)

**Layar (hal. 9)**

```text
[ VOLUME ]                                   [ 🔍 SEARCH ] [ ADD ]
PROJECT [ dropdown / searchable select ]

NO | WORK ITEM | VOLUME | UNIT | LENGTH (m) | WIDTH (m) | HEIGHT (m) | QUANTITY | UNIT WEIGHT (kg/m'; kg/m³) | TOTAL QUANTITY
```

Aturan & formula:

```text
QUANTITY        = LENGTH × WIDTH × HEIGHT            (bila ketiga dimensi diisi)
                = VOLUME                             (bila dimensi kosong)
TOTAL QUANTITY  = QUANTITY × UNIT WEIGHT             (default UNIT WEIGHT = 1)
```

- `WORK ITEM` dipilih dari AHSP (membawa `UNIT`) atau diketik manual.
- `VOLUME` = volume rencana item tersebut (dari BOQ/desain).
- Dimensi disimpan per baris; satu work item boleh punya banyak baris take-off
  (mis. kolom K1 di 3 zona) → `QUANTITY` total item = Σ baris.
- `UNIT WEIGHT` dipakai untuk pekerjaan berbobot (besi `kg/m'`, beton `kg/m³`). Kosong = 1.
- `UNIT` mengikuti satuan work item (m³, m², m', kg, ls). Validasi ringan: bila `UNIT = kg`
  dan `UNIT WEIGHT` diisi, tampilkan hint satuan.
- Perubahan volume mengubah `VOLUME/STATUS` di halaman BOQ dan angka RAB (harga tidak berubah).

### 6.7 RAB

**Layar (hal. 10)**

```text
[ RAB ]                                      [ 🔍 SEARCH ] [ ADD ]
PROJECT [ dropdown / searchable select ]

NO | WORK ITEM | VOLUME | UNIT | UNIT PRICE | TOTAL COST
```

Aturan & formula:

```text
TOTAL COST = VOLUME × UNIT PRICE        (dibulatkan ke rupiah)
TOTAL RAB  = Σ TOTAL COST
```

- `WORK ITEM` berasal dari item Volume proyek tersebut; `VOLUME` mengikuti volume aktual take-off.
- `UNIT PRICE` default = `UNIT PRICE OF WORK` (F) dari AHSP yang terhubung; boleh override
  (harga negosiasi). Nilai yang dipakai **disimpan** di baris RAB, bukan dihitung ulang saat dibaca.
- `ADD` = menambah baris RAB: pilih work item → volume & harga terisi → simpan.
- Baris tanpa harga AHSP ditandai `⚠ belum ada harga` dan dihitung 0 pada total, dengan banner
  peringatan di atas tabel.
- Footer tabel: `TOTAL RAB` (dan opsional subtotal per grup pekerjaan bila item punya grup).

### 6.8 Cost Estimate (grup menu)

Tidak ada layar di mockup; diturunkan dari data yang ada. Ditampilkan sebagai ringkasan
di halaman detail proyek:

```text
BUDGET          Rp1.000.000.000
RAB (rencana)   Rp  950.000.000
SPB (terbayar)  Rp  700.000.000
VARIANCE        Rp   50.000.000      (BUDGET - RAB)
SISA            Rp  300.000.000      (BUDGET - SPB)
```

`VARIANCE` negatif = RAB melebihi budget → badge merah.

---

## 7. Model Data

PostgreSQL. Nama tabel snake_case, plural. Semua tabel: `id uuid`, `created_at`, `updated_at`.
Uang disimpan `bigint` rupiah (tanpa sen). Koefisien/dimensi `numeric(14,4)`.

```text
users(id, email, name, password_hash, role, created_at, updated_at)

projects(
  id, name, section, plant, implementation_plan date, budget bigint,
  information text, status enum, deleted_at,
  progress_spb enum, progress_rab enum, progress_xxx enum
)

schedules(
  id, project_id, code, work_item, start_date, end_date, duration int,
  weight numeric, progress numeric, status enum
)

unit_prices(
  id, code unique, category enum, work_item, unit, unit_price bigint, information
)

ahsps(
  id, code unique, work_item, unit, overhead_profit_pct numeric(6,2) default 10,
  total_labor bigint, total_material bigint, total_equipment bigint,
  total_cost bigint, unit_price_of_work bigint, status enum
)

ahsp_items(
  id, ahsp_id, group enum('LABOR','MATERIALS','EQUIPMENT'), sort int,
  description, code, unit, coefficient numeric(14,4), unit_price bigint,
  cost bigint, unit_price_id nullable
)

boq_items(                      -- item pekerjaan proyek (dibuat dari layar Volume)
  id, project_id, ahsp_id nullable, work_item, unit, planned_volume numeric(14,4)
)

volumes(                        -- baris take-off
  id, boq_item_id, date, length numeric, width numeric, height numeric,
  quantity numeric, unit_weight numeric default 1, total_quantity numeric,
  description
)

rab_items(
  id, project_id, boq_item_id, work_item, unit,
  volume numeric(14,4), unit_price bigint, total_cost bigint
)
```

Relasi:

```text
Project 1─N Schedule
Project 1─N BoqItem 1─N Volume
Project 1─N RabItem N─1 BoqItem
Ahsp    1─N AhspItem N─1 UnitPrice (opsional)
BoqItem N─1 Ahsp (opsional, sumber harga)
```

Enum:

```text
project.status         : planning | on_going | completed | cancelled
project.progress_*     : pending | done
ahsp.status            : draft | final
unit_price.category    : LABOR | MATERIALS | EQUIPMENT | OTHER
schedule.status        : planning | on_going | done | late
```

Index minimal: `projects(deleted_at)`, `boq_items(project_id)`, `volumes(boq_item_id)`,
`rab_items(project_id)`, `ahsps(code)`, `unit_prices(code)`, `schedules(project_id)`.

---

## 8. Tech Stack

### Frontend

| Bagian | Pilihan |
|---|---|
| Build | Vite + React + TypeScript |
| UI | shadcn/ui + Tailwind CSS |
| Routing | TanStack Router (file-based, satu layout shell) |
| Data | TanStack Query |
| API client | `@hey-api/openapi-ts` generate dari `openapi.yaml` (client + TanStack Query hooks) |
| Form | react-hook-form + zod (schema dari generator bila tersedia) |
| Tabel | TanStack Table (header group untuk kolom PROGRESS/Volume/RAB) |
| Export | unduh file dari endpoint backend (bukan generate di browser) |

Aturan FE:

- Semua request lewat client hasil generate; tidak ada `fetch` manual di komponen.
- Query key dari helper hasil generate; invalidasi per entitas setelah mutation.
- Uang & desimal dikirim sebagai string dari API, diformat dengan `Intl.NumberFormat('id-ID')`.
- Tidak ada state global tambahan (Redux/Zustand) — cukup TanStack Query + local state.

### Backend

| Bagian | Pilihan |
|---|---|
| Bahasa | Go (≥ 1.26) |
| HTTP | Fiber v3 (middleware bawaan: requestid, recover, logger, cors) |
| ORM | GORM |
| Migrasi | `golang-migrate` (file SQL berurutan di `db/migrations`) |
| DB | PostgreSQL |
| API spec | `openapi.yaml` (OpenAPI 3.1) sebagai sumber kontrak |
| Auth | session cookie (httpOnly) + bcrypt |
| Export | Excel: `excelize`; PDF: template HTML → `chromium`/`wkhtmltopdf` (server-side) |

Aturan BE:

- Perhitungan uang di service layer, bukan di handler/DB trigger.
- GORM hanya untuk CRUD; query laporan pakai `db.Raw`/CTE agar tidak N+1.
- Migrasi tidak boleh dijalankan otomatis saat start di production (job terpisah).

### Struktur repo

```text
musashi/
├── docs/PRD.md
├── openapi.yaml
├── backend/
│   ├── cmd/api/main.go
│   ├── db/migrations/          # golang-migrate
│   ├── internal/
│   │   ├── auth/
│   │   ├── project/
│   │   ├── schedule/
│   │   ├── unitprice/
│   │   ├── ahsp/
│   │   ├── boq/
│   │   ├── volume/
│   │   ├── rab/
│   │   ├── export/
│   │   └── platform/           # db, config, middleware, httperr
│   └── go.mod
└── frontend/
    ├── src/routes/             # TanStack Router
    ├── src/components/ui/      # shadcn
    ├── src/api/                # hasil generate heyapi (git-ignored/committed? lihat §14)
    └── vite.config.ts
```

---

## 9. API

Kontrak ditulis di `openapi.yaml`, lalu:

```bash
npx @hey-api/openapi-ts -i openapi.yaml -o frontend/src/api -c @hey-api/client-fetch
```

Endpoint MVP (semua di bawah `/api/v1`, JSON):

```text
POST   /auth/login
POST   /auth/logout
GET    /auth/me

GET    /projects?q=&page=&per_page=
POST   /projects
GET    /projects/{id}
PUT    /projects/{id}
DELETE /projects/{id}

GET    /projects/{id}/schedules
POST   /projects/{id}/schedules
PUT    /schedules/{id}
DELETE /schedules/{id}

GET    /unit-prices?q=&category=
POST   /unit-prices
PUT    /unit-prices/{id}
DELETE /unit-prices/{id}

GET    /ahsps?q=
POST   /ahsps
GET    /ahsps/{id}
PUT    /ahsps/{id}
DELETE /ahsps/{id}

GET    /boq?q=                                  # daftar proyek + status volume/rab
GET    /projects/{id}/boq                       # item BOQ + volume + rab per item

GET    /projects/{id}/volumes
POST   /projects/{id}/volumes
PUT    /volumes/{id}
DELETE /volumes/{id}

GET    /projects/{id}/rab
POST   /projects/{id}/rab
PUT    /rab-items/{id}
DELETE /rab-items/{id}

GET    /projects/{id}/cost-estimate

GET    /projects/{id}/exports/volume.xlsx
GET    /projects/{id}/exports/rab.xlsx
GET    /projects/{id}/exports/rab.pdf
```

Konvensi:

- Error: `{ "error": { "code": "...", "message": "...", "fields": {...} } }`.
- List: `{ "data": [...], "meta": { "page": 1, "per_page": 25, "total": 120 } }`.
- Uang/desimal: `string` di JSON (`"1234567"`, `"0.4500"`).
- Tanggal: `YYYY-MM-DD`. Timestamp: RFC3339 UTC.
- 409 untuk pelanggaran unique/FK (mis. hapus unit price yang dipakai AHSP).

---

## 10. Export

| Export | Isi | Format |
|---|---|---|
| Volume | Header proyek + tabel take-off (work item, volume, unit, L, W, H, quantity, unit weight, total quantity) | `.xlsx` |
| RAB | Header proyek + tabel RAB + baris TOTAL RAB | `.xlsx` dan `.pdf` |

- Nama file: `Volume-{project}-{YYYYMMDD}.xlsx`, `RAB-{project}-{YYYYMMDD}.{ext}`.
- Angka diformat `#,##0` dan kolom uang `Rp#,##0`.
- PDF dibuat server-side dari template HTML yang sama dengan tampilan tabel (A4 landscape).
- Export mengikuti filter search yang aktif bila query dikirim (default: seluruh item proyek).

---

## 11. Non-Functional

- **Performa:** p95 API < 300 ms untuk data ≤ 500 item/proyek; halaman list < 1.000 baris per request
  (server-side pagination). Bundle awal < 2 MB compressed.
- **Kapasitas:** 50 concurrent user internal.
- **Keamanan:** session httpOnly + SameSite=Lax, CSRF token untuk mutation, password bcrypt,
  semua endpoint (kecuali login) butuh sesi. Audit ringan: `created_at/updated_at` + `updated_by`.
- **Integritas angka:** semua perhitungan uang integer rupiah; pembulatan half-up di service layer;
  perubahan harga master tidak mengubah dokumen historis.
- **Aksesibilitas:** semua form punya `<label>` terhubung, navigasi keyboard di modal, kontras
  teks putih di atas merah `#E4000F` ≥ 4.5:1 (pakai merah lebih gelap `#C4000D` untuk teks kecil bila gagal audit).
- **Deployment:**

```text
Caddy
 ├── /            → static React (Vite build)
 └── /api/*       → Go (port 8080)
                       └── PostgreSQL
```

Tanpa Node runtime di production.

---

## 12. Prioritas Implementasi

```text
1. Shell + design system (sidebar, topbar, pola halaman, tabel, modal)
2. Auth + users (seed admin)
3. Project (list + form + detail)
4. Unit Price
5. AHSP (komponen A–F + perhitungan)
6. Volume (take-off + dimensi) + BoqItem
7. RAB (harga dari AHSP, simpan snapshot)
8. BOQ hub (status + export Excel/PDF)
9. Schedule
10. Cost Estimate ringkasan
```

**Alur inti:**

```text
Project → Unit Price → AHSP → Volume (BoqItem) → RAB → BOQ (status + export)
```

---

## 13. Acceptance Criteria (ringkas)

- [ ] Sidebar & halaman 1:1 dengan mockup (label, warna, urutan grup).
- [ ] Semua tabel memakai kolom dan urutan kolom sesuai §6.
- [ ] Search di setiap list memfilter server-side (bukan hanya baris yang sudah dimuat).
- [ ] Form Project menyimpan 3 status progress (SPB/RAB/XXX) dan tampil di list.
- [ ] AHSP: `F = D + E` benar untuk minimal 3 kasus uji (labor only, material only, campuran).
- [ ] Volume: `QUANTITY` dari L×W×H dan `TOTAL QUANTITY = QUANTITY × UNIT WEIGHT` benar.
- [ ] RAB: `TOTAL COST = VOLUME × UNIT PRICE`, harga tersimpan saat simpan (tidak ikut berubah
      saat master harga diedit).
- [ ] Export Volume `.xlsx` dan RAB `.xlsx`/`.pdf` terbuka tanpa error dan angka sama dengan layar.
- [ ] Hapus Unit Price yang dipakai AHSP → 409 dengan pesan jelas.

---

## 14. Tambahan di Luar Mockup

Hal-hal yang tidak ada di PDF tetapi dibutuhkan agar aplikasi jalan. Semua ditandai agar mudah
dicabut bila tidak diinginkan:

| Tambahan | Alasan |
|---|---|
| Login + 2 role | data biaya proyek tidak boleh terbuka |
| Halaman detail proyek (ringkasan + Cost Estimate) | mockup hanya punya list & form |
| Layar Schedule | menu ada, layar belum digambar |
| `project.status` | dibutuhkan filter/urutan proyek aktif |
| Pagination & soft delete | volume data proyek bertambah tiap tahun |
| Kolom `created_at/updated_at/updated_by` | audit minimal |

---

## 15. Pertanyaan Terbuka

1. **SPB** pada kolom PROGRESS: *Surat Perintah Bayar* / *Sertifikat Progress*? Apakah nilainya
   rupiah (butuh input nominal) atau sekadar status selesai/belum?
2. **XXX** pada kolom PROGRESS dan checkbox form: metrik apa? Placeholder atau kolom ketiga
   yang sudah direncanakan (mis. *Faktur*/*Termin*)?
3. **IMPLEMENTATION PLAN** berupa satu tanggal atau rentang (start–end)? Mockup hanya satu field `DATE`.
4. **SECTION** dan **PLANT**: master data (dropdown) atau teks bebas?
5. **UNIT WEIGHT** `kg/m'; kg/m³`: satu kolom dengan pilihan satuan, atau dua kolom?
   Default per satuan work item?
6. Apakah `VOLUME` di layar Volume = volume rencana (desain) dan `TOTAL QUANTITY` = hasil take-off,
   atau sebaliknya?
7. Export RAB PDF: perlu kop/logo perusahaan dan kolom tanda tangan?
8. Apakah item BOQ perlu hierarki grup pekerjaan (Pekerjaan Persiapan → Struktur → …) atau flat?
