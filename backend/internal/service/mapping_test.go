package service_test

import (
	"context"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/service"
)

func TestDesaFaskesMappingVersioning(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	if seeded, err := store.SeedDemo(ctx); err != nil || !seeded {
		t.Fatalf("seed: %v seeded=%v", err, seeded)
	}
	master := service.NewMasterService(store)
	admin := &domain.User{ID: "usr-admin", Role: domain.RoleAdmin}

	mappings, err := master.DesaFaskesMappings(ctx)
	if err != nil || len(mappings) != 2 {
		t.Fatalf("mapping awal: %+v err=%v", mappings, err)
	}

	updated, err := master.UpsertDesaFaskesMapping(ctx, service.DesaFaskesMappingInput{
		DesaID: "wil-desa-dalam", FaskesID: "faskes-pkm-jereweh", BerlakuDari: "2026-07-01",
	}, admin)
	if err != nil {
		t.Fatal(err)
	}
	if updated == nil || updated.FaskesID != "faskes-pkm-jereweh" || !updated.Aktif {
		t.Fatalf("mapping baru salah: %+v", updated)
	}
	taliwang, err := store.ListActiveDesaIDsForFaskes(ctx, "faskes-pkm-taliwang")
	if err != nil {
		t.Fatal(err)
	}
	if len(taliwang) != 1 || taliwang[0] != "wil-desa-beru" {
		t.Fatalf("desa Taliwang setelah pindah: %+v", taliwang)
	}
	jereweh, err := store.ListActiveDesaIDsForFaskes(ctx, "faskes-pkm-jereweh")
	if err != nil {
		t.Fatal(err)
	}
	if len(jereweh) != 1 || jereweh[0] != "wil-desa-dalam" {
		t.Fatalf("desa Jereweh: %+v", jereweh)
	}

	if err := master.DeactivateDesaFaskesMapping(ctx, updated.ID, admin); err != nil {
		t.Fatal(err)
	}
	active, err := store.GetActiveFaskesForDesa(ctx, "wil-desa-dalam")
	if err != nil {
		t.Fatal(err)
	}
	if active != nil {
		t.Fatalf("mapping masih aktif: %+v", active)
	}
}
