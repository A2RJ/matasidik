PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version INTEGER PRIMARY KEY,
  applied_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  nama TEXT NOT NULL,
  role TEXT NOT NULL,
  faskes_id TEXT,
  bidang_id TEXT,
  aktif INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS wilayah (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL UNIQUE,
  nama TEXT NOT NULL,
  level TEXT NOT NULL CHECK (level IN ('desa', 'kecamatan', 'kabupaten')),
  parent_id TEXT REFERENCES wilayah(id)
);

CREATE TABLE IF NOT EXISTS faskes (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL UNIQUE,
  nama TEXT NOT NULL,
  wilayah_id TEXT REFERENCES wilayah(id),
  jenis TEXT NOT NULL DEFAULT 'puskesmas',
  aktif INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS bidang (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL UNIQUE,
  nama TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS periode (
  id TEXT PRIMARY KEY,
  tahun INTEGER NOT NULL,
  bulan INTEGER NOT NULL CHECK (bulan BETWEEN 1 AND 12),
  kode TEXT NOT NULL UNIQUE,
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'locked')),
  is_late_allowed INTEGER NOT NULL DEFAULT 0,
  UNIQUE (tahun, bulan)
);

CREATE TABLE IF NOT EXISTS master_release (
  id TEXT PRIMARY KEY,
  tahun INTEGER NOT NULL,
  status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
  published_at TEXT,
  published_by TEXT REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS label (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL,
  nama TEXT NOT NULL,
  tipe TEXT NOT NULL,
  parent_id TEXT REFERENCES label(id),
  master_release_id TEXT NOT NULL REFERENCES master_release(id),
  aggregation_method TEXT NOT NULL DEFAULT 'RATIO_OF_SUMS'
    CHECK (aggregation_method IN ('RATIO_OF_SUMS', 'SUM', 'WEIGHTED_AVERAGE', 'FORMULA', 'NONE')),
  allow_direct_input INTEGER NOT NULL DEFAULT 1,
  is_rollup_only INTEGER NOT NULL DEFAULT 0,
  pj_bidang_id TEXT REFERENCES bidang(id),
  iku_owner_id TEXT REFERENCES users(id),
  level INTEGER NOT NULL DEFAULT 1,
  path TEXT,
  aktif INTEGER NOT NULL DEFAULT 1,
  berlaku_dari TEXT,
  berlaku_sampai TEXT,
  UNIQUE (kode, master_release_id)
);

CREATE TABLE IF NOT EXISTS indikator (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL UNIQUE,
  nama TEXT NOT NULL,
  satuan TEXT NOT NULL,
  definisi TEXT,
  numerator_label TEXT,
  denominator_label TEXT,
  sifat_input TEXT NOT NULL DEFAULT 'puskesmas',
  frekuensi TEXT NOT NULL DEFAULT 'bulanan',
  pj_bidang_id TEXT REFERENCES bidang(id),
  pj_timker TEXT,
  aktif INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS indikator_label (
  indikator_id TEXT NOT NULL REFERENCES indikator(id),
  label_id TEXT NOT NULL REFERENCES label(id),
  bobot REAL NOT NULL DEFAULT 1,
  is_primary INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (indikator_id, label_id)
);

CREATE TABLE IF NOT EXISTS iku_komponen (
  id TEXT PRIMARY KEY,
  parent_label_id TEXT NOT NULL REFERENCES label(id),
  tipe_komponen TEXT NOT NULL CHECK (tipe_komponen IN ('INDIKATOR', 'CHILD_IKU')),
  komponen_id TEXT NOT NULL,
  bobot REAL NOT NULL DEFAULT 1,
  sort INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS target (
  id TEXT PRIMARY KEY,
  label_id TEXT REFERENCES label(id),
  indikator_id TEXT REFERENCES indikator(id),
  periode_id TEXT NOT NULL REFERENCES periode(id),
  scope_type TEXT NOT NULL DEFAULT 'KABUPATEN',
  scope_id TEXT,
  target_tahunan REAL NOT NULL,
  target_bulanan REAL NOT NULL,
  UNIQUE (label_id, indikator_id, periode_id, scope_type, scope_id)
);

CREATE TABLE IF NOT EXISTS workflow_version (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL,
  version_no INTEGER NOT NULL,
  steps_json TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'published',
  published_at TEXT
);

CREATE TABLE IF NOT EXISTS file_object (
  id TEXT PRIMARY KEY,
  original_name TEXT NOT NULL,
  sha256 TEXT NOT NULL,
  path TEXT NOT NULL,
  mime TEXT,
  size INTEGER NOT NULL,
  uploaded_by TEXT REFERENCES users(id),
  scan_status TEXT NOT NULL DEFAULT 'pending',
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS submission (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL UNIQUE,
  tipe TEXT NOT NULL DEFAULT 'BNBA',
  uploader_id TEXT NOT NULL REFERENCES users(id),
  faskes_id TEXT NOT NULL REFERENCES faskes(id),
  periode_id TEXT NOT NULL REFERENCES periode(id),
  file_id TEXT NOT NULL REFERENCES file_object(id),
  file_hash TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'draf',
  step TEXT NOT NULL DEFAULT 'DRAFT_PKM',
  master_release_id TEXT REFERENCES master_release(id),
  workflow_version_id TEXT REFERENCES workflow_version(id),
  idempotency_key TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_submission_idem
  ON submission(idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE INDEX IF NOT EXISTS ix_submission_period ON submission(periode_id, status);
CREATE INDEX IF NOT EXISTS ix_submission_faskes ON submission(faskes_id, periode_id);

CREATE TABLE IF NOT EXISTS submission_revision (
  id TEXT PRIMARY KEY,
  submission_id TEXT NOT NULL REFERENCES submission(id),
  revision_no INTEGER NOT NULL,
  parent_revision_id TEXT REFERENCES submission_revision(id),
  file_id TEXT NOT NULL REFERENCES file_object(id),
  created_by TEXT REFERENCES users(id),
  created_at TEXT NOT NULL,
  UNIQUE (submission_id, revision_no)
);

CREATE TABLE IF NOT EXISTS submission_group (
  id TEXT PRIMARY KEY,
  submission_id TEXT NOT NULL REFERENCES submission(id),
  owner_bidang_id TEXT REFERENCES bidang(id),
  primary_label_id TEXT REFERENCES label(id),
  status TEXT NOT NULL DEFAULT 'active',
  current_step TEXT NOT NULL DEFAULT 'DRAFT_PKM'
);

CREATE TABLE IF NOT EXISTS submission_row (
  id TEXT PRIMARY KEY,
  submission_id TEXT NOT NULL REFERENCES submission(id),
  group_id TEXT REFERENCES submission_group(id),
  row_no INTEGER NOT NULL,
  kode_iku_raw TEXT,
  kode_indikator_raw TEXT,
  matched_indikator_id TEXT REFERENCES indikator(id),
  matched_label_id TEXT REFERENCES label(id),
  scope_type TEXT NOT NULL DEFAULT 'FASKES',
  scope_code TEXT,
  nilai REAL,
  numerator REAL,
  denominator REAL,
  classification_status TEXT NOT NULL,
  error_message TEXT,
  created_at TEXT NOT NULL,
  UNIQUE (submission_id, row_no)
);
CREATE INDEX IF NOT EXISTS ix_row_submission ON submission_row(submission_id, classification_status);

CREATE TABLE IF NOT EXISTS submission_review (
  id TEXT PRIMARY KEY,
  submission_id TEXT NOT NULL REFERENCES submission(id),
  group_id TEXT REFERENCES submission_group(id),
  step TEXT NOT NULL,
  reviewer_id TEXT REFERENCES users(id),
  action TEXT NOT NULL,
  notes TEXT,
  details_json TEXT,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS job (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  payload_json TEXT NOT NULL,
  dedupe_key TEXT,
  status TEXT NOT NULL DEFAULT 'queued',
  attempts INTEGER NOT NULL DEFAULT 0,
  max_attempts INTEGER NOT NULL DEFAULT 5,
  run_after TEXT NOT NULL,
  locked_by TEXT,
  locked_at TEXT,
  heartbeat_at TEXT,
  last_error TEXT,
  progress_json TEXT,
  created_at TEXT NOT NULL,
  finished_at TEXT
);
CREATE INDEX IF NOT EXISTS ix_job_claim ON job(status, run_after, created_at);

CREATE TABLE IF NOT EXISTS indicator_value (
  id TEXT PRIMARY KEY,
  submission_row_id TEXT NOT NULL REFERENCES submission_row(id),
  indikator_id TEXT NOT NULL REFERENCES indikator(id),
  periode_id TEXT NOT NULL REFERENCES periode(id),
  faskes_id TEXT REFERENCES faskes(id),
  scope_type TEXT NOT NULL DEFAULT 'FASKES',
  scope_id TEXT,
  numerator REAL NOT NULL,
  denominator REAL NOT NULL,
  value REAL NOT NULL,
  status TEXT NOT NULL DEFAULT 'COMPLETE',
  data_version TEXT,
  created_at TEXT NOT NULL,
  UNIQUE (submission_row_id)
);
CREATE INDEX IF NOT EXISTS ix_indicator_value_period ON indicator_value(periode_id, scope_type, scope_id);

CREATE TABLE IF NOT EXISTS calculation_run (
  id TEXT PRIMARY KEY,
  periode_id TEXT NOT NULL REFERENCES periode(id),
  scope_type TEXT NOT NULL,
  scope_id TEXT,
  trigger TEXT NOT NULL,
  status TEXT NOT NULL,
  started_at TEXT NOT NULL,
  finished_at TEXT,
  error TEXT
);

CREATE TABLE IF NOT EXISTS iku_value (
  id TEXT PRIMARY KEY,
  label_id TEXT NOT NULL REFERENCES label(id),
  periode_id TEXT NOT NULL REFERENCES periode(id),
  scope_type TEXT NOT NULL,
  scope_id TEXT,
  numerator REAL NOT NULL DEFAULT 0,
  denominator REAL NOT NULL DEFAULT 0,
  value REAL,
  status TEXT NOT NULL DEFAULT 'NO_DATA',
  formula_version TEXT NOT NULL DEFAULT 'v1',
  calculation_run_id TEXT REFERENCES calculation_run(id),
  source_hash TEXT,
  data_version TEXT,
  updated_at TEXT NOT NULL,
  UNIQUE (label_id, periode_id, scope_type, scope_id)
);
CREATE INDEX IF NOT EXISTS ix_iku_value_period ON iku_value(periode_id, scope_type, scope_id);

CREATE TABLE IF NOT EXISTS iku_value_detail (
  id TEXT PRIMARY KEY,
  iku_value_id TEXT NOT NULL REFERENCES iku_value(id),
  component_type TEXT NOT NULL,
  component_id TEXT NOT NULL,
  component_label_id TEXT REFERENCES label(id),
  component_value REAL,
  numerator REAL,
  denominator REAL,
  weight REAL NOT NULL DEFAULT 1,
  contribution REAL
);
CREATE INDEX IF NOT EXISTS ix_detail_value ON iku_value_detail(iku_value_id);

CREATE TABLE IF NOT EXISTS dashboard_row (
  id TEXT PRIMARY KEY,
  kode TEXT NOT NULL UNIQUE,
  label_id TEXT REFERENCES label(id),
  indikator_id TEXT REFERENCES indikator(id),
  urutan INTEGER NOT NULL,
  target REAL NOT NULL DEFAULT 100,
  baseline_json TEXT,
  format_json TEXT
);

CREATE TABLE IF NOT EXISTS faktor (
  id TEXT PRIMARY KEY,
  iku_value_id TEXT NOT NULL REFERENCES iku_value(id),
  periode_id TEXT NOT NULL REFERENCES periode(id),
  pendukung TEXT,
  penghambat TEXT,
  created_by TEXT REFERENCES users(id),
  updated_at TEXT NOT NULL,
  UNIQUE (iku_value_id)
);

CREATE TABLE IF NOT EXISTS tindak_lanjut (
  id TEXT PRIMARY KEY,
  iku_value_id TEXT REFERENCES iku_value(id),
  periode_id TEXT NOT NULL REFERENCES periode(id),
  masalah TEXT NOT NULL,
  rekomendasi TEXT NOT NULL,
  pj_user_id TEXT REFERENCES users(id),
  status TEXT NOT NULL DEFAULT 'open'
    CHECK (status IN ('open', 'in_progress', 'done', 'cancelled')),
  created_by TEXT REFERENCES users(id),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS audit_log (
  id TEXT PRIMARY KEY,
  actor_id TEXT REFERENCES users(id),
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id TEXT,
  before_json TEXT,
  after_json TEXT,
  reason TEXT,
  request_id TEXT,
  created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS ix_audit_entity ON audit_log(entity_type, entity_id, created_at);
