# Skema Database MATA SIDIK

> **Status implementasi pilot.** Skema ini adalah kontrak database aplikasi pada implementasi awal. Pemilik proyek meminta **SQLite lebih dulu** dan **PASETO v4.local**. `docs/rules.md` §4 menyebut PostgreSQL 16 + PASETO untuk produksi; karena itu, migrasi ke PostgreSQL wajib dilakukan sebelum produksi dan tetap mengikuti aturan proyek.
>
> Cakupan: **31 tabel domain** + `schema_migrations`, mencakup master, submission, kalkulasi, faktor/TL, audit, auth, jabatan struktural, dan kategori dashboard dinamis.
>
> Diagram tabel, kolom, dan relasi: [`../assets/diagrams/schema-db.webp`](../assets/diagrams/schema-db.webp) · sumber HTML: [`../assets/diagrams/schema-db.html`](../assets/diagrams/schema-db.html) · generator: [`../assets/diagrams/gen-schema.py`](../assets/diagrams/gen-schema.py).
>
> **Catatan:** diagram `schema-db` belum diregenerasi setelah penambahan `dashboard_category`, `label.category_id`, `bidang.jenis`, dan `indikator.pj_timker_id`; generator sudah diperbarui, WebP menunggu build ulang.

## 1. Prinsip

1. **Satu nilai dasar, banyak label.** Nilai indikator disimpan sekali di `indicator_value`; label menghubungkan nilai ke dashboard/IKU.
2. **Kode adalah kontrak.** `kode_iku` dan `kode_indikator` stabil; nama/formula boleh berubah lewat versi master.
3. **Tidak ada auto-create master dari file upload.** Kode tidak dikenal masuk karantina (`submission_row.classification_status`).
4. **Angka harus bisa ditelusuri.** `dashboard_row → iku_value → iku_value_detail → indicator_value → submission_row → submission → file_object`.
5. **Final terkunci.** Perubahan pasca-final harus lewat mekanisme khusus dengan alasan + audit; implementasi pilot mencatat audit tanpa dual approval UI.
6. **Repository pattern.** Semua akses SQL berada di `backend/internal/repository/sqlite`; service tidak menulis SQL langsung.
7. **SQLite disiapkan untuk migrasi.** Nama tabel/kolom, status enum, dan relasi dibuat selaras dengan rancangan PostgreSQL. Tipe `TEXT`/`INTEGER`/`REAL` dipakai di pilot.

## 2. Konvensi

| Aspek | Konvensi |
| --- | --- |
| ID | UUID string (`TEXT`), dibuat aplikasi |
| Waktu | RFC3339 UTC (`TEXT`) |
| Boolean | `INTEGER` 0/1; di JSON menjadi `true/false` |
| Uang/desimal | `REAL`; pembulatan hanya di tampilan, bukan data mentah |
| JSON | `TEXT` berisi JSON valid (SQLite tidak punya JSONB) |
| Enum | Disimpan sebagai `TEXT`, divalidasi di domain + CHECK constraint untuk nilai kritis |
| Soft delete | `aktif = 0`; data histori tidak dihapus fisik |
| Nama tabel | `snake_case`, tunggal untuk entitas inti (`label`, `submission`) |
| Nama kolom | `snake_case`, bahasa Indonesia untuk domain (`nama`, `periode_id`) |

## 3. Kelompok tabel

| Kelompok | Tabel |
| --- | --- |
| Master inti | `users`, `positions`, `auth_session`, `wilayah`, `faskes`, `bidang`, `periode`, `master_release`, `workflow_version` |
| Taksonomi & target | `dashboard_category`, `label`, `indikator`, `indikator_label`, `iku_komponen`, `target`, `capaian_tahunan`, `dashboard_row` |
| Submission & file | `file_object`, `submission`, `submission_revision`, `submission_group`, `submission_row`, `submission_review`, `job` |
| Nilai & kalkulasi | `indicator_value`, `calculation_run`, `iku_value`, `iku_value_detail` |
| Faktor, TL & audit | `faktor`, `tindak_lanjut`, `audit_log` |
| Internal | `schema_migrations` |

---

## 4. Master inti

### 4.1 `users`

Pengguna, peran, username, password hash, dan posisi struktural. Autentikasi memakai PASETO v4.local + refresh cookie.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID pengguna |
| `nama` | TEXT | NOT NULL | Nama tampilan |
| `username` | TEXT | UNIQUE | Username login |
| `password_hash` | TEXT | | Hash bcrypt |
| `role` | TEXT | NOT NULL | Role dasar; role efektif dapat berasal dari posisi |
| `faskes_id` | TEXT | FK → `faskes.id` | Untuk operator Puskesmas |
| `bidang_id` | TEXT | FK → `bidang.id` | Untuk Kabid/PJ |
| `position_id` | TEXT | FK → `positions.id` | Jabatan struktural |
| `last_login_at` | TEXT | | Login terakhir |
| `aktif` | INTEGER | NOT NULL | Soft delete |
| `created_at` | TEXT | NOT NULL | Waktu dibuat |

