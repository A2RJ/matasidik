package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type SubmissionService struct {
	subs    repository.SubmissionRepository
	master  repository.MasterRepository
	audit   *AuditService
	calc    *CalculationService
	dataDir string
	maxSize int64
}

func NewSubmissionService(subs repository.SubmissionRepository, master repository.MasterRepository, audit *AuditService, calc *CalculationService, dataDir string, maxSize int64) *SubmissionService {
	return &SubmissionService{subs: subs, master: master, audit: audit, calc: calc, dataDir: dataDir, maxSize: maxSize}
}

type UploadInput struct {
	File           multipart.File
	Filename       string
	Size           int64
	PeriodeKode    string
	FaskesID       string
	IdempotencyKey string
	RequestID      string
}

func (s *SubmissionService) Upload(ctx context.Context, actor *domain.User, in UploadInput) (*domain.Submission, *domain.Job, error) {
	if actor.Role != domain.RolePKMOperator && actor.Role != domain.RolePKMPJ && actor.Role != domain.RoleAdmin {
		return nil, nil, domain.Forbidden("hanya Puskesmas/PJ atau admin yang boleh upload")
	}
	if in.Size > s.maxSize {
		return nil, nil, domain.NewError(http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "ukuran file melebihi batas")
	}
	faskesID := strings.TrimSpace(in.FaskesID)
	if faskesID == "" && actor.FaskesID != nil {
		faskesID = *actor.FaskesID
	}
	if faskesID == "" {
		return nil, nil, domain.Validation("faskes wajib dipilih", nil)
	}
	if actor.Role != domain.RoleAdmin && (actor.FaskesID == nil || *actor.FaskesID != faskesID) {
		return nil, nil, domain.Forbidden("upload hanya untuk faskes sendiri")
	}
	if _, err := s.master.GetFaskes(ctx, faskesID); err != nil {
		return nil, nil, err
	}
	period, err := s.master.GetPeriodByKode(ctx, strings.TrimSpace(in.PeriodeKode))
	if err != nil {
		return nil, nil, err
	}
	if period.Status != "open" {
		return nil, nil, domain.Conflict("periode sudah terkunci")
	}
	if key := strings.TrimSpace(in.IdempotencyKey); key != "" {
		if existing, err := s.subs.FindIdempotentSubmission(ctx, key); err == nil && existing != nil {
			return existing, nil, nil
		}
	}
	release, err := s.master.GetPublishedRelease(ctx, period.Tahun)
	if err != nil {
		return nil, nil, err
	}
	wf, err := s.master.GetWorkflowVersion(ctx, "wf-default-v1")
	if err != nil {
		return nil, nil, err
	}
	data, err := io.ReadAll(io.LimitReader(in.File, s.maxSize+1))
	if err != nil {
		return nil, nil, err
	}
	if int64(len(data)) > s.maxSize {
		return nil, nil, domain.NewError(http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "ukuran file melebihi batas")
	}
	ext := strings.ToLower(filepath.Ext(in.Filename))
	if ext != ".csv" && ext != ".xlsx" {
		return nil, nil, domain.Validation("format file harus .csv atau .xlsx", nil)
	}
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])
	relPath := filepath.Join(fmt.Sprintf("%04d", period.Tahun), fmt.Sprintf("%02d", period.Bulan), hash+ext)
	absPath := filepath.Join(s.dataDir, relPath)
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return nil, nil, err
	}
	if err := os.WriteFile(absPath, data, 0o644); err != nil {
		return nil, nil, err
	}
	mimeType := http.DetectContentType(data)
	if ext == ".csv" {
		mimeType = "text/csv"
	}
	sub, job, err := s.subs.CreateUpload(ctx, domain.CreateUploadParams{
		SubmissionID: uuid.NewString(), PeriodKode: period.Kode, UploaderID: actor.ID,
		FaskesID: faskesID, PeriodeID: period.ID, MasterReleaseID: release.ID,
		WorkflowVersionID: wf.ID, IdempotencyKey: in.IdempotencyKey,
		Now: now(),
		File: domain.FileObject{
			ID: uuid.NewString(), OriginalName: filepath.Base(in.Filename), SHA256: hash,
			Path: absPath, Mime: &mimeType, Size: int64(len(data)), UploadedBy: &actor.ID,
			ScanStatus: "clean", CreatedAt: now(),
		},
	})
	if err != nil {
		return nil, nil, err
	}
	_ = s.audit.Write(ctx, &actor.ID, "submission.created", "submission", sub.ID, nil, sub, "upload BNBA", in.RequestID)
	return sub, job, nil
}

type SubmissionDetail struct {
	Submission *domain.Submission          `json:"submission"`
	Rows       []domain.SubmissionRow      `json:"rows"`
	Reviews    []domain.SubmissionReview   `json:"reviews"`
	Revisions  []domain.SubmissionRevision `json:"revisions"`
	AllowedIKU map[string][]domain.Label   `json:"allowed_iku"`
	Blockers   int                         `json:"blockers"`
}

