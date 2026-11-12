import type {
  ClassificationStatus,
  Role,
  SubmissionStatus
} from '$lib/types';

export const ROLE_LABELS: Record<Role, string> = {
  pkm_operator: 'Puskesmas Operator',
  pkm_pj: 'PJ Program',
  dinas_verifier: 'Verifikator Dinas',
  kabid: 'Kabid',
  sekretaris: 'Sekretaris',
  iku_owner: 'IKU Owner',
  pimpinan: 'Pimpinan',
  admin: 'Admin'
};

export const ROLE_SHORT: Record<Role, string> = {
  pkm_operator: 'PKM',
  pkm_pj: 'PJ',
  dinas_verifier: 'Dinas',
  kabid: 'Kabid',
  sekretaris: 'Sekretaris',
  iku_owner: 'Owner',
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

export type Tone = 'neutral' | 'info' | 'success' | 'warning' | 'danger' | 'final';

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
  MATCHED: 'Matched',
  MISMATCH: 'Mismatch',
  UNKNOWN_IKU: 'IKU tidak dikenal',
  UNKNOWN_INDICATOR: 'Indikator tidak dikenal',
  IKU_ROLLUP_ONLY: 'IKU induk / rollup',
  AMBIGUOUS: 'Ambigu',
  DUPLICATE: 'Duplikat',
  INVALID_VALUE: 'Nilai tidak valid',
  SCOPE_DENIED: 'Scope ditolak'
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
  NO_DATA: 'Belum ada data',
  STALE: 'Perlu hitung ulang'
};

export const DISPLAY_STATUS_TONE: Record<string, Tone> = {
  ON_TARGET: 'success',
  BELOW_TARGET: 'warning',
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
  VERIFY_DINAS: 'Verifikasi Dinas',
  VALIDATE_KABID: 'Validasi Kabid',
  FINAL: 'Final',
  SPECIAL_CHANGE: 'Perubahan Khusus'
};

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
