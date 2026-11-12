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
	Categories         []domain.DashboardCategory    `json:"categories"`
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
	categories, err := s.repo.ListDashboardCategories(ctx)
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
	if err := s.attachIndicatorLabels(ctx, indicators); err != nil {
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
		Users: users, Faskes: faskes, Wilayah: wilayah, Bidang: bidang, Categories: categories,
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
	indicators, err := s.repo.ListIndicators(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.attachIndicatorLabels(ctx, indicators); err != nil {
		return nil, err
	}
	return indicators, nil
}

func (s *MasterService) attachIndicatorLabels(ctx context.Context, indicators []domain.Indikator) error {
	links, err := s.repo.ListAllIndicatorLabels(ctx)
	if err != nil {
		return err
	}
	byIndicator := map[string][]domain.IndikatorLabel{}
	for _, link := range links {
		byIndicator[link.IndikatorID] = append(byIndicator[link.IndikatorID], link)
	}
	for i := range indicators {
		for _, link := range byIndicator[indicators[i].ID] {
			indicators[i].LabelIDs = append(indicators[i].LabelIDs, link.LabelID)
			if link.IsPrimary {
				labelID := link.LabelID
				indicators[i].PrimaryLabelID = &labelID
			}
		}
	}
	return nil
}
func (s *MasterService) Targets(ctx context.Context, periodeID string) ([]domain.Target, error) {
	return s.repo.ListTargets(ctx, periodeID)
}
func (s *MasterService) Periods(ctx context.Context) ([]domain.Periode, error) {
	return s.repo.ListPeriods(ctx)
}

func (s *MasterService) Categories(ctx context.Context) ([]domain.DashboardCategory, error) {
	return s.repo.ListDashboardCategories(ctx)
}

type CategoryInput struct {
	Kode   string `json:"kode"`
	Nama   string `json:"nama"`
	Jenis  string `json:"jenis"`
	Urutan int    `json:"urutan"`
	Aktif  *bool  `json:"aktif"`
}

func (s *MasterService) CreateCategory(ctx context.Context, in CategoryInput, actor *domain.User) (*domain.DashboardCategory, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh menambah kategori")
	}
	kode := strings.ToUpper(strings.TrimSpace(in.Kode))
	nama := strings.TrimSpace(in.Nama)
	if kode == "" || nama == "" {
		return nil, domain.Validation("kode dan nama kategori wajib diisi", nil)
	}
	jenis := strings.ToUpper(strings.TrimSpace(in.Jenis))
	if jenis == "" {
		jenis = "LABEL"
	}
	if jenis != "LABEL" && jenis != "OWNER" {
		return nil, domain.Validation("jenis kategori harus LABEL atau OWNER", nil)
	}
	aktif := true
	if in.Aktif != nil {
		aktif = *in.Aktif
	}
	now := now()
	c := &domain.DashboardCategory{
		ID: uuid.NewString(), Kode: kode, Nama: nama, Jenis: jenis, Urutan: in.Urutan,
		Aktif: aktif, IsSystem: false, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.CreateDashboardCategory(ctx, c); err != nil {
		return nil, fmt.Errorf("simpan kategori: %w", err)
	}
	return c, nil
}

func (s *MasterService) UpdateCategory(ctx context.Context, id string, in CategoryInput, actor *domain.User) (*domain.DashboardCategory, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh mengubah kategori")
	}
	c, err := s.repo.GetDashboardCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if nama := strings.TrimSpace(in.Nama); nama != "" {
		c.Nama = nama
	}
	if c.IsSystem {
		if jenis := strings.ToUpper(strings.TrimSpace(in.Jenis)); jenis != "" && jenis != c.Jenis {
			return nil, domain.Validation("jenis kategori sistem tidak boleh diubah", nil)
		}
	} else if jenis := strings.ToUpper(strings.TrimSpace(in.Jenis)); jenis != "" {
		if jenis != "LABEL" && jenis != "OWNER" {
			return nil, domain.Validation("jenis kategori harus LABEL atau OWNER", nil)
		}
		c.Jenis = jenis
	}
	if in.Urutan != 0 {
		c.Urutan = in.Urutan
	}
	if in.Aktif != nil {
		c.Aktif = *in.Aktif
	}
	c.UpdatedAt = now()
	if err := s.repo.UpdateDashboardCategory(ctx, c); err != nil {
		return nil, fmt.Errorf("ubah kategori: %w", err)
	}
	return c, nil
}

