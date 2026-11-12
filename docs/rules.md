# Aturan Proyek MATA SIDIK

> **Status: MENGIKAT.** Dokumen ini adalah aturan proyek yang harus dipatuhi agent, developer, dan reviewer. Jika ada instruksi yang bertentangan dengan dokumen ini, **berhenti dan minta keputusan manusia**; jangan memilih tafsir sendiri.
>
> Sumber aturan: keputusan pemilik proyek dalam diskusi, `README.md`, `docs/01`–`docs/04`, dan `assets/diagrams/README.md`. Dokumen ini merangkum aturan yang tidak boleh dilanggar agar aplikasi tidak keluar jalur.

## 0. Cara memakai dokumen ini

1. Baca `README.md` di root terlebih dahulu.
2. Baca dokumen ini sampai selesai.
3. Baca seluruh `docs/01`–`docs/04`.
4. Baca `assets/diagrams/README.md` bila menyentuh diagram.
5. Baru baca kode/tests yang relevan.

Jika dua sumber bertentangan, **jangan mengambil kesimpulan diam-diam**. Tunjukkan konflik, lokasi sumber, dan minta keputusan.

---

## 1. Tujuan dan ruang lingkup produk

**MATA SIDIK** (Manajemen Data dan Sistem Informasi Dinas Kesehatan Kabupaten Sumbawa Barat) menggantikan rekap manual dengan alur terstandar:

**input Puskesmas/bidang → verifikasi Dinas → validasi Kabid/Sekretaris → pengolahan sistem → dashboard → keputusan pimpinan → tindak lanjut.**

Ruang lingkup wajib:

| No | Kelompok | Isi |
| --- | --- | --- |
| 1 | Pengelolaan data | Sumber, kamus data, standar indikator, periode, data induk |
| 2 | Pelaporan dan validasi | Kumpul, cek kelengkapan, validasi, kembalikan yang salah, setujui |
| 3 | Pengolahan dan analisis | Rekap, hitung indikator, tren, target vs realisasi, kesenjangan |
| 4 | Penyajian informasi | Dashboard pimpinan, laporan periodik, grafik, peta bila ada, ekspor |
| 5 | Tindak lanjut dan evaluasi | Masalah, rekomendasi, PJ, pantau penyelesaian |

Tujuan yang tidak boleh digeser:

- Data kesehatan terintegrasi, terstandar, dan dapat ditelusuri.
- Kelengkapan, konsistensi, ketepatan waktu, dan keterlacakan data.
- Pemantauan indikator per program, faskes, wilayah, dan periode.
- Laporan cepat tanpa rekap manual.
- Informasi strategis bagi Kepala Dinas/pimpinan.
- Fondasi analisis, integrasi, dan pemantauan tindak lanjut.

**Output acuan:** dashboard SPM dengan **3 baris acuan** yang setiap barisnya dapat ditelusuri sampai `FILE_OBJECT`. Tiga baris ini adalah acuan, bukan pengganti tabel sumber 12 indikator di `docs/03-contoh-data.md`.

---

## 2. Aturan perilaku agent

### 2.1 Larangan keras

- **Jangan berhalusinasi.** Jangan mengarang fakta, angka, tabel, endpoint, status, hasil tes, atau isi file.
- **Jangan mengambil kesimpulan diam-diam.** Jangan menyimpulkan requirement, penyebab bug, keputusan desain, atau angka yang tidak ada di sumber.
- **Jangan menebak keputusan bisnis.** Keputusan yang belum ditetapkan Dinas harus ditanyakan, bukan diisi dengan asumsi.
- **Jangan mengubah scope tanpa persetujuan.** Fitur, dependency, infrastruktur, atau refactor besar di luar permintaan dilarang.
- **Jangan menghapus atau menimpa pekerjaan pengguna** tanpa memahami dan mengonfirmasi.
- **Jangan commit/push kecuali diminta eksplisit.**
- **Jangan mengklaim sudah diuji** tanpa menjalankan pemeriksaan.
- **Jangan menyembunyikan error, konflik, atau keterbatasan.**
- **Jangan memakai Mermaid, Graphviz, atau ASCII art untuk flow.**
- **Jangan mengedit HTML/PNG diagram secara manual**; edit generator lalu rebuild.
- **Jangan mengarang API, tabel, kolom, atau enum baru** tanpa persetujuan.

### 2.2 Wajib bertanya (ASK)

Bertanya sebelum melanjutkan jika:

1. Instruksi ambigu atau punya lebih dari satu tafsir masuk akal.
2. Dua dokumen saling bertentangan.
3. Keputusan bisnis belum ada (lihat §14).
4. Tindakan berpotensi destruktif: hapus file, reset database, rewrite history, migration, mengubah data final.
5. Butuh dependency/layanan/infrastruktur baru.
6. Ada implikasi keamanan, privasi, atau data pribadi.
7. Perubahan memperluas scope atau mengubah kalkulasi/workflow.
8. Verifikasi tidak bisa dijalankan atau gagal.
9. Instruksi pengguna bertentangan dengan aturan ini.

Format pertanyaan: **Pertanyaan → Opsi → Dampak → Rekomendasi**. Jangan melanjutkan dengan asumsi sambil menunggu jawaban.

### 2.3 Label pernyataan

| Label | Arti | Cara menulis |
| --- | --- | --- |
| **Fakta** | Terverifikasi dari file/kode/data | Sertakan `path:line` atau nama tes |
| **Asumsi** | Belum terverifikasi | Awali “Asumsi:” dan minta konfirmasi |
| **Rekomendasi** | Usulan, bukan keputusan | Awali “Rekomendasi:” dan jelaskan trade-off |
| **Belum diketahui** | Tidak ada sumber | Tulis “Belum diketahui” dan ajukan pertanyaan |

### 2.4 Protokol kerja

1. Pahami tujuan dan batasan.
2. Baca sumber wajib.
3. Periksa kondisi nyata (kode, git, file, tes).
4. Buat rencana kecil dan sebutkan file yang akan diubah.
5. Tanya bila perlu.
6. Implementasi minimal sesuai pola yang ada.
7. Verifikasi build/tes/tautan/tabel/diagram.
8. Laporkan perubahan, bukti, verifikasi, asumsi, dan pertanyaan terbuka.

---

## 3. Struktur dokumentasi

Struktur wajib:

```text
README.md                  # wajib di root; pintu masuk Dinas + Developer + Agent
AGENTS.md                  # aturan operasional agent (dibaca otomatis oleh banyak tool)
docs/
  rules.md                 # aturan proyek yang mengikat (dokumen ini)
  01-gambaran-umum.md
  02-daftar-indikator-kinerja.md
  03-contoh-data.md
  04-arsitektur-iku.md
assets/
  diagrams/README.md       # katalog dan cara regenerasi diagram
```

Aturan:

- `README.md` **wajib** ada di root dan memuat bagian **Dinas** serta **Developer**.
- Semua dokumen proyek utama berada di `docs/`; jangan menaruh dokumen proyek baru di root.
- Pengecualian di root: `README.md` dan `AGENTS.md` (instruksi agent). `assets/diagrams/README.md` tetap di folder aset karena melekat pada diagram.
- README wajib memuat peta dokumentasi dan urutan baca agent.
- Setiap tautan relatif harus diuji. Dari `docs/`, aset memakai `../assets/...`; dari README, aset memakai `assets/...`.
- Jika dokumen dipindahkan, perbarui semua tautan, caption, dan referensi teks pada saat yang sama.
- Dokumen yang sudah usang harus diperbarui atau ditandai; jangan dibiarkan menyesatkan.

---

## 4. Batasan teknologi yang tidak boleh dilanggar

| Area | Keputusan | Larangan |
| --- | --- | --- |
| Backend | Go + **Gin** | Jangan ganti framework tanpa persetujuan |
| Frontend | **Svelte 5 + Vite**, build statis satu origin dengan API | Jangan menambah framework UI baru tanpa persetujuan |
| Database | **PostgreSQL 16** | Jangan menambah database lain |
| Arsitektur | **Modular monolith + worker**, satu image, dua command (`api`, `worker`) | Jangan memecah menjadi microservice tanpa persetujuan |
| Queue | Tabel `JOB` dengan `FOR UPDATE SKIP LOCKED`, heartbeat, retry, `dead` + alert | **Tanpa Redis**; jangan menambah broker |
| File | Volume lokal `/data/files`, nama `sha256`, metadata `FILE_OBJECT`, backup volume wajib | **Tanpa object storage** |
| Observability | Log JSON ke stdout, rotasi log container, `/metrics`, `/healthz`, `/readyz` | **Tanpa APM**; jangan menambah layanan observability |
| Deployment | Docker Compose: `api`, `worker`, `postgres`, reverse proxy, volume `pgdata` + `files` | Jangan menambah orchestrator/cloud service tanpa persetujuan |
| Auth | PASETO v4.local akses 30 menit + refresh rotasi (hash di DB), RBAC/scope, MFA TOTP untuk approver | Jangan memakai Keycloak atau auth provider eksternal tanpa persetujuan |
| API | Route group `/api/v1`; SPA lewat `NoRoute` fallback; file hanya lewat handler `/api/v1/files/{id}` | Jangan menyajikan file sebagai static |
| Worker | Worker tidak memakai Gin | Jangan mencampur HTTP handler ke worker |

