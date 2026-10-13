# PROSES BISNIS MATA SIDIK

> Sumber: `Proses Bisnis Mata Sidik.xlsx` — Sheet **Gambaran Umum** (A1:F52)
>
> Dokumen ini adalah ekstraksi 1 dari 3 sheet. Format tabel asli dipertahankan pada bagian Tabel Mentah di bawah.

## 1. Nama Aplikasi

**MATA SIDIK (Manajemen Data dan Sistem Informasi Dinas Kesehatan Kabupaten Sumbawa Barat)**

## 2. Maksud

| No | Uraian |
| --- | --- |
| 1. | Mewujudkan pengelolaan data kesehatan yang terintegrasi dan terstandar. |
| 2. | Meningkatkan kelengkapan, konsistensi, ketepatan waktu, dan keterlacakan data. |
| 3. | Memudahkan pemantauan indikator kesehatan berdasarkan program, fasilitas kesehatan, wilayah, dan periode. |
| 4. | Mempercepat penyusunan laporan serta mengurangi pekerjaan rekapitulasi manual. |
| 5. | Menyediakan informasi strategis bagi Kepala Dinas dan pimpinan terkait. |
| 6. | Menjadi fondasi bagi pengembangan analisis, integrasi aplikasi, dan pemantauan tindak lanjut pada tahap berikutnya. |

## 3. Tujuan

| No | Uraian |
| --- | --- |
| 1. | Mengintegrasikan data kesehatan dari Puskesmas, bidang, program, dan unit kerja terkait |
| 2. | Menstandarkan format, definisi indikator, periode pelaporan, dan sumber data |
| 3. | Memastikan data diperiksa dan divalidasi sebelum digunakan |
| 4. | Menyediakan dashboard capaian indikator dan analisis kinerja kesehatan |
| 5. | Mendukung perencanaan program, penganggaran, monitoring, evaluasi, serta pengambilan keputusan pimpinan |
| 6. | Menyediakan riwayat perubahan data dan pengendalian hak akses pengguna |

## 4. Ruang Lingkup

| No | Kelompok | Uraian |
| --- | --- | --- |
| 1. | Pengelolaan Data | Pendataan sumber, kamus data, standar indikator, periode pelaporan, dan pengelolaan data induk |
| 2. | Pelaporan dan validasi | Pengumpulan, pemeriksaan kelengkapan, validasi, pengembalian data yang salah, dan persetujuan data. |
| 3. | Pengolahan dan analisis | Rekapitulasi, perhitungan indikator, analisis tren, perbandingan target dan realisasi, serta identifikasi kesenjangan. |
| 4. | Penyajian informasi | Dashboard pimpinan, laporan periodik, grafik capaian, peta sebaran bila tersedia, dan ekspor laporan. |
| 5. | Tindak lanjut dan evaluasi | Identifikasi masalah, rekomendasi tindak lanjut, penetapan penanggung jawab, dan pemantauan penyelesaian. |

## 5. Flow Chart / Alur Proses

### 5.1 Pengumpulan Data, Input dan Pengiriman Data

User disini adalah Petugas Puskesmas, Petugas Dinas yang menjadi penanggung jawab program/indikator.

### 5.2 Pemeriksaan dan Verifikasi

User disini adalah petugas Dinas

> Jika data belum sesuai maka dikembalikan kepada user nomor 1, jika data sesuai maka disetujui sesuai kewenangan

### 5.3 Validasi

User disini adalah Kepala Bidang/Sekretaris Dinas

> Jika data masih perlu di cek, dikembalikan ke user nomor 2, jika dianggap valid maka disetujui menjadi data final yang tidak dapat dirubah kecuali melalui mekanisme yang ditentukan.

> Data final : perubahan memerlukan kewenangan khusus dan jejak audit.<br>Status tersebut sebaiknya dapat difilter menurut fasilitas kesehatan, bidang, program, indikator, dan periode pelaporan. Sistem juga perlu mencatat siapa yang mengirim, memeriksa, menyetujui, atau mengubah data beserta waktu tindakannya.

