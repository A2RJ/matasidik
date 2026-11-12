# AGENTS.md — Aturan Operasional Agent MATA SIDIK

> Dokumen ini mengikat semua agent yang bekerja di repositori ini. Aturan proyek yang lebih lengkap ada di [`docs/rules.md`](docs/rules.md). Jika isi dokumen ini bertentangan dengan `docs/rules.md`, **berhenti dan tanyakan manusia**; jangan memilih salah satu diam-diam.

## 1. Urutan baca wajib

Sebelum menganalisis, mengubah, atau menyimpulkan apa pun:

1. Baca [`README.md`](README.md) sampai selesai. README adalah pintu masuk dan memuat peta seluruh dokumentasi.
2. Baca [`docs/rules.md`](docs/rules.md) sampai selesai. Ini aturan proyek yang mengikat.
3. Baca **seluruh** file `.md` di [`docs/`](docs/):
   - [`docs/01-gambaran-umum.md`](docs/01-gambaran-umum.md)
   - [`docs/02-daftar-indikator-kinerja.md`](docs/02-daftar-indikator-kinerja.md)
   - [`docs/03-contoh-data.md`](docs/03-contoh-data.md)
   - [`docs/04-arsitektur-iku.md`](docs/04-arsitektur-iku.md)
4. Baca [`assets/diagrams/README.md`](assets/diagrams/README.md) jika tugas menyentuh diagram.
5. Baca kode, tes, konfigurasi, dan data yang relevan dengan tugas.

Jangan mengandalkan ringkasan sesi sebelumnya, ingatan model, atau potongan file. Jika file berubah, baca ulang. Sebutkan file yang sudah dibaca di laporan akhir.

## 2. Larangan keras

- **Jangan berhalusinasi.** Jangan mengarang fakta, angka, nama tabel, endpoint, status, aturan bisnis, hasil tes, atau isi file. Jika belum diverifikasi, tulis bahwa itu belum diverifikasi.
- **Jangan mengambil kesimpulan diam-diam.** Jangan menyimpulkan requirements, penyebab bug, keputusan desain, atau angka yang tidak ada di sumber. Jika terpaksa menyimpulkan, beri label **“Asumsi:”** dan minta konfirmasi sebelum dipakai.
- **Jangan mengisi keputusan bisnis yang masih terbuka.** Kategori dashboard ke-9, ambang status, SLA, target per wilayah, PII, retensi file, dan keputusan lain di `docs/rules.md` harus ditanyakan kepada Dinas, bukan ditebak.
- **Jangan mengubah scope tanpa persetujuan.** Jangan menambah fitur, dependency, infrastruktur, atau refactor besar yang tidak diminta.
- **Jangan menghapus atau menimpa pekerjaan pengguna.** File, perubahan, atau keputusan yang tidak dipahami harus diselidiki dulu; konfirmasi sebelum tindakan destruktif.
- **Jangan commit atau push kecuali diminta secara eksplisit.**
- **Jangan mengklaim sudah diuji/diverifikasi** tanpa benar-benar menjalankan pemeriksaan yang relevan.
- **Jangan menyembunyikan error atau konflik.** Laporkan kegagalan, keterbatasan, dan ketidakpastian apa adanya.
- **Jangan memakai Mermaid, Graphviz, atau ASCII art untuk visualisasi alur/flow.** Semua flow wajib memakai stack-diagram sesuai `docs/rules.md`.
- **Jangan mengedit HTML/PNG diagram secara manual.** Edit generator (`gen.py`, `gen-end-to-end.py`, `gen-teknis.py`), lalu rebuild.
- **Jangan mengarang API atau skema database.** Ikuti kontrak yang ada di `docs/04-arsitektur-iku.md`; jika butuh yang baru, ajukan usulan dan minta persetujuan.

## 3. ASK — kapan wajib bertanya

Berhenti dan ajukan pertanyaan sebelum melanjutkan jika:

1. Instruksi pengguna ambigu atau punya lebih dari satu tafsir yang masuk akal.
2. Dua dokumen saling bertentangan.
3. Keputusan bisnis belum ada (lihat daftar keputusan terbuka di `docs/rules.md`).
4. Tindakan berpotensi merusak: hapus file, reset database, rewrite history, migration, atau mengubah data final.
5. Dibutuhkan dependency, layanan, atau infrastruktur baru.
6. Ada implikasi keamanan, privasi, atau data pribadi.
7. Perubahan akan memperluas scope atau mengubah aturan kalkulasi/workflow.
8. Verifikasi tidak bisa dijalankan atau hasilnya gagal.
9. Instruksi pengguna bertentangan dengan `docs/rules.md`.

Format pertanyaan:

- **Pertanyaan:** satu pertanyaan spesifik.
- **Opsi:** 2–4 pilihan konkret.
- **Dampak:** konsekuensi tiap pilihan.
- **Rekomendasi:** pilihan yang disarankan dan alasannya.

Jangan melanjutkan dengan asumsi sambil menunggu jawaban. Jika pengguna tidak menjawab, hentikan pekerjaan yang bergantung pada keputusan itu.

## 4. Bedakan fakta, asumsi, dan rekomendasi

Setiap pernyataan penting harus jelas statusnya:

| Label | Arti | Cara menulis |
| --- | --- | --- |
| **Fakta** | Terverifikasi dari file/kode/data | Sertakan `path:line` atau nama tes |
| **Asumsi** | Belum terverifikasi, dipakai sementara | Awali dengan “Asumsi:” dan minta konfirmasi |
| **Rekomendasi** | Usulan, bukan keputusan | Awali dengan “Rekomendasi:” dan jelaskan trade-off |
| **Belum diketahui** | Tidak ada sumber | Tulis “Belum diketahui” dan ajukan pertanyaan |

Jangan mencampur opini dengan fakta. Jangan menyajikan rekomendasi sebagai keputusan final.

## 5. Protokol kerja

1. **Pahami tujuan.** Ulangi tujuan dan batasan dengan kata sendiri.
2. **Baca sumber wajib.** Ikuti urutan baca di bagian 1.
3. **Periksa kondisi nyata.** Lihat kode, status git, file, dan tes yang ada. Jangan mengandalkan deskripsi.
4. **Buat rencana kecil.** Sebutkan file yang akan diubah dan cara verifikasinya.
5. **Tanya bila perlu.** Gunakan aturan ASK di bagian 3.
6. **Implementasi minimal.** Ikuti pola, gaya, dan struktur yang sudah ada. Jangan refactor di luar kebutuhan.
7. **Verifikasi.** Jalankan pemeriksaan yang relevan: build, tes, link/anchor, tabel, diagram, dan pencarian Mermaid/ASCII flow.
8. **Laporkan.** Sebutkan file yang berubah, bukti, hasil verifikasi, asumsi, dan pertanyaan terbuka.

## 6. Ringkasan aturan proyek yang tidak boleh dilanggar

Ini ringkasan; detail dan rasional ada di [`docs/rules.md`](docs/rules.md).

- **Produk:** MATA SIDIK untuk Dinas Kesehatan Kab. Sumbawa Barat; output utama adalah dashboard SPM yang dapat ditelusuri sampai file upload.
- **Stack:** backend **Go + Gin**, frontend **Svelte 5 + Vite**, database **PostgreSQL 16**, satu image dengan dua command (`api`, `worker`), **tanpa Redis**, **tanpa object storage**, **tanpa APM**. **Pilot Nov 2026:** SQLite + PASETO v4.local atas instruksi pemilik proyek; produksi tetap PostgreSQL + PASETO/RBAC/MFA.
- **Upload:** background job di tabel `JOB` dengan `FOR UPDATE SKIP LOCKED`; file disimpan di volume lokal dengan nama `sha256`; API membalas `202 {submission_id, job_id}`; UI polling `GET /jobs/{id}`.
- **Klasifikasi:** wajib `kode_iku` + `kode_indikator` + `periode`; deterministik; baris ambigu/mismatch masuk karantina; IKU induk tidak menerima upload langsung.
- **Workflow:** default `DRAFT_PKM → VERIFY_DINAS → VALIDATE_KABID → FINAL`; dinamis tetapi versinya di-pin per submission; return wajib catatan; **tidak ada auto-approve**; final terkunci; perubahan pasca-final lewat `SPECIAL_CHANGE` + alasan + dual approval + audit.
- **Kalkulasi:** `target_bulanan = target_tahunan / 12`; `capaian = realisasi / target × 100`; simpan numerator/denominator; IKU induk memakai `RATIO_OF_SUMS` dari komponen eksplisit, **bukan rata-rata persentase**; provenance di `IKU_VALUE_DETAIL`; komponen hilang = `PARTIAL`, bukan nol.
- **Traceability:** `DASHBOARD_ROW → IKU_VALUE → IKU_VALUE_DETAIL → INDICATOR_VALUE → SUBMISSION_ROW → SUBMISSION_REVISION → SUBMISSION → FILE_OBJECT`; setiap baris SPM membawa `iku_value_id`, `indicator_value_id`, `source_hash`, `data_version`, dan `calculation_run_id`.
- **Dokumentasi:** `README.md` wajib di root; dokumen proyek lain di `docs/`; README memuat bagian Dinas dan Developer; semua tautan dan gambar harus valid.
- **Diagram:** semua flow memakai stack-diagram; sumber HTML + PNG; caption dan alt text; jangan mengedit hasil render manual.
- **Git:** jangan commit sebelum diminta; jangan campur perubahan yang tidak terkait.

## 7. Definition of done

Pekerjaan dianggap selesai hanya jika:

- Semua requirement yang diminta terpenuhi dan dapat ditelusuri ke sumber.
- Tidak ada klaim tanpa bukti; asumsi sudah dilabeli atau dikonfirmasi.
- Tidak ada keputusan bisnis yang ditebak.
- Dokumen dan diagram yang terdampak sudah diperbarui.
- Tautan, anchor, tabel, build, dan tes yang relevan sudah diverifikasi.
- Tidak ada Mermaid atau flow ASCII baru.
- Tidak ada commit/push kecuali diminta.
- Laporan akhir menyebutkan perubahan, verifikasi, asumsi, dan pertanyaan terbuka.

## 8. Konflik dan sumber kebenaran

- `docs/rules.md` adalah aturan proyek yang mengikat.
- `README.md` adalah pintu masuk dan ringkasan untuk Dinas + Developer.
- `docs/01`–`docs/04` adalah sumber detail bisnis dan teknis.
- `assets/diagrams/README.md` adalah sumber untuk katalog dan regenerasi diagram.
- Jika dua sumber bertentangan, **jangan memilih diam-diam**. Tunjukkan konflik, sebutkan lokasi masing-masing, dan minta keputusan.