> **Pilot implementasi (Nov 2026):** atas instruksi pemilik proyek, pilot memakai **SQLite** dan **PASETO v4.local**. Produksi tetap wajib **PostgreSQL 16 + PASETO/RBAC/MFA** sesuai tabel di atas; migrasi dan hardening auth adalah pekerjaan sebelum produksi. Skema pilot: [`schema.md`](schema.md).

**HTTP layer Gin:** `gin.New()` dengan middleware eksplisit (request ID, log JSON, recovery, metrics), binding DTO + validator, dan helper `respondError` untuk format `{code, message, details, request_id}`. Jangan menyimpan `*gin.Context` di service/worker; teruskan `c.Request.Context()`.

**Route group:** `/auth`, `/submissions`, `/jobs`, `/tasks`, `/iku`, `/files`, `/admin` di bawah `/api/v1`.

---

## 5. Upload, validasi, dan klasifikasi

### 5.1 Pipa upload background

1. Svelte mengirim multipart `POST /submissions`.
2. Gin memvalidasi auth + scope + periode.
3. API menghitung SHA-256 dan menyimpan file ke volume lokal.
4. API menulis `SUBMISSION` + `FILE_OBJECT` + `JOB` dalam **satu transaksi DB**.
5. API membalas `202 { submission_id, job_id }`.
6. Worker mengambil job dengan `FOR UPDATE SKIP LOCKED`.
7. Worker mem-parse Excel/CSV, cek header, batas ukuran, dan formula injection.
8. Worker memvalidasi baris + klasifikasi IKU.
9. Worker menulis `SUBMISSION_ROW` + `INDICATOR_VALUE` sementara.
10. Worker memperbarui `progress_json`; status job `DONE`/`FAILED`.
11. UI polling `GET /jobs/{id}` setiap ±2 detik dengan backoff; menampilkan langkah validasi → klasifikasi → simpan → selesai.
12. Jika selesai, UI menampilkan preview + tombol Kirim.

**Mengapa background:** request upload tidak boleh menunggu parsing/klasifikasi; file besar dan validasi baris harus diselesaikan worker.

### 5.2 Kontrak file dan klasifikasi

File agregat minimal memuat:

```text
periode;kode_faskes;kode_iku;kode_indikator;scope_type;scope_code;nilai;numerator;denominator;catatan
2026-05;PKM-TALIWANG;IKU-2026-001.01;IND-001;DESA;5201...;45;450;500;
```

Aturan:

- Upload dinamis **wajib** memuat `kode_iku` + `kode_indikator` + periode.
- `kode_iku` adalah konteks pengelompokan dan validasi routing.
- `kode_indikator` adalah identitas baris data.
- `periode` mengikat ke `MASTER_RELEASE` yang berlaku.
- `numerator`/`denominator` wajib untuk rasio/persentase; `nilai` boleh diturunkan.
- `scope_type`/`scope_code` menentukan desa/puskesmas/kecamatan/kabupaten. Jika kosong, pakai scope uploader **dengan peringatan**.
- Klasifikasi harus **deterministik**: `MATCHED`, `MISMATCH`, `IKU_ROLLUP_ONLY`, `AMBIGUOUS`, `DUPLICATE`.
- Baris ambigu/mismatch/duplikat masuk **karantina**; jangan menebak IKU.
- Hanya IKU leaf/direct-input yang boleh menerima upload. **IKU induk tidak menerima input langsung.**
- Jangan membuat master baru dari file upload. Kode tidak dikenal → karantina.
- BNBA line-level dipetakan lewat **profil dataset**, bukan per baris: profil menentukan kolom numerator/denominator dan indikator turunan.
- BNBA wajib memakai header template Puskesmas; tolak baris tanpa domisili (desa/puskesmas/kecamatan).
- Simpan file asli + hasil parsing untuk audit.
- Tampilkan pesan error per baris saat preview.

### 5.3 Preview dan koreksi IKU

- Preview menampilkan baris yang cocok, status, dan baris karantina.
- Dropdown **“Ubah IKU”** hanya berisi IKU leaf yang valid untuk indikator tersebut.
- Perubahan sebelum submit dicatat di audit.
- Setelah submit, perubahan mengikuti revisi + approval, bukan edit langsung.
- Contoh yang benar: R3 semula `IKU_ROLLUP_ONLY` (kode IKU induk) dikoreksi ke `.03` di preview, tercatat di `AUDIT_LOG`, lalu nilainya masuk ke `.03`.
- Jangan pernah menampilkan angka tanpa status klasifikasi/sumber.

