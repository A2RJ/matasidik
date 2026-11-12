# Rencana — BNBA, Sebaran Wilayah & Dashboard Dinamis

**Status:** DUA TRACK — **Track A (BNBA/sebaran wilayah)** `BLOCKED` menunggu contoh file BNBA `.xlsx`; **Track B (dashboard dinamis)** `IN PROGRESS` — fondasi kategori, CRUD, dashboard, dan template sudah diimplementasikan.
**Mode:** Implementasi bertahap; file Excel asli `Proses Bisnis Mata Sidik.xlsx` tidak diubah.
**Tanggal dibuat:** 2026-11-12
**Sumber requirement:** `Proses Bisnis Mata Sidik.xlsx`, sheet **Contoh Data**, baris 39; ekstraksi di `docs/03-contoh-data.md:61`. Pengelompokan dashboard: `docs/01-gambaran-umum.md:77-89`; label indikator: `docs/02-daftar-indikator-kinerja.md:15-31`.
**Terkait:** `README.md:109`; `docs/rules.md:195-196, 362-368`; `docs/04-arsitektur-iku.md:54-97,169`.
**Track:** A = BNBA & sebaran wilayah; B = pengelompokan dashboard dinamis.

---

## 1. Tujuan

Memungkinkan Puskesmas mengunggah file Excel BNBA **sesuai format yang sudah mereka miliki**, lalu sistem:

1. Memetakan data ibu hamil berdasarkan alamat domisili: **desa → puskesmas → kecamatan**.
2. Menghitung cakupan layanan sesuai standar per wilayah.
3. Menampilkan wilayah dengan jumlah/rasio ibu hamil **tidak mendapat pelayanan sesuai standar** terbanyak.
4. Menjaga agar setiap angka tetap dapat ditelusuri sampai file dan baris sumber.

Requirement ini **sudah tercatat** di dokumen proyek, tetapi **belum diimplementasikan** di pilot.

> Plan ini punya dua track yang berbagi fondasi dashboard dan master data:
> - **Track A — BNBA & sebaran wilayah:** menunggu file BNBA dan keputusan denominator.
> - **Track B — dashboard dinamis:** menjadikan pengelompokan/filter 8 kategori sebagai data, bukan hardcode; dapat dimulai tanpa menunggu file BNBA.
>
> Keduanya baru boleh bertemu di Fase 2 (dashboard generik + scope wilayah) agar tidak membangun dua dashboard terpisah.

---

## 2. Fakta kondisi aplikasi saat ini

| Area | Kondisi | Bukti |
| --- | --- | --- |
| Format upload | Hanya `.csv` dan `.xlsx`; `.xls` ditolak | `backend/internal/service/submission.go:92-93` |
| Parser | Hanya sheet pertama; header wajib `kode_iku`, `kode_indikator`, `periode` | `backend/internal/worker/parser.go:84,106` |
| Baris file | Satu baris file = satu baris submission; belum ada agregasi per desa | `parser.go:113-131` |
| Klasifikasi | Wajib kode IKU/indikator; denominator wajib > 0; satu baris → satu indikator | `backend/internal/worker/classifier.go:62,97,146-151` |
| Scope | Hanya `FASKES` yang divalidasi; `DESA` diterima tanpa validasi master | `classifier.go:78-83` |
| Kalkulasi | Hanya `KABUPATEN` dan `FASKES`; belum `DESA`/`KECAMATAN` | `backend/internal/service/calculation.go:52-62` |
| Atribusi faskes | Memakai faskes pengupload, bukan faskes domisili | `classifier.go:147-148`; `value_repo.go:27-29` |
| Dashboard | Detail hanya per puskesmas; belum per desa | `backend/internal/service/dashboard.go:227-272`; `frontend/src/routes/dashboard/spm/[kode]/+page.svelte:132-178` |
| Filter dashboard | Hanya periode, faskes, status | `frontend/src/routes/+page.svelte:180-193` |
| Master wilayah | Hanya 2 kecamatan dan 2 desa contoh | `backend/internal/repository/sqlite/migrations/002_seed.sql:10-19` |
| Mapping desa–puskesmas | Belum ada; `faskes.wilayah_id` hanya lokasi faskes | `migrations/001_init.sql:26-33` |
| Level wilayah | `wilayah` hanya desa/kecamatan/kabupaten; puskesmas adalah `faskes` | `migrations/001_init.sql:22` |
| PII | Kebijakan belum diputuskan; belum ada masking/audit unduh | `docs/rules.md:397,492`; `handler/submission.go:150-163` |
| Input manual | Pilot hanya menerima file; tidak ada form input angka manual | `handler/submission.go` |

