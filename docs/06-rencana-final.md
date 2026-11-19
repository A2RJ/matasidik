# Rencana Final MATA SIDIK

**Status:** `IMPLEMENTED (inti + UI role)` — Track B/C/D/E, penghapusan `iku_owner`, pengajuan IKU, approval Puskesmas, dan pembatasan tampilan per peran sudah diimplementasikan; sisa diagram WebP/regenerasi dokumentasi.
**Tanggal:** 2026-11-19
**Sumber keputusan:** `Proses Bisnis Mata Sidik.xlsx`, `README.md`, `docs/01`–`docs/04`, `docs/rules.md`, dan keputusan pemilik proyek dalam diskusi.
**Plan pendukung:** [`05-rencana-bnba-wilayah.md`](05-rencana-bnba-wilayah.md) untuk detail BNBA/sebaran wilayah dan dashboard dinamis.

> Dokumen ini adalah **satu sumber urutan pekerjaan**. Jika terjadi konflik dengan `docs/rules.md`, berhenti dan minta keputusan manusia.

---

## 0. Ringkasan

Ada lima track pekerjaan:

| Track | Isi | Status |
| --- | --- | --- |
| A | BNBA Excel line-level, domisili, sebaran wilayah | **BLOCKED** — menunggu format BNBA dan definisi denominator |
| B | Dashboard kategori dinamis + CRUD kategori/bidang + preview pemetaan | **SELESAI sebagian besar** |
| C | Laporan absolut bulanan dari IKU/SPM + target layanan + ekspor | **SELESAI** |
| D | Tata kelola pengajuan IKU: Operator → Puskesmas → Admin | **IMPLEMENTED** |
| E | Workflow BNBA: tambah approval Puskesmas sebelum Dinas | **IMPLEMENTED** |
| F | Tampilan per peran, route guard, label UI bersih | **IMPLEMENTED** |

Perubahan role yang sudah dikunci:

- Tambah role **Puskesmas approver** untuk menyetujui BNBA dan pengajuan IKU.
- Hapus role `iku_owner`.
- Hapus kolom `label.iku_owner_id`.
- Admin dapat CRUD IKU **global** langsung dari Master, dan tetap approve/reject pengajuan yang sudah lolos Puskesmas.

### Riwayat batch (ringkas)

| Batch | Perubahan | Dampak aturan/dokumen |
| --- | --- | --- |
| 1 | Kategori dashboard dinamis, CRUD kategori/bidang, preview pemetaan, filter **Semua**, laporan absolut bulanan + target layanan + ekspor | `docs/04`, `docs/05`, `docs/rules.md`, `docs/schema.md`, README |
| 2 | Role `pkm_approver`, hapus `iku_owner`, pengajuan IKU, approval Puskesmas, workflow `APPROVE_PKM`, Admin CRUD IKU global | `docs/rules.md` §6/§10, `docs/04` §4, `docs/06`, `docs/schema.md`, README B.4 |
| 3 | Sidebar/route guard per peran, aksi per peran, panduan per peran, pembersihan label UI, pembatasan endpoint faktor/TL | `docs/rules.md` §10.4, `docs/04` §4.5, `docs/implementation.md` §6.3, README B.4, `docs/05` checklist |

---

## 1. Prinsip yang tidak boleh dilanggar

1. Operator/Puskesmas **tidak boleh** langsung mempublikasikan IKU. Admin dapat CRUD IKU global langsung atau mempublikasikan pengajuan yang sudah lolos Puskesmas.
2. Approval bertingkat, tidak ada auto-approve.
3. Setiap transisi approval menulis aktor, waktu, aksi, catatan, diff, dan versi.
4. Data final terkunci; perubahan pasca-final lewat mekanisme khusus.
5. Tidak ada hard delete master; gunakan soft delete `aktif = 0`.
6. Satu nilai dasar, banyak label; nilai tidak diduplikasi.
7. Jangan mengarang format BNBA, denominator, target per wilayah, atau ambang status.
8. Jangan mengedit HTML/WebP diagram manual; edit generator dan rebuild.
9. Jangan commit sebelum diminta.

---

## 2. Status implementasi saat ini

### Sudah ada

