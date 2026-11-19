import {
  BookOpen,
  LayoutDashboard,
  Upload,
  FileText,
  Database,
  MessageSquare,
  ScrollText,
  Users
} from '@lucide/svelte';

export const APP_NAME = 'MATA SIDIK';
export const APP_SUBTITLE = 'Manajemen Data dan Sistem Informasi Dinas Kesehatan';
export const APP_REGION = 'Kabupaten Sumbawa Barat';

export const NAV_ITEMS: { href: string; label: string; icon: typeof LayoutDashboard; description: string }[] = [
  { href: '/', label: 'Dashboard', icon: LayoutDashboard, description: 'Capaian SPM & IKU' },
  { href: '/user-guide', label: 'Panduan', icon: BookOpen, description: 'Tata cara per peran' },
  { href: '/upload', label: 'Upload BNBA', icon: Upload, description: 'Kirim & periksa data' },
  { href: '/submissions', label: 'Submission', icon: FileText, description: 'Status & keputusan' },
  { href: '/iku-requests', label: 'Pengajuan IKU', icon: FileText, description: 'Ajukan & setujui IKU' },
  { href: '/master', label: 'Master Data', icon: Database, description: 'Indikator, IKU, target' },
  { href: '/factors', label: 'Faktor & TL', icon: MessageSquare, description: 'Analisis & tindak lanjut' },
  { href: '/audit', label: 'Jejak Audit', icon: ScrollText, description: 'Riwayat keputusan' },
  { href: '/admin', label: 'Administrasi', icon: Users, description: 'Faskes, user, posisi' }
];

export const PAGE_COPY = {
  dashboard: {
    title: 'Dashboard SPM & IKU',
    description:
      'Capaian final dari alur unggah, verifikasi, validasi, dan perhitungan. Setiap angka dapat ditelusuri sampai berkas asli.'
  },
  upload: {
    title: 'Upload BNBA',
    description:
      'Unggah berkas BNBA Puskesmas, pantau prosesnya, koreksi IKU pada pratinjau, lalu kirim ke alur persetujuan.'
  },
  submissions: {
    title: 'Submission',
    description: 'Pantau status, rincian baris, dan keputusan verifikasi/validasi.'
  },
  master: {
    title: 'Master Data',
    description: 'Indikator, IKU, target, dan periode yang menjadi acuan dashboard serta perhitungan.'
  },
  factors: {
    title: 'Faktor & Tindak Lanjut',
    description:
      'Faktor adalah analisis kualitatif pendukung/penghambat; TL adalah rencana aksi dengan PJ dan status.'
  },
  audit: {
    title: 'Jejak Audit',
    description: 'Siapa, kapan, apa, dan mengapa—setiap keputusan tercatat.'
  }
} as const;

export const TABLE_LABELS = {
  code: 'Kode',
  indicator: 'Indikator',
  label: 'Label/IKU',
  target: 'Target',
  baselinePrev: 'Capaian sebelumnya',
  baselineTW1: 'Capaian TW1',
  value: 'Capaian',
  status: 'Status',
  period: 'Periode',
  faskes: 'Faskes',
  uploader: 'Pengunggah',
  action: 'Aksi',
  numerator: 'Jumlah',
  denominator: 'Total',
  gap: 'Gap',
  updated: 'Diperbarui'
} as const;

export const EMPTY_STATE = {
  dashboard: 'Belum ada baris dashboard pada periode ini.',
  submissions: 'Belum ada submission pada periode ini.',
  factors: 'Belum ada nilai IKU untuk diisi faktor.',
  audit: 'Belum ada riwayat keputusan.',
  table: 'Tidak ada data.'
} as const;