---

## 6. Workflow, approval, dan audit

### 6.1 Status formal (7 nilai)

| Status | Arti | Boleh dilakukan oleh |
| --- | --- | --- |
| `draf` | Baru diinput, belum dikirim | Puskesmas/PJ |
| `dikirim` | Menunggu verifikasi | Sistem (terkunci untuk pengirim) |
| `diverifikasi` | Lolos cek Dinas, menunggu validasi | Verifikator |
| `dikembalikan_petugas` | Belum sesuai, kembali ke tahap 1 + catatan | Verifikator |
| `dikembalikan_verifikator` | Perlu cek, kembali ke tahap 2 | Kabid/Sekretaris |
| `final` | Sah, terkunci, masuk engine + dashboard | Kabid/Sekretaris |
| `diubah_khusus` | Perubahan pasca-final via kewenangan khusus, wajib audit | Kabid/Sekretaris + Admin |

### 6.2 Aturan workflow

- Default: `DRAFT_PKM → VERIFY_DINAS → VALIDATE_KABID → FINAL`.
- Approval bertingkat **dinamis tetapi berversi**; versi workflow di-pin per submission agar perubahan alur tidak mengubah pengajuan berjalan.
- Submission lintas bidang dipecah per pemilik data via `SUBMISSION_GROUP`.
- Return wajib catatan; final wajib catatan penerimaan.
- **Tidak ada auto-approve**, termasuk saat SLA lewat.
- Setiap transisi wajib menulis `AUDIT_LOG`: aktor, waktu, aksi, catatan, diff, versi.
- Final mengunci data dan memicu kalkulasi.
- Perubahan pasca-final hanya lewat `SPECIAL_CHANGE`: alasan wajib, **dual approval** (Kabid/Sekretaris + Admin), dan audit.
- Jangan mengubah data final langsung di database atau UI.

---

## 7. Kalkulasi dan provenance

### 7.1 Rumus dasar

- `target_bulanan = target_tahunan / 12` (contoh 2000 → 166,67; target bulanan persentase = 8,33%).
- `capaian_% = realisasi / target × 100` per periode.
- Simpan **numerator dan denominator**, bukan hanya persentase.
- Tampilkan pembulatan yang konsisten (contoh 2 desimal), tetapi **jangan membulatkan data mentah** yang dipakai kalkulasi.

### 7.2 IKU induk

- IKU induk memakai `RATIO_OF_SUMS`: `Σ numerator / Σ denominator` dari komponen eksplisit `IKU_KOMPONEN`.
- **Bukan rata-rata persentase anak.**
- Setiap anak menyimpan kontribusi di `IKU_VALUE_DETAIL` (`component_type = CHILD_IKU`, numerator, denominator, bobot).
- IKU leaf diturunkan dari `INDICATOR_VALUE`; IKU induk diturunkan dari detail anak.
- Status nilai: `COMPLETE`, `PARTIAL`, `NO_DATA`, `STALE`.
- Komponen belum lengkap = `PARTIAL`, **bukan nol**.
- Jangan double counting; jangan memakai komponen yang tidak eksplisit.

### 7.3 Contoh acuan (wajib konsisten)

| Komponen | Perhitungan | Hasil |
| --- | --- | --- |
| `.01` Ibu Hamil | 450 / 500 | **90,00%** |
| `.02` Ibu Bersalin | 380 / 400 | **95,00%** |
| `.03` Bayi Baru Lahir | 20 / 20 | **100,00%** |
| Induk `IKU-2026-001` | (450 + 380 + 20) / (500 + 400 + 20) = 850 / 920 | **92,39%** |
| Target | 100% | gap **-7,61 pp** |

Contoh koreksi: jika R3 dibiarkan di IKU induk, `.03` = `NO_DATA` dan induk = 830/900 = 92,22%; setelah dikoreksi ke `.03`, induk = 850/920 = 92,39%. Ini contoh untuk menjelaskan pentingnya klasifikasi, bukan data produksi.

### 7.4 Pemicu recalculation

- Submission menjadi final.
- Koreksi pasca-final disetujui.
- Target berubah.
- Formula/hierarki IKU dipublikasikan.
- Periode dibuka kembali.

Hanya subtree IKU terdampak yang dihitung ulang. Dashboard membaca nilai **materialized**, bukan menghitung ulang saat request. Setiap run menyimpan `calculation_run_id` dan `source_hash` agar idempoten.

---

## 8. Dashboard SPM dan traceability

### 8.1 Kategori dashboard

