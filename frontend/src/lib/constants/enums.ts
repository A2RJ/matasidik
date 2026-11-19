import type {
  ClassificationStatus,
  Role,
  SubmissionStatus
} from '$lib/types';

export const ROLE_LABELS: Record<Role, string> = {
  pkm_operator: 'Puskesmas Operator',
  pkm_pj: 'PJ Program',
  pkm_approver: 'Puskesmas Approver',
  dinas_verifier: 'Verifikator Dinas',
  kabid: 'Kabid',
  sekretaris: 'Sekretaris',
  pimpinan: 'Pimpinan',
  admin: 'Admin'
};

export const ROLE_SHORT: Record<Role, string> = {
  pkm_operator: 'PKM',
  pkm_pj: 'PJ',
  pkm_approver: 'Puskesmas',
  dinas_verifier: 'Dinas',
  kabid: 'Kabid',
  sekretaris: 'Sekretaris',
  pimpinan: 'Pimpinan',
  admin: 'Admin'
};

export const SUBMISSION_STATUS_LABELS: Record<SubmissionStatus, string> = {
  draf: 'Draf',
  dikirim: 'Dikirim',
  diverifikasi: 'Diverifikasi',
  dikembalikan_petugas: 'Dikembalikan ke Petugas',
  dikembalikan_verifikator: 'Dikembalikan ke Verifikator',
  final: 'Final · Terkunci',
  diubah_khusus: 'Diubah Khusus'
};

export type Tone = 'neutral' | 'info' | 'success' | 'warning' | 'danger' | 'final' | 'dark' | 'slate';

export const SUBMISSION_STATUS_TONE: Record<SubmissionStatus, Tone> = {
  draf: 'neutral',
  dikirim: 'info',
  diverifikasi: 'info',
  dikembalikan_petugas: 'warning',
  dikembalikan_verifikator: 'warning',
  final: 'final',
  diubah_khusus: 'danger'
};

export const CLASSIFICATION_LABELS: Record<ClassificationStatus, string> = {
  MATCHED: 'Sesuai',
  MISMATCH: 'Tidak sesuai',
  UNKNOWN_IKU: 'IKU tidak dikenal',
  UNKNOWN_INDICATOR: 'Indikator tidak dikenal',
  IKU_ROLLUP_ONLY: 'IKU induk',
  AMBIGUOUS: 'Perlu dipilih',
  DUPLICATE: 'Duplikat',
  INVALID_VALUE: 'Nilai tidak valid',
  SCOPE_DENIED: 'Di luar kewenangan'
};

export const CLASSIFICATION_TONE: Record<ClassificationStatus, Tone> = {
  MATCHED: 'success',
  MISMATCH: 'danger',
  UNKNOWN_IKU: 'danger',
  UNKNOWN_INDICATOR: 'danger',
  IKU_ROLLUP_ONLY: 'warning',
  AMBIGUOUS: 'warning',
  DUPLICATE: 'warning',
  INVALID_VALUE: 'danger',
  SCOPE_DENIED: 'danger'
};

export const DISPLAY_STATUS_LABELS: Record<string, string> = {
  ON_TARGET: 'Sesuai target',
  BELOW_TARGET: 'Di bawah target',
  NO_TARGET: 'Belum ada target',
  NO_DATA: 'Belum ada data',
  STALE: 'Perlu hitung ulang'
};

export const DISPLAY_STATUS_TONE: Record<string, Tone> = {
  ON_TARGET: 'success',
  BELOW_TARGET: 'warning',
  NO_TARGET: 'neutral',
  NO_DATA: 'neutral',
  STALE: 'warning'
};

export const VALUE_STATUS_LABELS: Record<string, string> = {
  COMPLETE: 'Lengkap',
  PARTIAL: 'Sebagian',
  NO_DATA: 'Belum ada data',
  STALE: 'Stale'
};

export const VALUE_STATUS_TONE: Record<string, Tone> = {
  COMPLETE: 'success',
  PARTIAL: 'warning',
  NO_DATA: 'neutral',
  STALE: 'warning'
};

export const TL_STATUS_LABELS: Record<string, string> = {
  open: 'Terbuka',
  in_progress: 'Dikerjakan',
  done: 'Selesai',
  cancelled: 'Dibatalkan'
};

export const TL_STATUS_TONE: Record<string, Tone> = {
  open: 'warning',
  in_progress: 'info',
  done: 'success',
  cancelled: 'neutral'
};

export const WORKFLOW_STEP_LABELS: Record<string, string> = {
  DRAFT_PKM: 'Draf Puskesmas',
  APPROVE_PKM: 'Persetujuan Puskesmas',
  VERIFY_DINAS: 'Verifikasi Dinas',
  VALIDATE_KABID: 'Validasi Kabid',
  FINAL: 'Final',
  SPECIAL_CHANGE: 'Perubahan Khusus'
};

export const SCOPE_TYPE_LABELS: Record<string, string> = {
  KABUPATEN: 'Kabupaten',
  KECAMATAN: 'Kecamatan',
  FASKES: 'Puskesmas/UPTD',
  DESA: 'Desa'
};

export const AGGREGATION_LABELS: Record<string, string> = {
  RATIO_OF_SUMS: 'Rasio jumlah',
  SUM: 'Penjumlahan',
  WEIGHTED_AVERAGE: 'Rata-rata berbobot',
  FORMULA: 'Rumus khusus',
  NONE: 'Tanpa agregasi'
};