**Fondasi yang sudah bisa dipakai ulang:**

- `indicator_value` dan `iku_value` sudah punya `scope_type`/`scope_id` (`001_init.sql:239-286`).
- `wilayah` sudah punya hierarki desa → kecamatan → kabupaten.
- Background job, file hashing, workflow tiga level, audit, pohon IKU, `RATIO_OF_SUMS`, dan trace sampai `FILE_OBJECT` tidak perlu dirombak.

---

## 3. Kesenjangan utama

### P0 — Upload BNBA existing kemungkinan gagal

Parser mewajibkan tiga kolom kode. File BNBA Puskesmas belum tentu memilikinya. Parser juga hanya membaca sheet pertama dan tidak menangani title row/merged cell/multi-sheet.

### P0 — Desa/kecamatan belum menjadi unit analisis

Kalkulasi hanya membuat nilai KABUPATEN dan FASKES. Tidak ada `iku_value` per desa/kecamatan, sehingga dashboard tidak bisa menjawab “desa mana”.

### P0 — “Tidak terlayani” tidak bisa dihitung tanpa denominator

Jika BNBA hanya memuat ibu hamil yang sudah dilayani, sistem tidak punya total sasaran per desa. “Tidak terlayani” tidak boleh dianggap nol; harus `NO_DATA` sampai denominator tersedia.

### P0 — Risiko double counting dan salah atribusi

- Baris per desa + baris total faskes bisa dijumlahkan dua kali di level kabupaten.
- Nilai faskes memakai faskes pengupload, bukan faskes domisili/pelayanan.

### P1 — Satu baris BNBA bisa menghasilkan beberapa indikator

Model sekarang satu `submission_row` → satu `indicator_value`. Padahal satu ibu hamil bisa berkontribusi ke beberapa indikator (K1, K4, ANC, persalinan, dll.).

### P1 — Master dan definisi belum siap

- Master wilayah belum lengkap.
- Belum ada mapping desa → puskesmas wilayah kerja.
- `indikator` belum punya kolom `level_wilayah`, padahal diwajibkan `docs/rules.md:359`.
- Definisi “sesuai standar” masih kosong/terlalu umum.
- Ada 3 baris indikator ibu hamil identik di katalog (`docs/02:29-31`).

### P0 jika data individu — PII

BNBA line-level berpotensi memuat nama, NIK, alamat, tanggal lahir, dan data kehamilan. Kebijakan PII belum diputuskan. Rekomendasi: agregasi di ingest dan jangan simpan identitas individu.

---

## 4. Keputusan yang harus diambil sebelum coding

