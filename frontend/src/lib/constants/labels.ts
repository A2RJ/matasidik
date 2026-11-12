import {
  LayoutDashboard,
  Upload,
  FileText,
  Database,
  MessageSquare,
  ScrollText,
  Users
} from '@lucide/svelte';
import type { Role } from '$lib/types';

export const APP_NAME = 'MATA SIDIK';
export const APP_SUBTITLE = 'Manajemen Data dan Sistem Informasi Dinas Kesehatan';
export const APP_REGION = 'Kabupaten Sumbawa Barat';

export const NAV_ITEMS: { href: string; label: string; icon: typeof LayoutDashboard; description: string; roles?: Role[] }[] = [
  { href: '/', label: 'Dashboard', icon: LayoutDashboard, description: 'Capaian SPM & IKU' },
  { href: '/upload', label: 'Upload BNBA', icon: Upload, description: 'Kirim & preview data', roles: ['pkm_operator', 'pkm_pj', 'admin'] },
  { href: '/submissions', label: 'Submission', icon: FileText, description: 'Workflow & revisi', roles: ['pkm_operator', 'pkm_pj', 'dinas_verifier', 'kabid', 'sekretaris', 'pimpinan', 'admin'] },
  { href: '/master', label: 'Master Data', icon: Database, description: 'Indikator, label, target', roles: ['kabid', 'sekretaris', 'admin'] },
  { href: '/factors', label: 'Faktor & TL', icon: MessageSquare, description: 'Analisis & tindak lanjut', roles: ['kabid', 'sekretaris', 'pimpinan', 'admin'] },
  { href: '/audit', label: 'Audit Log', icon: ScrollText, description: 'Jejak keputusan', roles: ['pimpinan', 'admin'] },
  { href: '/admin', label: 'Administrasi', icon: Users, description: 'Faskes, user, posisi', roles: ['admin'] }
];

export const PAGE_COPY = {
  dashboard: {
    title: 'Dashboard SPM & IKU',
    description:
      'Capaian final dari alur upload, verifikasi, validasi, dan kalkulasi. Setiap angka dapat ditelusuri sampai file asli.'
  },
  upload: {
    title: 'Upload BNBA',
    description:
      'Unggah file BNBA Puskesmas, pantau proses background, koreksi IKU di preview, lalu kirim ke workflow.'
  },
  submissions: {
    title: 'Submission & Workflow',
    description: 'Pantau status, preview baris, dan keputusan verifikasi/validasi.'
  },
  master: {
    title: 'Master Data',
    description: 'Label, indikator, target, dan periode yang mengendalikan dashboard serta kalkulasi.'
  },
  factors: {
    title: 'Faktor & Tindak Lanjut',
    description:
      'Faktor adalah analisis kualitatif pendukung/penghambat; TL adalah rencana aksi ber-PJ dan berstatus.'
  },
  audit: {
    title: 'Audit Log',
    description: 'Siapa, kapan, apa, dan mengapa—setiap keputusan tercatat.'
  }
} as const;

export const TABLE_LABELS = {
  code: 'Kode',
  indicator: 'Indikator',
  label: 'Label',
  target: 'Target',
  baselinePrev: 'Capaian sebelumnya',
  baselineTW1: 'Capaian TW1',
  value: 'Capaian',
  status: 'Status',
  period: 'Periode',
  faskes: 'Faskes',
  uploader: 'Uploader',
  action: 'Aksi',
  numerator: 'Numerator',
  denominator: 'Denominator',
  gap: 'Gap',
  updated: 'Diperbarui'
} as const;

export const EMPTY_STATE = {
  dashboard: 'Belum ada baris dashboard pada periode ini.',
  submissions: 'Belum ada submission. Mulai dari Upload BNBA.',
  factors: 'Belum ada nilai IKU untuk diisi faktor.',
  audit: 'Belum ada audit log.',
  table: 'Tidak ada data.'
} as const;