### 4.2 `positions`

Jabatan struktural yang dapat diisi user berbeda tanpa mengubah kode workflow. Role efektif user diambil dari posisi.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID posisi |
| `kode` | TEXT | UNIQUE | Kode posisi |
| `nama` | TEXT | NOT NULL | Nama jabatan |
| `role` | TEXT | NOT NULL | Role efektif (`dinas_verifier`, `kabid`, dll.) |
| `level` | INTEGER | NOT NULL | Level approval |
| `bidang_id` | TEXT | FK → `bidang.id` | Scope bidang |
| `faskes_id` | TEXT | FK → `faskes.id` | Scope faskes |
| `aktif` | INTEGER | NOT NULL | Soft delete |
| `created_at` | TEXT | NOT NULL | Waktu dibuat |

### 4.3 `auth_session`

Sesi refresh token. Access token PASETO tidak disimpan; refresh token disimpan sebagai hash.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID sesi |
| `user_id` | TEXT | FK → `users.id` | Pemilik sesi |
| `refresh_hash` | TEXT | UNIQUE | SHA-256 refresh token |
| `user_agent` | TEXT | | User agent |
| `created_at` | TEXT | NOT NULL | Dibuat |
| `expires_at` | TEXT | NOT NULL | Kedaluwarsa |
| `revoked_at` | TEXT | | Revoke |

### 4.4 `wilayah`

Hierarki **desa → kecamatan → kabupaten**.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID wilayah |
| `kode` | TEXT | UNIQUE, NOT NULL | Kode wilayah |
| `nama` | TEXT | NOT NULL | Nama wilayah |
| `level` | TEXT | CHECK `desa/kecamatan/kabupaten` | Level hierarki |
| `parent_id` | TEXT | FK → `wilayah.id` | Parent hierarki |

### 4.5 `faskes`

Puskesmas/faskes pengirim data.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID faskes |
| `kode` | TEXT | UNIQUE, NOT NULL | Contoh `PKM-TALIWANG` |
| `nama` | TEXT | NOT NULL | Nama faskes |
| `wilayah_id` | TEXT | FK → `wilayah.id` | Wilayah faskes |
| `jenis` | TEXT | NOT NULL | `puskesmas` dan sejenisnya |
| `aktif` | INTEGER | NOT NULL | Soft delete |

### 4.6 `bidang`

Bidang/sekretariat/timkerja pemilik data.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID bidang |
| `kode` | TEXT | UNIQUE, NOT NULL | Kode bidang |
| `nama` | TEXT | NOT NULL | Nama bidang |
| `jenis` | TEXT | NOT NULL | `bidang`, `sekretariat`, atau `timker` |
| `parent_id` | TEXT | FK → `bidang.id` | Parent untuk timker |
| `aktif` | INTEGER | NOT NULL | Soft delete |

### 4.7 `periode`

Bulan pelaporan; upload dan perhitungan terikat ke periode.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID periode |
| `tahun` | INTEGER | NOT NULL | Tahun |
| `bulan` | INTEGER | NOT NULL, CHECK 1–12 | Bulan |
| `kode` | TEXT | UNIQUE, NOT NULL | `YYYY-MM`, contoh `2026-05` |
| `status` | TEXT | CHECK `open/locked` | Periode terkunci tidak menerima upload |
| `is_late_allowed` | INTEGER | NOT NULL | Kebijakan keterlambatan |

### 4.8 `master_release`

Bundel master per tahun. Upload mengikat ke release yang berlaku.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID release |
| `tahun` | INTEGER | NOT NULL | Tahun master |
| `status` | TEXT | CHECK `draft/published/archived` | Hanya `published` yang dipakai upload |
| `published_at` | TEXT |  | Waktu publikasi |
| `published_by` | TEXT | FK → `users.id` | Aktor publikasi |

### 4.9 `workflow_version`

Versi workflow yang di-pin ke submission. Perubahan workflow tidak mengubah pengajuan berjalan.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID versi |
| `kode` | TEXT | NOT NULL | Contoh `default_3_level` |
| `version_no` | INTEGER | NOT NULL | Nomor versi |
| `steps_json` | TEXT | NOT NULL | Definisi step/role |
| `status` | TEXT | NOT NULL | `published` |
| `published_at` | TEXT |  | Waktu publikasi |

---

## 5. Taksonomi, indikator, dan target

### 5.0 `dashboard_category`

