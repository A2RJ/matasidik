package service_test

import (
	"context"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/service"
)

func TestIKUChangeApprovalFlow(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	operator, err := store.GetUser(ctx, "usr-siti")
	if err != nil {
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
	iku := service.NewIKUChangeService(store, service.NewAuditService(store))

	request, err := iku.Create(ctx, operator, service.IKUChangeInput{
		Action: "CREATE", Kode: "IKU-TEST-NEW", Nama: "IKU Test Baru",
		CategoryID: "cat-iku", AggregationMethod: "RATIO_OF_SUMS",
		AllowDirectInput: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if request.Status != service.IKUStatusMenungguPuskesmas {
		t.Fatalf("status awal: %s", request.Status)
	}

	list, err := iku.List(ctx, approver, domain.IKUChangeFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("antrean Puskesmas: %d", len(list))
	}
	if _, err := iku.Decide(ctx, approver, request.ID, service.IKUChangeDecisionInput{Action: "approve"}); err != nil {
		t.Fatal(err)
	}
	pendingAdmin, err := store.GetIKUChangeRequest(ctx, request.ID)
	if err != nil {
		t.Fatal(err)
	}
	if pendingAdmin.Status != service.IKUStatusMenungguAdmin {
		t.Fatalf("status setelah Puskesmas: %s", pendingAdmin.Status)
	}

	if _, err := iku.Decide(ctx, admin, request.ID, service.IKUChangeDecisionInput{
		Action: "approve", Notes: "Disetujui Admin",
	}); err != nil {
		t.Fatal(err)
	}
	done, err := store.GetIKUChangeRequest(ctx, request.ID)
	if err != nil {
		t.Fatal(err)
	}
	if done.Status != service.IKUStatusDisetujui || done.AppliedLabelID == nil {
		t.Fatalf("pengajuan belum dipublikasikan: %+v", done)
	}
	labels, err := store.ListLabels(ctx)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, label := range labels {
		if label.Kode == "IKU-TEST-NEW" && label.Aktif {
			found = true
		}
	}
	if !found {
		t.Fatal("label hasil approval tidak ditemukan")
	}
	reviews, err := store.ListIKUChangeReviews(ctx, request.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 2 || reviews[0].Step != "PKM" || reviews[1].Step != "ADMIN" {
		t.Fatalf("review tidak lengkap: %+v", reviews)
	}
}

func TestIKUChangeRejectsSelfApproval(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	operator, err := store.GetUser(ctx, "usr-siti")
	if err != nil {
		t.Fatal(err)
	}
	iku := service.NewIKUChangeService(store, service.NewAuditService(store))
	request, err := iku.Create(ctx, operator, service.IKUChangeInput{
		Action: "CREATE", Kode: "IKU-TEST-SELF", Nama: "IKU Self",
		CategoryID: "cat-iku", AggregationMethod: "RATIO_OF_SUMS", AllowDirectInput: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	// Operator bukan pkm_approver, jadi tidak bisa memutus tahap Puskesmas.
	if _, err := iku.Decide(ctx, operator, request.ID, service.IKUChangeDecisionInput{Action: "approve"}); err == nil {
		t.Fatal("operator tidak boleh approve pengajuannya sendiri")
	}
}
