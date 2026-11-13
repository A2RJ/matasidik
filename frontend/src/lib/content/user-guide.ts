import type { Role } from '$lib/types';

export type GuideTabID = 'ringkasan' | 'puskesmas' | 'verifikator' | 'kabid' | 'pimpinan';

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
  label: string;
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
  access: GuideAccess[];
  featureLinks: GuideFeatureLink[];
  relations: GuideRelation[];
  screenshots: GuideScreenshot[];
  notes: string[];
}

export const GUIDE_TABS: GuideTab[] = [
  { id: 'ringkasan', label: 'Ringkasan & peta peran', short: 'Ringkasan' },
  { id: 'puskesmas', label: 'Puskesmas / PJ Program', short: 'Puskesmas' },
  { id: 'verifikator', label: 'Verifikator Dinas', short: 'Verifikator' },
  { id: 'kabid', label: 'Kabid / Sekretaris', short: 'Kabid' },
  { id: 'pimpinan', label: 'Pimpinan', short: 'Pimpinan' }
];

export const DEMO_ACCOUNTS = [
  { username: 'siti', role: 'Puskesmas Operator · PKM Taliwang', use: 'Upload, preview, kirim' },
  { username: 'pjgizi', role: 'PJ Program · Timker Gizi', use: 'Alur input dengan scope faskes pada akun' },
  { username: 'rina', role: 'Verifikator Dinas', use: 'Verifikasi, approve/return' },
  { username: 'kabid', role: 'Kabid Kesehatan Masyarakat', use: 'Validasi final, faktor' },
  { username: 'sekretaris', role: 'Sekretaris Dinas', use: 'Validasi final, faktor' },
  { username: 'pimpinan', role: 'Kepala Dinas', use: 'Dashboard, arahan, TL, audit' }
];

export const FLOW_LEGEND = [
  { swatch: 'green', text: 'Kotak hijau = titik mulai alur role pada modul ini' },
  { swatch: 'blue', text: 'Alur utama: masuk → unggah → verifikasi → validasi → hitung → dashboard' },
  { swatch: 'dashed', text: 'Garis putus-putus: pengembalian data atau dukungan lintas peran' }
];

const flowImage = '/user-guide/alur-panduan.webp';
const flowAlt =
  'Alur lengkap MATA SIDIK: masuk dan sesi, unggah BNBA Puskesmas, pratinjau dan koreksi IKU, verifikasi Dinas, validasi final Kabid/Sekretaris, kalkulasi dan dashboard, serta arahan dan TL Pimpinan.';

export const FLOW_IMAGE = { src: flowImage, alt: flowAlt };

