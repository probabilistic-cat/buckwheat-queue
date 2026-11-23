package test

import (
	"buckwheat-queue/internal/config"
	"buckwheat-queue/internal/storage/postgres"
	"testing"
	"time"
)

func SetTimezone(t testing.TB) {
	t.Helper()

	_, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("Set timezone: %v", err)
	}
}

func GetConfig(t testing.TB) config.Config {
	t.Helper()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load config: %v", err)
	}

	if cfg.App.Env != config.AppEnvTest {
		t.Fatalf("Wrong env: %s", cfg.App.Env)
	}

	return cfg
}

func GetDB(t testing.TB, cfg config.Config) *postgres.DB {
	t.Helper()

	db, err := postgres.NewDB(t.Context(), cfg.DB, false)
	if err != nil {
		t.Fatalf("Connect to DB: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}
