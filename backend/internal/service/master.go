package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type MasterService struct {
	repo repository.MasterRepository
}

type Bootstrap struct {
	Users              []domain.User                 `json:"users"`
	Faskes             []domain.Faskes               `json:"faskes"`
	Wilayah            []domain.Wilayah              `json:"wilayah"`
	Bidang             []domain.Bidang               `json:"bidang"`
	Periods            []domain.Periode              `json:"periods"`
	DefaultPeriodID    string                        `json:"default_period_id"`
	Labels             []domain.Label                `json:"labels"`
	Indicators         []domain.Indikator            `json:"indicators"`
	Targets            []domain.Target               `json:"targets"`
	DashboardRows      []domain.DashboardRow         `json:"dashboard_rows"`
	Roles              []domain.Role                 `json:"roles"`
	SubmissionStatuses []domain.SubmissionStatus     `json:"submission_statuses"`
	Classifications    []domain.ClassificationStatus `json:"classifications"`
	ValueStatuses      []domain.ValueStatus          `json:"value_statuses"`
	TLStatuses         []domain.TLStatus             `json:"tl_statuses"`
}

func NewMasterService(repo repository.MasterRepository) *MasterService {
	return &MasterService{repo: repo}
}

func (s *MasterService) Bootstrap(ctx context.Context) (*Bootstrap, error) {
	users, err := s.repo.(repository.UserRepository).ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	faskes, err := s.repo.ListFaskes(ctx)
	if err != nil {
		return nil, err
	}
	wilayah, err := s.repo.ListWilayah(ctx)
	if err != nil {
		return nil, err
	}
	bidang, err := s.repo.ListBidang(ctx)
	if err != nil {
		return nil, err
	}
	periods, err := s.repo.ListPeriods(ctx)
	if err != nil {
		return nil, err
	}
	labels, err := s.repo.ListLabels(ctx)
	if err != nil {
		return nil, err
	}
	indicators, err := s.repo.ListIndicators(ctx)
	if err != nil {
		return nil, err
	}
	targets, err := s.repo.ListTargets(ctx, "")
	if err != nil {
		return nil, err
	}
	rows, err := s.repo.ListDashboardRows(ctx)
	if err != nil {
		return nil, err
	}
	defaultPeriod, err := s.repo.DefaultPeriodID(ctx)
	if err != nil {
		return nil, err
	}
	return &Bootstrap{
		Users: users, Faskes: faskes, Wilayah: wilayah, Bidang: bidang,
		Periods: periods, DefaultPeriodID: defaultPeriod, Labels: labels, Indicators: indicators, Targets: targets,
		DashboardRows: rows, Roles: domain.AllRoles,
		SubmissionStatuses: domain.AllSubmissionStatuses,
		Classifications:    domain.AllClassificationStatuses,
		ValueStatuses:      domain.AllValueStatuses,
		TLStatuses:         domain.AllTLStatuses,
	}, nil
}

func (s *MasterService) Labels(ctx context.Context) ([]domain.Label, error) {
	return s.repo.ListLabels(ctx)
}
func (s *MasterService) Indicators(ctx context.Context) ([]domain.Indikator, error) {
	return s.repo.ListIndicators(ctx)
}
func (s *MasterService) Targets(ctx context.Context, periodeID string) ([]domain.Target, error) {
	return s.repo.ListTargets(ctx, periodeID)
}
func (s *MasterService) Periods(ctx context.Context) ([]domain.Periode, error) {
	return s.repo.ListPeriods(ctx)
}

type CreateIndicatorInput struct {
	Kode             string   `json:"kode"`
	Nama             string   `json:"nama"`
	Satuan           string   `json:"satuan"`
	Definisi         string   `json:"definisi"`
	NumeratorLabel   string   `json:"numerator_label"`
	DenominatorLabel string   `json:"denominator_label"`
	SifatInput       string   `json:"sifat_input"`
	Frekuensi        string   `json:"frekuensi"`
	PJBidangID       string   `json:"pj_bidang_id"`
	PJTimker         string   `json:"pj_timker"`
	LabelIDs         []string `json:"label_ids"`
}

func (s *MasterService) CreateIndicator(ctx context.Context, in CreateIndicatorInput, actor *domain.User) (*domain.Indikator, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh menambah indikator")
	}
	if strings.TrimSpace(in.Kode) == "" || strings.TrimSpace(in.Nama) == "" || strings.TrimSpace(in.Satuan) == "" {
		return nil, domain.Validation("kode, nama, dan satuan wajib diisi", nil)
	}
	if in.SifatInput == "" {
		in.SifatInput = "puskesmas"
	}
	if in.Frekuensi == "" {
		in.Frekuensi = "bulanan"
	}
	ind := &domain.Indikator{
		ID:               uuid.NewString(),
		Kode:             strings.TrimSpace(in.Kode),
		Nama:             strings.TrimSpace(in.Nama),
		Satuan:           strings.TrimSpace(in.Satuan),
		Definisi:         strPtr(in.Definisi),
		NumeratorLabel:   strPtr(in.NumeratorLabel),
		DenominatorLabel: strPtr(in.DenominatorLabel),
		SifatInput:       in.SifatInput,
		Frekuensi:        in.Frekuensi,
		PJBidangID:       strPtr(in.PJBidangID),
		PJTimker:         strPtr(in.PJTimker),
		Aktif:            true,
	}
	if err := s.repo.CreateIndicator(ctx, ind, in.LabelIDs); err != nil {
		return nil, fmt.Errorf("simpan indikator: %w", err)
	}
	return ind, nil
}

