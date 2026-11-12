package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"matasidik/backend/internal/config"
	"matasidik/backend/internal/repository/sqlite"
	"matasidik/backend/internal/worker"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	cfg := config.Load()
	store, err := sqlite.New(cfg.DSN(), cfg.DataDir)
	if err != nil {
		slog.Error("buka database", "error", err)
		os.Exit(1)
	}
	defer store.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	host, _ := os.Hostname()
	w := worker.New("worker-"+host, store, store, store)
	go w.Run(ctx)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	cancel()
	slog.Info("worker berhenti")
}