Sumber `docs/01-gambaran-umum.md` menyebut **9 kategori**:

1. IKU (Tujuan/Sasaran)
2. IKK
3. SPM
4. RIBK
5. Program
6. Kegiatan
7. Sub Kegiatan
8. Bidang/Sekretariat/Tim Kerja
9. Baris `........................` — **belum didefinisikan**

**Jangan mengisi kategori ke-9 dengan tebakan.** Delapan kategori pertama pasti; kategori ke-9 menunggu keputusan Dinas.

### 8.2 Tampilan dan filter

- Setiap kartu menampilkan **Target, Capaian, Status, Tren**.
- Filter minimal: faskes, bidang, program, indikator, periode, wilayah.
- Status dashboard pilot dihitung dari **target periode**: capaian ≥ target = sesuai target; di bawah target = perlu perhatian. Ambang warna/status khusus belum diputuskan Dinas.
- Capaian tahun sebelumnya diisi admin per indikator/IKU untuk scope kabupaten. Dashboard periode tahun Y membaca capaian tahun Y-1; nilai tidak disalin otomatis dan tidak menimpa histori.
- Dashboard tidak boleh menampilkan angka yang tidak punya provenance.

### 8.3 Tiga baris SPM acuan

| Baris | Indikator | IKU leaf | Baris file | Num/den | Capaian Mei 2026 | Status |
| --- | --- | --- | --- | ---: | ---: | --- |
| SPM 01 | Pelayanan Kesehatan Ibu Hamil | `.01` | R1 | 450/500 | **90,00%** | di bawah target |
| SPM 02 | Pelayanan Kesehatan Ibu Bersalin | `.02` | R2 | 380/400 | **95,00%** | di bawah target |
| SPM 03 | Pelayanan Kesehatan Bayi Baru Lahir | `.03` | R3 | 20/20 | **100,00%** | sesuai target |

Ketiganya berasal dari satu file: `FILE_OBJECT` hash `9f2c…a71`, `SUBMISSION S-2026-05-001`, `REVISION 1`, uploader Siti (PKM Taliwang), periode Mei 2026.

### 8.4 Rantai trace wajib

`DASHBOARD_ROW → IKU_VALUE → IKU_VALUE_DETAIL (untuk agregat/induk) → INDICATOR_VALUE → SUBMISSION_ROW → SUBMISSION_REVISION → SUBMISSION → FILE_OBJECT`

Setiap baris SPM membawa:

- `iku_value_id`
- `indicator_value_id`
- `source_hash`
- `data_version`
- `calculation_run_id`

API drill-down:

- `GET /api/v1/dashboard/spm?periode=2026-05`
- `GET /api/v1/dashboard/spm/{kode}/trace?periode=2026-05`
- `GET /api/v1/files/{id}`

Nilai SPM adalah **proyeksi** dari nilai terhitung, bukan angka yang diketik ulang. Jika `IKU_VALUE` belum ada, status `NO_DATA`, bukan nol. Jika komponen belum lengkap, `PARTIAL`. Jika data berubah tapi kalkulasi belum selesai, `STALE`.

---

## 9. Model domain dan master data

### 9.1 Tiga konsep yang harus dipisahkan

| Konsep | Fungsi | Contoh |
| --- | --- | --- |
| **Label** | Taksonomi dan tampilan dashboard | `IKU`, `Indikator Sasaran`, `IKK`, `SPM`, `RIBK`, `Program` |
| **IKU node** | Label bertipe IKU yang punya hierarki, formula, target, pemilik | `IKU-2026-001` → `IKU-2026-001.01` |
| **Indikator** | Definisi data yang diukur | `IND-001` |

- Satu indikator dapat menempel ke banyak label; nilainya tetap **satu** di `INDICATOR_VALUE`.
- Label hanya menentukan dashboard dan tanggung jawab, bukan menduplikasi nilai.
- Taksonomi label dipisahkan dari graf kalkulasi (`IKU_KOMPONEN`).

### 9.2 Master berversi

- `MASTER_RELEASE` membundel master per tahun; upload mengikat ke release yang berlaku.
- Tambah indikator = **versi baru**; nonaktifkan = `berlaku_sampai` + `status_aktif=false`.
- **Jangan hapus fisik** data master agar tren historis utuh.
- Setiap versi indikator wajib mengisi: satuan baku, definisi + numerator/denominator, label (many-to-many via `INDIKATOR_LABEL`), sifat input, frekuensi, level wilayah, PJ bidang/timker.
- Kode: `IKU-2026-001`, `IND-001`; kode stabil, versi terpisah.

### 9.3 Agregasi wilayah dan sifat data

