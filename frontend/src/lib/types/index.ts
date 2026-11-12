export type Role =
  | 'pkm_operator'
  | 'pkm_pj'
  | 'dinas_verifier'
  | 'kabid'
  | 'sekretaris'
  | 'iku_owner'
  | 'pimpinan'
  | 'admin';

export type SubmissionStatus =
  | 'draf'
  | 'dikirim'
  | 'diverifikasi'
  | 'dikembalikan_petugas'
  | 'dikembalikan_verifikator'
  | 'final'
  | 'diubah_khusus';

export type ClassificationStatus =
  | 'MATCHED'
  | 'MISMATCH'
  | 'UNKNOWN_IKU'
  | 'UNKNOWN_INDICATOR'
  | 'IKU_ROLLUP_ONLY'
  | 'AMBIGUOUS'
  | 'DUPLICATE'
  | 'INVALID_VALUE'
  | 'SCOPE_DENIED';

export interface User {
  id: string;
  nama: string;
  username: string;
  role: Role;
  faskes_id?: string;
  bidang_id?: string;
  position_id?: string;
  position_nama?: string;
  last_login_at?: string;
  aktif: boolean;
  created_at: string;
}

export interface Position {
  id: string;
  kode: string;
  nama: string;
  role: Role;
  level: number;
  bidang_id?: string;
  faskes_id?: string;
  aktif: boolean;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  access_expires_at: string;
  user: User;
  acting_as?: User;
}

export interface AdminOverview {
  users: User[];
  faskes: Faskes[];
  positions: Position[];
  bidang: Bidang[];
}

export interface Faskes {
  id: string;
  kode: string;
  nama: string;
  wilayah_id?: string;
  jenis: string;
  aktif: boolean;
}

export interface Bidang {
  id: string;
  kode: string;
  nama: string;
}

export interface Wilayah {
  id: string;
  kode: string;
  nama: string;
  level: string;
  parent_id?: string;
}

export interface Periode {
  id: string;
  tahun: number;
  bulan: number;
  kode: string;
  status: string;
  is_late_allowed: boolean;
}

export interface Label {
  id: string;
  kode: string;
  nama: string;
  tipe: string;
  parent_id?: string;
  master_release_id: string;
  aggregation_method: string;
  allow_direct_input: boolean;
  is_rollup_only: boolean;
  pj_bidang_id?: string;
  iku_owner_id?: string;
  level: number;
  path?: string;
  aktif: boolean;
}

export interface Indicator {
  id: string;
  kode: string;
  nama: string;
  satuan: string;
  definisi?: string;
  numerator_label?: string;
  denominator_label?: string;
  sifat_input: string;
  frekuensi: string;
  pj_bidang_id?: string;
  pj_timker?: string;
  aktif: boolean;
}

export interface Target {
  id: string;
  label_id?: string;
  indikator_id?: string;
  periode_id: string;
  scope_type: string;
  scope_id?: string;
  target_tahunan: number;
  target_bulanan: number;
}

export interface CapaianTahunan {
  id: string;
  label_id?: string;
  indikator_id?: string;
  tahun: number;
  scope_type: string;
  scope_id?: string;
  nilai: number;
  catatan?: string;
  updated_by?: string;
  created_at: string;
  updated_at: string;
  label_kode?: string;
  label_nama?: string;
  indikator_kode?: string;
  indikator_nama?: string;
  updated_by_nama?: string;
}

export interface DashboardRow {
  id: string;
  kode: string;
  label_id?: string;
  indikator_id?: string;
  urutan: number;
  target: number;
  baseline_json?: string;
}

export interface DashboardItem {
  kode: string;
  indikator: string;
  label_id?: string;
  indikator_id?: string;
  label_kode?: string;
  label_nama?: string;
  iku_value_id?: string;
  indicator_value_id?: string;
  calculation_run_id?: string;
  source_hash?: string;
  data_version?: string;
  target: number;
  baseline_prev?: number;
  baseline_year?: number;
  baseline_tw1?: number;
  numerator?: number;
  denominator?: number;
  value?: number;
  value_status: string;
  display_status: string;
  gap?: number;
  updated_at?: string;
}

export interface DashboardSummary {
  total: number;
  on_target: number;
  below_target: number;
  no_data: number;
  average: number;
  parent_value?: number;
  parent_status: string;
}

export interface SPMFaskesDetail {
  faskes_id: string;
  faskes_kode: string;
  faskes_nama: string;
  uploaded: boolean;
  submission_id?: string;
  submission_kode?: string;
  file_id?: string;
  file_name?: string;
  uploader_nama?: string;
  numerator?: number;
  denominator?: number;
  value?: number;
  status: string;
  updated_at?: string;
}

