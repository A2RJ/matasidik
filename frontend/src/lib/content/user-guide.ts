import type { Role } from '$lib/types';

export type GuideTabID = 'ringkasan' | 'puskesmas' | 'approver' | 'verifikator' | 'kabid' | 'pimpinan';

export interface GuideTab {
  id: GuideTabID;
  label: string;
  short: string;
}

export interface HighlightBox {
  left: number;
  top: number;
  width: number;
  height: number;
  label?: string;
}

export interface GuideStep {
  title: string;
  body: string;
  where?: string;
}

export interface GuideRelation {
  actor: string;
  relation: string;
}

export interface GuideScreenshot {
  src: string;
  alt: string;
  caption: string;
  detail: string;
  highlight: string;
  inline?: 'scroll' | 'full';
  access?: string;
  stage?: string;
  size?: 'compact' | 'wide';
}

export interface GuideAccess {
  menu: string;
  access: string;
  detail: string;
}

export interface GuideFeatureLink {
  feature: string;
  relation: string;
}

export interface GuideWorkflowDetail {
  stage: string;
  actions: string;
  data: string;
  next: string;
}

export interface MenuCatalogRow {
  menu: string;
  route: string;
  function: string;
  data: string;
  related: string;
}

export interface MenuAccessRow {
  menu: string;
  pkm: string;
  verifier: string;
  kabid: string;
  pimpinan: string;
}

export interface DataLinkRow {
  from: string;
  to: string;
  via: string;
  note: string;
}

export interface RoleGuide {
  id: Exclude<GuideTabID, 'ringkasan'>;
  title: string;
  level: string;
  roles: Role[];
  summary: string;
  startLabel: string;
  accounts: string[];
  flowOverlay: HighlightBox;
  steps: GuideStep[];
  workflowDetail: GuideWorkflowDetail[];
  access: GuideAccess[];
  featureLinks: GuideFeatureLink[];
  relations: GuideRelation[];
  screenshots: GuideScreenshot[];
  notes: string[];
}

export const GUIDE_TABS: GuideTab[] = [
  { id: 'ringkasan', label: 'Ringkasan & peta peran', short: 'Ringkasan' },
  { id: 'puskesmas', label: 'Puskesmas / PJ Program', short: 'Puskesmas' },
  { id: 'approver', label: 'Kepala Puskesmas', short: 'Kepala PKM' },
  { id: 'verifikator', label: 'Verifikator Dinas', short: 'Verifikator' },
  { id: 'kabid', label: 'Kabid / Sekretaris', short: 'Kabid' },
  { id: 'pimpinan', label: 'Pimpinan', short: 'Pimpinan' }
];

export const DEMO_ACCOUNTS = [
  { username: 'siti', role: 'Puskesmas Operator · PKM Taliwang', use: 'Upload, pratinjau, kirim' },
  { username: 'kepalapkm', role: 'Kepala Puskesmas · PKM Taliwang', use: 'Persetujuan BNBA & pengajuan IKU' },
  { username: 'pjgizi', role: 'PJ Program · Timker Gizi', use: 'Alur input dengan cakupan faskes pada akun' },
  { username: 'rina', role: 'Verifikator Dinas', use: 'Verifikasi, setujui/kembalikan' },
  { username: 'kabid', role: 'Kabid Kesehatan Masyarakat', use: 'Validasi final, faktor' },
  { username: 'sekretaris', role: 'Sekretaris Dinas', use: 'Validasi final, faktor' },
  { username: 'pimpinan', role: 'Kepala Dinas', use: 'Dashboard, arahan, TL, jejak audit' }
];

export const FLOW_LEGEND = [
  { swatch: 'blue', text: 'Alur utama: masuk → unggah → verifikasi → validasi → hitung → dashboard' },
  { swatch: 'dashed', text: 'Garis putus-putus: pengembalian data atau dukungan lintas peran' }
];

const flowImage = '/user-guide/alur-panduan.webp';
const flowAlt =
  'Alur lengkap MATA SIDIK: masuk dan sesi, unggah BNBA Puskesmas, pratinjau dan koreksi IKU, verifikasi Dinas, validasi final Kabid/Sekretaris, perhitungan dan dashboard, serta arahan dan TL Pimpinan.';

export const FLOW_IMAGE = { src: flowImage, alt: flowAlt };

export const MENU_CATALOG: MenuCatalogRow[] = [
  {
    menu: 'Dashboard',
    route: '/',
    function:
      'Ringkasan IKU induk, jumlah sesuai/di bawah target, filter periode–kelompok–faskes–status, grafik capaian, tabel SPM, laporan absolut bulanan, dan ekspor.',
    data:
      'Nilai final IKU, capaian, target, capaian tahunan, dan riwayat perhitungan.',
    related:
      'Mengambil hasil akhir dari Upload → Submission → Persetujuan → Perhitungan. Setiap baris dapat dibuka ke Detail SPM.'
  },
  {
    menu: 'Detail SPM',
    route: '/dashboard/spm/{kode}',
    function:
      'Menampilkan rincian satu baris: capaian, target, gap, jumlah/total, cakupan upload per puskesmas, grafik, berkas sumber, dan komponen angka.',
    data:
      'Capaian, target, jumlah/total, berkas sumber, pengunggah, dan komponen pembentuk angka.',
    related:
      'Subhalaman Dashboard. Menghubungkan angka dashboard ke Submission, berkas asli, dan jejak audit.'
  },
  {
    menu: 'Upload BNBA',
    route: '/upload',
    function:
      'Memilih berkas .csv/.xlsx, periode, dan faskes; memantau proses; menampilkan pratinjau baris; mengoreksi IKU rincian; mengirim ke alur persetujuan.',
    data:
      'Berkas BNBA, baris data, hasil klasifikasi, status proses, dan catatan audit.',
    related:
      'Bergantung pada Master Data. Setelah dikirim, data masuk ke Submission dan alur persetujuan. Baris bermasalah harus diperbaiki sebelum dikirim.'
  },
  {
    menu: 'Submission',
    route: '/submissions',
    function:
      'Daftar pengajuan per periode; detail baris, klasifikasi, catatan, riwayat revisi, berkas asli; tindakan kirim/setujui/kembalikan sesuai status dan peran.',
    data:
      'Status pengajuan, baris data, catatan keputusan, riwayat revisi, dan berkas asli.',
    related:
      'Pusat alur persetujuan. Menerima hasil Upload, menjadi sumber nilai setelah final, dan dapat ditelusuri dari Dashboard.'
  },
  {
    menu: 'Pengajuan IKU',
    route: '/iku-requests',
    function:
      'Operator Puskesmas mengajukan tambah, ubah, atau nonaktifkan IKU; Puskesmas memeriksa; Admin mempublikasikan.',
    data:
      'Usulan IKU, riwayat revisi, catatan pemeriksaan, dan status pengajuan.',
    related:
      'Terhubung ke Master Data. IKU baru berlaku setelah disetujui Puskesmas dan Admin.'
  },
  {
    menu: 'Master Data',
    route: '/master',
    function:
      'Mengelola indikator, IKU, kategori dashboard, bidang/timker, target, capaian tahunan, periode, dan pratinjau pemetaan.',
    data:
      'Kode, definisi, satuan, kategori, penanggung jawab, target, dan capaian tahunan.',
    related:
      'Menjadi acuan Upload (kode dan periode), Perhitungan (komponen), dan Dashboard (kategori/target).'
  },
  {
    menu: 'Faktor & TL',
    route: '/factors',
    function:
      'Mengisi faktor pendukung/penghambat dan mengelola tindak lanjut: masalah → rekomendasi → PJ → status.',
    data:
      'Nilai IKU, faktor pendukung/penghambat, dan tindak lanjut.',
    related:
      'Memakai nilai final dari Dashboard. Faktor tidak mengubah angka; TL menjadi jembatan keputusan Pimpinan ke eksekusi Kabid/Sekretaris.'
  },
  {
    menu: 'Jejak Audit',
    route: '/audit',
    function:
      'Menampilkan riwayat keputusan: waktu, aktor, tindakan, bagian, dan alasan/catatan.',
    data:
      'Riwayat keputusan dari upload, koreksi IKU, pengiriman, persetujuan/pengembalian, faktor, TL, dan perubahan master.',
    related:
      'Lintas semua fitur. Menjadi bukti keputusan dan penghubung antara aktor, data, dan waktu.'
  },
  {
    menu: 'Administrasi',
    route: '/admin',
    function:
      'Mengelola pengguna, faskes, dan jabatan struktural.',
    data:
      'Akun, peran, jabatan, faskes, dan bidang.',
    related:
      'Menentukan peran dan cakupan akses yang dipakai di seluruh menu.'
  },
  {
    menu: 'Profil',
    route: '/profile',
    function:
      'Mengubah nama tampilan dan password sendiri; menampilkan peran, jabatan, faskes, bidang, dan login terakhir.',
    data:
      'Nama tampilan, password, peran, jabatan, faskes, dan bidang.',
    related:
      'Terhubung ke autentikasi dan seluruh menu yang memakai identitas/peran pengguna.'
  },
  {
    menu: 'Panduan',
    route: '/user-guide',
    function:
      'Dokumentasi in-app per peran: alur, langkah, hak akses menu, keterkaitan fitur/data, status, screenshot, dan batasan pilot.',
    data:
      'Tidak menulis data bisnis. Membaca konten panduan dan aset screenshot.',
    related:
      'Menjadi titik rujukan untuk semua menu dan peran. Tidak menggantikan kontrol akses aplikasi.'
  }
];

