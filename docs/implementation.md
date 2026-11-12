# Implementasi MATA SIDIK (Pilot)

> Dokumen ini menjelaskan struktur, keputusan teknis, cara menjalankan, dan batasan implementasi pilot. Skema database rinci ada di [`schema.md`](schema.md); aturan proyek ada di [`rules.md`](rules.md).

## 1. Keputusan pilot

| Area | Implementasi pilot | Catatan |
| --- | --- | --- |
| Backend | Go + Gin + PASETO | Selaras `docs/rules.md` §4 |
| Database | **SQLite** (`modernc.org/sqlite`) | Diminta pemilik proyek untuk phase awal; produksi wajib PostgreSQL sesuai `docs/rules.md` |
| Frontend | Svelte 5 + SvelteKit + Vite | Build statis, satu origin dengan API |
| UI | shadcn-svelte pattern + neobrutalism.dev | Komponen di `frontend/src/lib/components/ui`; sidebar neobrutalism |
| Auth | **PASETO v4.local + refresh cookie** | Access token 30 menit, refresh HttpOnly 7 hari; RBAC + actingAs admin |
| Arsitektur | Modular monolith + worker | `cmd/api` dan `cmd/worker`; API bisa menjalankan worker inline untuk dev |
| Pola akses data | Repository pattern | Service tidak menulis SQL; implementasi SQLite di `internal/repository/sqlite` |
| File | Volume lokal `data/files` | Nama file = SHA-256 + ekstensi; metadata di `file_object` |

> **Catatan:** `docs/rules.md` §4 menyebut PostgreSQL 16 + PASETO. Pilot memakai SQLite + PASETO atas instruksi pemilik proyek; produksi tetap wajib PostgreSQL dan hardening auth (rotasi kunci, CSRF, MFA approver).

## 2. Struktur proyek

```text
backend/
  cmd/
    api/main.go                 # HTTP server, wiring, graceful shutdown, worker inline
    worker/main.go              # worker background terpisah
  internal/
    config/                     # env + DSN SQLite
    domain/                     # entitas, enum, error, DTO
    repository/
      repository.go             # interface repository
      sqlite/                   # implementasi SQLite + migrasi + seed + testdata
    service/                    # use case: upload, workflow, kalkulasi, dashboard, faktor/TL
    worker/                     # parser CSV/XLSX, classifier, job loop
    httpapi/
      router.go                 # route /api/v1 + SPA fallback
      handler/                  # handler Gin per domain
      middleware/               # request id, auth PASETO, CORS
      response/                 # format response/error konsisten
  testdata/bnba-2026-05.csv     # contoh file upload
  go.mod

frontend/
  src/
    app.css                     # tema neobrutalism + Tailwind
    lib/
      api/client.ts             # API client typed + Bearer PASETO + refresh
      components/
        ui/                     # Button, Card, Badge, Input, Select, Table, Modal, dll.
        layout/                 # Sidebar fixed, UserMenu/actingAs
        dashboard/              # StatCard, MultiBarChart
        upload/                 # PreviewRow
      constants/                # enum, label UI, route
      stores/                   # auth, bootstrap, periode
      types/                    # tipe domain untuk FE
      utils/                    # format, cn, options
    routes/
      login/+page.svelte        # login + diagram alur approval
      +page.svelte              # dashboard SPM & IKU + multi bar chart
      dashboard/spm/[kode]/+page.svelte # detail SPM + provenance + chart per faskes
      upload/+page.svelte       # upload → job → preview → submit
      submissions/+page.svelte
      submissions/[id]/+page.svelte
      master/+page.svelte
      factors/+page.svelte
      audit/+page.svelte
      admin/+page.svelte        # puskesmas, user, jabatan
  svelte.config.js
  vite.config.ts
  tailwind.config.js
  components.json               # konfigurasi shadcn-svelte

assets/diagrams/
  gen-schema.py                 # generator diagram skema database
  schema-db.html / .png         # hasil diagram
docs/
  schema.md                     # penjelasan database lengkap
  implementation.md             # dokumen ini
```

## 3. Cara menjalankan

### 3.1 Backend

```bash
cd backend
go mod tidy
export APP_AUTH_KEY=<64 hex chars = 32 byte>   # opsional; jika kosong API memakai kunci acak sementara
APP_WORKER_INLINE=true APP_ADDR=:8080 go run ./cmd/api
```

API: `http://localhost:8080` · health: `/healthz` · bootstrap: `/api/v1/meta/bootstrap`.

Database dan seed demo dibuat otomatis di `backend/data/matasidik.db`. Seed hanya berjalan jika `S-2026-05-001` belum ada. Reset:

```bash
rm -rf backend/data
```

