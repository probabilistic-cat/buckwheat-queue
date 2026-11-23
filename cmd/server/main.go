package main

import (
	"buckwheat-queue/internal/config"
	"buckwheat-queue/internal/log"
	"buckwheat-queue/internal/storage/postgres"
	"buckwheat-queue/internal/storage/postgres/migrations"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := getConfig()

	logger := log.NewLogger(cfg.App.LogLevel)

	err := commands(ctx, cfg)
	if err != nil {
		logger.Error("Main error.", "Error", err)
		os.Exit(1)
	}

	// shutdown
	<-ctx.Done()
	logger.Info("Shutdown signal received")
}

func getConfig() config.Config {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	return cfg
}

func commands(ctx context.Context, cfg config.Config) error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Usage: app [app|migration]")
	}

	switch os.Args[1] {
	case "serve":
		err := app(ctx, cfg)
		if err != nil {
			return err
		}
	case "migration":
		err := migration(ctx, cfg)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown command: %s", os.Args[1])
	}

	return nil
}

func app(ctx context.Context, cfg config.Config) error {
	db, err := postgres.NewDB(ctx, cfg.DB, false)
	if err != nil {
		return fmt.Errorf("Create db pool: %w", err)
	}
	defer db.Close()

	// taskRepo := postgres.NewTaskRepository(db)

	return nil
}

func migration(ctx context.Context, cfg config.Config) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: app migration [up|down|version]")
	}

	action := os.Args[2]
	if err := migrations.Run(ctx, cfg.DB.DBMigrationURL, action); err != nil {
		return fmt.Errorf("Run migration with action %s. Error: %w", action, err)
	}

	fmt.Println("Migration command completed")

	return nil
}