- `dashboard_category` + 8 kategori seed.
- CRUD kategori, bidang/timker, label/IKU, indikator.
- `label.category_id`; label/IKU wajib kategori.
- Indikator multi-label + label utama + PJ bidang/timker.
- Preview Pemetaan di header Master.
- Filter dashboard per kategori, opsi **Semua**, filter faskes/status yang bisa kembali ke “Semua”.
- Laporan absolut bulanan: `target_layanan_tahunan`, `target_layanan_bulanan`, endpoint `dashboard/spm/monthly`, ekspor XLSX/PDF, section di dashboard.
- Template BNBA aggregate menambahkan `kode_kategori` opsional.
- Badge status capaian dashboard memakai centang hijau / X merah.
- Perbaikan layout tabel lebar (`nowrap` + scroll horizontal).

### Sudah diimplementasikan pada batch ini

- Role `pkm_approver` + position/akun demo `kepalapkm`.
- Tabel `iku_change_request`, `iku_change_revision`, `iku_change_review`.
- API/UI pengajuan IKU dan approval Puskesmas/Admin.
- Step `APPROVE_PKM` pada workflow BNBA + UI persetujuan.
- Penghapusan role `iku_owner` dan kolom `label.iku_owner_id` (migrasi 008).
- CRUD label langsung dari Master ditutup untuk Operator/PJ; Admin tetap dapat CRUD IKU global. Operator memakai menu Pengajuan IKU.
- Pembatasan tampilan per peran: sidebar, route guard, tombol aksi, dan tab panduan.
- Dashboard Operator/PJ/Kepala Puskesmas dibatasi ke faskes sendiri, termasuk grafik capaian.
- Pembersihan label UI dari istilah tech stack dan internal.
- Perbaikan mode **lihat sebagai**: tetap aktif saat reload/refresh, pulih setelah refresh token, dan kembali ke Admin lewat tombol **Kembali**.
- Navigasi daftar submission memakai router SPA, bukan reload penuh.
- Scope `SubmissionDetail`, `Revisions`, `File`, dan `Job` dibatasi faskes untuk Operator/PJ/Kepala Puskesmas.
- `GET /factors` dan `GET /tindak-lanjut` dibatasi untuk Kabid, Sekretaris, Pimpinan, dan Admin.

### Dampak aturan dan file

- Aturan akses: `frontend/src/lib/constants/access.ts` (sumber tunggal).
- Penjaga rute: `frontend/src/routes/+layout.svelte`; URL di luar peran menampilkan **Akses terbatas**.
- Sidebar: `frontend/src/lib/components/layout/sidebar.svelte`.
- Label: `frontend/src/lib/constants/labels.ts`, `frontend/src/lib/constants/enums.ts`.
- Halaman: dashboard, detail SPM, upload, submissions, submission detail, iku-requests, master, factors, audit, admin, profile, user-guide.
- Backend: `backend/internal/httpapi/handler/factor.go`, `backend/internal/httpapi/handler/dashboard.go`, `backend/internal/httpapi/handler/submission.go`, `backend/internal/service/dashboard.go`.
- Sesi/mode lihat sebagai: `frontend/src/lib/stores/auth.ts`, `frontend/src/lib/auth/token.ts`.
- Navigasi SPA: `frontend/src/routes/submissions/+page.svelte`.
- Panduan: `frontend/src/lib/content/user-guide.ts`, `frontend/src/routes/user-guide/+page.svelte`.
- Dokumen aturan: `README.md` B.4, `docs/rules.md` §10.1/§10.4, `docs/04-arsitektur-iku.md` §4.5, `docs/implementation.md` §6.3.

### Belum selesai

- Regenerasi diagram WebP untuk alur bisnis, alur peran, workflow, dan skema.
- BNBA line-level/domisili (Track A) karena format asli belum tersedia.
- Resubmit/revisi pengajuan IKU yang ditolak belum punya layar khusus; operator membuat pengajuan baru.

---

## 3. Track A — BNBA Excel, domisili, dan sebaran wilayah

### Tujuan

Puskesmas mengunggah BNBA sesuai format yang sudah ada; sistem memetakan ibu hamil per desa/puskesmas/kecamatan dan menandai wilayah dengan ibu hamil yang tidak mendapat pelayanan sesuai standar.

### Keputusan yang sudah ada

