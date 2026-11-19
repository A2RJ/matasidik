INSERT OR IGNORE INTO positions
  (id, kode, nama, role, level, bidang_id, faskes_id, aktif, created_at)
VALUES
  ('pos-kepala-pkm-taliwang', 'KEPALA-PKM-TALIWANG', 'Kepala Puskesmas Taliwang', 'pkm_approver', 2, NULL, 'faskes-pkm-taliwang', 1, '2026-11-19T00:00:00Z');

INSERT OR IGNORE INTO users
  (id, nama, role, faskes_id, bidang_id, aktif, created_at)
VALUES
  ('usr-kepala-pkm', 'Kepala PKM Taliwang', 'pkm_approver', 'faskes-pkm-taliwang', NULL, 1, '2026-11-19T00:00:00Z');

UPDATE users
SET username = 'kepalapkm',
    password_hash = '$2a$12$ZyNszKJDkWth/Q5bv8xBJ..ZV6LCXbIrsaS4TkRs3TWa7iR8PKCZy',
    position_id = 'pos-kepala-pkm-taliwang'
WHERE id = 'usr-kepala-pkm';

INSERT OR IGNORE INTO workflow_version
  (id, kode, version_no, steps_json, status, published_at)
VALUES
  ('wf-default-v2', 'default_pkm_approval', 2,
   '[{"kode":"DRAFT_PKM","urutan":1,"role":"pkm_operator"},{"kode":"APPROVE_PKM","urutan":2,"role":"pkm_approver"},{"kode":"VERIFY_DINAS","urutan":3,"role":"dinas_verifier"},{"kode":"VALIDATE_KABID","urutan":4,"role":"kabid"},{"kode":"FINAL","urutan":5,"role":"sistem"}]',
   'published', '2026-11-19T00:00:00Z');

CREATE TABLE IF NOT EXISTS iku_change_request (
  id TEXT PRIMARY KEY,
  action TEXT NOT NULL CHECK (action IN ('CREATE', 'UPDATE', 'DEACTIVATE')),
  label_id TEXT REFERENCES label(id),
  requested_by TEXT NOT NULL REFERENCES users(id),
  faskes_id TEXT NOT NULL REFERENCES faskes(id),
  master_release_id TEXT NOT NULL REFERENCES master_release(id),
  payload_json TEXT NOT NULL,
  current_revision_no INTEGER NOT NULL DEFAULT 1,
  status TEXT NOT NULL CHECK (status IN (
    'draf', 'menunggu_puskesmas', 'dikembalikan_puskesmas',
    'menunggu_admin', 'dikembalikan_admin', 'disetujui', 'ditolak', 'dibatalkan'
  )),
  applied_label_id TEXT REFERENCES label(id),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS ix_iku_change_status ON iku_change_request(status, created_at);
CREATE INDEX IF NOT EXISTS ix_iku_change_faskes ON iku_change_request(faskes_id, status);
CREATE INDEX IF NOT EXISTS ix_iku_change_requester ON iku_change_request(requested_by, created_at);

CREATE TABLE IF NOT EXISTS iku_change_revision (
  id TEXT PRIMARY KEY,
  request_id TEXT NOT NULL REFERENCES iku_change_request(id),
  revision_no INTEGER NOT NULL,
  payload_json TEXT NOT NULL,
  before_json TEXT,
  after_json TEXT,
  created_by TEXT NOT NULL REFERENCES users(id),
  created_at TEXT NOT NULL,
  UNIQUE (request_id, revision_no)
);
CREATE INDEX IF NOT EXISTS ix_iku_change_revision ON iku_change_revision(request_id, revision_no);

CREATE TABLE IF NOT EXISTS iku_change_review (
  id TEXT PRIMARY KEY,
  request_id TEXT NOT NULL REFERENCES iku_change_request(id),
  revision_no INTEGER NOT NULL,
  step TEXT NOT NULL CHECK (step IN ('PKM', 'ADMIN')),
  reviewer_id TEXT NOT NULL REFERENCES users(id),
  action TEXT NOT NULL CHECK (action IN ('approve', 'reject', 'return')),
  notes TEXT,
  created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS ix_iku_change_review ON iku_change_review(request_id, created_at);