- Alur agregasi: **Desa → Puskesmas → Kecamatan → Kabupaten**.
- Kasus kematian ibu/balita wajib bisa di-drill-down.
- UHH hanya level kabupaten, input 1×/tahun.
- BNBA line-level dipetakan lewat profil dataset, bukan per baris.
- Scope kosong memakai scope uploader dengan peringatan.

---

## 10. RBAC, keamanan, dan audit

### 10.1 Matriks hak akses

| Aksi | Puskesmas/PJ | Verifikator Dinas | Kabid/Sekretaris | Pimpinan | Admin |
| --- | --- | --- | --- | --- | --- |
| Input/kirim + upload BNBA | ✅ | ❌ | ❌ | ❌ | ✅ (bantuan) |
| Verifikasi/kembalikan ke pengirim | ❌ | ✅ | ❌ | ❌ | ❌ |
| Validasi final/kembalikan ke verifikator | ❌ | ❌ | ✅ | ❌ | ❌ |
| Ubah data final (khusus + audit) | ❌ | ❌ | ✅ terbatas | ❌ | ✅ dual approval |
| Isi faktor pendukung/penghambat | ❌ | ❌ | ✅ wajib | ❌ | ❌ |
| Lihat dashboard semua filter | milik sendiri | bidangnya | bidangnya | ✅ semua | ✅ |
| Beri arahan / tetapkan TL | ❌ | ❌ | terima | ✅ beri | ❌ |
| Kelola master + user | ❌ | ❌ | ❌ | ❌ | ✅ |

### 10.2 Keamanan

- PASETO v4.local akses 30 menit + refresh rotasi; hash refresh token di DB.
- RBAC/scope di route group; cek kepemilikan submission untuk `/jobs`.
- MFA TOTP untuk approver.
- Semua user dapat mengganti password sendiri dengan verifikasi password lama; admin tetap dapat melakukan reset. Perubahan password mencabut seluruh refresh session user.
- File hanya lewat handler dengan auth/scope; jangan static.
- Validasi multipart dengan `http.MaxBytesReader`; tolak sebelum file disimpan.
- Cegah formula injection saat parse Excel/CSV.
- Jangan menaruh secret/kredensial di repo.
- PII di BNBA **belum diputuskan**; jangan berasumsi boleh menyimpan/menampilkan semua kolom.

### 10.3 Audit

Setiap keputusan/transisi wajib mencatat: **siapa, kapan, apa, catatan, diff, versi**. Audit tidak boleh dihapus. Data final terkunci; perubahan lewat mekanisme khusus.

---

## 11. Aturan diagram dan visualisasi

### 11.1 Larangan

- **Tidak boleh Mermaid.**
- **Tidak boleh Graphviz/DOT.**
- **Tidak boleh ASCII art untuk flow.**
- Semua visualisasi alur/flow/arsitektur memakai **stack-diagram**.

### 11.2 Diagram yang ada

| Diagram | Dipakai di | Isi |
| --- | --- | --- |
| `alur-bisnis` | README A.2, `docs/01` §5 | Enam tahap bisnis + pengembalian + tindak lanjut |
| `status-data` | README A.2 | Siklus status + dua jalur pengembalian |
| `alur-peran` | README B.2 | Swimlane lima peran |
| `model-data` | README B.3 | Master vs transaksi vs audit |
| `alur-data-iku` | README A.7, `docs/04` §10 | End-to-end IKU + 3 baris SPM |
| `alur-upload` | `docs/04` §3.2 | Upload background |
| `alur-workflow` | `docs/04` §4.1 | State machine approval |
| `alur-provenance` | `docs/04` §5.3 | Rollup IKU induk |
| `alur-trace-spm` | `docs/04` §5.5 | Trace SPM sampai file |
| `arsitektur-aplikasi` | `docs/04` §6.1, README B.8 | Arsitektur aplikasi |
| `schema-db` | `docs/schema.md` | Katalog tabel, kolom, dan relasi database |

### 11.3 Aturan pembuatan dan pemeliharaan