export interface SPMDetail {
  row: DashboardItem;
  faskes: SPMFaskesDetail[];
  provenance?: IKUValueDetail[];
}

export interface ChartFaskes {
  id: string;
  kode: string;
  nama: string;
}

export interface ChartSeries {
  kode: string;
  nama: string;
  values: (number | null)[];
}

export interface IKUChart {
  faskes: ChartFaskes[];
  series: ChartSeries[];
}

export interface Bootstrap {
  users: User[];
  faskes: Faskes[];
  wilayah: Wilayah[];
  bidang: Bidang[];
  periods: Periode[];
  default_period_id?: string;
  labels: Label[];
  indicators: Indicator[];
  targets: Target[];
  dashboard_rows: DashboardRow[];
  roles: Role[];
  submission_statuses: SubmissionStatus[];
  classifications: ClassificationStatus[];
  value_statuses: string[];
  tl_statuses: string[];
}

export interface Submission {
  id: string;
  kode: string;
  tipe: string;
  uploader_id: string;
  faskes_id: string;
  periode_id: string;
  file_id: string;
  file_hash: string;
  status: SubmissionStatus;
  step: string;
  master_release_id?: string;
  workflow_version_id?: string;
  created_at: string;
  updated_at: string;
  uploader_nama?: string;
  faskes_nama?: string;
  periode_kode?: string;
}

export interface SubmissionRow {
  id: string;
  submission_id: string;
  row_no: number;
  kode_iku_raw?: string;
  kode_indikator_raw?: string;
  matched_indikator_id?: string;
  matched_label_id?: string;
  scope_type: string;
  scope_code?: string;
  nilai?: number;
  numerator?: number;
  denominator?: number;
  classification_status: ClassificationStatus;
  error_message?: string;
  indikator_kode?: string;
  indikator_nama?: string;
  iku_kode?: string;
  iku_nama?: string;
}

export interface SubmissionReview {
  id: string;
  step: string;
  reviewer_id?: string;
  reviewer_nama?: string;
  action: string;
  notes?: string;
  created_at: string;
}

export interface SubmissionRevision {
  id: string;
  revision_no: number;
  file_id: string;
  created_at: string;
}

export interface SubmissionDetail {
  submission: Submission;
  rows: SubmissionRow[];
  reviews: SubmissionReview[] | null;
  revisions: SubmissionRevision[];
  allowed_iku: Record<string, Label[]>;
  blockers: number;
}

export interface Job {
  id: string;
  type: string;
  status: string;
  attempts: number;
  progress_json?: string;
  last_error?: string;
  created_at: string;
  finished_at?: string;
}

export interface IKUValue {
  id: string;
  label_id: string;
  periode_id: string;
  scope_type: string;
  scope_id?: string;
  numerator: number;
  denominator: number;
  value?: number;
  status: string;
  formula_version: string;
  calculation_run_id?: string;
  source_hash?: string;
  data_version?: string;
  updated_at: string;
  label_kode?: string;
  label_nama?: string;
  label_tipe?: string;
}

export interface IKUValueDetail {
  id: string;
  iku_value_id: string;
  component_type: string;
  component_id: string;
  component_label_id?: string;
  component_value?: number;
  numerator?: number;
  denominator?: number;
  weight: number;
  contribution?: number;
  label_kode?: string;
  label_nama?: string;
  indikator_kode?: string;
  indikator_nama?: string;
}

export interface SPMTrace {
  row: DashboardItem;
  iku_value?: IKUValue;
  details?: IKUValueDetail[];
  indicator_value?: {
    id: string;
    submission_row_id: string;
    indikator_id: string;
    numerator: number;
    denominator: number;
    value: number;
  };
  submission_row?: SubmissionRow;
  submission?: Submission;
  file?: {
    id: string;
    original_name: string;
    sha256: string;
    size: number;
    created_at: string;
  };
  reviews?: SubmissionReview[];
}

export interface Factor {
  id: string;
  iku_value_id: string;
  periode_id: string;
  pendukung?: string;
  penghambat?: string;
  created_by?: string;
  updated_at: string;
}

export interface FactorRow {
  iku_value: IKUValue;
  faktor?: Factor;
}

export interface TindakLanjut {
  id: string;
  iku_value_id?: string;
  periode_id: string;
  masalah: string;
  rekomendasi: string;
  pj_user_id?: string;
  status: string;
  created_by?: string;
  created_at: string;
  updated_at: string;
  pj_nama?: string;
  label_kode?: string;
  label_nama?: string;
}

export interface AuditLog {
  id: string;
  actor_id?: string;
  actor_nama?: string;
  action: string;
  entity_type: string;
  entity_id?: string;
  reason?: string;
  request_id?: string;
  impersonator_id?: string;
  impersonator_nama?: string;
  created_at: string;
}