Kategori dashboard yang dapat di-CRUD admin. `jenis = LABEL` dipakai untuk taksonomi label/IKU; `jenis = OWNER` dipakai untuk filter Bidang/Sekretariat/Tim Kerja dan tidak memiliki label.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID kategori |
| `kode` | TEXT | UNIQUE, NOT NULL | `IKU`, `IKK`, `SPM`, `RIBK`, `PROGRAM`, `KEGIATAN`, `SUB_KEGIATAN`, `BIDANG_TIMKER` |
| `nama` | TEXT | NOT NULL | Nama tampilan filter |
| `jenis` | TEXT | CHECK `LABEL/OWNER` | Jenis kategori |
| `urutan` | INTEGER | NOT NULL | Urutan filter |
| `aktif` | INTEGER | NOT NULL | Soft delete; kategori yang dipakai tetap disimpan |
| `is_system` | INTEGER | NOT NULL | Kategori seed tidak boleh dihapus fisik |
| `config_json` | TEXT | | Konfigurasi tambahan |
| `created_at` | TEXT | NOT NULL | Dibuat |
| `updated_at` | TEXT | NOT NULL | Diubah |

### 5.1 `label`

Taksonomi sekaligus node IKU. **`parent_id` adalah hierarki**, bukan parsing titik pada kode.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID label |
| `kode` | TEXT | UNIQUE per release | `IKU-2026-001`, `IKU-2026-001.01` |
| `nama` | TEXT | NOT NULL | Nama label |
| `tipe` | TEXT | NOT NULL | Kode kategori legacy; disinkronkan dari `category_id` |
| `category_id` | TEXT | FK → `dashboard_category.id` | Kategori label/IKU; wajib untuk label baru |
| `parent_id` | TEXT | FK → `label.id` | Parent hierarki |
| `master_release_id` | TEXT | FK → `master_release.id` | Release master |
| `aggregation_method` | TEXT | CHECK | `RATIO_OF_SUMS`, `SUM`, `WEIGHTED_AVERAGE`, `FORMULA`, `NONE` |
| `allow_direct_input` | INTEGER | NOT NULL | Leaf yang boleh menerima upload |
| `is_rollup_only` | INTEGER | NOT NULL | Induk murni rollup |
| `pj_bidang_id` | TEXT | FK → `bidang.id` | Pemilik data |
| `iku_owner_id` | TEXT | FK → `users.id` | Pemilik publikasi |
| `level` | INTEGER | NOT NULL | Kedalaman hierarki |
| `path` | TEXT |  | Path bantu tampilan |
| `aktif` | INTEGER | NOT NULL | Soft delete |
| `berlaku_dari` | TEXT |  | Awal berlaku |
| `berlaku_sampai` | TEXT |  | Akhir berlaku |

### 5.2 `indikator`

Definisi indikator yang diukur.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID indikator |
| `kode` | TEXT | UNIQUE, NOT NULL | `IND-001` |
| `nama` | TEXT | NOT NULL | Nama indikator |
| `satuan` | TEXT | NOT NULL | Satuan baku |
| `definisi` | TEXT |  | Definisi operasional |
| `numerator_label` | TEXT |  | Label pembilang |
| `denominator_label` | TEXT |  | Label penyebut |
| `sifat_input` | TEXT | NOT NULL | `puskesmas`, `kabupaten` |
| `frekuensi` | TEXT | NOT NULL | `bulanan`, `tahunan` |
| `pj_bidang_id` | TEXT | FK → `bidang.id` | PJ data level bidang/sekretariat |
| `pj_timker_id` | TEXT | FK → `bidang.id` | PJ data level timker (`bidang.jenis = timker`) |
| `pj_timker` | TEXT |  | Nama timker legacy; diisi dari join `pj_timker_id` |
| `aktif` | INTEGER | NOT NULL | Soft delete |

### 5.3 `indikator_label`

Many-to-many indikator ↔ label. **Bukan** graf kalkulasi.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `indikator_id` | TEXT | PK, FK → `indikator.id` | Indikator |
| `label_id` | TEXT | PK, FK → `label.id` | Label/IKU |
| `bobot` | REAL | NOT NULL | Default 1 |
| `is_primary` | INTEGER | NOT NULL | Label utama |

### 5.4 `iku_komponen`

Graf kalkulasi induk. Hanya komponen eksplisit yang dihitung; mencegah double counting.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID komponen |
| `parent_label_id` | TEXT | FK → `label.id` | IKU induk |
| `tipe_komponen` | TEXT | CHECK `CHILD_IKU/INDIKATOR` | Jenis komponen |
| `komponen_id` | TEXT | NOT NULL | ID label anak/indikator |
| `bobot` | REAL | NOT NULL | Bobot |
| `sort` | INTEGER | NOT NULL | Urutan tampilan |

### 5.5 `target`

Target per periode/scope. `target_bulanan = target_tahunan / 12`.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID target |
| `label_id` | TEXT | FK → `label.id` | Target IKU |
| `indikator_id` | TEXT | FK → `indikator.id` | Target indikator |
| `periode_id` | TEXT | FK → `periode.id` | Periode |
| `scope_type` | TEXT | NOT NULL | `KABUPATEN`, `FASKES`, `KECAMATAN`, `DESA` |
| `scope_id` | TEXT |  | ID scope |
| `target_tahunan` | REAL | NOT NULL | Target capaian/persentase tahunan |
| `target_bulanan` | REAL | NOT NULL | `target_tahunan / 12` |
| `target_layanan_tahunan` | REAL |  | Target layanan absolut tahunan (mis. 2000 ibu hamil) |
| `target_layanan_bulanan` | REAL |  | `target_layanan_tahunan / 12` |

