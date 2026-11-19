package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

const (
	IKUStatusMenungguPuskesmas     = "menunggu_puskesmas"
	IKUStatusDikembalikanPuskesmas = "dikembalikan_puskesmas"
	IKUStatusMenungguAdmin         = "menunggu_admin"
	IKUStatusDikembalikanAdmin     = "dikembalikan_admin"
	IKUStatusDisetujui             = "disetujui"
	IKUStatusDitolak               = "ditolak"
	IKUStatusDibatalkan            = "dibatalkan"
)

type IKUChangeService struct {
	repo  repository.Repository
	audit *AuditService
}

func NewIKUChangeService(repo repository.Repository, audit *AuditService) *IKUChangeService {
	return &IKUChangeService{repo: repo, audit: audit}
}

type IKUChangeInput struct {
	Action            string `json:"action"`
	LabelID           string `json:"label_id"`
	Kode              string `json:"kode"`
	Nama              string `json:"nama"`
	CategoryID        string `json:"category_id"`
	ParentID          string `json:"parent_id"`
	MasterReleaseID   string `json:"master_release_id"`
	AggregationMethod string `json:"aggregation_method"`
	AllowDirectInput  bool   `json:"allow_direct_input"`
	IsRollupOnly      bool   `json:"is_rollup_only"`
	PJBidangID        string `json:"pj_bidang_id"`
	Notes             string `json:"notes"`
	RequestID         string `json:"-"`
}

type IKUChangeDecisionInput struct {
	Action    string `json:"action"`
	Notes     string `json:"notes"`
	RequestID string `json:"-"`
}