- Diagram alur wajib memperlihatkan: **upload → validasi/klasifikasi → preview → workflow → kalkulasi → output SPM**, dengan 3 baris SPM yang dapat ditelusuri ke file.
- Setiap diagram punya satu HTML mandiri + satu PNG.
- Sertakan caption, alt text, dan tautan ke HTML sumber.
- Dari `docs/`, gambar memakai `../assets/diagrams/...`; dari README, `assets/diagrams/...`.
- Generator: `gen.py` (alur bisnis, status, peran, model data), `gen-end-to-end.py` (alur data IKU), `gen-teknis.py` (lima diagram teknis).
- Build: `build.sh` (fetch ikon → generate → inline font/ikon → render PNG).
- **Jangan mengedit HTML/PNG hasil render manual.** Edit generator lalu jalankan build.
- Jangan mengubah diagram tanpa memperbarui dokumen yang memakainya, dan sebaliknya.
- Ikon memakai Lucide; font Geist + JetBrains Mono; jangan mencampur gaya diagram lain.
- Setelah render, periksa label terpotong, panah menabrak teks, dan kontras/keterbacaan.
- Pengecualian: **pohon taksonomi** (`docs/04` §2.1) dan **route/package tree** (`docs/04` §6.2.1) bukan flow; keduanya tetap code block agar mudah disalin.

---

## 12. Aturan dokumentasi dan penulisan

- Bahasa dokumen: **Indonesia**; istilah teknis/kode boleh Inggris dan harus konsisten.
- Gunakan tabel untuk data terstruktur; jangan menaruh data tabular dalam prosa panjang.
- Setiap klaim angka harus punya sumber (`docs/...`, tabel, atau diagram).
- Tabel mentah Excel dipertahankan apa adanya; jangan menghapus kolom kosong karena itu penanda PR.
- Jangan mengisi definisi operasional, satuan, atau PJ yang kosong dengan tebakan.
- Tandai keputusan terbuka dengan jelas: **“Belum diputuskan”** + siapa yang harus memutuskan.
- Caption gambar harus menjelaskan isi dan sumber HTML.
- Tautan relatif harus valid; anchor heading harus cocok.
- README wajib memuat peta dokumentasi dan bagian Dinas + Developer.
- Setiap perubahan aturan harus memperbarui `docs/rules.md`; jangan menyimpan aturan hanya di chat.
- Jangan menambahkan dokumen baru di root selain pengecualian §3.

---

## 13. Keputusan final yang sudah dikunci

| No | Keputusan | Konsekuensi |
| --- | --- | --- |
| 1 | IKU adalah label many-to-many, bukan kolom tunggal | Satu indikator bisa punya banyak label; nilai tetap satu |
| 2 | Taksonomi label dipisah dari graf kalkulasi | `IKU_KOMPONEN` menentukan angka induk |
| 3 | Upload memakai tiga kunci: `kode_iku`, `kode_indikator`, `periode` | Kode IKU saja tidak cukup |
| 4 | Klasifikasi deterministik; ragu masuk karantina | Tidak ada data masuk IKU secara diam-diam |
| 5 | Hanya IKU leaf/direct-input yang boleh menerima upload | IKU induk murni rollup |
| 6 | IKU induk memakai komponen eksplisit + `RATIO_OF_SUMS` | Mencegah double counting dan rata-rata persentase |
| 7 | Workflow bertingkat dinamis tetapi berversi | Versi di-pin per submission |
| 8 | Setiap keputusan menyimpan aktor, waktu, catatan, diff, versi | Audit dan pertanggungjawaban |
| 9 | Modular monolith + worker, PostgreSQL, volume lokal | Tanpa Redis, tanpa object storage |
| 10 | Backend Go + Gin, frontend Svelte, auth tanpa Keycloak | Deployment sederhana, satu origin |
| 11 | Default approval tiga level | `DRAFT_PKM → VERIFY_DINAS → VALIDATE_KABID → FINAL` |
| 12 | Semua flow memakai stack-diagram | Tanpa Mermaid/Graphviz/ASCII flow |
| 13 | README di root; dokumen proyek di `docs/` | Tautan dan aset harus disesuaikan |
| 14 | R3 contoh dikoreksi ke `.03` | Capaian 90,00/95,00/100,00; induk 92,39% |
| 15 | Tiga baris SPM harus traceable sampai `FILE_OBJECT` | Setiap baris membawa ID kontrak dan `source_hash` |

---

## 14. Keputusan terbuka — jangan ditebak