func (s *MasterService) DeactivateCategory(ctx context.Context, id string, actor *domain.User) (*domain.DashboardCategory, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh menonaktifkan kategori")
	}
	c, err := s.repo.GetDashboardCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Aktif = false
	c.UpdatedAt = now()
	if err := s.repo.UpdateDashboardCategory(ctx, c); err != nil {
		return nil, fmt.Errorf("nonaktifkan kategori: %w", err)
	}
	return c, nil
}

type BidangInput struct {
	Kode     string  `json:"kode"`
	Nama     string  `json:"nama"`
	Jenis    string  `json:"jenis"`
	ParentID *string `json:"parent_id"`
	Aktif    *bool   `json:"aktif"`
}

func (s *MasterService) CreateBidang(ctx context.Context, in BidangInput, actor *domain.User) (*domain.Bidang, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh menambah bidang/timker")
	}
	b, err := s.bidangInput(ctx, in, nil)
	if err != nil {
		return nil, err
	}
	b.ID = uuid.NewString()
	b.Aktif = true
	if in.Aktif != nil {
		b.Aktif = *in.Aktif
	}
	if err := s.repo.CreateBidang(ctx, b); err != nil {
		return nil, fmt.Errorf("simpan bidang: %w", err)
	}
	return b, nil
}

func (s *MasterService) UpdateBidang(ctx context.Context, id string, in BidangInput, actor *domain.User) (*domain.Bidang, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh mengubah bidang/timker")
	}
	existing, err := s.repo.GetBidang(ctx, id)
	if err != nil {
		return nil, err
	}
	b, err := s.bidangInput(ctx, in, existing)
	if err != nil {
		return nil, err
	}
	b.ID = id
	if in.Aktif != nil {
		b.Aktif = *in.Aktif
	}
	if err := s.repo.UpdateBidang(ctx, b); err != nil {
		return nil, fmt.Errorf("ubah bidang: %w", err)
	}
	return b, nil
}

func (s *MasterService) bidangInput(ctx context.Context, in BidangInput, existing *domain.Bidang) (*domain.Bidang, error) {
	b := &domain.Bidang{}
	if existing != nil {
		*b = *existing
	}
	if kode := strings.ToUpper(strings.TrimSpace(in.Kode)); kode != "" {
		b.Kode = kode
	}
	if nama := strings.TrimSpace(in.Nama); nama != "" {
		b.Nama = nama
	}
	if b.Kode == "" || b.Nama == "" {
		return nil, domain.Validation("kode dan nama bidang wajib diisi", nil)
	}
	jenis := strings.ToLower(strings.TrimSpace(in.Jenis))
	if jenis != "" {
		if jenis != "bidang" && jenis != "sekretariat" && jenis != "timker" {
			return nil, domain.Validation("jenis bidang harus bidang, sekretariat, atau timker", nil)
		}
		b.Jenis = jenis
	}
	if b.Jenis == "" {
		b.Jenis = "bidang"
	}
	if in.ParentID != nil {
		if *in.ParentID == "" {
			b.ParentID = nil
		} else {
			parent, err := s.repo.GetBidang(ctx, *in.ParentID)
			if err != nil {
				return nil, err
			}
			if parent.Jenis == "timker" {
				return nil, domain.Validation("parent bidang tidak boleh timker", nil)
			}
			b.ParentID = in.ParentID
		}
	}
	return b, nil
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
	PJTimkerID       string   `json:"pj_timker_id"`
	PJTimker         string   `json:"pj_timker"`
	LabelIDs         []string `json:"label_ids"`
	PrimaryLabelID   string   `json:"primary_label_id"`
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
	labelIDs, primaryLabelID, err := s.validateIndicatorLabels(ctx, in.LabelIDs, in.PrimaryLabelID)
	if err != nil {
		return nil, err
	}
	pjTimker := strings.TrimSpace(in.PJTimker)
	if id := strings.TrimSpace(in.PJTimkerID); id != "" {
		b, err := s.repo.GetBidang(ctx, id)
		if err != nil {
			return nil, err
		}
		if b.Jenis != "timker" {
			return nil, domain.Validation("pj_timker_id harus menunjuk timker", nil)
		}
		pjTimker = b.Nama
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
		PJTimkerID:       strPtr(in.PJTimkerID),
		PJTimker:         strPtr(pjTimker),
		Aktif:            true,
	}
	if err := s.repo.CreateIndicator(ctx, ind, labelIDs, primaryLabelID); err != nil {
		return nil, fmt.Errorf("simpan indikator: %w", err)
	}
	ind.LabelIDs = labelIDs
	ind.PrimaryLabelID = &primaryLabelID
	return ind, nil
}