type CreateLabelInput struct {
	Kode              string `json:"kode"`
	Nama              string `json:"nama"`
	Tipe              string `json:"tipe"`
	ParentID          string `json:"parent_id"`
	MasterReleaseID   string `json:"master_release_id"`
	AggregationMethod string `json:"aggregation_method"`
	AllowDirectInput  bool   `json:"allow_direct_input"`
	IsRollupOnly      bool   `json:"is_rollup_only"`
	PJBidangID        string `json:"pj_bidang_id"`
	IKUOwnerID        string `json:"iku_owner_id"`
}

func (s *MasterService) CreateLabel(ctx context.Context, in CreateLabelInput, actor *domain.User) (*domain.Label, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh menambah label")
	}
	if strings.TrimSpace(in.Kode) == "" || strings.TrimSpace(in.Nama) == "" {
		return nil, domain.Validation("kode dan nama label wajib diisi", nil)
	}
	if in.Tipe == "" {
		in.Tipe = "IKU"
	}
	if in.AggregationMethod == "" {
		in.AggregationMethod = string(domain.AggRatioOfSums)
	}
	if in.MasterReleaseID == "" {
		return nil, domain.Validation("master_release_id wajib diisi", nil)
	}
	level := 1
	var path *string
	if in.ParentID != "" {
		parent, err := s.repo.GetLabelByID(ctx, in.ParentID)
		if err != nil {
			return nil, err
		}
		level = parent.Level + 1
		p := in.Kode
		if parent.Path != nil {
			p = *parent.Path + "/" + in.Kode
		}
		path = &p
	} else {
		path = &in.Kode
	}
	label := &domain.Label{
		ID: uuid.NewString(), Kode: strings.TrimSpace(in.Kode), Nama: strings.TrimSpace(in.Nama),
		Tipe: in.Tipe, ParentID: strPtr(in.ParentID), MasterReleaseID: in.MasterReleaseID,
		AggregationMethod: domain.AggregationMethod(in.AggregationMethod),
		AllowDirectInput:  in.AllowDirectInput, IsRollupOnly: in.IsRollupOnly,
		PJBidangID: strPtr(in.PJBidangID), IKUOwnerID: strPtr(in.IKUOwnerID),
		Level: level, Path: path, Aktif: true,
	}
	if err := s.repo.CreateLabel(ctx, label); err != nil {
		return nil, fmt.Errorf("simpan label: %w", err)
	}
	return label, nil
}

type UpsertTargetInput struct {
	LabelID       string  `json:"label_id"`
	IndikatorID   string  `json:"indikator_id"`
	PeriodeID     string  `json:"periode_id"`
	ScopeType     string  `json:"scope_type"`
	ScopeID       string  `json:"scope_id"`
	TargetTahunan float64 `json:"target_tahunan"`
}

func (s *MasterService) UpsertTarget(ctx context.Context, in UpsertTargetInput, actor *domain.User) (*domain.Target, error) {
	if actor.Role != domain.RoleAdmin && actor.Role != domain.RoleKabid && actor.Role != domain.RoleSekretaris {
		return nil, domain.Forbidden("hanya admin/kabid/sekretaris yang boleh mengubah target")
	}
	if in.PeriodeID == "" || (in.LabelID == "" && in.IndikatorID == "") {
		return nil, domain.Validation("periode dan label/indikator target wajib diisi", nil)
	}
	if in.ScopeType == "" {
		in.ScopeType = "KABUPATEN"
	}
	if in.ScopeID == "" {
		in.ScopeID = "wil-kab-swb"
	}
	target := &domain.Target{
		ID: uuid.NewString(), LabelID: strPtr(in.LabelID), IndikatorID: strPtr(in.IndikatorID),
		PeriodeID: in.PeriodeID, ScopeType: in.ScopeType, ScopeID: strPtr(in.ScopeID),
		TargetTahunan: in.TargetTahunan, TargetBulanan: in.TargetTahunan / 12,
	}
	if err := s.repo.UpsertTarget(ctx, target); err != nil {
		return nil, err
	}
	return target, nil
}

func (s *MasterService) kabupatenScope(ctx context.Context) (string, error) {
	wilayah, err := s.repo.ListWilayah(ctx)
	if err != nil {
		return "", err
	}
	for _, w := range wilayah {
		if w.Level == "kabupaten" {
			return w.ID, nil
		}
	}
	return "", nil
}