| # | Pertanyaan | Opsi | Dampak | Rekomendasi |
| --- | --- | --- | --- | --- |
| 1 | Apakah BNBA memuat seluruh ibu hamil atau hanya yang dilayani? | A. Seluruh ibu hamil + status layanan; B. Hanya yang dilayani; C. Campuran per Puskesmas | A: denominator bisa dihitung. B: butuh sasaran eksternal. C: profil harus menangani dua jenis file | **A**; jika B, siapkan sasaran per desa |
| 2 | Apa definisi “sesuai standar”? | A. Kolom status di BNBA; B. Hitung dari item layanan (K1/K4/10T dll.); C. Definisi baru Dinas | A: cepat, validitas tergantung Puskesmas. B: lebih auditabel. C: perlu waktu | **B** jika kolom item tersedia; A untuk pilot |
| 3 | “Puskesmas” di rantai domisili itu apa? | A. Puskesmas wilayah kerja domisili; B. Faskes yang melayani; C. Keduanya | A: cocok untuk “wilayah keberadaan”. B: menjawab layanan. C: paling informatif | **C**; default A untuk peta wilayah, B untuk performa layanan |
| 4 | Cukup agregat per desa atau perlu individu? | A. Agregat; B. Agregat + pseudonim; C. Data individu lengkap | A: risiko PII rendah. B: menengah. C: tinggi | **A** untuk fitur ini |
| 5 | Target/sasaran per desa tersedia? | A. Dinas menyediakan; B. Proporsional dari kabupaten; C. Tanpa target | A: akurat. B: asumsi berisiko. C: terbatas | **A**; C sementara |
| 6 | Output yang diinginkan? | A. Tabel ranking + bar chart; B. Peta; C. Keduanya | A: cepat. B: butuh GeoJSON/library. C: lebih lama | **A dulu**, peta setelah data valid |
| 7 | Tingkat dinamisme kategori | A. `label_type` saja; B. `label_type` + `dashboard_facet`; C. Hardcode 8 kategori; D. Satu tabel `dashboard_category` dengan CRUD | D: kategori, urutan, dan status aktif dikelola admin; A/B: butuh layer tambahan; C: kaku | **Diputuskan: D** — CRUD kategori; satu tabel dengan `jenis` `LABEL`/`OWNER` |
| 8 | Model Bidang/Timker | A. Satu tabel `bidang` dengan `jenis`; B. Tabel `timker` terpisah + FK; C. Tetap teks di indikator | A: perubahan lebih kecil; B: paling rapi untuk grouping; C: hanya layak pilot | **Diputuskan: A** — satu tabel `bidang` + `jenis` + `parent_id` |
| 9 | Arti kategori 8 | A. Daftar/filter indikator per owner; B. Satu angka komposit; C. Rata-rata persentase indikator | A: aman dan tidak menyesatkan; B: perlu formula eksplisit; C: dilarang karena merata-ratakan persentase | **Diputuskan: A** — daftar/filter per owner; bukan angka komposit |
| 10 | Indikator memilih label spesifik atau kategori langsung? | A. Label/IKU spesifik; kategori mengikuti label; B. Kategori langsung via `indikator_category`; C. Keduanya | A: satu sumber kebenaran, cocok dengan `indikator_label`; B: form lebih sederhana, tetapi menambah sumber kategori; C: paling fleksibel, paling kompleks | **Diputuskan: A** — label/IKU spesifik; kategori mengikuti label; root label per kategori untuk penandaan kategori saja |
| 11 | “Indikator Sasaran” diperlakukan bagaimana? | A. Kategori terpisah; B. Bagian dari kategori IKU; C. Belum diputuskan | A: filter terpisah; B: sesuai pengelompokan sheet 1 “IKU (Tujuan dan Sasaran)”; C: menunda | **Diputuskan: B** — Indikator Sasaran adalah bagian dari kategori IKU; bukan kategori filter terpisah |

**Keputusan Track B yang sudah dikonfirmasi (2026-11-12):**

- Kategori dashboard dapat di-CRUD admin; satu tabel `dashboard_category` dengan `jenis` (`LABEL`/`OWNER`), `urutan`, `aktif`, `is_system`.
- Setiap label/IKU wajib punya satu kategori; indikator memilih banyak label melalui `indikator_label`, sehingga bisa muncul di beberapa kategori.
- Indikator memilih **label/IKU spesifik**, bukan kategori langsung; kategori mengikuti label. Untuk penandaan kategori saja, gunakan root label per kategori.
- “Indikator Sasaran” adalah bagian dari kategori IKU, bukan kategori filter terpisah.
- Bidang/Timker memakai satu tabel `bidang` dengan `jenis` (`bidang`/`sekretariat`/`timker`) dan `parent_id`; indikator memakai `pj_bidang_id` dan `pj_timker_id`.
- Kategori 8 adalah daftar/filter per owner, bukan angka komposit/rata-rata.
- File Excel asli tidak diubah; implementasi dilakukan di kode dan dokumen proyek.

**Jangan mengisi keputusan ini dengan asumsi.** Jika belum ada jawaban, pekerjaan yang bergantung padanya harus berhenti.

---

## 5. Track A — Rencana implementasi BNBA & sebaran wilayah

### Fase 0 — Inventarisasi dan keputusan (belum coding)

