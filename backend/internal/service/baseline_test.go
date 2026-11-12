package service_test

import (
	"context"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/service"
)

func TestCapaianTahunanCRUDAndDashboardBaseline(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	master := service.NewMasterService(store)
	admin, err := store.GetUser(ctx, "usr-admin")
	if err != nil {
		t.Fatal(err)
	}
	siti, err := store.GetUser(ctx, "usr-siti")
	if err != nil {
		t.Fatal(err)
	}

	list, err := master.CapaianTahunan(ctx, 2025)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 12 {
		t.Fatalf("seed capaian tahunan = %d, want 12", len(list))
	}
	target := list[0]
	if target.LabelID == nil {
		t.Fatal("capaian seed tidak punya label")
	}
	value := 88.5
	updated, err := master.UpdateCapaianTahunan(ctx, target.ID, service.CapaianTahunanInput{
		LabelID: *target.LabelID, Tahun: 2025, Nilai: &value, Catatan: "koreksi baseline",
	}, admin)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Nilai != value || updated.Catatan == nil {
		t.Fatalf("update capaian gagal: %+v", updated)
	}

	dashboard := service.NewDashboardService(store, store, store)
	items, _, err := dashboard.SPM(ctx, domain.DashboardFilter{PeriodeID: "per-2026-05"})
	if err != nil {
		t.Fatal(err)
	}
	var found *domain.DashboardSPMItem
	for i := range items {
		if items[i].LabelID != nil && *items[i].LabelID == *target.LabelID {
			found = &items[i]
			break
		}
	}
	if found == nil || found.BaselinePrev == nil || *found.BaselinePrev != value {
		t.Fatalf("dashboard tidak memakai baseline baru: %+v", found)
	}
	if found.BaselineYear != 2025 {
		t.Fatalf("baseline_year = %d, want 2025", found.BaselineYear)
	}

	filtered, _, err := dashboard.SPM(ctx, domain.DashboardFilter{
		PeriodeID: "per-2026-05", FaskesID: "faskes-pkm-taliwang",
	})
	if err != nil {
		t.Fatal(err)
	}
	for i := range filtered {
		if filtered[i].LabelID != nil && *filtered[i].LabelID == *target.LabelID {
			if filtered[i].BaselinePrev == nil || *filtered[i].BaselinePrev != value {
				t.Fatalf("baseline hilang saat filter faskes: %+v", filtered[i])
			}
		}
	}

	if err := master.DeleteCapaianTahunan(ctx, updated.ID, admin); err != nil {
		t.Fatal(err)
	}
	items, _, err = dashboard.SPM(ctx, domain.DashboardFilter{PeriodeID: "per-2026-05"})
	if err != nil {
		t.Fatal(err)
	}
	for i := range items {
		if items[i].LabelID != nil && *items[i].LabelID == *target.LabelID && items[i].BaselinePrev != nil {
			t.Fatal("baseline masih muncul setelah dihapus")
		}
	}

	if _, err := master.CreateCapaianTahunan(ctx, service.CapaianTahunanInput{
		LabelID: *target.LabelID, Tahun: 2025, Nilai: &value,
	}, siti); err == nil {
		t.Fatal("non-admin seharusnya tidak boleh membuat capaian")
	}
}
