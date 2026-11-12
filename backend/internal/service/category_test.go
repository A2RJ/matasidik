package service_test

import (
	"context"
	"testing"

	"matasidik/backend/internal/service"
)

func TestCategoryMasterAndIndicatorLabels(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	admin, err := store.GetUser(ctx, "usr-admin")
	if err != nil {
		t.Fatal(err)
	}
	master := service.NewMasterService(store)

	categories, err := master.Categories(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(categories) < 8 {
		t.Fatalf("kategori seed kurang: %d", len(categories))
	}
	var ikkCategoryID string
	for _, category := range categories {
		if category.Kode == "IKK" {
			ikkCategoryID = category.ID
		}
	}
	if ikkCategoryID == "" {
		t.Fatal("kategori IKK tidak ditemukan")
	}

	ikkLabel, err := master.CreateLabel(ctx, service.CreateLabelInput{
		Kode:              "IKK-TEST-1",
		Nama:              "IKK Test",
		CategoryID:        ikkCategoryID,
		MasterReleaseID:   "mr-2026",
		AggregationMethod: "RATIO_OF_SUMS",
		AllowDirectInput:  true,
	}, admin)
	if err != nil {
		t.Fatal(err)
	}
	if ikkLabel.CategoryKode == nil || *ikkLabel.CategoryKode != "IKK" {
		t.Fatalf("kategori label salah: %+v", ikkLabel)
	}

	indicator, err := master.CreateIndicator(ctx, service.CreateIndicatorInput{
		Kode:             "IND-TEST-1",
		Nama:             "Indikator Test",
		Satuan:           "%",
		SifatInput:       "puskesmas",
		Frekuensi:        "bulanan",
		PJBidangID:       "bid-kesmas",
		PJTimkerID:       "bid-kesga",
		LabelIDs:         []string{"lbl-iku-001-01", ikkLabel.ID},
		PrimaryLabelID:   ikkLabel.ID,
		NumeratorLabel:   "pembilang",
		DenominatorLabel: "penyebut",
	}, admin)
	if err != nil {
		t.Fatal(err)
	}
	if indicator.PJTimkerID == nil || *indicator.PJTimkerID != "bid-kesga" {
		t.Fatalf("PJ timker tidak tersimpan: %+v", indicator)
	}
	if indicator.PJTimker == nil || *indicator.PJTimker != "Timker Kesehatan Keluarga" {
		t.Fatalf("nama PJ timker tidak di-join: %+v", indicator)
	}

	links, err := store.ListIndicatorLabels(ctx, indicator.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 2 {
		t.Fatalf("jumlah label indikator: %d", len(links))
	}
	primary := 0
	for _, link := range links {
		if link.IsPrimary {
			primary++
			if link.LabelID != ikkLabel.ID {
				t.Fatalf("label utama salah: %s", link.LabelID)
			}
		}
	}
	if primary != 1 {
		t.Fatalf("label utama harus tepat satu: %d", primary)
	}
	if len(indicator.LabelIDs) != 2 || indicator.PrimaryLabelID == nil || *indicator.PrimaryLabelID != ikkLabel.ID {
		t.Fatalf("label indikator tidak lengkap pada response: %+v", indicator)
	}
	listed, err := master.Indicators(ctx)
	if err != nil {
		t.Fatal(err)
	}
	foundPreview := false
	for _, item := range listed {
		if item.ID == indicator.ID {
			foundPreview = true
			if len(item.LabelIDs) != 2 || item.PrimaryLabelID == nil {
				t.Fatalf("preview indikator tidak membawa label: %+v", item)
			}
		}
	}
	if !foundPreview {
		t.Fatal("indikator test tidak ada di daftar preview")
	}

	dashboard := service.NewDashboardService(store, store, store)
	ikkResult, err := dashboard.Category(ctx, "IKK", domainFilter("per-2026-05"))
	if err != nil {
		t.Fatal(err)
	}
	if len(ikkResult.Items) != 1 || ikkResult.Items[0].LabelKode == nil || *ikkResult.Items[0].LabelKode != "IKK-TEST-1" {
		t.Fatalf("dashboard IKK tidak memuat label: %+v", ikkResult.Items)
	}

	ownerResult, err := dashboard.Category(ctx, "BIDANG_TIMKER", domainFilter("per-2026-05"))
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, item := range ownerResult.Items {
		if item.Kode == "IND-TEST-1" {
			found = true
			if item.OwnerNama == nil || *item.OwnerNama != "Timker Kesehatan Keluarga" {
				t.Fatalf("owner indikator salah: %+v", item)
			}
		}
	}
	if !found {
		t.Fatal("indikator test tidak muncul di kategori owner")
	}

	allResult, err := dashboard.Category(ctx, "ALL", domainFilter("per-2026-05"))
	if err != nil {
		t.Fatal(err)
	}
	if len(allResult.Items) == 0 {
		t.Fatal("filter Semua tidak mengembalikan baris")
	}
	foundLabel, foundIndicator := false, false
	for _, item := range allResult.Items {
		if item.LabelKode != nil && *item.LabelKode == "IKK-TEST-1" {
			foundLabel = true
		}
		if item.Kode == "IND-TEST-1" {
			foundIndicator = true
		}
	}
	if !foundLabel || !foundIndicator {
		t.Fatalf("filter Semua tidak lengkap: label=%v indikator=%v", foundLabel, foundIndicator)
	}
}

func TestCreateCategoryAndDeactivate(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	admin, err := store.GetUser(ctx, "usr-admin")
	if err != nil {
		t.Fatal(err)
	}
	master := service.NewMasterService(store)

	created, err := master.CreateCategory(ctx, service.CategoryInput{
		Kode: "TEST_CAT", Nama: "Kategori Test", Jenis: "LABEL", Urutan: 99,
	}, admin)
	if err != nil {
		t.Fatal(err)
	}
	if created.Kode != "TEST_CAT" || !created.Aktif {
		t.Fatalf("kategori baru tidak valid: %+v", created)
	}
	updated, err := master.UpdateCategory(ctx, created.ID, service.CategoryInput{
		Nama: "Kategori Test Baru", Urutan: 100, Aktif: boolPtr(false),
	}, admin)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Nama != "Kategori Test Baru" || updated.Aktif {
		t.Fatalf("kategori tidak ter-update: %+v", updated)
	}
}

func TestDeactivateCategoryWithLabelIsSoftDelete(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	admin, err := store.GetUser(ctx, "usr-admin")
	if err != nil {
		t.Fatal(err)
	}
	master := service.NewMasterService(store)

	category, err := master.CreateCategory(ctx, service.CategoryInput{
		Kode: "SOFT_CAT", Nama: "Kategori Soft Delete", Jenis: "LABEL", Urutan: 50,
	}, admin)
	if err != nil {
		t.Fatal(err)
	}
	label, err := master.CreateLabel(ctx, service.CreateLabelInput{
		Kode: "SOFT-LABEL-1", Nama: "Label Soft Delete", CategoryID: category.ID,
		MasterReleaseID: "mr-2026", AggregationMethod: "RATIO_OF_SUMS", AllowDirectInput: true,
	}, admin)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := master.DeactivateCategory(ctx, category.ID, admin); err != nil {
		t.Fatal(err)
	}

	categories, err := master.Categories(ctx)
	if err != nil {
		t.Fatal(err)
	}
	foundCategory := false
	for _, c := range categories {
		if c.ID == category.ID {
			foundCategory = true
			if c.Aktif {
				t.Fatal("kategori masih aktif setelah soft delete")
			}
		}
	}
	if !foundCategory {
		t.Fatal("kategori historis hilang setelah soft delete")
	}
	labels, err := store.ListLabels(ctx)
	if err != nil {
		t.Fatal(err)
	}
	foundLabel := false
	for _, l := range labels {
		if l.ID == label.ID {
			foundLabel = true
			if !l.Aktif {
				t.Fatal("label ikut terhapus saat kategori di-soft-delete")
			}
		}
	}
	if !foundLabel {
		t.Fatal("label historis hilang setelah kategori di-soft-delete")
	}

	if _, err := master.CreateLabel(ctx, service.CreateLabelInput{
		Kode: "SOFT-LABEL-2", Nama: "Label Baru", CategoryID: category.ID,
		MasterReleaseID: "mr-2026", AggregationMethod: "RATIO_OF_SUMS", AllowDirectInput: true,
	}, admin); err == nil {
		t.Fatal("label baru masih bisa memakai kategori nonaktif")
	}
	if _, err := master.CreateIndicator(ctx, service.CreateIndicatorInput{
		Kode: "IND-SOFT-1", Nama: "Indikator Soft", Satuan: "%",
		SifatInput: "puskesmas", Frekuensi: "bulanan", LabelIDs: []string{label.ID},
	}, admin); err == nil {
		t.Fatal("indikator baru masih bisa memakai label dari kategori nonaktif")
	}
}

func boolPtr(v bool) *bool { return &v }