func (s *IKUChangeService) Create(ctx context.Context, actor *domain.User, in IKUChangeInput) (*domain.IKUChangeRequest, error) {
	if actor.Role != domain.RolePKMOperator {
		return nil, domain.Forbidden("hanya operator Puskesmas yang boleh mengajukan IKU")
	}
	if actor.FaskesID == nil || strings.TrimSpace(*actor.FaskesID) == "" {
		return nil, domain.Validation("operator harus terikat ke faskes", nil)
	}
	action := strings.ToUpper(strings.TrimSpace(in.Action))
	if action != "CREATE" && action != "UPDATE" && action != "DEACTIVATE" {
		return nil, domain.Validation("action harus CREATE, UPDATE, atau DEACTIVATE", nil)
	}
	releaseID := strings.TrimSpace(in.MasterReleaseID)
	var before *domain.Label
	if action == "UPDATE" || action == "DEACTIVATE" {
		label, err := s.repo.GetLabelByID(ctx, strings.TrimSpace(in.LabelID))
		if err != nil {
			return nil, err
		}
		before = label
		releaseID = label.MasterReleaseID
	}
	if action == "CREATE" || action == "UPDATE" {
		if strings.TrimSpace(in.Kode) == "" || strings.TrimSpace(in.Nama) == "" {
			return nil, domain.Validation("kode dan nama IKU wajib diisi", nil)
		}
		if strings.TrimSpace(in.CategoryID) == "" {
			return nil, domain.Validation("kategori IKU wajib diisi", nil)
		}
		category, err := s.repo.GetDashboardCategoryByID(ctx, strings.TrimSpace(in.CategoryID))
		if err != nil {
			return nil, err
		}
		if category.Jenis != "LABEL" || !category.Aktif {
			return nil, domain.Validation("kategori IKU harus aktif dan bertipe LABEL", nil)
		}
		if action == "CREATE" {
			if _, err := s.repo.GetLabelByKode(ctx, strings.TrimSpace(in.Kode), releaseID); err == nil {
				return nil, domain.Conflict("kode IKU sudah dipakai pada release ini")
			}
		}
		if parentID := strings.TrimSpace(in.ParentID); parentID != "" {
			if _, err := s.repo.GetLabelByID(ctx, parentID); err != nil {
				return nil, err
			}
		}
	}
	if releaseID == "" {
		release, err := s.repo.GetPublishedRelease(ctx, time.Now().Year())
		if err != nil {
			return nil, domain.Validation("master release untuk pengajuan IKU belum tersedia", nil)
		}
		releaseID = release.ID
	}
	payload := domain.IKUChangePayload{
		Kode: strings.TrimSpace(in.Kode), Nama: strings.TrimSpace(in.Nama),
		CategoryID: strings.TrimSpace(in.CategoryID), ParentID: strings.TrimSpace(in.ParentID),
		MasterReleaseID: releaseID, AggregationMethod: strings.TrimSpace(in.AggregationMethod),
		AllowDirectInput: in.AllowDirectInput, IsRollupOnly: in.IsRollupOnly,
		PJBidangID: strings.TrimSpace(in.PJBidangID),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request := &domain.IKUChangeRequest{
		ID: uuid.NewString(), Action: action, RequestedBy: actor.ID,
		FaskesID: *actor.FaskesID, MasterReleaseID: releaseID,
		PayloadJSON: string(payloadJSON), CurrentRevisionNo: 1,
		Status: IKUStatusMenungguPuskesmas, CreatedAt: now(), UpdatedAt: now(),
	}
	if before != nil {
		request.LabelID = &before.ID
	}
	if err := s.repo.CreateIKUChangeRequest(ctx, request); err != nil {
		return nil, fmt.Errorf("simpan pengajuan IKU: %w", err)
	}
	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(payload)
	revision := &domain.IKUChangeRevision{
		ID: uuid.NewString(), RequestID: request.ID, RevisionNo: 1,
		PayloadJSON: string(payloadJSON), BeforeJSON: optionalJSON(beforeJSON),
		AfterJSON: optionalJSON(afterJSON), CreatedBy: actor.ID, CreatedAt: now(),
	}
	if err := s.repo.AddIKUChangeRevision(ctx, revision); err != nil {
		return nil, fmt.Errorf("simpan revisi pengajuan IKU: %w", err)
	}
	_ = s.audit.Write(ctx, &actor.ID, "iku.change_requested", "iku_change_request", request.ID, nil, request, in.Notes, in.RequestID)
	return s.repo.GetIKUChangeRequest(ctx, request.ID)
}

func (s *IKUChangeService) List(ctx context.Context, actor *domain.User, f domain.IKUChangeFilter) ([]domain.IKUChangeRequest, error) {
	switch actor.Role {
	case domain.RoleAdmin:
		// semua
	case domain.RolePKMApprover:
		if actor.FaskesID == nil {
			return nil, domain.Forbidden("approver Puskesmas tidak terikat faskes")
		}
		f.FaskesID = *actor.FaskesID
	case domain.RolePKMOperator:
		f.RequestedBy = actor.ID
	default:
		return nil, domain.Forbidden("role tidak boleh melihat pengajuan IKU")
	}
	return s.repo.ListIKUChangeRequests(ctx, f)
}

func (s *IKUChangeService) Detail(ctx context.Context, actor *domain.User, id string) (*domain.IKUChangeRequest, []domain.IKUChangeRevision, []domain.IKUChangeReview, error) {
	request, err := s.repo.GetIKUChangeRequest(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := s.canView(actor, request); err != nil {
		return nil, nil, nil, err
	}
	revisions, err := s.repo.ListIKUChangeRevisions(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	reviews, err := s.repo.ListIKUChangeReviews(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}
	return request, revisions, reviews, nil
}

func (s *IKUChangeService) Decide(ctx context.Context, actor *domain.User, id string, in IKUChangeDecisionInput) (*domain.IKUChangeRequest, error) {
	request, err := s.repo.GetIKUChangeRequest(ctx, id)
	if err != nil {
		return nil, err
	}
	action := strings.ToLower(strings.TrimSpace(in.Action))
	if action != "approve" && action != "reject" && action != "return" {
		return nil, domain.Validation("action harus approve, reject, atau return", nil)
	}
	notes := strings.TrimSpace(in.Notes)
	if (action == "reject" || action == "return") && notes == "" {
		return nil, domain.Validation("catatan wajib diisi saat reject/return", nil)
	}
	var nextStatus string
	switch request.Status {
	case IKUStatusMenungguPuskesmas:
		if actor.Role != domain.RolePKMApprover {
			return nil, domain.Forbidden("hanya approver Puskesmas yang boleh memutus tahap ini")
		}
		if actor.FaskesID == nil || *actor.FaskesID != request.FaskesID {
			return nil, domain.Forbidden("pengajuan di luar faskes approver")
		}
		if actor.ID == request.RequestedBy {
			return nil, domain.Forbidden("tidak boleh menyetujui pengajuan sendiri")
		}
		if action == "approve" {
			nextStatus = IKUStatusMenungguAdmin
		} else {
			nextStatus = IKUStatusDikembalikanPuskesmas
		}
		_ = s.repo.AddIKUChangeReview(ctx, &domain.IKUChangeReview{
			RequestID: id, RevisionNo: request.CurrentRevisionNo, Step: "PKM",
			ReviewerID: actor.ID, Action: action, Notes: strPtr(notes), CreatedAt: now(),
		})
	case IKUStatusMenungguAdmin:
		if actor.Role != domain.RoleAdmin {
			return nil, domain.Forbidden("hanya Admin yang boleh memutus tahap final")
		}
		if action == "approve" {
			label, err := s.apply(ctx, request)
			if err != nil {
				return nil, err
			}
			nextStatus = IKUStatusDisetujui
			request.AppliedLabelID = &label.ID
		} else {
			nextStatus = IKUStatusDikembalikanAdmin
		}
		_ = s.repo.AddIKUChangeReview(ctx, &domain.IKUChangeReview{
			RequestID: id, RevisionNo: request.CurrentRevisionNo, Step: "ADMIN",
			ReviewerID: actor.ID, Action: action, Notes: strPtr(notes), CreatedAt: now(),
		})
	default:
		return nil, domain.Conflict("pengajuan IKU tidak berada pada tahap approval")
	}
	if err := s.repo.UpdateIKUChangeRequestStatus(ctx, id, nextStatus, request.AppliedLabelID, now()); err != nil {
		return nil, err
	}
	_ = s.audit.Write(ctx, &actor.ID, "iku.change_"+action, "iku_change_request", id, request, nil, notes, in.RequestID)
	return s.repo.GetIKUChangeRequest(ctx, id)
}

func (s *IKUChangeService) apply(ctx context.Context, request *domain.IKUChangeRequest) (*domain.Label, error) {
	var payload domain.IKUChangePayload
	if err := json.Unmarshal([]byte(request.PayloadJSON), &payload); err != nil {
		return nil, domain.Validation("payload pengajuan IKU tidak valid", nil)
	}
	var category *domain.DashboardCategory
	if request.Action != "DEACTIVATE" {
		var err error
		category, err = s.repo.GetDashboardCategoryByID(ctx, payload.CategoryID)
		if err != nil {
			return nil, err
		}
		if category.Jenis != "LABEL" || !category.Aktif {
			return nil, domain.Validation("kategori IKU tidak aktif", nil)
		}
	}
	switch request.Action {
	case "CREATE":
		level := 1
		path := payload.Kode
		if parentID := payload.ParentID; parentID != "" {
			parent, err := s.repo.GetLabelByID(ctx, parentID)
			if err != nil {
				return nil, err
			}
			level = parent.Level + 1
			if parent.Path != nil {
				path = *parent.Path + "/" + payload.Kode
			}
		}
		label := &domain.Label{
			ID: uuid.NewString(), Kode: payload.Kode, Nama: payload.Nama,
			Tipe: category.Kode, CategoryID: &category.ID,
			ParentID: strPtr(payload.ParentID), MasterReleaseID: request.MasterReleaseID,
			AggregationMethod: domain.AggregationMethod(payload.AggregationMethod),
			AllowDirectInput:  payload.AllowDirectInput, IsRollupOnly: payload.IsRollupOnly,
			PJBidangID: strPtr(payload.PJBidangID), Level: level, Path: strPtr(path), Aktif: true,
		}
		if label.AggregationMethod == "" {
			label.AggregationMethod = domain.AggRatioOfSums
		}
		if err := s.repo.CreateLabel(ctx, label); err != nil {
			return nil, err
		}
		return s.repo.GetLabelByID(ctx, label.ID)
	case "UPDATE", "DEACTIVATE":
		label, err := s.repo.GetLabelByID(ctx, deref(request.LabelID))
		if err != nil {
			return nil, err
		}
		if request.Action == "DEACTIVATE" {
			label.Aktif = false
		} else {
			label.Nama = payload.Nama
			label.Kode = payload.Kode
			label.Tipe = category.Kode
			label.CategoryID = &category.ID
			label.ParentID = strPtr(payload.ParentID)
			label.AggregationMethod = domain.AggregationMethod(payload.AggregationMethod)
			label.AllowDirectInput = payload.AllowDirectInput
			label.IsRollupOnly = payload.IsRollupOnly
			label.PJBidangID = strPtr(payload.PJBidangID)
		}
		if err := s.repo.UpdateLabel(ctx, label); err != nil {
			return nil, err
		}
		return s.repo.GetLabelByID(ctx, label.ID)
	default:
		return nil, domain.Validation("action IKU tidak dikenal", nil)
	}
}

func (s *IKUChangeService) canView(actor *domain.User, request *domain.IKUChangeRequest) error {
	switch actor.Role {
	case domain.RoleAdmin:
		return nil
	case domain.RolePKMApprover:
		if actor.FaskesID != nil && *actor.FaskesID == request.FaskesID {
			return nil
		}
	case domain.RolePKMOperator:
		if actor.ID == request.RequestedBy {
			return nil
		}
	}
	return domain.Forbidden("tidak boleh melihat pengajuan IKU ini")
}

func optionalJSON(value []byte) *string {
	if len(value) == 0 || string(value) == "null" {
		return nil
	}
	text := string(value)
	return &text
}