- Template BNBA aggregate sekarang memiliki kolom opsional `kode_kategori`.
- `kode_iku` adalah kode label/IKU leaf.
- `kode_kategori` divalidasi terhadap kategori label.
- Tiga baris duplikat ibu hamil di Excel **tidak diimpor terpisah**; laporan memakai baris dashboard/IKU.

### Masih terbuka

- Contoh file BNBA `.xlsx`/`.xls` asli dari Puskesmas.
- Apakah BNBA memuat seluruh ibu hamil atau hanya yang dilayani.
- Definisi “sesuai standar”.
- Denominator/sasaran per desa.
- PII dan retensi file.
- Format multi-Puskesmas dan profil dataset.
- Peta vs tabel/ranking.
- Target per wilayah/faskes.

### Fase

1. Inventaris format dan anonimisasi sampel.
2. Definisikan standar, denominator, PII, dan output.
3. Buat profil dataset + parser multi-sheet.
4. Agregasi per desa/kecamatan, anti-double-count.
5. Dashboard wilayah + trace ke file.
6. Hardening (dedupe, peta/ekspor bila disetujui).

---

## 4. Track B — Dashboard kategori dinamis

### Keputusan final

- Kategori disimpan di `dashboard_category`; dapat di-CRUD Admin.
- `jenis`: `LABEL` atau `OWNER`.
- 8 kategori seed: IKU, IKK, SPM, RIBK, Program, Kegiatan, Sub Kegiatan, Bidang/Sekretariat/Tim Kerja.
- Kategori ke-9 tetap kosong sampai Dinas memutuskan.
- Kategori yang dipakai memakai soft delete `aktif = 0`.
- `bidang` memakai satu tabel dengan `jenis` (`bidang`/`sekretariat`/`timker`) dan `parent_id`.
- Indikator memakai `pj_bidang_id` dan `pj_timker_id`.
- Setiap label/IKU wajib satu kategori.
- Indikator memilih **label/IKU spesifik**, bukan kategori langsung.
- “Indikator Sasaran” adalah bagian dari IKU, bukan kategori filter terpisah.
- Kategori 8 adalah daftar/filter per owner, bukan angka komposit/rata-rata persentase.
- Kolom Indikator di semua filter Kelompok hanya menampilkan teks, tanpa badge.
- Filter Kelompok memiliki opsi **Semua**; mode ini menampilkan gabungan label semua kategori + owner, dengan kolom Kelompok.

### Sisa pekerjaan

- Isi label untuk IKK/RIBK/Program/Kegiatan/Sub Kegiatan; saat ini masih kosong.
- Isi target per kategori/label selain SPM.
- Regenerasi diagram skema setelah penambahan tabel/kolom.
- Audit data label dari sheet 2 dan hapus duplikasi katalog.

---

## 5. Track C — Laporan absolut bulanan

### Keputusan final

- `Target 2026 = 2000` adalah **target layanan tahunan**, bukan target persentase.
- `Target Bulanan = target layanan tahunan / 12`.
- Status bulanan membandingkan **jumlah layanan bulan itu** dengan **target bulanan**.
- Laporan ditampilkan sebagai section + ekspor XLSX/PDF.
- Data diambil dari IKU/SPM/dashboard, bukan tiga baris duplikat Excel.
- Bulan tanpa data = `NO_DATA`, bukan nol.
- Target layanan disimpan terpisah dari target capaian persentase:
  - `target_layanan_tahunan`
  - `target_layanan_bulanan`

### Sudah diimplementasikan

- Migrasi kolom target layanan.
- Seed target layanan 2000 untuk `IKU-2026-001.01`.
- API `GET /api/v1/dashboard/spm/monthly?tahun=2026`.
- API ekspor `GET /api/v1/dashboard/spm/monthly/export?tahun=2026&format=xlsx|pdf`.
- Section **Data Absolut Kabupaten** di dashboard.
- Field target layanan tahunan di Master Target.

### Catatan data

- Demo saat ini hanya memiliki nilai Mei 2026; bulan lain `NO_DATA`.
- Nilai Mei demo 450 lebih besar dari target bulanan 166,67 sehingga status `ON_TARGET`. Jika 450 adalah angka kumulatif, sumber datanya harus dibedakan; jangan mengubah makna data diam-diam.
- Target layanan indikator selain 01 harus diisi manual lewat Master Target.

---

## 6. Track D — Tata kelola pengajuan IKU

### 6.1 Role

Tambah role:

| Kode usulan | Nama tampilan | Scope |
| --- | --- | --- |
| `pkm_approver` | Puskesmas / Kepala Puskesmas | Faskes sendiri |

Hak akses:

- Approve/reject submission BNBA dari faskes sendiri.
- Approve/reject pengajuan IKU dari faskes sendiri.
- Tidak boleh approve pengajuannya sendiri.
- Tidak bisa upload BNBA.
- Tidak bisa mempublikasikan IKU.

Kode role yang dipakai: `pkm_approver`; label UI: **Puskesmas Approver**; jabatan: **Kepala Puskesmas**.

### 6.2 Tabel

#### `iku_change_request`

| Kolom | Relasi | Keterangan |
| --- | --- | --- |
| `id` | PK | ID pengajuan |
| `action` | CHECK `CREATE`/`UPDATE`/`DEACTIVATE` | Jenis pengajuan |
| `label_id` | FK → `label.id`, nullable | IKU lama untuk update/deactivate |
| `requested_by` | FK → `users.id` | Operator Puskesmas pengaju |
| `faskes_id` | FK → `faskes.id` | Faskes asal |
| `master_release_id` | FK → `master_release.id` | Release tujuan |
| `current_revision_no` | INTEGER | Revisi terakhir |
| `status` | TEXT | Status pengajuan |
| `applied_label_id` | FK → `label.id`, nullable | IKU hasil publish |
| `created_at` / `updated_at` | TEXT | Waktu |

Status:

`draf`, `menunggu_puskesmas`, `dikembalikan_puskesmas`, `menunggu_admin`, `dikembalikan_admin`, `disetujui`, `ditolak`, `dibatalkan`.

#### `iku_change_revision`

Immutable snapshot payload setiap kali operator submit/revisi.

| Kolom | Relasi |
| --- | --- |
| `id` | PK |
| `request_id` | FK → `iku_change_request.id` |
| `revision_no` | INTEGER |
| `payload_json` | TEXT |
| `before_json` / `after_json` | TEXT |
| `created_by` | FK → `users.id` |
| `created_at` | TEXT |

#### `iku_change_review`

Riwayat keputusan per tahap.

| Kolom | Relasi |
| --- | --- |
| `id` | PK |
| `request_id` | FK → `iku_change_request.id` |
| `revision_no` | INTEGER |
| `step` | CHECK `PKM`/`ADMIN` |
| `reviewer_id` | FK → `users.id` |
| `action` | CHECK `approve`/`reject`/`return` |
| `notes` | TEXT |
| `created_at` | TEXT |

### 6.3 Alur

1. Operator Puskesmas membuat pengajuan `CREATE`/`UPDATE`/`DEACTIVATE` untuk faskesnya.
2. Status `menunggu_puskesmas`.
3. Puskesmas approver faskes yang sama:
   - approve → `menunggu_admin`;
   - reject/return → `dikembalikan_puskesmas` + catatan wajib.
4. Operator dapat revisi dan submit ulang; revisi immutable.
5. Admin:
   - approve → status `disetujui`, lalu apply ke `label` dalam transaksi;
   - reject/return → `dikembalikan_admin` + catatan wajib.
6. Hanya setelah Admin approve, IKU aktif dan muncul di dashboard.
7. `DEACTIVATE` hanya mengubah `label.aktif = 0`; tidak ada hard delete.

### 6.4 Hak akses

| Aksi | Operator Puskesmas | Kepala Puskesmas | Admin |
| --- | --- | --- | --- |
| Buat & kirim pengajuan | ✅ faskes sendiri | ❌ | ❌ |
| Approve/reject Puskesmas | ❌ | ✅ faskes sendiri | ❌ |
| Approve/reject final + publish | ❌ | ❌ | ✅ |
| CRUD IKU global langsung | ❌ | ❌ | ✅ |
| Lihat riwayat | ✅ milik sendiri | ✅ faskes sendiri | ✅ semua |

Tidak ada self-approval. Admin tidak boleh mempublikasikan pengajuan yang belum lolos Puskesmas.

### 6.5 API dan UI

Endpoint yang dipakai:

- `POST /api/v1/iku/change-requests`
- `GET /api/v1/iku/change-requests`
- `GET /api/v1/iku/change-requests/{id}`
- `POST /api/v1/iku/change-requests/{id}/decision`
- Admin CRUD global: `POST/PATCH/DELETE /api/v1/master/labels` (DELETE = nonaktifkan)

