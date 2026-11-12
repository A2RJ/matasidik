package service_test

import (
	"context"
	"testing"

	"matasidik/backend/internal/config"
	"matasidik/backend/internal/service"
)

func TestProfileAndPasswordSelfService(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	auth, err := service.NewAuthService(config.Load(), store, store)
	if err != nil {
		t.Fatal(err)
	}
	login, err := auth.Login(ctx, "siti", "password123", "test")
	if err != nil {
		t.Fatal(err)
	}

	updated, err := auth.UpdateProfile(ctx, login.User, service.UpdateProfileInput{Nama: "Siti Nurhaliza"})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Nama != "Siti Nurhaliza" || updated.Username != "siti" || updated.FaskesID == nil {
		t.Fatalf("profil tidak diperbarui dengan benar: %+v", updated)
	}

	if err := auth.ChangePassword(ctx, updated, service.ChangePasswordInput{
		CurrentPassword: "salah-password", NewPassword: "passwordBaru123",
	}); err == nil {
		t.Fatal("password lama salah seharusnya ditolak")
	}
	if err := auth.ChangePassword(ctx, updated, service.ChangePasswordInput{
		CurrentPassword: "password123", NewPassword: "password123",
	}); err == nil {
		t.Fatal("password baru yang sama seharusnya ditolak")
	}
	if err := auth.ChangePassword(ctx, updated, service.ChangePasswordInput{
		CurrentPassword: "password123", NewPassword: "passwordBaru123",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := auth.Login(ctx, "siti", "password123", "test"); err == nil {
		t.Fatal("login dengan password lama seharusnya gagal")
	}
	newLogin, err := auth.Login(ctx, "siti", "passwordBaru123", "test")
	if err != nil {
		t.Fatalf("login dengan password baru gagal: %v", err)
	}
	if newLogin.User.Nama != "Siti Nurhaliza" {
		t.Fatalf("nama hasil update tidak terbaca: %+v", newLogin.User)
	}
	if _, err := auth.Refresh(ctx, login.RefreshToken, "test"); err == nil {
		t.Fatal("refresh token lama seharusnya dicabut setelah ganti password")
	}
}
