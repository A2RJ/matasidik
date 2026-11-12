package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
)

const submissionSelect = `SELECT s.id, s.kode, s.tipe, s.uploader_id, s.faskes_id, s.periode_id,
	s.file_id, s.file_hash, s.status, s.step, s.master_release_id, s.workflow_version_id,
	s.idempotency_key, s.created_at, s.updated_at,
	u.nama AS uploader_nama, f.nama AS faskes_nama, p.kode AS periode_kode
	FROM submission s
	JOIN users u ON u.id = s.uploader_id
	JOIN faskes f ON f.id = s.faskes_id
	JOIN periode p ON p.id = s.periode_id`

func (s *Store) FindIdempotentSubmission(ctx context.Context, key string) (*domain.Submission, error) {
	if strings.TrimSpace(key) == "" {
		return nil, nil
	}
	var out domain.Submission
	err := s.db.GetContext(ctx, &out, submissionSelect+` WHERE s.idempotency_key = ?`, key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *Store) CreateUpload(ctx context.Context, p domain.CreateUploadParams) (*domain.Submission, *domain.Job, error) {
	if p.SubmissionID == "" {
		p.SubmissionID = uuid.NewString()
	}
	if p.File.ID == "" {
		p.File.ID = uuid.NewString()
	}
	jobID := uuid.NewString()
	if p.Kode == "" {
		var count int
		if err := s.db.GetContext(ctx, &count, `SELECT COUNT(1) FROM submission WHERE periode_id = ?`, p.PeriodeID); err != nil {
			return nil, nil, err
		}
		p.Kode = fmt.Sprintf("S-%s-%03d", p.PeriodKode, count+1)
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `INSERT INTO file_object
		(id, original_name, sha256, path, mime, size, uploaded_by, scan_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.File.ID, p.File.OriginalName, p.File.SHA256, p.File.Path, p.File.Mime,
		p.File.Size, p.UploaderID, "clean", p.Now); err != nil {
		return nil, nil, fmt.Errorf("simpan file: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO submission
		(id, kode, tipe, uploader_id, faskes_id, periode_id, file_id, file_hash, status, step,
		 master_release_id, workflow_version_id, idempotency_key, created_at, updated_at)
		VALUES (?, ?, 'BNBA', ?, ?, ?, ?, ?, 'draf', 'DRAFT_PKM', ?, ?, ?, ?, ?)`,
		p.SubmissionID, p.Kode, p.UploaderID, p.FaskesID, p.PeriodeID, p.File.ID, p.File.SHA256,
		p.MasterReleaseID, p.WorkflowVersionID, nullString(&p.IdempotencyKey), p.Now, p.Now); err != nil {
		return nil, nil, fmt.Errorf("simpan submission: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO submission_revision
		(id, submission_id, revision_no, parent_revision_id, file_id, created_by, created_at)
		VALUES (?, ?, 1, NULL, ?, ?, ?)`,
		uuid.NewString(), p.SubmissionID, p.File.ID, p.UploaderID, p.Now); err != nil {
		return nil, nil, fmt.Errorf("simpan revisi: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO submission_group
		(id, submission_id, owner_bidang_id, primary_label_id, status, current_step)
		VALUES (?, ?, NULL, NULL, 'active', 'DRAFT_PKM')`,
		uuid.NewString(), p.SubmissionID); err != nil {
		return nil, nil, fmt.Errorf("simpan group: %w", err)
	}
	payload, _ := json.Marshal(map[string]string{"submission_id": p.SubmissionID, "file_id": p.File.ID})
	if _, err := tx.ExecContext(ctx, `INSERT INTO job
		(id, type, payload_json, status, attempts, max_attempts, run_after, progress_json, created_at)
		VALUES (?, 'parse_submission', ?, 'queued', 0, 5, ?, ?, ?)`,
		jobID, string(payload), p.Now, `{"step":"queued","percent":0}`, p.Now); err != nil {
		return nil, nil, fmt.Errorf("simpan job: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}
	submission, err := s.GetSubmission(ctx, p.SubmissionID)
	if err != nil {
		return nil, nil, err
	}
	job, err := s.GetJob(ctx, jobID)
	if err != nil {
		return nil, nil, err
	}
	return submission, job, nil
}

func (s *Store) GetSubmission(ctx context.Context, id string) (*domain.Submission, error) {
	var out domain.Submission
	err := s.db.GetContext(ctx, &out, submissionSelect+` WHERE s.id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("submission")
	}
	return &out, err
}

func (s *Store) GetSubmissionByKode(ctx context.Context, kode string) (*domain.Submission, error) {
	var out domain.Submission
	err := s.db.GetContext(ctx, &out, submissionSelect+` WHERE s.kode = ?`, kode)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("submission")
	}
	return &out, err
}

func (s *Store) GetSubmissionByFileID(ctx context.Context, fileID string) (*domain.Submission, error) {
	var out domain.Submission
	err := s.db.GetContext(ctx, &out, submissionSelect+` WHERE s.file_id = ? ORDER BY s.created_at DESC LIMIT 1`, fileID)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("submission")
	}
	return &out, err
}

func (s *Store) GetFileObject(ctx context.Context, id string) (*domain.FileObject, error) {
	var out domain.FileObject
	err := s.db.GetContext(ctx, &out, `SELECT id, original_name, sha256, path, mime, size,
		uploaded_by, scan_status, created_at FROM file_object WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("file")
	}
	return &out, err
}

func (s *Store) ListSubmissions(ctx context.Context, f domain.SubmissionFilter) ([]domain.Submission, error) {
	q := submissionSelect + ` WHERE 1=1`
	args := []any{}
	if f.Status != "" {
		q += ` AND s.status = ?`
		args = append(args, f.Status)
	}
	if f.PeriodeID != "" {
		q += ` AND s.periode_id = ?`
		args = append(args, f.PeriodeID)
	}
	if f.FaskesID != "" {
		q += ` AND s.faskes_id = ?`
		args = append(args, f.FaskesID)
	}
	q += ` ORDER BY s.created_at DESC`
	if f.Limit > 0 {
		q += ` LIMIT ?`
		args = append(args, f.Limit)
	}
	var out []domain.Submission
	err := s.db.SelectContext(ctx, &out, q, args...)
	return out, err
}

func (s *Store) GetRows(ctx context.Context, submissionID string) ([]domain.SubmissionRow, error) {
	var out []domain.SubmissionRow
	err := s.db.SelectContext(ctx, &out, `SELECT sr.id, sr.submission_id, sr.group_id, sr.row_no,
		sr.kode_iku_raw, sr.kode_indikator_raw, sr.matched_indikator_id, sr.matched_label_id,
		sr.scope_type, sr.scope_code, sr.nilai, sr.numerator, sr.denominator,
		sr.classification_status, sr.error_message, sr.created_at,
		i.kode AS indikator_kode, i.nama AS indikator_nama,
		l.kode AS iku_kode, l.nama AS iku_nama
		FROM submission_row sr
		LEFT JOIN indikator i ON i.id = sr.matched_indikator_id
		LEFT JOIN label l ON l.id = sr.matched_label_id
		WHERE sr.submission_id = ? ORDER BY sr.row_no`, submissionID)
	return out, err
}

func (s *Store) GetRow(ctx context.Context, rowID string) (*domain.SubmissionRow, error) {
	var out domain.SubmissionRow
	err := s.db.GetContext(ctx, &out, `SELECT sr.id, sr.submission_id, sr.group_id, sr.row_no,
		sr.kode_iku_raw, sr.kode_indikator_raw, sr.matched_indikator_id, sr.matched_label_id,
		sr.scope_type, sr.scope_code, sr.nilai, sr.numerator, sr.denominator,
		sr.classification_status, sr.error_message, sr.created_at,
		i.kode AS indikator_kode, i.nama AS indikator_nama,
		l.kode AS iku_kode, l.nama AS iku_nama
		FROM submission_row sr
		LEFT JOIN indikator i ON i.id = sr.matched_indikator_id
		LEFT JOIN label l ON l.id = sr.matched_label_id
		WHERE sr.id = ?`, rowID)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("baris submission")
	}
	return &out, err
}

func (s *Store) UpdateRowClassification(ctx context.Context, rowID, labelID string, status domain.ClassificationStatus, errMsg string) error {
	var label any
	if strings.TrimSpace(labelID) != "" {
		label = labelID
	}
	var message any
	if strings.TrimSpace(errMsg) != "" {
		message = errMsg
	}
	_, err := s.db.ExecContext(ctx, `UPDATE submission_row
		SET matched_label_id = ?, classification_status = ?, error_message = ?
		WHERE id = ?`, label, status, message, rowID)
	return err
}

func (s *Store) UpsertIndicatorValueForRow(ctx context.Context, row *domain.SubmissionRow, sub *domain.Submission) error {
	if row.MatchedIndikatorID == nil || row.Numerator == nil || row.Denominator == nil {
		return domain.Validation("baris belum memiliki indikator/numerator/denominator", nil)
	}
	value := 0.0
	if *row.Denominator > 0 {
		value = *row.Numerator / *row.Denominator * 100
	}
	var existingID string
	_ = s.db.GetContext(ctx, &existingID, `SELECT id FROM indicator_value WHERE submission_row_id = ?`, row.ID)
	if existingID == "" {
		_, err := s.db.ExecContext(ctx, `INSERT INTO indicator_value
			(id, submission_row_id, indikator_id, periode_id, faskes_id, scope_type, scope_id,
			 numerator, denominator, value, status, data_version, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'COMPLETE', 'preview', ?)`,
			uuid.NewString(), row.ID, *row.MatchedIndikatorID, sub.PeriodeID, sub.FaskesID,
			row.ScopeType, nullString(row.ScopeCode), *row.Numerator, *row.Denominator, value, nowRFC3339())
		return err
	}
	_, err := s.db.ExecContext(ctx, `UPDATE indicator_value SET indikator_id = ?, faskes_id = ?,
		scope_type = ?, scope_id = ?, numerator = ?, denominator = ?, value = ?, status = 'COMPLETE'
		WHERE id = ?`, *row.MatchedIndikatorID, sub.FaskesID, row.ScopeType, nullString(row.ScopeCode),
		*row.Numerator, *row.Denominator, value, existingID)
	return err
}

func (s *Store) UpdateSubmissionStatus(ctx context.Context, id string, status domain.SubmissionStatus, step domain.WorkflowStep) error {
	_, err := s.db.ExecContext(ctx, `UPDATE submission SET status = ?, step = ?, updated_at = ? WHERE id = ?`,
		status, step, nowRFC3339(), id)
	return err
}

func (s *Store) CreateReview(ctx context.Context, review *domain.SubmissionReview) error {
	if review.ID == "" {
		review.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO submission_review
		(id, submission_id, group_id, step, reviewer_id, action, notes, details_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		review.ID, review.SubmissionID, review.GroupID, review.Step, review.ReviewerID,
		review.Action, review.Notes, review.DetailsJSON, review.CreatedAt)
	return err
}

func (s *Store) ListReviews(ctx context.Context, submissionID string) ([]domain.SubmissionReview, error) {
	var out []domain.SubmissionReview
	err := s.db.SelectContext(ctx, &out, `SELECT r.id, r.submission_id, r.group_id, r.step,
		r.reviewer_id, r.action, r.notes, r.details_json, r.created_at,
		u.nama AS reviewer_nama
		FROM submission_review r
		LEFT JOIN users u ON u.id = r.reviewer_id
		WHERE r.submission_id = ? ORDER BY r.created_at`, submissionID)
	return out, err
}

func (s *Store) ListRevisions(ctx context.Context, submissionID string) ([]domain.SubmissionRevision, error) {
	var out []domain.SubmissionRevision
	err := s.db.SelectContext(ctx, &out, `SELECT id, submission_id, revision_no, parent_revision_id,
		file_id, created_by, created_at FROM submission_revision
		WHERE submission_id = ? ORDER BY revision_no`, submissionID)
	return out, err
}

func (s *Store) InsertSubmissionRows(ctx context.Context, submissionID string, rows []domain.SubmissionRowInput, values []domain.IndicatorValueInput) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM indicator_value WHERE submission_row_id IN
		(SELECT id FROM submission_row WHERE submission_id = ?)`, submissionID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM submission_row WHERE submission_id = ?`, submissionID); err != nil {
		return err
	}
	for _, row := range rows {
		id := row.ID
		if id == "" {
			id = uuid.NewString()
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO submission_row
			(id, submission_id, group_id, row_no, kode_iku_raw, kode_indikator_raw,
			 matched_indikator_id, matched_label_id, scope_type, scope_code, nilai, numerator, denominator,
			 classification_status, error_message, created_at)
			VALUES (?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, submissionID, row.RowNo, nullString(&row.KodeIKURaw), nullString(&row.KodeIndikatorRaw),
			nullString(&row.MatchedIndikatorID), nullString(&row.MatchedLabelID),
			row.ScopeType, nullString(&row.ScopeCode), row.Value, row.Numerator, row.Denominator,
			row.Classification, nullString(&row.ErrorMessage), nowRFC3339()); err != nil {
			return fmt.Errorf("simpan baris %d: %w", row.RowNo, err)
		}
	}
	for _, v := range values {
		if _, err := tx.ExecContext(ctx, `INSERT INTO indicator_value
			(id, submission_row_id, indikator_id, periode_id, faskes_id, scope_type, scope_id,
			 numerator, denominator, value, status, data_version, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			uuid.NewString(), v.SubmissionRowID, v.IndikatorID, v.PeriodeID, nullString(&v.FaskesID),
			v.ScopeType, nullString(&v.ScopeID), v.Numerator, v.Denominator, v.Value,
			v.Status, nullString(&v.DataVersion), nowRFC3339()); err != nil {
			return fmt.Errorf("simpan nilai indikator: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Store) FindDuplicateRow(ctx context.Context, submissionID, labelID, indikatorID, scopeCode string) (bool, error) {
	var count int
	err := s.db.GetContext(ctx, &count, `SELECT COUNT(1) FROM submission_row
		WHERE submission_id = ? AND matched_label_id = ? AND matched_indikator_id = ?
		  AND COALESCE(scope_code,'') = COALESCE(?,'') AND classification_status = 'MATCHED'`,
		submissionID, labelID, indikatorID, scopeCode)
	return count > 0, err
}

func (s *Store) CountBlockingRows(ctx context.Context, submissionID string) (int, error) {
	var count int
	err := s.db.GetContext(ctx, &count, `SELECT COUNT(1) FROM submission_row
		WHERE submission_id = ? AND classification_status <> 'MATCHED'`, submissionID)
	return count, err
}
