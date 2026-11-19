package sqlite

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
)

const ikuChangeSelect = `SELECT r.id, r.action, r.label_id, r.requested_by, r.faskes_id,
	r.master_release_id, r.payload_json, r.current_revision_no, r.status,
	r.applied_label_id, r.created_at, r.updated_at,
	u.nama AS requested_by_nama, f.nama AS faskes_nama,
	l.kode AS label_kode, l.nama AS label_nama, a.kode AS applied_label_kode
	FROM iku_change_request r
	JOIN users u ON u.id = r.requested_by
	JOIN faskes f ON f.id = r.faskes_id
	LEFT JOIN label l ON l.id = r.label_id
	LEFT JOIN label a ON a.id = r.applied_label_id`

func (s *Store) CreateIKUChangeRequest(ctx context.Context, request *domain.IKUChangeRequest) error {
	if request.ID == "" {
		request.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO iku_change_request
		(id, action, label_id, requested_by, faskes_id, master_release_id, payload_json,
		 current_revision_no, status, applied_label_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		request.ID, request.Action, nullString(request.LabelID), request.RequestedBy,
		request.FaskesID, request.MasterReleaseID, request.PayloadJSON,
		request.CurrentRevisionNo, request.Status, nullString(request.AppliedLabelID),
		request.CreatedAt, request.UpdatedAt)
	return err
}

func (s *Store) GetIKUChangeRequest(ctx context.Context, id string) (*domain.IKUChangeRequest, error) {
	var out domain.IKUChangeRequest
	err := s.db.GetContext(ctx, &out, ikuChangeSelect+` WHERE r.id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("pengajuan IKU")
	}
	return &out, err
}

func (s *Store) ListIKUChangeRequests(ctx context.Context, f domain.IKUChangeFilter) ([]domain.IKUChangeRequest, error) {
	q := ikuChangeSelect + ` WHERE 1 = 1`
	args := []any{}
	if f.Status != "" {
		q += ` AND r.status = ?`
		args = append(args, f.Status)
	}
	if f.FaskesID != "" {
		q += ` AND r.faskes_id = ?`
		args = append(args, f.FaskesID)
	}
	if f.RequestedBy != "" {
		q += ` AND r.requested_by = ?`
		args = append(args, f.RequestedBy)
	}
	q += ` ORDER BY r.updated_at DESC`
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 100
	}
	q += ` LIMIT ?`
	args = append(args, f.Limit)
	var out []domain.IKUChangeRequest
	err := s.db.SelectContext(ctx, &out, q, args...)
	return out, err
}

func (s *Store) UpdateIKUChangeRequestPayload(ctx context.Context, id, payloadJSON string, revisionNo int, updatedAt string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE iku_change_request
		SET payload_json = ?, current_revision_no = ?, updated_at = ? WHERE id = ?`,
		payloadJSON, revisionNo, updatedAt, id)
	return err
}

func (s *Store) UpdateIKUChangeRequestStatus(ctx context.Context, id, status string, appliedLabelID *string, updatedAt string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE iku_change_request
		SET status = ?, applied_label_id = ?, updated_at = ? WHERE id = ?`,
		status, nullString(appliedLabelID), updatedAt, id)
	return err
}

func (s *Store) AddIKUChangeRevision(ctx context.Context, revision *domain.IKUChangeRevision) error {
	if revision.ID == "" {
		revision.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO iku_change_revision
		(id, request_id, revision_no, payload_json, before_json, after_json, created_by, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		revision.ID, revision.RequestID, revision.RevisionNo, revision.PayloadJSON,
		nullString(revision.BeforeJSON), nullString(revision.AfterJSON),
		revision.CreatedBy, revision.CreatedAt)
	return err
}

func (s *Store) ListIKUChangeRevisions(ctx context.Context, requestID string) ([]domain.IKUChangeRevision, error) {
	var out []domain.IKUChangeRevision
	err := s.db.SelectContext(ctx, &out, `SELECT id, request_id, revision_no, payload_json,
		before_json, after_json, created_by, created_at
		FROM iku_change_revision WHERE request_id = ? ORDER BY revision_no`, requestID)
	return out, err
}

func (s *Store) AddIKUChangeReview(ctx context.Context, review *domain.IKUChangeReview) error {
	if review.ID == "" {
		review.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO iku_change_review
		(id, request_id, revision_no, step, reviewer_id, action, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		review.ID, review.RequestID, review.RevisionNo, review.Step,
		review.ReviewerID, review.Action, nullString(review.Notes), review.CreatedAt)
	return err
}

func (s *Store) ListIKUChangeReviews(ctx context.Context, requestID string) ([]domain.IKUChangeReview, error) {
	var out []domain.IKUChangeReview
	err := s.db.SelectContext(ctx, &out, `SELECT r.id, r.request_id, r.revision_no, r.step,
		r.reviewer_id, r.action, r.notes, r.created_at, u.nama AS reviewer_nama
		FROM iku_change_review r JOIN users u ON u.id = r.reviewer_id
		WHERE r.request_id = ? ORDER BY r.created_at`, requestID)
	return out, err
}