> Target layanan absolut dipisahkan dari target capaian persentase. Laporan bulanan absolut memakai `target_layanan_*`; dashboard persentase memakai `target`/`dashboard_row`.

### 5.6 `capaian_tahunan`

Capaian historis tahunan per label/IKU atau indikator, dipakai sebagai pembanding dashboard. Dashboard periode 2026 membaca baris dengan `tahun = 2025`. Untuk pilot, scope dibatasi `KABUPATEN`; CRUD hanya admin.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID capaian |
| `label_id` | TEXT | FK → `label.id` | IKU/label; salah satu dengan `indikator_id` |
| `indikator_id` | TEXT | FK → `indikator.id` | Indikator; salah satu dengan `label_id` |
| `tahun` | INTEGER | NOT NULL | Tahun capaian, mis. 2025 |
| `scope_type` | TEXT | NOT NULL | `KABUPATEN` untuk pilot |
| `scope_id` | TEXT |  | ID kabupaten |
| `nilai` | REAL | NOT NULL | Nilai capaian |
| `catatan` | TEXT |  | Catatan admin |
| `updated_by` | TEXT | FK → `users.id` | Admin terakhir |
| `created_at` | TEXT | NOT NULL | Dibuat |
| `updated_at` | TEXT | NOT NULL | Diperbarui |

Unique key: `(tahun, scope_type, scope_id, label_id, indikator_id)` dengan `COALESCE` untuk kolom nullable.

### 5.7 `dashboard_row`

Kontrak baris output SPM. Nilai tidak disimpan di sini; nilai dibaca dari `iku_value`.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID baris |
| `kode` | TEXT | UNIQUE, NOT NULL | `01`–`12` |
| `label_id` | TEXT | FK → `label.id` | Node nilai |
| `indikator_id` | TEXT | FK → `indikator.id` | Indikator |
| `urutan` | INTEGER | NOT NULL | Urutan tampil |
| `target` | REAL | NOT NULL | Target tampilan |
| `baseline_json` | TEXT |  | Capaian TW1; baseline tahunan ada di `capaian_tahunan` |
| `format_json` | TEXT |  | Format tampilan |

---

## 6. Submission, file, dan job

### 6.1 `file_object`

Metadata file asli di volume lokal. Nama file = `sha256` + ekstensi.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID file |
| `original_name` | TEXT | NOT NULL | Nama asli |
| `sha256` | TEXT | NOT NULL | Hash konten |
| `path` | TEXT | NOT NULL | Path volume lokal |
| `mime` | TEXT |  | MIME |
| `size` | INTEGER | NOT NULL | Ukuran byte |
| `uploaded_by` | TEXT | FK → `users.id` | Uploader |
| `scan_status` | TEXT | NOT NULL | `pending/clean/infected` |
| `created_at` | TEXT | NOT NULL | Waktu simpan |

### 6.2 `submission`

Satu pengajuan BNBA.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID submission |
| `kode` | TEXT | UNIQUE, NOT NULL | `S-2026-05-001` |
| `tipe` | TEXT | NOT NULL | `BNBA` |
| `uploader_id` | TEXT | FK → `users.id` | Uploader |
| `faskes_id` | TEXT | FK → `faskes.id` | Faskes |
| `periode_id` | TEXT | FK → `periode.id` | Periode |
| `file_id` | TEXT | FK → `file_object.id` | File aktif |
| `file_hash` | TEXT | NOT NULL | Hash file |
| `status` | TEXT | NOT NULL | 7 status formal |
| `step` | TEXT | NOT NULL | `DRAFT_PKM`, `VERIFY_DINAS`, `VALIDATE_KABID`, `FINAL`, `SPECIAL_CHANGE` |
| `master_release_id` | TEXT | FK → `master_release.id` | Release dipin |
| `workflow_version_id` | TEXT | FK → `workflow_version.id` | Workflow dipin |
| `idempotency_key` | TEXT | UNIQUE partial | Cegah upload ganda |
| `created_at` | TEXT | NOT NULL | Waktu dibuat |
| `updated_at` | TEXT | NOT NULL | Waktu diubah |

### 6.3 `submission_revision`

Revisi immutable; file lama tidak ditimpa.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID revisi |
| `submission_id` | TEXT | FK → `submission.id` | Submission |
| `revision_no` | INTEGER | UNIQUE per submission | Nomor revisi |
| `parent_revision_id` | TEXT | FK → `submission_revision.id` | Revisi sebelumnya |
| `file_id` | TEXT | FK → `file_object.id` | File revisi |
| `created_by` | TEXT | FK → `users.id` | Aktor |
| `created_at` | TEXT | NOT NULL | Waktu |