export const MENU_ACCESS: MenuAccessRow[] = [
  {
    menu: 'Dashboard + Detail SPM',
    pkm: 'Milik faskes sendiri',
    verifier: 'Bidang*',
    kabid: 'Bidang*',
    pimpinan: 'Semua',
  },
  {
    menu: 'Upload BNBA',
    pkm: 'Ya',
    verifier: 'Tidak',
    kabid: 'Tidak',
    pimpinan: 'Tidak',
  },
  {
    menu: 'Submission',
    pkm: 'Milik faskes',
    verifier: 'Verifikasi',
    kabid: 'Validasi',
    pimpinan: 'Lihat',
  },
  {
    menu: 'Pengajuan IKU',
    pkm: 'Operator: ajukan',
    verifier: 'Tidak',
    kabid: 'Tidak',
    pimpinan: 'Tidak',
  },
  {
    menu: 'Master Data',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Lihat + target*',
    pimpinan: 'Tidak',
  },
  {
    menu: 'Faktor & TL',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Faktor + status TL',
    pimpinan: 'TL + status',
  },
  {
    menu: 'Jejak Audit',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Tidak',
    pimpinan: 'Ya',
  },
  {
    menu: 'Administrasi',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Tidak',
    pimpinan: 'Tidak',
  },
  {
    menu: 'Profil',
    pkm: 'Ya',
    verifier: 'Ya',
    kabid: 'Ya',
    pimpinan: 'Ya',
  },
  {
    menu: 'Panduan',
    pkm: 'Publik',
    verifier: 'Publik',
    kabid: 'Publik',
    pimpinan: 'Publik',
  }
];

export const DATA_FLOW_LINKS: DataLinkRow[] = [
  {
    from: 'Master Data',
    to: 'Upload BNBA',
    via: 'kode IKU · kode indikator · periode · versi master',
    note:
      'Baris hanya sah jika kode dan relasinya ada di master. Master berubah lewat versi baru, bukan penghapusan data lama.'
  },
  {
    from: 'Upload BNBA',
    to: 'Submission',
    via: 'berkas BNBA · baris data · status proses',
    note:
      'Berkas diperiksa di latar belakang. Pratinjau menampilkan baris, klasifikasi, dan catatan yang perlu diperbaiki.'
  },
  {
    from: 'Submission',
    to: 'Persetujuan',
    via: 'status · catatan keputusan · jejak audit',
    note:
      'Draf → kirim → persetujuan Puskesmas → verifikasi Dinas → validasi → final. Pengembalian wajib catatan; final mengunci data.'
  },
  {
    from: 'Final',
    to: 'Perhitungan',
    via: 'nilai IKU · rincian komponen',
    note:
      'Final memicu perhitungan. IKU induk dihitung dari komponen yang eksplisit, bukan rata-rata persentase.'
  },
  {
    from: 'Perhitungan',
    to: 'Dashboard',
    via: 'nilai final · capaian · sumber angka',
    note:
      'Dashboard membaca nilai final. Baris yang belum memiliki nilai final berstatus belum ada data, bukan nol.'
  },
  {
    from: 'Dashboard',
    to: 'Faktor & TL',
    via: 'nilai IKU · faktor · tindak lanjut',
    note:
      'Faktor adalah analisis kualitatif; TL adalah rencana aksi dengan PJ. Keduanya tidak mengubah perhitungan.'
  },
  {
    from: 'Semua keputusan',
    to: 'Jejak Audit',
    via: 'waktu · aktor · tindakan · catatan',
    note:
      'Jejak audit menjadi bukti siapa melakukan apa, kapan, dan mengapa.'
  }
];