UI:

- Master → Label & IKU: Admin melihat **IKU Global**, Edit, dan Nonaktifkan. Operator/PJ tidak melihat aksi CRUD langsung.
- Menu **Pengajuan IKU**: Operator Puskesmas membuat pengajuan; Puskesmas approver memutus faskes sendiri; Admin memutus tahap final.
- Tombol Setujui/Kembalikan/Tolak dengan catatan wajib saat mengembalikan/menolak.
- Timeline review dan jejak audit.
- Panduan in-app menyaring tab peran; Admin dapat membuka semua.

### 6.6 Penghapusan `iku_owner`

Keputusan: hapus role dan kolom.

- Hapus `RoleIKUOwner` dari enum dan `AllRoles`.
- Hapus `iku_owner` dari frontend Role/label role.
- Hapus `iku_owner_id` dari `domain.Label`, input service, query repository.
- Migrasi baru `ALTER TABLE label DROP COLUMN iku_owner_id`.
- Migrasi 001/002 tidak diubah karena sudah menjadi riwayat.
- Backup database sebelum migrasi; data kolom akan hilang permanen.
- Tidak ada akun/position `iku_owner` yang perlu dimigrasikan.

---

## 7. Track E — Workflow BNBA dengan approval Puskesmas

### Alur baru

1. `DRAFT_PKM` — Operator/PJ upload dan preview.
2. `APPROVE_PKM` — Puskesmas approver faskes sendiri.
3. `VERIFY_DINAS` — Petugas Dinas.
4. `VALIDATE_KABID` — Kabid/Sekretaris.
5. `FINAL` — sistem mengunci dan menghitung.

Aturan:

- Puskesmas reject → kembali ke operator (`dikembalikan_petugas`) dengan catatan.
- Tidak ada auto-approve.
- Operator tidak boleh menyetujui submission sendiri.
- Submission baru memakai workflow v2; submission lama tetap memakai versi yang di-pin.
- Puskesmas approver hanya melihat submission faskesnya.
- Admin dapat melihat semua, tetapi bukan approver default.

### Keputusan status

Dipakai **step-only**: status tetap `dikirim`, step berubah dari `APPROVE_PKM` ke `VERIFY_DINAS`. Enum 7 status formal tidak diubah. Filter antrean menggunakan `step`.

### UI

- Filter submission: Menunggu Puskesmas, Menunggu Dinas, Menunggu Kabid, Final.
- Detail submission menampilkan tombol Setujui/Kembalikan untuk `pkm_approver`.
- Timeline menampilkan aktor Puskesmas.
- Upload hanya untuk operator/PJ; approver tidak bisa upload.

---

## 8. Model data dan migrasi

### Tabel baru

- `iku_change_request`
- `iku_change_revision`
- `iku_change_review`

### Perubahan role

- Tambah `pkm_approver`.
- Hapus `iku_owner`.

### Perubahan kolom

- Hapus `label.iku_owner_id`.
- Tambah target layanan (sudah ada di 007):
  - `target_layanan_tahunan`
  - `target_layanan_bulanan`

### Urutan migrasi yang dipakai

1. `008_drop_iku_owner.sql` — hapus kolom `label.iku_owner_id`.
2. `009_iku_change_and_pkm_approval.sql` — role/position/user demo `pkm_approver`, tabel pengajuan IKU, dan workflow v2.

Jangan mengubah migrasi yang sudah pernah diterapkan; gunakan migrasi baru.

---

## 9. Diagram dan dokumentasi

Karena aturan diagram sekarang memakai **HTML + WebP**:

- Edit generator, bukan hasil render:
  - `assets/diagrams/gen.py` — `alur-bisnis`, `alur-peran`.
  - `assets/diagrams/gen-teknis.py` — `alur-workflow`.
  - `assets/diagrams/gen-schema.py` — kolom/tabel baru.
- Ubah label role menjadi **Operator Puskesmas/PJ — pengupload data BNBA**.
- Tambah lane/node **Puskesmas approve** sebelum Dinas.
- Tambah langkah `APPROVE_PKM` pada `alur-workflow`.
- Tambah diagram baru untuk pengajuan IKU: Operator → Puskesmas → Admin.
- Rebuild `build.sh` → WebP.
- Update `assets/diagrams/README.md`, README A.2/B.2/B.8, `docs/01`, `docs/04`, `docs/rules`, `docs/schema`, `docs/implementation`.

