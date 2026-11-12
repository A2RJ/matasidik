package sqlite

import (
	"context"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
)

func (s *Store) InsertAudit(ctx context.Context, log *domain.AuditLog) error {
	if log.ID == "" {
		log.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_log
		(id, actor_id, impersonator_id, action, entity_type, entity_id, before_json, after_json, reason, request_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		log.ID, log.ActorID, log.ImpersonatorID, log.Action, log.EntityType, log.EntityID, log.BeforeJSON,
		log.AfterJSON, log.Reason, log.RequestID, log.CreatedAt)
	return err
}

func (s *Store) ListAudit(ctx context.Context, f domain.AuditFilter) ([]domain.AuditLog, error) {
	q := `SELECT a.id, a.actor_id, a.impersonator_id, a.action, a.entity_type, a.entity_id, a.before_json,
		a.after_json, a.reason, a.request_id, a.created_at, u.nama AS actor_nama,
		im.nama AS impersonator_nama
		FROM audit_log a
		LEFT JOIN users u ON u.id = a.actor_id
		LEFT JOIN users im ON im.id = a.impersonator_id WHERE 1=1`
	args := []any{}
	if f.EntityType != "" {
		q += ` AND a.entity_type = ?`
		args = append(args, f.EntityType)
	}
	if f.EntityID != "" {
		q += ` AND a.entity_id = ?`
		args = append(args, f.EntityID)
	}
	q += ` ORDER BY a.created_at DESC`
	if f.Limit > 0 {
		q += ` LIMIT ?`
		args = append(args, f.Limit)
	}
	var out []domain.AuditLog
	err := s.db.SelectContext(ctx, &out, q, args...)
	return out, err
}