export const AGGREGATION_OPTIONS = Object.entries(AGGREGATION_LABELS).map(([value, label]) => ({
  value,
  label
}));

export const IKU_ACTION_LABELS: Record<string, string> = {
  CREATE: 'Tambah IKU',
  UPDATE: 'Ubah IKU',
  DEACTIVATE: 'Nonaktifkan IKU'
};

export const IKU_STATUS_LABELS: Record<string, string> = {
  draf: 'Draf',
  menunggu_puskesmas: 'Menunggu Puskesmas',
  dikembalikan_puskesmas: 'Dikembalikan Puskesmas',
  menunggu_admin: 'Menunggu Admin',
  dikembalikan_admin: 'Dikembalikan Admin',
  disetujui: 'Disetujui',
  ditolak: 'Ditolak',
  dibatalkan: 'Dibatalkan'
};

export const REVIEW_ACTION_LABELS: Record<string, string> = {
  submit: 'Dikirim',
  approve: 'Disetujui',
  return: 'Dikembalikan',
  reject: 'Ditolak'
};

export const CATEGORY_TYPE_LABELS: Record<string, string> = {
  LABEL: 'Kelompok IKU',
  OWNER: 'Bidang/Timker'
};

export const ORGANIZATION_TYPE_LABELS: Record<string, string> = {
  bidang: 'Bidang',
  sekretariat: 'Sekretariat',
  timker: 'Tim Kerja'
};

export const JOB_STATUS_LABELS: Record<string, string> = {
  queued: 'Menunggu',
  processing: 'Diproses',
  retrying: 'Diproses ulang',
  done: 'Selesai',
  failed: 'Gagal',
  dead: 'Gagal permanen'
};

export const AUDIT_ENTITY_LABELS: Record<string, string> = {
  submission: 'Submission',
  submission_row: 'Baris submission',
  iku_value: 'Nilai IKU',
  tindak_lanjut: 'Tindak lanjut',
  label: 'IKU',
  indikator: 'Indikator',
  user: 'Pengguna',
  faskes: 'Faskes',
  position: 'Jabatan',
  dashboard_category: 'Kategori dashboard',
  bidang: 'Bidang/Timker',
  target: 'Target',
  capaian_tahunan: 'Capaian tahunan',
  iku_change_request: 'Pengajuan IKU'
};

export const AUDIT_ACTION_LABELS: Record<string, string> = {
  'submission.created': 'Submission dibuat',
  'submission.submitted': 'Submission dikirim',
  'submission.verified': 'Submission diverifikasi',
  'submission.final': 'Submission difinalkan',
  'submission.finalized': 'Submission difinalkan',
  'submission.row_reclassified': 'IKU baris diperbarui',
  'faktor.upsert': 'Faktor disimpan',
  'tindak_lanjut.created': 'Tindak lanjut dibuat',
  'tindak_lanjut.updated': 'Tindak lanjut diperbarui',
  'profile.updated': 'Profil diperbarui',
  'profile.password_changed': 'Password diubah',
  'auth.acting_as': 'Mode lihat sebagai dimulai',
  'iku.change_requested': 'Pengajuan IKU dikirim',
  'iku.change_approve': 'Pengajuan IKU disetujui',
  'iku.change_return': 'Pengajuan IKU dikembalikan',
  'iku.change_reject': 'Pengajuan IKU ditolak',
  'master.label_created': 'IKU dibuat',
  'master.label_updated': 'IKU diperbarui',
  'master.label_deactivated': 'IKU dinonaktifkan',
  'master.category_created': 'Kategori dibuat',
  'master.category_updated': 'Kategori diperbarui',
  'master.category_deactivated': 'Kategori dinonaktifkan',
  'master.bidang_created': 'Bidang/Timker dibuat',
  'master.bidang_updated': 'Bidang/Timker diperbarui',
  'master.indicator_created': 'Indikator dibuat',
  'master.target_upserted': 'Target disimpan',
  'master.capaian_tahunan_created': 'Capaian tahunan dibuat',
  'master.capaian_tahunan_updated': 'Capaian tahunan diperbarui',
  'master.capaian_tahunan_deleted': 'Capaian tahunan dihapus',
  'admin.user_created': 'Pengguna dibuat',
  'admin.user_updated': 'Pengguna diperbarui',
  'admin.user_password_reset': 'Password pengguna direset',
  'admin.faskes_created': 'Faskes dibuat',
  'admin.faskes_updated': 'Faskes diperbarui',
  'admin.position_created': 'Jabatan dibuat',
  'admin.position_updated': 'Jabatan diperbarui'
};

export function auditActionLabel(action: string): string {
  if (AUDIT_ACTION_LABELS[action]) return AUDIT_ACTION_LABELS[action];
  const words = action.split('.').pop()?.replaceAll('_', ' ') ?? action;
  return words.charAt(0).toUpperCase() + words.slice(1);
}

export function auditEntityLabel(entity: string): string {
  return AUDIT_ENTITY_LABELS[entity] ?? entity.replaceAll('_', ' ');
}

export const SCOPE_TYPES = ['KABUPATEN', 'KECAMATAN', 'FASKES', 'DESA'] as const;

export const DASHBOARD_CATEGORIES = [
  'IKU',
  'IKK',
  'SPM',
  'RIBK',
  'Program',
  'Kegiatan',
  'Sub Kegiatan',
  'Bidang/Sekretariat/Tim Kerja'
] as const;

export const STATUS_TARGET_NOTE =
  'Status dihitung dari target periode: capaian ≥ target = sesuai target; di bawah target = perlu perhatian.';