### 5.4 Integrasi dan Pengolahan

Penyimpanan terpusat, rekapitulasi, perhitungan indikator, dan analisis

> Pada tahap ini akan muncul hasil pengolahan data berupa tren capaian indikator dan permintaan faktor pendukung dan faktor penghambat capaian indikator yang perlu dilengkapi oleh kepala bidang/sekretaris dinas.

### 5.5 Dashboard dan Pelaporan

Informasi capaian, tren, target, kesenjangan, dan laporan periodik

bagian pelaporan akan menampilkan data capaian yang bisa di filter sesuai pengelompokan berikut :

| No | Pengelompokan Laporan |
| --- | --- |
| 1. | Capaian Kinerja Indikator Kinerja Utama (IKU) (Indikator Tujuan dan Sasaran) |
| 2. | Capaian Kinerja Indikator Kinerja Kunci (IKK) |
| 3. | Capaian Kinerja Standar Pelayanan Minimal (SPM) |
| 4. | Capaian Kinerja Indikator Kinerja RIBK |
| 5. | Capaian Kinerja Indikator Kinerja Program |
| 6. | Capaian Kinerja Indikator Kinerja Kegiatan |
| 7. | Capaian Kinerja Indikator Kinerja Sub Kegiatan |
| 8. | Capaian Kinerja Indikator Kinerja Bidang/Sekretariat/Tim Kerja |
| 9. | ........................ |

### 5.6 Pengambilan Keputusan

Arahan pimpinan, rencana tindak lanjut, monitoring, dan evaluasi

> Arahan Pimpinan ditujukan kepada Kepala Bidang/Sekretaris Dinas

---

## Lampiran: Tabel Mentah Excel (format asli dipertahankan)

Struktur kolom Excel asli: `A=No` \| `B=Aspek` \| `C=:` \| `D=Sub-No` \| `E-F=Uraian (merge)`.
 Sel merge E:F digabung menjadi satu kolom Uraian. Baris kosong dipertahankan apa adanya.