| # | Pertanyaan | Status | Tindakan |
| --- | --- | --- | --- |
| 1 | Kategori dashboard ke-9 (`........................`) | Belum diputuskan | Tanya Dinas; definisikan atau hapus |
| 2 | Ambang warna/status khusus (✅/🟡/🔴) | Belum diputuskan | Dashboard pilot memakai status dasar berbasis target (sesuai/di bawah target); ambang warna formal menunggu Dinas |
| 3 | SLA input/verifikasi/validasi per periode | Belum diputuskan | Tanya Dinas; tanpa auto-approve |
| 4 | Target sama untuk semua faskes atau berbeda per wilayah | Belum diputuskan | Tanya Dinas |
| 5 | PII di BNBA dan kebijakan masking/enkripsi | Belum diputuskan | Tanya Dinas; batasi akses sementara |
| 6 | Retensi file dan backup | Belum diputuskan | Tanya Dinas; minimal sesuai audit |
| 7 | Format BNBA final | Belum diputuskan | Tentukan template resmi sebelum implementasi parser |
| 8 | Duplikasi 3 baris ibu hamil di katalog | Belum diputuskan | Gabung atau bedakan (TW1/Mei/absolut?) |
| 9 | Satuan baku (`%` vs `Persentase` vs `(%)`, dll.) | Belum diputuskan | Normalisasi sebelum engine |
| 10 | PJ Data untuk semua indikator | Belum lengkap | Lengkapi Bidang + Timker/Subbag |
| 11 | Level approval tambahan di luar default 3 | Belum diputuskan | Konfigurasi hanya jika disetujui |
| 12 | Integrasi SIKN | Belum diputuskan | Jangan implementasi tanpa spesifikasi |
| 13 | Ekspor PDF/Excel dan peta | Prioritas P3, detail belum | Tanya Dinas |
| 14 | Notifikasi SLA | Prioritas P3, kanal belum | Tanya Dinas |

Jika pekerjaan menyentuh salah satu item di atas, **berhenti dan tanyakan**. Jangan memilih default yang terlihat masuk akal.

---

## 15. Checklist verifikasi sebelum menyatakan selesai

- [ ] Requirement yang diminta terpenuhi dan dapat ditelusuri ke sumber.
- [ ] Tidak ada klaim tanpa bukti; asumsi dilabeli atau dikonfirmasi.
- [ ] Tidak ada keputusan bisnis yang ditebak.
- [ ] Tautan file, anchor, dan gambar valid.
- [ ] Tabel konsisten jumlah kolom.
- [ ] Tidak ada Mermaid/Graphviz/ASCII flow baru.
- [ ] Diagram dirender ulang bila generator berubah; PNG ada.
- [ ] Kalkulasi sesuai §7 dan traceability sesuai §8.
- [ ] Build/tes yang relevan dijalankan.
- [ ] `git diff --check` bersih.
- [ ] Tidak ada commit/push kecuali diminta.
- [ ] Laporan menyebutkan perubahan, verifikasi, asumsi, dan pertanyaan terbuka.

---

## 16. Definisi istilah

| Istilah | Arti |
| --- | --- |
| **IKU** | Indikator Kinerja Utama; label dan node hierarki yang dapat dihitung |
| **IKU leaf** | IKU paling bawah/direct-input yang boleh menerima upload |
| **IKU induk** | IKU rollup yang dihitung dari `IKU_KOMPONEN`, tidak menerima upload langsung |
| **IKK** | Indikator Kinerja Kegiatan |
| **SPM** | Standar Pelayanan Minimal; salah satu label dashboard |
| **RIBK** | Rencana Induk Bidang Kesehatan |
| **BNBA** | Buku/berkas data yang diupload Puskesmas untuk input agregat |
| **Provenance** | Asal-usul angka: komponen, file, baris, aktor, versi |
| **RATIO_OF_SUMS** | `Σ numerator / Σ denominator`, bukan rata-rata persentase |
| **Karantina** | Tempat baris ambigu/mismatch/duplikat yang tidak boleh masuk IKU |
| **`SPECIAL_CHANGE`** | Perubahan pasca-final dengan alasan + dual approval + audit |
| **`FILE_OBJECT`** | Metadata file asli: hash SHA-256, nama, path, uploader |
| **`SUBMISSION_REVISION`** | Revisi submission yang immutable untuk audit |
| **`DASHBOARD_ROW`** | Baris tampilan SPM yang membawa ID kontrak untuk drill-down |

---

## 17. Referensi

- [`README.md`](../README.md) — pintu masuk dan ringkasan Dinas + Developer.
- [`docs/01-gambaran-umum.md`](01-gambaran-umum.md) — proses bisnis, ruang lingkup, tabel mentah.
- [`docs/02-daftar-indikator-kinerja.md`](02-daftar-indikator-kinerja.md) — katalog indikator.
- [`docs/03-contoh-data.md`](03-contoh-data.md) — contoh dashboard SPM dan data absolut.
- [`docs/04-arsitektur-iku.md`](04-arsitektur-iku.md) — rancangan teknis lengkap.
- [`docs/schema.md`](schema.md) — skema database lengkap.
- [`docs/implementation.md`](implementation.md) — struktur dan cara run implementasi pilot.
- [`assets/diagrams/README.md`](../assets/diagrams/README.md) — katalog dan regenerasi diagram.
- [`AGENTS.md`](../AGENTS.md) — aturan operasional agent.
