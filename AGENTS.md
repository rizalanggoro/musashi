# AGENTS.md

## Git (RULE — user-mandated)

- **Jangan pernah commit atau push kecuali user memintanya secara eksplisit.**
  "Build hijau", "sudah selesai", atau kebiasaan commit sebelumnya BUKAN izin.
- Saat commit diminta: pesan dalam **bahasa Inggris**, format **Conventional Commits**
  (`feat:`, `fix:`, `refactor:`, `chore:`, `test:`, `docs:`, `perf:`, `build:` …),
  subject line singkat, huruf kecil, tanpa titik di akhir.
- Jangan pernah `git push` tanpa permintaan eksplisit, termasuk setelah commit.

## Database (RULE — user-mandated)

- **Jangan gunakan Docker untuk database.** User menggunakan **DBngin** untuk menjalankan PostgreSQL secara lokal.
- Database lokal PostgreSQL berjalan pada port default (`localhost:5432`) dengan nama database `musashi`.
- Konfigurasi koneksi PostgreSQL dikelola melalui `backend/.env` (`DATABASE_URL`).