export const MENU_CATALOG: MenuCatalogRow[] = [
  {
    menu: 'Dashboard',
    route: '/',
    function:
      'Ringkasan IKU induk, jumlah sesuai/di bawah target, filter periode–kelompok–faskes–status, chart multi-series, tabel SPM, laporan absolut bulanan, dan ekspor.',
    data:
      'Membaca IKU_VALUE, DASHBOARD_ROW, INDICATOR_VALUE, target, capaian tahunan, label/indikator, dan metadata provenance.',
    related:
      'Mengambil hasil akhir dari Upload → Submission → Workflow → Calculation. Setiap baris dapat dibuka ke Detail SPM.'
  },
  {
    menu: 'Detail SPM',
    route: '/dashboard/spm/{kode}',
    function:
      'Menampilkan provenance satu baris: capaian, target, gap, num/den, cakupan upload per puskesmas, chart, submission/file sumber, dan komponen angka.',
    data:
      'Membaca rantai DASHBOARD_ROW → IKU_VALUE → IKU_VALUE_DETAIL → INDICATOR_VALUE → SUBMISSION_ROW → SUBMISSION_REVISION → FILE_OBJECT.',
    related:
      'Subhalaman Dashboard. Menghubungkan angka dashboard ke Submission, file asli, dan audit.'
  },
  {
    menu: 'Upload BNBA',
    route: '/upload',
    function:
      'Memilih file .csv/.xlsx, periode, dan faskes; memantau job background; menampilkan preview baris; mengoreksi IKU leaf; mengirim ke workflow.',
    data:
      'Menulis SUBMISSION, FILE_OBJECT, SUBMISSION_ROW, INDICATOR_VALUE sementara, JOB, dan AUDIT_LOG. Membaca master indikator/IKU, periode, dan release.',
    related:
      'Bergantung pada Master Data. Setelah dikirim, data masuk ke Submission dan workflow verifikasi. Baris bermasalah menjadi blocker.'
  },
  {
    menu: 'Submission',
    route: '/submissions',
    function:
      'Daftar pengajuan per periode; detail baris, klasifikasi, blocker, timeline, revisi, file asli; aksi kirim/approve/return sesuai status dan role.',
    data:
      'Membaca SUBMISSION, SUBMISSION_ROW, SUBMISSION_REVIEW, SUBMISSION_REVISION, FILE_OBJECT. Menulis status, review, dan audit.',
    related:
      'Pusat workflow. Menerima hasil Upload, menjadi sumber nilai Calculation setelah final, dan dapat ditelusuri dari Dashboard.'
  },
  {
    menu: 'Master Data',
    route: '/master',
    function:
      'Mengelola indikator, label/IKU, kategori dashboard, bidang/timker, target, capaian tahunan, periode, dan pratinjau pemetaan.',
    data:
      'Menulis master indikator/label/kategori/target/periode. Menjadi kontrak kode, definisi, satuan, kategori, PJ, dan target.',
    related:
      'Mengendalikan Upload (kode dan release), Calculation (formula/komponen), dan Dashboard (kategori/target).'
  },
  {
    menu: 'Faktor & TL',
    route: '/factors',
    function:
      'Mengisi faktor pendukung/penghambat dan mengelola tindak lanjut: masalah → rekomendasi → PJ → status.',
    data:
      'Membaca IKU_VALUE. Menulis FAKTOR dan TINDAK_LANJUT, serta AUDIT_LOG.',
    related:
      'Memakai nilai final dari Dashboard/Calculation. Faktor tidak mengubah angka; TL menjadi jembatan keputusan Pimpinan ke eksekusi Kabid/Sekretaris.'
  },
  {
    menu: 'Audit Log',
    route: '/audit',
    function:
      'Menampilkan jejak append-only: waktu, aktor, aksi, entitas, alasan/catatan, dan request ID.',
    data:
      'Membaca AUDIT_LOG. Audit ditulis oleh upload, reklasifikasi, submit, approve/return, faktor, TL, dan perubahan master.',
    related:
      'Lintas semua fitur. Menjadi bukti keputusan dan penghubung antara aktor, data, dan waktu.'
  },
  {
    menu: 'Profil',
    route: '/profile',
    function:
      'Mengubah nama tampilan dan password sendiri; menampilkan role, jabatan, faskes, bidang, dan login terakhir.',
    data:
      'Menulis USER dan mencabut sesi refresh saat password berubah.',
    related:
      'Terhubung ke autentikasi dan seluruh menu yang memakai identitas/role pengguna.'
  },
  {
    menu: 'Panduan',
    route: '/user-guide',
    function:
      'Dokumentasi in-app per peran: alur, langkah, hak akses menu, keterkaitan fitur/data, status, screenshot, dan batasan pilot.',
    data:
      'Tidak menulis data bisnis. Membaca konten panduan dan aset screenshot.',
    related:
      'Menjadi titik rujukan untuk semua menu dan role. Tidak menggantikan kontrol akses aplikasi.'
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
    menu: 'Audit Log',
    pkm: 'Tidak',
    verifier: 'Tidak',
    kabid: 'Tidak',
    pimpinan: 'Ya',
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
    via: 'kode_iku · kode_indikator · periode · master_release',
    note:
      'Baris hanya sah jika kode dan relasinya ada di master. Master berubah lewat versi, bukan hapus fisik.'
  },
  {
    from: 'Upload BNBA',
    to: 'Submission',
    via: 'SUBMISSION · SUBMISSION_ROW · FILE_OBJECT · JOB',
    note:
      'Upload tidak menunggu parsing di request. Worker mengisi baris, klasifikasi, dan blocker; preview membaca hasilnya.'
  },
  {
    from: 'Submission',
    to: 'Workflow',
    via: 'SUBMISSION_STATUS · SUBMISSION_REVIEW · AUDIT_LOG',
    note:
      'Draf → kirim → verifikasi → validasi → final. Return wajib catatan; final mengunci data.'
  },
  {
    from: 'Workflow final',
    to: 'Calculation',
    via: 'IKU_VALUE · IKU_VALUE_DETAIL · CALCULATION_RUN',
    note:
      'Final memicu kalkulasi. IKU induk memakai RATIO_OF_SUMS dari komponen eksplisit, bukan rata-rata persentase.'
  },
  {
    from: 'Calculation',
    to: 'Dashboard',
    via: 'DASHBOARD_ROW · IKU_VALUE · INDICATOR_VALUE · provenance',
    note:
      'Dashboard membaca nilai materialized. Baris tanpa IKU_VALUE final berstatus belum ada data, bukan nol.'
  },
  {
    from: 'Dashboard',
    to: 'Faktor & TL',
    via: 'IKU_VALUE · FAKTOR · TINDAK_LANJUT',
    note:
      'Faktor adalah analisis kualitatif; TL adalah rencana aksi ber-PJ. Keduanya tidak mengubah kalkulasi.'
  },
  {
    from: 'Semua keputusan',
    to: 'Audit Log',
    via: 'AUDIT_LOG: aktor · waktu · aksi · catatan · diff · request_id',
    note:
      'Audit append-only menjadi bukti siapa melakukan apa, kapan, dan mengapa.'
  }
];