Worker terpisah (opsional, untuk produksi pilot):

```bash
cd backend
APP_WORKER_INLINE=false go run ./cmd/api
go run ./cmd/worker
```

### 3.2 Frontend

```bash
cd frontend
npm install
npm run dev
```

Vite berjalan di `http://localhost:5173` dan mem-proxy `/api` ke `http://localhost:8080`.

Build produksi:

```bash
cd frontend
npm run build
```

Go menyajikan `frontend/build` melalui env `APP_WEB_DIR` (default `../frontend/build`). SPA fallback ditangani `router.NoRoute`.

### 3.3 Test & pemeriksaan

```bash
cd backend && go test ./...
cd backend && go build ./...
cd frontend && npm run check
cd frontend && npm run build
```

## 4. API pilot

Semua route di bawah `/api/v1`. Autentikasi memakai `Authorization: Bearer <PASETO>`; refresh token dikirim sebagai cookie HttpOnly.

| Method | Path | Fungsi |
| --- | --- | --- |
| POST | `/auth/login` | Login username/password; set refresh cookie |
| POST | `/auth/refresh` | Rotasi refresh + access PASETO baru |
| POST | `/auth/logout` | Revoke sesi + hapus cookie |
| GET | `/auth/me` | User aktif + actingAs |
| POST | `/auth/acting-as` | Admin bertindak sebagai user lain |
| GET | `/profile` | Profil user aktif + actingAs |
| PATCH | `/profile` | Ubah nama tampilan sendiri |
| POST | `/profile/change-password` | Ganti password sendiri; cabut semua sesi |
| GET | `/meta/bootstrap` | Master, periode, label, indikator, target, enum |
| POST | `/submissions` | Upload multipart; balas `202 {submission, job}` |
| GET | `/submissions` | Daftar submission + filter periode/status/faskes |
| GET | `/submissions/{id}` | Detail baris, review, revisi, IKU yang diizinkan |
| PATCH | `/submissions/{id}/rows/{rowId}` | Ubah IKU leaf sebelum submit |
| POST | `/submissions/{id}/submit` | Kirim ke `VERIFY_DINAS` |
| POST | `/submissions/{id}/decision` | `approve`/`return` + catatan |
| GET | `/jobs/{id}` | Polling progress background |
| GET | `/dashboard/spm` | Baris SPM + summary berbasis target |
| GET | `/dashboard/spm/export?format=xlsx\|pdf` | Export seluruh output SPM |
| GET | `/dashboard/spm/{kode}/detail` | Detail faskes, upload coverage, provenance |
| GET | `/dashboard/iku-by-faskes` | Data multi-series chart antar faskes |
| GET | `/dashboard/spm/{kode}/trace` | Trace sampai file asli |
| GET | `/iku/values` | Nilai IKU per periode/scope |
| GET | `/iku/{kode}/values` | Nilai IKU per kode |
| GET | `/iku-values/{id}/detail` | Provenance komponen |
| GET | `/factors` | Nilai IKU + faktor |
| PUT | `/factors/{ikuValueID}` | Upsert faktor (Kabid/Sekretaris/Admin) |
| GET | `/tindak-lanjut` | Daftar TL |
| POST | `/tindak-lanjut` | Buat TL (Pimpinan/Admin) |
| PATCH | `/tindak-lanjut/{id}` | Ubah status/field TL |
| GET | `/audit` | Audit log |
| GET | `/files/{id}` | Unduh file asli |
| GET | `/admin/overview` | Ringkasan user, faskes, jabatan, bidang |
| GET/POST/PATCH | `/admin/users`, `/admin/faskes`, `/admin/positions` | Administrasi user, puskesmas, jabatan |
| GET | `/master/labels` / `/master/indicators` / `/master/targets` | Master |
| POST | `/master/labels` / `/master/indicators` | Tambah master (Admin) |
| PUT | `/master/targets` | Upsert target (Admin/Kabid/Sekretaris) |
| GET | `/master/capaian-tahunan?tahun=2025` | Capaian tahunan per label/IKU/indikator |
| POST | `/master/capaian-tahunan` | Tambah/upsert capaian tahunan (Admin) |
| PATCH | `/master/capaian-tahunan/{id}` | Ubah capaian tahunan (Admin) |
| DELETE | `/master/capaian-tahunan/{id}` | Hapus capaian tahunan (Admin) |

Error konsisten: `{code, message, details, request_id}`.

### 4.1 Auth, actingAs, dan posisi

