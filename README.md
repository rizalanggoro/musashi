# MUSASHI

ERP monitoring proyek konstruksi. PRD: [`docs/PRD.md`](docs/PRD.md). Kontrak API: [`openapi.yaml`](openapi.yaml).

```text
frontend/   Vite + React + TS, shadcn/ui, Tailwind v4, TanStack Router + Query, heyapi client
backend/    Go 1.26 + Fiber v3, GORM, golang-migrate, PostgreSQL
```

## Menjalankan

```bash
# 1. database
docker compose up -d db

# 2. backend
cd backend
go run ./cmd/migrate        # terapkan migrasi
go run ./cmd/api            # http://localhost:8080/api/v1/health

# 3. frontend
cd frontend
npm install
npm run dev                 # http://localhost:5173 (proxy /api → :8080)
```

Konfigurasi backend lewat env: lihat `backend/.env.example`.

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
