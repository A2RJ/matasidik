# Rencana — BNBA, Sebaran Wilayah & Dashboard Dinamis

**Status:** DUA TRACK — **Track A (BNBA/sebaran wilayah)** `IN PROGRESS` — template, parser multi-sheet, count/scalar, dan periode tahunan sudah diimplementasikan; keputusan denominator/definisi dan target per wilayah masih terbuka. **Track B (dashboard dinamis)** `IN PROGRESS` — fondasi kategori, CRUD, dashboard, dan template sudah diimplementasikan.
**Mode:** Implementasi bertahap; file Excel asli `Proses Bisnis Mata Sidik.xlsx` tidak diubah.
**Tanggal dibuat:** 2026-11-12
**Pembaruan:** 2026-11-20 — analisis format matriks desa × IKU, edge case sampai kalkulasi dashboard, template v1 review [`assets/templates/bnba/arsip/template-bnba-final.xlsx`](../assets/templates/bnba/arsip/template-bnba-final.xlsx), contoh v1 [`assets/templates/bnba/arsip/contoh-bnba-terisi.xlsx`](../assets/templates/bnba/arsip/contoh-bnba-terisi.xlsx), template v2 terpisah di `assets/templates/bnba/arsip/v2/`, template dinamis yang diterima aplikasi [`assets/templates/bnba/template-bnba-dinamis.xlsx`](../assets/templates/bnba/template-bnba-dinamis.xlsx), serta satu template final semua satuan & periode [`BNBA-Semua-Satuan.xlsx`](../assets/templates/bnba/final/BNBA-Semua-Satuan.xlsx) (§10.14–§10.19). Tidak ada template tahunan terpisah; data tahunan memakai `periode=YYYY`. Alur upload → dashboard rasio, count/scalar, dan periode tahunan sudah diimplementasikan; format harian baru tersedia di template, belum dihitung.
**Sumber requirement:** `Proses Bisnis Mata Sidik.xlsx`, sheet **Contoh Data**, baris 39; ekstraksi di `docs/03-contoh-data.md:61`. Pengelompokan dashboard: `docs/01-gambaran-umum.md:77-89`; label indikator: `docs/02-daftar-indikator-kinerja.md:15-31`.
**Terkait:** `README.md:109`; `docs/rules.md:195-196, 362-368`; `docs/04-arsitektur-iku.md:54-97,169`.
**Track:** A = BNBA & sebaran wilayah; B = pengelompokan dashboard dinamis.
**Plan final konsolidasi:** [`06-rencana-final.md`](06-rencana-final.md).

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
| Parser | Dulu hanya sheet pertama; sekarang membaca `Rasio`/`Nilai`/`Data`; header wajib `kode_iku`, `kode_indikator`, `periode` | `backend/internal/worker/parser.go`; §10.18 |
| Baris file | Satu baris file = satu baris submission; belum ada agregasi per desa | `parser.go:113-131` |
| Klasifikasi | Wajib kode IKU/indikator; denominator wajib > 0; satu baris → satu indikator | `backend/internal/worker/classifier.go:62,97,146-151` |
| Scope | `DESA`/`KECAMATAN`/`KABUPATEN` divalidasi dan di-resolve ke master; `FASKES` hanya legacy | `classifier.go` scope resolution |
| Kalkulasi | `KABUPATEN`, `KECAMATAN`, `DESA`, dan `FASKES` pengunggah | `calculation.go` hierarchy |
| Atribusi faskes | Memakai faskes pengupload, bukan faskes domisili | `classifier.go:147-148`; `value_repo.go:27-29` |
| Dashboard | Detail per faskes dan per wilayah (desa/kecamatan/kabupaten) | `dashboard.go`; halaman detail SPM |
| Filter dashboard | Periode, kategori, faskes, wilayah, status | `frontend/src/routes/+page.svelte` |
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

Parser mewajibkan tiga kolom kode. File BNBA Puskesmas belum tentu memilikinya. Parser sekarang membaca beberapa sheet data, tetapi belum menangani title row/merged cell.

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
- Mapping desa → puskesmas domisili sudah tersedia (§10.19); master wilayah masih 2 desa contoh.
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
4. Lengkapi master wilayah; mapping desa → puskesmas sudah tersedia (§10.19).
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

## 10. Analisis format matriks desa × IKU (referensi 2026-11-19)

### 10.1 Status dan batas analisis

Bagian ini menganalisis **skenario format**, bukan file BNBA resmi. Gambar yang dilampirkan adalah **referensi tampilan**, bukan format final dan bukan data produksi.

| Hal | Status | Catatan |
| --- | --- | --- |
| Gambar 1 | **Fakta referensi** | Rekap per kecamatan; header dua tingkat: Total Ibu, Ibu (Hamil/Bersalin/Nifas), Meninggal (Hamil/Bersalin/Nifas), Bayi (Kelahiran/Kematian); nilai contoh 147, 115, 0, dst. |
| Skenario teks pengguna | **Asumsi yang dianalisis** | Baris = desa; kolom = kode IKU; satu sel = satu angka untuk IKU tersebut. Bentuk ini berbeda dari gambar. |
| File BNBA asli | **Belum diterima** | Track A tetap `BLOCKED`; analisis ini tidak menggantikan inventaris file asli. |
| Angka pada gambar | **Ilustrasi** | Tidak dipakai sebagai data produksi atau fixture tanpa verifikasi. |
| Periode | **Belum diketahui** | Apakah satu file = satu bulan, kumulatif tahun berjalan, atau snapshot? |
| Arti “Total Ibu” | **Belum diketahui** | Apakah sasaran ibu hamil, ibu bersalin, atau populasi perempuan? Tidak otomatis menjadi denominator semua kolom. |
| Definisi “sesuai standar” | **Belum diputuskan** | Lihat §4 keputusan #2; jangan menurunkan rumus dari judul kolom. |
| PII | **Belum diputuskan** | Format agregat per desa lebih rendah risikonya; data individu tetap terlarang sampai ada keputusan. |

**Kesimpulan awal:** format matriks **bisa dipetakan** ke dashboard per desa/puskesmas/periode, tetapi **tidak bisa menjadi kontrak final apa adanya**. Alasannya bukan sekadar parsing: satu sel tidak membawa numerator, denominator, satuan, periode, dan identitas indikator yang dibutuhkan engine IKU. Detailnya di §10.3–§10.9.

---

### 10.2 Dua bentuk data yang perlu dibedakan

| Aspek | Matriks lebar (wide) | Baris panjang (long/tidy) |
| --- | --- | --- |
| Bentuk | Baris = desa; kolom = IKU; sel = angka | Satu baris = satu kombinasi periode + wilayah + IKU/indikator + nilai |
| Ramah dilihat manusia | Tinggi; mirip rekap Puskesmas | Sedang; lebih mirip tabel data |
| Identitas sel | Hanya dari posisi kolom | Eksplisit di `kode_iku`, `kode_indikator`, `scope_type`, `scope_code`, `periode` |
| Numerator/denominator | Tidak ada tempat eksplisit | Dua kolom terpisah |
| Satuan | Tercampur dalam header | Dari master indikator, bukan tebakan dari sel |
| IKU baru | Menambah kolom; template lama rusak | Menambah baris; template lama tetap terbaca |
| Total/subtotal | Mudah ikut terjumlah dua kali | Bisa dilarang dan divalidasi per baris |
| Trace ke sel | Sulit; perlu simpan sheet/baris/kolom | Langsung per baris; tetap perlu simpan sel asal untuk audit |
| Cocok untuk template final | Tidak, kecuali dengan profil khusus | **Ya** |

**Rekomendasi:** kontrak upload final memakai **long/tidy**. Tampilan matriks tetap boleh ada sebagai **tampilan UI** atau **format impor sementara** lewat profil, tetapi bukan format penyimpanan dan bukan satu-satunya format upload.

---

### 10.3 Apakah format pada gambar bisa dipetakan?

Bisa, dengan syarat berikut. Jika salah satu tidak terpenuhi, baris/sel harus masuk karantina, bukan ditebak.

1. **Kode, bukan hanya nama kolom.** Header `IKU-ibu-meninggal` atau `Ibu Hamil` tidak cukup jika tidak ada kode master yang stabil. Gunakan `kode_iku` dan `kode_indikator` resmi.
2. **Satu sel harus jelas perannya:** numerator, denominator, atau nilai jadi. `Total Ibu = 147` bisa menjadi denominator untuk satu indikator, tetapi belum tentu untuk indikator lain.
3. **Wilayah memakai kode.** Nama desa saja tidak cukup; desa bisa bernama sama di kecamatan berbeda, berubah nama, atau ditulis dengan variasi.
4. **Periode eksplisit.** Satu file = satu bulan, atau ada kolom `periode`. Jangan menaruh 12 bulan sebagai 12 kolom dalam satu file yang diproses sebagai satu submission.
5. **Total/subtotal dipisahkan.** Baris `Jumlah`, `Total`, atau kolom total tidak boleh ikut dihitung.
6. **Master wilayah dan mapping desa → puskesmas domisili tersedia.** Mapping berversi dipakai untuk agregasi faskes; `faskes.wilayah_id` bukan mapping desa.

Contoh masalah denominator: pada gambar, `Total Ibu` 147 dan `Hamil` 115 terlihat seperti denominator/numerator. Tetapi `Meninggal Bersalin` 0 belum tentu memakai denominator 147; bisa jadi denominatornya jumlah persalinan, jumlah kelahiran hidup, atau populasi. **Jangan memakai satu kolom Total Ibu sebagai denominator universal.** Definisi numerator/denominator tiap indikator harus ada di master atau profil.