func (s *SubmissionService) Detail(ctx context.Context, id string) (*SubmissionDetail, error) {
	sub, err := s.subs.GetSubmission(ctx, id)
	if err != nil {
		return nil, err
	}
	rows, err := s.subs.GetRows(ctx, id)
	if err != nil {
		return nil, err
	}
	reviews, err := s.subs.ListReviews(ctx, id)
	if err != nil {
		return nil, err
	}
	revisions, err := s.subs.ListRevisions(ctx, id)
	if err != nil {
		return nil, err
	}
	allowed := map[string][]domain.Label{}
	for _, row := range rows {
		if row.MatchedIndikatorID == nil {
			continue
		}
		labels, _ := s.allowedIKUs(ctx, *row.MatchedIndikatorID, sub.MasterReleaseID)
		allowed[row.ID] = labels
	}
	blockers, _ := s.subs.CountBlockingRows(ctx, id)
	return &SubmissionDetail{Submission: sub, Rows: rows, Reviews: reviews, Revisions: revisions, AllowedIKU: allowed, Blockers: blockers}, nil
}

func (s *SubmissionService) allowedIKUs(ctx context.Context, indikatorID string, releaseID *string) ([]domain.Label, error) {
	links, err := s.master.ListIndicatorLabels(ctx, indikatorID)
	if err != nil {
		return nil, err
	}
	labels, err := s.master.ListLabels(ctx)
	if err != nil {
		return nil, err
	}
	byID := map[string]domain.Label{}
	for _, l := range labels {
		byID[l.ID] = l
	}
	var out []domain.Label
	for _, link := range links {
		l, ok := byID[link.LabelID]
		if !ok || !l.AllowDirectInput || l.IsRollupOnly || !l.Aktif {
			continue
		}
		if releaseID != nil && l.MasterReleaseID != *releaseID {
			continue
		}
		out = append(out, l)
	}
	return out, nil
}

type UpdateRowInput struct {
	IKUKode string `json:"iku_kode"`
	Reason  string `json:"reason"`
}