### 6.4 `submission_group`

Pecahan approval per pemilik data. Pada pilot dibuat satu group per submission.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID group |
| `submission_id` | TEXT | FK → `submission.id` | Submission |
| `owner_bidang_id` | TEXT | FK → `bidang.id` | Pemilik |
| `primary_label_id` | TEXT | FK → `label.id` | IKU utama |
| `status` | TEXT | NOT NULL | `active` |
| `current_step` | TEXT | NOT NULL | Step berjalan |

### 6.5 `submission_row`

Baris hasil parsing + klasifikasi deterministik.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID baris |
| `submission_id` | TEXT | FK → `submission.id` | Submission |
| `group_id` | TEXT | FK → `submission_group.id` | Group |
| `row_no` | INTEGER | UNIQUE per submission | Nomor baris data (R1, R2, …) |
| `kode_iku_raw` | TEXT |  | Kode IKU dari file |
| `kode_indikator_raw` | TEXT |  | Kode indikator dari file |
| `matched_indikator_id` | TEXT | FK → `indikator.id` | Hasil pencocokan |
| `matched_label_id` | TEXT | FK → `label.id` | IKU leaf hasil pencocokan |
| `scope_type` | TEXT | NOT NULL | `FASKES`, `DESA`, `KECAMATAN`, `KABUPATEN` |
| `scope_code` | TEXT |  | Kode scope |
| `nilai` | REAL |  | Nilai baris |
| `numerator` | REAL |  | Pembilang |
| `denominator` | REAL |  | Penyebut |
| `classification_status` | TEXT | NOT NULL | `MATCHED`, `MISMATCH`, `UNKNOWN_IKU`, `UNKNOWN_INDICATOR`, `IKU_ROLLUP_ONLY`, `AMBIGUOUS`, `DUPLICATE`, `INVALID_VALUE`, `SCOPE_DENIED` |
| `error_message` | TEXT |  | Pesan per baris |
| `created_at` | TEXT | NOT NULL | Waktu |

### 6.6 `submission_review`

Keputusan tiap level approval.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID review |
| `submission_id` | TEXT | FK → `submission.id` | Submission |
| `group_id` | TEXT | FK → `submission_group.id` | Group |
| `step` | TEXT | NOT NULL | Step |
| `reviewer_id` | TEXT | FK → `users.id` | Aktor |
| `action` | TEXT | NOT NULL | `submit`, `approve`, `return` |
| `notes` | TEXT |  | Wajib saat return/final |
| `details_json` | TEXT |  | Diff/komentar |
| `created_at` | TEXT | NOT NULL | Waktu |

### 6.7 `job`

Antrean background SQLite. Klaim job memakai transaksi `BEGIN IMMEDIATE`; bukan `SKIP LOCKED` PostgreSQL.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID job |
| `type` | TEXT | NOT NULL | `parse_submission` |
| `payload_json` | TEXT | NOT NULL | `submission_id`, `file_id` |
| `dedupe_key` | TEXT |  | Opsional |
| `status` | TEXT | NOT NULL | `queued`, `running`, `done`, `failed`, `dead` |
| `attempts` | INTEGER | NOT NULL | Percobaan |
| `max_attempts` | INTEGER | NOT NULL | Default 5 |
| `run_after` | TEXT | NOT NULL | Backoff |
| `locked_by` | TEXT |  | Worker |
| `locked_at` | TEXT |  | Waktu lock |
| `heartbeat_at` | TEXT |  | Heartbeat |
| `last_error` | TEXT |  | Error terakhir |
| `progress_json` | TEXT |  | Progress UI |
| `created_at` | TEXT | NOT NULL | Dibuat |
| `finished_at` | TEXT |  | Selesai |

---

## 7. Nilai dan kalkulasi

### 7.1 `indicator_value`

Nilai dasar satu indikator dari satu baris submission.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID nilai |
| `submission_row_id` | TEXT | UNIQUE, FK → `submission_row.id` | Baris sumber |
| `indikator_id` | TEXT | FK → `indikator.id` | Indikator |
| `periode_id` | TEXT | FK → `periode.id` | Periode |
| `faskes_id` | TEXT | FK → `faskes.id` | Faskes |
| `scope_type` | TEXT | NOT NULL | Scope |
| `scope_id` | TEXT |  | ID scope |
| `numerator` | REAL | NOT NULL | Jangan dibulatkan |
| `denominator` | REAL | NOT NULL | Jangan dibulatkan |
| `value` | REAL | NOT NULL | `numerator/denominator×100` |
| `status` | TEXT | NOT NULL | `COMPLETE`, `PARTIAL`, `NO_DATA`, `STALE` |
| `data_version` | TEXT |  | Versi data |
| `created_at` | TEXT | NOT NULL | Waktu |