func (s *MasterService) validateIndicatorLabels(ctx context.Context, labelIDs []string, primaryLabelID string) ([]string, string, error) {
	seen := map[string]bool{}
	clean := make([]string, 0, len(labelIDs))
	for _, id := range labelIDs {
		id = strings.TrimSpace(id)
		if id == "" || seen[id] {
			continue
		}
		label, err := s.repo.GetLabelByID(ctx, id)
		if err != nil {
			return nil, "", domain.Validation("label indikator tidak ditemukan", map[string]any{"label_id": id})
		}
		if !label.Aktif {
			return nil, "", domain.Validation("label indikator tidak aktif", map[string]any{"label_id": id})
		}
		if label.CategoryID == nil {
			return nil, "", domain.Validation("label indikator belum punya kategori", map[string]any{"label_id": id})
		}
		category, err := s.repo.GetDashboardCategoryByID(ctx, *label.CategoryID)
		if err != nil {
			return nil, "", domain.Validation("kategori label indikator tidak ditemukan", map[string]any{"label_id": id})
		}
		if !category.Aktif {
			return nil, "", domain.Validation("kategori label indikator tidak aktif", map[string]any{"label_id": id})
		}
		seen[id] = true
		clean = append(clean, id)
	}
	if len(clean) == 0 {
		return nil, "", domain.Validation("minimal satu label indikator wajib dipilih", nil)
	}
	primary := strings.TrimSpace(primaryLabelID)
	if primary == "" {
		primary = clean[0]
	}
	if !seen[primary] {
		return nil, "", domain.Validation("label utama harus termasuk label yang dipilih", nil)
	}
	return clean, primary, nil
}

type CreateLabelInput struct {
	Kode              string `json:"kode"`
	Nama              string `json:"nama"`
	Tipe              string `json:"tipe"`
	CategoryID        string `json:"category_id"`
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
	category, err := s.resolveCategory(ctx, in.CategoryID, in.Tipe)
	if err != nil {
		return nil, err
	}
	if category.Jenis != "LABEL" {
		return nil, domain.Validation("kategori owner tidak bisa dipakai sebagai label taksonomi", nil)
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
		Tipe: category.Kode, CategoryID: &category.ID, ParentID: strPtr(in.ParentID),
		MasterReleaseID:   in.MasterReleaseID,
		AggregationMethod: domain.AggregationMethod(in.AggregationMethod),
		AllowDirectInput:  in.AllowDirectInput, IsRollupOnly: in.IsRollupOnly,
		PJBidangID: strPtr(in.PJBidangID), IKUOwnerID: strPtr(in.IKUOwnerID),
		Level: level, Path: path, Aktif: true,
	}
	if err := s.repo.CreateLabel(ctx, label); err != nil {
		return nil, fmt.Errorf("simpan label: %w", err)
	}
	return s.repo.GetLabelByID(ctx, label.ID)
}

func cleanOptionalID(id *string) *string {
	if id == nil || strings.TrimSpace(*id) == "" {
		return nil
	}
	value := strings.TrimSpace(*id)
	return &value
}

func (s *MasterService) resolveCategory(ctx context.Context, categoryID, tipe string) (*domain.DashboardCategory, error) {
	id := strings.TrimSpace(categoryID)
	var (
		category *domain.DashboardCategory
		err      error
	)
	if id != "" {
		category, err = s.repo.GetDashboardCategoryByID(ctx, id)
	} else {
		kode := strings.ToUpper(strings.TrimSpace(tipe))
		if kode == "" {
			return nil, domain.Validation("category_id wajib diisi", nil)
		}
		category, err = s.repo.GetDashboardCategoryByKode(ctx, kode)
	}
	if err != nil {
		return nil, err
	}
	if !category.Aktif {
		return nil, domain.Validation("kategori tidak aktif", nil)
	}
	return category, nil
}