Jika sel berisi persentase jadi (mis. `78,2%`), sel itu **tidak boleh** menjadi satu-satunya sumber untuk IKU induk, karena `RATIO_OF_SUMS` membutuhkan numerator dan denominator. Menyimpan persentase saja melanggar aturan `docs/rules.md` §7.1.

---

### 10.4 Rantai pemetaan ke dashboard

Setiap sel matriks harus menjadi satu baris data setelah proses *unpivot*:

`SEL EXCEL → SUBMISSION_ROW → INDICATOR_VALUE → IKU_VALUE (leaf) → IKU_VALUE (induk) → baris dashboard`

Agar rantai ini utuh, `SUBMISSION_ROW` perlu menyimpan lokasi asal sel (sheet, baris, kolom) selain `row_no`. Saat ini yang tersedia baru `row_no`; lihat §10.9.

| Dimensi dashboard | Data yang dibutuhkan | Cara memperoleh |
| --- | --- | --- |
| Per desa | `scope_type = DESA`, `scope_code` = kode desa resmi | Simpan di `indicator_value.scope_type/scope_id`; validasi ke master `wilayah` |
| Per puskesmas (domisili) | Mapping desa → puskesmas wilayah kerja | Tabel mapping berversi; **bukan** `faskes.wilayah_id` (itu lokasi faskes) |
| Per puskesmas (pelayanan) | Kode faskes yang melayani | Perlu kolom terpisah atau satu file per faskes pelayanan; jangan mencampur dengan domisili |
| Per bulan | `periode_id` | Satu file = satu bulan; `periode` wajib `YYYY-MM` |
| Per triwulan | Tiga periode bulanan | Jumlahkan numerator/denominator, baru hitung rasio; bukan rata-rata persentase |
| Per tahun | Dua belas periode bulanan | Sama seperti triwulan; bulan yang tidak ada = `PARTIAL`, bukan nol |
| IKU leaf | `kode_iku` leaf + `kode_indikator` yang terhubung | Klasifikasi deterministik; induk/rollup ditolak |
| IKU induk | `IKU_KOMPONEN` eksplisit | `RATIO_OF_SUMS` atau `SUM` sesuai metode; bukan rata-rata persentase |

**Catatan penting:** jika hanya tersedia baris kecamatan (seperti gambar), sistem **tidak boleh** membagi angka kecamatan ke desa secara proporsional tanpa keputusan Dinas. Tampilkan sampai kecamatan, dan tandai desa sebagai `NO_DATA`.

---

### 10.5 Edge case dari file sampai kalkulasi IKU

Status: `P0` = harus ditutup sebelum pilot BNBA; `P1` = hardening; `Keputusan` = menunggu Dinas/pemilik proyek.

#### A. Struktur file dan format

| # | Edge case | Dampak | Penanganan yang diusulkan | Status |
| --- | --- | --- | --- | --- |
| A1 | Header dua tingkat / merged cell seperti gambar | Parser hanya membaca satu baris header; kolom `Hamil` bisa hilang atau salah nama | Profil menentukan `header_rows`; isi merged cell ke bawah; validasi hasil | P0 |
| A2 | Baris judul/catatan sebelum header | Baris judul dianggap header; file gagal | Deteksi header atau `header_row` di profil; jangan menebak | P0 |
| A3 | Banyak sheet (Petunjuk/Data/Rekap) | Dulu hanya sheet pertama dibaca; bisa membaca rekap bukan data | Parser memilih `Rasio`/`Nilai`/`Data`; sheet ringkasan diabaikan (§10.18) | Selesai |
| A4 | Baris/kolom total, subtotal, `Jumlah` | Double counting di kabupaten | Tandai dan kecualikan; validasi total untuk peringatan, bukan untuk penjumlahan | P0 |
| A5 | Dua kolom dengan kode IKU sama | Nilai ganda / ambigu | Duplikat mapping ditolak; karantina | P0 |
| A6 | IKU baru/tahun baru menambah kolom | Template lama rusak; kolom bergeser | Long/tidy sebagai kontrak; wide hanya lewat profil berversi | P0 |
| A7 | Formula, hidden row/column, filter | Nilai cached/blank; data tidak terlihat | Baca nilai cached; tolak formula yang tidak bisa diverifikasi; laporkan hidden row | P1 |
| A8 | `.xls`, `.xlsm`, file password | Ditolak parser; macro tidak aman | Dukung `.xls` bila diperlukan; tolak macro/password dengan pesan jelas | P1 |
| A9 | Satu file berisi semua Puskesmas | Model submission satu faskes; approval per faskes | Pecah per faskes atau tolak; jangan memakai satu faskes untuk semua baris | P0 |
| A10 | CSV delimiter/encoding/koma desimal | Kolom bergeser; angka salah | Normalisasi encoding; delimiter eksplisit; parser angka sadar locale | P1 |

#### B. Identitas wilayah dan puskesmas

| # | Edge case | Dampak | Penanganan yang diusulkan | Status |
| --- | --- | --- | --- | --- |
| B1 | Nama desa bervariasi/typo | Baris tidak cocok atau salah desa | Wajib kode desa; nama hanya untuk verifikasi; alias hanya lewat master | P0 |
| B2 | Nama desa sama di kecamatan berbeda | Data masuk desa yang salah | Kunci `kode_wilayah`; validasi parent kecamatan | P0 |
| B3 | Baris kecamatan diberi label desa | Scope salah; drill-down menyesatkan | `level_wilayah` eksplisit; validasi level dengan master | P0 |
| B4 | `Tanpa Kecamatan`, `Lainnya`, baris kosong | Wilayah fiktif; total kacau | Karantina; jangan membuat wilayah baru | P0 |
| B5 | Desa dimekarkan/diganti nama | Data historis pindah wilayah | Mapping berversi dengan `berlaku_dari/sampai`; histori tetap | P1 |
| B6 | Pasien domisili A dilayani Puskesmas B | Double count atau atribusi salah | Simpan domisili dan faskes pelayanan sebagai dimensi terpisah; definisikan dedupe | Keputusan |
| B7 | Desa dilayani lebih dari satu Puskesmas | Dua file menjumlahkan angka yang sama | Aturan mapping/overlap; karantina jika ambigu | Keputusan |
| B8 | Master wilayah baru 2 desa dan 2 faskes | Hampir semua baris nyata akan unknown | Lengkapi master resmi sebelum pilot; jangan mengisi kode tebakan | P0 |
| B9 | Mapping desa → puskesmas | Tidak bisa menjawab per-puskesmas domisili | Tabel `desa_faskes_mapping` berversi sudah dibuat (§10.19); `faskes.wilayah_id` bukan mapping | Selesai |

#### C. Periode

| # | Edge case | Dampak | Penanganan yang diusulkan | Status |
| --- | --- | --- | --- | --- |
| C1 | `periode` kosong/tidak cocok dengan submission | Baris invalid | Wajib `YYYY-MM`; harus sama dengan submission | P0 |
| C2 | Sel tanggal Excel (`01/05/2026`) | Parser gagal | Normalisasi tanggal; validasi setelah parse | P1 |
| C3 | Upload triwulan/tahun dicampur dengan bulanan | Double counting | Input resmi bulanan; triwulan/tahun diturunkan sistem; jika terpaksa, `periode_type` dan tidak boleh dicampur | P0 |
| C4 | Data kumulatif vs transaksi bulanan | Penjumlahan tahunan salah | Tandai sifat data per indikator (`snapshot`/`incremental`); jangan menebak | Keputusan |
| C5 | Bulan yang sama diupload dua kali | Nilai dobel | Revisi/supersede, bukan penjumlahan; satu final aktif per kunci | P0 |
| C6 | Satu sheet berisi 12 bulan | Satu submission hanya satu periode | Satu file = satu bulan; atau pecah menjadi 12 submission | P0 |

#### D. Nilai dan satuan

| # | Edge case | Dampak | Penanganan yang diusulkan | Status |
| --- | --- | --- | --- | --- |
| D1 | Sel kosong vs `0` | `0` bisa berarti nol kasus; kosong bisa berarti belum lapor | Blank = `NO_DATA`/karantina; `0` = nol yang dilaporkan; jangan isi otomatis | P0 |
| D2 | `1.234,56` / `1,234.56` / `1.234` | Angka salah 1000× atau gagal parse | Sel numerik; parser locale eksplisit; larang pemisah ribuan di template | P1 |
| D3 | `45%` vs `0,45` vs `45` | Selisih 100× | Normalisasi berdasarkan satuan master; simpan nilai mentah | P0 |
| D4 | Count (45 kasus) vs ratio (45%) | Dashboard menampilkan 4500% atau salah status | `tipe_nilai` per indikator: `COUNT`/`RATIO`; agregasi berbeda | P0 |
| D5 | Denominator nol/hilang | `NO_DATA` atau invalid | Untuk rasio: tolak/`NO_DATA`; jangan anggap nol | P0 |
| D6 | Numerator > denominator | Bisa valid untuk kunjungan, mencurigakan untuk cakupan | Peringatan/outlier, bukan tolak otomatis kecuali definisi indikator | P1 |
| D7 | Negatif/outlier ekstrem | Merusak total | Tolak negatif; flag outlier dan minta konfirmasi | P1 |
| D8 | `nilai` eksplisit `0` | Kode saat ini menghitung ulang dari num/den | Perbaiki: `0` eksplisit harus dihormati; jangan ubah diam-diam | P1 |
| D9 | Pembulatan di file | Total tidak cocok | Simpan raw; pembulatan hanya tampilan; toleransi selisih terdokumentasi | P1 |