| No | Aspek | : | Sub-No | Uraian |
| --- | --- | --- | --- | --- |
| 1 | Nama Aplikasi | : |  | MATA SIDIK (Manajemen Data dan Sistem Informasi Dinas Kesehatan Kabupaten Sumbawa Barat) |
| 2 | Maksud | : | 1. | Mewujudkan pengelolaan data kesehatan yang terintegrasi dan terstandar. |
|  |  |  | 2. | Meningkatkan kelengkapan, konsistensi, ketepatan waktu, dan keterlacakan data. |
|  |  |  | 3. | Memudahkan pemantauan indikator kesehatan berdasarkan program, fasilitas kesehatan, wilayah, dan periode. |
|  |  |  | 4. | Mempercepat penyusunan laporan serta mengurangi pekerjaan rekapitulasi manual. |
|  |  |  | 5. | Menyediakan informasi strategis bagi Kepala Dinas dan pimpinan terkait. |
|  |  |  | 6. | Menjadi fondasi bagi pengembangan analisis, integrasi aplikasi, dan pemantauan tindak lanjut pada tahap berikutnya. |
| 3 | Tujuan | : | 1. | Mengintegrasikan data kesehatan dari Puskesmas, bidang, program, dan unit kerja terkait |
|  |  |  | 2. | Menstandarkan format, definisi indikator, periode pelaporan, dan sumber data |
|  |  |  | 3. | Memastikan data diperiksa dan divalidasi sebelum digunakan |
|  |  |  | 4. | Menyediakan dashboard capaian indikator dan analisis kinerja kesehatan |
|  |  |  | 5. | Mendukung perencanaan program, penganggaran, monitoring, evaluasi, serta pengambilan keputusan pimpinan |
|  |  |  | 6. | Menyediakan riwayat perubahan data dan pengendalian hak akses pengguna |
| 4 | Ruang Lingkup | : | 1. | Pengelolaan Data |
|  |  |  |  | Pendataan sumber, kamus data, standar indikator, periode pelaporan, dan pengelolaan data induk |
|  |  |  | 2. | Pelaporan dan validasi |
|  |  |  |  | Pengumpulan, pemeriksaan kelengkapan, validasi, pengembalian data yang salah, dan persetujuan data. |
|  |  |  | 3. | Pengolahan dan analisis |
|  |  |  |  | Rekapitulasi, perhitungan indikator, analisis tren, perbandingan target dan realisasi, serta identifikasi kesenjangan. |
|  |  |  | 4. | Penyajian informasi |
|  |  |  |  | Dashboard pimpinan, laporan periodik, grafik capaian, peta sebaran bila tersedia, dan ekspor laporan. |
|  |  |  | 5. | Tindak lanjut dan evaluasi |
|  |  |  |  | Identifikasi masalah, rekomendasi tindak lanjut, penetapan penanggung jawab, dan pemantauan penyelesaian. |
| 5 | Flow Chart | : | 1. | Pengumpulan Data, Input dan Pengiriman Data |
|  |  |  |  | User disini adalah Petugas Puskesmas, Petugas Dinas yang menjadi penanggung jawab program/indikator. |
|  |  |  | 2. | Pemeriksaan dan Verifikasi |
|  |  |  |  | User disini adalah petugas Dinas |
|  |  |  |  | Jika data belum sesuai maka dikembalikan kepada user nomor 1, jika data sesuai maka disetujui sesuai kewenangan |
|  |  |  | 3. | Validasi |
|  |  |  |  | User disini adalah Kepala Bidang/Sekretaris Dinas |
|  |  |  |  | Jika data masih perlu di cek, dikembalikan ke user nomor 2, jika dianggap valid maka disetujui menjadi data final yang tidak dapat dirubah kecuali melalui mekanisme yang ditentukan. |
|  |  |  |  | Data final : perubahan memerlukan kewenangan khusus dan jejak audit.<br>Status tersebut sebaiknya dapat difilter menurut fasilitas kesehatan, bidang, program, indikator, dan periode pelaporan. Sistem juga perlu mencatat siapa yang mengirim, memeriksa, menyetujui, atau mengubah data beserta waktu tindakannya. |
|  |  |  | 4. | Integrasi dan Pengolahan |
|  |  |  |  | Penyimpanan terpusat, rekapitulasi, perhitungan indikator, dan analisis |
|  |  |  |  | Pada tahap ini akan muncul hasil pengolahan data berupa tren capaian indikator dan permintaan faktor pendukung dan faktor penghambat capaian indikator yang perlu dilengkapi oleh kepala bidang/sekretaris dinas. |
|  |  |  | 5. | Dashboard dan pelaporan |
|  |  |  |  | Informasi capaian, tren, target, kesenjangan, dan laporan periodik |
|  |  |  |  | bagian pelaporan akan menampilkan data capaian yang bisa di filter sesuai pengelompokan berikut : |
|  |  |  |  | 1. Capaian Kinerja Indikator Kinerja Utama (IKU) (Indikator Tujuan dan Sasaran) |
|  |  |  |  | 2. Capaian Kinerja Indikator Kinerja Kunci (IKK) |
|  |  |  |  | 3. Capaian Kinerja Standar Pelayanan Minimal (SPM) |
|  |  |  |  | 4. Capaian Kinerja Indikator Kinerja RIBK |
|  |  |  |  | 5. Capaian Kinerja Indikator Kinerja Program |
|  |  |  |  | 6. Capaian Kinerja Indikator Kinerja Kegiatan |
|  |  |  |  | 7. Capaian Kinerja Indikator Kinerja Sub Kegiatan |
|  |  |  |  | 8. Capaian Kinerja Indikator Kinerja Bidang/Sekretariat/Tim Kerja |
|  |  |  |  | 9. ........................ |
|  |  |  | 6. | Pengambilan Keputusan |
|  |  |  |  | Arahan pimpinan, rencana tindak lanjut, monitoring, dan evaluasi |
|  |  |  |  | Arahan Pimpinan ditujukan kepada Kepala Bidang/Sekretaris Dinas |
