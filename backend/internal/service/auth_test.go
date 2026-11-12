package service_test

import (
	"context"
	"testing"

	"matasidik/backend/internal/config"
	"matasidik/backend/internal/service"
)

func TestPasetoLoginRefreshAndActingAs(t *testing.T) {
	ctx := context.Background()
	store := newTestStore(t)
	cfg := config.Load()
	auth, err := service.NewAuthService(cfg, store, store)
	if err != nil {
		t.Fatal(err)
	}
	result, err := auth.Login(ctx, "siti", "password123", "test")
	if err != nil {
		t.Fatal(err)
	}
	if result.User.Username != "siti" || result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatalf("login tidak lengkap: %+v", result)
	}
	parsed, err := auth.Parse(ctx, result.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.User.ID != result.User.ID {
		t.Fatalf("subject token salah: %s", parsed.User.ID)
	}
	refreshed, err := auth.Refresh(ctx, result.RefreshToken, "test")
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.User.ID != result.User.ID || refreshed.RefreshToken == result.RefreshToken {
		t.Fatalf("refresh tidak rotasi: %+v", refreshed)
	}

	adminLogin, err := auth.Login(ctx, "admin", "password123", "test")
	if err != nil {
		t.Fatal(err)
	}
	acting, err := auth.ActingAs(ctx, adminLogin.User, result.User.ID)
	if err != nil {
		t.Fatal(err)
	}
	ac, err := auth.Parse(ctx, acting.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if ac.User.ID != result.User.ID || ac.Impersonator == nil || ac.Impersonator.ID != adminLogin.User.ID {
		t.Fatalf("actingAs tidak valid: %+v", ac)
	}
}