func (s *SubmissionService) UpdateRow(ctx context.Context, actor *domain.User, submissionID, rowID string, in UpdateRowInput, requestID string) (*domain.SubmissionRow, error) {
	sub, err := s.subs.GetSubmission(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if sub.Status != domain.StatusDraf && sub.Status != domain.StatusDikembalikanPetugas && sub.Status != domain.StatusDikembalikanVerifikator {
		return nil, domain.Conflict("perubahan IKU hanya boleh sebelum submit")
	}
	if actor.Role != domain.RoleAdmin && actor.ID != sub.UploaderID {
		return nil, domain.Forbidden("hanya uploader yang boleh mengubah baris ini")
	}
	row, err := s.subs.GetRow(ctx, rowID)
	if err != nil {
		return nil, err
	}
	if row.SubmissionID != submissionID || row.MatchedIndikatorID == nil {
		return nil, domain.Validation("baris tidak valid", nil)
	}
	label, err := s.master.GetLabelByKode(ctx, strings.TrimSpace(in.IKUKode), deref(sub.MasterReleaseID))
	if err != nil {
		return nil, err
	}
	if !label.AllowDirectInput || label.IsRollupOnly {
		return nil, domain.Validation("IKU induk/rollup tidak boleh menerima input langsung", map[string]any{"iku_kode": in.IKUKode})
	}
	links, err := s.master.ListIndicatorLabels(ctx, *row.MatchedIndikatorID)
	if err != nil {
		return nil, err
	}
	ok := false
	for _, link := range links {
		if link.LabelID == label.ID {
			ok = true
			break
		}
	}
	if !ok {
		return nil, domain.Validation("IKU tidak terhubung ke indikator baris ini", map[string]any{"iku_kode": in.IKUKode})
	}
	before := row
	if err := s.subs.UpdateRowClassification(ctx, rowID, label.ID, domain.ClassMatched, ""); err != nil {
		return nil, err
	}
	after, _ := s.subs.GetRow(ctx, rowID)
	if after != nil {
		if err := s.subs.UpsertIndicatorValueForRow(ctx, after, sub); err != nil {
			return nil, err
		}
	}
	_ = s.audit.Write(ctx, &actor.ID, "submission.row_reclassified", "submission_row", rowID, before, after, in.Reason, requestID)
	return after, nil
}

func (s *SubmissionService) Submit(ctx context.Context, actor *domain.User, id, requestID string) (*domain.Submission, error) {
	sub, err := s.subs.GetSubmission(ctx, id)
	if err != nil {
		return nil, err
	}
	if actor.Role != domain.RoleAdmin && actor.ID != sub.UploaderID {
		return nil, domain.Forbidden("hanya uploader yang boleh mengirim submission")
	}
	if sub.Status != domain.StatusDraf && sub.Status != domain.StatusDikembalikanPetugas && sub.Status != domain.StatusDikembalikanVerifikator {
		return nil, domain.Conflict("submission tidak dalam status draf/dikembalikan")
	}
	blockers, err := s.subs.CountBlockingRows(ctx, id)
	if err != nil {
		return nil, err
	}
	if blockers > 0 {
		return nil, domain.Validation("masih ada baris yang perlu dikoreksi", map[string]any{"blockers": blockers})
	}
	if err := s.subs.UpdateSubmissionStatus(ctx, id, domain.StatusDikirim, domain.StepVerifyDinas); err != nil {
		return nil, err
	}
	_ = s.subs.CreateReview(ctx, &domain.SubmissionReview{
		SubmissionID: id, Step: string(domain.StepDraftPKM), ReviewerID: &actor.ID,
		Action: "submit", Notes: strPtr("Dikirim ke verifikasi Dinas"), CreatedAt: now(),
	})
	updated, _ := s.subs.GetSubmission(ctx, id)
	_ = s.audit.Write(ctx, &actor.ID, "submission.submitted", "submission", id, sub, updated, "", requestID)
	return updated, nil
}

type DecisionInput struct {
	Action     string `json:"action"`
	Notes      string `json:"notes"`
	ReasonCode string `json:"reason_code"`
	RequestID  string `json:"-"`
}

func (s *SubmissionService) Decision(ctx context.Context, actor *domain.User, id string, in DecisionInput) (*domain.Submission, error) {
	sub, err := s.subs.GetSubmission(ctx, id)
	if err != nil {
		return nil, err
	}
	action := strings.ToLower(strings.TrimSpace(in.Action))
	if action != "approve" && action != "return" {
		return nil, domain.Validation("action harus approve atau return", nil)
	}
	notes := strings.TrimSpace(in.Notes)
	if action == "return" && notes == "" {
		return nil, domain.Validation("catatan wajib diisi saat return", nil)
	}
	var nextStatus domain.SubmissionStatus
	var nextStep domain.WorkflowStep
	var step string
	switch sub.Status {
	case domain.StatusDikirim:
		if actor.Role != domain.RoleVerifier && actor.Role != domain.RoleAdmin {
			return nil, domain.Forbidden("hanya verifikator Dinas yang boleh memutuskan tahap ini")
		}
		step = string(domain.StepVerifyDinas)
		if action == "approve" {
			nextStatus, nextStep = domain.StatusDiverifikasi, domain.StepValidateKabid
		} else {
			nextStatus, nextStep = domain.StatusDikembalikanPetugas, domain.StepDraftPKM
		}
	case domain.StatusDiverifikasi:
		if actor.Role != domain.RoleKabid && actor.Role != domain.RoleSekretaris && actor.Role != domain.RoleAdmin {
			return nil, domain.Forbidden("hanya Kabid/Sekretaris yang boleh memutuskan tahap ini")
		}
		step = string(domain.StepValidateKabid)
		if action == "approve" {
			if notes == "" {
				return nil, domain.Validation("catatan penerimaan wajib diisi saat final", nil)
			}
			nextStatus, nextStep = domain.StatusFinal, domain.StepFinal
		} else {
			nextStatus, nextStep = domain.StatusDikembalikanVerifikator, domain.StepVerifyDinas
		}
	default:
		return nil, domain.Conflict("submission tidak berada pada tahap approval")
	}
	if err := s.subs.UpdateSubmissionStatus(ctx, id, nextStatus, nextStep); err != nil {
		return nil, err
	}
	_ = s.subs.CreateReview(ctx, &domain.SubmissionReview{
		SubmissionID: id, Step: step, ReviewerID: &actor.ID, Action: action,
		Notes: strPtr(notes), DetailsJSON: strPtr(in.ReasonCode), CreatedAt: now(),
	})
	updated, _ := s.subs.GetSubmission(ctx, id)
	_ = s.audit.Write(ctx, &actor.ID, "submission."+action, "submission", id, sub, updated, notes, in.RequestID)
	if nextStatus == domain.StatusFinal {
		_, _ = s.calc.CalculatePeriod(ctx, sub.PeriodeID, "submission.final")
	}
	return updated, nil
}

func (s *SubmissionService) List(ctx context.Context, f domain.SubmissionFilter) ([]domain.Submission, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 100
	}
	return s.subs.ListSubmissions(ctx, f)
}

func (s *SubmissionService) Revisions(ctx context.Context, id string) ([]domain.SubmissionRevision, error) {
	return s.subs.ListRevisions(ctx, id)
}

func (s *SubmissionService) Reviews(ctx context.Context, id string) ([]domain.SubmissionReview, error) {
	return s.subs.ListReviews(ctx, id)
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