1. Terima 1–3 contoh file BNBA `.xlsx`/`.xls` asli (anonim).
2. Catat: jumlah sheet, baris header, kolom domisili, kolom status layanan, jumlah baris, format tanggal/angka, merged cell.
3. Tentukan definisi “sesuai standar” dan sumber denominator.
4. Lengkapi master wilayah + mapping desa → puskesmas.
5. Putuskan PII dan tingkat detail output.
6. Baru buat desain profil BNBA dan kontrak API.

### Fase 1 — Profil BNBA + ingest agregat per desa

- Registry profil per Puskesmas/format: mapping header, kolom desa/kecamatan, status layanan, indikator turunan, versi.
- Parser: multi-sheet, deteksi header, dukungan `.xls` bila diperlukan, validasi domisili.
- Preview agregat per desa + ringkasan baris karantina (bukan ribuan baris individu).
- Simpan agregat per desa/indikator; file asli tetap untuk audit.
- Jangan simpan identitas individu kecuali disetujui.

### Fase 2 — Kalkulasi dan dashboard wilayah

- Tambah scope `DESA` dan `KECAMATAN` ke kalkulasi.
- Cegah double counting: satu indikator hanya diinput pada satu level, atau agregasi eksplisit.
- `source_hash` per scope dihitung dari file yang relevan.
- Endpoint dashboard wilayah + drill-down.
- UI: tabel ranking desa/kecamatan, bar chart, filter wilayah.
- Detail SPM: tambah cakupan per desa.

### Fase 3 — Hardening

- Dedupe lintas bulan bila diperlukan.
- Peta choropleth bila disetujui (butuh GeoJSON + library).
- Export per wilayah.
- Uji performa file besar.
- Update dokumen dan diagram (stack-diagram, bukan Mermaid).

### Roadmap gabungan

| Fase | Track A — BNBA & wilayah | Track B — Dashboard dinamis | Dependensi |
| --- | --- | --- | --- |
| 0 | Inventaris file BNBA + keputusan denominator | Audit label + keputusan kategori/owner | — |
| 1 | Profil BNBA + ingest agregat desa | `dashboard_category` + owner/bidang + master UI | Fase 0 |
| 2 | Kalkulasi DESA/KECAMATAN + geo view | Dashboard generik kategori + filter owner | Fase 1A & 1B |
| 3 | Geo dashboard + hardening | Hardening taksonomi + audit | Fase 2 |

---

## 5A. Track B — Pengelompokan dashboard dinamis

### Tujuan

- 8 kategori menjadi data master, bukan daftar hardcode di frontend/backend.
- Admin dapat menambah, mengganti nama, mengatur urutan, atau menonaktifkan kategori tanpa deploy.
- Kategori ke-9 dapat ditambahkan setelah Dinas memutuskan; jangan diisi sekarang.
- Bidang/Timker diperlakukan sebagai dimensi kepemilikan (PJ), bukan label taksonomi.
- Dashboard memakai satu mekanisme filter untuk semua kategori, dengan bentuk query yang sesuai jenisnya.

### Prinsip

1. Pisahkan **kategori dashboard** (`dashboard_category`), **label** (node spesifik), dan **owner** (bidang/timker).
2. `dashboard_category.kode` adalah kontrak stabil; `nama` dan `urutan` boleh berubah.
3. Kategori sistem (8 kategori) tidak boleh dihapus; gunakan `aktif = 0`.
4. Kategori 8 memakai `jenis = OWNER`, bukan label.
5. Jangan memaksakan semua kategori memakai `dashboard_row` SPM; `dashboard_row` tetap kontrak khusus SPM.
6. Kategori 8 adalah daftar/filter per owner; jangan merata-ratakan persentase.
7. CRUD kategori wajib menulis audit; kategori yang sudah dipakai tidak boleh dihapus fisik.

### Desain data (usulan, perlu persetujuan)