#### E. Klasifikasi IKU/indikator

| # | Edge case | Dampak | Penanganan yang diusulkan | Status |
| --- | --- | --- | --- | --- |
| E1 | Kolom hanya berisi kode IKU, tanpa indikator | Satu IKU bisa punya banyak indikator; ambigu | Wajib `kode_indikator`, atau aturan 1:1 yang disetujui; jika tidak, karantina | P0 |
| E2 | Kolom IKU induk/rollup | Double counting | Tolak/abaikan dengan peringatan; hanya leaf yang boleh input | P0 |
| E3 | Satu IKU punya beberapa indikator | Satu sel tidak bisa mengisi semuanya | Pisahkan kolom/baris per indikator; jangan pilih diam-diam | P0 |
| E4 | Satu indikator punya banyak label IKU | Nilai bisa dobel jika diulang per kolom | Nilai tetap satu; label many-to-many; duplikat sel ditolak | P1 |
| E5 | Kode tidak dikenal/nonaktif | Data masuk master liar | Karantina; jangan auto-create | P0 |
| E6 | Kode berubah/rename | File lama tidak cocok | Versi master + profil; alias hanya lewat master | P1 |
| E7 | `kode_kategori` tidak cocok | Salah kelompok | Tolak baris (sudah ada di classifier) | P0 |
| E8 | Sel yang sama muncul di dua kolom | Duplikat | `DUPLICATE`; karantina | P0 |

#### F. Agregasi dan kalkulasi

| # | Edge case | Dampak | Penanganan yang diusulkan | Status |
| --- | --- | --- | --- | --- |
| F1 | Baris desa + baris puskesmas + baris kabupaten di satu file | Double counting | Simpan hanya grain terendah; level atas diturunkan; total tidak diimpor | P0 |
| F2 | Dua submission final untuk periode/faskes/indikator sama | Nilai dijumlahkan dua kali | Satu final aktif per kunci, atau mekanisme supersede/revisi | P0 |
| F3 | Rata-rata persentase | Capaian menyesatkan | `RATIO_OF_SUMS`; jumlahkan num/den dulu | P0 |
| F4 | Indikator count memakai mesin rasio | 45 kasus menjadi 4500% | Agregasi `SUM` untuk count; pisahkan dari rasio | P0 |
| F5 | Bulan tidak lengkap dalam triwulan/tahun | Seolah-olah nol | `PARTIAL`/`NO_DATA`, bukan nol; jangan annualize | P0 |
| F6 | Snapshot `Total Ibu` dijumlahkan lintas bulan | Perempuan yang sama dihitung berkali-kali | Definisi snapshot vs kohort; jangan jumlahkan snapshot | Keputusan |
| F7 | Faskes pengupload dianggap faskes domisili | Salah atribusi | Pisahkan `faskes_id` pengupload dan dimensi wilayah/domisili | P0 |
| F8 | Induk dengan anak count dan ratio | Rasio tidak bermakna | Larang campuran atau formula eksplisit | Keputusan |
| F9 | Target per desa tidak ada | Status tidak bisa dihitung | `NO_TARGET`, bukan nol; target kabupaten tidak dibagi proporsional tanpa keputusan | P0 |
| F10 | `source_hash` desa mengambil semua file periode | Trace salah | Hash per scope dari file yang benar-benar menyumbang scope itu | P1 |
| F11 | Filter scope mengabaikan DESA/KECAMATAN | Semua nilai terjumlah ke tiap desa | Filter `scope_type`/`scope_id` secara eksplisit | P0 |
| F12 | `scope_code` tidak di-resolve ke master | Drill-down tidak cocok | Resolve ke `wilayah.id`; simpan kode asli + ID | P0 |

#### G. Dashboard dan pengalaman pengguna

| # | Edge case | Dampak | Penanganan yang diusulkan | Status |
| --- | --- | --- | --- | --- |
| G1 | Detail hanya per faskes | Tidak bisa menjawab desa | Tambah tabel/ drill-down desa dan kecamatan | P1 |
| G2 | Filter hanya periode bulanan | Tidak bisa triwulan/tahun | Tambah grain periode; turunkan dari bulanan | P1 |
| G3 | “Tidak terlayani” tanpa denominator | Bisa dianggap nol | Tampilkan `NO_DATA`/“belum ada sasaran”, bukan 0 | P0 |
| G4 | Ranking desa berdasarkan jumlah vs rasio | Menyesatkan | Tampilkan keduanya; definisikan urutan dan denominator | Keputusan |
| G5 | Denominator sangat kecil | Persentase ekstrem | Tampilkan jumlah mentah + flag; jangan sembunyikan | P1 |
| G6 | `0` dan `NO_DATA` tampil sama | Salah tafsir | Badge/status berbeda; legend jelas | P0 |
| G7 | Drill-down tidak sampai sel | Audit lemah | Simpan sheet/baris/kolom asal di `SUBMISSION_ROW` | P1 |
| G8 | Data `PARTIAL`/`STALE` | Angka terlihat final | Badge status dan tanggal; jangan sembunyikan | P1 |
| G9 | Ekspor menghilangkan status/satuan | Laporan salah tafsir | Sertakan satuan, status, dan catatan | P1 |
| G10 | Scope akses desa | Operator bisa melihat desa lain | Batasi sesuai faskes/role; Dinas tetap luas | P1 |

#### H. Tata kelola

| # | Edge case | Dampak | Penanganan yang diusulkan | Status |
| --- | --- | --- | --- | --- |
| H1 | Data individu (nama/NIK) | Risiko PII | Agregat per desa; jangan simpan identitas tanpa keputusan | P0 |
| H2 | Siapa yang boleh membuat profil format | Profil liar | Admin/Dinas membuat, audit; Puskesmas memilih profil aktif | P1 |
| H3 | Template berubah | File lama gagal | `template_version` + `profile_version` di-pin ke submission | P1 |
| H4 | Master wilayah/IKU berubah | Histori berubah | Snapshot versi; mapping berversi | P0 |
| H5 | Retensi file | Audit hilang | Ikuti keputusan retensi Dinas; backup | Keputusan |

---

### 10.6 Mengapa satu kolom per IKU tidak cukup

| Kebutuhan engine | Contoh | Akibat jika hanya satu kolom |
| --- | --- | --- |
| Numerator + denominator | Ibu hamil dilayani 115 dari sasaran 147 | Hanya 115; tidak bisa menghitung 78,2% atau “32 tidak terlayani” |
| Satuan/tipe nilai | 45 kasus kematian vs 45% cakupan | Dashboard bisa menampilkan 4500% atau status salah |
| Identitas indikator | Satu IKU punya dua indikator | Tidak jelas nilai masuk indikator mana |
| Leaf vs induk | Kolom `IKU-2026-001` (induk) dan `.01` (leaf) | Induk dan anak dijumlahkan dua kali |
| Periode | Satu file berisi beberapa bulan | Satu submission hanya punya satu periode |
| Wilayah | Nama desa vs kode desa | Salah desa, terutama nama ganda |
| Trace | Sel C5 dari sheet `Data` | Hanya bisa trace ke baris, bukan sel |
| Metode agregasi | Count memakai `SUM`, rasio memakai `RATIO_OF_SUMS` | Semua dipaksa menjadi rasio |

**Konsekuensi:** format `desa × IKU` hanya layak sebagai **tampilan** atau **format impor dengan profil**, bukan sebagai kontrak final. Untuk menjadi kontrak final, setiap kolom IKU harus dipecah menjadi minimal `numerator` dan `denominator` (atau `nilai` untuk count), ditambah periode dan scope.

---

### 10.7 Rekomendasi template final

#### A. Kontrak utama: long/tidy

Template resmi yang diusulkan tetap memakai bentuk long, dengan kolom berikut. Nama kolom mengikuti kontrak yang sudah ada agar kompatibel; kolom baru ditandai.

| Kolom | Wajib | Isi | Catatan |
| --- | --- | --- | --- |
| `periode` | Ya | `YYYY-MM` | Satu file = satu bulan; harus sama dengan submission |
| `kode_faskes` | Ya | Kode faskes pengupload | Satu file satu faskes |
| `kode_kategori` | Tidak | Kode kategori dashboard | Jika diisi, harus cocok dengan label |
| `kode_iku` | Ya | Kode IKU **leaf**; boleh beberapa dengan pemisah `;` | Induk/rollup ditolak; label pertama = konteks |
| `kode_indikator` | Ya | Kode indikator | Harus terhubung ke `kode_iku` |
| `scope_type` | Ya | `DESA`/`KECAMATAN`/`KABUPATEN` | Faskes bukan scope; faskes adalah pengunggah (`kode_faskes`) |
| `scope_code` | Ya | Kode wilayah resmi | Di-resolve ke master; bukan nama |
| `nama_wilayah` | Tidak | Nama untuk verifikasi manusia | Tidak dipakai sebagai kunci |
| `numerator` | Ya untuk rasio/count | Angka mentah | Jangan dibulatkan |
| `denominator` | Ya untuk rasio | Angka mentah | Untuk count: sesuai master (`1`/kosong) — keputusan |
| `nilai` | Tidak | Nilai jadi bila sumber memang persentase | Jika diisi, divalidasi; jangan menimpa num/den |
| `catatan` | Tidak | Catatan baris | Boleh kosong |

