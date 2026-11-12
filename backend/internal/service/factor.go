package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type FactorService struct {
	repo  repository.FactorRepository
	audit *AuditService
}

func NewFactorService(repo repository.FactorRepository, audit *AuditService) *FactorService {
	return &FactorService{repo: repo, audit: audit}
}

type FactorRow struct {
	IKUValue domain.IKUValue `json:"iku_value"`
	Faktor   *domain.Faktor  `json:"faktor,omitempty"`
}

func (s *FactorService) List(ctx context.Context, periodeID, scopeType, scopeID string) ([]FactorRow, error) {
	values, err := s.repo.ListIKUValuesForFactors(ctx, periodeID, scopeType, scopeID)
	if err != nil {
		return nil, err
	}
	out := make([]FactorRow, 0, len(values))
	for _, v := range values {
		if v.Status == domain.ValueNoData {
			continue
		}
		f, err := s.repo.GetFaktor(ctx, v.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, FactorRow{IKUValue: v, Faktor: f})
	}
	return out, nil
}

type UpsertFactorInput struct {
	IKUValueID string `json:"iku_value_id"`
	PeriodeID  string `json:"periode_id"`
	ScopeType  string `json:"scope_type"`
	ScopeID    string `json:"scope_id"`
	Pendukung  string `json:"pendukung"`
	Penghambat string `json:"penghambat"`
	RequestID  string `json:"-"`
}

func (s *FactorService) Upsert(ctx context.Context, actor *domain.User, in UpsertFactorInput) (*domain.Faktor, error) {
	if actor.Role != domain.RoleKabid && actor.Role != domain.RoleSekretaris && actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya Kabid/Sekretaris yang boleh mengisi faktor")
	}
	if strings.TrimSpace(in.IKUValueID) == "" {
		return nil, domain.Validation("iku_value_id wajib diisi", nil)
	}
	if in.PeriodeID == "" {
		return nil, domain.Validation("periode_id wajib diisi", nil)
	}
	if in.ScopeType == "" {
		in.ScopeType = "KABUPATEN"
	}
	if in.ScopeID == "" {
		in.ScopeID = "wil-kab-swb"
	}
	values, err := s.repo.ListIKUValuesForFactors(ctx, in.PeriodeID, in.ScopeType, in.ScopeID)
	if err != nil {
		return nil, err
	}
	found := false
	for _, v := range values {
		if v.ID == in.IKUValueID {
			found = true
			break
		}
	}
	if !found {
		return nil, domain.Validation("nilai IKU tidak ditemukan pada periode/scope tersebut", nil)
	}
	f := &domain.Faktor{
		IKUValueID: in.IKUValueID, PeriodeID: in.PeriodeID,
		Pendukung: strPtr(in.Pendukung), Penghambat: strPtr(in.Penghambat),
		CreatedBy: &actor.ID, UpdatedAt: now(),
	}
	if err := s.repo.UpsertFaktor(ctx, f); err != nil {
		return nil, err
	}
	_ = s.audit.Write(ctx, &actor.ID, "faktor.upsert", "iku_value", in.IKUValueID, nil, f, "faktor kualitatif", in.RequestID)
	return f, nil
}

type TindakLanjutInput struct {
	IKUValueID  string `json:"iku_value_id"`
	PeriodeID   string `json:"periode_id"`
	Masalah     string `json:"masalah"`
	Rekomendasi string `json:"rekomendasi"`
	PJUserID    string `json:"pj_user_id"`
	Status      string `json:"status"`
	RequestID   string `json:"-"`
}

func (s *FactorService) ListTL(ctx context.Context, periodeID string) ([]domain.TindakLanjut, error) {
	return s.repo.ListTindakLanjut(ctx, periodeID)
}

func (s *FactorService) CreateTL(ctx context.Context, actor *domain.User, in TindakLanjutInput) (*domain.TindakLanjut, error) {
	if actor.Role != domain.RolePimpinan && actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya pimpinan yang boleh menetapkan TL")
	}
	if strings.TrimSpace(in.Masalah) == "" || strings.TrimSpace(in.Rekomendasi) == "" || in.PeriodeID == "" {
		return nil, domain.Validation("masalah, rekomendasi, dan periode wajib diisi", nil)
	}
	status := domain.TLStatus(in.Status)
	if status == "" {
		status = domain.TLOpen
	}
	tl := &domain.TindakLanjut{
		ID: uuid.NewString(), IKUValueID: strPtr(in.IKUValueID), PeriodeID: in.PeriodeID,
		Masalah: in.Masalah, Rekomendasi: in.Rekomendasi, PJUserID: strPtr(in.PJUserID),
		Status: status, CreatedBy: &actor.ID, CreatedAt: now(), UpdatedAt: now(),
	}
	if err := s.repo.CreateTindakLanjut(ctx, tl); err != nil {
		return nil, err
	}
	_ = s.audit.Write(ctx, &actor.ID, "tindak_lanjut.created", "tindak_lanjut", tl.ID, nil, tl, "", in.RequestID)
	return tl, nil
}

func (s *FactorService) UpdateTL(ctx context.Context, actor *domain.User, id string, in TindakLanjutInput) (*domain.TindakLanjut, error) {
	if actor.Role != domain.RolePimpinan && actor.Role != domain.RoleKabid && actor.Role != domain.RoleSekretaris && actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("role tidak boleh mengubah TL")
	}
	tl, err := s.repo.GetTindakLanjut(ctx, id)
	if err != nil {
		return nil, err
	}
	before := *tl
	if in.Masalah != "" {
		tl.Masalah = in.Masalah
	}
	if in.Rekomendasi != "" {
		tl.Rekomendasi = in.Rekomendasi
	}
	if in.PJUserID != "" {
		tl.PJUserID = strPtr(in.PJUserID)
	}
	if in.Status != "" {
		status := domain.TLStatus(in.Status)
		valid := false
		for _, s := range domain.AllTLStatuses {
			if status == s {
				valid = true
			}
		}
		if !valid {
			return nil, domain.Validation("status TL tidak valid", nil)
		}
		tl.Status = status
	}
	tl.UpdatedAt = now()
	if err := s.repo.UpdateTindakLanjut(ctx, tl); err != nil {
		return nil, err
	}
	_ = s.audit.Write(ctx, &actor.ID, "tindak_lanjut.updated", "tindak_lanjut", id, before, tl, "", in.RequestID)
	return tl, nil
}