type UpdateLabelInput struct {
	Nama             string  `json:"nama"`
	CategoryID       string  `json:"category_id"`
	ParentID         *string `json:"parent_id"`
	PJBidangID       *string `json:"pj_bidang_id"`
	IKUOwnerID       *string `json:"iku_owner_id"`
	AllowDirectInput *bool   `json:"allow_direct_input"`
	IsRollupOnly     *bool   `json:"is_rollup_only"`
	Aktif            *bool   `json:"aktif"`
}

func (s *MasterService) UpdateLabel(ctx context.Context, id string, in UpdateLabelInput, actor *domain.User) (*domain.Label, error) {
	if actor.Role != domain.RoleAdmin {
		return nil, domain.Forbidden("hanya admin yang boleh mengubah label")
	}
	label, err := s.repo.GetLabelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if nama := strings.TrimSpace(in.Nama); nama != "" {
		label.Nama = nama
	}
	if strings.TrimSpace(in.CategoryID) != "" {
		category, err := s.repo.GetDashboardCategoryByID(ctx, strings.TrimSpace(in.CategoryID))
		if err != nil {
			return nil, err
		}
		if category.Jenis != "LABEL" {
			return nil, domain.Validation("kategori owner tidak bisa dipakai sebagai label taksonomi", nil)
		}
		if !category.Aktif {
			return nil, domain.Validation("kategori tidak aktif", nil)
		}
		label.CategoryID = &category.ID
		label.Tipe = category.Kode
	}
	if in.ParentID != nil {
		label.ParentID = cleanOptionalID(in.ParentID)
	}
	if in.PJBidangID != nil {
		label.PJBidangID = cleanOptionalID(in.PJBidangID)
	}
	if in.IKUOwnerID != nil {
		label.IKUOwnerID = cleanOptionalID(in.IKUOwnerID)
	}
	if in.AllowDirectInput != nil {
		label.AllowDirectInput = *in.AllowDirectInput
	}
	if in.IsRollupOnly != nil {
		label.IsRollupOnly = *in.IsRollupOnly
	}
	if in.Aktif != nil {
		label.Aktif = *in.Aktif
	}
	if err := s.repo.UpdateLabel(ctx, label); err != nil {
		return nil, fmt.Errorf("ubah label: %w", err)
	}
	return s.repo.GetLabelByID(ctx, id)
}

type UpsertTargetInput struct {
	LabelID              string   `json:"label_id"`
	IndikatorID          string   `json:"indikator_id"`
	PeriodeID            string   `json:"periode_id"`
	ScopeType            string   `json:"scope_type"`
	ScopeID              string   `json:"scope_id"`
	TargetTahunan        float64  `json:"target_tahunan"`
	TargetLayananTahunan *float64 `json:"target_layanan_tahunan"`
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
	layananTahunan := in.TargetLayananTahunan
	if layananTahunan == nil {
		if existing, err := s.repo.ListTargets(ctx, in.PeriodeID); err == nil {
			for _, item := range existing {
				if deref(item.LabelID) != in.LabelID || deref(item.IndikatorID) != in.IndikatorID {
					continue
				}
				if item.ScopeType != in.ScopeType || deref(item.ScopeID) != in.ScopeID {
					continue
				}
				layananTahunan = item.TargetLayananTahunan
				break
			}
		}
	}
	target := &domain.Target{
		ID: uuid.NewString(), LabelID: strPtr(in.LabelID), IndikatorID: strPtr(in.IndikatorID),
		PeriodeID: in.PeriodeID, ScopeType: in.ScopeType, ScopeID: strPtr(in.ScopeID),
		TargetTahunan: in.TargetTahunan, TargetBulanan: in.TargetTahunan / 12,
		TargetLayananTahunan: layananTahunan,
	}
	if layananTahunan != nil {
		monthly := *layananTahunan / 12
		target.TargetLayananBulanan = &monthly
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