### 7.2 `calculation_run`

Satu eksekusi kalkulasi per periode/scope.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID run |
| `periode_id` | TEXT | FK → `periode.id` | Periode |
| `scope_type` | TEXT | NOT NULL | `KABUPATEN`/`FASKES` |
| `scope_id` | TEXT |  | ID scope |
| `trigger` | TEXT | NOT NULL | `submission.final`, `seed`, `target` |
| `status` | TEXT | NOT NULL | `running`, `done`, `failed` |
| `started_at` | TEXT | NOT NULL | Mulai |
| `finished_at` | TEXT |  | Selesai |
| `error` | TEXT |  | Error |

### 7.3 `iku_value`

Nilai IKU materialized. Dashboard membaca tabel ini, bukan menghitung ulang saat request.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID nilai IKU |
| `label_id` | TEXT | FK → `label.id` | Node IKU |
| `periode_id` | TEXT | FK → `periode.id` | Periode |
| `scope_type` | TEXT | NOT NULL | Scope |
| `scope_id` | TEXT |  | ID scope |
| `numerator` | REAL | NOT NULL | Σ numerator |
| `denominator` | REAL | NOT NULL | Σ denominator |
| `value` | REAL |  | `Σnum/Σden×100` |
| `status` | TEXT | NOT NULL | `COMPLETE`, `PARTIAL`, `NO_DATA`, `STALE` |
| `formula_version` | TEXT | NOT NULL | `v1` |
| `calculation_run_id` | TEXT | FK → `calculation_run.id` | Run |
| `source_hash` | TEXT |  | Hash file sumber |
| `data_version` | TEXT |  | `v2026-05-xxxx` |
| `updated_at` | TEXT | NOT NULL | Waktu |
| — | — | UNIQUE (`label_id`,`periode_id`,`scope_type`,`scope_id`) | Satu nilai per kombinasi |

### 7.4 `iku_value_detail`

Provenance komponen angka induk/anak.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID detail |
| `iku_value_id` | TEXT | FK → `iku_value.id` | Nilai IKU |
| `component_type` | TEXT | NOT NULL | `CHILD_IKU`, `INDIKATOR` |
| `component_id` | TEXT | NOT NULL | ID label/indikator |
| `component_label_id` | TEXT | FK → `label.id` | Label komponen |
| `component_value` | REAL |  | Nilai komponen |
| `numerator` | REAL |  | Numerator komponen |
| `denominator` | REAL |  | Denominator komponen |
| `weight` | REAL | NOT NULL | Bobot |
| `contribution` | REAL |  | Nilai × bobot |

---

## 8. Faktor, tindak lanjut, dan audit

### 8.1 `faktor`

Analisis kualitatif. **Bukan** input angka dan **bukan** TL; tidak mengubah kalkulasi.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID faktor |
| `iku_value_id` | TEXT | UNIQUE, FK → `iku_value.id` | Nilai IKU |
| `periode_id` | TEXT | FK → `periode.id` | Periode |
| `pendukung` | TEXT |  | Faktor pendukung |
| `penghambat` | TEXT |  | Faktor penghambat |
| `created_by` | TEXT | FK → `users.id` | Pengisi |
| `updated_at` | TEXT | NOT NULL | Waktu |

### 8.2 `tindak_lanjut`

Rencana aksi ber-PJ dan berstatus.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID TL |
| `iku_value_id` | TEXT | FK → `iku_value.id` | Nilai IKU terkait |
| `periode_id` | TEXT | FK → `periode.id` | Periode |
| `masalah` | TEXT | NOT NULL | Masalah |
| `rekomendasi` | TEXT | NOT NULL | Rekomendasi |
| `pj_user_id` | TEXT | FK → `users.id` | Penanggung jawab |
| `status` | TEXT | CHECK `open/in_progress/done/cancelled` | Status |
| `created_by` | TEXT | FK → `users.id` | Pembuat |
| `created_at` | TEXT | NOT NULL | Dibuat |
| `updated_at` | TEXT | NOT NULL | Diubah |

### 8.3 `audit_log`

Audit append-only. Aplikasi hanya `INSERT`; tidak ada endpoint update/delete.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `id` | TEXT | PK | UUID audit |
| `actor_id` | TEXT | FK → `users.id` | Aktor |
| `action` | TEXT | NOT NULL | Contoh `submission.finalized` |
| `entity_type` | TEXT | NOT NULL | `submission`, `submission_row`, `iku_value`, `faktor`, `tindak_lanjut` |
| `entity_id` | TEXT |  | ID entitas |
| `before_json` | TEXT |  | Snapshot sebelum |
| `after_json` | TEXT |  | Snapshot sesudah |
| `reason` | TEXT |  | Alasan/catatan |
| `request_id` | TEXT |  | Korelasi request |
| `impersonator_id` | TEXT | FK → `users.id` | Admin yang melakukan actingAs |
| `created_at` | TEXT | NOT NULL | Waktu |