func (s *MasterService) CapaianTahunan(ctx context.Context, tahun int) ([]domain.CapaianTahunan, error) {
	scopeID, err := s.kabupatenScope(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListCapaianTahunan(ctx, tahun, "KABUPATEN", scopeID)
}

type CapaianTahunanInput struct {
	LabelID     string   `json:"label_id"`
	IndikatorID string   `json:"indikator_id"`
	Tahun       int      `json:"tahun"`
	ScopeType   string   `json:"scope_type"`
	ScopeID     string   `json:"scope_id"`
	Nilai       *float64 `json:"nilai"`
	Catatan     string   `json:"catatan"`
}

func validateCapaianTarget(labelID, indikatorID string) error {
	labelID = strings.TrimSpace(labelID)
	indikatorID = strings.TrimSpace(indikatorID)
	if (labelID == "") == (indikatorID == "") {
		return domain.Validation("pilih salah satu: label/IKU atau indikator", nil)
	}
	return nil
}

func validateCapaianYear(tahun int) error {
	if tahun < 2000 || tahun > 2100 {
		return domain.Validation("tahun capaian harus antara 2000 dan 2100", nil)
	}
	return nil
}

func (s *MasterService) CreateCapaianTahunan(ctx context.Context, in CapaianTahunanInput, actor *domain.User) (*domain.CapaianTahunan, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh mengisi capaian tahunan")
	}
	if err := validateCapaianTarget(in.LabelID, in.IndikatorID); err != nil {
		return nil, err
	}
	if err := validateCapaianYear(in.Tahun); err != nil {
		return nil, err
	}
	if in.Nilai == nil || *in.Nilai < 0 {
		return nil, domain.Validation("nilai capaian wajib diisi dan tidak boleh negatif", nil)
	}
	scopeType := strings.TrimSpace(in.ScopeType)
	if scopeType == "" {
		scopeType = "KABUPATEN"
	}
	if scopeType != "KABUPATEN" {
		return nil, domain.Validation("scope capaian saat ini hanya KABUPATEN", nil)
	}
	scopeID := strings.TrimSpace(in.ScopeID)
	if scopeID == "" {
		var err error
		scopeID, err = s.kabupatenScope(ctx)
		if err != nil {
			return nil, err
		}
	}
	value := *in.Nilai
	c := &domain.CapaianTahunan{
		LabelID:     strPtr(strings.TrimSpace(in.LabelID)),
		IndikatorID: strPtr(strings.TrimSpace(in.IndikatorID)),
		Tahun:       in.Tahun,
		ScopeType:   scopeType,
		ScopeID:     strPtr(scopeID),
		Nilai:       value,
		Catatan:     strPtr(strings.TrimSpace(in.Catatan)),
		UpdatedBy:   &actor.ID,
		CreatedAt:   now(),
		UpdatedAt:   now(),
	}
	if err := s.repo.UpsertCapaianTahunan(ctx, c); err != nil {
		return nil, fmt.Errorf("simpan capaian tahunan: %w", err)
	}
	return s.repo.GetCapaianTahunan(ctx, c.ID)
}

func (s *MasterService) UpdateCapaianTahunan(ctx context.Context, id string, in CapaianTahunanInput, actor *domain.User) (*domain.CapaianTahunan, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh mengubah capaian tahunan")
	}
	c, err := s.repo.GetCapaianTahunan(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.LabelID) != "" || strings.TrimSpace(in.IndikatorID) != "" {
		if err := validateCapaianTarget(in.LabelID, in.IndikatorID); err != nil {
			return nil, err
		}
		c.LabelID = strPtr(strings.TrimSpace(in.LabelID))
		c.IndikatorID = strPtr(strings.TrimSpace(in.IndikatorID))
	}
	if in.Tahun != 0 {
		if err := validateCapaianYear(in.Tahun); err != nil {
			return nil, err
		}
		c.Tahun = in.Tahun
	}
	if in.Nilai != nil {
		if *in.Nilai < 0 {
			return nil, domain.Validation("nilai capaian tidak boleh negatif", nil)
		}
		c.Nilai = *in.Nilai
	}
	c.Catatan = strPtr(strings.TrimSpace(in.Catatan))
	if strings.TrimSpace(in.ScopeType) != "" && strings.TrimSpace(in.ScopeType) != "KABUPATEN" {
		return nil, domain.Validation("scope capaian saat ini hanya KABUPATEN", nil)
	}
	if strings.TrimSpace(in.ScopeID) != "" {
		c.ScopeID = strPtr(strings.TrimSpace(in.ScopeID))
	}
	c.UpdatedBy = &actor.ID
	c.UpdatedAt = now()
	if err := s.repo.UpsertCapaianTahunan(ctx, c); err != nil {
		return nil, fmt.Errorf("ubah capaian tahunan: %w", err)
	}
	return s.repo.GetCapaianTahunan(ctx, c.ID)
}

func (s *MasterService) DeleteCapaianTahunan(ctx context.Context, id string, actor *domain.User) error {
	if actor.Role != domain.RoleAdmin {
		return domain.Forbidden("hanya admin yang boleh menghapus capaian tahunan")
	}
	return s.repo.DeleteCapaianTahunan(ctx, id)
}
