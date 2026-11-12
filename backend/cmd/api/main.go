package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"matasidik/backend/internal/config"
	"matasidik/backend/internal/httpapi"
	"matasidik/backend/internal/httpapi/handler"
	"matasidik/backend/internal/repository/sqlite"
	"matasidik/backend/internal/service"
	"matasidik/backend/internal/worker"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	cfg := config.Load()
	store, err := sqlite.New(cfg.DSN(), cfg.DataDir)
	if err != nil {
		slog.Error("buka database", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	ctx := context.Background()
	audit := service.NewAuditService(store)
	auth, err := service.NewAuthService(cfg, store, store)
	if err != nil {
		slog.Error("inisialisasi auth", "error", err)
		os.Exit(1)
	}
	admin := service.NewAdminService(store)
	master := service.NewMasterService(store)
	calc := service.NewCalculationService(store, store)

	seeded, err := store.SeedDemo(ctx)
	if err != nil {
		slog.Error("seed demo", "error", err)
		os.Exit(1)
	}
	if seeded {
		if _, err := calc.CalculatePeriod(ctx, "per-2026-05", "seed"); err != nil {
			slog.Error("kalkulasi seed", "error", err)
		}
	}

	subs := service.NewSubmissionService(store, store, audit, calc, cfg.DataDir, cfg.MaxUploadBytes())
	dashboard := service.NewDashboardService(store, store, store)
	factor := service.NewFactorService(store, audit)
	h := handler.New(cfg, store, auth, admin, master, subs, calc, dashboard, factor, audit)
	router := httpapi.NewRouter(cfg, h)

	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if cfg.WorkerInline {
		host, _ := os.Hostname()
		go worker.New("inline-"+host, store, store, store).Run(runCtx)
		slog.Info("worker inline aktif")
	}

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	go func() {
		slog.Info("API berjalan", "addr", cfg.Addr, "web_dir", cfg.WebDir)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server", "error", err)
			cancel()
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	slog.Info("API berhenti")
}