### 8.4 `schema_migrations`

Internal migrasi.

| Kolom | Tipe | Constraint | Keterangan |
| --- | --- | --- | --- |
| `version` | INTEGER | PK | Versi migrasi |
| `applied_at` | TEXT | NOT NULL | Waktu apply |

---

## 9. Relasi utama

| Dari | Ke | Kardinalitas | Kolom |
| --- | --- | --- | --- |
| `users` | `positions` | N:1 | `position_id` |
| `auth_session` | `users` | N:1 | `user_id` |
| `users` | `faskes` | N:1 | `faskes_id` |
| `users` | `bidang` | N:1 | `bidang_id` |
| `faskes` | `wilayah` | N:1 | `wilayah_id` |
| `label` | `dashboard_category` | N:1 | `category_id` |
| `label` | `label` | N:1 | `parent_id` |
| `label` | `master_release` | N:1 | `master_release_id` |
| `bidang` | `bidang` | N:1 | `parent_id` |
| `indikator` | `bidang` | N:1 | `pj_timker_id` |
| `indikator_label` | `indikator` | N:1 | `indikator_id` |
| `indikator_label` | `label` | N:1 | `label_id` |
| `iku_komponen` | `label` | N:1 | `parent_label_id` |
| `target` | `label`/`indikator`/`periode` | N:1 | FK |
| `capaian_tahunan` | `label`/`indikator`/`users` | N:1 | FK |
| `dashboard_row` | `label`/`indikator` | N:1 | FK |
| `file_object` | `users` | N:1 | `uploaded_by` |
| `submission` | `users`/`faskes`/`periode`/`file_object` | N:1 | FK |
| `submission_revision` | `submission`/`file_object` | N:1 | FK |
| `submission_group` | `submission`/`bidang`/`label` | N:1 | FK |
| `submission_row` | `submission`/`indikator`/`label` | N:1 | FK |
| `submission_review` | `submission`/`users` | N:1 | FK |
| `indicator_value` | `submission_row`/`indikator`/`periode`/`faskes` | N:1 | FK |
| `iku_value` | `label`/`periode`/`calculation_run` | N:1 | FK |
| `iku_value_detail` | `iku_value`/`label` | N:1 | FK |
| `faktor` | `iku_value`/`periode`/`users` | N:1 | FK |
| `tindak_lanjut` | `iku_value`/`periode`/`users` | N:1 | FK |
| `audit_log` | `users` | N:1 | `actor_id` |

## 10. Index dan constraint penting

| Index/constraint | Tujuan |
| --- | --- |
| `ux_users_username` | Username login unik |
| `ix_label_category` | Filter label per kategori |
| `ix_category_urutan` | Urutan/aktif kategori dashboard |
| `ix_bidang_jenis` | Filter bidang/sekretariat/timker |
| `ix_indikator_timker` | Filter indikator per timker |
| `ux_capaian_tahunan_key` | Cegah duplikasi capaian per tahun/label/indikator/scope |
| `ix_capaian_tahunan_tahun` | Cari capaian per tahun/scope |
| `ix_auth_session_user` | Cari sesi per user |
| `ix_audit_impersonator` | Audit tindakan actingAs |
| `ux_submission_idem` | Cegah upload ganda dengan `Idempotency-Key` |
| `ix_submission_period` | Daftar submission per periode/status |
| `ix_submission_faskes` | Filter faskes/periode |
| `ix_row_submission` | Preview baris per submission/status |
| `ix_job_claim` | Klaim job worker |
| `ix_indicator_value_period` | Kalkulasi per periode/scope |
| `ix_iku_value_period` | Dashboard per periode/scope |
| `ix_detail_value` | Provenance per nilai IKU |
| `ix_audit_entity` | Riwayat audit per entitas |
| `UNIQUE (label_id, periode_id, scope_type, scope_id)` di `iku_value` | Idempoten kalkulasi |
| `UNIQUE (submission_row_id)` di `indicator_value` | Satu nilai dasar per baris |
| `UNIQUE (iku_value_id)` di `faktor` | Satu faktor per nilai |

## 11. Status dan enum

| Enum | Nilai |
| --- | --- |
| Role | `pkm_operator`, `pkm_pj`, `dinas_verifier`, `kabid`, `sekretaris`, `iku_owner`, `pimpinan`, `admin` |
| Category jenis | `LABEL`, `OWNER` |
| Bidang jenis | `bidang`, `sekretariat`, `timker` |
| Submission status | `draf`, `dikirim`, `diverifikasi`, `dikembalikan_petugas`, `dikembalikan_verifikator`, `final`, `diubah_khusus` |
| Workflow step | `DRAFT_PKM`, `VERIFY_DINAS`, `VALIDATE_KABID`, `FINAL`, `SPECIAL_CHANGE` |
| Classification | `MATCHED`, `MISMATCH`, `UNKNOWN_IKU`, `UNKNOWN_INDICATOR`, `IKU_ROLLUP_ONLY`, `AMBIGUOUS`, `DUPLICATE`, `INVALID_VALUE`, `SCOPE_DENIED` |
| Value status | `COMPLETE`, `PARTIAL`, `NO_DATA`, `STALE` |
| TL status | `open`, `in_progress`, `done`, `cancelled` |
| Aggregation | `RATIO_OF_SUMS`, `SUM`, `WEIGHTED_AVERAGE`, `FORMULA`, `NONE` |