- Login: `POST /auth/login` dengan `username` + `password`; access token PASETO v4.local (30 menit) di Bearer, refresh token HttpOnly (7 hari) disimpan sebagai hash di `auth_session`.
- Refresh: `POST /auth/refresh` merotasi refresh token dan menerbitkan access token baru.
- `GET /auth/me` mengembalikan user aktif dan `acting_as` bila admin sedang bertindak sebagai user lain.
- Profil mandiri: semua user terautentikasi dapat membuka `/profile`, mengubah nama tampilan, dan mengganti password sendiri. Ganti password wajib memverifikasi password saat ini, minimal 8 karakter, lalu mencabut seluruh refresh session user tersebut.
- Peran, username, jabatan, faskes, dan status aktif tetap hanya dapat diubah admin. Saat actingAs, perubahan profil/password dinonaktifkan agar admin kembali ke akunnya sendiri.
- `POST /auth/acting-as` hanya admin; token target membawa claim `act` berisi admin. Backend memvalidasi `act` adalah admin aktif. Tidak boleh actingAs sesama admin.
- Role efektif diambil dari `positions.role` (fallback `users.role`). Admin dapat mengganti posisi user tanpa mengubah kode workflow; approval berikutnya langsung memakai role baru.
- Faskes dan user adalah entitas terpisah. User Puskesmas punya `faskes_id`; user Dinas tidak. Admin mengelola keduanya di `/admin`.
- Password awal pilot: `password123` (bcrypt cost 12). Wajib diganti sebelum produksi.
- Halaman login menampilkan fungsi aplikasi di kiri dan form login di kanan. Alur empat langkah ditampilkan sebagai satu baris ringkas di kiri (HTML, tanpa gambar) dan tanpa label jabatan. State machine approval lengkap tetap ada di `assets/diagrams/alur-workflow.png` dan ditampilkan di README B.8b.

### 4.2 Capaian tahun sebelumnya

- Capaian historis disimpan di tabel `capaian_tahunan`, bukan di `dashboard_row.baseline_json`.
- Scope pilot: `KABUPATEN`, per label/IKU atau indikator, dengan tahun capaian. CRUD hanya untuk admin di tab **Capaian tahunan** halaman Master.
- Dashboard periode tahun Y membaca `capaian_tahunan` dengan `tahun = Y - 1`; header kolom ikut berubah (mis. periode 2026 → **Capaian 2025**).
- Nilai tidak disalin otomatis dan tidak menimpa histori. Jika baseline tahun sebelumnya belum diisi, sel dashboard kosong (`—`), bukan nol.
- `Capaian TW1` masih dibaca dari `dashboard_row.baseline_json`; belum masuk CRUD karena bukan baseline tahun sebelumnya.

## 5. Alur upload & workflow

1. `POST /submissions` menyimpan file + `submission` + `submission_revision` + `job` dalam satu transaksi.
2. Worker mengambil job, parse CSV/XLSX, validasi header, klasifikasi deterministik.
3. Baris ditulis ke `submission_row`; yang `MATCHED` menghasilkan `indicator_value`.
4. Frontend polling `GET /jobs/{id}`, lalu menampilkan preview.
5. Koreksi IKU lewat `PATCH`; backend memvalidasi label leaf + relasi indikator dan membuat/memperbarui `indicator_value`.
6. `POST submit` → `dikirim`; verifikator approve → `diverifikasi`; Kabid approve + catatan → `final`.
7. Final memicu `CalculationService.CalculatePeriod` dan mengisi `iku_value` / `iku_value_detail`.

## 6. Klasifikasi

| Kondisi | Status |
| --- | --- |
| Kode IKU + indikator valid dan terhubung | `MATCHED` |
| Indikator valid tapi tidak terhubung ke IKU | `MISMATCH` |
| Kode IKU tidak dikenal | `UNKNOWN_IKU` |
| Kode indikator tidak dikenal | `UNKNOWN_INDICATOR` |
| IKU induk / `is_rollup_only` | `IKU_ROLLUP_ONLY` |
| Duplikat dalam file | `DUPLICATE` |
| Numerator/denominator tidak valid | `INVALID_VALUE` |
| Scope faskes berbeda dari pengupload | `SCOPE_DENIED` |

Baris non-`MATCHED` menjadi blocker submit sampai dikoreksi.

## 7. Kalkulasi

- Leaf: `value = Σ numerator / Σ denominator × 100` dari `indicator_value` yang final.
- Induk: `RATIO_OF_SUMS` dari anak eksplisit; **bukan rata-rata persentase**.
- Status: `COMPLETE`, `PARTIAL`, `NO_DATA`, `STALE`.
- Setiap run menyimpan `calculation_run_id`, `source_hash`, `data_version`.
- Contoh seed: `.01` 90,00%; `.02` 95,00%; `.03` 100,00%; induk 850/920 = 92,39%.
- Kalkulasi dilakukan untuk scope `KABUPATEN` dan tiap `FASKES` yang punya data final.