Kolom baru yang diusulkan (perlu persetujuan): `template_version` di sheet `Petunjuk`, serta `tipe_nilai`/`level_wilayah` di master indikator — **bukan** diisi bebas oleh pengunggah.

Template contoh yang sudah dibuat dari master pilot: [`assets/templates/bnba/arsip/template-bnba-final.xlsx`](../assets/templates/bnba/arsip/template-bnba-final.xlsx). Rincian sheet dan hasil verifikasinya ada di §10.14. Template v2 yang dipisah per pola input ada di §10.15.

#### B. Aturan template

1. **Satu file = satu faskes + satu bulan.** Jangan menggabungkan semua Puskesmas atau 12 bulan dalam satu file.
2. **Satu baris = satu indikator + satu scope.** Tidak ada matriks, merged cell, total, subtotal, atau formula.
3. **Kode resmi, nama hanya untuk manusia.** `scope_code`, `kode_iku`, `kode_indikator` harus ada di master.
4. **Blank ≠ 0.** Blank berarti belum ada data; `0` berarti nol yang dilaporkan.
5. **Hanya leaf yang boleh diinput.** IKU induk dihitung dari komponen.
6. **Rasio wajib numerator/denominator.** Count mengikuti `tipe_nilai` master; jangan menebak dari judul kolom.
7. **Tidak ada angka persentase tanpa num/den** untuk indikator rasio.
8. **Data validation + proteksi header.** Dropdown kode dari sheet `Referensi`; header dilindungi.
9. **Sheet `Petunjuk`** memuat versi template, periode, faskes, tanggal generate, dan kontak.
10. **Sheet `Referensi`** memuat daftar kode aktif, satuan, tipe nilai, dan level wilayah.
11. **Template digenerate server** dari master aktif untuk periode/faskes yang dipilih, bukan file statis.
12. **Versi template/profil di-pin** ke submission agar file lama bisa direproduksi.

#### C. Jika Puskesmas tetap memakai matriks (opsional, sementara)

Format matriks hanya diterima lewat **profil impor**:

- Satu file = satu faskes + satu periode; metadata di sheet `Petunjuk` atau blok atas.
- Header memuat kode, bukan hanya label: mis. `IKU-2026-001.01 | IND-001 | NUM` dan `... | DEN`.
- Tidak ada merged cell; header dua baris diperbolehkan hanya jika profil mendefinisikannya.
- Tidak ada baris/kolom total.
- Preview menampilkan hasil *unpivot* per baris + sel bermasalah; pengguna mengonfirmasi sebelum submit.
- Profil menyimpan: sheet, baris header, kolom scope/level, mapping kolom → `kode_iku`/`kode_indikator`/peran (`NUM`/`DEN`/`VALUE`), versi, dan faskes.
- Profil dibuat/disetujui Admin/Dinas dan di-audit; Puskesmas tidak membuat profil liar.
- **Jangan menyimpan bentuk matriks**; yang disimpan tetap baris panjang.

**Catatan:** profil per Puskesmas menambah beban perawatan. Karena itu, format matriks hanya untuk masa transisi; target akhirnya tetap template long.

#### D. Alternatif UI yang lebih baik

Jika tujuan utamanya adalah pengalaman mengisi seperti matriks, buat **grid input di aplikasi** (baris desa, kolom indikator, sel numerator/denominator). Data tetap disimpan long. Ini memberi tampilan yang familiar tanpa kerapuhan parsing Excel. Ini pekerjaan UI terpisah dan belum masuk scope pilot.

---

### 10.8 Aturan agregasi bulan, triwulan, tahun, dan wilayah

| Dimensi | Grain data | Rumus | Jika data kurang |
| --- | --- | --- | --- |
| Bulan | Satu periode | Rasio: `num/den×100`; count: `num` | `NO_DATA` |
| Triwulan | Tiga bulan | Jumlahkan `num` dan `den` tiga bulan, lalu rasio; count: jumlahkan | `PARTIAL` bila ada bulan kosong |
| Tahun | Dua belas bulan | Sama seperti triwulan | `PARTIAL` bila ada bulan kosong; jangan annualize |
| Desa | Baris desa | Jumlahkan baris desa; rasio dihitung setelah jumlah | `NO_DATA` per desa yang tidak ada |
| Puskesmas (domisili) | Mapping desa → puskesmas | Jumlahkan desa dalam wilayah kerja | `NO_DATA` bila mapping tidak ada |
| Puskesmas (pelayanan) | Faskes pelayanan | Jumlahkan baris dengan faskes pelayanan | `NO_DATA` bila kolom tidak ada |
| Kecamatan | Desa → kecamatan | Jumlahkan desa | `PARTIAL` bila desa kurang |
| Kabupaten | Desa/kecamatan/faskes | Jumlahkan **satu level** saja; jangan campur | `PARTIAL` bila komponen kurang |
| IKU leaf | Indikator | Rasio/count sesuai tipe | `NO_DATA`/`PARTIAL` |
| IKU induk | `IKU_KOMPONEN` | `RATIO_OF_SUMS` atau `SUM` sesuai metode | `PARTIAL` bila anak kurang |

Aturan wajib:

- **Jangan merata-ratakan persentase.** Jumlahkan numerator/denominator lebih dulu.
- **Jangan menjumlahkan snapshot lintas bulan.** Jika `Total Ibu` adalah snapshot, triwulan/tahun butuh metode khusus (mis. bulan terakhir, rata-rata, atau kohort) yang harus diputuskan Dinas.
- **Jangan mencampur input bulanan dan triwulanan.** Pilih satu sumber; turunkan sisanya.
- **Jangan membagi target kabupaten ke desa secara proporsional** tanpa keputusan. Tanpa target desa, status `NO_TARGET`.
- **Jangan menganggap desa yang tidak muncul sebagai nol.** Itu `NO_DATA` atau “belum lapor”.
- **Jangan memakai angka desa + angka puskesmas + angka kabupaten sekaligus.** Simpan grain terendah, turunkan ke atas.

---

### 10.9 Dampak ke implementasi saat ini

Temuan berikut adalah **fakta kode** yang harus diselesaikan sebelum BNBA per desa dapat dihitung dengan benar. Ini bukan daftar perubahan yang sudah dilakukan.

| # | Temuan | Bukti | Dampak | Status |
| --- | --- | --- | --- | --- |
| 1 | Kalkulasi hanya membuat scope `KABUPATEN` dan `FASKES` | `backend/internal/service/calculation.go:52-62` | Tidak ada nilai per desa/kecamatan | P0 |
| 2 | Query nilai untuk scope non-FASKES tidak memfilter `scope_type/scope_id` | `backend/internal/repository/sqlite/value_repo.go:18-30` | Jika DESA ditambahkan, setiap desa bisa menjumlahkan semua nilai | P0 |
| 3 | Input langsung selalu dihitung sebagai rasio; `aggregation_method` label tidak dipakai | `backend/internal/service/calculation.go:137-157,198-200` | Count 45 kasus bisa menjadi 4500%; `SUM` tidak berfungsi | P0 |
| 4 | Satu `submission_row` menghasilkan satu `indicator_value` | `backend/internal/worker/classifier.go:158-163`; `migrations/001_init.sql:239-254` | Matriks harus di-unpivot menjadi banyak baris; belum ada | P0 untuk wide |
| 5 | Nilai final dari beberapa submission dijumlahkan tanpa dedupe/supersede | `value_repo.go:18-35` | Upload ulang bulan yang sama bisa double count | P0 |
| 6 | Hanya `FASKES` yang divalidasi; `DESA` diterima mentah | `classifier.go:78-83`; `001_init.sql:196-197` | Desa di luar wilayah kerja bisa masuk; kode tidak di-resolve | P0 |
| 7 | `faskes.wilayah_id` adalah lokasi faskes, bukan mapping desa | `001_init.sql:26-33`; `011_desa_faskes_mapping.sql` | Sudah diganti mapping berversi `desa_faskes_mapping` | Selesai |
| 8 | Master wilayah hanya 2 kecamatan/2 desa contoh | `002_seed.sql:10-15` | Data nyata akan karantina massal | P0 |
| 9 | Detail dashboard hanya per faskes | `backend/internal/service/dashboard.go:585-630` | Tidak ada tabel/ drill-down desa | P1 |
| 10 | Filter dashboard hanya periode/faskes/status | `frontend/src/routes/+page.svelte:33-35,103-113`; `dashboard.go:39-42` | Belum ada triwulan/tahun/desa | P1 |
| 11 | Laporan bulanan memakai `numerator` IKU, bukan agregasi 12 bulan | `dashboard.go:180-203` | Triwulan/tahun belum dihitung | P1 |
| 12 | Parser multi-sheet sudah ada; header baris pertama non-kosong, wajib tiga kolom | `backend/internal/worker/parser.go`; §10.18 | File matriks/gambar tetap gagal | P1 untuk file matriks |
| 13 | Template hanya header CSV statis | `backend/internal/httpapi/handler/master.go:282-285` | Tidak ada validasi kode/dropdown/versi | P1 |
| 14 | Belum ada kolom level/tipe nilai di indikator | `001_init.sql:80-91` | Tidak bisa membedakan count/rasio/level | P1 |
| 15 | Belum ada trace sel asal (sheet/baris/kolom) | `001_init.sql:187-205` | Drill-down hanya sampai baris hasil, bukan sel Excel | P1 |
| 16 | `nilai` eksplisit `0` dihitung ulang; persen tidak dinormalisasi | `classifier.go:102-104`; `parser.go:169-178` | Nilai 0 bisa berubah; 0,45 vs 45 tidak konsisten | P1 |
| 17 | `source_hash` untuk scope non-FASKES mengambil semua file periode | `value_repo.go:159-170` | Trace desa tidak akurat | P1 |
| 18 | Submission hanya satu faskes | `001_init.sql:145-161` | File gabungan semua Puskesmas tidak bisa diproses adil | P1 |