| Tabel | Field kunci | Fungsi |
| --- | --- | --- |
| `dashboard_category` | `id`, `kode`, `nama`, `jenis` (`LABEL`/`OWNER`), `urutan`, `aktif`, `is_system`, `config_json` | Kategori dashboard yang dapat di-CRUD; 8 seed |
| `label` | `category_id` FK, `parent_id`, ... | Node taksonomi; ganti `tipe` teks |
| `bidang` | `id`, `kode`, `nama`, `jenis` (`bidang`/`sekretariat`/`timker`), `parent_id`, `aktif` | Dimensi owner; timker adalah child bidang |
| `indikator` | `pj_bidang_id` FK, `pj_timker_id` FK | PJ data; keduanya menunjuk `bidang` |

Catatan:

- Satu tabel `dashboard_category` cukup untuk MVP; tidak perlu `dashboard_facet` terpisah. Jika nanti nama filter perlu berbeda dari nama kategori, tambahkan layer baru saat itu.
- Jangan mencampur `bidang` dan `timker` tanpa `jenis`; seed sekarang mencampurnya (`bid-gizi` = “Timker Gizi”).
- Kategori `jenis = OWNER` tidak punya label; query memakai `bidang`/`timker`.
- `label.parent_id` hanya hierarki tampilan; rollup Program/Kegiatan/Sub Kegiatan tetap butuh `iku_komponen`.
- `label.tipe` lama dimigrasikan ke `category_id`; nilai yang tidak dikenal dilaporkan, bukan dibuat otomatis.

### API dan UI

- `GET /api/v1/meta/bootstrap` mengirim `categories`, `bidang` (dengan `jenis`), `labels`, dan `indicators`.
- Admin CRUD: `/admin/categories` dan `/admin/bidang`; kategori sistem tidak boleh dihapus, hanya dinonaktifkan.
- Dashboard: `GET /api/v1/dashboard?category=IKK&periode=...&scope_type=...` atau endpoint setara; `category=SPM` tetap memakai `dashboard_row`.
- `GET /api/v1/dashboard/owners?bidang_id=...&timker_id=...` untuk kategori 8.
- Frontend membangun opsi filter dari bootstrap, bukan dari array hardcode.
- Form indikator: multi-select label + pilih label utama + PJ bidang/timker.
- Form label: pilih `category` dari master; bukan teks bebas.
- Master UI: tab **Kategori** (CRUD) dan tab **Bidang/Timker** (CRUD).

### Penempatan kategori pada IKU dan indikator

| Objek | Relasi | Form/UI | Contoh |
| --- | --- | --- | --- |
| `label` (node IKU/taksonomi) | satu `category_id` | Dropdown **Kategori** saat tambah/edit IKU/label; wajib diisi | `IKU-2026-001` → kategori `IKU` |
| `indikator` | banyak label via `indikator_label` | Multi-select **Label Indikator** dikelompokkan per kategori + pilih label utama | Indikator ibu hamil → label IKU `.01`, IKK, RIBK |
| `indikator` | `pj_bidang_id`, `pj_timker_id` | Dropdown PJ terpisah dari label | Bidang Kesmas, Timker Kesga |
| Dashboard | turunan dari label | Filter kategori | Filter `IKK` menampilkan label/indikator kategori IKK |

Aturan:

- Setiap label/IKU **wajib** punya satu kategori; jangan simpan label tanpa kategori.
- Satu indikator boleh punya banyak label, sehingga muncul di banyak kategori; nilainya tetap **satu** di `indicator_value`.
- “Label Indikator” di sheet 2 dipetakan ke `indikator_label`; kategori berasal dari label yang dipilih.
- Jika Dinas hanya ingin menandai kategori tanpa label spesifik, gunakan **root label per kategori** (mis. `IKK-2026`) yang dibuat admin. Jangan membuat label otomatis dari file upload.
- `indikator_category` (kategori langsung tanpa label) **tidak dipakai**; keputusan sudah dikunci ke label/IKU spesifik.
- “Indikator Sasaran” adalah label di dalam kategori IKU; bukan kategori filter terpisah.

### Fase B

| Fase | Isi | Gate |
| --- | --- | --- |
| B0 | Audit label sheet 2; putuskan isi kategori/owner; audit data kosong | Dinas menyetujui daftar kategori dan definisi |
| B1 | Migrasi `label.tipe` → `category_id`; seed 8 kategori; CRUD master; bidang + `jenis` | Master bisa dikelola tanpa deploy |
| B2 | Dashboard generik kategori + owner; filter dinamis; empty state | 8 kategori tampil dari data, bukan hardcode |
| B3 | Audit, deaktivasi, migrasi, dokumentasi | Tidak ada kategori liar; histori tetap tampil |

