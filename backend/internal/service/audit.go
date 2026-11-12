package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type AuditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) *AuditService { return &AuditService{repo: repo} }

func (s *AuditService) Write(ctx context.Context, actorID *string, action, entityType, entityID string, before, after any, reason, requestID string) error {
	return s.WriteWithImpersonator(ctx, actorID, nil, action, entityType, entityID, before, after, reason, requestID)
}

func (s *AuditService) WriteWithImpersonator(ctx context.Context, actorID, impersonatorID *string, action, entityType, entityID string, before, after any, reason, requestID string) error {
	if impersonatorID == nil {
		if u := ImpersonatorFromContext(ctx); u != nil {
			impersonatorID = &u.ID
		}
	}
	log := &domain.AuditLog{
		ID:             uuid.NewString(),
		ActorID:        actorID,
		ImpersonatorID: impersonatorID,
		Action:         action,
		EntityType:     entityType,
		EntityID:       &entityID,
		Reason:         strPtr(reason),
		RequestID:      strPtr(requestID),
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
	}
	if before != nil {
		if b, err := json.Marshal(before); err == nil {
			s := string(b)
			log.BeforeJSON = &s
		}
	}
	if after != nil {
		if b, err := json.Marshal(after); err == nil {
			s := string(b)
			log.AfterJSON = &s
		}
	}
	return s.repo.InsertAudit(ctx, log)
}

func (s *AuditService) List(ctx context.Context, f domain.AuditFilter) ([]domain.AuditLog, error) {
	if f.Limit <= 0 || f.Limit > 500 {
		f.Limit = 200
	}
	return s.repo.ListAudit(ctx, f)
}

func strPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