**Kesimpulan teknis:** menambahkan format matriks tanpa memperbaiki #1–#8 akan menghasilkan dashboard yang terlihat benar tetapi salah hitung. Urutan yang aman: perbaiki fondasi kalkulasi dan master, lalu template long, baru adapter wide.

---

### 10.10 Keputusan baru yang harus diambil

| # | Pertanyaan | Opsi | Dampak | Rekomendasi |
| --- | --- | --- | --- | --- |
| 12 | Format final: long, wide, atau keduanya? | A. Long saja; B. Wide saja; C. Long resmi + wide lewat profil | A: paling aman; B: rapuh; C: fleksibel tetapi ada beban profil | **C** untuk transisi, target akhir A |
| 13 | Cara membedakan count dan rasio | A. Kolom baru `tipe_nilai`; B. Pakai `satuan`; C. Tebak dari nilai | A: eksplisit; B: rawan teks bebas; C: dilarang | **A** |
| 14 | Sifat data bulanan | A. Incremental; B. Snapshot; C. Campuran per indikator | Menentukan boleh/tidaknya penjumlahan triwulan/tahun | **C** dengan field eksplisit |
| 15 | Arti “per puskesmas” | A. Domisili; B. Pelayanan; C. Keduanya | Menentukan kebutuhan mapping dan kolom faskes pelayanan | **C**; default A untuk peta wilayah |
| 16 | Input triwulan/tahunan | A. Hanya bulanan, turunkan; B. Boleh langsung | A: konsisten; B: risiko double count | **A** |
| 17 | Sumber denominator/sasaran | A. Kolom BNBA; B. Sheet `Sasaran`; C. Master sasaran terpisah | A: cepat tetapi tidak semua indikator punya; B: jelas per file; C: paling terkontrol | **C/B** sesuai keputusan Dinas |
| 18 | Mapping desa → puskesmas | A. Tabel berversi; B. Tebak dari faskes uploader; C. Tidak ada | A: benar; B: salah atribusi; C: tidak bisa per puskesmas | **A** |
| 19 | Prioritas adapter wide | A. Setelah long stabil; B. Bersamaan; C. Tidak perlu | A: risiko rendah; B: kompleks; C: Puskesmas harus ubah format | **A** |

**Jangan mengisi keputusan 12–19 dengan asumsi.** Jika Dinas/pemilik proyek belum memutuskan, pekerjaan yang bergantung padanya berhenti.

**Keputusan pilot 2026-11-20:** #13 memakai `tipe_nilai` eksplisit; #16 mengizinkan periode tahunan `YYYY` (triwulan belum); #18 memakai tabel mapping desa → puskesmas domisili; format final memakai satu template `Rasio`/`Nilai`. Konfirmasi Dinas untuk produksi tetap diperlukan.

---

### 10.11 Rencana implementasi jika disetujui

| Fase | Pekerjaan | Gate |
| --- | --- | --- |
| A0 | Putuskan §10.10; lengkapi master wilayah, tipe nilai, dan level wilayah. Mapping desa → puskesmas sudah tersedia (§10.19) | Dinas menyetujui definisi dan master |
| A1 | Perbaiki fondasi: filter scope, agregasi `SUM`/`RATIO_OF_SUMS`, dedupe/supersede, source hash per scope | Tes kalkulasi count/rasio dan scope lulus |
| A2 | Template long + validasi + `template_version` + generator dari master | File contoh lolos preview tanpa karantina massal |
| A3 | Kalkulasi `DESA`/`KECAMATAN` + dashboard drill-down + periode triwulan/tahun | Angka desa → puskesmas → kabupaten konsisten |
| A4 | Adapter wide opsional + profil berversi + preview unpivot | File matriks contoh bisa dipetakan tanpa double count |
| A5 | Hardening: PII, retensi, performa, ekspor, dokumentasi, diagram | Uji beban, audit, dan acceptance criteria lulus |

Urutan ini sengaja menempatkan **perbaikan kalkulasi sebelum format matriks**. Membangun parser matriks lebih dulu hanya akan mempercepat masuknya angka yang belum bisa dihitung dengan benar.

---

### 10.12 Tes dan fixture yang wajib

| # | Skenario | Ekspektasi |
| --- | --- | --- |
| 1 | Long: 2 desa × 2 indikator (1 rasio, 1 count) × 1 bulan | Leaf, induk, desa, puskesmas, kabupaten konsisten; count tidak menjadi persen |
| 2 | Matriks dengan header merged dan baris total | Unpivot benar; total tidak masuk; sel ambigu karantina |
| 3 | Blank vs `0` | Blank `NO_DATA`; `0` tetap nol |
| 4 | Dua submission final periode/faskes/indikator sama | Tidak double count; supersede/revisi sesuai aturan |
| 5 | Triwulan/tahun dengan satu bulan kosong | `PARTIAL`, bukan nol |
| 6 | Nama desa typo/nama ganda | Karantina; tidak menebak |
| 7 | Count 45 kasus | Dashboard menampilkan 45 kasus, bukan 4500% |
| 8 | Kolom IKU induk | Ditolak/diabaikan; tidak double count |
| 9 | Pasien domisili A, faskes pelayanan B | Atribusi domisili dan pelayanan terpisah |
| 10 | Drill-down | Sampai sheet/baris/kolom asal |
| 11 | `45%` vs `0,45` | Ternormalisasi sesuai satuan; tidak 100× |
| 12 | Desa di luar wilayah kerja puskesmas | `SCOPE_DENIED`/karantina |

Fixture harus **anonim** dan disimpan di `backend/testdata/`; jangan memakai data pasien nyata. File contoh pada gambar tidak boleh dianggap fixture produksi.

---

### 10.13 Ringkasan

| Pertanyaan | Jawaban |
| --- | --- |
| Apakah format desa × IKU bisa dipetakan? | Bisa untuk tampilan/agregasi jika kode, scope, periode, numerator/denominator, dan aturan total jelas. Tidak bisa menjadi kontrak final apa adanya. |
| Mengapa tidak bisa apa adanya? | Satu sel tidak membawa num/den, satuan, periode, identitas indikator, dan trace sel; count vs rasio bercampur; IKU induk bisa double count. |
| Template final yang disarankan | Long/tidy: satu baris = satu periode + scope + IKU/indikator + numerator/denominator. |
| Template v1 review | [`assets/templates/bnba/arsip/template-bnba-final.xlsx`](../assets/templates/bnba/arsip/template-bnba-final.xlsx) — `Data` (long), `Matriks` (referensi desa × IKU), `Referensi`, `Contoh`, `Sasaran`; lihat §10.14. |
| Contoh v1 terisi | [`assets/templates/bnba/arsip/contoh-bnba-terisi.xlsx`](../assets/templates/bnba/arsip/contoh-bnba-terisi.xlsx) — 2 desa × 12 indikator, 24 baris, angka fiktif. |
| Template v2 terpisah | `assets/templates/bnba/arsip/v2/` — `BNBA-1-Kejadian-Bulanan`, `BNBA-2-Cakupan-Bulanan`, `BNBA-3-Tahunan`; eksplorasi UX, lihat §10.15. |
| Template dinamis (legacy) | [`assets/templates/bnba/template-bnba-dinamis.xlsx`](../assets/templates/bnba/template-bnba-dinamis.xlsx) — IKU/indikator dinamis + scope wilayah desa/kecamatan/kabupaten; faskes adalah pengunggah; lihat §10.16. Template utama sekarang [`BNBA-Semua-Satuan.xlsx`](../assets/templates/bnba/final/BNBA-Semua-Satuan.xlsx) (§10.18). |
| Bagaimana dengan format Puskesmas yang sudah ada? | Adapter wide lewat profil berversi; unpivot di preview; jangan simpan wide. |
| Per desa/puskesmas/bulan/triwulan/tahun? | Simpan grain terendah; turunkan ke atas. Triwulan/tahun = jumlah num/den, bukan rata-rata persentase. |
| Blocker terbesar | Denominator/definisi, master wilayah + mapping, agregasi count vs rasio, filter scope, dedupe/supersede. |
| Status | **Analisis selesai; implementasi belum dimulai.** Track A tetap `BLOCKED` sampai keputusan §10.10 dan file asli diterima. |

---

### 10.14 Template final yang dihasilkan (v1 review, 2026-11-19)

> **Catatan:** bagian ini adalah **v1 review** dengan satu `Data` berisi `N/D` untuk semua indikator. Implementasi template terpisah yang lebih baru ada di §10.15.

**Fakta:** file template sudah dibuat di [`assets/templates/bnba/arsip/template-bnba-final.xlsx`](../assets/templates/bnba/arsip/template-bnba-final.xlsx) dari master pilot `backend/data/matasidik.db`. Tidak ada kode aplikasi yang diubah.