### Acceptance criteria Track B

- Menambah/mengubah kategori lewat admin langsung muncul di filter tanpa deploy.
- Menonaktifkan kategori menyembunyikannya dari pilihan baru, tetapi label historis tetap tampil.
- Kategori sistem tidak bisa dihapus fisik; CRUD tercatat di audit.
- 8 kategori di-seed; kategori ke-9 tidak diisi sampai diputuskan.
- Label/IKU tidak bisa disimpan tanpa kategori; tabel IKU menampilkan kategori.
- Form indikator bisa memilih banyak label yang dikelompokkan per kategori dan menentukan label utama.
- Indikator muncul di semua kategori label yang dipilih; nilainya tetap satu.
- Pemetaan “Label Indikator” sheet 2 dapat direpresentasikan tanpa mencampur Bidang/Timker.
- Kategori 8 memakai dimensi owner, bukan label; tidak ada rata-rata persentase.
- Dashboard SPM/IKU yang ada tetap berjalan.
- Tidak ada daftar 8 kategori hardcode di frontend.

---

## 6. Checklist saat file BNBA `.xlsx` diterima

- [ ] Buka semua sheet, catat nama dan kegunaannya.
- [ ] Identifikasi baris header sebenarnya.
- [ ] Identifikasi kolom: nama/identitas, alamat, desa, kecamatan, puskesmas, status layanan, tanggal.
- [ ] Cek apakah ada kolom “sesuai standar” atau item layanan terpisah.
- [ ] Cek apakah ada baris total/agregat yang bisa menyebabkan double counting.
- [ ] Cek apakah ada data ibu hamil yang belum dilayani.
- [ ] Cek variasi nama desa/kecamatan/puskesmas.
- [ ] Cek format `.xls` vs `.xlsx`, merged cell, formula, dan jumlah baris.
- [ ] Anonimkan file sebelum disimpan sebagai fixture tes.
- [ ] Bandingkan dengan master wilayah yang ada.

---

## 7. Acceptance criteria

- File BNBA format Puskesmas bisa diunggah tanpa mengubah format file.
- Baris tanpa domisili valid ditolak/karantina, bukan ditebak.
- Setiap desa/kecamatan memiliki nilai cakupan yang bisa ditelusuri ke file dan baris sumber.
- “Tidak terlayani” dihitung dari denominator yang jelas; jika tidak ada denominator → `NO_DATA`, bukan nol.
- Tidak ada double counting antara level desa, faskes, dan kabupaten.
- Tidak ada identitas individu yang tersimpan/terekspos tanpa keputusan PII.
- Dashboard bisa menjawab: desa/kecamatan mana, berapa ibu hamil, berapa yang terlayani standar, berapa yang tidak, dan dari file mana.
- Data final tetap terkunci; perubahan lewat revisi/audit.
- Tidak ada Mermaid/ASCII flow; diagram baru memakai stack-diagram.
- Build/test relevan dijalankan sebelum menyatakan selesai.

---

## 8. File yang kemungkinan tersentuh nanti

