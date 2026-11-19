package service_test

import (
	"context"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/service"
)

func TestSubmissionPuskesmasApprovalStep(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	if _, err := store.SeedDemo(ctx); err != nil {
		t.Fatal(err)
	}
	approver, err := store.GetUser(ctx, "usr-kepala-pkm")
	if err != nil {
		t.Fatal(err)
	}
	admin, err := store.GetUser(ctx, "usr-admin")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission
		(id, kode, tipe, uploader_id, faskes_id, periode_id, file_id, file_hash, status, step,
		 master_release_id, workflow_version_id, created_at, updated_at)
		SELECT 'sub-test-pkm', 'S-TEST-PKM', 'BNBA', uploader_id, faskes_id, periode_id, file_id,
		 file_hash, 'dikirim', 'APPROVE_PKM', master_release_id, 'wf-default-v2',
		 '2026-05-10T08:00:00Z', '2026-05-10T08:00:00Z'
		FROM submission WHERE id = 'sub-demo-001'`); err != nil {
		t.Fatal(err)
	}
	audit := service.NewAuditService(store)
	calc := service.NewCalculationService(store, store)
	subs := service.NewSubmissionService(store, store, audit, calc, t.TempDir(), 15<<20)

	if _, err := subs.Decision(ctx, admin, "sub-test-pkm", service.DecisionInput{Action: "approve"}); err == nil {
		t.Fatal("admin tidak boleh melewati approval Puskesmas")
	}
	updated, err := subs.Decision(ctx, approver, "sub-test-pkm", service.DecisionInput{Action: "approve"})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status != domain.StatusDikirim || updated.Step != domain.StepVerifyDinas {
		t.Fatalf("setelah approval Puskesmas: %s/%s", updated.Status, updated.Step)
	}
	reviews, err := store.ListReviews(ctx, "sub-test-pkm")
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 1 || reviews[0].Step != "APPROVE_PKM" {
		t.Fatalf("review Puskesmas tidak tercatat: %+v", reviews)
	}
}

func TestSubmissionPuskesmasReturn(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	if _, err := store.SeedDemo(ctx); err != nil {
		t.Fatal(err)
	}
	approver, err := store.GetUser(ctx, "usr-kepala-pkm")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.DB().ExecContext(ctx, `INSERT INTO submission
		(id, kode, tipe, uploader_id, faskes_id, periode_id, file_id, file_hash, status, step,
		 master_release_id, workflow_version_id, created_at, updated_at)
		SELECT 'sub-test-return', 'S-TEST-RETURN', 'BNBA', uploader_id, faskes_id, periode_id, file_id,
		 file_hash, 'dikirim', 'APPROVE_PKM', master_release_id, 'wf-default-v2',
		 '2026-05-10T08:00:00Z', '2026-05-10T08:00:00Z'
		FROM submission WHERE id = 'sub-demo-001'`); err != nil {
		t.Fatal(err)
	}
	audit := service.NewAuditService(store)
	calc := service.NewCalculationService(store, store)
	subs := service.NewSubmissionService(store, store, audit, calc, t.TempDir(), 15<<20)
	updated, err := subs.Decision(ctx, approver, "sub-test-return", service.DecisionInput{
		Action: "return", Notes: "Data belum lengkap",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status != domain.StatusDikembalikanPetugas || updated.Step != domain.StepDraftPKM {
		t.Fatalf("hasil return Puskesmas: %s/%s", updated.Status, updated.Step)
	}
}