**File contoh terisi:** [`assets/templates/bnba/arsip/contoh-bnba-terisi.xlsx`](../assets/templates/bnba/arsip/contoh-bnba-terisi.xlsx) — 2 desa × 12 indikator = **24 baris** di sheet `Data`, dengan `Matriks` terisi untuk 2 desa. Semua angka fiktif. Ini **lengkap untuk master pilot saat ini** (2 desa), bukan seluruh desa Kabupaten Sumbawa Barat.

| Sheet | Isi | Status |
| --- | --- | --- |
| `Data` | Kontrak unggah long/tidy; header sesuai parser pilot + `nama_wilayah`; hanya header, sengaja tanpa baris data | **Sheet unggah** |
| `Petunjuk` | Identitas template, cara isi, aturan wajib, kesiapan kode, keputusan terbuka | Bacaan |
| `Matriks` | Tampilan desa × IKU/indikator dengan pasangan `N`/`D`; 2 desa contoh + 20 baris kosong | **Referensi, belum untuk unggah** |
| `Referensi` | Kode kategori, IKU, IKU leaf, indikator, wilayah, faskes, periode, kode scope, tipe nilai | Validasi |
| `Contoh` | 4 baris contoh long/tidy dengan angka fiktif | Ilustrasi |
| `Sasaran` | Usulan sheet denominator per desa; belum dibaca aplikasi | Diskusi |

Isi penting:

- `Data` adalah sheet pertama agar parser pilot yang hanya membaca sheet pertama tetap bisa memakainya. Tidak ada catatan di bawah header `Data` karena parser akan membacanya sebagai baris data.
- `Matriks` memakai pasangan `N`/`D` untuk setiap IKU/indikator, bukan satu kolom per IKU. Ini menjaga numerator/denominator tetap eksplisit.
- `Matriks` menaruh **nama wilayah di kolom pertama** dan **kode wilayah di kolom kedua**; kode tetap kunci input, nama otomatis dari `Referensi`.
- Header `Matriks` memuat `kode_iku · kode_indikator — nama`. Pada master pilot saat ini setiap leaf punya satu indikator, sehingga satu grup per leaf; jika nanti satu IKU punya beberapa indikator, tambahkan grup kolom per indikator.
- `Referensi` hanya memuat kode yang **benar-benar ada** di master pilot. Kode contoh seperti `IKU-ibu-meninggal` tidak dibuat karena belum ada di master; kode baru harus ditambahkan lewat Master/Pengajuan IKU.
- `Contoh` memakai angka fiktif (`10/20`, `5/10`, `8/10`, `20/40`) dan diberi label “ILUSTRASI — angka fiktif”; jangan dipakai sebagai data produksi.
- `Sasaran` hanya untuk mendiskusikan denominator; belum ada parser/kalkulasi yang membacanya.
- Template v1 memakai `scope_type=DESA` dan kode desa untuk contoh utama. Pada kontrak final, **FASKES bukan scope**: faskes adalah pengunggah/pemilik data di `kode_faskes`; `scope_type` hanya wilayah.

**Verifikasi yang dijalankan:**

- 4 baris dari sheet `Contoh` disalin sementara ke sheet `Data` pada salinan file, lalu diuji dengan parser + classifier pilot.
- Hasil: 4 baris terbaca, 4/4 `MATCHED`, 4 `indicator_value` terbentuk. Uji ini memakai database seed pilot dan file sementara; bukan data produksi.
- File contoh terisi `contoh-bnba-terisi.xlsx` diuji dengan parser + classifier pilot: **24 baris terbaca, 24/24 `MATCHED`**, 24 `indicator_value` terbentuk. Angka tetap fiktif dan tidak disimpan sebagai data produksi.
- `unzip -t` pada file `.xlsx` tidak menemukan error arsip.
- Template asli tetap memiliki `Data` kosong; mengunggah template kosong akan berhenti dengan pesan “tidak memiliki baris data”, yang memang diharapkan sampai diisi.

**Yang belum bisa dilakukan template ini tanpa perubahan kode:**

- Menghitung `iku_value` per desa/kecamatan (kalkulasi pilot baru `KABUPATEN`/`FASKES`).
- Menghitung indikator count dengan `SUM` (mesin pilot selalu rasio).
- Mengunggah sheet `Matriks` langsung (butuh adapter wide/unpivot).
- Memvalidasi `scope_code` desa ke master dan membatasi desa di luar wilayah kerja.
- Menyimpan `template_version`/profil ke submission.
- Membaca sheet `Sasaran`.

**Cara melihat:** buka [`assets/templates/bnba/arsip/template-bnba-final.xlsx`](../assets/templates/bnba/arsip/template-bnba-final.xlsx); mulai dari `Petunjuk`, lalu `Data`, `Matriks`, dan `Referensi`. Untuk melihat versi yang sudah terisi, buka [`assets/templates/bnba/arsip/contoh-bnba-terisi.xlsx`](../assets/templates/bnba/arsip/contoh-bnba-terisi.xlsx). Untuk uji upload, isi `Data` minimal satu baris yang valid, pilih periode/faskes yang sama di aplikasi, lalu unggah.

---

### 10.15 Template v2 — tiga jalur input (2026-11-19)

**Fakta:** template terpisah sudah dibuat di `assets/templates/bnba/arsip/v2/`. Ini adalah **eksplorasi UX** untuk pola input yang lebih sederhana; **belum** bisa diunggah dengan parser pilot saat ini. Kontrak generik yang sudah diterima aplikasi ada di §10.16. File v1 tetap disimpan sebagai riwayat review.

| File | Pola | Kolom `Data` | Satuan hasil | Status master |
| --- | --- | --- | --- | --- |
| [`BNBA-1-Kejadian-Bulanan.xlsx`](../assets/templates/bnba/arsip/v2/BNBA-1-Kejadian-Bulanan.xlsx) | Kejadian/count bulanan | `periode`, `kode_faskes`, `kode_desa`, `kode_indikator`, `jumlah`, `catatan` | `Kasus` | Belum ada indikator kejadian di master pilot |
| [`BNBA-2-Cakupan-Bulanan.xlsx`](../assets/templates/bnba/arsip/v2/BNBA-2-Cakupan-Bulanan.xlsx) | Cakupan/rasio bulanan | `periode`, `kode_faskes`, `kode_desa`, `kode_indikator`, `numerator`, `denominator`, `catatan` | `%` | 12 indikator pilot tersedia |
| [`BNBA-3-Tahunan.xlsx`](../assets/templates/bnba/arsip/v2/BNBA-3-Tahunan.xlsx) | Tahunan | `tahun`, `scope_type`, `kode_scope`, `kode_indikator`, `nilai`, `catatan` | Sesuai indikator | Belum ada indikator tahunan di master pilot |

Prinsip implementasi:

- **Disable by omission:** kolom yang tidak relevan tidak dibuat. Template kejadian tidak punya `numerator`/`denominator`; template cakupan tidak punya `jumlah`; template tahunan v2 (arsip) tidak punya `bulan`/`kode_desa`.
- **Data minimal:** `kode_iku`, nama indikator, satuan, `scope_type`/`scope_code` diturunkan sistem dari master; tidak ditampilkan di `Data` agar sederhana.
- **Referensi lengkap:** setiap file punya sheet `Referensi` berisi kode, nama, `kode_iku`, kategori, `satuan_master`, `satuan_hasil`, `tipe_input`, `frekuensi`, `arah`, `agregasi`, dan `scope`.
- **Validasi:** dropdown periode/faskes/desa/indikator, validasi angka, `denominator > 0` untuk cakupan, conditional formatting untuk nilai kosong/negatif.
- **Proteksi:** sheet dilindungi; hanya sel input yang tidak terkunci. Ini mengurangi salah edit, bukan kontrol keamanan.
- **Tanpa formula di baris data:** baris kosong tetap benar-benar kosong; tidak ada 500 baris formula yang bisa dibaca parser sebagai data.
- **Tanpa macro:** `.xlsx`, bukan `.xlsm`; tidak ada VBA.
- **Generator:** `assets/templates/bnba/arsip/v2/gen_templates.py` (memerlukan `openpyxl`; membaca `backend/data/matasidik.db`).

Catatan penting:

- [`BNBA-2-Cakupan-Bulanan.xlsx`](../assets/templates/bnba/arsip/v2/BNBA-2-Cakupan-Bulanan.xlsx) memakai 12 indikator pilot. `IND-001`–`IND-003` memiliki `satuan_master = orang` (numerator), tetapi hasil IKU/SPM-nya `%`; template menampilkan `satuan_hasil = %`. Normalisasi satuan master masih perlu keputusan Dinas.
- `BNBA-1` dan `BNBA-3` adalah **kerangka siap pakai**, tetapi belum punya indikator valid karena master pilot belum memuat indikator kejadian/tahunan. Kode `CONTOH-*` di sheet `Contoh` hanya ilustrasi, bukan master.
- Ketiga template **belum bisa diunggah apa adanya** dengan parser pilot: parser masih mewajibkan `kode_iku` + `scope_type`/`scope_code`, dan kalkulasi masih ratio-only. Penyesuaian ada di §10.9–§10.11.
- `Referensi`/`Contoh` tidak dibaca parser; hanya sheet `Data` yang menjadi kontrak unggah.

**Verifikasi yang sudah dijalankan:** ketiga file `.xlsx` valid (`unzip -t`), sheet `Data` berada di urutan pertama, proteksi aktif, dropdown/validasi terpasang, dan baris data kosong tidak berisi formula. **Uji upload end-to-end belum dijalankan** karena parser pilot masih mewajibkan `kode_iku`/`scope_type`/`scope_code`; itu bagian dari penyesuaian kode A1/A2.