| File | Perubahan |
| --- | --- |
| `backend/internal/worker/parser.go` | Multi-sheet, deteksi header, mapping profil |
| `backend/internal/worker/classifier.go` | Klasifikasi per desa, validasi scope, profil |
| `backend/internal/service/submission.go` | Profil upload, preview agregat |
| `backend/internal/service/calculation.go` | Scope DESA/KECAMATAN, anti-double-count |
| `backend/internal/service/dashboard.go` | Dashboard wilayah + query generik per kategori/owner |
| `backend/internal/domain/models.go`, `service/master.go` | Model `DashboardCategory`, CRUD kategori/bidang |
| `backend/internal/repository/sqlite/migrations/` | Tabel profil/coverage, `dashboard_category`; kolom `bidang.jenis`/`parent_id`, `indikator.pj_timker_id`; backfill `label.tipe` |
| `backend/internal/repository/sqlite/*_repo.go` | Query agregasi wilayah |
| `backend/internal/httpapi/handler/` | Endpoint dashboard wilayah, kategori, dan profil |
| `frontend/src/routes/upload/+page.svelte` | Pilih/deteksi profil, preview agregat |
| `frontend/src/routes/+page.svelte` | Filter/tab wilayah + filter kategori dinamis |
| `frontend/src/routes/dashboard/spm/[kode]/+page.svelte` | Tabel per desa/kecamatan |
| `frontend/src/routes/master/+page.svelte` | Tab Kategori (CRUD) + Bidang/Timker; form multi-label + PJ |
| `frontend/src/routes/admin/+page.svelte` | Kelola profil BNBA & coverage |
| `docs/02`, `docs/schema.md`, `docs/implementation.md`, `docs/04` | Taksonomi, skema/API, dan alur baru |
| `assets/diagrams/` | Diagram alur BNBA → agregasi desa (stack-diagram) |
| `backend/internal/worker/*_test.go`, fixtures | Tes format BNBA asli |

---

## 9. Guardrail

- Jangan mengubah data final langsung.
- Jangan mengarang definisi “sesuai standar”, target per desa, atau denominator.
- Jangan menambah tabel/API/enum tanpa persetujuan.
- Jangan menambah dependency peta/library tanpa keputusan.
- Jangan mengisi kategori dashboard ke-9 dengan “wilayah”; ini filter/drill-down SPM, bukan label baru.
- Jangan mengisi kategori ke-9 dengan tebakan; dynamic capability tidak berarti datanya boleh dikarang.
- Jangan hardcode daftar 8 kategori di frontend/backend; ambil dari master.
- Jangan menghapus kategori sistem; gunakan `aktif = 0` agar histori tetap tampil.
- CRUD kategori dan bidang/timker wajib menulis audit; jangan mengubah data master tanpa jejak.
- Jangan menyimpan label/IKU tanpa kategori; kategori adalah bagian dari kontrak taksonomi.
- Jangan mencampur label/category dengan PJ Bidang/Timker; keduanya dimensi berbeda.
- Jangan mengubah file Excel asli `Proses Bisnis Mata Sidik.xlsx`; perubahan hanya di plan/docs setelah disetujui.
- Jangan merata-ratakan persentase lintas indikator untuk kategori owner; gunakan daftar/filter atau formula eksplisit.
- Jangan memakai Mermaid/Graphviz/ASCII untuk flow.
- Jangan commit/push sebelum diminta.
- Jangan menghapus jalur upload agregat yang sudah ada; fitur BNBA baru harus bisa hidup berdampingan.

---

## 10. Status berikutnya

**Track A berhenti di sini** sampai ada contoh file BNBA `.xlsx` dan keputusan denominator/definisi pada tabel bagian 4.

**Track B sudah berjalan.** Progress implementasi:

- [x] Migrasi `dashboard_category` + seed 8 kategori.
- [x] `label.category_id` dan backfill IKU/SPM.
- [x] `bidang.jenis`/`parent_id` dan `indikator.pj_timker_id`.
- [x] CRUD kategori dan bidang/timker (API + master UI).
- [x] Form IKU/label memilih kategori.
- [x] Form indikator multi-label + label utama + PJ bidang/timker.
- [x] Tombol header **Pratinjau Pemetaan** (indikator + label/kategori + PJ), terpisah dari tab CRUD.
- [x] Filter dashboard per kategori + opsi **Semua** + query owner.
- [x] Template BNBA menambahkan `kode_kategori` opsional dan validasi classifier.
- [ ] Audit/isi label sheet 2 untuk kategori selain IKU/SPM.
- [ ] Target per kategori/label selain SPM.
- [ ] Regenerasi diagram skema setelah perubahan tabel.

Urutan lanjutan:

1. **B0** — audit label sheet 2, isi kategori yang masih kosong.
2. **A0/A1** — saat file BNBA tiba, jalankan inventaris dan ingest.
3. **Fase 2** — geo view BNBA di atas dashboard kategori yang sudah ada.
4. **Fase 3** — hardening dan regenerasi diagram/dokumen.

Setelah itu, perbarui plan ini dengan temuan format aktual dan keputusan Dinas.
