package sqlite

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

//go:embed testdata/bnba-2026-05.csv
var demoFS embed.FS

// SeedDemo membuat satu submission final awal Mei 2026 agar dashboard
// langsung dapat dieksplorasi. File asli tetap memakai kode IKU induk untuk R3;
// koreksi ke .03 dilakukan di preview dan disimpan pada SUBMISSION_ROW.
func (s *Store) SeedDemo(ctx context.Context) (bool, error) {
	var count int
	if err := s.db.GetContext(ctx, &count, `SELECT COUNT(1) FROM submission WHERE kode = 'S-2026-05-001'`); err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}
	content, err := demoFS.ReadFile("testdata/bnba-2026-05.csv")
	if err != nil {
		return false, err
	}
	sum := sha256.Sum256(content)
	hash := hex.EncodeToString(sum[:])
	relPath := filepath.Join("2026", "05", hash+".csv")
	absPath := filepath.Join(s.dataDir, relPath)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return false, err
	}
	if err := os.WriteFile(absPath, content, 0o644); err != nil {
		return false, err
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	fileID := uuid.NewString()
	if _, err := tx.ExecContext(ctx, `INSERT INTO file_object
		(id, original_name, sha256, path, mime, size, uploaded_by, scan_status, created_at)
		VALUES (?, 'BNBA_2026-05.csv', ?, ?, 'text/csv', ?, 'usr-siti', 'clean', '2026-05-10T08:00:00Z')`,
		fileID, hash, absPath, len(content)); err != nil {
		return false, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO submission
		(id, kode, tipe, uploader_id, faskes_id, periode_id, file_id, file_hash, status, step,
		 master_release_id, workflow_version_id, created_at, updated_at)
		VALUES ('sub-demo-001', 'S-2026-05-001', 'BNBA', 'usr-siti', 'faskes-pkm-taliwang',
		 'per-2026-05', ?, ?, 'final', 'FINAL', 'mr-2026', 'wf-default-v1',
		 '2026-05-10T08:00:00Z', '2026-05-10T10:31:00Z')`,
		fileID, hash); err != nil {
		return false, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO submission_revision
		(id, submission_id, revision_no, parent_revision_id, file_id, created_by, created_at)
		VALUES ('rev-demo-001', 'sub-demo-001', 1, NULL, ?, 'usr-siti', '2026-05-10T08:00:00Z')`,
		fileID); err != nil {
		return false, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO submission_group
		(id, submission_id, owner_bidang_id, primary_label_id, status, current_step)
		VALUES ('grp-demo-001', 'sub-demo-001', 'bid-kesmas', 'lbl-iku-001', 'active', 'FINAL')`); err != nil {
		return false, err
	}

	rows := []struct {
		id, rawIKU, rawInd, indID, labelID string
		rowNo                              int
		num, den                           float64
		errMsg                             string
	}{
		{"row-r1", "IKU-2026-001.01", "IND-001", "ind-001", "lbl-iku-001-01", 1, 450, 500, ""},
		{"row-r2", "IKU-2026-001.02", "IND-002", "ind-002", "lbl-iku-001-02", 2, 380, 400, ""},
		{"row-r3", "IKU-2026-001", "IND-003", "ind-003", "lbl-iku-001-03", 3, 20, 20, "dikoreksi dari IKU induk ke .03 di preview"},
	}
	for _, r := range rows {
		if _, err := tx.ExecContext(ctx, `INSERT INTO submission_row
			(id, submission_id, group_id, row_no, kode_iku_raw, kode_indikator_raw,
			 matched_indikator_id, matched_label_id, scope_type, scope_code, nilai, numerator, denominator,
			 classification_status, error_message, created_at)
			VALUES (?, 'sub-demo-001', 'grp-demo-001', ?, ?, ?, ?, ?, 'FASKES', 'faskes-pkm-taliwang',
			 ?, ?, ?, 'MATCHED', ?, '2026-05-10T08:01:00Z')`,
			r.id, r.rowNo, r.rawIKU, r.rawInd, r.indID, r.labelID,
			(r.num/r.den)*100, r.num, r.den, nullString(&r.errMsg)); err != nil {
			return false, err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO indicator_value
			(id, submission_row_id, indikator_id, periode_id, faskes_id, scope_type, scope_id,
			 numerator, denominator, value, status, data_version, created_at)
			VALUES (?, ?, ?, 'per-2026-05', 'faskes-pkm-taliwang', 'FASKES', 'faskes-pkm-taliwang',
			 ?, ?, ?, 'COMPLETE', 'v2026-05-001', '2026-05-10T08:01:00Z')`,
			uuid.NewString(), r.id, r.indID, r.num, r.den, (r.num/r.den)*100); err != nil {
			return false, err
		}
	}

	reviews := []struct {
		id, step, reviewer, action, notes, at string
	}{
		{"rv-demo-submit", "DRAFT_PKM", "usr-siti", "submit", "Data Mei 2026 · 3 baris. R3 dikoreksi ke IKU .03.", "2026-05-10T10:02:00Z"},
		{"rv-demo-verify", "VERIFY_DINAS", "usr-rina", "approve", "Klasifikasi sesuai: R1 .01 450/500; R2 .02 380/400; R3 .03 20/20. Lanjut ke Kabid.", "2026-05-10T10:15:00Z"},
		{"rv-demo-final", "VALIDATE_KABID", "usr-kabid", "approve", "Diterima. Catatan penerimaan: data final, lanjut kalkulasi dan publikasi IKU.", "2026-05-10T10:30:00Z"},
	}
	for _, r := range reviews {
		if _, err := tx.ExecContext(ctx, `INSERT INTO submission_review
			(id, submission_id, group_id, step, reviewer_id, action, notes, details_json, created_at)
			VALUES (?, 'sub-demo-001', 'grp-demo-001', ?, ?, ?, ?, NULL, ?)`,
			r.id, r.step, r.reviewer, r.action, r.notes, r.at); err != nil {
			return false, err
		}
	}

	audits := []struct {
		id, actor, action, reason, at string
	}{
		{"aud-demo-1", "usr-siti", "submission.created", "Upload BNBA Mei 2026", "2026-05-10T08:00:00Z"},
		{"aud-demo-2", "usr-siti", "submission.row_reclassified", "R3 diubah dari IKU induk ke .03", "2026-05-10T10:01:00Z"},
		{"aud-demo-3", "usr-siti", "submission.submitted", "Kirim ke verifikasi Dinas", "2026-05-10T10:02:00Z"},
		{"aud-demo-4", "usr-rina", "submission.verified", "Verifikasi disetujui", "2026-05-10T10:15:00Z"},
		{"aud-demo-5", "usr-kabid", "submission.finalized", "Validasi final disetujui", "2026-05-10T10:30:00Z"},
	}
	for _, a := range audits {
		if _, err := tx.ExecContext(ctx, `INSERT INTO audit_log
			(id, actor_id, action, entity_type, entity_id, before_json, after_json, reason, request_id, created_at)
			VALUES (?, ?, ?, 'submission', 'sub-demo-001', NULL, NULL, ?, 'seed', ?)`,
			a.id, a.actor, a.action, a.reason, a.at); err != nil {
			return false, err
		}
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}
