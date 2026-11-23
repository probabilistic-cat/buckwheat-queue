package migrations

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed files/*.sql
var migrationFiles embed.FS

func Run(ctx context.Context, dbMigrationURL string, action string) error {
	dbMigration, err := getDB(ctx, dbMigrationURL)
	if err != nil {
		return err
	}
	defer dbMigration.Close()

	source, err := iofs.New(migrationFiles, "files")
	if err != nil {
		return fmt.Errorf("Create migration source: %w", err)
	}

	driver, err := pgx.WithInstance(dbMigration, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("Create migration driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "pgx5", driver)
	if err != nil {
		return fmt.Errorf("Create migrator: %w", err)
	}
	defer migrator.Close()

	err = doMigration(migrator, action)

	return err
}

func getDB(ctx context.Context, dbMigrationURL string) (*sql.DB, error) {
	dbMigration, err := sql.Open("pgx", dbMigrationURL)
	if err != nil {
		return nil, fmt.Errorf("Open database: %w", err)
	}

	if err := dbMigration.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Ping database: %w", err)
	}

	return dbMigration, nil
}

func doMigration(migrator *migrate.Migrate, action string) error {
	var err error = nil

	switch action {
	case "up":
		err = migrator.Up()

	case "down":
		err = migrator.Steps(-1)

	case "version":
		version, dirty, err := migrator.Version()

		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("No migrations applied")
			return nil
		}

		if err != nil {
			return fmt.Errorf("Get migration version: %w", err)
		}

		fmt.Printf("Version: %d, dirty: %v\n", version, dirty)
		return nil

	default:
		return fmt.Errorf("Unknown migration action: %s", action)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("No migrations to apply")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Migration failed: %w", err)
	}

	return nil
}
