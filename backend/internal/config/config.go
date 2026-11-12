package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Addr         string
	DBPath       string
	DataDir      string
	WebDir       string
	WorkerInline bool
	MaxUploadMB  int64
	AuthKeyHex   string
	CookieSecure bool
}

func Load() Config {
	return Config{
		Addr:         env("APP_ADDR", ":8080"),
		DBPath:       env("APP_DB", "data/matasidik.db"),
		DataDir:      env("APP_DATA_DIR", "data/files"),
		WebDir:       env("APP_WEB_DIR", "../frontend/build"),
		WorkerInline: envBool("APP_WORKER_INLINE", true),
		MaxUploadMB:  envInt("APP_MAX_UPLOAD_MB", 15),
		AuthKeyHex:   strings.TrimSpace(os.Getenv("APP_AUTH_KEY")),
		CookieSecure: envBool("APP_COOKIE_SECURE", false),
	}
}

func (c Config) DSN() string {
	path := strings.TrimSpace(c.DBPath)
	if path == "" {
		path = "data/matasidik.db"
	}
	return fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", path)
}

func (c Config) MaxUploadBytes() int64 { return c.MaxUploadMB << 20 }

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func envInt(key string, fallback int64) int64 {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func envFloat(key string, fallback float64) float64 {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return f
}