export const ROLE_GUIDES: Record<Exclude<GuideTabID, 'ringkasan'>, RoleGuide> = {
  puskesmas: {
    id: 'puskesmas',
    title: 'Puskesmas & PJ Program',
    level: 'Level 1 · Input dan pengiriman data',
    roles: ['pkm_operator', 'pkm_pj'],
    summary:
      'Menyiapkan berkas BNBA, mengunggah ke MATA SIDIK, memeriksa pratinjau, mengoreksi IKU rincian, lalu mengirim pengajuan untuk disetujui Kepala Puskesmas sebelum verifikasi Dinas.',
    startLabel: 'Input BNBA',
    accounts: ['siti · Puskesmas Operator · PKM Taliwang', 'pjgizi · PJ Program · cakupan faskes pada akun'],
    flowOverlay: { left: 1.19, top: 20.45, width: 19.05, height: 75 },
    steps: [
      {
        title: 'Masuk dan pastikan periode',
        where: '/login → /',
        body:
          'Masuk dengan akun Puskesmas/PJ. Setelah login, dashboard menampilkan periode aktif. Pilih periode yang sama dengan periode data BNBA sebelum membuka menu Upload BNBA.'
      },
      {
        title: 'Siapkan berkas BNBA',
        where: 'Template BNBA di halaman Upload BNBA',
        body:
          'Gunakan berkas .csv atau .xlsx. Kolom wajib: periode, kode_faskes, kode_iku, kode_indikator, scope_type, scope_code, numerator, denominator. Kolom nilai dan catatan boleh diisi; kode_kategori opsional. Jangan memakai kode IKU induk untuk baris unggahan.'
      },
      {
        title: 'Unggah berkas',
        where: '/upload',
        body:
          'Pilih berkas BNBA, periode, dan faskes pengunggah. Akun Puskesmas dengan faskes otomatis memakai faskes sendiri; pastikan faskes sesuai dengan akun. Klik Upload & proses.'
      },
      {
        title: 'Tunggu proses berkas',
        where: 'Kartu Proses berkas',
        body:
          'Sistem menyimpan berkas dan memprosesnya; status diperbarui otomatis. Langkah yang tampil: Validasi berkas → Klasifikasi → Simpan baris → Selesai. Berkas tidak perlu diunggah ulang jika halaman di-refresh; buka kembali dari daftar pengajuan.'
      },
      {
        title: 'Periksa pratinjau dan koreksi IKU',
        where: 'Kartu Langkah 3 · Pratinjau',
        body:
          'Baris yang sesuai berwarna hijau. Baris yang tidak sesuai, memakai IKU induk, perlu dipilih ulang, duplikat, nilainya tidak valid, atau di luar kewenangan menjadi baris bermasalah. Gunakan pilihan Ubah IKU hanya untuk memilih IKU rincian yang terhubung ke indikator baris tersebut. Perubahan sebelum dikirim tercatat pada jejak audit.'
      },
      {
        title: 'Kirim untuk persetujuan',
        where: 'Tombol Kirim untuk persetujuan',
        body:
          'Tombol Kirim aktif hanya jika tidak ada baris bermasalah. Setelah dikirim, status berubah menjadi Dikirim dan pengajuan masuk ke persetujuan Puskesmas sebelum verifikasi Dinas. Data tidak bisa diubah langsung sampai ada keputusan pengembalian.'
      },
      {
        title: 'Tindak lanjuti hasil verifikasi',
        where: '/submissions → detail submission',
        body:
          'Jika verifikator mengembalikan, status menjadi Dikembalikan ke Petugas. Baca catatan, perbaiki baris yang diminta (jika masih diizinkan), lalu kirim ulang. Jika permintaan berasal dari Kabid/Sekretaris, statusnya Dikembalikan ke Verifikator dan verifikator yang memproses ulang.'
      }
    ],
    workflowDetail: [
      {
        stage: 'Persiapan akun & periode',
        actions:
          'Pastikan akun aktif, memiliki faskes/cakupan yang benar, dan periode sama dengan periode berkas. Periode terkunci tidak menerima unggahan.',
        data: 'Akun · faskes · periode',
        next: 'Siap unggah'
      },
      {
        stage: 'Penyiapan berkas BNBA',
        actions:
          'Gunakan template resmi .csv/.xlsx; isi kode_iku rincian, kode_indikator, periode, numerator, dan denominator. Jangan mengubah header atau memakai IKU induk.',
        data: 'Berkas BNBA · kode kontrak',
        next: 'Siap diunggah'
      },
      {
        stage: 'Unggah & proses latar',
        actions:
          'Pilih berkas, periode, dan faskes; klik Upload & proses. Sistem menyimpan berkas dan memantau validasi → klasifikasi → simpan baris.',
        data: 'pengajuan · berkas · JOB',
        next: 'Pratinjau'
      },
      {
        stage: 'Pratinjau & koreksi IKU',
        actions:
          'Periksa setiap baris, jumlah/total, dan status klasifikasi. Baris bermasalah harus nol. Gunakan Ubah IKU hanya ke rincian yang valid, lalu Simpan per baris.',
        data: 'baris pengajuan · nilai indikator · audit',
        next: 'Kirim'
      },
      {
        stage: 'Kirim ke verifikasi',
        actions:
          'Klik Kirim. Status menjadi Dikirim, data terkunci untuk pengirim, dan pengajuan masuk ke persetujuan Puskesmas. Setelah disetujui, pengajuan diteruskan ke verifikator Dinas.',
        data: 'Status Dikirim · catatan pemeriksaan',
        next: 'Verifikator'
      },
      {
        stage: 'Tindak lanjut pengembalian',
        actions:
          'Baca catatan, perbaiki baris yang diizinkan, lalu kirim ulang. Jangan mengedit baris setelah status Dikirim; perubahan mengikuti pengembalian/revisi.',
        data: 'Status Dikembalikan · revisi · audit',
        next: 'Verifikator'
      }
    ],
    access: [
      {
        menu: 'Dashboard + Detail SPM',
        access: 'Milik faskes sendiri',
        detail: 'Melihat capaian final faskes dan sumber angka baris, tanpa melihat faskes lain.'
      },
      {
        menu: 'Upload BNBA',
        access: 'Ya',
        detail: 'Upload berkas, pantau proses, pratinjau, koreksi IKU rincian, lalu kirim.'
      },
      {
        menu: 'Submission',
        access: 'Milik faskes',
        detail: 'Daftar/detail, riwayat, revisi, dan kirim ulang setelah dikembalikan.'
      },
      {
        menu: 'Profil & Panduan',
        access: 'Ya',
        detail: 'Nama/password sendiri dan dokumentasi in-app.'
      },
      {
        menu: 'Master, Faktor & TL, Jejak Audit',
        access: 'Tidak',
        detail: 'Bukan kewenangan peran input; faktor/TL dan master dipegang peran lain.'
      }
    ],
    featureLinks: [
      {
        feature: 'Master Data → Upload BNBA',
        relation:
          'Master menentukan kode IKU/indikator, satuan, target, periode, dan release yang membuat baris upload sah atau masuk karantina.'
      },
      {
        feature: 'Upload BNBA → Submission',
        relation:
          'Berkas dan proses menjadi pengajuan; pratinjau membaca baris pengajuan dan nilai indikator sementara sebelum dikirim.'
      },
      {
        feature: 'Submission → Persetujuan',
        relation:
          'Kirim mengubah status menjadi Dikirim. Pengembalian dari verifikator membuka perbaikan; pengembalian dari Kabid/Sekretaris diteruskan ke verifikator.'
      },
      {
        feature: 'Dashboard + Detail SPM',
        relation:
          'Setelah final, nilai muncul di dashboard dan dapat ditelusuri kembali ke baris pengajuan dan berkas upload.'
      },
      {
        feature: 'Jejak Audit',
        relation:
          'Upload, koreksi IKU, dan kirim ulang tercatat pada jejak audit, meskipun peran ini tidak membuka menu Jejak Audit.'
      }
    ],
    relations: [
      {
        actor: 'Verifikator Dinas',
        relation:
          'Menerima pengajuan berstatus Dikirim. Verifikator dapat setujui ke Kabid/Sekretaris atau kembalikan ke Puskesmas dengan catatan.'
      },
      {
        actor: 'Kabid / Sekretaris',
        relation:
          'Menerima pengajuan setelah verifikasi. Keputusan final atau kembalikan ke verifikator tidak dikirim langsung ke Puskesmas.'
      },
      {
        actor: 'Sistem',
        relation:
          'Mengklasifikasi baris, menyimpan nilai, dan menghitung IKU setelah data final. Puskesmas tidak mengisi angka dashboard secara manual.'
      },
      {
        actor: 'Pimpinan',
        relation:
          'Melihat output final di dashboard dan memberi arahan/TL. Pimpinan tidak mengubah data mentah Puskesmas.'
      }
    ],
    screenshots: [
      {
        src: '/user-guide/01-upload-form.webp',
        stage: 'Unggah & proses latar',
        size: 'wide',
        access: 'Menu Upload BNBA → pilih periode dan faskes → pilih berkas → klik Upload & proses.',
        alt: 'Form Upload BNBA dengan pilihan berkas, periode, faskes, dan tombol Upload & proses.',
        caption: 'Form upload BNBA',
        highlight:
          'Kotak hijau menandai form yang harus diisi peran Puskesmas/PJ: berkas, periode, faskes, lalu Upload & proses.',
        detail:
          'Screenshot diambil dari halaman /upload pada aplikasi pilot. Seluruh form dan tabel di sekitarnya tetap ditampilkan agar konteks tidak hilang.'
      },
      {
        src: '/user-guide/02-upload-preview.webp',
        stage: 'Pratinjau & koreksi IKU',
        size: 'wide',
        access: 'Menu Upload BNBA → unggah berkas → tunggu proses → bagian Langkah 3 · Pratinjau.',
        alt: 'Pratinjau BNBA dengan baris hasil pembacaan, status klasifikasi, pilihan IKU, dan tombol Kirim untuk persetujuan.',
        caption: 'Pratinjau dan tombol kirim',
        highlight:
          'Kotak hijau menandai tombol Kirim untuk persetujuan. Periksa seluruh baris dan pastikan baris bermasalah nol sebelum menekan tombol.',
        detail:
          'Pratinjau memperlihatkan baris, numerator/denominator, status klasifikasi, IKU, dan catatan. Dropdown Ubah IKU hanya tersedia untuk IKU rincian yang valid.'
      },
      {
        src: '/user-guide/03-submission-list.webp',
        stage: 'Tindak lanjut pengembalian',
        size: 'wide',
        access: 'Menu Submission → pilih periode → tabel Submission per periode → aksi Buka.',
        alt: 'Tabel Submission per periode dengan kolom kode, faskes, pengunggah, periode, status, identitas berkas, diperbarui, dan aksi Buka.',
        caption: 'Daftar pengajuan dan aksi buka',
        highlight:
          'Kotak hijau menandai aksi Buka pada tabel. Klik Buka untuk melihat baris, riwayat, revisi, dan tombol yang tersedia untuk peran Anda.',
        detail:
          'Peran Puskesmas hanya melihat pengajuan dari faskes sendiri. Filter periode di atas tabel membantu menemukan pengajuan yang dikembalikan.'
      }
    ],
    notes: [
      'Pilot belum menyediakan input angka manual; data masuk melalui berkas BNBA .csv atau .xlsx.',
      'Pengajuan BNBA disetujui Kepala Puskesmas faskes yang sama sebelum diteruskan ke Verifikator Dinas. Pengaju tidak dapat menyetujui pengajuannya sendiri.',
      'IKU induk tidak boleh menerima unggahan langsung. Nilainya dihitung dari komponen anak yang eksplisit.',
      'Baris yang ambigu, tidak dikenal, duplikat, atau cakupannya di luar kewenangan masuk karantina/baris bermasalah; sistem tidak menebak IKU.',
      'Format BNBA final dan kebijakan data pribadi belum diputuskan Dinas. Jangan menambahkan kolom identitas pribadi tanpa keputusan.'
    ]
  },
  approver: {
    id: 'approver',
    title: 'Kepala Puskesmas',
    level: 'Persetujuan Puskesmas',
    roles: ['pkm_approver'],
    summary:
      'Memeriksa dan menyetujui atau mengembalikan pengajuan BNBA serta pengajuan IKU dari faskes sendiri sebelum diteruskan ke Dinas atau Admin.',
    startLabel: 'Persetujuan',
    accounts: ['kepalapkm · Kepala PKM Taliwang'],
    flowOverlay: { left: 1.19, top: 20.45, width: 19.05, height: 75 },
    steps: [
      {
        title: 'Buka pengajuan BNBA',
        where: '/submissions',
        body:
          'Masuk dengan akun Kepala Puskesmas. Pilih periode, lalu cari pengajuan berstatus Dikirim pada tahap Persetujuan Puskesmas. Daftar hanya memuat pengajuan dari faskes sendiri.'
      },
      {
        title: 'Periksa isi pengajuan',
        where: 'Detail pengajuan',
        body:
          'Periksa baris data, klasifikasi, jumlah/total, catatan, riwayat, dan berkas asli. Pastikan tidak ada baris yang masih perlu diperbaiki sebelum memberi keputusan.'
      },
      {
        title: 'Setujui atau kembalikan',
        where: 'Panel Aksi → Setujui / Kembalikan',
        body:
          'Setujui untuk meneruskan pengajuan ke Verifikator Dinas. Kembalikan dengan catatan jika data belum lengkap. Pengajuan yang dikembalikan kembali ke Operator Puskesmas untuk diperbaiki.'
      },
      {
        title: 'Periksa pengajuan IKU',
        where: '/iku-requests',
        body:
          'Buka menu Pengajuan IKU untuk melihat usulan tambah, ubah, atau nonaktifkan IKU dari Operator Puskesmas faskes sendiri. Periksa kode, nama, kategori, dan cara hitung.'
      },
      {
        title: 'Putuskan pengajuan IKU',
        where: 'Kolom Keputusan pada daftar pengajuan',
        body:
          'Setujui untuk meneruskan ke Admin, atau kembalikan/tolak dengan catatan. Pengajuan yang disetujui Puskesmas belum aktif sampai Admin mempublikasikan.'
      }
    ],
    workflowDetail: [
      {
        stage: 'Buka antrean Puskesmas',
        actions:
          'Pilih periode dan cari pengajuan berstatus Dikirim pada tahap Persetujuan Puskesmas.',
        data: 'pengajuan · periode · status',
        next: 'Periksa isi'
      },
      {
        stage: 'Periksa kelengkapan',
        actions:
          'Periksa baris, klasifikasi, jumlah/total, catatan, riwayat, dan berkas asli. Jangan menyetujui jika masih ada baris bermasalah.',
        data: 'baris pengajuan · berkas · catatan',
        next: 'Keputusan'
      },
      {
        stage: 'Setujui atau kembalikan',
        actions:
          'Setujui untuk meneruskan ke Verifikator Dinas. Kembalikan dengan catatan untuk memperbaiki data di Puskesmas.',
        data: 'status · catatan · jejak audit',
        next: 'Verifikator / Operator'
      },
      {
        stage: 'Putuskan pengajuan IKU',
        actions:
          'Periksa usulan IKU faskes sendiri. Setujui untuk diteruskan ke Admin, atau kembalikan/tolak dengan catatan.',
        data: 'usulan IKU · catatan · status',
        next: 'Admin / Operator'
      }
    ],
    access: [
      {
        menu: 'Dashboard + Detail SPM',
        access: 'Milik faskes sendiri',
        detail: 'Melihat capaian final faskes dan sumber angka tanpa melihat faskes lain.'
      },
      {
        menu: 'Submission',
        access: 'Persetujuan',
        detail: 'Menyetujui atau mengembalikan pengajuan BNBA faskes sendiri pada tahap Persetujuan Puskesmas.'
      },
      {
        menu: 'Pengajuan IKU',
        access: 'Persetujuan',
        detail: 'Menyetujui atau mengembalikan usulan IKU faskes sendiri; tidak mempublikasikan IKU.'
      },
      {
        menu: 'Profil & Panduan',
        access: 'Ya',
        detail: 'Nama/password sendiri dan dokumentasi in-app.'
      },
      {
        menu: 'Upload, Master, Faktor & TL, Jejak Audit',
        access: 'Tidak',
        detail: 'Bukan kewenangan Kepala Puskesmas pada alur ini.'
      }
    ],
    featureLinks: [
      {
        feature: 'Operator → Kepala Puskesmas',
        relation:
          'Operator Puskesmas mengirim pengajuan BNBA; Kepala Puskesmas memeriksa dan meneruskan atau mengembalikan.'
      },
      {
        feature: 'Kepala Puskesmas → Verifikator Dinas',
        relation:
          'Pengajuan yang disetujui masuk ke antrean Verifikator Dinas. Pengembalian kembali ke Operator Puskesmas.'
      },
      {
        feature: 'Pengajuan IKU → Admin',
        relation:
          'Usulan IKU yang disetujui Kepala Puskesmas masuk ke Admin untuk dipublikasikan.'
      }
    ],
    relations: [
      {
        actor: 'Operator Puskesmas',
        relation:
          'Mengirim pengajuan BNBA dan usulan IKU faskes. Kepala Puskesmas dapat mengembalikan dengan catatan.'
      },
      {
        actor: 'Verifikator Dinas',
        relation:
          'Menerima pengajuan BNBA setelah disetujui Kepala Puskesmas.'
      },
      {
        actor: 'Admin',
        relation:
          'Menerima usulan IKU setelah disetujui Kepala Puskesmas dan mempublikasikannya.'
      }
    ],
    screenshots: [],
    notes: [
      'Kepala Puskesmas hanya dapat memutus pengajuan faskes sendiri.',
      'Pengaju tidak dapat menyetujui pengajuannya sendiri.',
      'Kepala Puskesmas tidak mengunggah BNBA, tidak mengelola master, dan tidak mempublikasikan IKU.',
      'Pengembalian wajib catatan. Tidak ada persetujuan otomatis.'
    ]
  },
  verifikator: {
    id: 'verifikator',
    title: 'Verifikator Dinas',
    level: 'Level 2 · Pemeriksaan dan verifikasi',
    roles: ['dinas_verifier'],
    summary:
      'Memeriksa pengajuan yang sudah disetujui Kepala Puskesmas, menilai kelengkapan baris dan berkas, lalu meneruskan ke Kabid/Sekretaris atau mengembalikan ke Puskesmas dengan catatan.',
    startLabel: 'Verifikasi',
    accounts: ['rina · Verifikator Dinas · Bidang Kesehatan Masyarakat'],
    flowOverlay: { left: 20.83, top: 20.45, width: 17.86, height: 75 },
    steps: [
      {
        title: 'Buka daftar pengajuan',
        where: '/submissions',
        body:
          'Masuk dengan akun Verifikator Dinas, pilih periode, lalu buka tabel Submission per periode. Peran ini melihat pengajuan lintas faskes sesuai kewenangannya, bukan hanya satu faskes.'
      },
      {
        title: 'Pilih pengajuan berstatus Dikirim',
        where: 'Aksi Buka',
        body:
          'Cari baris dengan status Dikirim pada tahap Verifikasi Dinas. Klik Buka atau klik baris untuk masuk ke detail. Pengajuan berstatus lain belum/tidak berada di tahap verifikasi.'
      },
      {
        title: 'Periksa isi dan bukti',
        where: 'Detail pengajuan',
        body:
          'Periksa ringkasan status alur, jumlah baris bermasalah, baris hasil pembacaan, klasifikasi, IKU, numerator/denominator, catatan, riwayat, revisi, identitas berkas, dan tombol Unduh berkas asli. Jangan menyetujui baris yang masih bermasalah.'
      },
      {
        title: 'Setujui atau kembalikan',
        where: 'Panel Aksi → Setujui / Kembalikan',
        body:
          'Setujui mengubah status menjadi Diverifikasi dan meneruskan ke Kabid/Sekretaris. Pengembalian wajib catatan dan mengubah status menjadi Dikembalikan ke Petugas. Pilih salah satu sesuai hasil pemeriksaan.'
      },
      {
        title: 'Tangani pengembalian dari Kabid/Sekretaris',
        where: 'Status Dikembalikan ke Verifikator',
        body:
          'Jika Kabid/Sekretaris meminta pemeriksaan ulang, pengajuan kembali ke tahap verifikasi. Periksa catatan, lakukan verifikasi ulang, lalu setujui atau kembalikan lagi. Peran ini tidak menetapkan final.'
      }
    ],
    workflowDetail: [
      {
        stage: 'Buka daftar & filter',
        actions:
          'Pilih periode, cari pengajuan berstatus Dikirim, lalu buka detail. Pastikan status dan tahap alur sesuai.',
        data: 'pengajuan · periode',
        next: 'Detail pengajuan'
      },
      {
        stage: 'Periksa kelengkapan',
        actions:
          'Periksa baris, klasifikasi, IKU, jumlah/total, baris bermasalah, identitas/berkas, riwayat, revisi, catatan, dan unduhan berkas asli.',
        data: 'baris pengajuan · berkas · pemeriksaan',
        next: 'Keputusan'
      },
      {
        stage: 'Setujui',
        actions:
          'Klik Setujui tanpa menetapkan final. Status menjadi Diverifikasi dan tugas berpindah ke Kabid/Sekretaris.',
        data: 'Status Diverifikasi · jejak audit',
        next: 'Kabid/Sekretaris'
      },
      {
        stage: 'Kembalikan ke Puskesmas',
        actions:
          'Klik Kembalikan, isi catatan wajib, lalu kirim. Status menjadi Dikembalikan ke Petugas dan Puskesmas memperbaiki.',
        data: 'Catatan · status · audit',
        next: 'Puskesmas'
      },
      {
        stage: 'Tangani pengembalian dari Kabid',
        actions:
          'Jika Kabid/Sekretaris mengembalikan, status menjadi Dikembalikan ke Verifikator. Cek ulang lalu setujui/kembalikan lagi.',
        data: 'Status pengembalian · audit',
        next: 'Kabid/Sekretaris'
      }
    ],
    access: [
      {
        menu: 'Dashboard + Detail SPM',
        access: 'Bidang*',
        detail: 'Melihat capaian sesuai kebijakan bidang; pada pilot batas bidang belum sepenuhnya ditegakkan.'
      },
      {
        menu: 'Submission',
        access: 'Verifikasi',
        detail: 'Melihat daftar/detail dan menjalankan Setujui/Kembalikan pada status Dikirim.'
      },
      {
        menu: 'Upload BNBA',
        access: 'Tidak',
        detail: 'Sumber data berasal dari Puskesmas; verifikator tidak mengunggah.'
      },
      {
        menu: 'Master Data',
        access: 'Tidak',
        detail: 'Master di luar kewenangan verifikator; verifikator memakai kode yang berlaku.'
      },
      {
        menu: 'Faktor & TL',
        access: 'Tidak',
        detail: 'Faktor diisi Kabid/Sekretaris; TL ditetapkan Pimpinan dan dieksekusi Kabid/Sekretaris.'
      },
      {
        menu: 'Jejak Audit',
        access: 'Tidak',
        detail: 'Jejak audit ditulis sistem; menu Jejak Audit hanya untuk peran pengawas.'
      },
      {
        menu: 'Profil & Panduan',
        access: 'Ya',
        detail: 'Nama/password sendiri dan dokumentasi in-app.'
      }
    ],
    featureLinks: [
      {
        feature: 'Upload/Puskesmas → Submission',
        relation:
          'Pengajuan berstatus Dikirim pada tahap Verifikasi Dinas adalah input utama verifikasi. Verifikator tidak membuat atau mengedit baris unggahan.'
      },
      {
        feature: 'Submission → Kabid/Sekretaris',
        relation:
          'Persetujuan meneruskan data ke validasi final. Pengembalian mengembalikan ke Puskesmas dengan catatan wajib.'
      },
      {
        feature: 'Submission → Jejak Audit',
        relation:
          'Setiap persetujuan/pengembalian menulis catatan pemeriksaan dan jejak audit berisi aktor, waktu, tindakan, dan catatan.'
      },
      {
        feature: 'Dashboard + Detail SPM',
        relation:
          'Setelah Kabid/Sekretaris final, nilai masuk perhitungan dan dashboard. Verifikator tidak menghitung ulang angka.'
      }
    ],
    relations: [
      {
        actor: 'Puskesmas / PJ Program',
        relation:
          'Mengirim pengajuan. Verifikator dapat mengembalikannya dengan catatan; Puskesmas memperbaiki lalu mengirim ulang.'
      },
      {
        actor: 'Kabid / Sekretaris',
        relation:
          'Menerima hasil verifikasi. Kabid/Sekretaris dapat langsung final atau mengembalikan ke verifikator untuk dicek ulang.'
      },
      {
        actor: 'Pimpinan',
        relation:
          'Tidak menerima tugas verifikasi. Pimpinan membaca dashboard final dan memberi arahan/TL setelah data menjadi final.'
      }
    ],
    screenshots: [
      {
        src: '/user-guide/04-verif-list.webp',
        stage: 'Buka daftar & filter',
        size: 'wide',
        access: 'Menu Submission → pilih periode → tabel Submission per periode → aksi Buka.',
        alt: 'Tabel pengajuan untuk verifikator dengan kolom status dan aksi Buka yang ditandai kotak hijau.',
        caption: 'Daftar pengajuan untuk verifikator',
        highlight:
          'Kotak hijau menandai aksi Buka pada tabel. Dari sini verifikator masuk ke detail dan melihat tombol Setujui/Kembalikan yang sesuai status.',
        detail:
          'Tabel tetap ditampilkan penuh: kode, faskes, pengunggah, periode, status, berkas, waktu diperbarui, dan aksi.'
      },
      {
        src: '/user-guide/05-verif-detail.webp',
        stage: 'Setujui',
        size: 'compact',
        access: 'Menu Submission → buka pengajuan berstatus Dikirim → panel Aksi → Setujui/Kembalikan.',
        alt: 'Panel aksi verifikator dengan tombol Setujui dan Kembalikan yang ditandai kotak hijau.',
        caption: 'Panel keputusan verifikator',
        highlight:
          'Kotak hijau menandai tombol Setujui dan Kembalikan. Pengembalian wajib catatan; Setujui meneruskan ke Kabid/Sekretaris.',
        detail:
          'Panel Aksi adalah titik keputusan verifikator. Persetujuan meneruskan ke validasi; pengembalian mengembalikan ke Puskesmas dengan catatan wajib.'
      },
      {
        src: '/user-guide/05b-verif-rows.webp',
        stage: 'Periksa kelengkapan',
        size: 'wide',
        access: 'Menu Submission → buka pengajuan → bagian Pratinjau baris / Baris hasil pembacaan & klasifikasi.',
        alt: 'Tabel baris hasil pembacaan dan klasifikasi untuk verifikator.',
        caption: 'Periksa baris & klasifikasi',
        highlight:
          'Tabel ini adalah komponen yang harus diperiksa sebelum menekan Setujui/Kembalikan: nilai, jumlah/total, status klasifikasi, IKU, dan catatan.',
        detail:
          'Verifikator memeriksa setiap baris; baris yang sesuai dapat diproses, sedangkan baris bermasalah harus diperbaiki atau ditolak.'
      }
    ],
    notes: [
      'Verifikasi bukan final. Hanya Kabid/Sekretaris yang boleh menetapkan status final.',
      'Pengembalian wajib catatan. Tidak ada persetujuan otomatis, termasuk jika batas waktu layanan terlewati.',
      'Puskesmas tidak bisa mengedit baris setelah status Dikirim; perubahan mengikuti alur pengembalian/revisi.',
      'Pada pilot, pembatasan dashboard per bidang belum sepenuhnya ditegakkan; jangan memakai tampilan pilot sebagai acuan kontrol akses.',
      'Pilot belum menerapkan batas waktu layanan dan notifikasi; batas waktu masih keputusan terbuka Dinas.'
    ]
  },
  kabid: {
    id: 'kabid',
    title: 'Kabid / Sekretaris',
    level: 'Level 3 · Validasi final dan faktor',
    roles: ['kabid', 'sekretaris'],
    summary:
      'Memvalidasi data yang sudah diverifikasi, menetapkan final dengan catatan penerimaan, lalu mengisi faktor pendukung/penghambat. Final mengunci data dan memicu perhitungan.',
    startLabel: 'Validasi final',
    accounts: ['kabid · Kabid Kesehatan Masyarakat', 'sekretaris · Sekretaris Dinas'],
    flowOverlay: { left: 39.29, top: 20.45, width: 17.86, height: 75 },
    steps: [
      {
        title: 'Buka pengajuan berstatus Diverifikasi',
        where: '/submissions',
        body:
          'Masuk dengan akun kabid atau sekretaris. Pilih periode dan cari pengajuan berstatus Diverifikasi; status ini berarti pemeriksaan Dinas sudah selesai dan menunggu validasi final.'
      },
      {
        title: 'Periksa data dan jejak',
        where: 'Detail pengajuan',
        body:
          'Periksa baris, klasifikasi, IKU, nilai, baris bermasalah, berkas asli, riwayat, dan revisi. Pastikan catatan verifikator dan data yang diajukan sudah sesuai sebelum memberi keputusan.'
      },
      {
        title: 'Tetapkan final atau kembalikan',
        where: 'Panel Aksi → Setujui / Kembalikan',
        body:
          'Persetujuan final wajib catatan penerimaan. Setelah final, data terkunci dan sistem memicu perhitungan IKU/SPM. Pengembalian mengembalikan ke verifikator dengan catatan; bukan langsung ke Puskesmas.'
      },
      {
        title: 'Isi faktor pendukung dan penghambat',
        where: '/factors',
        body:
          'Setelah ada nilai final, buka Faktor & TL dan pilih periode. Gunakan aksi Isi/Ubah untuk menulis faktor. Faktor adalah analisis kualitatif: tidak mengubah angka, tidak menggantikan tindak lanjut.'
      },
      {
        title: 'Tindak lanjuti TL',
        where: 'Tab Tindak Lanjut',
        body:
          'TL dibuat Pimpinan. Kabid/Sekretaris menerima penugasan dan dapat mengubah status TL (Terbuka, Dikerjakan, Selesai, Dibatalkan) melalui dropdown pada tabel. TL berisi masalah → rekomendasi → PJ → status.'
      }
    ],
    workflowDetail: [
      {
        stage: 'Buka pengajuan Diverifikasi',
        actions:
          'Pilih periode, cari status Diverifikasi, lalu buka detail. Pastikan verifikasi Dinas sudah selesai.',
        data: 'pengajuan · status',
        next: 'Validasi'
      },
      {
        stage: 'Validasi isi & jejak',
        actions:
          'Periksa baris, klasifikasi, IKU, nilai, baris bermasalah, berkas, riwayat, revisi, dan catatan verifikator sebelum memutuskan.',
        data: 'Baris · berkas · pemeriksaan',
        next: 'Final / Dikembalikan'
      },
      {
        stage: 'Tetapkan final',
        actions:
          'Klik Setujui, isi catatan penerimaan, lalu simpan. Data terkunci dan sistem memicu perhitungan IKU/SPM.',
        data: 'Status Final · nilai IKU',
        next: 'Dashboard'
      },
      {
        stage: 'Kembalikan ke verifikator',
        actions:
          'Klik Kembalikan, isi catatan perbaikan, lalu kirim. Status menjadi Dikembalikan ke Verifikator, bukan ke Puskesmas.',
        data: 'Catatan · status · audit',
        next: 'Verifikator'
      },
      {
        stage: 'Faktor & TL',
        actions:
          'Pilih periode, isi faktor pendukung/penghambat, dan perbarui status TL yang ditugaskan. Faktor tidak mengubah angka.',
        data: 'faktor · tindak lanjut',
        next: 'Pimpinan/pantau'
      }
    ],
    access: [
      {
        menu: 'Dashboard + Detail SPM',
        access: 'Bidang*',
        detail: 'Melihat capaian bidang dan menelusuri sumber angka; batas bidang pilot belum sepenuhnya ditegakkan.'
      },
      {
        menu: 'Submission',
        access: 'Validasi',
        detail: 'Persetujuan final dengan catatan penerimaan atau pengembalian ke verifikator dengan catatan.'
      },
      {
        menu: 'Master Data',
        access: 'Lihat + target*',
        detail: 'Melihat master; target dapat diubah sesuai kewenangan. Penambahan atau perubahan master lain di luar kewenangan peran ini.'
      },
      {
        menu: 'Faktor & TL',
        access: 'Faktor + status TL',
        detail: 'Mengisi faktor pendukung/penghambat dan memperbarui status TL yang ditugaskan.'
      },
      {
        menu: 'Upload BNBA',
        access: 'Tidak',
        detail: 'Input berkas milik Puskesmas/PJ.'
      },
      {
        menu: 'Jejak Audit',
        access: 'Tidak',
        detail: 'Jejak audit hanya dibuka peran pengawas; keputusan tetap tercatat.'
      },
      {
        menu: 'Profil & Panduan',
        access: 'Ya',
        detail: 'Nama/password sendiri dan dokumentasi in-app.'
      }
    ],
    featureLinks: [
      {
        feature: 'Submission (Verifikasi) → Validasi final',
        relation:
          'Data berstatus Diverifikasi masuk ke validasi. Final wajib catatan dan langsung mengunci data.'
      },
      {
        feature: 'Final → Perhitungan → Dashboard',
        relation:
          'Final memicu perhitungan IKU/SPM. Dashboard membaca nilai final dan sumber angka, bukan menghitung ulang saat halaman dibuka.'
      },
      {
        feature: 'Dashboard → Faktor & TL',
        relation:
          'Faktor diisi dari nilai IKU final. Faktor tidak mengubah angka; TL adalah rencana aksi ber-PJ.'
      },
      {
        feature: 'Master Data (target) → Dashboard',
        relation:
          'Target dan kategori dari master menentukan status/capaian yang tampil di dashboard.'
      },
      {
        feature: 'Semua keputusan → Jejak Audit',
        relation:
          'Persetujuan, pengembalian, faktor, dan perubahan TL tercatat pada jejak audit meskipun peran ini tidak membuka menu Jejak Audit.'
      }
    ],
    relations: [
      {
        actor: 'Verifikator Dinas',
        relation:
          'Mengirim data setelah verifikasi. Kabid/Sekretaris dapat mengembalikan ke verifikator jika masih perlu dicek.'
      },
      {
        actor: 'Puskesmas / PJ Program',
        relation:
          'Tidak menerima pengembalian langsung dari Kabid/Sekretaris. Perbaikan dari Puskesmas melewati verifikator.'
      },
      {
        actor: 'Sistem',
        relation:
          'Final memicu perhitungan. Dashboard membaca nilai final dan sumber angka sampai berkas; bukan angka yang diketik ulang.'
      },
      {
        actor: 'Pimpinan',
        relation:
          'Memberi arahan dan menetapkan TL. Kabid/Sekretaris mengeksekusi TL dan mengisi faktor; Pimpinan memantau statusnya.'
      }
    ],
    screenshots: [
      {
        src: '/user-guide/06-kabid-detail.webp',
        stage: 'Tetapkan final',
        size: 'compact',
        access: 'Menu Submission → buka pengajuan berstatus Diverifikasi → panel Aksi → Setujui/Kembalikan.',
        alt: 'Panel aksi Kabid/Sekretaris dengan tombol Setujui dan Kembalikan yang ditandai kotak hijau.',
        caption: 'Panel validasi final Kabid/Sekretaris',
        highlight:
          'Kotak hijau menandai tombol Setujui dan Kembalikan. Persetujuan wajib catatan penerimaan; pengembalian wajib catatan perbaikan.',
        detail:
          'Panel Aksi adalah titik final. Setujui mengunci data dan memicu perhitungan; Kembalikan mengembalikan ke verifikator, bukan ke Puskesmas.'
      },
      {
        src: '/user-guide/06b-kabid-rows.webp',
        stage: 'Validasi isi & jejak',
        size: 'wide',
        access: 'Menu Submission → buka pengajuan → bagian Pratinjau baris / Baris hasil pembacaan & klasifikasi.',
        alt: 'Tabel baris hasil pembacaan dan klasifikasi untuk validasi Kabid/Sekretaris.',
        caption: 'Validasi baris sebelum final',
        highlight:
          'Tabel ini harus diperiksa sebelum final: nilai, jumlah/total, status klasifikasi, IKU, dan catatan. Final tidak boleh dilakukan jika masih ada baris bermasalah.',
        detail:
          'Kabid/Sekretaris memeriksa baris yang sudah diverifikasi. Setelah disetujui final, data terkunci dan perhitungan dijalankan.'
      },
      {
        src: '/user-guide/07-factors-action.webp',
        stage: 'Faktor & TL',
        size: 'wide',
        access: 'Menu Faktor & TL → tab Faktor → tabel IKU → aksi Isi/Ubah.',
        alt: 'Tabel Faktor & TL dengan aksi Isi atau Ubah yang ditandai kotak hijau.',
        caption: 'Isi faktor dan pantau TL',
        highlight:
          'Kotak hijau menandai aksi Isi/Ubah pada tabel faktor. Faktor tidak mengubah perhitungan; TL dikelola pada tab terpisah.',
        detail:
          'Tabel faktor menampilkan IKU, capaian, jumlah/total, faktor pendukung, faktor penghambat, dan aksi. Tab Tindak Lanjut menampilkan masalah, rekomendasi, PJ, status, dan waktu pembaruan.'
      }
    ],
    notes: [
      'Final wajib catatan penerimaan. Setelah final, data tidak boleh diubah langsung dari aplikasi.',
      'Pengembalian dari tahap ini kembali ke verifikator, bukan langsung ke Puskesmas.',
      'Pada pilot, pembatasan dashboard per bidang belum sepenuhnya ditegakkan; jangan memakai tampilan pilot sebagai acuan kontrol akses.',
      'Pilot belum memiliki layar perubahan pasca-final dan persetujuan ganda; hal ini masih menjadi batasan yang harus diselesaikan sebelum produksi.',
      'Ambang warna/status formal dan batas waktu layanan belum diputuskan Dinas.'
    ]
  },
  pimpinan: {
    id: 'pimpinan',
    title: 'Pimpinan',
    level: 'Observer & pengambil keputusan',
    roles: ['pimpinan'],
    summary:
      'Membaca dashboard final, menelusuri angka sampai berkas, melihat faktor, memberi arahan, menetapkan tindak lanjut, dan memantau penyelesaiannya.',
    startLabel: 'Dashboard SPM',
    accounts: ['pimpinan · Kepala Dinas'],
    flowOverlay: { left: 76.19, top: 20.45, width: 22.62, height: 75 },
    steps: [
      {
        title: 'Baca dashboard',
        where: '/',
        body:
          'Pilih periode, kelompok (misalnya SPM), faskes, dan status. Ringkasan menampilkan IKU induk, jumlah sesuai target, di bawah target, dan belum ada data. Tabel menampilkan target, capaian, status, gap, dan aksi Detail.'
      },
      {
        title: 'Telusuri angka',
        where: 'Detail SPM',
        body:
          'Klik baris atau tombol Detail. Halaman detail menampilkan sumber angka, cakupan upload per puskesmas, grafik antar faskes, pengajuan/berkas sumber, dan komponen angka. Rantai yang dijamin: dashboard → nilai IKU → nilai indikator → baris pengajuan → pengajuan → berkas.'
      },
      {
        title: 'Gunakan faktor sebagai konteks',
        where: '/factors',
        body:
          'Buka Faktor & TL untuk membaca faktor pendukung/penghambat yang diisi Kabid/Sekretaris. Faktor adalah analisis kualitatif dan tidak mengubah nilai capaian.'
      },
      {
        title: 'Tetapkan tindak lanjut',
        where: 'Faktor & TL → TL baru',
        body:
          'Buat TL dengan masalah, rekomendasi, PJ, dan status. TL adalah rencana aksi ber-PJ, bukan catatan faktor. Pimpinan dapat memperbarui status TL dan memantau penyelesaiannya.'
      },
      {
        title: 'Periksa jejak keputusan',
        where: '/audit',
        body:
          'Jejak Audit menampilkan waktu, aktor, tindakan, bagian, dan alasan/catatan. Pimpinan dapat memfilter berdasarkan bagian. Jejak audit adalah catatan keputusan, bukan catatan teknis.'
      }
    ],
    workflowDetail: [
      {
        stage: 'Baca dashboard',
        actions:
          'Pilih periode, kelompok, faskes, dan status. Baca ringkasan, grafik, tabel SPM, laporan absolut, dan ekspor.',
        data: 'nilai IKU · baris dashboard · SPM',
        next: 'Detail'
      },
      {
        stage: 'Telusuri detail',
        actions:
          'Klik baris atau Detail. Periksa sumber angka, cakupan upload, grafik faskes, pengajuan/berkas sumber, dan komponen angka.',
        data: 'baris dashboard → … → berkas',
        next: 'Faktor/TL'
      },
      {
        stage: 'Baca faktor',
        actions:
          'Buka Faktor & TL, baca faktor pendukung/penghambat sebagai konteks. Faktor tidak mengubah perhitungan.',
        data: 'faktor · nilai IKU',
        next: 'Arahan'
      },
      {
        stage: 'Tetapkan TL',
        actions:
          'Buat TL dengan masalah → rekomendasi → PJ → status. Pantau status dan koordinasi eksekusi dengan Kabid/Sekretaris.',
        data: 'tindak lanjut · audit',
        next: 'Kabid/Sekretaris'
      },
      {
        stage: 'Jejak Audit',
        actions:
          'Buka Jejak Audit, filter bagian, lalu periksa aktor, waktu, tindakan, dan catatan.',
        data: 'jejak audit',
        next: 'Evaluasi'
      }
    ],
    access: [
      {
        menu: 'Dashboard + Detail SPM',
        access: 'Semua',
        detail: 'Melihat seluruh faskes, kategori, filter, grafik, ekspor, dan sumber angka sampai berkas.'
      },
      {
        menu: 'Submission',
        access: 'Lihat',
        detail: 'Melihat status dan detail untuk konteks; tidak melakukan persetujuan/pengembalian.'
      },
      {
        menu: 'Faktor & TL',
        access: 'TL + status',
        detail: 'Membuat TL (masalah → rekomendasi → PJ → status) dan memperbarui status; faktor dibaca.'
      },
      {
        menu: 'Jejak Audit',
        access: 'Ya',
        detail: 'Membaca jejak keputusan.'
      },
      {
        menu: 'Upload BNBA',
        access: 'Tidak',
        detail: 'Input data milik Puskesmas/PJ; Pimpinan tidak mengunggah.'
      },
      {
        menu: 'Master Data',
        access: 'Tidak',
        detail: 'Master di luar kewenangan Pimpinan; Pimpinan membaca hasilnya di dashboard.'
      },
      {
        menu: 'Profil & Panduan',
        access: 'Ya',
        detail: 'Nama/password sendiri dan dokumentasi in-app.'
      }
    ],
    featureLinks: [
      {
        feature: 'Dashboard → Detail SPM',
        relation:
          'Angka dashboard dapat ditelusuri ke nilai IKU, nilai indikator, baris pengajuan, revisi, dan berkas asli.'
      },
      {
        feature: 'Faktor & TL',
        relation:
          'Faktor memberi konteks kualitatif; TL menerjemahkan masalah menjadi aksi ber-PJ. Keduanya tidak mengubah perhitungan.'
      },
      {
        feature: 'TL → Kabid/Sekretaris',
        relation:
          'Pimpinan menetapkan arah dan PJ; Kabid/Sekretaris mengeksekusi dan memperbarui status yang dipantau Pimpinan.'
      },
      {
        feature: 'Jejak Audit',
        relation:
          'Keputusan verifikasi/validasi dan perubahan TL dapat diperiksa ulang lewat jejak audit.'
      },
      {
        feature: 'Ekspor Dashboard',
        relation:
          'Ekspor memakai data dashboard yang sudah final, bukan angka yang diketik ulang.'
      }
    ],
    relations: [
      {
        actor: 'Puskesmas / PJ Program',
        relation:
          'Menyediakan data mentah lewat upload. Pimpinan tidak mengubah data tersebut; arahan diberikan setelah data final.'
      },
      {
        actor: 'Verifikator Dinas',
        relation:
          'Memeriksa data. Pimpinan tidak menggantikan verifikasi; fokus Pimpinan adalah keputusan dan tindak lanjut.'
      },
      {
        actor: 'Kabid / Sekretaris',
        relation:
          'Menerima arahan/TL dan mengisi faktor. Pimpinan menetapkan arah, Kabid/Sekretaris mengeksekusi, lalu Pimpinan memantau status.'
      }
    ],
    screenshots: [
      {
        src: '/user-guide/08-dashboard.webp',
        stage: 'Baca dashboard',
        size: 'wide',
        alt: 'Tabel Capaian SPM per indikator dengan kolom target, capaian, status, gap, dan aksi Detail yang ditandai kotak hijau.',
        caption: 'Tabel Capaian SPM dan aksi Detail',
        highlight:
          'Kotak hijau menandai aksi Detail pada tabel. Gunakan aksi ini untuk masuk ke sumber angka dan cakupan upload per puskesmas.',
        detail:
          'Screenshot difokuskan pada tabel Capaian SPM per indikator agar seluruh baris dan kolom aksi terbaca. Tabel menampilkan target, capaian, status berdasarkan target, gap, dan aksi Detail.',
        inline: 'full',
        access:
          'Dashboard → pilih periode dan Kelompok SPM → tabel Capaian SPM per indikator → klik Detail pada baris indikator.'
      },
      {
        src: '/user-guide/09-dashboard-detail.webp',
        stage: 'Telusuri detail',
        size: 'wide',
        alt: 'Tabel cakupan dan sumber angka per puskesmas dengan tautan pengajuan/berkas sumber yang ditandai kotak hijau.',
        caption: 'Cakupan & sumber angka per puskesmas',
        highlight:
          'Kotak hijau menandai pengajuan/berkas sumber. Tabel ini menjawab dari faskes, pengajuan, dan berkas mana angka dashboard berasal.',
        detail:
          'Screenshot difokuskan pada tabel cakupan & sumber angka per puskesmas: status upload, capaian, jumlah/total, pengajuan/berkas sumber, pengunggah, dan waktu pembaruan.',
        access:
          'Dashboard → pilih Kelompok SPM → klik Detail pada baris indikator → halaman Detail SPM → bagian Cakupan & sumber angka per puskesmas.'
      },
      {
        src: '/user-guide/12-factors-tl.webp',
        stage: 'Tetapkan TL',
        size: 'wide',
        access: 'Menu Faktor & TL → pilih periode → tab Tindak Lanjut → tabel TL → dropdown status.',
        alt: 'Tab Tindak Lanjut dengan TL baru dan dropdown status yang ditandai kotak hijau.',
        caption: 'Tindak lanjut dan status',
        highlight:
          'Kotak hijau menandai dropdown status TL. Pimpinan membuat TL lewat tombol TL baru, lalu memantau perubahan status bersama PJ.',
        detail:
          'Tab Tindak Lanjut menampilkan masalah, rekomendasi, PJ, status, dan waktu pembaruan. Faktor kualitatif diisi Kabid/Sekretaris; keduanya tidak mengubah perhitungan angka.'
      }
    ],
    notes: [
      'Pimpinan bukan penyetuju data. Persetujuan/pengembalian tetap di Verifikator dan Kabid/Sekretaris.',
      'Status dashboard dihitung dari target periode: capaian ≥ target = sesuai target; di bawah target = perlu perhatian. Ambang warna khusus belum diputuskan.',
      'Ekspor PDF/Excel tersedia untuk peran tertentu; peta, batas waktu layanan, dan integrasi SIKN belum diimplementasikan.',
      'Angka tanpa sumber yang jelas tidak boleh ditampilkan; jika nilai IKU belum ada, statusnya belum ada data.'
    ]
  },

};