## 8. Reusable frontend

- `components/ui`: Button, Card, Badge, Input, Select, Textarea, Modal, DataTable, PageHeader, EmptyState, Spinner, ProgressSteps, StatusBadge.
- `components/layout`: Sidebar, UserMenu (actingAs).
- `components/dashboard`: StatCard, MultiBarChart.
- `components/upload`: PreviewRow.
- `constants/enums.ts`: role, status, klasifikasi, value status, TL status, label, tone.
- `constants/labels.ts`: navigasi, page copy, table labels, empty state.
- `constants/routes.ts`: route helpers.
- `stores/session.ts`: bootstrap, actors, periode, selected period.
- `api/client.ts`: client typed + actor header.

Komponen mengikuti struktur shadcn-svelte (`components.json`, `$lib/components/ui`, `tailwind-variants`, helper `cn`) dan divariasikan dengan gaya neobrutalism.dev (border tebal, shadow keras, warna kontras). CLI shadcn-svelte tidak dijalankan agar tema neobrutalism tidak tertimpa.

**Chart.** [ChartCN](https://www.shadcndeck.com/chartcn) adalah generator React/Recharts, sehingga tidak kompatibel langsung dengan Svelte 5. Visualisasi yang diminta diimplementasikan sebagai `MultiBarChart.svelte` dengan pola ChartCN: grouped multi-series bar, legend seluruh IKU leaf, label nilai, dan tooltip native. Data berasal dari `GET /api/v1/dashboard/iku-by-faskes`.

Jangan menulis label status langsung di halaman; gunakan map di `constants/enums.ts`.

## 9. Repository pattern

```text
handler → service → repository interface → sqlite implementation
```

- `internal/repository/repository.go` berisi kontrak (`UserRepository`, `MasterRepository`, `SubmissionRepository`, `JobRepository`, `ValueRepository`, `FactorRepository`, `AuditRepository`).
- `internal/repository/sqlite` mengimplementasikan semuanya dalam satu `Store`, tetapi service hanya menerima interface yang dibutuhkan.
- Service tidak mengimpor `gin` dan tidak menulis SQL.
- Worker memakai repository/service domain, bukan handler HTTP.
- Transaksi multi-tabel (upload, simpan baris, kalkulasi) berada di repository.

## 10. Special change dan notifikasi SLA

**Approval UI `SPECIAL_CHANGE`** adalah layar khusus untuk mengubah data yang sudah `final`:

1. Pengaju memilih baris/nilai final dan mengisi **alasan**.
2. Sistem menampilkan **diff before/after**, file sumber, dan dampak ke IKU induk.
3. Perubahan harus disetujui **Kabid/Sekretaris + Admin** (dual approval); satu approval saja tidak cukup.
4. Setelah disetujui, sistem membuat versi/revisi baru, menghitung ulang subtree terdampak, dan menulis audit.
5. Data final lama tidak dihapus; histori tetap dapat ditelusuri.

**Notifikasi SLA** adalah pengingat/eskalasi ketika tahap input, verifikasi, atau validasi melewati batas waktu yang ditetapkan Dinas. Notifikasi **tidak boleh auto-approve**; hanya mengingatkan approver dan atasan. Nilai SLA (berapa hari per tahap) masih keputusan terbuka Dinas, sehingga layar/notifikasi belum diimplementasikan.

Keduanya sudah disiapkan secara domain/audit, tetapi belum memiliki layar khusus pada pilot ini.

## 11. Batasan pilot dan pekerjaan produksi

- SQLite single-writer; job queue memakai `BEGIN IMMEDIATE`, bukan `FOR UPDATE SKIP LOCKED`.
- Auth pilot memakai PASETO v4.local + refresh cookie; produksi wajib rotasi kunci, CSRF, MFA approver, dan hardening cookie.
- Belum ada layar approval khusus `SPECIAL_CHANGE` dan notifikasi SLA; mekanisme audit/alasan sudah tersedia, tetapi dual approval formal menunggu keputusan Dinas.
- Audit belum hash chain.
- PII BNBA belum diputuskan; jangan menambah kolom PII tanpa keputusan Dinas.
- Status dashboard dihitung dari target periode: `ON_TARGET` bila capaian ≥ target, `BELOW_TARGET` bila di bawah target; tidak memakai ambang warna yang belum diputuskan Dinas.
- Ekspor PDF/Excel, peta, notifikasi SLA, dan integrasi SIKN belum diimplementasikan.
- Migrasi PostgreSQL wajib mempertahankan kontrak API/UI dan enum.
