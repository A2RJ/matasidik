# Daftar Indikator Kinerja

> Sumber: `Proses Bisnis Mata Sidik.xlsx` — Sheet **Daftar Indikator Kinerja** (A1:H62)
>
> Dokumen ini adalah ekstraksi 2 dari 3 sheet. Seluruh tabel dan kolom dipertahankan apa adanya.
> Kolom kosong (mis. Definisi Operasional yang masih kosong di file sumber) tetap ditampilkan.
>
> **Catatan implementasi:** nilai seperti “Bidang Kesmas” dan “Timker Kesga” pada kolom **Label Indikator** adalah dimensi kepemilikan/PJ, bukan label taksonomi. Di aplikasi, keduanya dipetakan ke kolom **PJ Data** (`indikator.pj_bidang_id` dan `indikator.pj_timker_id`); label/IKU hanya menyimpan kategori seperti IKU, IKK, SPM, RIBK, Program, Kegiatan, dan Sub Kegiatan.

## Tabel Indikator

Keterangan kolom PJ Data di file sumber terbagi dua sub-kolom:

- **Bidang / Sekretariat**
- **Timker / Sub Bagian**

| No. | Nama Indikator Kinerja | Satuan | Definisi Operasional | Label Indikator | Sifat Data | Bidang / Sekretariat | Timker / Sub Bagian |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | Usia Harapan Hidup (UHH) | (tahun) |  | IKU; Indikator Sasaran; IKK; RIBK; Bidang Kesmas | Inputan Langsung UHH saja, hanya ada nilai Kabupaten. Input 1 kali setahun. | Bidang Kesmas | - |
| 2 | Jumlah Kematian Ibu | Kasus |  | IKU; Indikator Sasaran; IKK; RIBK; Bidang Kesmas; Timker Kesga | Langsung Jumlah Kasus kematian Ibu, tapi di breakdown inputan data yang bisa di filter berdasarkan Desa, Puskesmas, Kecamatan, baru angka Kabupaten.<br>Inputan dimulai oleh Petugas Puskesmas, diverifikasi oleh petugas kabupaten, divalidasi oleh kepala bidang. | Bidang Kesmas | Timker Kesga |
| 3 | Jumlah Kematian Balita | Kasus |  | IKU; Indikator Sasaran; IKK; RIBK; Bidang Kesmas; Timker Kesga | Langsung Jumlah Kasus kematian Ibu, tapi di breakdown inputan data yang bisa di filter berdasarkan Desa, Puskesmas, Kecamatan, baru angka Kabupaten.<br>Inputan dimulai oleh Petugas Puskesmas, diverifikasi oleh petugas kabupaten, divalidasi oleh kepala bidang. | Bidang Kesmas | Timker Kesga |
| 4 | Prevalensi Stunting | % |  | IKU; Indikator Sasaran; IKK; RIBK; Bidang Kesmas; Timker Gizi |  |  |  |
| 5 | Cakupan penerima pemeriksaan kesehatan gratis | % |  | IKU; Indikator Sasaran; IKK; RIBK; Bidang ...... |  |  |  |
| 6 | Cakupan imunisasi bayi lengkap | Persentase |  |  |  |  |  |
| 7 | Cakupan kepesertaan jaminan kesehatan nasional (JKN) | (%) |  |  |  |  |  |
| 8 | Cakupan kepesertaan aktif Jaminan Kesehatan Nasional (JKN) | % |  |  |  |  |  |
| 9 | Angka keberhasilan pengobatan tuberkulosis (treatment success rate) | % |  |  |  |  |  |
| 10 | Persentase pelayanan kesehatan penderita hipertensi sesuai standar | % |  |  |  |  |  |
| 11 | Proporsi penduduk dengan aktivitas fisik cukup | % |  |  |  |  |  |
| 12 | Nilai AKIP perangkat Daerah | Poin |  |  |  |  |  |
|  | Persentase pelayanan kesehatan ibu hamil sesuai standar | % |  |  |  |  |  |
|  | Persentase pelayanan kesehatan ibu hamil sesuai standar | % |  |  |  |  |  |
|  | Persentase pelayanan kesehatan ibu hamil sesuai standar | % |  |  |  |  |  |
|  | Persentase pelayanan kesehatan balita sesuai standar | % |  |  |  |  |  |
|  | Persentase Pelayanan kesehatan pada usia pendidikan dasar sesuai standar | % |  |  |  |  |  |
|  | Persentase Pelayanan kesehatan pada usia produktif sesuai standar | % |  |  |  |  |  |
|  | Persentase warga negara usia 60 tahun ke atas mendapatkan pelayanan kesehatan sesuai standar | % |  |  |  |  |  |
|  | Persentase pelayanan kesehatan penderita hipertensi sesuai standar | % |  |  |  |  |  |
|  | Persentase pelayanan kesehatan penderita diabetes melitus sesuai standar | % |  |  |  |  |  |
|  | Persentase pelayanan kesehatan orang dengan gangguan jiwa berat sesuai standar | % |  |  |  |  |  |
|  | Persentase Pelayanan kesehatan orang terduga tuberkulosis sesuai standar | % |  |  |  |  |  |
|  | Persentase Pelayanan kesehatan orang dengan risiko terinfeksi virus Yang melemahkan daya tahan tubuh manusia (Human Immunodeficiency Virus) sesuai standar | % |  |  |  |  |  |
|  | Realisai Kinerja Perangkat Daerah (%) |  |  |  |  |  |  |
|  | Realisasi Keuangan Perangkat Daerah (%) |  |  |  |  |  |  |
|  | Indeks Profesionalitas ASN Perangkat Daerah (Poin) |  |  |  |  |  |  |
|  | Prevalensi depresi di umur lebih dari 15 Tahun (%) |  |  |  |  |  |  |
|  | Persentase Desa/Kelurahan Sanitasi Total Berbasis Masyarakat (%) |  |  |  |  |  |  |
|  | Prevalensi Obesitas lebih dari 18 tahun (%) |  |  |  |  |  |  |
|  | Persentase Fasilitas Pelayanan Kesehatan sesuai standar (%) |  |  |  |  |  |  |
|  | Tingkat kepuasan pasien terhadap layanan kesehatan (%) |  |  |  |  |  |  |
|  | Persentase tempat pengelolaan pangan memenuhi syarat (%) |  |  |  |  |  |  |
|  | Persentase fasilitas kesehatan yang terintegrasi dalam Sistem Informasi Kesehatan Nasional (%) |  |  |  |  |  |  |
|  | Persentase fasilitas pelayanan kesehatan terakreditasi paripurna (%) |  |  |  |  |  |  |
|  | Persentase lanjut usia yang mandiri (%) |  |  |  |  |  |  |
|  | Persentase merokok penduduk 10-21 tahun (%) |  |  |  |  |  |  |
|  | Cakupan Penemuan Kasus TB (%) |  |  |  |  |  |  |
|  | Cakupan Layanan Kartu Sumbawa Barat Maju Layanan Maju Kesehatan (%) |  |  |  |  |  |  |
|  | Rasio Tenaga Kesehatan dan Tenaga Medis terhadap populasi (Rasio) |  |  |  |  |  |  |
|  | Persentase Puskesmas dengan SDMK Sesuai Standar (%) |  |  |  |  |  |  |
|  | Persentase fasilitas pelayanan kesehatan dengan perbekalan kesehatan sesuai standar (%) |  |  |  |  |  |  |
|  | Persentase Penerbitan dan Pemenuhan Komitmen SPPIRT Sesuai Standar (%) |  |  |  |  |  |  |
|  | Persentase Pangan Industri Rumah Tangga Aman dan Bermutu (%) |  |  |  |  |  |  |
|  | Persentase Sarana IRTP yang Memenuhi Ketentuan (%) |  |  |  |  |  |  |
|  | Persentase Apotek dan Toko Obat yang Mampu Memelihara Persyaratan Perizinan (%) |  |  |  |  |  |  |
|  | Persentase Rekomendasi Hasil Pengawasan Sediaan Farmasi dan Pangan Olahan dari Lintas Sektor yang Ditindaklanjuti oleh Pemerintah Daerah (%) |  |  |  |  |  |  |
|  | Persentase tingkat Rumah Tangga berPrilaku Hidup Bersih dan Sehat (PHBS) (%) |  |  |  |  |  |  |
|  | Jumlah Dokumen DPA-SKPD dan Laporan Hasil Koordinasi Penyusunan Dokumen DPA-SKPD (Dokumen) |  |  |  |  |  |  |
|  | Jumlah Dokumen Hasil Penyelenggaraan Walidata Pendukung Statistik Sektoral Daerah (Dokumen) |  |  |  |  |  |  |
|  | Jumlah Laporan Evaluasi Kinerja Perangkat Daerah (Laporan) |  |  |  |  |  |  |
|  | Jumlah Laporan Capaian Kinerja dan Ikhtisar Realisasi Kinerja SKPD dan Laporan Hasil Koordinasi Penyusunan Laporan Capaian Kinerja dan Ikhtisar Realisasi Kinerja SKPD (Laporan) |  |  |  |  |  |  |
|  | Jumlah Dokumen Perencanaan Perangkat Daerah (Dokumen) |  |  |  |  |  |  |
|  | Jumlah Dokumen Perubahan DPA-SKPD dan Laporan Hasil Koordinasi Penyusunan Dokumen Perubahan DPA-SKPD (Dokumen) |  |  |  |  |  |  |
|  | Jumlah Dokumen Perubahan RKA-SKPD dan Laporan Hasil Koordinasi Penyusunan Dokumen Perubahan RKA-SKPD (Dokumen) |  |  |  |  |  |  |
|  | Jumlah Berita Acara Hasil Forum Perangkat Daerah Berdasarkan Bidang Urusan yang Diampu dalam Rangka Penyusunan Dokumen Perencanaan Perangkat Daerah (Berita Acara) |  |  |  |  |  |  |
|  | Jumlah Dokumen RKA-SKPD dan Laporan Hasil Koordinasi Penyusunan Dokumen RKA-SKPD (Dokumen) |  |  |  |  |  |  |
|  | dst belum selesai |  |  |  |  |  |  |

> Baris terakhir "dst belum selesai" menunjukkan daftar masih berlanjut di file sumber.
> Total baris data pada tabel di atas: 58 baris (No. 1–12 terisi nomor, sisanya tanpa nomor + 1 baris penutup).