export const ROLE_MATRIX = [
  {
    action: 'Input/kirim + upload BNBA',
    pkm: 'Ya',
    verifier: 'Tidak',
    kabid: 'Tidak',
    pimpinan: 'Tidak',
  },
  {
    action: 'Verifikasi / kembalikan ke pengirim',
    pkm: 'Tidak',
    verifier: 'Ya',
    kabid: 'Tidak',
    pimpinan: 'Tidak',
  },
  {
    action: 'Validasi final / kembalikan ke verifikator',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Ya',
    pimpinan: 'Tidak',
  },
  {
    action: 'Isi faktor pendukung/penghambat',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Ya (wajib)',
    pimpinan: 'Tidak',
  },
  {
    action: 'Beri arahan / tetapkan TL',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Eksekusi',
    pimpinan: 'Ya',
  },
  {
    action: 'Lihat dashboard',
    pkm: 'Milik sendiri',
    verifier: 'Bidangnya',
    kabid: 'Bidangnya',
    pimpinan: 'Semua',
  },
  {
    action: 'Lihat jejak audit',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Tidak',
    pimpinan: 'Ya',
  },
  {
    action: 'Perubahan pasca-final',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Terbatas',
    pimpinan: 'Tidak',
  }
];

export const STATUS_ROWS = [
  {
    status: 'draf',
    meaning: 'Baru diinput, belum dikirim',
    actor: 'Puskesmas/PJ',
    next: 'Upload/pratinjau → kirim'
  },
  {
    status: 'dikirim',
    meaning: 'Menunggu persetujuan Puskesmas atau verifikasi Dinas',
    actor: 'Puskesmas Approver / Verifikator',
    next: 'Setujui → lanjut; Kembalikan → Dikembalikan ke Petugas'
  },
  {
    status: 'diverifikasi',
    meaning: 'Lolos verifikasi, menunggu validasi',
    actor: 'Kabid/Sekretaris',
    next: 'Setujui → final; Kembalikan → Dikembalikan ke Verifikator'
  },
  {
    status: 'dikembalikan_petugas',
    meaning: 'Kembali ke Puskesmas + catatan',
    actor: 'Puskesmas/PJ',
    next: 'Perbaiki → kirim ulang'
  },
  {
    status: 'dikembalikan_verifikator',
    meaning: 'Kembali ke verifikator + catatan',
    actor: 'Verifikator',
    next: 'Cek ulang → setujui/kembalikan'
  },
  {
    status: 'final',
    meaning: 'Sah, terkunci, masuk perhitungan + dashboard',
    actor: 'Sistem',
    next: 'Tidak ada edit langsung; perubahan khusus + audit'
  },
  {
    status: 'diubah_khusus',
    meaning: 'Perubahan pasca-final dengan alasan + audit',
    actor: 'Kabid/Sekretaris + persetujuan kedua',
    next: 'Persetujuan ganda; layar pilot belum tersedia'
  }
];