## 12. Data seed pilot

Seed idempotent di `backend/internal/repository/sqlite/migrations/002_seed.sql` dan `demo_seed.go` mengisi:

- 7 pengguna awal (Siti, Rina, dr. Ahmad, Sekretaris, Kepala Dinas, Admin, PJ Gizi) dengan username dan password hash bcrypt.
- 7 jabatan struktural + tabel `auth_session`.
- Login awal: `admin`, `siti`, `rina`, `kabid`, `pimpinan`, `sekretaris`, `pjgizi` — password `password123` (ganti sebelum produksi).
- 2 faskes (PKM Taliwang, PKM Jereweh), wilayah contoh, dan bidang/timker dengan `jenis`.
- 2 periode (2026-05, 2026-06).
- 1 master release 2026 + workflow `default_3_level`.
- 8 `dashboard_category`: IKU, IKK, SPM, RIBK, Program, Kegiatan, Sub Kegiatan, dan Bidang/Sekretariat/Tim Kerja (`OWNER`). Kategori ke-9 tidak di-seed.
- Label seed di-backfill ke kategori IKU/SPM; label baru wajib memilih kategori.
- IKU `IKU-2026-001` (induk) dengan anak `.01`, `.02`, `.03`.
- 12 indikator SPM + 12 baris `dashboard_row` (3 contoh terhubung ke IKU, 9 baris kontrak).
- 12 baris `capaian_tahunan` 2025 hasil migrasi baseline JSON; dashboard 2026 memakainya sebagai pembanding.
- 1 submission final contoh `S-2026-05-001` dengan 3 baris, 3 `indicator_value`, review, dan audit.
- Kalkulasi seed menghasilkan `.01` 90,00%, `.02` 95,00%, `.03` 100,00%, induk 92,39%.

## 13. Alur trace yang dijamin

```text
dashboard_row
  → iku_value
  → iku_value_detail
  → indicator_value
  → submission_row
  → submission_revision
  → submission
  → file_object
```

Setiap `dashboard_row` membawa `iku_value_id`; setiap `iku_value` membawa `calculation_run_id`, `source_hash`, dan `data_version`. API trace: `GET /api/v1/dashboard/spm/{kode}/trace?periode=2026-05`.

## 14. Catatan SQLite → PostgreSQL

| SQLite pilot | PostgreSQL produksi |
| --- | --- |
| `TEXT` ID | `uuid` |
| `TEXT` waktu | `timestamptz` |
| `TEXT` JSON | `jsonb` |
| `INTEGER` boolean | `boolean` |
| `REAL` | `numeric`/`double precision` |
| `BEGIN IMMEDIATE` untuk job | `FOR UPDATE SKIP LOCKED` |
| Tanpa schema | schema `app` atau `public` |
| Audit tanpa hash chain | tambah `prev_hash`/`hash` sesuai `docs/04` |

Migrasi wajib mempertahankan nama kolom dan enum agar kontrak API/UI tidak berubah.

## 15. Keamanan dan batasan pilot

- **Auth PASETO v4.local.** Access token 30 menit; refresh token HttpOnly 7 hari disimpan sebagai hash di `auth_session`; role/scope dan actingAs admin divalidasi backend.
- **File tidak disajikan static.** Unduhan lewat `GET /api/v1/files/{id}`; pada produksi wajib ditambah cek auth/scope.
- **PII BNBA belum diputuskan.** Jangan menambah kolom PII baru tanpa keputusan Dinas; batasi akses dan masking pada produksi.
- **Audit belum hash chain.** Pilot mencatat `before_json`/`after_json`; produksi wajib menambah `prev_hash`/`hash`.
- **Belum ada dual approval UI.** `SPECIAL_CHANGE` dan final lock sudah dicatat, tetapi mekanisme dual approval formal menunggu keputusan Dinas.
- **Status dashboard.** Dihitung dari target periode: capaian ≥ target = sesuai target; di bawah target = perlu perhatian. Ambang warna khusus belum diputuskan.

## 16. Cara reset & seed

```bash
cd backend
rm -rf data
APP_WORKER_INLINE=true APP_ADDR=:8080 go run ./cmd/api
```

Database dibuat ulang di `backend/data/matasidik.db`; file demo ditulis ke `backend/data/files/2026/05/`. Seed hanya berjalan jika `S-2026-05-001` belum ada.
