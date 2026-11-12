package service_test

import (
	"context"
	"testing"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/service"
)

func TestAdminPositionChangesEffectiveRole(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	admin := service.NewAdminService(store)
	adminUser, err := store.GetUser(ctx, "usr-admin")
	if err != nil {
		t.Fatal(err)
	}
	rina, err := store.GetUser(ctx, "usr-rina")
	if err != nil {
		t.Fatal(err)
	}
	if rina.Role != domain.RoleVerifier {
		t.Fatalf("role awal rina: %s", rina.Role)
	}
	updated, err := admin.UpdateUser(ctx, rina.ID, service.UpdateUserInput{
		Nama: rina.Nama, Username: rina.Username, PositionID: "pos-kabid-kesmas",
		BidangID: "bid-kesmas",
	}, adminUser)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Role != domain.RoleKabid || updated.PositionNama == nil {
		t.Fatalf("role efektif tidak berubah: %+v", updated)
	}
}