export const OPEN_DECISIONS = [
  { topic: 'Kategori dashboard ke-9', status: 'Belum diputuskan', action: 'Tanyakan Dinas; jangan diisi dengan tebakan.' },
  { topic: 'Ambang warna/status khusus', status: 'Belum diputuskan', action: 'Dashboard pilot memakai perbandingan target.' },
  { topic: 'batas waktu layanan input/verifikasi/validasi', status: 'Belum diputuskan', action: 'Belum ada notifikasi/persetujuan otomatis.' },
  { topic: 'Target per wilayah/faskes', status: 'Belum diputuskan', action: 'Jangan menganggap target sama untuk semua faskes.' },
  { topic: 'Data pribadi di BNBA', status: 'Belum diputuskan', action: 'Jangan menambah kolom identitas pribadi.' },
  { topic: 'Penyimpanan berkas dan cadangan', status: 'Belum diputuskan', action: 'Minimal ikuti kebutuhan audit.' },
  { topic: 'Format BNBA final', status: 'Belum diputuskan', action: 'Template resmi belum dikunci.' },
  { topic: 'Duplikasi 3 baris ibu hamil', status: 'Belum diputuskan', action: 'Gabung atau bedakan (TW1/Mei/absolut?).' },
  { topic: 'Satuan baku indikator', status: 'Belum diputuskan', action: 'Normalisasi sebelum sistem produksi.' },
  { topic: 'PJ Data semua indikator', status: 'Belum lengkap', action: 'Lengkapi Bidang + Timker/Subbag.' },
  { topic: 'Tingkat persetujuan tambahan', status: 'Belum diputuskan', action: 'Konfigurasi hanya jika disetujui.' },
  { topic: 'Integrasi SIKN', status: 'Belum diputuskan', action: 'Jangan implementasi tanpa spesifikasi.' },
  { topic: 'Ekspor peta / PDF / Excel', status: 'Sebagian tersedia', action: 'Peta dan detail prioritas belum final.' },
  { topic: 'Notifikasi batas waktu layanan', status: 'Belum diputuskan', action: 'Kanal dan aturan belum ada.' }
];