export const ROLE_GUIDES: Record<Exclude<GuideTabID, 'ringkasan'>, RoleGuide> = {
  puskesmas: {
    id: 'puskesmas',
    title: 'Puskesmas & PJ Program',
    level: 'Level 1 · Input dan pengiriman data',
    roles: ['pkm_operator', 'pkm_pj'],
    summary:
      'Menyiapkan file BNBA, mengunggah ke MATA SIDIK, memeriksa preview, mengoreksi IKU leaf, lalu mengirim submission ke verifikator Dinas.',
    startLabel: 'Mulai di Input BNBA',
    accounts: ['siti · Puskesmas Operator · PKM Taliwang', 'pjgizi · PJ Program · scope faskes pada akun'],
    flowOverlay: {
      left: 4.76,
      top: 36.36,
      width: 11.9,
      height: 26.14,
      label: 'Mulai Puskesmas / PJ'
    },
    steps: [
      {
        title: 'Masuk dan pastikan periode',
        where: '/login → /',
        body:
          'Masuk dengan akun Puskesmas/PJ. Setelah login, dashboard menampilkan periode aktif. Pilih periode yang sama dengan periode data BNBA sebelum membuka menu Upload BNBA.'
      },
      {
        title: 'Siapkan file BNBA',
        where: 'Template CSV di halaman Upload BNBA',
        body:
          'Gunakan file .csv atau .xlsx. Kolom wajib: periode, kode_faskes, kode_iku, kode_indikator, scope_type, scope_code, numerator, denominator. Kolom nilai dan catatan boleh diisi; kode_kategori opsional. Jangan memakai kode IKU induk untuk baris upload.'
      },
      {
        title: 'Unggah file',
        where: '/upload',
        body:
          'Pilih file BNBA, periode, dan faskes pengupload. Akun Puskesmas dengan faskes otomatis memakai faskes sendiri; pastikan faskes sesuai dengan akun. Klik Upload & proses.'
      },
      {
        title: 'Tunggu proses background',
        where: 'Kartu Proses background',
        body:
          'API menyimpan file dan job, lalu UI polling status. Langkah yang tampil: Validasi file → Klasifikasi → Simpan baris → Selesai. File tidak perlu diunggah ulang jika halaman di-refresh; buka kembali dari daftar submission.'
      },
      {
        title: 'Periksa preview dan koreksi IKU',
        where: 'Kartu Langkah 3 · Preview',
        body:
          'Baris MATCHED berwarna hijau. Baris MISMATCH, IKU_ROLLUP_ONLY, AMBIGUOUS, DUPLICATE, INVALID_VALUE, atau SCOPE_DENIED menjadi blocker. Gunakan dropdown Ubah IKU hanya untuk memilih IKU leaf yang terhubung ke indikator baris tersebut. Perubahan sebelum submit tercatat di audit.'
      },
      {
        title: 'Kirim ke workflow',
        where: 'Tombol Kirim ke workflow',
        body:
          'Tombol Kirim aktif hanya jika tidak ada blocker. Setelah dikirim, status berubah menjadi Dikirim dan submission masuk ke Verifikator Dinas. Data tidak bisa diubah langsung sampai ada keputusan return.'
      },
      {
        title: 'Tindak lanjuti hasil verifikasi',
        where: '/submissions → detail submission',
        body:
          'Jika verifikator mengembalikan, status menjadi Dikembalikan ke Petugas. Baca catatan, perbaiki baris yang diminta ( jika masih diizinkan ), lalu kirim ulang. Jika permintaan berasal dari Kabid/Sekretaris, statusnya Dikembalikan ke Verifikator dan verifikator yang memproses ulang.'
      }
    ],
    access: [
      {
        menu: 'Dashboard + Detail SPM',
        access: 'Milik faskes sendiri',
        detail: 'Melihat capaian final faskes dan provenance baris, tanpa melihat faskes lain.'
      },
      {
        menu: 'Upload BNBA',
        access: 'Ya',
        detail: 'Upload file, pantau job, preview, koreksi IKU leaf, lalu kirim.'
      },
      {
        menu: 'Submission',
        access: 'Milik faskes',
        detail: 'Daftar/detail, timeline, revisi, dan kirim ulang setelah dikembalikan.'
      },
      {
        menu: 'Profil & Panduan',
        access: 'Ya',
        detail: 'Nama/password sendiri dan dokumentasi in-app.'
      },
      {
        menu: 'Master, Faktor & TL, Audit',
        access: 'Tidak',
        detail: 'Bukan kewenangan role input; faktor/TL dan master dipegang role lain.'
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
          'File dan job menjadi submission; preview membaca SUBMISSION_ROW dan INDICATOR_VALUE sementara sebelum dikirim.'
      },
      {
        feature: 'Submission → Workflow',
        relation:
          'Kirim mengubah status menjadi Dikirim. Return dari verifikator membuka perbaikan; return dari Kabid/Sekretaris dikembalikan ke verifikator.'
      },
      {
        feature: 'Dashboard + Detail SPM',
        relation:
          'Setelah final, nilai muncul di dashboard dan dapat ditelusuri kembali ke baris submission dan file upload.'
      },
      {
        feature: 'Audit Log',
        relation:
          'Upload, reklasifikasi IKU, dan kirim ulang tercatat di audit, meskipun role ini tidak membuka menu Audit.'
      }
    ],
    relations: [
      {
        actor: 'Verifikator Dinas',
        relation:
          'Menerima submission berstatus Dikirim. Verifikator dapat approve ke Kabid/Sekretaris atau return ke Puskesmas dengan catatan.'
      },
      {
        actor: 'Kabid / Sekretaris',
        relation:
          'Menerima submission setelah verifikasi. Keputusan final atau return ke verifikator tidak dikirim langsung ke Puskesmas.'
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
        alt: 'Form Upload BNBA dengan pilihan file, periode, faskes, dan tombol Upload & proses.',
        caption: 'Form upload BNBA',
        highlight:
          'Kotak hijau menandai form yang harus diisi role Puskesmas/PJ: file, periode, faskes, lalu Upload & proses.',
        detail:
          'Screenshot diambil dari halaman /upload pada aplikasi pilot. Seluruh form dan tabel di sekitarnya tetap ditampilkan agar konteks tidak hilang.'
      },
      {
        src: '/user-guide/02-upload-preview.webp',
        alt: 'Preview BNBA dengan baris hasil parsing, status klasifikasi, pilihan IKU, dan tombol Kirim ke workflow.',
        caption: 'Preview dan tombol kirim',
        highlight:
          'Kotak hijau menandai tombol Kirim ke workflow. Periksa seluruh baris dan pastikan blocker nol sebelum menekan tombol.',
        detail:
          'Preview memperlihatkan baris, numerator/denominator, status klasifikasi, IKU, dan catatan. Dropdown Ubah IKU hanya tersedia untuk IKU leaf yang valid.'
      },
      {
        src: '/user-guide/03-submission-list.webp',
        alt: 'Tabel Submission per periode dengan kolom kode, faskes, uploader, periode, status, hash, diperbarui, dan aksi Buka.',
        caption: 'Daftar submission dan aksi buka',
        highlight:
          'Kotak hijau menandai aksi Buka pada tabel. Klik Buka untuk melihat baris, timeline, revisi, dan tombol yang tersedia untuk role Anda.',
        detail:
          'Role Puskesmas hanya melihat submission dari faskes sendiri. Filter periode di atas tabel membantu menemukan pengajuan yang dikembalikan.'
      }
    ],
    notes: [
      'Pilot belum menyediakan input angka manual; data masuk melalui file BNBA .csv atau .xlsx.',
      'IKU induk/rollup tidak boleh menerima upload langsung. Nilainya dihitung dari komponen anak yang eksplisit.',
      'Baris yang ambigu, tidak dikenal, duplikat, atau scope-nya di luar kewenangan masuk karantina/blocker; sistem tidak menebak IKU.',
      'Format BNBA final dan kebijakan PII belum diputuskan Dinas. Jangan menambahkan kolom identitas pribadi tanpa keputusan.'
    ]
  },
  verifikator: {
    id: 'verifikator',
    title: 'Verifikator Dinas',
    level: 'Level 2 · Pemeriksaan dan verifikasi',
    roles: ['dinas_verifier'],
    summary:
      'Memeriksa submission yang dikirim Puskesmas, menilai kelengkapan baris dan file, lalu meneruskan ke Kabid/Sekretaris atau mengembalikan ke Puskesmas dengan catatan.',
    startLabel: 'Mulai di Verifikasi',
    accounts: ['rina · Verifikator Dinas · Bidang Kesehatan Masyarakat'],
    flowOverlay: {
      left: 23.81,
      top: 36.36,
      width: 11.9,
      height: 26.14,
      label: 'Mulai Verifikator'
    },
    steps: [
      {
        title: 'Buka daftar submission',
        where: '/submissions',
        body:
          'Masuk dengan akun dinas_verifier, pilih periode, lalu buka tabel Submission per periode. Role ini melihat pengajuan lintas faskes sesuai kewenangannya, bukan hanya satu faskes.'
      },
      {
        title: 'Pilih submission berstatus Dikirim',
        where: 'Aksi Buka',
        body:
          'Cari baris dengan status Dikirim. Klik Buka atau klik baris untuk masuk ke detail. Submission berstatus lain belum/tidak berada di tahap verifikasi.'
      },
      {
        title: 'Periksa isi dan bukti',
        where: 'Detail submission',
        body:
          'Periksa ringkasan status workflow, jumlah blocker, baris hasil parsing, klasifikasi, IKU, numerator/denominator, catatan, timeline, revisi, hash file, dan tombol Unduh file asli. Jangan menyetujui baris yang masih blocker.'
      },
      {
        title: 'Setujui atau kembalikan',
        where: 'Panel Aksi → Approve / Return',
        body:
          'Approve mengubah status menjadi Diverifikasi dan meneruskan ke Kabid/Sekretaris. Return wajib catatan dan mengubah status menjadi Dikembalikan ke Petugas. Pilih salah satu sesuai hasil pemeriksaan.'
      },
      {
        title: 'Tangani pengembalian dari Kabid/Sekretaris',
        where: 'Status Dikembalikan ke Verifikator',
        body:
          'Jika Kabid/Sekretaris meminta pemeriksaan ulang, submission kembali ke tahap verifikasi. Periksa catatan, lakukan verifikasi ulang, lalu approve atau return lagi. Role ini tidak menetapkan final.'
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
        detail: 'Melihat daftar/detail dan menjalankan Approve/Return pada status Dikirim.'
      },
      {
        menu: 'Upload BNBA',
        access: 'Tidak',
        detail: 'Sumber data berasal dari Puskesmas; verifikator tidak mengunggah.'
      },
      {
        menu: 'Master Data',
        access: 'Tidak',
        detail: 'Master di luar kewenangan verifikator; verifikator memakai kontrak kode yang ada.'
      },
      {
        menu: 'Faktor & TL',
        access: 'Tidak',
        detail: 'Faktor diisi Kabid/Sekretaris; TL ditetapkan Pimpinan dan dieksekusi Kabid/Sekretaris.'
      },
      {
        menu: 'Audit Log',
        access: 'Tidak',
        detail: 'Audit ditulis sistem; menu Audit hanya untuk role pengawas.'
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
          'Submission berstatus Dikirim adalah input utama verifikasi. Verifikator tidak membuat atau mengedit baris upload.'
      },
      {
        feature: 'Submission → Kabid/Sekretaris',
        relation:
          'Approve meneruskan data ke validasi final. Return mengembalikan ke Puskesmas dengan catatan wajib.'
      },
      {
        feature: 'Submission → Audit',
        relation:
          'Setiap approve/return menulis SUBMISSION_REVIEW dan AUDIT_LOG berisi aktor, waktu, aksi, dan catatan.'
      },
      {
        feature: 'Dashboard + Detail SPM',
        relation:
          'Setelah Kabid/Sekretaris final, nilai masuk kalkulasi dan dashboard. Verifikator tidak menghitung ulang angka.'
      }
    ],
    relations: [
      {
        actor: 'Puskesmas / PJ Program',
        relation:
          'Mengirim submission. Verifikator dapat mengembalikannya dengan catatan; Puskesmas memperbaiki lalu mengirim ulang.'
      },
      {
        actor: 'Kabid / Sekretaris',
        relation:
          'Menerima hasil verifikasi. Kabid/Sekretaris dapat langsung final atau mengembalikan ke verifikator untuk dicek ulang.'
      },
      {
        actor: 'Pimpinan',
        relation:
          'Tidak menerima task verifikasi. Pimpinan membaca dashboard final dan memberi arahan/TL setelah data menjadi final.'
      }
    ],
    screenshots: [
      {
        src: '/user-guide/04-verif-list.webp',
        alt: 'Tabel submission untuk verifikator dengan kolom status dan aksi Buka yang ditandai kotak hijau.',
        caption: 'Daftar submission untuk verifikator',
        highlight:
          'Kotak hijau menandai aksi Buka pada tabel. Dari sini verifikator masuk ke detail dan melihat tombol Approve/Return yang sesuai status.',
        detail:
          'Tabel tetap ditampilkan penuh: kode, faskes, uploader, periode, status, hash file, waktu diperbarui, dan aksi.'
      },
      {
        src: '/user-guide/05-verif-detail.webp',
        alt: 'Detail submission verifikator dengan panel aksi Approve dan Return serta tabel baris hasil parsing.',
        caption: 'Panel keputusan verifikator',
        highlight:
          'Kotak hijau menandai panel Aksi berisi Approve dan Return. Return wajib catatan; Approve meneruskan ke Kabid/Sekretaris.',
        detail:
          'Di halaman yang sama terlihat status workflow, blocker, baris hasil parsing, timeline keputusan, revisi, dan unduhan file asli.'
      }
    ],
    notes: [
      'Verifikasi bukan final. Hanya Kabid/Sekretaris yang boleh menetapkan status final.',
      'Return wajib catatan. Tidak ada auto-approve, termasuk jika SLA terlewati.',
      'Puskesmas tidak bisa mengedit baris setelah status Dikirim; perubahan mengikuti alur return/revisi.',
      'Pada pilot, pembatasan dashboard per bidang belum sepenuhnya ditegakkan; lihat catatan implementasi pilot di ringkasan.',
      'Pilot belum menerapkan SLA dan notifikasi; batas waktu masih keputusan terbuka Dinas.'
    ]
  },
  kabid: {
    id: 'kabid',
    title: 'Kabid / Sekretaris',
    level: 'Level 3 · Validasi final dan faktor',
    roles: ['kabid', 'sekretaris'],
    summary:
      'Memvalidasi data yang sudah diverifikasi, menetapkan final dengan catatan penerimaan, lalu mengisi faktor pendukung/penghambat. Final mengunci data dan memicu kalkulasi.',
    startLabel: 'Mulai di Validasi final',
    accounts: ['kabid · Kabid Kesehatan Masyarakat', 'sekretaris · Sekretaris Dinas'],
    flowOverlay: {
      left: 42.26,
      top: 36.36,
      width: 11.9,
      height: 26.14,
      label: 'Mulai Kabid / Sekretaris'
    },
    steps: [
      {
        title: 'Buka submission berstatus Diverifikasi',
        where: '/submissions',
        body:
          'Masuk dengan akun kabid atau sekretaris. Pilih periode dan cari submission berstatus Diverifikasi; status ini berarti pemeriksaan Dinas sudah selesai dan menunggu validasi final.'
      },
      {
        title: 'Periksa data dan jejak',
        where: 'Detail submission',
        body:
          'Periksa baris, klasifikasi, IKU, nilai, blocker, file asli, timeline, dan revisi. Pastikan catatan verifikator dan data yang diajukan sudah sesuai sebelum memberi keputusan.'
      },
      {
        title: 'Tetapkan final atau kembalikan',
        where: 'Panel Aksi → Approve / Return',
        body:
          'Approve final wajib catatan penerimaan. Setelah final, data terkunci dan sistem memicu kalkulasi IKU/SPM. Return mengembalikan ke verifikator dengan catatan; bukan langsung ke Puskesmas.'
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
    access: [
      {
        menu: 'Dashboard + Detail SPM',
        access: 'Bidang*',
        detail: 'Melihat capaian bidang dan menelusuri provenance; batas bidang pilot belum sepenuhnya ditegakkan.'
      },
      {
        menu: 'Submission',
        access: 'Validasi',
        detail: 'Approve final dengan catatan penerimaan atau Return ke verifikator dengan catatan.'
      },
      {
        menu: 'Master Data',
        access: 'Lihat + target*',
        detail: 'Melihat master; target dapat diubah sesuai kewenangan, CRUD master lain di luar kewenangan role ini.'
      },
      {
        menu: 'Faktor & TL',
        access: 'Faktor + status TL',
        detail: 'Mengisi faktor pendukung/penghambat dan memperbarui status TL yang ditugaskan.'
      },
      {
        menu: 'Upload BNBA',
        access: 'Tidak',
        detail: 'Input file milik Puskesmas/PJ.'
      },
      {
        menu: 'Audit Log',
        access: 'Tidak',
        detail: 'Audit hanya dibuka role pengawas; keputusan tetap tercatat.'
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
        feature: 'Final → Calculation → Dashboard',
        relation:
          'Final memicu kalkulasi IKU/SPM. Dashboard membaca nilai materialized dan provenance, bukan hitung ulang saat request.'
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
        feature: 'Semua keputusan → Audit Log',
        relation:
          'Approve, return, faktor, dan perubahan TL menulis audit meskipun role ini tidak membuka menu Audit.'
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
          'Tidak menerima return langsung dari Kabid/Sekretaris. Perbaikan dari Puskesmas melewati verifikator.'
      },
      {
        actor: 'Sistem',
        relation:
          'Final memicu kalkulasi. Dashboard membaca nilai materialized dan rantai provenance sampai file; bukan angka yang diketik ulang.'
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
        alt: 'Detail submission Kabid dengan panel Approve dan Return serta catatan penerimaan.',
        caption: 'Validasi final Kabid/Sekretaris',
        highlight:
          'Kotak hijau menandai panel Aksi final. Approve wajib catatan penerimaan; Return wajib catatan perbaikan.',
        detail:
          'Halaman detail menampilkan status workflow, blocker, baris hasil parsing, timeline, revisi, dan file asli. Setelah final, data terkunci.'
      },
      {
        src: '/user-guide/07-factors-action.webp',
        alt: 'Tabel Faktor & TL dengan aksi Isi atau Ubah yang ditandai kotak hijau.',
        caption: 'Isi faktor dan pantau TL',
        highlight:
          'Kotak hijau menandai aksi Isi/Ubah pada tabel faktor. Faktor tidak mengubah kalkulasi; TL dikelola pada tab terpisah.',
        detail:
          'Tabel faktor menampilkan IKU, capaian, num/den, faktor pendukung, faktor penghambat, dan aksi. Tab Tindak Lanjut menampilkan masalah, rekomendasi, PJ, status, dan waktu pembaruan.'
      }
    ],
    notes: [
      'Final wajib catatan penerimaan. Setelah final, data tidak boleh diubah langsung dari UI.',
      'Return dari tahap ini kembali ke verifikator, bukan langsung ke Puskesmas.',
      'Pada pilot, pembatasan dashboard per bidang belum sepenuhnya ditegakkan; lihat catatan implementasi pilot di ringkasan.',
      'Pilot belum memiliki layar SPECIAL_CHANGE/dual approval; perubahan pasca-final masih menjadi batasan yang harus diselesaikan sebelum produksi.',
      'Ambang warna/status formal dan SLA belum diputuskan Dinas.'
    ]
  },
  pimpinan: {
    id: 'pimpinan',
    title: 'Pimpinan',
    level: 'Observer & pengambil keputusan',
    roles: ['pimpinan'],
    summary:
      'Membaca dashboard final, menelusuri angka sampai file, melihat faktor, memberi arahan, menetapkan tindak lanjut, dan memantau penyelesaiannya.',
    startLabel: 'Mulai di Dashboard SPM',
    accounts: ['pimpinan · Kepala Dinas'],
    flowOverlay: {
      left: 80.65,
      top: 36.36,
      width: 13.69,
      height: 26.14,
      label: 'Mulai Pimpinan'
    },
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
          'Klik baris atau tombol Detail. Halaman detail menampilkan provenance, cakupan upload per puskesmas, chart antar faskes, submission/file sumber, dan komponen angka. Rantai yang dijamin: dashboard → IKU_VALUE → INDICATOR_VALUE → SUBMISSION_ROW → SUBMISSION → FILE_OBJECT.'
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
          'Audit Log menampilkan waktu, aktor, aksi, entitas, alasan/catatan, dan request ID. Pimpinan dapat memfilter berdasarkan tipe entitas. Audit adalah catatan bisnis append-only, bukan log teknis.'
      }
    ],
    access: [
      {
        menu: 'Dashboard + Detail SPM',
        access: 'Semua',
        detail: 'Melihat seluruh faskes, kategori, filter, chart, ekspor, dan provenance sampai file.'
      },
      {
        menu: 'Submission',
        access: 'Lihat',
        detail: 'Melihat status dan detail untuk konteks; tidak melakukan approve/return.'
      },
      {
        menu: 'Faktor & TL',
        access: 'TL + status',
        detail: 'Membuat TL (masalah → rekomendasi → PJ → status) dan memperbarui status; faktor dibaca.'
      },
      {
        menu: 'Audit Log',
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
          'Angka dashboard dapat ditelusuri ke IKU_VALUE, INDICATOR_VALUE, baris submission, revisi, dan file asli.'
      },
      {
        feature: 'Faktor & TL',
        relation:
          'Faktor memberi konteks kualitatif; TL menerjemahkan masalah menjadi aksi ber-PJ. Keduanya tidak mengubah kalkulasi.'
      },
      {
        feature: 'TL → Kabid/Sekretaris',
        relation:
          'Pimpinan menetapkan arah dan PJ; Kabid/Sekretaris mengeksekusi dan memperbarui status yang dipantau Pimpinan.'
      },
      {
        feature: 'Audit Log',
        relation:
          'Keputusan verifikasi/validasi dan perubahan TL dapat diperiksa ulang lewat audit append-only.'
      },
      {
        feature: 'Ekspor Dashboard',
        relation:
          'Ekspor memakai data dashboard yang sudah materialized, bukan angka yang diketik ulang.'
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
        alt: 'Tabel Capaian SPM per indikator dengan kolom target, capaian, status, gap, dan aksi Detail yang ditandai kotak hijau.',
        caption: 'Tabel Capaian SPM dan aksi Detail',
        highlight:
          'Kotak hijau menandai aksi Detail pada tabel. Gunakan aksi ini untuk masuk ke provenance dan cakupan upload per puskesmas.',
        detail:
          'Screenshot difokuskan pada tabel Capaian SPM per indikator agar seluruh baris dan kolom aksi terbaca. Tabel menampilkan target, capaian, status berdasarkan target, gap, dan aksi Detail.',
        inline: 'full',
        access:
          'Dashboard → pilih periode dan Kelompok SPM → tabel Capaian SPM per indikator → klik Detail pada baris indikator.'
      },
      {
        src: '/user-guide/09-dashboard-detail.webp',
        alt: 'Tabel cakupan dan provenance per puskesmas dengan tautan submission/berkas sumber yang ditandai kotak hijau.',
        caption: 'Cakupan & provenance per puskesmas',
        highlight:
          'Kotak hijau menandai submission/berkas sumber. Tabel ini menjawab dari faskes, submission, dan file mana angka dashboard berasal.',
        detail:
          'Screenshot difokuskan pada tabel cakupan & provenance per puskesmas: status upload, capaian, num/den, submission/berkas sumber, pengunggah, dan waktu pembaruan.',
        access:
          'Dashboard → pilih Kelompok SPM → klik Detail pada baris indikator → halaman Detail SPM → bagian Cakupan & provenance per puskesmas.'
      },
      {
        src: '/user-guide/12-factors-tl.webp',
        alt: 'Tab Tindak Lanjut dengan TL baru dan dropdown status yang ditandai kotak hijau.',
        caption: 'Tindak lanjut dan status',
        highlight:
          'Kotak hijau menandai dropdown status TL. Pimpinan membuat TL lewat tombol TL baru, lalu memantau perubahan status bersama PJ.',
        detail:
          'Tab Tindak Lanjut menampilkan masalah, rekomendasi, PJ, status, dan waktu pembaruan. Faktor kualitatif diisi Kabid/Sekretaris; keduanya tidak mengubah kalkulasi angka.'
      }
    ],
    notes: [
      'Pimpinan bukan approver data. Approve/Return tetap di Verifikator dan Kabid/Sekretaris.',
      'Status dashboard dihitung dari target periode: capaian ≥ target = sesuai target; di bawah target = perlu perhatian. Ambang warna khusus belum diputuskan.',
      'Ekspor PDF/Excel tersedia untuk role tertentu; peta, SLA, dan integrasi SIKN belum diimplementasikan.',
      'Angka tanpa provenance tidak boleh ditampilkan; jika IKU_VALUE belum ada, statusnya belum ada data.'
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
    action: 'Lihat audit',
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
    next: 'Upload/preview → kirim'
  },
  {
    status: 'dikirim',
    meaning: 'Menunggu verifikasi Dinas',
    actor: 'Verifikator',
    next: 'Approve → diverifikasi; Return → dikembalikan_petugas'
  },
  {
    status: 'diverifikasi',
    meaning: 'Lolos verifikasi, menunggu validasi',
    actor: 'Kabid/Sekretaris',
    next: 'Approve → final; Return → dikembalikan_verifikator'
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
    next: 'Cek ulang → approve/return'
  },
  {
    status: 'final',
    meaning: 'Sah, terkunci, masuk kalkulasi + dashboard',
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
  { topic: 'SLA input/verifikasi/validasi', status: 'Belum diputuskan', action: 'Belum ada notifikasi/auto-approve.' },
  { topic: 'Target per wilayah/faskes', status: 'Belum diputuskan', action: 'Jangan menganggap target sama untuk semua faskes.' },
  { topic: 'PII di BNBA', status: 'Belum diputuskan', action: 'Jangan menambah kolom identitas pribadi.' },
  { topic: 'Retensi file dan backup', status: 'Belum diputuskan', action: 'Minimal ikuti kebutuhan audit.' },
  { topic: 'Format BNBA final', status: 'Belum diputuskan', action: 'Template resmi belum dikunci.' },
  { topic: 'Duplikasi 3 baris ibu hamil', status: 'Belum diputuskan', action: 'Gabung atau bedakan (TW1/Mei/absolut?).' },
  { topic: 'Satuan baku indikator', status: 'Belum diputuskan', action: 'Normalisasi sebelum engine produksi.' },
  { topic: 'PJ Data semua indikator', status: 'Belum lengkap', action: 'Lengkapi Bidang + Timker/Subbag.' },
  { topic: 'Level approval tambahan', status: 'Belum diputuskan', action: 'Konfigurasi hanya jika disetujui.' },
  { topic: 'Integrasi SIKN', status: 'Belum diputuskan', action: 'Jangan implementasi tanpa spesifikasi.' },
  { topic: 'Ekspor peta / PDF / Excel', status: 'Sebagian tersedia', action: 'Peta dan detail prioritas belum final.' },
  { topic: 'Notifikasi SLA', status: 'Belum diputuskan', action: 'Kanal dan aturan belum ada.' }
];