---

## 10. Roadmap eksekusi

| Fase | Pekerjaan | Dependensi |
| --- | --- | --- |
| 1 | Backup DB; hapus role/kolom `iku_owner` | Persetujuan migrasi |
| 2 | Tambah role `pkm_approver` + position/user demo | Fase 1 |
| 3 | Tabel + API + service pengajuan IKU | Fase 2 |
| 4 | UI operator/Puskesmas/Admin + RBAC | Fase 3 |
| 5 | Workflow BNBA `APPROVE_PKM` + UI antrean | Fase 2 |
| 6 | Diagram WebP + dokumentasi | Fase 4–5 |
| 7 | Tes, build, audit, verifikasi | Semua |

Track A (BNBA line-level/domisili) tetap berjalan terpisah dan **blocked** sampai format BNBA diterima.

---

## 11. Acceptance criteria

- Hanya Operator Puskesmas yang bisa mengajukan IKU; Operator/PJ tidak bisa langsung mempublikasikan IKU.
- Pengajuan IKU hanya bisa diputus Puskesmas dari faskes yang sama.
- Untuk pengajuan operator, Admin hanya bisa approve setelah Puskesmas approve.
- Admin dapat CRUD IKU global langsung tanpa pengajuan.
- Reject/return wajib catatan; operator bisa revisi.
- IKU hasil pengajuan baru aktif setelah Admin approve; IKU yang dibuat Admin langsung aktif setelah disimpan.
- Semua transisi memiliki audit dan review.
- BNBA baru: operator upload → Puskesmas approve → Dinas verify → Kabid validate → final.
- Submission lama tidak berubah workflow.
- Tidak ada self-approval.
- `iku_owner` tidak lagi ada di role/kolom/dokumen.
- Diagram WebP dan dokumentasi diperbarui.
- `go test ./...`, `go build ./...`, `npm run check`, `npm run build`, dan `git diff --check` lulus.

---

## 12. Keputusan terbuka — jangan ditebak

| # | Pertanyaan | Status |
| --- | --- | --- |
| 1 | Nama final role `pkm_approver` vs `kepala_puskesmas` | Dipakai `pkm_approver` |
| 2 | Status setelah approve Puskesmas: `disetujui_pkm` atau step-only | Dipakai step-only |
| 3 | Field IKU yang boleh diubah operator | Belum final |
| 4 | Pengajuan IKU baru berlaku global atau per faskes | Perlu penegasan |
| 5 | Format BNBA asli + definisi denominator/PII | Belum ada |
| 6 | Target per faskes/wilayah | Belum diputuskan |
| 7 | Regenerasi diagram dan schema-db | Belum dijalankan |
| 8 | Backup/retensi sebelum drop `iku_owner_id` | Wajib disiapkan |

---

## 13. Verifikasi wajib

- Baca ulang `docs/rules.md` setelah perubahan workflow/RBAC.
- Uji migrasi pada database kosong dan database existing.
- Uji role/scope: operator, Puskesmas approver, Dinas, Kabid, Admin.
- Uji self-approval ditolak.
- Uji reject/return dan revisi.
- Uji publish IKU dan audit.
- Uji submission BNBA v2 dan kompatibilitas v1.
- Periksa semua tautan, anchor, tabel, dan diagram WebP.
- Tidak ada Mermaid/ASCII flow baru.
- Tidak ada commit sebelum diminta.

---

## 14. Referensi

- [`README.md`](../README.md)
- [`docs/rules.md`](rules.md)
- [`docs/01-gambaran-umum.md`](01-gambaran-umum.md)
- [`docs/02-daftar-indikator-kinerja.md`](02-daftar-indikator-kinerja.md)
- [`docs/03-contoh-data.md`](03-contoh-data.md)
- [`docs/04-arsitektur-iku.md`](04-arsitektur-iku.md)
- [`docs/05-rencana-bnba-wilayah.md`](05-rencana-bnba-wilayah.md)
- [`docs/schema.md`](schema.md)
- [`docs/implementation.md`](implementation.md)
- [`assets/diagrams/README.md`](../assets/diagrams/README.md)