---

### 10.16 Template dinamis yang diterima aplikasi (2026-11-19)

**Fakta:** template generik sudah dibuat di [`assets/templates/bnba/template-bnba-dinamis.xlsx`](../assets/templates/bnba/template-bnba-dinamis.xlsx). Ini adalah **kontrak upload yang diterima parser/classifier pilot** untuk baris rasio, dan tidak mengunci daftar indikator.

| Kolom `Data` | Isi |
| --- | --- |
| `periode` | `YYYY-MM`; harus sama dengan periode upload |
| `kode_faskes` | Kode faskes pengunggah |
| `kode_kategori` | Opsional; divalidasi terhadap kategori label |
| `kode_iku` | Satu atau beberapa kode IKU leaf, dipisahkan `;`; label pertama = utama/konteks; **dinamis dari master** |
| `kode_indikator` | Indikator; **dinamis dari master**; harus terhubung ke `kode_iku` |
| `scope_type` | `DESA`, `KECAMATAN`, `KABUPATEN` |
| `scope_code` | Kode wilayah sesuai level |
| `nilai` | Opsional; untuk rasio biarkan kosong |
| `numerator` | Angka mentah |
| `denominator` | Wajib > 0 untuk rasio |
| `catatan` | Catatan baris |

Prinsip:

- **Indikator dinamis.** Sheet `Data` tidak berisi daftar indikator hardcode. `kode_iku`/`kode_indikator` berasal dari master; jika master berubah, regenerate template. `Referensi` memuat mapping **semua label** IKU ↔ indikator.
- **Scope level eksplisit.** `scope_type` menentukan level filter: desa, kecamatan, kabupaten. `scope_code` adalah kode wilayah resmi level tersebut. **Faskes bukan scope**; faskes adalah pengunggah/pemilik data di `kode_faskes`.
- **Many-to-many label.** Satu indikator boleh punya banyak label. `kode_iku` di `Data` boleh berisi beberapa kode dengan pemisah `;` (contoh `IKU-A;IKK-B`); label pertama = utama/konteks, semua label harus leaf dan terhubung ke indikator. Jangan menduplikasi baris per label.
- **PJ terpisah.** `Bidang/Timker` adalah dimensi PJ (`pj_bidang_id`/`pj_timker_id`), bukan label IKU. Template menampilkannya di tabel PJ, bukan di `kode_iku`.
- **Trace label.** `SUBMISSION_ROW.matched_label_id` menyimpan label pertama (utama/konteks); label lain tetap terhubung lewat `indicator_label` dan muncul di `IKU_VALUE` masing-masing.
- **Validasi kode_iku.** Kolom `kode_iku` tidak memakai dropdown Excel karena boleh berisi beberapa kode dengan `;`; validasi setiap kode dilakukan classifier/server.
- **Satu grain terendah.** Simpan satu baris pada level terendah yang tersedia; level atas diturunkan. Jangan menulis desa + kecamatan + kabupaten sekaligus untuk indikator/periode yang sama.
- **Filter puskesmas.** Puskesmas domisili dihitung dari mapping `desa_faskes_mapping`; fallback ke `kode_faskes` pengunggah hanya bila desa belum punya nilai. Jangan menebak dari `faskes.wilayah_id`.
- **Data minimal.** Nama wilayah/indikator dan satuan diturunkan sistem dari master; tidak perlu diisi ulang di `Data`.

**Verifikasi:** 3 baris `Contoh` (DESA, KECAMATAN, KABUPATEN) disalin sementara ke `Data`, lalu diuji dengan parser + classifier pilot: **3/3 `MATCHED`**, 3 `indicator_value` terbentuk. Artinya template ini diterima aplikasi saat ini untuk baris rasio. Uji tambahan `IKU-2026-001.01; IKK-TEST` pada indikator yang punya dua label juga `MATCHED`; label pertama menjadi `matched_label_id`. Tabel `Referensi` sekarang memuat semua pasangan `indikator_label` (termasuk `is_primary` dan `is_input`); pada master pilot saat ini masih 1:1, sehingga baru tampil satu baris per indikator.

**Yang sudah aktif di aplikasi (2026-11-20):**

- `scope_code` DESA/KECAMATAN/KABUPATEN di-resolve ke master `wilayah`; `kode_faskes` divalidasi terhadap faskes pengunggah.
- Kalkulasi membuat `iku_value` per DESA/KECAMATAN/KABUPATEN dan per FASKES pengunggah.
- Agregasi memakai grain terendah per indikator untuk mencegah double counting; KABUPATEN memakai fallback FASKES hanya untuk data legacy.
- Filter dashboard desa/kecamatan/kabupaten dan filter faskes pengunggah sudah terpisah.
- Detail SPM menampilkan tabel wilayah (desa/kecamatan/kabupaten) selain tabel faskes.
- `kode_iku` multi-label (`;`) diterima; label pertama menjadi konteks.

**Yang belum aktif/diputuskan:**

- Indikator `COUNT`/`SCALAR` (kasus/tahun/poin) belum didukung; mesin masih rasio.
- Kebijakan data yang bercampur level (DESA + KECAMATAN untuk indikator/periode sama) belum final; implementasi memilih level terendah.
- Dedupe/supersede lintas submission belum ada; mapping desa sudah tersedia (§10.19).
- Target per desa/kecamatan belum tersedia; dashboard memakai fallback target kabupaten.

**Hubungan dengan v2:** §10.15 adalah eksplorasi UX (template terpisah per pola). Setelah §10.18, parser membaca sheet `Rasio`/`Nilai`/`Data`; template utama adalah `BNBA-Semua-Satuan.xlsx`, sedangkan `template-bnba-dinamis.xlsx` tetap legacy untuk rasio.

Generator: [`assets/templates/bnba/gen_dinamis.py`](../assets/templates/bnba/gen_dinamis.py) (memerlukan `openpyxl`; membaca `backend/data/matasidik.db`).

---

### 10.17 Implementasi upload → dashboard (2026-11-20)

**Status:** alur teknis untuk rasio, count/scalar, dan periode tahunan sudah diimplementasikan; keputusan bisnis denominator/definisi, target per desa, dan dedupe tetap terbuka.

| Lapisan | Perubahan |
| --- | --- |
| Upload/classifier | `scope_type` hanya DESA/KECAMATAN/KABUPATEN untuk template baru; `scope_code` di-resolve ke `wilayah.id`; `kode_faskes` divalidasi; `kode_iku` multi `;`; `tipe_nilai` menentukan validasi `RATIO` vs `COUNT`/`SCALAR`; `frekuensi` harus cocok dengan periode; FASKES legacy tetap diterima |
| Kalkulasi | `iku_value` per DESA/KECAMATAN/KABUPATEN/FASKES pengunggah; agregasi memilih level terendah; rasio `RATIO_OF_SUMS`, count `SUM`, scalar `LAST_VALUE`; `source_hash` per scope; KABUPATEN fallback FASKES untuk seed legacy |
| Dashboard | Filter `scope_type`/`scope_id` terpisah dari `faskes_id`; detail SPM menampilkan per wilayah; frontend punya selector wilayah |
| Upload UI | Template CSV header final + unduhan XLSX; preview menampilkan scope dan input kode IKU multi `;` |
| Test | `TestCalculateScopeHierarchy`, `TestClassifyMultipleIKUCodes`, `TestParseXLSXReadsRasioAndNilaiSheets`, `TestClassifyCountAndScalarWithoutDenominator`, `TestCalculateCountScalarAnnual`, `TestUploadGeneralTemplateRatioAndCount`; `go test ./...`, frontend check/build |
| E2E API | Upload XLSX `Rasio` + `Nilai` bulanan (rasio + count) dan `Nilai` tahunan (scalar), preview `MATCHED`, approval Puskesmas → Dinas → Kabid final; `iku_value` count/scalar menyimpan nilai tanpa denominator |

**Belum dikerjakan:** format harian, target per desa, dedupe/supersede lintas submission, diagram/dokumentasi skema penuh. Mapping desa sudah tersedia (§10.19).

---

### 10.18 Satu template untuk semua satuan dan periode (2026-11-20)

**Keputusan pemilik proyek:** tidak ada template tahunan terpisah. Data bulanan dan tahunan memakai satu file; kolom `periode` menentukan granularitas. Format harian tersedia di template, tetapi perhitungan harian belum aktif. File tahunan lama diarsipkan di `assets/templates/bnba/arsip/BNBA-Tahunan.xlsx` dan tidak dipakai.

**Fakta:** template final ada di [`BNBA-Semua-Satuan.xlsx`](../assets/templates/bnba/final/BNBA-Semua-Satuan.xlsx):

| Sheet | Isi |
| --- | --- |
| `Petunjuk` | Aturan pola nilai, periode, satuan, dan cara isi |
| `Rasio` | `numerator` + `denominator` untuk `%`, `Persentase`, `(%)`, `Rasio` |
| `Nilai` | `nilai` untuk `Kasus`, `Dokumen`, `Laporan`, `Berita Acara`, `orang` (count), `(tahun)`, `Poin`, dan indeks/scalar |
| `Referensi` | Satuan, tipe nilai, pola, format periode, kategori, faskes, scope, indikator, mapping label, PJ |
| `Contoh` | Contoh rasio, count, dan scalar |

Pemetaan satuan dari sheet 2:

| Satuan | Tipe | Field | Agregasi |
| --- | --- | --- | --- |
| `(tahun)` | `SCALAR` | `nilai` | `LAST_VALUE` |
| `Kasus` | `COUNT` | `nilai` | `SUM` |
| `%`, `Persentase`, `(%)` | `RATIO` | `numerator + denominator` | `RATIO_OF_SUMS` |
| `Poin` | `SCALAR` | `nilai` | `LAST_VALUE` |
| `Rasio` | `RATIO` (non-persen) | `numerator + denominator` | `RATIO_OF_SUMS` |
| `Dokumen`, `Laporan`, `Berita Acara` | `COUNT` | `nilai` | `SUM` |
| `orang` | `COUNT/RATIO` | `nilai` atau numerator/denominator | `SUM`/`RATIO_OF_SUMS` |

**Keputusan teknis:** `indikator.tipe_nilai` eksplisit (`RATIO`/`COUNT`/`SCALAR`) ditambahkan lewat migrasi [`010_value_types_and_annual_period.sql`](../backend/internal/repository/sqlite/migrations/010_value_types_and_annual_period.sql). `periode.bulan = 0` menandai tahunan; periode tahunan 2026 memakai kode `2026`. Migrasi ini untuk pilot SQLite; produksi PostgreSQL memerlukan migrasi setara. Normalisasi teks satuan tetap perlu konfirmasi Dinas untuk produksi.

Aturan:

- Parser membaca sheet data `Rasio`, `Nilai`, atau `Data`; `Petunjuk`/`Referensi`/`Contoh` tidak dibaca sebagai data.
- `RATIO` wajib `numerator` dan `denominator` (> 0). `COUNT`/`SCALAR` wajib `nilai`; denominator tidak diwajibkan.
- `periode` bulanan `YYYY-MM`; tahunan `YYYY`; format `YYYY-MM-DD` tersedia di template tetapi belum dihitung.
- `frekuensi` indikator harus cocok dengan periode: indikator tahunan ditolak pada submission bulanan, dan sebaliknya.
- Count diagregasi `SUM`; scalar memakai `LAST_VALUE`; rasio persen memakai `RATIO_OF_SUMS` dan `×100`; satuan `Rasio` non-persen memakai `num/den` tanpa `×100`.
- `kode_iku` boleh multi `;`; `Bidang/Timker` tetap PJ, bukan label.
- `scope_type` hanya DESA/KECAMATAN/KABUPATEN; faskes adalah pengunggah.
- `kode_kategori` opsional; `Referensi` memuat master, tipe nilai, periode, dan scope.
- Generator: [`assets/templates/bnba/final/gen_template.py`](../assets/templates/bnba/final/gen_template.py); hasil juga disalin ke `frontend/static/templates/`.

**Verifikasi:** tes parser multi-sheet, klasifikasi count/scalar tanpa denominator, penolakan indikator tahunan pada periode bulanan, kalkulasi `SUM`/`LAST_VALUE`, dashboard count/scalar, dan upload template gabungan lulus. Detail perintah ada di laporan pekerjaan; `go test ./...`, `npm run check`, dan `npm run build` lulus.

**Legacy:** `template-bnba-dinamis.xlsx` tetap bisa diunggah untuk rasio. Template baru adalah format utama; file tahunan terpisah tidak dipakai.

### 10.19 Mapping desa → puskesmas domisili (2026-11-20)

**Fakta:** mapping disimpan di tabel [`desa_faskes_mapping`](../backend/internal/repository/sqlite/migrations/011_desa_faskes_mapping.sql) dan merupakan keputusan #18 (tabel berversi, bukan tebakan dari faskes uploader).

| Kolom | Fungsi |
| --- | --- |
| `desa_id` | FK ke `wilayah` level desa |
| `faskes_id` | FK ke puskesmas domisili |
| `berlaku_dari` / `berlaku_sampai` | Versi berlaku |
| `aktif` | Hanya satu mapping aktif per desa |
| `created_by` | Admin pembuat |

Seed pilot: `Desa Dalam` dan `Desa Beru` → `PKM Taliwang`. `faskes.wilayah_id` tetap lokasi faskes, bukan mapping desa.

Aturan agregasi:

- Scope `FASKES` mengambil nilai `DESA` yang desanya dipetakan ke faskes tersebut, lalu fallback ke nilai `FASKES` uploader hanya jika indikator itu belum punya nilai desa.
- `KECAMATAN`/`KABUPATEN` tetap memakai nilai desa; mapping tidak menggandakan nilai.
- `source_hash` scope `FASKES` mencakup file dari desa yang dipetakan, bukan hanya file yang diunggah faskes domisili.
- Mapping tidak mengubah otorisasi upload; ia menentukan atribusi/agregasi domisili. Validasi scope uploader tetap mengikuti aturan yang ada.

API dan UI:

- `GET /api/v1/master/mapping-desa`, `POST /api/v1/master/mapping-desa`, `DELETE /api/v1/master/mapping-desa/{id}`.
- Master → tab **Mapping Desa** menampilkan desa, puskesmas domisili, tanggal berlaku, status; admin dapat menambah versi baru atau menonaktifkan.
- Menyimpan mapping baru menutup mapping aktif lama dan menyimpan baris baru; histori tetap tersimpan. Tanggal berlaku tidak boleh melebihi hari ini.

**Verifikasi:** tes `TestDesaFaskesMappingVersioning` (pindah domisili, daftar desa per faskes, nonaktifkan) dan `TestCalculateScopeHierarchy` (nilai desa yang diunggah faskes lain tetap masuk faskes domisili, tidak masuk faskes lain) lulus.

**Belum dikerjakan:** antarmuka edit tanggal berlaku historis, validasi tumpang tindih tanggal, dan pemetaan desa lintas kecamatan yang belum ada di master wilayah.

---

## 11. Status berikutnya

**Track A:** alur teknis upload → kalkulasi per wilayah → dashboard untuk rasio, count/scalar, dan periode tahunan sudah diimplementasikan (§10.17–§10.19). Keputusan bisnis denominator/definisi, target per desa, dan dedupe lintas submission masih terbuka; mapping desa sudah diimplementasikan (§10.19). Template v1 review ada di §10.14, template v2 eksplorasi di `assets/templates/bnba/arsip/v2/` (§10.15), dan template dinamis yang diterima parser pilot ada di [`assets/templates/bnba/template-bnba-dinamis.xlsx`](../assets/templates/bnba/template-bnba-dinamis.xlsx) (§10.16).

**Track B sudah berjalan.** Progress implementasi:

- [x] Migrasi `dashboard_category` + seed 8 kategori.
- [x] `label.category_id` dan backfill IKU/SPM.
- [x] `bidang.jenis`/`parent_id` dan `indikator.pj_timker_id`.
- [x] CRUD kategori dan bidang/timker (API + master UI).
- [x] Form IKU/label memilih kategori.
- [x] Form indikator multi-label + label utama + PJ bidang/timker.
- [x] Tombol header **Pratinjau Pemetaan** (indikator + label/kategori + PJ), terpisah dari tab CRUD.
- [x] Filter dashboard per kategori + opsi **Semua** + query owner.
- [x] Laporan absolut bulanan dari IKU/SPM + target layanan + ekspor XLSX/PDF.
- [x] Field target layanan tahunan di Master Target.
- [x] Role `pkm_approver`, pengajuan IKU, dan approval BNBA sebelum Dinas.
- [x] Penghapusan role/kolom `iku_owner`.
- [x] Tampilan per peran: sidebar, route guard, tombol aksi, dan tab panduan.
- [x] Pembersihan label UI dari istilah tech stack/internal.
- [x] Pembatasan `GET /factors` dan `GET /tindak-lanjut` ke Kabid/Sekretaris/Pimpinan/Admin.
- [x] Template BNBA menambahkan `kode_kategori` opsional dan validasi classifier.
- [ ] Audit/isi label sheet 2 untuk kategori selain IKU/SPM.
- [ ] Target per kategori/label selain SPM.
- [ ] Regenerasi diagram skema setelah perubahan tabel.

Urutan lanjutan:

1. **B0** — audit label sheet 2, isi kategori yang masih kosong.
2. **Putuskan §10.10** — format final (long/wide), tipe nilai count/rasio, snapshot/incremental, arti per-puskesmas, input bulanan, dan sumber denominator.
3. **A0** — lengkapi master wilayah, level wilayah, dan tipe nilai indikator; mapping desa sudah tersedia (§10.19).
4. **A1** — perbaiki fondasi kalkulasi (filter scope, `SUM` vs `RATIO_OF_SUMS`, dedupe/supersede, source hash per scope) sebelum menambah format baru.
5. **Review template** — untuk upload yang diterima aplikasi sekarang, gunakan [`BNBA-Semua-Satuan.xlsx`](../assets/templates/bnba/final/BNBA-Semua-Satuan.xlsx) (§10.18). `template-bnba-dinamis.xlsx` tetap legacy untuk rasio. Setelah format disetujui Dinas, **A2** — finalisasi template + validasi + versi template; jalankan inventaris file BNBA asli saat tersedia.
6. **A3** — kalkulasi `DESA`/`KECAMATAN`, dashboard drill-down, dan agregasi triwulan/tahun.
7. **A4** — adapter wide opsional dan profil berversi bila Dinas memutuskan perlu.
8. **Fase 2/3** — geo view, hardening, dan regenerasi diagram/dokumen.

Setelah itu, perbarui plan ini dengan temuan format aktual dan keputusan Dinas.
