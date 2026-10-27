# MUSASHI

ERP monitoring proyek konstruksi. PRD: [`docs/PRD.md`](docs/PRD.md). Kontrak API: [`openapi.yaml`](openapi.yaml).

```text
frontend/   Vite + React + TS, shadcn/ui, Tailwind v4, TanStack Router + Query, heyapi client
backend/    Go 1.26 + Fiber v3, GORM, golang-migrate, PostgreSQL (DBngin / local)
```

## Menjalankan

### 1. Database (DBngin)
1. Buka aplikasi **DBngin** dan nyalakan service **PostgreSQL** (default port `5432`).
2. Buat database bernama `musashi` (misalnya via TablePlus, pgAdmin, atau DBngin CLI).
3. Buat file `backend/.env` dari salinan `backend/.env.example`:
   ```bash
   cp backend/.env.example backend/.env
   ```
   Sesuaikan `DATABASE_URL` jika port atau password PostgreSQL di DBngin Anda berbeda (misal tanpa password: `postgres://postgres@localhost:5432/musashi?sslmode=disable`).

### 2. Backend
```bash
cd backend
go run ./cmd/migrate        # terapkan migrasi tabel dan seed admin
go run ./cmd/api            # http://localhost:8080/api/v1/health
```

### 3. Frontend
```bash
cd frontend
npm install
npm run dev                 # http://localhost:5173 (proxy /api → :8080)
```

Default login Administrator:
- **Email:** `admin@musashi.local`
- **Password:** `admin123`

---

## Perubahan API

`openapi.yaml` adalah sumber kontrak. Setelah mengubahnya:

```bash
cd frontend && npm run api     # regenerate src/api (client + TanStack Query hooks)
```

Jangan edit file di `frontend/src/api/` — semuanya hasil generate.

## Build

```bash
cd backend  && go build ./...
cd frontend && npm run build
```
